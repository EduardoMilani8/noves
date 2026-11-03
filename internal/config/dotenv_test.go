package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseDotEnv(t *testing.T) {
	in := `
# comentário
NOVES_FB_DSN=noves:noves@localhost:3050/erp.fdb
export NOVES_PG_DSN="postgres://u:p@h/db?sslmode=disable"
NOVES_LANG='en'
linha inválida
`
	got, err := ParseDotEnv(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"NOVES_FB_DSN": "noves:noves@localhost:3050/erp.fdb",
		"NOVES_PG_DSN": "postgres://u:p@h/db?sslmode=disable",
		"NOVES_LANG":   "en",
	}
	if len(got) != len(want) {
		t.Fatalf("got %v", got)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
}

func TestLoadDotEnvNaoSobrescreve(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("NOVES_TESTE_A=arquivo\nNOVES_TESTE_B=arquivo\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("NOVES_TESTE_A", "ambiente")
	t.Setenv("NOVES_TESTE_B", "")
	os.Unsetenv("NOVES_TESTE_B")

	if err := LoadDotEnv(path); err != nil {
		t.Fatal(err)
	}
	if v := os.Getenv("NOVES_TESTE_A"); v != "ambiente" {
		t.Errorf("A = %q, o ambiente deveria vencer", v)
	}
	if v := os.Getenv("NOVES_TESTE_B"); v != "arquivo" {
		t.Errorf("B = %q", v)
	}
	if err := LoadDotEnv(filepath.Join(dir, "nao-existe")); err != nil {
		t.Errorf("arquivo ausente não deveria ser erro: %v", err)
	}
}
