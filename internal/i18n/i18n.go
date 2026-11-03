// Package i18n traduz as mensagens do noves. Português é o padrão; --lang en
// troca para inglês.
package i18n

import (
	"fmt"
	"strings"
)

// Lang é um idioma suportado.
type Lang string

// Idiomas suportados.
const (
	PT Lang = "pt"
	EN Lang = "en"
)

// Parse interpreta valores como "pt", "pt-BR", "en" ou "en_US".
func Parse(s string) (Lang, error) {
	switch l := strings.ToLower(s); {
	case l == "" || strings.HasPrefix(l, "pt"):
		return PT, nil
	case strings.HasPrefix(l, "en"):
		return EN, nil
	default:
		return "", fmt.Errorf("idioma não suportado: %q (use pt ou en)", s)
	}
}

// Messages devolve textos traduzidos.
type Messages struct{ lang Lang }

// New cria um tradutor para o idioma.
func New(l Lang) Messages { return Messages{lang: l} }

// T traduz a chave, formatando os argumentos como fmt.Sprintf. Chave
// desconhecida volta como está, para ficar visível em vez de sumir.
func (m Messages) T(key string, args ...any) string {
	s, ok := catalog[m.lang][key]
	if !ok {
		s, ok = catalog[PT][key]
	}
	if !ok {
		s = key
	}
	if len(args) > 0 {
		return fmt.Sprintf(s, args...)
	}
	return s
}

var catalog = map[Lang]map[string]string{
	PT: {
		"doctor.connected":      "conexão ok",
		"doctor.failed":         "falha na conexão",
		"doctor.skipped":        "não configurado (defina %s ou use %s)",
		"doctor.version":        "versão",
		"doctor.ods":            "ODS",
		"doctor.dialect":        "dialeto",
		"doctor.charset":        "charset",
		"doctor.page_size":      "página",
		"doctor.database":       "banco",
		"doctor.user":           "usuário",
		"doctor.read_ok":        "leitura ok em %d tabelas",
		"doctor.read_denied":    "sem permissão de leitura em %d de %d tabelas: %s",
		"doctor.read_error":     "não foi possível conferir a leitura",
		"doctor.nothing":        "nenhum banco configurado: defina NOVES_FB_DSN e/ou NOVES_PG_DSN (ou use --fb/--pg)",
		"doctor.summary_ok":     "tudo certo",
		"doctor.summary_issues": "%d problema(s) encontrado(s)",
	},
	EN: {
		"doctor.connected":      "connection ok",
		"doctor.failed":         "connection failed",
		"doctor.skipped":        "not configured (set %s or use %s)",
		"doctor.version":        "version",
		"doctor.ods":            "ODS",
		"doctor.dialect":        "dialect",
		"doctor.charset":        "charset",
		"doctor.page_size":      "page size",
		"doctor.database":       "database",
		"doctor.user":           "user",
		"doctor.read_ok":        "read access ok on %d tables",
		"doctor.read_denied":    "no read permission on %d of %d tables: %s",
		"doctor.read_error":     "could not check read access",
		"doctor.nothing":        "no database configured: set NOVES_FB_DSN and/or NOVES_PG_DSN (or use --fb/--pg)",
		"doctor.summary_ok":     "all good",
		"doctor.summary_issues": "%d issue(s) found",
	},
}
