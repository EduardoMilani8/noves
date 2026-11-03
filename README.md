# noves

> A prova dos nove de migrações Firebird → PostgreSQL.
> **Converta com o que quiser. Prove com o noves.**

[English](README.en.md)

> ⚠️ Em desenvolvimento. Acompanhe em [docs/PROGRESSO.md](docs/PROGRESSO.md).

## O problema

O Brasil tem uma das maiores bases de Firebird do mundo, e muitas software houses
Delphi + Firebird estão migrando para PostgreSQL. Já existem conversores, como o
[firebird2postgres](https://github.com/SamuelTM/firebird2postgres) e o DBConvert.
O que falta é uma forma automatizada de **provar** que o PostgreSQL resultante
tem os mesmos dados e o mesmo comportamento.

O noves não converte nada. Ele:

1. **avalia o risco antes** (`noves assess`): lê o catálogo do Firebird e aponta
   o que não tem equivalente direto no PostgreSQL;
2. **prova os dados depois** (`noves data`): compara tabela a tabela e mostra
   exatamente a linha e a coluna que divergiram;
3. **mede a aplicação** (`noves scan`): acha o SQL dentro de `.dfm`, `.lfm` e `.pas`;
4. **prova o comportamento** (`noves replay`): reexecuta o uso real capturado
   pelo trace do Firebird nos dois bancos e compara os resultados;
5. **acompanha até a virada** (`noves sync`): mantém o PostgreSQL atualizado
   enquanto o Firebird continua em produção.

## Por que "noves"?

Vem da **prova dos nove** (ou "noves fora"), o truque que se aprende na escola
para conferir uma conta: some os algarismos de cada número, descarte os noves e
veja se o resultado bate dos dois lados. É rápido, não refaz a conta inteira e
pega quase todo erro.

O noves faz o mesmo com bancos: calcula uma impressão digital de cada tabela nos
dois lados. Se bate, a tabela está provada. Se não bate, ele desce até achar a
linha que diverge.

## Instalação

Ainda não há release. Para compilar a partir do código:

```bash
go install github.com/EduardoMilani8/noves/cmd/noves@latest
```

O binário é estático (`CGO_ENABLED=0`): funciona em Windows, Linux e macOS sem
instalar cliente de banco nenhum.

## Início rápido

```bash
make up                   # sobe Firebird 5 + PostgreSQL com um ERP sintético
cp .env.example .env      # DSNs do ambiente de teste
go run ./cmd/noves doctor # confere conexões, versões e permissões
```

## Comandos

| Comando | Status |
|---|---|
| `noves version` | disponível |
| `noves doctor` | disponível |
| `noves assess` | planejado (v0.1.0) |
| `noves data` | planejado (v0.2.0) |
| `noves scan` | planejado (v0.3.0) |
| `noves replay` | planejado (v0.4.0) |
| `noves sync` | planejado (v0.5.0) |

## Contribuindo

Veja [CONTRIBUTING.md](CONTRIBUTING.md).

## Licença

[MIT](LICENSE)
