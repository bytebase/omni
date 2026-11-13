-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseRowFormat default#1
> exec
| DROP DATABASE IF EXISTS obrf_unstated
< ok 0 0
> exec
| CREATE DATABASE obrf_unstated
< ok 1 0
> exec
| USE obrf_unstated
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE obrf_unstated.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obrf_default
< ok 0 0
> exec
| CREATE DATABASE obrf_default
< ok 1 0
> exec
| USE obrf_default
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=DEFAULT
< ok 0 0
> query
| SHOW CREATE TABLE obrf_default.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obrf_dynamic
< ok 0 0
> exec
| CREATE DATABASE obrf_dynamic
< ok 1 0
> exec
| USE obrf_dynamic
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=DYNAMIC
< ok 0 0
> query
| SHOW CREATE TABLE obrf_dynamic.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obrf_compact
< ok 0 0
> exec
| CREATE DATABASE obrf_compact
< ok 1 0
> exec
| USE obrf_compact
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPACT
< ok 0 0
> query
| SHOW CREATE TABLE obrf_compact.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPACT COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obrf_redundant
< ok 0 0
> exec
| CREATE DATABASE obrf_redundant
< ok 1 0
> exec
| USE obrf_redundant
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=REDUNDANT
< ok 0 0
> query
| SHOW CREATE TABLE obrf_redundant.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = REDUNDANT COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obrf_compressed
< ok 0 0
> exec
| CREATE DATABASE obrf_compressed
< ok 1 0
> exec
| USE obrf_compressed
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED
< ok 0 0
> query
| SHOW CREATE TABLE obrf_compressed.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPRESSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obrf_condensed
< ok 0 0
> exec
| CREATE DATABASE obrf_condensed
< ok 1 0
> exec
| USE obrf_condensed
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=CONDENSED
< ok 0 0
> query
| SHOW CREATE TABLE obrf_condensed.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = CONDENSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_rowformat_reset_compact_a000152c_to
< ok 0 0
> exec
| CREATE DATABASE pa_rowformat_reset_compact_a000152c_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_rowformat_reset_compact_a000152c_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_rowformat_reset_compact_a000152c_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_rowformat_reset_compact_a000152c_from
< ok 0 0
> exec
| CREATE DATABASE pa_rowformat_reset_compact_a000152c_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_rowformat_reset_compact_a000152c_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPACT
< ok 0 0
> query
| SHOW CREATE TABLE pa_rowformat_reset_compact_a000152c_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPACT COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_rowformat_reset_compact_a000152c
< ok 0 0
> exec
| CREATE DATABASE pa_rowformat_reset_compact_a000152c DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_rowformat_reset_compact_a000152c
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| ALTER TABLE `pa_rowformat_reset_compact_a000152c`.`t` ROW_FORMAT=DEFAULT
< ok 0 0
> query
| SHOW CREATE TABLE pa_rowformat_reset_compact_a000152c.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_rowformat_reset_compressed_9b75c244_to
< ok 0 0
> exec
| CREATE DATABASE pa_rowformat_reset_compressed_9b75c244_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_rowformat_reset_compressed_9b75c244_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_rowformat_reset_compressed_9b75c244_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_rowformat_reset_compressed_9b75c244_from
< ok 0 0
> exec
| CREATE DATABASE pa_rowformat_reset_compressed_9b75c244_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_rowformat_reset_compressed_9b75c244_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED
< ok 0 0
> query
| SHOW CREATE TABLE pa_rowformat_reset_compressed_9b75c244_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPRESSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_rowformat_reset_compressed_9b75c244
< ok 0 0
> exec
| CREATE DATABASE pa_rowformat_reset_compressed_9b75c244 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_rowformat_reset_compressed_9b75c244
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED
< ok 0 0
> exec
| ALTER TABLE `pa_rowformat_reset_compressed_9b75c244`.`t` ROW_FORMAT=DEFAULT
< ok 0 0
> query
| SHOW CREATE TABLE pa_rowformat_reset_compressed_9b75c244.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
