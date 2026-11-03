package core

// DatabaseInfo descreve o servidor e o banco a que o noves se conectou.
// Campos que não se aplicam a um produto ficam vazios (ODS e Dialect só
// existem no Firebird).
type DatabaseInfo struct {
	Product  string // "Firebird" ou "PostgreSQL"
	Version  string // ex.: "5.0.4", "18.1"
	Database string // caminho do arquivo ou nome do banco
	User     string
	Charset  string // NONE quando o banco Firebird não define charset
	ODS      string // ex.: "13.1"
	Dialect  int
	PageSize int
}

// ReadCheck é o resultado da conferência de permissão de leitura.
type ReadCheck struct {
	Tables     int      // tabelas de usuário encontradas
	Unreadable []string // tabelas sem permissão de SELECT
}

// OK indica que todas as tabelas podem ser lidas.
func (r ReadCheck) OK() bool { return len(r.Unreadable) == 0 }
