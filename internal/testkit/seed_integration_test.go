//go:build integration

package testkit

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

// Confere que make up deixou o PostgreSQL no estado esperado: cópia completa
// e exatamente as divergências plantadas.
func TestSeedAplicouDivergencias(t *testing.T) {
	ctx := context.Background()
	_, pgDSN := DSNs()
	pg, err := pgx.Connect(ctx, pgDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer pg.Close(ctx)

	counts := map[string]int{
		"clientes": 5000, "produtos": 1000, "pedidos": 20000, "itens": 49999,
		"configuracoes": 5, "movimento_estoque": 3000,
	}
	for table, want := range counts {
		var got int
		if err := pg.QueryRow(ctx, "SELECT count(*) FROM "+pgx.Identifier{table}.Sanitize()).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("%s: %d linhas, want %d", table, got, want)
		}
	}

	checks := []struct {
		name, query string
		want        any
	}{
		{"decimal alterado", "SELECT limite::text FROM clientes WHERE id = 42", "57.55"},
		{"codificação dupla", "SELECT nome FROM clientes WHERE id = 13", "JOSÃ‰ SANTOS"},
		{"char sem espaços", "SELECT codigo FROM produtos WHERE id = 77", "P00077"},
		{"'' virou NULL", "SELECT obs IS NULL FROM clientes WHERE id = 300", true},
		{"char com espaços preservado", "SELECT codigo FROM produtos WHERE id = 78", "P00078    "},
		{"'' preservado", "SELECT obs FROM clientes WHERE id = 3", ""},
		{"mojibake de origem preservado", "SELECT nome FROM clientes WHERE id = 997", "JOÃƒO RIBEIRO"},
		{"WIN1252 decodificado", "SELECT nome FROM clientes WHERE id = 1", "JOSÉ FALCÃO"},
	}
	for _, c := range checks {
		var got any
		if err := pg.QueryRow(ctx, c.query).Scan(&got); err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
}
