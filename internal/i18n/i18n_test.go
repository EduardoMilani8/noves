package i18n

import "testing"

// Toda mensagem em português precisa existir em inglês, e vice-versa.
func TestCatalogosCompletos(t *testing.T) {
	for k := range catalog[PT] {
		if _, ok := catalog[EN][k]; !ok {
			t.Errorf("chave %q sem tradução em inglês", k)
		}
	}
	for k := range catalog[EN] {
		if _, ok := catalog[PT][k]; !ok {
			t.Errorf("chave %q sem tradução em português", k)
		}
	}
}

func TestParse(t *testing.T) {
	for in, want := range map[string]Lang{"": PT, "pt": PT, "pt-BR": PT, "en": EN, "EN_us": EN} {
		got, err := Parse(in)
		if err != nil || got != want {
			t.Errorf("Parse(%q) = %q, %v", in, got, err)
		}
	}
	if _, err := Parse("fr"); err == nil {
		t.Error("esperava erro para fr")
	}
}

func TestT(t *testing.T) {
	if got := New(EN).T("doctor.read_ok", 3); got != "read access ok on 3 tables" {
		t.Errorf("got %q", got)
	}
	if got := New(PT).T("chave.inexistente"); got != "chave.inexistente" {
		t.Errorf("got %q", got)
	}
}
