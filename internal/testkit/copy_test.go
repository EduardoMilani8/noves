package testkit

import (
	"bytes"
	"encoding/csv"
	"math/big"
	"strings"
	"testing"
	"time"
)

func TestFormatValue(t *testing.T) {
	ts := time.Date(2015, 1, 1, 8, 19, 43, 680_000_000, time.UTC)
	tests := []struct {
		name   string
		v      any
		pgType string
		want   string
	}{
		{"null", nil, "integer", `\N`},
		{"string vazia não é null", []byte{}, "character varying(200)", ""},
		{"win1252 vira utf-8", []byte{'J', 'O', 0xC3, 'O'}, "character varying(60)", "JOÃO"},
		{"mojibake é preservado", []byte{'J', 'O', 0xC3, 0x83, 'O'}, "character varying(60)", "JOÃƒO"},
		{"char mantém espaços", []byte("P00077    "), "character varying(10)", "P00077    "},
		{"bytea em hex", []byte{0x00, 0xFF}, "bytea", `\x00ff`},
		{"numeric exato", "57.54", "numeric(15,2)", "57.54"},
		{"blob de texto chega como string crua", string([]byte{'A', 0xC7, 0xDA, 'C', 'A', 'R'}), "text", "AÇÚCAR"},
		{"inteiro", int64(-42), "integer", "-42"},
		{"double sem perda", 1.0 / 3, "double precision", "0.3333333333333333"},
		{"int128", big.NewInt(170), "numeric(38,0)", "170"},
		{"boolean", true, "boolean", "true"},
		{"date", ts, "date", "2015-01-01"},
		{"time", ts, "time without time zone", "08:19:43.68"},
		{"timestamp", ts, "timestamp without time zone", "2015-01-01 08:19:43.68"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatValue(tt.v, tt.pgType)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("FormatValue(%v, %q) = %q, want %q", tt.v, tt.pgType, got, tt.want)
			}
		})
	}
}

func TestFormatValueTipoDesconhecido(t *testing.T) {
	if _, err := FormatValue(struct{}{}, "text"); err == nil {
		t.Fatal("esperava erro para tipo desconhecido")
	}
}

// O COPY só distingue string vazia de NULL se o CSV não puser aspas no marcador e
// mantiver o campo vazio como vazio.
func TestCSVDistingueVazioDeNull(t *testing.T) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write([]string{"", nullMarker, "a\r\nb"}); err != nil {
		t.Fatal(err)
	}
	w.Flush()
	got := buf.String()
	if !strings.HasPrefix(got, `,\N,"a`) {
		t.Fatalf("CSV inesperado: %q", got)
	}
}
