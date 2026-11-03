# CLAUDE.md — convenções do noves

Leia `docs/PLANO.md` (plano completo) e `docs/PROGRESSO.md` (onde paramos) antes de trabalhar.

## O projeto
CLI em Go que **não converte** bancos: avalia o risco de migrar Firebird → PostgreSQL e prova, depois, que dados e comportamento batem. Lema: "Converta com o que quiser. Prove com o noves."

## Como trabalhar com o Eduardo
- Ele é estudante de Ciência da Computação e quer aprender: em cada decisão técnica importante, explique o porquê em 1 a 3 frases.
- Fase por fase. No fim de cada uma: resumo curto + como testar em até 3 comandos, e espere "ok" (a menos que ele diga "pode seguir direto").
- Se algo do plano estiver errado ou for inviável, pare, explique e proponha alternativa. Não contorne em silêncio.

## Commits
- Conventional Commits: tipo em inglês, descrição em português no imperativo. Ex.: `feat(assess): detecta SUSPEND em procedures`.
- Um commit por mudança lógica; o teste vai junto com o código que testa.
- Nunca misture refatoração com funcionalidade; nunca commite código que não compila.
- Antes de cada commit: `gofmt -l .`, `go vet ./...`, `golangci-lint run`, `go test ./...` (atalho: `make check`).

## GitHub
- Uma branch por fase (`fase-N-nome`), PR com `gh pr create` (o que mudou + como testar, `Closes #N`).
- Esperar CI verde e fazer merge com `gh pr merge --rebase --delete-branch`.
- Nunca force push na `main`.
- Milestone por versão; issues por tarefa da fase.
- Ao fechar a fase: CHANGELOG (Keep a Changelog), tag `vX.Y.Z`, conferir a release do GoReleaser.

## Código
- Go estável instalado; binário com `CGO_ENABLED=0`.
- Arquitetura hexagonal: `internal/core` não importa driver de banco; portas em `internal/ports`; adaptadores em `internal/adapters/{firebird,postgres}`.
- Não invente APIs: confira com `go doc` ou no código do módulo.
- Detecção em SQL/PSQL só por tokens (`internal/psql`), nunca regex no texto cru.
- Dependências mínimas e com licença permissiva (nada de LGPL/GPL, ex.: ACBr).
- Mensagens em português por padrão; `--lang en` para inglês.
- Saída: 0 = ok; 1 = achados/divergências acima de `--fail-on`; 2 = erro de uso ou conexão.

## Segurança
- Nunca versionar credenciais (`.env` no `.gitignore`, `.env.example` versionado).
- Mascarar senhas de DSN em logs e mensagens (`internal/dsn`).
- Todos os comandos são somente leitura, exceto `sync install/run/uninstall`, que mostram o DDL (`--dry-run` é o padrão) e só aplicam com `--yes`.
- Amostras de dados mascaradas por padrão (LGPD); `--show-samples` exibe.

## Ambiente de teste
- `make up` sobe Firebird 5 + PostgreSQL (`deploy/docker-compose.yml`), copia os dados e injeta as divergências; `make down` derruba.
- Testes de integração: build tag `integration`, `make test-integration`.
- As divergências injetadas estão em `deploy/postgres/divergencias.sql`; os testes exigem encontrar exatamente elas.
