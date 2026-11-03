/*
 * ERP sintetico do noves: dados.
 *
 * Tudo e deterministico (MOD sobre o indice, nada de RAND), para que os testes
 * de integracao possam apontar linhas especificas. Volume aproximado:
 *   CLIENTES 5.000 | PRODUTOS 1.000 | PEDIDOS 20.000 | ITENS 50.000
 *
 * Bytes WIN1252 usados com ASCII_CHAR:
 *   195 = A til   199 = C cedilha   201 = E agudo   205 = I agudo
 *   211 = O agudo 212 = O circunflexo  218 = U agudo  131 = f (mojibake)
 */

SET TERM ^ ;

/* Desliga a auditoria durante a carga (lido pelo trigger PEDIDOS_AUDITORIA). */
EXECUTE BLOCK AS
DECLARE VARIABLE R INTEGER;
BEGIN
    R = RDB$SET_CONTEXT('USER_SESSION', 'NOVES_CARGA', '1');
END^

/* CLIENTES */
EXECUTE BLOCK AS
DECLARE VARIABLE I INTEGER = 1;
DECLARE VARIABLE PRIMEIRO VARCHAR(30) CHARACTER SET NONE;
DECLARE VARIABLE SOBRENOME VARCHAR(30) CHARACTER SET NONE;
BEGIN
    WHILE (I <= 5000) DO
    BEGIN
        PRIMEIRO = CASE MOD(I, 12)
            WHEN 0 THEN 'JO' || ASCII_CHAR(195) || 'O'
            WHEN 1 THEN 'JOS' || ASCII_CHAR(201)
            WHEN 2 THEN 'MARIA'
            WHEN 3 THEN 'ANT' || ASCII_CHAR(212) || 'NIO'
            WHEN 4 THEN 'FRANCISCO'
            WHEN 5 THEN 'ANA'
            WHEN 6 THEN 'LU' || ASCII_CHAR(205) || 'S'
            WHEN 7 THEN 'CONCEI' || ASCII_CHAR(199) || ASCII_CHAR(195) || 'O'
            WHEN 8 THEN 'PAULO'
            WHEN 9 THEN 'SEBASTI' || ASCII_CHAR(195) || 'O'
            WHEN 10 THEN 'RAIMUNDA'
            ELSE 'CARLOS'
        END;

        /* Alguns nomes ja chegaram corrompidos no legado: UTF-8 gravado como WIN1252. */
        IF (MOD(I, 997) = 0) THEN
            PRIMEIRO = 'JO' || ASCII_CHAR(195) || ASCII_CHAR(131) || 'O';

        SOBRENOME = CASE MOD(I * 7, 10)
            WHEN 0 THEN 'SILVA'
            WHEN 1 THEN 'SANTOS'
            WHEN 2 THEN 'OLIVEIRA'
            WHEN 3 THEN 'SOUZA'
            WHEN 4 THEN 'PEREIRA'
            WHEN 5 THEN 'LIMA'
            WHEN 6 THEN 'GON' || ASCII_CHAR(199) || 'ALVES'
            WHEN 7 THEN 'FALC' || ASCII_CHAR(195) || 'O'
            WHEN 8 THEN 'ARA' || ASCII_CHAR(218) || 'JO'
            ELSE 'RIBEIRO'
        END;

        INSERT INTO CLIENTES (ID, NOME, CPF_CNPJ, ATIVO, LIMITE, DT_CADASTRO, UF, OBS)
        VALUES (
            :I,
            :PRIMEIRO || ' ' || :SOBRENOME,
            CASE WHEN MOD(:I, 10) = 0
                THEN LPAD(CAST(MOD(:I * 31, 100) AS VARCHAR(10)), 2, '0') || '.' || LPAD(CAST(MOD(:I * 7919, 1000) AS VARCHAR(10)), 3, '0') || '.'
                     || LPAD(CAST(MOD(:I * 104729, 1000) AS VARCHAR(10)), 3, '0') || '/0001-' || LPAD(CAST(MOD(:I, 100) AS VARCHAR(10)), 2, '0')
                ELSE LPAD(CAST(MOD(:I * 7919, 1000) AS VARCHAR(10)), 3, '0') || '.' || LPAD(CAST(MOD(:I * 104729, 1000) AS VARCHAR(10)), 3, '0') || '.'
                     || LPAD(CAST(MOD(:I * 31, 1000) AS VARCHAR(10)), 3, '0') || '-' || LPAD(CAST(MOD(:I, 100) AS VARCHAR(10)), 2, '0')
            END,
            CASE WHEN MOD(:I, 17) = 0 THEN 'N' ELSE 'S' END,
            CAST(MOD(:I * 137, 1000000) AS NUMERIC(15,2)) / 100,
            TIMESTAMP '2015-01-01 08:00:00' + :I * 0.0137,
            CASE MOD(:I, 5) WHEN 0 THEN 'RS' WHEN 1 THEN 'SC' WHEN 2 THEN 'PR' WHEN 3 THEN 'SP' ELSE 'MG' END,
            CASE MOD(:I, 3)
                WHEN 0 THEN ''
                WHEN 1 THEN NULL
                ELSE 'CLIENTE DESDE ' || (2015 + MOD(:I, 10))
            END
        );
        I = I + 1;
    END
END^

/* PRODUTOS */
EXECUTE BLOCK AS
DECLARE VARIABLE J INTEGER = 1;
DECLARE VARIABLE NOME VARCHAR(30) CHARACTER SET NONE;
BEGIN
    WHILE (J <= 1000) DO
    BEGIN
        NOME = CASE MOD(J, 8)
            WHEN 0 THEN 'CAF' || ASCII_CHAR(201)
            WHEN 1 THEN 'A' || ASCII_CHAR(199) || ASCII_CHAR(218) || 'CAR'
            WHEN 2 THEN 'FEIJ' || ASCII_CHAR(195) || 'O'
            WHEN 3 THEN 'P' || ASCII_CHAR(195) || 'O DE QUEIJO'
            WHEN 4 THEN 'MA' || ASCII_CHAR(199) || ASCII_CHAR(195)
            WHEN 5 THEN ASCII_CHAR(211) || 'LEO DE SOJA'
            WHEN 6 THEN 'SAB' || ASCII_CHAR(195) || 'O EM BARRA'
            ELSE 'ARROZ'
        END;

        INSERT INTO PRODUTOS (ID, CODIGO, DESCRICAO, UNIDADE, PRECO, CUSTO, ESTOQUE, ATIVO)
        VALUES (
            :J,
            'P' || LPAD(CAST(:J AS VARCHAR(10)), 5, '0'),
            :NOME || ' - ITEM ' || :J || ASCII_CHAR(13) || ASCII_CHAR(10)
                || 'Linha de produto ' || MOD(:J, 20) || ', embalagem ' || (1 + MOD(:J, 6)) || ' kg.',
            CASE WHEN MOD(:J, 4) = 0 THEN 'KG' ELSE 'UN' END,
            CAST(MOD(:J * 311, 50000) + 99 AS NUMERIC(15,4)) / 100,
            CASE WHEN MOD(:J, 9) = 0 THEN NULL
                 ELSE CAST(MOD(:J * 311, 50000) + 99 AS NUMERIC(15,4)) * 6 / 1000 END,
            CAST(MOD(:J * 13, 500) AS NUMERIC(15,3)) + 0.125,
            CASE WHEN MOD(:J, 23) = 0 THEN 'N' ELSE 'S' END
        );
        J = J + 1;
    END
END^

/* PEDIDOS e ITENS (o trigger ITENS_TOTAL calcula PEDIDOS.TOTAL) */
EXECUTE BLOCK AS
DECLARE VARIABLE K INTEGER = 1;
DECLARE VARIABLE N INTEGER;
BEGIN
    WHILE (K <= 20000) DO
    BEGIN
        INSERT INTO PEDIDOS (ID, CLIENTE_ID, DT_PEDIDO, HR_PEDIDO, STATUS, TOTAL, OBS)
        VALUES (
            :K,
            1 + MOD(:K * 7, 5000),
            DATE '2020-01-01' + MOD(:K * 3, 2190),
            TIME '08:00:00' + MOD(:K * 97, 36000),
            CASE MOD(:K, 10) WHEN 0 THEN 'C' WHEN 1 THEN 'A' WHEN 2 THEN 'A' ELSE 'F' END,
            0,
            CASE MOD(:K, 4)
                WHEN 0 THEN NULL
                WHEN 1 THEN 'ENTREGA AP' || ASCII_CHAR(211) || 'S 18H'
                WHEN 2 THEN ''
                ELSE 'RETIRA NO BALC' || ASCII_CHAR(195) || 'O'
            END
        );

        N = 1;
        WHILE (N <= 1 + MOD(K, 4)) DO
        BEGIN
            INSERT INTO ITENS (PEDIDO_ID, ITEM, PRODUTO_ID, QTDE, VL_UNIT, DESCONTO)
            VALUES (
                :K,
                :N,
                1 + MOD(:K * 31 + :N * 17, 1000),
                CAST(1 + MOD(:K + :N, 10) AS NUMERIC(15,3)),
                CAST(MOD(:K * 13 + :N, 9000) + 50 AS NUMERIC(15,4)) / 100,
                CAST(MOD(:K, 7) AS DOUBLE PRECISION) / 3
            );
            N = N + 1;
        END
        K = K + 1;
    END
END^

/* MOVIMENTO_ESTOQUE (sem PK) e CONFIGURACOES (PK de texto) */
EXECUTE BLOCK AS
DECLARE VARIABLE M INTEGER = 1;
BEGIN
    WHILE (M <= 3000) DO
    BEGIN
        INSERT INTO MOVIMENTO_ESTOQUE (PRODUTO_ID, DT_MOVTO, TIPO, QTDE)
        VALUES (1 + MOD(:M * 11, 1000), TIMESTAMP '2024-01-01 00:00:00' + :M * 0.25,
                CASE WHEN MOD(:M, 3) = 0 THEN 'S' ELSE 'E' END,
                CAST(1 + MOD(:M, 50) AS NUMERIC(15,3)));
        M = M + 1;
    END
END^

SET TERM ; ^

INSERT INTO CONFIGURACOES (CHAVE, VALOR, "Descricao") VALUES ('EMPRESA.NOME', 'COMERCIAL NOVES LTDA', 'Razao social');
INSERT INTO CONFIGURACOES (CHAVE, VALOR, "Descricao") VALUES ('EMPRESA.UF', 'RS', 'Estado da matriz');
INSERT INTO CONFIGURACOES (CHAVE, VALOR, "Descricao") VALUES ('NFE.SERIE', '1', NULL);
INSERT INTO CONFIGURACOES (CHAVE, VALOR, "Descricao") VALUES ('NFE.AMBIENTE', '2', 'Homologacao');
INSERT INTO CONFIGURACOES (CHAVE, VALOR, "Descricao") VALUES ('estoque.permite_negativo', 'N', NULL);

/* Generators continuam de onde a carga parou. */
SET GENERATOR GEN_CLIENTES_ID TO 5000;
SET GENERATOR GEN_PRODUTOS_ID TO 1000;
SET GENERATOR GEN_PEDIDOS_ID TO 20000;

COMMIT;

/* Religa a auditoria e registra o fim da carga pelo caminho normal. */
SET TERM ^ ;
EXECUTE BLOCK AS
DECLARE VARIABLE R INTEGER;
BEGIN
    R = RDB$SET_CONTEXT('USER_SESSION', 'NOVES_CARGA', NULL);
END^
SET TERM ; ^

EXECUTE PROCEDURE SP_DEFINE_USUARIO_APP('CARGA_INICIAL');
EXECUTE PROCEDURE SP_LOG('*', 'I', 'Carga inicial do ERP sintetico concluida');

COMMIT;
