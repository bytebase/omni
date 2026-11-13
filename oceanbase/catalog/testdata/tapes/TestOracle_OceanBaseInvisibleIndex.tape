-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseInvisibleIndex default#1
> exec
| DROP DATABASE IF EXISTS obidx_invisible
< ok 0 0
> exec
| CREATE DATABASE obidx_invisible
< ok 1 0
> exec
| USE obidx_invisible
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT, KEY in_create (a) INVISIBLE, KEY in_create_comment (b) /*!80000 INVISIBLE */)
< ok 0 0
> query
| SHOW CREATE TABLE obidx_invisible.t
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
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `in_create` (`a`) BLOCK_SIZE 16384 LOCAL /*!80000 INVISIBLE */,
|   KEY `in_create_comment` (`b`) BLOCK_SIZE 16384 LOCAL /*!80000 INVISIBLE */
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = ? AND SEQ_IN_INDEX = 1
: string
| obidx_invisible
: string
| in_create
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| NO
> query ^
: string
| obidx_invisible
: string
| in_create_comment
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| NO
> exec
| ALTER TABLE t ADD KEY in_alter (c) INVISIBLE
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = ? AND SEQ_IN_INDEX = 1
: string
| obidx_invisible
: string
| in_alter
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> exec
| ALTER TABLE t ADD KEY in_alter_comment (a, b) /*!80000 INVISIBLE */
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = ? AND SEQ_IN_INDEX = 1
: string
| obidx_invisible
: string
| in_alter_comment
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> exec
| ALTER TABLE t ALTER INDEX in_alter INVISIBLE
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = ? AND SEQ_IN_INDEX = 1
: string
| obidx_invisible
: string
| in_alter
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| NO
> exec
| ALTER TABLE t ALTER INDEX in_create VISIBLE
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = ? AND SEQ_IN_INDEX = 1
: string
| obidx_invisible
: string
| in_create
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
