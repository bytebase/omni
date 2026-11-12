-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseIndexScopeEcho default#1
> exec
| DROP DATABASE IF EXISTS obscope_unpartitioned
< ok 0 0
> exec
| CREATE DATABASE obscope_unpartitioned
< ok 1 0
> exec
| USE obscope_unpartitioned
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))
< ok 0 0
> query
| SHOW CREATE TABLE obscope_unpartitioned.t
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
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obscope_partitioned
< ok 0 0
> exec
| CREATE DATABASE obscope_partitioned
< ok 1 0
> exec
| USE obscope_partitioned
< ok 0 0
> exec
| CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), KEY ka (a)) PARTITION BY KEY(id) PARTITIONS 2
< ok 0 0
> query
| SHOW CREATE TABLE obscope_partitioned.t
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
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key(`id`)
| (partition `p0`,
| partition `p1`)
> exec
| DROP DATABASE IF EXISTS obscope_global_written
< ok 0 0
> exec
| CREATE DATABASE obscope_global_written
< ok 1 0
> exec
| USE obscope_global_written
< ok 0 0
> exec
| CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), KEY ka (a) GLOBAL) PARTITION BY KEY(id) PARTITIONS 2
< ok 0 0
> query
| SHOW CREATE TABLE obscope_global_written.t
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
|   KEY `ka` (`a`) BLOCK_SIZE 16384 GLOBAL
|  
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key(`id`)
| (partition `p0`,
| partition `p1`)
> exec
| DROP DATABASE IF EXISTS obscope_repeated_scope
< ok 0 0
> exec
| CREATE DATABASE obscope_repeated_scope
< ok 1 0
> exec
| USE obscope_repeated_scope
< ok 0 0
> exec
| CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), KEY ka (a) LOCAL GLOBAL, KEY kb (a) GLOBAL LOCAL) PARTITION BY KEY(id) PARTITIONS 2
< ok 0 0
> query
| SHOW CREATE TABLE obscope_repeated_scope.t
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
|   KEY `ka` (`a`) BLOCK_SIZE 16384 GLOBAL
|  ,
|   KEY `kb` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key(`id`)
| (partition `p0`,
| partition `p1`)
> exec
| DROP DATABASE IF EXISTS obscope_explicit_partitions
< ok 0 0
> exec
| CREATE DATABASE obscope_explicit_partitions
< ok 1 0
> exec
| USE obscope_explicit_partitions
< ok 0 0
> exec
| CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), KEY ka (a)) PARTITION BY KEY(id) (PARTITION p0, PARTITION p1)
< ok 0 0
> query
| SHOW CREATE TABLE obscope_explicit_partitions.t
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
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key(`id`)
| (partition `p0`,
| partition `p1`)
> exec
| DROP DATABASE IF EXISTS obscope_subpartition_template
< ok 0 0
> exec
| CREATE DATABASE obscope_subpartition_template
< ok 1 0
> exec
| USE obscope_subpartition_template
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY KEY(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1)
< ok 0 0
> query
| SHOW CREATE TABLE obscope_subpartition_template.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by hash(`a`) subpartition by key(`b`) subpartition template (
| subpartition `s0`,
| subpartition `s1`)
| (partition `p0`,
| partition `p1`)
