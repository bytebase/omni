-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_NumberIsDecimal default#1
> exec
| DROP DATABASE IF EXISTS zz_number
< ok 0 0
> exec
| CREATE DATABASE zz_number
< ok 1 0
> exec
| USE zz_number
< ok 0 0
> exec
| CREATE TABLE g (a NUMBER)
< ok 0 0
> query
| SHOW CREATE TABLE zz_number.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` decimal(10,0) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_number
< ok 1 0
> exec
| CREATE DATABASE zz_number
< ok 1 0
> exec
| USE zz_number
< ok 0 0
> exec
| CREATE TABLE g (a NUMBER(5))
< ok 0 0
> query
| SHOW CREATE TABLE zz_number.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` decimal(5,0) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_number
< ok 1 0
> exec
| CREATE DATABASE zz_number
< ok 1 0
> exec
| USE zz_number
< ok 0 0
> exec
| CREATE TABLE g (a NUMBER(5,2))
< ok 0 0
> query
| SHOW CREATE TABLE zz_number.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` decimal(5,2) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_number
< ok 1 0
> exec
| CREATE DATABASE zz_number
< ok 1 0
> exec
| USE zz_number
< ok 0 0
> exec
| CREATE TABLE g (a NUMBER(5,2) UNSIGNED)
< ok 0 0
> query
| SHOW CREATE TABLE zz_number.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` decimal(5,2) unsigned DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
