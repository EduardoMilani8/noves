//go:build integration

package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/EduardoMilani8/noves/internal/testkit"
)

func TestDoctorNoAmbienteDeTeste(t *testing.T) {
	fb, pg := testkit.DSNs()
	var out, errb bytes.Buffer
	code := execute(context.Background(), []string{"doctor", "--fb", fb, "--pg", pg}, &out, &errb)
	if code != 0 {
		t.Fatalf("code=%d\nstdout:\n%s\nstderr:\n%s", code, out.String(), errb.String())
	}
	for _, want := range []string{"WIN1252", "dialeto", "ODS", "UTF8", "tudo certo"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("saída sem %q:\n%s", want, out.String())
		}
	}
	if strings.Contains(out.String(), ":noves@") {
		t.Errorf("senha vazou:\n%s", out.String())
	}
}
