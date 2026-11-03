# Projeto: noves — a prova dos nove de migrações Firebird → PostgreSQL

Você vai construir do zero, até a v1.0.0, uma CLI open source em Go chamada **noves**, publicada no meu GitHub, com histórico de commits limpo e separado por mudança lógica. Leia tudo antes de começar.

## 1. Quem sou eu e como trabalhar comigo
- Sou estudante de Ciência da Computação e trabalho com integração de dados (Go, Python, PostgreSQL, Oracle). Uso arquitetura hexagonal em Go.
- Quero aprender com o projeto: ao tomar uma decisão técnica importante, explique em 1 a 3 frases o porquê.
- Trabalhe fase por fase (seção 7). Ao fim de cada fase, me mostre um resumo curto (o que foi feito e como testar em até 3 comandos) e espere meu "ok" antes de seguir, a menos que eu diga "pode seguir direto".

## 2. Regras de trabalho (sempre)
1. Antes de tudo, confira o ambiente: `go version`, `docker info`, `docker compose version`, `git config user.name`, `gh auth status`. Se faltar algo, pare e me diga exatamente o que instalar ou configurar.
2. Logo no começo, salve este documento inteiro em `docs/PLANO.md` e crie um `CLAUDE.md` na raiz com as convenções deste prompt, para você reler em sessões futuras. Mantenha o checklist de progresso em `docs/PROGRESSO.md`.
3. Commits: Conventional Commits, com tipo em inglês e descrição em português no imperativo (ex.: `feat(assess): detecta SUSPEND em procedures`).
   - Um commit por mudança lógica.
   - O teste vai no mesmo commit do código que ele testa.
   - Nunca misture refatoração com funcionalidade e nunca commite código que não compila.
   - Antes de cada commit, rode `gofmt`, `go vet ./...`, `golangci-lint run` e `go test ./...`.
4. GitHub:
   - Use uma branch por fase (ex.: `fase-1-assess`) e abra PR com `gh pr create` explicando o que mudou e como testar.
   - Espere o CI passar e faça merge com `gh pr merge --rebase --delete-branch`, que preserva os commits separados.
   - Nunca faça force push na `main`.
   - Crie uma milestone por versão (via `gh api`) e issues para as tarefas de cada fase, fechando-as pelo PR (`Closes #N`).
5. Ao fechar cada fase: atualize o `CHANGELOG.md` (formato Keep a Changelog), crie a tag (`v0.1.0`, `v0.2.0`...) e confira se a release foi publicada pelo GoReleaser.
6. Não invente APIs: confira assinaturas com `go doc` ou no código do módulo antes de usar. Se algo deste plano estiver errado ou for inviável, pare, explique e proponha uma alternativa. Não contorne em silêncio.
7. Segurança:
   - Nunca versione credenciais: use `.env` (no `.gitignore`) e um `.env.example`.
   - Mascare senhas de DSNs nos logs.
   - Todos os comandos são somente leitura nos bancos, exceto `sync install/run/uninstall`. Esses mostram o DDL antes (`--dry-run` é o padrão) e só aplicam com `--yes`.
8. Licenças: dependências mínimas e só com licença permissiva. Não copie código de projetos com licença incompatível (o ACBr, por exemplo, é LGPL).

## 3. Contexto: por que o projeto existe
- O Brasil é considerado a maior base de usuários de Firebird do mundo, e muitas software houses Delphi + Firebird estão migrando para PostgreSQL.
- Conversores já existem:
  - o open source firebird2postgres (https://github.com/SamuelTM/firebird2postgres);
  - o DBConvert, comercial, que deixa procedures e triggers para reescrita manual.
- O que não existe é uma forma automatizada de provar que o PostgreSQL resultante tem os mesmos dados e o mesmo comportamento.
- A Firebird Foundation publicou um guia com 23 problemas reais dessa migração: https://firebirdsql.org/img/articles/pgfb/firebird_postgresql_migration_pt_br_v2.pdf
  - Leia antes da Fase 1: ele é a base do catálogo de regras.
  - Se não conseguir acessar, use a lista da Fase 1.
- O noves NÃO converte bancos. Ele avalia o risco antes, prova os dados e o comportamento depois e aponta exatamente onde divergiu. Lema: "Converta com o que quiser. Prove com o noves."
- O nome vem da "prova dos nove" (noves fora). Explique isso no README.

## 4. Decisões técnicas
- Go: a versão estável instalada (1.24 ou mais nova). Binário único com `CGO_ENABLED=0`, porque muitos servidores Firebird rodam em Windows e o noves precisa funcionar sem instalar nada.
- Arquitetura hexagonal:
  - núcleo de domínio sem dependência de banco;
  - portas (interfaces) para catálogo, leitura de linhas e execução de consultas;
  - adaptadores para Firebird e PostgreSQL;
  - caminho aberto para outras origens no futuro (Oracle, SQL Server), mas sem implementá-las agora.
- Dependências sugeridas:
  - `github.com/nakagami/firebirdsql` (driver puro Go, Firebird 2.5+);
  - `github.com/jackc/pgx/v5`;
  - `github.com/spf13/cobra`;
  - `golang.org/x/text` (charmaps e normalização Unicode);
  - `modernc.org/sqlite` (armazenamento temporário puro Go, só se precisar);
  - hash: `crypto/sha256` da stdlib.
- Configuração: flags, variáveis de ambiente (`NOVES_FB_DSN`, `NOVES_PG_DSN`) e um arquivo opcional `noves.yaml` (mapeamento de tabelas e colunas, objetos ignorados, normalização por coluna, pesos da estimativa).
- Saídas:
  - HTML autocontido (CSS e JS inline, sem CDN), porque precisa abrir offline e ir por e-mail;
  - JSON com `schema_version`;
  - Markdown resumido, para comentário em PR.
- Mensagens e relatórios em português por padrão, com `--lang en`.
- Códigos de saída: 0 = ok; 1 = achados ou divergências acima do limite (`--fail-on`); 2 = erro de uso ou conexão.
- Paralelismo com `--jobs`, memória constante (streaming) e progresso no stderr.

## 5. Estrutura sugerida
    cmd/noves/                    main e comandos (cobra)
    internal/core/                domínio: inventário, achados, severidade, diff
    internal/ports/               interfaces
    internal/adapters/firebird/   catálogo RDB$, leitura de linhas, trace
    internal/adapters/postgres/   catálogo, leitura, execução
    internal/psql/                tokenizer de SQL/PSQL do Firebird
    internal/rules/               catálogo de regras e estimativa de esforço
    internal/assess/  internal/scan/  internal/datadiff/  internal/replay/  internal/sync/
    internal/normalize/           normalização de valores entre os bancos
    internal/report/              templates (embed) HTML/MD e JSON
    internal/testkit/             helpers de teste (inclui cópia simples FB→PG só para testes)
    deploy/                       docker-compose e scripts SQL
    docs/                         PLANO, PROGRESSO, ARQUITETURA, comandos/, REGRAS.md (gerado)

## 6. Ambiente de teste (é a vitrine do projeto)
- `deploy/docker-compose.yml` com as imagens oficiais `firebirdsql/firebird:5` e `postgres`, com healthcheck. A imagem do Firebird executa os scripts de `/docker-entrypoint-initdb.d/`.
- Crie um "ERP sintético" no Firebird que imite o legado brasileiro típico:
  - banco com charset NONE ou WIN1252 e textos acentuados, incluindo alguns nomes com mojibake ("JOÃƒO");
  - tabela CLIENTES: ID por generator + trigger BI com GEN_ID, NOME, CPF_CNPJ, ATIVO CHAR(1) 'S'/'N', LIMITE NUMERIC(15,2), DT_CADASTRO;
  - tabela PRODUTOS: descrição em BLOB SUB_TYPE 1 e uma coluna COMPUTED BY;
  - tabelas PEDIDOS e ITENS, com FK;
  - domínios e uma exceção;
  - uma procedure selecionável com SUSPEND e uma com UPDATE OR INSERT;
  - um log com IN AUTONOMOUS TRANSACTION, um EXECUTE BLOCK e uma GTT;
  - um trigger com vários eventos e POSITION;
  - uso de RDB$SET_CONTEXT e consultas com FIRST/SKIP;
  - dezenas de milhares de linhas geradas com EXECUTE BLOCK.
- Crie a versão "migrada" no PostgreSQL: schema escrito à mão + cópia dos dados pelo testkit.
- Crie um script que injeta divergências conhecidas no PostgreSQL:
  - linha faltando;
  - decimal alterado;
  - nome com codificação dupla;
  - CHAR sem espaços;
  - '' virando NULL.
- Os testes de integração exigem que o noves encontre exatamente essas divergências e nenhuma outra.
- Testes de integração com build tag `integration`, rodados por `make test-integration` e no CI.

## 7. Fases

### Fase 0 — Fundação (sem tag)
- Repositório:
  - caminho do módulo `github.com/<meu-usuário>/noves` (descubra com `gh api user -q .login`);
  - criação com `gh repo create noves --public --source . --remote origin --push`, com descrição e topics (firebird, postgresql, migration, database, delphi, go, cli, brasil);
  - milestones v0.1.0 a v1.0.0;
  - se possível, proteja a `main` exigindo CI verde.
- Arquivos:
  - `go.mod` e estrutura de pastas;
  - LICENSE (MIT), README.md (PT) e README.en.md em esqueleto;
  - CONTRIBUTING.md, CODE_OF_CONDUCT.md, SECURITY.md, CHANGELOG.md;
  - `.editorconfig`, `.gitignore`, `.golangci.yml`;
  - Makefile com `build`, `test`, `lint`, `test-integration`, `up`, `down`;
  - templates de issue e PR, Dependabot (gomod e actions).
- CI no GitHub Actions: lint, testes unitários com `-race`, build cruzado (incluindo `GOOS=windows`) e testes de integração com docker compose. Release com GoReleaser disparada por tag.
- Ambiente de teste da seção 6.
- Comandos `noves version` e `noves doctor`: testa as conexões, mostra versão do Firebird, ODS, dialeto e charset e a versão do PostgreSQL, e confere a permissão de leitura.
- Exemplo de sequência de commits:
  1. `chore: inicializa módulo Go e estrutura de diretórios`
  2. `docs: adiciona plano do projeto e CLAUDE.md`
  3. `docs: adiciona README, licença MIT e guias da comunidade`
  4. `chore: adiciona Makefile, golangci-lint e editorconfig`
  5. `build: adiciona docker-compose com Firebird 5 e PostgreSQL`
  6. `test: adiciona ERP sintético no Firebird`
  7. `test: adiciona banco migrado com divergências injetadas`
  8. `feat(cli): adiciona comandos version e doctor`
  9. `ci: adiciona lint, testes e integração`
  10. `ci: adiciona release com GoReleaser`

### Fase 1 — `noves assess` (v0.1.0)
Objetivo: relatório de avaliação do banco Firebird antes da migração.

1. Catálogo pelas tabelas de sistema, consciente da versão (2.5 a 5): confira se a tabela ou coluna existe antes de usar e degrade com elegância. Fontes:
   - RDB$DATABASE e MON$DATABASE (dialeto, ODS, tamanho de página);
   - RDB$GET_CONTEXT('SYSTEM','ENGINE_VERSION');
   - RDB$RELATIONS (tabelas, views, GTTs);
   - RDB$RELATION_FIELDS e RDB$FIELDS (tipos, escala, charset, COMPUTED BY, domínios);
   - RDB$CHARACTER_SETS e RDB$COLLATIONS;
   - índices e constraints;
   - RDB$GENERATORS;
   - RDB$TRIGGERS (decodifique RDB$TRIGGER_TYPE multievento e os triggers de banco e de DDL);
   - RDB$PROCEDURES (selecionável x executável, código-fonte);
   - RDB$FUNCTIONS (PSQL, UDF, UDR), RDB$PACKAGES e RDB$EXCEPTIONS.

   Contagem de linhas só com `--counts`, porque COUNT(*) é caro no Firebird.
2. Tokenizer de SQL/PSQL: comentários, strings com '' de escape, identificadores entre aspas, palavras-chave sem diferenciar maiúsculas. A detecção é por sequência de tokens, nunca regex no texto cru, para não acusar nada dentro de string ou comentário.
3. Catálogo de regras como dados. Cada regra tem ID, título PT/EN, severidade, explicação, recomendação para o PostgreSQL e referência ao item do guia, com detector e testes positivo e negativo. Regras iniciais:
   - FB-SUSPEND (média): procedure selecionável vira função RETURNS TABLE; no PostgreSQL o resultado é materializado inteiro antes de voltar.
   - FB-EXECUTE-BLOCK (alta): bloco anônimo com parâmetros e retorno não tem equivalente (DO não recebe nem devolve valores).
   - FB-AUTONOMOUS-TX (alta): IN AUTONOMOUS TRANSACTION não existe no PostgreSQL puro.
   - FB-UPDATE-OR-INSERT (baixa): vira INSERT ... ON CONFLICT, que exige constraint única.
   - FB-MERGE, FB-FIRST-SKIP e FB-ROWS (baixa): FIRST/SKIP e ROWS viram LIMIT/OFFSET.
   - FB-EXCEPTION e FB-WHEN-HANDLER (média): exceções nomeadas e WHEN SQLCODE/GDSCODE viram RAISE e EXCEPTION WHEN.
   - FB-EXECUTE-STATEMENT (média; alta com ON EXTERNAL).
   - FB-POST-EVENT (média): vira NOTIFY e muda quem escuta na aplicação.
   - FB-CONTEXT-VARS (média): RDB$GET_CONTEXT e RDB$SET_CONTEXT viram set_config e current_setting, e quebram com pool de conexões.
   - FB-GTT, FB-PACKAGE, FB-UDF-UDR e FB-DB-TRIGGER (alta; o último cobre ON CONNECT, ON TRANSACTION etc.).
   - FB-COMPUTED-BY (baixa): vira coluna gerada.
   - FB-GENERATOR-TRIGGER (info): trigger BI que só preenche o ID com GEN_ID vira identity ou DEFAULT nextval.
   - FB-TRIGGER-ORDER (média): POSITION e multievento; no PostgreSQL, triggers do mesmo tipo disparam em ordem alfabética do nome.
   - FB-CHARSET-NONE (alta) e FB-COLLATION-CI-AI (média): PT_BR e WIN_PTBR pedem collation ICU não determinística.
   - FB-DIALECT-1 (alta), FB-DECFLOAT-INT128 (média), FB-BLOB (baixa) e FB-CHAR-FLAG (info: CHAR(1) 'S'/'N').
   - FB-IDENTIFIER-CASE e FB-RESERVED-WORD (média): nomes que mudam de caixa ou são reservados no PostgreSQL.
   - FB-SYSTEM-TABLES (média): uso de RDB$ ou MON$ no código.
   - FB-FUNCTIONS (baixa): IIF, DECODE, LIST, DATEADD, DATEDIFF, GEN_UUID e funções comuns de UDF.
   - FB-RUNTIME-VALIDATION (info): no PostgreSQL, nome errado dentro de função só aparece na execução; recomendar plpgsql_check.
4. Classifique cada procedure como somente leitura ou com escrita, com análise transitiva das chamadas. A Fase 4 usa isso.
5. Estimativa de esforço configurável (pesos por tipo de objeto, tamanho do código e severidade). Ela devolve uma faixa de dias-pessoa e um nível A/B/C, com a fórmula documentada e o aviso de que é estimativa.
6. Relatório com:
   - resumo do banco e achados por severidade;
   - principais riscos explicados;
   - lista filtrável por objeto;
   - estimativa e metodologia.

   Gere `docs/REGRAS.md` a partir do catálogo com `go generate`.

Aceite: no ERP sintético, o assess encontra exatamente os casos plantados; há golden tests do JSON; o HTML abre offline.

### Fase 2 — `noves data` (v0.2.0)
Objetivo: provar que os dados batem.

1. Mapeamento de tabelas e colunas (maiúsculas no Firebird x minúsculas no PostgreSQL, com overrides no `noves.yaml`). Relatar tabelas e colunas sem par.
2. Normalização canônica por tipo antes de comparar:
   - inteiros como estão;
   - NUMERIC sem zeros à direita;
   - float com tolerância configurável;
   - CHAR sem os espaços à direita (configurável);
   - texto decodificado do charset de origem e normalizado em NFC (NONE assume WIN1252 por padrão, configurável);
   - BLOB binário comparado com bytea;
   - DATE, TIME e TIMESTAMP em ISO, truncados à precisão comum (o Firebird guarda décimos de milissegundo);
   - tipos com fuso em UTC;
   - booleanos, com mapeamento opcional 'S'/'N' → true/false;
   - '' x NULL reportado como categoria própria.
3. Algoritmo:
   a) Caminho rápido: ler as duas tabelas em streaming, sem ORDER BY, calculando a contagem e uma impressão digital independente da ordem (soma mod 2^128 do hash de cada linha canônica, com prefixo de tamanho em cada campo). Se forem iguais, a tabela está ok.
   b) Diferentes e PK numérica simples: bisseção por faixas de PK (WHERE pk BETWEEN) até blocos pequenos. Nesses blocos, diff linha a linha, separando linhas faltando, sobrando e alteradas (com as colunas que mudaram).
   c) PK composta ou de texto: diff por chave em armazenamento temporário. Não dependa de ORDER BY em texto, porque as collations dos dois bancos ordenam diferente.
   d) Sem PK: só contagem e impressão digital, com aviso.
4. Auditoria de encoding:
   - leia os bytes crus (conexão com charset NONE ou CAST para OCTETS);
   - acuse bytes inválidos e mojibake ("Ã§", "Ã£", "Ã©"...) na origem e codificação dupla no destino;
   - mascare as amostras por padrão, porque podem conter CPF e nomes (LGPD); `--show-samples` exibe.
5. Opções `--tables`, `--exclude`, `--where` por tabela e `--sample` (percentual) para bancos grandes.

Aceite: nos testes de integração, encontra exatamente as divergências injetadas e nenhuma outra; uma tabela de 1 milhão de linhas roda sem estourar a memória.

### Fase 3 — `noves scan` (v0.3.0)
Objetivo: medir a migração do lado da aplicação.

1. Varrer projetos Delphi/Lazarus:
   - .dfm e .lfm em texto: propriedades SQL.Strings, SQL.Text e CommandText, com concatenação '+' e caracteres #nn;
   - .pas: literais passados a SQL.Add, SQL.Text, ExecSQL, Open etc.;
   - detectar DFM binário (assinatura TPF0) e converter ou avisar.
2. Aplicar o mesmo tokenizer e as mesmas regras ao SQL extraído, com arquivo e linha.
3. Acusar várias transações por conexão, já que o PostgreSQL só tem uma transação por conexão:
   - mais de um componente de transação (TFDTransaction, TpFIBTransaction, TIBTransaction, TUniTransaction, TSQLTransaction) ligado à mesma conexão;
   - datasets com Transaction e UpdateTransaction diferentes.
4. `noves assess --scan ./src` junta os achados da aplicação no mesmo relatório.

Aceite: fixtures realistas de .dfm, .lfm e .pas nos testes.

### Fase 4 — `noves replay` (v0.4.0)
Objetivo: provar que o comportamento bate usando o uso real. Pré-requisito: `noves data` sem divergências.

1. `noves replay trace-config` gera a configuração do trace do Firebird (statements e procedures finalizados, com parâmetros). O usuário captura com o fbtracemgr.
2. `noves replay parse trace.log -o workload.jsonl`: parser tolerante do log do fbtracemgr para Firebird 3, 4 e 5 (e 2.5 se der). Gere logs reais nos testes rodando o fbtracemgr dentro do container.
3. Deduplicar por impressão digital da instrução (literais trocados por ?) e amostrar N execuções de cada.
4. Classificar em:
   - somente leitura: SELECT sem WITH LOCK/FOR UPDATE e procedures marcadas como leitura na Fase 1;
   - escrita: reportada como "não verificada".
5. Traduzir para o PostgreSQL as construções comuns:
   - FIRST/SKIP/ROWS;
   - caixa dos identificadores;
   - parâmetros ? e :p para $n;
   - IIF, LIST e FROM RDB$DATABASE;
   - EXECUTE PROCEDURE e SELECT FROM procedure.

   Inclua um `rewrites.yaml` com traduções manuais por impressão digital.
6. Executar nos dois bancos em transações somente leitura com timeout:
   - comparar em ordem se houver ORDER BY e como multiconjunto se não houver;
   - mesma normalização da Fase 2;
   - tempo pela mediana de algumas execuções.
7. Relatório com:
   - % do workload provado equivalente;
   - divergências com exemplo;
   - erros no PostgreSQL (função inexistente etc.);
   - consultas muito mais lentas no PostgreSQL (ex.: 10x).

Aceite: um workload de exemplo no ERP sintético, com as divergências plantadas detectadas.

### Fase 5 — `noves sync` (v0.5.0)
Objetivo: manter o PostgreSQL atualizado até a virada, com CDC por triggers.

1. `sync install`:
   - cria no Firebird a tabela NOVES$LOG (id por sequence, tabela, PK, operação, timestamp);
   - cria triggers AFTER INSERT/UPDATE/DELETE por tabela, com POSITION alta para rodar por último;
   - num UPDATE que muda a PK, registra a PK antiga e a nova.
2. `sync run`:
   - lê lotes do log em ordem e agrupa por tabela e PK;
   - relê a linha atual no Firebird e faz upsert no PostgreSQL, ou DELETE se a linha sumiu;
   - depois apaga do log exatamente os ids processados. Não use "maior id processado": transações que comitam fora de ordem criariam buracos;
   - deve ser idempotente;
   - aplique com `session_replication_role = replica`, para não disparar triggers do destino, e documente o privilégio necessário.
3. `sync status` (atraso, pendências, taxa) e `sync uninstall`.
4. Documente a ordem segura: install → cópia inicial com o conversor → run → `noves data` → virada.

Aceite: teste ponta a ponta que instala, altera o Firebird em paralelo, sincroniza e termina com `noves data` sem divergências.

### Fase 6 — Polimento e v1.0.0
- README completo em português (e versão em inglês): o problema, o nome, instalação, início rápido em 3 comandos com o docker compose, exemplos de cada comando, GIF e badges (CI, release, licença, Go Report Card).
- Relatório de exemplo, gerado a partir do ERP sintético, publicado no GitHub Pages e linkado no README.
- GoReleaser:
  - binários para Windows, Linux e macOS (amd64 e arm64);
  - zip para Windows e checksums;
  - notas de release a partir do CHANGELOG.
- Matriz de CI: Firebird 3, 4 e 5; PostgreSQL 15 (mínimo suportado) e a última versão estável.
- Cobertura de pelo menos 70% nos pacotes do núcleo, documentação de cada comando em `docs/comandos/` e o guia "como adicionar uma regra" no CONTRIBUTING.
- A v1.0.0 está pronta quando:
  - os cinco comandos têm testes;
  - o CI está verde na matriz;
  - a release está publicada com binários;
  - o relatório de exemplo está no ar;
  - CHANGELOG e documentação estão completos.

## 8. Fora do escopo
- Converter schema, procedures ou dados. No teste, o schema do PostgreSQL é escrito à mão e a cópia é um helper do testkit.
- Interface gráfica e outros bancos de origem (deixe só a arquitetura preparada).

Comece agora: confira o ambiente (regra 1) e execute a Fase 0.
