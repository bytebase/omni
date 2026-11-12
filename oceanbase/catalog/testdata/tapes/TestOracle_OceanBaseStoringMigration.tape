-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseStoringMigration default#1
> exec
| DROP DATABASE IF EXISTS obstore_plain
< ok 0 0
> exec
| CREATE DATABASE obstore_plain
< ok 1 0
> exec
| USE obstore_plain
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE obstore_plain.t
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
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `obstore_plain`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `obstore_plain`.`t` ADD KEY `k` (`a`) STORING (`b`)
< ok 0 0
> query
| SHOW CREATE TABLE `obstore_plain`.`t`
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
|   KEY `k` (`a`) STORING (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obstore_two
< ok 0 0
> exec
| CREATE DATABASE obstore_two
< ok 1 0
> exec
| USE obstore_two
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE obstore_two.t
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
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `obstore_two`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `obstore_two`.`t` ADD KEY `k` (`a`) STORING (`b`,`c`)
< ok 0 0
> query
| SHOW CREATE TABLE `obstore_two`.`t`
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
| DROP DATABASE IF EXISTS obstore_comment
< ok 0 0
> exec
| CREATE DATABASE obstore_comment
< ok 1 0
> exec
| USE obstore_comment
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE obstore_comment.t
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
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `obstore_comment`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `obstore_comment`.`t` ADD KEY `k` (`a`) STORING (`b`) COMMENT 'x'
< ok 0 0
> query
| SHOW CREATE TABLE `obstore_comment`.`t`
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
|   KEY `k` (`a`) STORING (`b`) BLOCK_SIZE 16384 LOCAL COMMENT 'x'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obstore_global
< ok 0 0
> exec
| CREATE DATABASE obstore_global
< ok 1 0
> exec
| USE obstore_global
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE obstore_global.t
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
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `obstore_global`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `obstore_global`.`t` ADD KEY `k` (`a`) STORING (`b`) GLOBAL
< ok 0 0
> query
| SHOW CREATE TABLE `obstore_global`.`t`
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
|   KEY `k` (`a`) STORING (`b`) BLOCK_SIZE 16384 GLOBAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obstore_unique
< ok 0 0
> exec
| CREATE DATABASE obstore_unique
< ok 1 0
> exec
| USE obstore_unique
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE obstore_unique.t
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
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `obstore_unique`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `obstore_unique`.`t` ADD UNIQUE KEY `k` (`a`) STORING (`b`)
< ok 0 0
> query
| SHOW CREATE TABLE `obstore_unique`.`t`
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
|   UNIQUE KEY `k` (`a`) STORING (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
