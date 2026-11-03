# Progresso

Checklist vivo do plano (`docs/PLANO.md`). Marque ao concluir.

## Fase 0 — Fundação
- [x] Módulo Go e estrutura de diretórios
- [x] Repositório público, topics e milestones v0.1.0–v1.0.0
- [x] Plano, CLAUDE.md e este checklist
- [ ] README (PT/EN), LICENSE, guias da comunidade e CHANGELOG
- [ ] Makefile, golangci-lint, editorconfig, templates e Dependabot
- [ ] docker-compose com Firebird 5 e PostgreSQL
- [ ] ERP sintético no Firebird
- [ ] Banco migrado no PostgreSQL com divergências injetadas
- [ ] `noves version` e `noves doctor`
- [ ] CI (lint, testes com -race, build cruzado, integração)
- [ ] Release com GoReleaser
- [ ] Proteção da `main` exigindo CI verde

## Fase 1 — `noves assess` (v0.1.0)
- [ ] Catálogo consciente da versão (2.5 a 5)
- [ ] Tokenizer SQL/PSQL
- [ ] Catálogo de regras com testes positivo e negativo
- [ ] Classificação leitura/escrita de procedures (transitiva)
- [ ] Estimativa de esforço
- [ ] Relatórios HTML/JSON/MD e `docs/REGRAS.md` gerado

## Fase 2 — `noves data` (v0.2.0)
- [ ] Mapeamento de tabelas e colunas
- [ ] Normalização canônica
- [ ] Impressão digital, bisseção e diff por chave
- [ ] Auditoria de encoding
- [ ] Opções `--tables`, `--exclude`, `--where`, `--sample`

## Fase 3 — `noves scan` (v0.3.0)
- [ ] Extração de SQL de .dfm/.lfm/.pas
- [ ] Regras aplicadas ao SQL extraído
- [ ] Várias transações por conexão
- [ ] `assess --scan`

## Fase 4 — `noves replay` (v0.4.0)
- [ ] trace-config e parser do fbtracemgr
- [ ] Deduplicação, classificação e tradução
- [ ] Execução comparada e relatório

## Fase 5 — `noves sync` (v0.5.0)
- [ ] install / run / status / uninstall
- [ ] Teste ponta a ponta

## Fase 6 — v1.0.0
- [ ] README completo, GIF, badges
- [ ] Relatório de exemplo no GitHub Pages
- [ ] Matriz de CI (FB 3/4/5 × PG 15/última)
- [ ] Cobertura ≥ 70% no núcleo e docs dos comandos
