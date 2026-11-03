package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/EduardoMilani8/noves/internal/core"
	"github.com/EduardoMilani8/noves/internal/ports"
)

var _ ports.Inspector = (*Inspector)(nil)

// Inspector implementa ports.Inspector para o PostgreSQL.
type Inspector struct {
	conn *pgx.Conn
}

// Open conecta ao PostgreSQL.
func Open(ctx context.Context, dsn string) (*Inspector, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Inspector{conn: conn}, nil
}

// Close encerra a conexão.
func (i *Inspector) Close() error { return i.conn.Close(context.Background()) }

// Info lê versão, banco, usuário e encoding.
func (i *Inspector) Info(ctx context.Context) (core.DatabaseInfo, error) {
	info := core.DatabaseInfo{Product: "PostgreSQL"}
	err := i.conn.QueryRow(ctx, `
		SELECT current_setting('server_version'), current_database(), current_user,
		       pg_encoding_to_char(d.encoding), current_setting('block_size')::int
		  FROM pg_database d
		 WHERE d.datname = current_database()`).
		Scan(&info.Version, &info.Database, &info.User, &info.Charset, &info.PageSize)
	if err != nil {
		return info, fmt.Errorf("lendo versão do PostgreSQL: %w", err)
	}
	return info, nil
}

// CheckRead usa has_table_privilege, que já resolve papéis herdados.
func (i *Inspector) CheckRead(ctx context.Context) (core.ReadCheck, error) {
	var res core.ReadCheck
	rows, err := i.conn.Query(ctx, `
		SELECT n.nspname || '.' || c.relname, has_table_privilege(c.oid, 'SELECT')
		  FROM pg_class c
		  JOIN pg_namespace n ON n.oid = c.relnamespace
		 WHERE c.relkind IN ('r', 'p')
		   AND n.nspname NOT IN ('pg_catalog', 'information_schema')
		   AND n.nspname NOT LIKE 'pg\_toast%'
		 ORDER BY 1`)
	if err != nil {
		return res, fmt.Errorf("listando tabelas: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var ok bool
		if err := rows.Scan(&name, &ok); err != nil {
			return res, err
		}
		res.Tables++
		if !ok {
			res.Unreadable = append(res.Unreadable, name)
		}
	}
	return res, rows.Err()
}
