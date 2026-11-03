package firebird

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	// Registra o driver "firebirdsql" (puro Go, sem fbclient).
	_ "github.com/nakagami/firebirdsql"

	"github.com/EduardoMilani8/noves/internal/core"
	"github.com/EduardoMilani8/noves/internal/ports"
)

var _ ports.Inspector = (*Inspector)(nil)

// Inspector implementa ports.Inspector para o Firebird, lendo as tabelas de
// sistema RDB$ e MON$.
type Inspector struct {
	db *sql.DB
}

// Open conecta ao Firebird e confere a conexão.
func Open(ctx context.Context, dsn string) (*Inspector, error) {
	db, err := sql.Open("firebirdsql", WithDefaults(dsn))
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return &Inspector{db: db}, nil
}

// WithDefaults acrescenta ao DSN as opções que o noves espera quando o
// usuário não as definiu. timezone=UTC evita que o driver monte datas no fuso
// local, onde horários do início do horário de verão não existem.
func WithDefaults(dsn string) string {
	base, query, _ := strings.Cut(dsn, "?")
	q, err := url.ParseQuery(query)
	if err != nil {
		return dsn
	}
	if q.Get("timezone") != "" {
		return dsn
	}
	q.Set("timezone", "UTC")
	return base + "?" + q.Encode()
}

// Close encerra a conexão.
func (i *Inspector) Close() error { return i.db.Close() }

// Info lê versão do motor, ODS, dialeto, tamanho de página e charset.
func (i *Inspector) Info(ctx context.Context) (core.DatabaseInfo, error) {
	info := core.DatabaseInfo{Product: "Firebird"}

	var charset sql.NullString
	err := i.db.QueryRowContext(ctx, `
		SELECT RDB$GET_CONTEXT('SYSTEM', 'ENGINE_VERSION'), CURRENT_USER, TRIM(RDB$CHARACTER_SET_NAME)
		  FROM RDB$DATABASE`).Scan(&info.Version, &info.User, &charset)
	if err != nil {
		return info, fmt.Errorf("lendo RDB$DATABASE: %w", err)
	}
	info.User = strings.TrimSpace(info.User)
	info.Charset = "NONE"
	if charset.Valid && charset.String != "" {
		info.Charset = charset.String
	}

	var odsMajor, odsMinor int
	err = i.db.QueryRowContext(ctx, `
		SELECT MON$ODS_MAJOR, MON$ODS_MINOR, MON$SQL_DIALECT, MON$PAGE_SIZE, MON$DATABASE_NAME
		  FROM MON$DATABASE`).Scan(&odsMajor, &odsMinor, &info.Dialect, &info.PageSize, &info.Database)
	if err != nil {
		return info, fmt.Errorf("lendo MON$DATABASE: %w", err)
	}
	info.ODS = fmt.Sprintf("%d.%d", odsMajor, odsMinor)
	info.Database = strings.TrimSpace(info.Database)
	return info, nil
}

// CheckRead tenta um SELECT vazio em cada tabela de usuário. É mais fiel do
// que consultar RDB$USER_PRIVILEGES, que não resolve sozinho papéis, PUBLIC,
// dono do banco e RDB$ADMIN.
func (i *Inspector) CheckRead(ctx context.Context) (core.ReadCheck, error) {
	var res core.ReadCheck
	tables, err := i.userTables(ctx)
	if err != nil {
		return res, err
	}
	res.Tables = len(tables)
	for _, t := range tables {
		if !i.canRead(ctx, t) {
			res.Unreadable = append(res.Unreadable, t)
		}
	}
	return res, nil
}

func (i *Inspector) canRead(ctx context.Context, table string) bool {
	rows, err := i.db.QueryContext(ctx, "SELECT 1 FROM "+QuoteIdent(table)+" WHERE 1 = 0")
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Err() == nil
}

func (i *Inspector) userTables(ctx context.Context) ([]string, error) {
	rows, err := i.db.QueryContext(ctx, `
		SELECT TRIM(RDB$RELATION_NAME)
		  FROM RDB$RELATIONS
		 WHERE COALESCE(RDB$SYSTEM_FLAG, 0) = 0
		   AND RDB$VIEW_BLR IS NULL
		 ORDER BY 1`)
	if err != nil {
		return nil, fmt.Errorf("listando tabelas: %w", err)
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

// QuoteIdent põe o nome entre aspas duplas, como o Firebird espera para nomes
// com caixa mista ou caracteres especiais.
func QuoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}
