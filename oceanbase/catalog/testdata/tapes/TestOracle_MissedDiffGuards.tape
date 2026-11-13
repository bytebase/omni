-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_MissedDiffGuards default#1
> exec
| DROP DATABASE IF EXISTS mdg_a
< ok 0 0
> exec
| CREATE DATABASE mdg_a
< ok 1 0
> exec
| USE mdg_a
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10)) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> query
| SHOW CREATE TABLE mdg_a.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(10) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS mdg_b
< ok 0 0
> exec
| CREATE DATABASE mdg_b
< ok 1 0
> exec
| USE mdg_b
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10)) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> query
| SHOW CREATE TABLE mdg_b.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS mdg_bin
< ok 0 0
> exec
| CREATE DATABASE mdg_bin
< ok 1 0
> exec
| USE mdg_bin
< ok 0 0
> exec
| CREATE TABLE t (b VARCHAR(10) BINARY, p VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE mdg_bin.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `p` varchar(10) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS mdg_year0
< ok 0 0
> exec
| CREATE DATABASE mdg_year0
< ok 1 0
> exec
| USE mdg_year0
< ok 0 0
> exec
| CREATE TABLE t (n YEAR DEFAULT 0, s YEAR DEFAULT '0')
< ok 0 0
> query
| SHOW CREATE TABLE mdg_year0.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `n` year(4) DEFAULT '0000',
|   `s` year(4) DEFAULT '2000'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
