-- Versão "migrada" do ERP sintético, escrita à mão como faria um conversor.
-- Decisões típicas de conversor que o noves precisa enxergar:
--   * nomes em minúsculas;
--   * CHAR(10) do código do produto virou varchar(10);
--   * domínios e generators somem; COMPUTED BY vira coluna gerada;
--   * a GTT TMP_RELATORIO não é migrada.

CREATE TABLE clientes (
    id          integer PRIMARY KEY,
    nome        varchar(60) NOT NULL,
    cpf_cnpj    varchar(18),
    ativo       char(1) NOT NULL DEFAULT 'S' CHECK (ativo IN ('S', 'N')),
    limite      numeric(15,2) DEFAULT 0,
    dt_cadastro timestamp NOT NULL DEFAULT now(),
    uf          char(2),
    obs         varchar(200)
);

CREATE TABLE produtos (
    id        integer PRIMARY KEY,
    codigo    varchar(10) NOT NULL UNIQUE,
    descricao text,
    unidade   char(2) DEFAULT 'UN',
    preco     numeric(15,4) NOT NULL,
    custo     numeric(15,4),
    estoque   numeric(15,3) DEFAULT 0,
    margem    numeric GENERATED ALWAYS AS (CASE WHEN custo > 0 THEN (preco - custo) * 100 / custo END) STORED,
    ativo     char(1) NOT NULL DEFAULT 'S' CHECK (ativo IN ('S', 'N'))
);

CREATE TABLE pedidos (
    id         integer PRIMARY KEY,
    cliente_id integer NOT NULL REFERENCES clientes (id),
    dt_pedido  date NOT NULL,
    hr_pedido  time,
    status     char(1) NOT NULL DEFAULT 'A' CHECK (status IN ('A', 'F', 'C')),
    total      numeric(15,2) DEFAULT 0,
    obs        varchar(200)
);

CREATE TABLE itens (
    pedido_id  integer NOT NULL REFERENCES pedidos (id) ON DELETE CASCADE,
    item       smallint NOT NULL,
    produto_id integer NOT NULL REFERENCES produtos (id),
    qtde       numeric(15,3) NOT NULL,
    vl_unit    numeric(15,4) NOT NULL,
    desconto   double precision DEFAULT 0,
    PRIMARY KEY (pedido_id, item)
);

CREATE INDEX ix_pedidos_data ON pedidos (dt_pedido);

CREATE TABLE configuracoes (
    chave     varchar(40) PRIMARY KEY,
    valor     varchar(200),
    descricao varchar(100)
);

CREATE TABLE movimento_estoque (
    produto_id integer NOT NULL,
    dt_movto   timestamp NOT NULL,
    tipo       char(1) NOT NULL,
    qtde       numeric(15,3) NOT NULL
);

CREATE TABLE log_auditoria (
    id       bigint PRIMARY KEY,
    dt_hora  timestamp NOT NULL DEFAULT now(),
    usuario  varchar(63),
    tabela   varchar(63),
    operacao char(1),
    detalhe  varchar(500)
);

CREATE VIEW v_clientes_ativos (id, nome, uf, limite) AS
    SELECT id, nome, uf, limite FROM clientes WHERE ativo = 'S';
