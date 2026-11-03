// Package dsn mascara senhas em strings de conexão antes de exibi-las.
package dsn

import (
	"regexp"
	"strings"
)

const mask = "***"

// keyword cobre password=... no formato chave=valor do libpq e em query strings.
var keyword = regexp.MustCompile(`(?i)(\bpassword\s*=\s*)('(?:[^'\\]|\\.)*'|[^\s&]+)`)

// Mask troca a senha do DSN por ***. Entende URLs (postgres://user:senha@host),
// o formato do driver Firebird (user:senha@host/caminho) e chave=valor
// (password=senha).
func Mask(s string) string {
	s = keyword.ReplaceAllString(s, "${1}"+mask)
	at, colon := userinfo(s)
	if colon < 0 {
		return s
	}
	return s[:colon+1] + mask + s[at:]
}

// Scrub remove de um texto qualquer (por exemplo, a mensagem de erro de um
// driver) a senha contida no DSN.
func Scrub(text, dsn string) string {
	if p := password(dsn); p != "" {
		text = strings.ReplaceAll(text, p, mask)
	}
	return strings.ReplaceAll(text, dsn, Mask(dsn))
}

func password(s string) string {
	if m := keyword.FindStringSubmatch(s); m != nil {
		return strings.Trim(m[2], "'")
	}
	at, colon := userinfo(s)
	if colon < 0 {
		return ""
	}
	return s[colon+1 : at]
}

// userinfo localiza "usuario:senha@" e devolve os índices do @ e dos
// dois-pontos (ou -1 se não houver senha).
func userinfo(s string) (at, colon int) {
	start := 0
	if i := strings.Index(s, "://"); i >= 0 {
		start = i + 3
	}
	rest := s[start:]
	if slash := strings.Index(rest, "/"); slash >= 0 {
		// A senha pode conter "@" sem escape; a barra do caminho não.
		if a := strings.LastIndex(rest[:slash], "@"); a >= 0 {
			rest = rest[:a+1]
		}
	}
	a := strings.LastIndex(rest, "@")
	if a < 0 || strings.ContainsAny(rest[:a], " =") {
		return -1, -1
	}
	c := strings.Index(rest[:a], ":")
	if c < 0 {
		return -1, -1
	}
	return start + a, start + c
}
