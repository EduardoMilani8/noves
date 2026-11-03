package testkit

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/text/encoding/charmap"
)

// nullMarker representa NULL no COPY em CSV. Com ele, um campo vazio sem aspas
// continua sendo string vazia, e ” não vira NULL por acidente na cópia.
const nullMarker = `\N`

// Column é uma coluna do destino que recebe dados da origem.
type Column struct {
	Name   string // nome no PostgreSQL
	Source string // nome exato no Firebird
	PGType string // format_type() do PostgreSQL
}

// CopyTable copia uma tabela do Firebird para a tabela de mesmo nome no
// PostgreSQL, em streaming, via COPY FROM STDIN em CSV.
//
// A conexão Firebird precisa usar charset=None: assim o driver devolve os
// bytes crus (com os espaços de CHAR) e o texto é decodificado aqui como
// WIN1252, que é o que o legado grava até em colunas NONE.
func CopyTable(ctx context.Context, fb *sql.DB, pg *pgx.Conn, table string) (int64, error) {
	cols, err := matchColumns(ctx, fb, pg, table)
	if err != nil {
		return 0, err
	}

	src := make([]string, len(cols))
	dst := make([]string, len(cols))
	for i, c := range cols {
		src[i] = `"` + c.Source + `"`
		dst[i] = pgx.Identifier{c.Name}.Sanitize()
	}

	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(src, ", "), strings.ToUpper(table))
	pr, pw := io.Pipe()
	go func() {
		pw.CloseWithError(writeCSV(ctx, pw, fb, query, cols))
	}()

	copySQL := fmt.Sprintf(`COPY %s (%s) FROM STDIN WITH (FORMAT csv, NULL '%s')`,
		pgx.Identifier{table}.Sanitize(), strings.Join(dst, ", "), nullMarker)
	tag, err := pg.PgConn().CopyFrom(ctx, pr, copySQL)
	if err != nil {
		// Destrava o goroutine se o COPY falhou antes de ler tudo.
		pr.CloseWithError(err)
		return 0, fmt.Errorf("copiando %s para o PostgreSQL: %w", table, err)
	}
	return tag.RowsAffected(), nil
}

func writeCSV(ctx context.Context, w io.Writer, fb *sql.DB, query string, cols []Column) error {
	rows, err := fb.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("lendo do Firebird: %w", err)
	}
	defer rows.Close()

	cw := csv.NewWriter(w)
	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	rec := make([]string, len(cols))
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return err
		}
		for i, v := range vals {
			s, err := FormatValue(v, cols[i].PGType)
			if err != nil {
				return fmt.Errorf("coluna %s: %w", cols[i].Name, err)
			}
			rec[i] = s
		}
		if err := cw.Write(rec); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	cw.Flush()
	return cw.Error()
}

// FormatValue converte um valor lido do Firebird (driver nakagami, conexão
// charset=None) para o texto que o COPY em CSV do PostgreSQL entende.
func FormatValue(v any, pgType string) (string, error) {
	switch x := v.(type) {
	case nil:
		return nullMarker, nil
	case []byte:
		if pgType == "bytea" {
			return `\x` + hex.EncodeToString(x), nil
		}
		return decodeWin1252(x)
	case string:
		// NUMERIC/DECIMAL chegam como texto ASCII exato; BLOB de texto chega
		// como string com os bytes crus. Decodificar ASCII não muda nada.
		return decodeWin1252([]byte(x))
	case int64:
		return strconv.FormatInt(x, 10), nil
	case float64:
		// 'g' com precisão -1 é a menor representação que volta ao mesmo float.
		return strconv.FormatFloat(x, 'g', -1, 64), nil
	case bool:
		return strconv.FormatBool(x), nil
	case *big.Int:
		return x.String(), nil
	case time.Time:
		switch {
		case pgType == "date":
			return x.Format("2006-01-02"), nil
		case strings.HasPrefix(pgType, "timestamp"):
			return x.Format("2006-01-02 15:04:05.999999"), nil
		default: // time
			return x.Format("15:04:05.999999"), nil
		}
	default:
		return "", fmt.Errorf("tipo não suportado: %T", v)
	}
}

func decodeWin1252(b []byte) (string, error) {
	s, err := charmap.Windows1252.NewDecoder().Bytes(b)
	return string(s), err
}

// matchColumns casa as colunas do destino (sem as geradas) com as da origem,
// ignorando maiúsculas e minúsculas.
func matchColumns(ctx context.Context, fb *sql.DB, pg *pgx.Conn, table string) ([]Column, error) {
	src := map[string]string{}
	rows, err := fb.QueryContext(ctx,
		`SELECT TRIM(RDB$FIELD_NAME) FROM RDB$RELATION_FIELDS WHERE RDB$RELATION_NAME = ?`,
		strings.ToUpper(table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name any
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		s := asString(name)
		src[strings.ToLower(s)] = s
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	pgRows, err := pg.Query(ctx, `
		SELECT a.attname, format_type(a.atttypid, a.atttypmod)
		  FROM pg_attribute a
		 WHERE a.attrelid = $1::regclass AND a.attnum > 0
		   AND NOT a.attisdropped AND a.attgenerated = ''
		 ORDER BY a.attnum`, table)
	if err != nil {
		return nil, err
	}
	defer pgRows.Close()
	var cols []Column
	for pgRows.Next() {
		var c Column
		if err := pgRows.Scan(&c.Name, &c.PGType); err != nil {
			return nil, err
		}
		s, ok := src[strings.ToLower(c.Name)]
		if !ok {
			return nil, fmt.Errorf("coluna %s.%s não existe no Firebird", table, c.Name)
		}
		c.Source = s
		cols = append(cols, c)
	}
	return cols, pgRows.Err()
}

func asString(v any) string {
	if b, ok := v.([]byte); ok {
		return string(b)
	}
	return fmt.Sprint(v)
}
