-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseIndexOptions default#1
> exec
| DROP DATABASE IF EXISTS obopt_storing
< ok 0 0
> exec
| CREATE DATABASE obopt_storing
< ok 1 0
> exec
| USE obopt_storing
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b, c))
< ok 0 0
> query
| SHOW CREATE TABLE obopt_storing.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   KEY `k` (`a`) STORING (`b`, `c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obopt_twice
< ok 0 0
> exec
| CREATE DATABASE obopt_twice
< ok 1 0
> exec
| USE obopt_twice
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b) STORING (c))
< ok 0 0
> query
| SHOW CREATE TABLE obopt_twice.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   KEY `k` (`a`) STORING (`b`, `c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obopt_case
< ok 0 0
> exec
| CREATE DATABASE obopt_case
< ok 1 0
> exec
| USE obopt_case
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, KEY k (a) STORING (B))
< ok 0 0
> query
| SHOW CREATE TABLE obopt_case.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   KEY `k` (`a`) STORING (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obopt_parallel
< ok 0 0
> exec
| CREATE DATABASE obopt_parallel
< ok 1 0
> exec
| USE obopt_parallel
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a) PARALLEL 3)
< ok 0 0
> query
| SHOW CREATE TABLE obopt_parallel.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obopt_noparallel
< ok 0 0
> exec
| CREATE DATABASE obopt_noparallel
< ok 1 0
> exec
| USE obopt_noparallel
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a) NOPARALLEL)
< ok 0 0
> query
| SHOW CREATE TABLE obopt_noparallel.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obopt_rowid
< ok 0 0
> exec
| CREATE DATABASE obopt_rowid
< ok 1 0
> exec
| USE obopt_rowid
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a) WITH ROWID)
< ok 0 0
> query
| SHOW CREATE TABLE obopt_rowid.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
