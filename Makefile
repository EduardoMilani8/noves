COMPOSE  := docker compose -f deploy/docker-compose.yml
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT   ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE     ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS  := -s -w \
	-X github.com/EduardoMilani8/noves/internal/buildinfo.Version=$(VERSION) \
	-X github.com/EduardoMilani8/noves/internal/buildinfo.Commit=$(COMMIT) \
	-X github.com/EduardoMilani8/noves/internal/buildinfo.Date=$(DATE)

.PHONY: build test lint fmt vet check test-integration up down logs seed

build: ## Compila o binário estático em bin/noves
	CGO_ENABLED=0 go build -trimpath -ldflags '$(LDFLAGS)' -o bin/noves ./cmd/noves

test: ## Testes de unidade
	go test -race ./...

lint: ## golangci-lint
	golangci-lint run

fmt: ## Formata o código
	gofmt -w .

vet:
	go vet ./...

check: ## Tudo o que roda antes de cada commit
	@test -z "$$(gofmt -l .)" || (gofmt -l . && echo 'rode make fmt' && exit 1)
	go vet ./...
	golangci-lint run
	go test ./...

test-integration: ## Testes de integração (exige make up)
	go test -race -tags integration -count=1 ./...

up: ## Sobe Firebird + PostgreSQL de teste, copia os dados e injeta as divergências
	$(COMPOSE) up -d --wait
	$(MAKE) seed

seed: ## Recopia Firebird -> PostgreSQL e reaplica as divergências
	go run ./internal/testkit/cmd/seed

down: ## Derruba o ambiente e apaga os volumes
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f
