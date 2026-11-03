# Política de segurança

## Versões suportadas

Enquanto o projeto estiver antes da v1.0.0, só a última release recebe correções.

## Como relatar uma vulnerabilidade

Não abra issue pública. Use o
[relato privado de vulnerabilidades do GitHub](https://github.com/EduardoMilani8/noves/security/advisories/new).

Responderemos em até 7 dias.

## O que o noves promete

- **Somente leitura:** todos os comandos só leem os bancos, exceto
  `noves sync install/run/uninstall`. Esses mostram o DDL antes (`--dry-run` é o
  padrão) e só aplicam com `--yes`.
- **Senhas mascaradas:** DSNs aparecem com a senha trocada por `***` em logs e
  relatórios.
- **Dados pessoais mascarados:** amostras de dados nos relatórios (que podem
  conter CPF e nomes) são mascaradas por padrão, por causa da LGPD. Só aparecem
  com `--show-samples`.

Se encontrar um caso em que alguma dessas promessas falha, trate como
vulnerabilidade.
