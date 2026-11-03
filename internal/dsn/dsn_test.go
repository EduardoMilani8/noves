package dsn

import (
	"errors"
	"strings"
	"testing"
)

func TestMask(t *testing.T) {
	tests := []struct{ in, want string }{
		{"postgres://noves:segredo@localhost:5433/erp?sslmode=disable", "postgres://noves:***@localhost:5433/erp?sslmode=disable"},
		{"noves:segredo@localhost:3050/var/lib/firebird/data/erp.fdb", "noves:***@localhost:3050/var/lib/firebird/data/erp.fdb"},
		{"sysdba:masterkey@srv/C:/dados/ERP.FDB", "sysdba:***@srv/C:/dados/ERP.FDB"},
		{"user:p@ss@host/db", "user:***@host/db"},
		{"host=localhost user=noves password=segredo dbname=erp", "host=localhost user=noves password=*** dbname=erp"},
		{"host=localhost password='com espaco' dbname=erp", "host=localhost password=*** dbname=erp"},
		{"postgres://noves@localhost/erp?password=segredo&sslmode=disable", "postgres://noves@localhost/erp?password=***&sslmode=disable"},
		{"postgres://noves@localhost/erp", "postgres://noves@localhost/erp"},
		{"localhost:3050/C:/dados/ERP.FDB", "localhost:3050/C:/dados/ERP.FDB"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := Mask(tt.in); got != tt.want {
			t.Errorf("Mask(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestScrub(t *testing.T) {
	d := "noves:s3gr3do@localhost:3050/erp.fdb"
	err := errors.New("falha em " + d + ": senha s3gr3do recusada")
	got := Scrub(err.Error(), d)
	if strings.Contains(got, "s3gr3do") {
		t.Errorf("senha vazou: %q", got)
	}
}
