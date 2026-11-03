package doctor

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/EduardoMilani8/noves/internal/core"
	"github.com/EduardoMilani8/noves/internal/i18n"
	"github.com/EduardoMilani8/noves/internal/ports"
)

type fakeInspector struct {
	info    core.DatabaseInfo
	read    core.ReadCheck
	readErr error
	closed  bool
}

func (f *fakeInspector) Info(context.Context) (core.DatabaseInfo, error) { return f.info, nil }
func (f *fakeInspector) CheckRead(context.Context) (core.ReadCheck, error) {
	return f.read, f.readErr
}
func (f *fakeInspector) Close() error { f.closed = true; return nil }

func opener(f *fakeInspector, err error) Opener {
	return func(context.Context, string) (ports.Inspector, error) {
		if err != nil {
			return nil, err
		}
		return f, nil
	}
}

func fbInfo() core.DatabaseInfo {
	return core.DatabaseInfo{Product: "Firebird", Version: "5.0.4", ODS: "13.1", Dialect: 3,
		Charset: "WIN1252", PageSize: 8192, Database: "/var/lib/firebird/data/erp.fdb", User: "NOVES"}
}

func TestTudoCerto(t *testing.T) {
	fb := &fakeInspector{info: fbInfo(), read: core.ReadCheck{Tables: 8}}
	res := Run(context.Background(), []Target{
		{Name: "Firebird", DSN: "noves:segredo@localhost/erp.fdb", Open: opener(fb, nil)},
		{Name: "PostgreSQL", EnvVar: "NOVES_PG_DSN", Flag: "--pg"},
	})
	if code := ExitCode(res); code != 0 {
		t.Errorf("ExitCode = %d, want 0", code)
	}
	if !fb.closed {
		t.Error("a conexão não foi fechada")
	}

	var buf bytes.Buffer
	Render(&buf, i18n.New(i18n.PT), res)
	out := buf.String()
	for _, want := range []string{"noves:***@localhost/erp.fdb", "✓ conexão ok", "5.0.4", "13.1", "WIN1252",
		"leitura ok em 8 tabelas", "PostgreSQL: não configurado (defina NOVES_PG_DSN ou use --pg)", "tudo certo"} {
		if !strings.Contains(out, want) {
			t.Errorf("saída sem %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "segredo") {
		t.Errorf("senha vazou na saída:\n%s", out)
	}
}

func TestFalhaDeConexaoSaiCom2ESemSenha(t *testing.T) {
	d := "noves:segredo@localhost/erp.fdb"
	res := Run(context.Background(), []Target{
		{Name: "Firebird", DSN: d, Open: opener(nil, errors.New("login falhou para "+d))},
	})
	if code := ExitCode(res); code != 2 {
		t.Errorf("ExitCode = %d, want 2", code)
	}
	var buf bytes.Buffer
	Render(&buf, i18n.New(i18n.EN), res)
	if out := buf.String(); strings.Contains(out, "segredo") || !strings.Contains(out, "connection failed") {
		t.Errorf("saída inesperada:\n%s", out)
	}
}

func TestSemPermissaoSaiCom1(t *testing.T) {
	fb := &fakeInspector{info: fbInfo(), read: core.ReadCheck{Tables: 8, Unreadable: []string{"CLIENTES"}}}
	res := Run(context.Background(), []Target{{Name: "Firebird", DSN: "x:y@h/db", Open: opener(fb, nil)}})
	if code := ExitCode(res); code != 1 {
		t.Errorf("ExitCode = %d, want 1", code)
	}
	var buf bytes.Buffer
	Render(&buf, i18n.New(i18n.PT), res)
	if !strings.Contains(buf.String(), "sem permissão de leitura em 1 de 8 tabelas: CLIENTES") {
		t.Errorf("saída:\n%s", buf.String())
	}
}

func TestNadaConfiguradoSaiCom2(t *testing.T) {
	res := Run(context.Background(), []Target{{Name: "Firebird"}, {Name: "PostgreSQL"}})
	if code := ExitCode(res); code != 2 {
		t.Errorf("ExitCode = %d, want 2", code)
	}
}
