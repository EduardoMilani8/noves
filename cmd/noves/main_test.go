package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func run(t *testing.T, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	t.Setenv("NOVES_FB_DSN", "")
	t.Setenv("NOVES_PG_DSN", "")
	var out, errb bytes.Buffer
	code = execute(context.Background(), args, &out, &errb)
	return code, out.String(), errb.String()
}

func TestVersion(t *testing.T) {
	code, out, _ := run(t, "version")
	if code != 0 || !strings.HasPrefix(out, "noves ") {
		t.Errorf("code=%d out=%q", code, out)
	}
}

func TestDoctorSemConfiguracaoSaiCom2(t *testing.T) {
	code, out, _ := run(t, "doctor")
	if code != 2 || !strings.Contains(out, "nenhum banco configurado") {
		t.Errorf("code=%d out=%q", code, out)
	}
}

func TestDoctorEmIngles(t *testing.T) {
	_, out, _ := run(t, "doctor", "--lang", "en")
	if !strings.Contains(out, "no database configured") {
		t.Errorf("out=%q", out)
	}
}

func TestErroDeUsoSaiCom2(t *testing.T) {
	if code, _, _ := run(t, "comando-que-nao-existe"); code != 2 {
		t.Errorf("code=%d", code)
	}
	if code, _, errOut := run(t, "doctor", "--lang", "fr"); code != 2 || !strings.Contains(errOut, "idioma") {
		t.Errorf("code=%d err=%q", code, errOut)
	}
}

func TestDoctorConexaoRecusadaSaiCom2(t *testing.T) {
	code, out, _ := run(t, "doctor", "--pg", "postgres://u:segredo@127.0.0.1:1/x?connect_timeout=2", "--timeout", "5s")
	if code != 2 {
		t.Errorf("code=%d out=%q", code, out)
	}
	if strings.Contains(out, "segredo") {
		t.Errorf("senha vazou: %q", out)
	}
}
