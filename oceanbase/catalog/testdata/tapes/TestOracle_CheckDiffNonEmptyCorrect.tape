-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_CheckDiffNonEmptyCorrect default#1
> exec
| DROP DATABASE IF EXISTS chk_nf
< ok 0 0
> exec
| CREATE DATABASE chk_nf
< ok 1 0
> exec
| USE chk_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE chk_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_nt
< ok 0 0
> exec
| CREATE DATABASE chk_nt
< ok 1 0
> exec
| USE chk_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_nf
< ok 1 0
> exec
| CREATE DATABASE chk_nf
< ok 1 0
> exec
| USE chk_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_nt
< ok 1 0
> exec
| CREATE DATABASE chk_nt
< ok 1 0
> exec
| USE chk_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE chk_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_nf
< ok 1 0
> exec
| CREATE DATABASE chk_nf
< ok 1 0
> exec
| USE chk_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_nt
< ok 1 0
> exec
| CREATE DATABASE chk_nt
< ok 1 0
> exec
| USE chk_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 10))
< ok 0 0
> query
| SHOW CREATE TABLE chk_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck` CHECK ((`a` > 10))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_nf
< ok 1 0
> exec
| CREATE DATABASE chk_nf
< ok 1 0
> exec
| USE chk_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b))
< ok 0 0
> query
| SHOW CREATE TABLE chk_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck` CHECK ((`a` < `b`))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_nt
< ok 1 0
> exec
| CREATE DATABASE chk_nt
< ok 1 0
> exec
| USE chk_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE chk_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck` CHECK ((`a` < `b`))  NOT ENFORCED
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
