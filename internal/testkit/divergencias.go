package testkit

// Divergencia descreve uma diferença plantada de propósito no PostgreSQL por
// deploy/postgres/divergencias.sql.
type Divergencia struct {
	Tabela string
	Chave  string // PK em texto; composta separada por vírgula
	Coluna string // vazia quando a linha inteira falta
	Tipo   string
}

// Tipos de divergência plantados.
const (
	LinhaFaltando    = "linha-faltando"
	ValorAlterado    = "valor-alterado"
	CodificacaoDupla = "codificacao-dupla"
	CharSemEspacos   = "char-sem-espacos"
	VazioVirouNull   = "vazio-virou-null"
)

// ArquivoDivergente é o script, relativo à raiz do repositório, que planta as divergências.
const ArquivoDivergente = "deploy/postgres/divergencias.sql"

// Divergencias são exatamente as diferenças que os testes de integração
// precisam encontrar, e nenhuma outra.
var Divergencias = []Divergencia{
	{Tabela: "itens", Chave: "1234,2", Tipo: LinhaFaltando},
	{Tabela: "clientes", Chave: "42", Coluna: "limite", Tipo: ValorAlterado},
	{Tabela: "clientes", Chave: "13", Coluna: "nome", Tipo: CodificacaoDupla},
	{Tabela: "produtos", Chave: "77", Coluna: "codigo", Tipo: CharSemEspacos},
	{Tabela: "clientes", Chave: "300", Coluna: "obs", Tipo: VazioVirouNull},
}
