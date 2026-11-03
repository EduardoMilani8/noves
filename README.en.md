# noves

> The "casting out nines" check for Firebird → PostgreSQL migrations.
> **Convert with whatever you like. Prove it with noves.**

[Português](README.md)

> ⚠️ Work in progress. See [docs/PROGRESSO.md](docs/PROGRESSO.md) (in Portuguese).

## The problem

Brazil has one of the largest Firebird user bases in the world, and many
Delphi + Firebird software houses are moving to PostgreSQL. Converters already
exist, such as [firebird2postgres](https://github.com/SamuelTM/firebird2postgres)
and DBConvert. What is missing is an automated way to **prove** that the
resulting PostgreSQL database holds the same data and behaves the same way.

noves converts nothing. It:

1. **assesses risk before** (`noves assess`): reads the Firebird catalog and
   flags what has no direct PostgreSQL equivalent;
2. **proves the data after** (`noves data`): compares table by table and points
   to the exact row and column that diverged;
3. **measures the application** (`noves scan`): finds SQL inside `.dfm`, `.lfm`
   and `.pas` files;
4. **proves behavior** (`noves replay`): replays real workload captured by the
   Firebird trace on both databases and compares the results;
5. **keeps up until cutover** (`noves sync`): keeps PostgreSQL up to date while
   Firebird is still in production.

## Why "noves"?

It comes from *prova dos nove* ("casting out nines"), the school trick for
checking arithmetic: add up the digits of each number, drop the nines and see if
both sides match. It is quick, does not redo the whole computation and catches
almost every mistake.

noves does the same for databases: it computes a fingerprint of each table on
both sides. If it matches, the table is proven. If not, it narrows down to the
row that differs.

## Installation

No release yet. To build from source:

```bash
go install github.com/EduardoMilani8/noves/cmd/noves@latest
```

The binary is static (`CGO_ENABLED=0`): it runs on Windows, Linux and macOS
without any database client installed.

## Quick start

```bash
make up                                # Firebird 5 + PostgreSQL with a synthetic ERP
cp .env.example .env                   # test environment DSNs
go run ./cmd/noves doctor --lang en    # checks connections, versions and permissions
```

## License

[MIT](LICENSE)
