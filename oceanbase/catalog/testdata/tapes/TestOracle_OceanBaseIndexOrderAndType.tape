-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseIndexOrderAndType default#1
> exec
| DROP DATABASE IF EXISTS obidx_order_type
< ok 0 0
> exec
| CREATE DATABASE obidx_order_type
< ok 1 0
> exec
| USE obidx_order_type
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ab (a, b DESC), KEY ha (a) USING HASH, KEY bb (b) USING BTREE, UNIQUE KEY ub (b DESC))
< ok 0 0
> query
| SHOW CREATE TABLE obidx_order_type.t
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
|   UNIQUE KEY `ub` (`b`) BLOCK_SIZE 16384 LOCAL,
|   KEY `ab` (`a`, `b`) BLOCK_SIZE 16384 LOCAL,
|   KEY `ha` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `bb` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT INDEX_NAME, COLUMN_NAME, COLLATION, INDEX_TYPE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'obidx_order_type' AND TABLE_NAME = 't' ORDER BY BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: bytes
| BTREE
< row
: bytes
| ab
: bytes
| a
: bytes
| A
: bytes
| BTREE
< row
: bytes
| ab
: bytes
| b
: bytes
| A
: bytes
| BTREE
< row
: bytes
| bb
: bytes
| b
: bytes
| A
: bytes
| BTREE
< row
: bytes
| ha
: bytes
| a
: bytes
| A
: bytes
| HASH
< row
: bytes
| ub
: bytes
| b
: bytes
| A
: bytes
| BTREE
