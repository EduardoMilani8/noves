package buildinfo

import (
	"strings"
	"testing"
)

func TestStringUsaLdflags(t *testing.T) {
	old := Version
	t.Cleanup(func() { Version = old })
	Version = "v1.2.3"

	if got := String(); !strings.HasPrefix(got, "noves v1.2.3 (commit ") {
		t.Errorf("String() = %q", got)
	}
}
