# Contribuindo com o noves

Obrigado pelo interesse! Issues e PRs são bem-vindos, em português ou inglês.

## Preparando o ambiente

Requisitos: Go (versão em `go.mod`), Docker com Compose e
[golangci-lint](https://golangci-lint.run/).

```bash
make build             # compila em bin/noves
make test              # testes de unidade
make up                # sobe Firebird + PostgreSQL de teste
make test-integration  # testes de integração (exige make up)
make down              # derruba o ambiente
```

## Antes de abrir um PR

Rode `make check` (gofmt, go vet, golangci-lint e testes). O CI roda o mesmo.

## Commits

Usamos [Conventional Commits](https://www.conventionalcommits.org/pt-br/), com o
tipo em inglês e a descrição em português, no imperativo:

```
feat(assess): detecta SUSPEND em procedures
fix(data): trunca TIMESTAMP para a precisão comum
```

- Um commit por mudança lógica, com o teste junto do código que ele testa.
- Não misture refatoração com funcionalidade.
- Todo commit precisa compilar.

## Arquitetura

O núcleo (`internal/core`) não conhece banco nenhum. Ele fala com o mundo por
interfaces em `internal/ports`, implementadas pelos adaptadores em
`internal/adapters`. Para suportar uma nova origem (Oracle, SQL Server...), basta
um novo adaptador.

## Como adicionar uma regra

Guia completo a partir da v0.1.0, quando o catálogo de regras existir.

## Licenças

Só aceitamos dependências com licença permissiva (MIT, BSD, Apache-2.0...). Não
copie código de projetos com licença incompatível, como LGPL ou GPL.
