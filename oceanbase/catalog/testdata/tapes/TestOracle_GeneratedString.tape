-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_GeneratedString default#1
> exec
| DROP DATABASE IF EXISTS zz_generated_string
< ok 0 0
> exec
| CREATE DATABASE zz_generated_string
< ok 1 0
> exec
| USE zz_generated_string
< ok 0 0
> exec
| CREATE TABLE g (a VARCHAR(10) AS ('it''s'))
< ok 0 0
> query
| SHOW CREATE TABLE zz_generated_string.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` varchar(10) GENERATED ALWAYS AS ('it\'s') VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_generated_string
< ok 1 0
> exec
| CREATE DATABASE zz_generated_string
< ok 1 0
> exec
| USE zz_generated_string
< ok 0 0
> exec
| CREATE TABLE g (a VARCHAR(10) AS (_latin1 X'27'))
< ok 0 0
> query
| SHOW CREATE TABLE zz_generated_string.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` varchar(10) GENERATED ALWAYS AS ('\'') VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_generated_string
< ok 1 0
> exec
| CREATE DATABASE zz_generated_string
< ok 1 0
> exec
| USE zz_generated_string
< ok 0 0
> exec
| CREATE TABLE g (a VARCHAR(10) AS (_latin1 X'5c'))
< ok 0 0
> query
| SHOW CREATE TABLE zz_generated_string.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` varchar(10) GENERATED ALWAYS AS ('\\') VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_generated_string
< ok 1 0
> exec
| CREATE DATABASE zz_generated_string
< ok 1 0
> exec
| USE zz_generated_string
< ok 0 0
> exec
| CREATE TABLE g (a VARCHAR(10) AS ('a\\b'))
< ok 0 0
> query
| SHOW CREATE TABLE zz_generated_string.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` varchar(10) GENERATED ALWAYS AS ('a\\b') VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
