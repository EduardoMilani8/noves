package testkit

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"

	// Registra o driver "firebirdsql".
	_ "github.com/nakagami/firebirdsql"
)

// DSNs padrão do ambiente de teste (deploy/docker-compose.yml).
const (
	DefaultFirebirdDSN = "noves:noves@localhost:3050/var/lib/firebird/data/erp.fdb"
	DefaultPostgresDSN = "postgres://noves:noves@localhost:5433/erp?sslmode=disable"
)

// Tables são as tabelas copiadas, na ordem das chaves estrangeiras.
var Tables = []string{
	"clientes", "produtos", "pedidos", "itens",
	"configuracoes", "movimento_estoque", "log_auditoria",
}

// DSNs devolve os DSNs do ambiente de teste, lendo NOVES_FB_DSN e NOVES_PG_DSN
// quando definidos.
func DSNs() (fb, pg string) {
	fb, pg = os.Getenv("NOVES_FB_DSN"), os.Getenv("NOVES_PG_DSN")
	if fb == "" {
		fb = DefaultFirebirdDSN
	}
	if pg == "" {
		pg = DefaultPostgresDSN
	}
	return fb, pg
}

// RawFirebirdDSN acrescenta ao DSN as opções que a cópia exige: charset=None,
// para receber bytes crus, e timezone=UTC, para que horários que não existem no
// fuso local (início do horário de verão) não sejam "corrigidos" pelo Go.
func RawFirebirdDSN(dsn string) string {
	sep := "?"
	if strings.Contains(dsn, "?") {
		sep = "&"
	}
	q := url.Values{"charset": {"None"}, "timezone": {"UTC"}}
	return dsn + sep + q.Encode()
}

// Seed apaga os dados do PostgreSQL de teste, copia tudo do Firebird e aplica
// as divergências de divergenciasSQL (conteúdo de deploy/postgres/divergencias.sql).
func Seed(ctx context.Context, fbDSN, pgDSN, divergenciasSQL string, logf func(string, ...any)) error {
	fb, err := sql.Open("firebirdsql", RawFirebirdDSN(fbDSN))
	if err != nil {
		return err
	}
	defer fb.Close()
	if err := fb.PingContext(ctx); err != nil {
		return fmt.Errorf("conectando ao Firebird: %w", err)
	}

	pg, err := pgx.Connect(ctx, pgDSN)
	if err != nil {
		return fmt.Errorf("conectando ao PostgreSQL: %w", err)
	}
	defer pg.Close(ctx)

	quoted := make([]string, len(Tables))
	for i, t := range Tables {
		quoted[i] = pgx.Identifier{t}.Sanitize()
	}
	if _, err := pg.Exec(ctx, "TRUNCATE "+strings.Join(quoted, ", ")+" CASCADE"); err != nil {
		return err
	}

	for _, t := range Tables {
		n, err := CopyTable(ctx, fb, pg, t)
		if err != nil {
			return err
		}
		logf("%-18s %7d linhas", t, n)
	}

	if divergenciasSQL != "" {
		if _, err := pg.Exec(ctx, divergenciasSQL); err != nil {
			return fmt.Errorf("aplicando divergências: %w", err)
		}
		logf("divergências aplicadas: %d", len(Divergencias))
	}
	_, err = pg.Exec(ctx, "ANALYZE")
	return err
}
