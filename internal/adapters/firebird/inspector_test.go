package firebird

import "testing"

func TestWithDefaults(t *testing.T) {
	tests := []struct{ in, want string }{
		{"u:p@h/db.fdb", "u:p@h/db.fdb?timezone=UTC"},
		{"u:p@h/db.fdb?charset=WIN1252", "u:p@h/db.fdb?charset=WIN1252&timezone=UTC"},
		{"u:p@h/db.fdb?timezone=America%2FSao_Paulo", "u:p@h/db.fdb?timezone=America%2FSao_Paulo"},
	}
	for _, tt := range tests {
		if got := WithDefaults(tt.in); got != tt.want {
			t.Errorf("WithDefaults(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestQuoteIdent(t *testing.T) {
	if got := QuoteIdent(`A"B`); got != `"A""B"` {
		t.Errorf("got %s", got)
	}
}
