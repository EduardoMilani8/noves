-- Divergências conhecidas, aplicadas depois da cópia dos dados.
-- Os testes de integração exigem que o noves encontre exatamente estas cinco
-- e nenhuma outra. Se mudar algo aqui, mude também internal/testkit/divergencias.go.

-- 1. Linha faltando (PK composta).
DELETE FROM itens WHERE pedido_id = 1234 AND item = 2;

-- 2. Decimal alterado: 57.54 -> 57.55.
UPDATE clientes SET limite = limite + 0.01 WHERE id = 42;

-- 3. Nome com codificação dupla: 'JOSÉ SANTOS' -> 'JOSÃ‰ SANTOS'.
UPDATE clientes SET nome = convert_from(convert_to(nome, 'UTF8'), 'WIN1252') WHERE id = 13;

-- 4. CHAR sem os espaços à direita: 'P00077    ' -> 'P00077'.
UPDATE produtos SET codigo = rtrim(codigo) WHERE id = 77;

-- 5. '' virando NULL.
UPDATE clientes SET obs = NULL WHERE id = 300;
