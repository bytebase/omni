-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseRangeListSubpartitions default#1
> exec
| DROP DATABASE IF EXISTS obrlsub_0
< ok 0 0
> exec
| CREATE DATABASE obrlsub_0
< ok 1 0
> exec
| USE obrlsub_0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE)) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_0.t
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
|  partition by range(`a`) subpartition by range(`b`) subpartition template (
| subpartition `s0` values less than (5),
| subpartition `s1` values less than (MAXVALUE))
| (partition `p0` values less than (10),
| partition `p1` values less than (20))
> query ^
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
|  partition by range(`a`) subpartition by range(`b`) subpartition template (
| subpartition `s0` values less than (5),
| subpartition `s1` values less than (MAXVALUE))
| (partition `p0` values less than (10),
| partition `p1` values less than (20))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_0
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0ss0
< row
: bytes
| p0
: bytes
| p0ss1
< row
: bytes
| p1
: bytes
| p1ss0
< row
: bytes
| p1
: bytes
| p1ss1
> exec
| DROP DATABASE IF EXISTS obrlsub_0
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_1
< ok 0 0
> exec
| CREATE DATABASE obrlsub_1
< ok 1 0
> exec
| USE obrlsub_1
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN (1, 2), SUBPARTITION s1 VALUES IN (DEFAULT)) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_1.t
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
|  partition by range(`a`) subpartition by list(`b`) subpartition template (
| subpartition `s0` values in (1,2),
| subpartition `s1` values in (DEFAULT))
| (partition `p0` values less than (10))
> query ^
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
|  partition by range(`a`) subpartition by list(`b`) subpartition template (
| subpartition `s0` values in (1,2),
| subpartition `s1` values in (DEFAULT))
| (partition `p0` values less than (10))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_1
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0ss0
< row
: bytes
| p0
: bytes
| p0ss1
> exec
| DROP DATABASE IF EXISTS obrlsub_1
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_2
< ok 0 0
> exec
| CREATE DATABASE obrlsub_2
< ok 1 0
> exec
| USE obrlsub_2
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES IN (1) (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (10)), PARTITION p1 VALUES IN (2) (SUBPARTITION s2 VALUES LESS THAN (100)))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_2.t
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
|  partition by list(`a`) subpartition by range(`b`)
| (partition `p0` values in (1) (
| subpartition `s0` values less than (5),
| subpartition `s1` values less than (10)),
| partition `p1` values in (2) (
| subpartition `s2` values less than (100)))
> query ^
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
|  partition by list(`a`) subpartition by range(`b`)
| (partition `p0` values in (1) (
| subpartition `s0` values less than (5),
| subpartition `s1` values less than (10)),
| partition `p1` values in (2) (
| subpartition `s2` values less than (100)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_2
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| s0
< row
: bytes
| p0
: bytes
| s1
< row
: bytes
| p1
: bytes
| s2
> exec
| DROP DATABASE IF EXISTS obrlsub_2
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_3
< ok 0 0
> exec
| CREATE DATABASE obrlsub_3
< ok 1 0
> exec
| USE obrlsub_3
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(10), c DATE) PARTITION BY RANGE COLUMNS(a) SUBPARTITION BY RANGE COLUMNS(b, c) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN ('m', '2020-01-01'), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE, MAXVALUE)) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_3.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL,
|   `c` date DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range columns(`a`) subpartition by range columns(`b`, `c`) subpartition template (
| subpartition `s0` values less than ('m','2020-01-01'),
| subpartition `s1` values less than (MAXVALUE,MAXVALUE))
| (partition `p0` values less than (10))
> query ^
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL,
|   `c` date DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range columns(`a`) subpartition by range columns(`b`, `c`) subpartition template (
| subpartition `s0` values less than ('m','2020-01-01'),
| subpartition `s1` values less than (MAXVALUE,MAXVALUE))
| (partition `p0` values less than (10))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_3
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0ss0
< row
: bytes
| p0
: bytes
| p0ss1
> exec
| DROP DATABASE IF EXISTS obrlsub_3
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_4
< ok 0 0
> exec
| CREATE DATABASE obrlsub_4
< ok 1 0
> exec
| USE obrlsub_4
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(10), c INT) PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY LIST COLUMNS(b, c) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN (('x', 1), ('y', 2)), SUBPARTITION s1 VALUES IN (DEFAULT))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by hash(`a`) subpartition by list columns(`b`, `c`) subpartition template (
| subpartition `s0` values in (('x',1),('y',2)),
| subpartition `s1` values in (DEFAULT))
| (partition `p0`,
| partition `p1`)
> query ^
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by hash(`a`) subpartition by list columns(`b`, `c`) subpartition template (
| subpartition `s0` values in (('x',1),('y',2)),
| subpartition `s1` values in (DEFAULT))
| (partition `p0`,
| partition `p1`)
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_4
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0ss0
< row
: bytes
| p0
: bytes
| p0ss1
< row
: bytes
| p1
: bytes
| p1ss0
< row
: bytes
| p1
: bytes
| p1ss1
> exec
| DROP DATABASE IF EXISTS obrlsub_4
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_5
< ok 0 0
> exec
| CREATE DATABASE obrlsub_5
< ok 1 0
> exec
| USE obrlsub_5
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY KEY(a) SUBPARTITION BY LIST(b) (PARTITION p0 (SUBPARTITION s0 VALUES IN (1), SUBPARTITION s1 VALUES IN (2)), PARTITION p1 (SUBPARTITION s2 VALUES IN (3)))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_5.t
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
|  partition by key(`a`) subpartition by list(`b`)
| (partition `p0` (
| subpartition `s0` values in (1),
| subpartition `s1` values in (2)),
| partition `p1` (
| subpartition `s2` values in (3)))
> query ^
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
|  partition by key(`a`) subpartition by list(`b`)
| (partition `p0` (
| subpartition `s0` values in (1),
| subpartition `s1` values in (2)),
| partition `p1` (
| subpartition `s2` values in (3)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_5
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| s0
< row
: bytes
| p0
: bytes
| s1
< row
: bytes
| p1
: bytes
| s2
> exec
| DROP DATABASE IF EXISTS obrlsub_5
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_6
< ok 0 0
> exec
| CREATE DATABASE obrlsub_6
< ok 1 0
> exec
| USE obrlsub_6
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b + 1) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5)) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_6.t
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
|  partition by range(`a`) subpartition by range((`b` + 1)) subpartition template (
| subpartition `s0` values less than (5))
| (partition `p0` values less than (10))
> query ^
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
|  partition by range(`a`) subpartition by range((`b` + 1)) subpartition template (
| subpartition `s0` values less than (5))
| (partition `p0` values less than (10))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_6
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0ss0
> exec
| DROP DATABASE IF EXISTS obrlsub_6
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_7
< ok 0 0
> exec
| CREATE DATABASE obrlsub_7
< ok 1 0
> exec
| USE obrlsub_7
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS(a) SUBPARTITION BY RANGE COLUMNS(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE)) (PARTITION p0 VALUES IN (1, 2), PARTITION p1 VALUES IN (3))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_7.t
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
|  partition by list columns(`a`) subpartition by range columns(`b`) subpartition template (
| subpartition `s0` values less than (5),
| subpartition `s1` values less than (MAXVALUE))
| (partition `p0` values in (1,2),
| partition `p1` values in (3))
> query ^
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
|  partition by list columns(`a`) subpartition by range columns(`b`) subpartition template (
| subpartition `s0` values less than (5),
| subpartition `s1` values less than (MAXVALUE))
| (partition `p0` values in (1,2),
| partition `p1` values in (3))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_7
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0ss0
< row
: bytes
| p0
: bytes
| p0ss1
< row
: bytes
| p1
: bytes
| p1ss0
< row
: bytes
| p1
: bytes
| p1ss1
> exec
| DROP DATABASE IF EXISTS obrlsub_7
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_8
< ok 0 0
> exec
| CREATE DATABASE obrlsub_8
< ok 1 0
> exec
| USE obrlsub_8
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5)), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_8.t
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
|  partition by range(`a`) subpartition by range(`b`)
| (partition `p0` values less than (10) (
| subpartition `s0` values less than (5)),
| partition `p1` values less than (20) (
| subpartition `P8192` values less than (MAXVALUE)))
> query ^
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
|  partition by range(`a`) subpartition by range(`b`)
| (partition `p0` values less than (10) (
| subpartition `s0` values less than (5)),
| partition `p1` values less than (20) (
| subpartition `P8192` values less than (MAXVALUE)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_8
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| s0
< row
: bytes
| p1
: bytes
| P8192
> exec
| DROP DATABASE IF EXISTS obrlsub_8
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_9
< ok 0 0
> exec
| CREATE DATABASE obrlsub_9
< ok 1 0
> exec
| USE obrlsub_9
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_9.t
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
|  partition by range(`a`) subpartition by range(`b`)
| (partition `p0` values less than (10) (
| subpartition `P8192` values less than (MAXVALUE)))
> query ^
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
|  partition by range(`a`) subpartition by range(`b`)
| (partition `p0` values less than (10) (
| subpartition `P8192` values less than (MAXVALUE)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_9
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| P8192
> exec
| DROP DATABASE IF EXISTS obrlsub_9
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_10
< ok 0 0
> exec
| CREATE DATABASE obrlsub_10
< ok 1 0
> exec
| USE obrlsub_10
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES IN (1)), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_10.t
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
|  partition by range(`a`) subpartition by list(`b`)
| (partition `p0` values less than (10) (
| subpartition `s0` values in (1)),
| partition `p1` values less than (20) (
| subpartition `P8192` values in (DEFAULT)))
> query ^
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
|  partition by range(`a`) subpartition by list(`b`)
| (partition `p0` values less than (10) (
| subpartition `s0` values in (1)),
| partition `p1` values less than (20) (
| subpartition `P8192` values in (DEFAULT)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_10
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| s0
< row
: bytes
| p1
: bytes
| P8192
> exec
| DROP DATABASE IF EXISTS obrlsub_10
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_11
< ok 0 0
> exec
| CREATE DATABASE obrlsub_11
< ok 1 0
> exec
| USE obrlsub_11
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST COLUMNS(b) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_11.t
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
|  partition by range(`a`) subpartition by list columns(`b`)
| (partition `p0` values less than (10) (
| subpartition `P8192` values in (DEFAULT)))
> query ^
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
|  partition by range(`a`) subpartition by list columns(`b`)
| (partition `p0` values less than (10) (
| subpartition `P8192` values in (DEFAULT)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_11
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| P8192
> exec
| DROP DATABASE IF EXISTS obrlsub_11
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_12
< ok 0 0
> exec
| CREATE DATABASE obrlsub_12
< ok 1 0
> exec
| USE obrlsub_12
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) SUBPARTITION BY RANGE(b) (PARTITION p0 (SUBPARTITION s0 VALUES LESS THAN (5)), PARTITION p1)
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_12.t
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
|  partition by hash(`a`) subpartition by range(`b`)
| (partition `p0` (
| subpartition `s0` values less than (5)),
| partition `p1` (
| subpartition `P8192` values less than (MAXVALUE)))
> query ^
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
|  partition by hash(`a`) subpartition by range(`b`)
| (partition `p0` (
| subpartition `s0` values less than (5)),
| partition `p1` (
| subpartition `P8192` values less than (MAXVALUE)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_12
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| s0
< row
: bytes
| p1
: bytes
| P8192
> exec
| DROP DATABASE IF EXISTS obrlsub_12
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_13
< ok 0 0
> exec
| CREATE DATABASE obrlsub_13
< ok 1 0
> exec
| USE obrlsub_13
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY LIST(b)
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_13.t
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
|  partition by hash(`a`) subpartition by list(`b`)
| (partition `p0` (
| subpartition `P8192` values in (DEFAULT)),
| partition `p1` (
| subpartition `P8193` values in (DEFAULT)))
> query ^
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
|  partition by hash(`a`) subpartition by list(`b`)
| (partition `p0` (
| subpartition `P8192` values in (DEFAULT)),
| partition `p1` (
| subpartition `P8193` values in (DEFAULT)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_13
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| P8192
< row
: bytes
| p1
: bytes
| P8193
> exec
| DROP DATABASE IF EXISTS obrlsub_13
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_14
< ok 0 0
> exec
| CREATE DATABASE obrlsub_14
< ok 1 0
> exec
| USE obrlsub_14
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN MAXVALUE) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_14.t
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
|  partition by range(`a`) subpartition by range(`b`) subpartition template (
| subpartition `s0` values less than (5),
| subpartition `s1` values less than (MAXVALUE))
| (partition `p0` values less than (10))
> query ^
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
|  partition by range(`a`) subpartition by range(`b`) subpartition template (
| subpartition `s0` values less than (5),
| subpartition `s1` values less than (MAXVALUE))
| (partition `p0` values less than (10))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_14
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0ss0
< row
: bytes
| p0
: bytes
| p0ss1
> exec
| DROP DATABASE IF EXISTS obrlsub_14
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_15
< ok 0 0
> exec
| CREATE DATABASE obrlsub_15
< ok 1 0
> exec
| USE obrlsub_15
< ok 0 0
> exec
| CREATE TABLE t (a INT, b DATE) PARTITION BY RANGE(a) SUBPARTITION BY RANGE COLUMNS(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN ('2020-01-01'), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE)) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_15.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` date DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(`a`) subpartition by range columns(`b`) subpartition template (
| subpartition `s0` values less than ('2020-01-01'),
| subpartition `s1` values less than (MAXVALUE))
| (partition `p0` values less than (10))
> query ^
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` date DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(`a`) subpartition by range columns(`b`) subpartition template (
| subpartition `s0` values less than ('2020-01-01'),
| subpartition `s1` values less than (MAXVALUE))
| (partition `p0` values less than (10))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_15
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0ss0
< row
: bytes
| p0
: bytes
| p0ss1
> exec
| DROP DATABASE IF EXISTS obrlsub_15
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_16
< ok 0 0
> exec
| CREATE DATABASE obrlsub_16
< ok 1 0
> exec
| USE obrlsub_16
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b + 1) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5 + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_16.t
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
|  partition by range(`a`) subpartition by range((`b` + 1))
| (partition `p0` values less than (10) (
| subpartition `s0` values less than (6)))
> query ^
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
|  partition by range(`a`) subpartition by range((`b` + 1))
| (partition `p0` values less than (10) (
| subpartition `s0` values less than (6)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_16
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| s0
> exec
| DROP DATABASE IF EXISTS obrlsub_16
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_17
< ok 0 0
> exec
| CREATE DATABASE obrlsub_17
< ok 1 0
> exec
| USE obrlsub_17
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST(a) SUBPARTITION BY LIST(b) (PARTITION p0 VALUES IN (1, 2) (SUBPARTITION s0 VALUES IN (1, 2), SUBPARTITION s1 VALUES IN (3)))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_17.t
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
|  partition by list(`a`) subpartition by list(`b`)
| (partition `p0` values in (1,2) (
| subpartition `s0` values in (1,2),
| subpartition `s1` values in (3)))
> query ^
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
|  partition by list(`a`) subpartition by list(`b`)
| (partition `p0` values in (1,2) (
| subpartition `s0` values in (1,2),
| subpartition `s1` values in (3)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_17
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| s0
< row
: bytes
| p0
: bytes
| s1
> exec
| DROP DATABASE IF EXISTS obrlsub_17
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_18
< ok 0 0
> exec
| CREATE DATABASE obrlsub_18
< ok 1 0
> exec
| USE obrlsub_18
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(10), c INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST COLUMNS(b, c) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES IN (('x', 1), ('y', 2)), SUBPARTITION s1 VALUES IN (DEFAULT)))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_18.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(`a`) subpartition by list columns(`b`, `c`)
| (partition `p0` values less than (10) (
| subpartition `s0` values in (('x',1),('y',2)),
| subpartition `s1` values in (DEFAULT)))
> query ^
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(`a`) subpartition by list columns(`b`, `c`)
| (partition `p0` values less than (10) (
| subpartition `s0` values in (('x',1),('y',2)),
| subpartition `s1` values in (DEFAULT)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_18
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| s0
< row
: bytes
| p0
: bytes
| s1
> exec
| DROP DATABASE IF EXISTS obrlsub_18
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_19
< ok 0 0
> exec
| CREATE DATABASE obrlsub_19
< ok 1 0
> exec
| USE obrlsub_19
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE COLUMNS(b, c) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5, MAXVALUE)) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_19.t
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
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(`a`) subpartition by range columns(`b`, `c`) subpartition template (
| subpartition `s0` values less than (5,MAXVALUE))
| (partition `p0` values less than (10))
> query ^
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
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(`a`) subpartition by range columns(`b`, `c`) subpartition template (
| subpartition `s0` values less than (5,MAXVALUE))
| (partition `p0` values less than (10))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_19
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0ss0
> exec
| DROP DATABASE IF EXISTS obrlsub_19
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlsub_20
< ok 0 0
> exec
| CREATE DATABASE obrlsub_20
< ok 1 0
> exec
| USE obrlsub_20
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN (5, MAXVALUE))
< ok 0 0
> query
| SHOW CREATE TABLE obrlsub_20.t
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
|  partition by range columns(`a`, `b`)
| (partition `p0` values less than (5,MAXVALUE))
> query ^
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
|  partition by range columns(`a`, `b`)
| (partition `p0` values less than (5,MAXVALUE))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrlsub_20
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
|
> exec
| DROP DATABASE IF EXISTS obrlsub_20
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrladd_0
< ok 0 0
> exec
| CREATE DATABASE obrladd_0
< ok 1 0
> exec
| USE obrladd_0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE)) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obrladd_0.t
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
|  partition by range(`a`) subpartition by range(`b`) subpartition template (
| subpartition `s0` values less than (5),
| subpartition `s1` values less than (MAXVALUE))
| (partition `p0` values less than (10))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))
< ok 0 0
> query
| SHOW CREATE TABLE obrladd_0.t
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
|  partition by range(`a`) subpartition by range(`b`) subpartition template (
| subpartition `s0` values less than (5),
| subpartition `s1` values less than (MAXVALUE))
| (partition `p0` values less than (10),
| partition `p2` values less than (30))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrladd_0
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0ss0
< row
: bytes
| p0
: bytes
| p0ss1
< row
: bytes
| p2
: bytes
| p2ss0
< row
: bytes
| p2
: bytes
| p2ss1
> exec
| DROP DATABASE IF EXISTS obrladd_0
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrladd_1
< ok 0 0
> exec
| CREATE DATABASE obrladd_1
< ok 1 0
> exec
| USE obrladd_1
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN (1, 2), SUBPARTITION s1 VALUES IN (DEFAULT)) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obrladd_1.t
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
|  partition by range(`a`) subpartition by list(`b`) subpartition template (
| subpartition `s0` values in (1,2),
| subpartition `s1` values in (DEFAULT))
| (partition `p0` values less than (10))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION z0 VALUES IN (5), SUBPARTITION z1 VALUES IN (DEFAULT)))
< ok 0 0
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))
< ok 0 0
> query
| SHOW CREATE TABLE obrladd_1.t
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
|  partition by range(`a`) subpartition by list(`b`)
| (partition `p0` values less than (10) (
| subpartition `p0ss0` values in (1,2),
| subpartition `p0ss1` values in (DEFAULT)),
| partition `p1` values less than (20) (
| subpartition `z0` values in (5),
| subpartition `z1` values in (DEFAULT)),
| partition `p2` values less than (30) (
| subpartition `p2ss0` values in (1,2),
| subpartition `p2ss1` values in (DEFAULT)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrladd_1
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0ss0
< row
: bytes
| p0
: bytes
| p0ss1
< row
: bytes
| p1
: bytes
| z0
< row
: bytes
| p1
: bytes
| z1
< row
: bytes
| p2
: bytes
| p2ss0
< row
: bytes
| p2
: bytes
| p2ss1
> exec
| DROP DATABASE IF EXISTS obrladd_1
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrladd_2
< ok 0 0
> exec
| CREATE DATABASE obrladd_2
< ok 1 0
> exec
| USE obrladd_2
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE)))
< ok 0 0
> query
| SHOW CREATE TABLE obrladd_2.t
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
|  partition by range(`a`) subpartition by range(`b`)
| (partition `p0` values less than (10) (
| subpartition `s0` values less than (5),
| subpartition `s1` values less than (MAXVALUE)))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SHOW CREATE TABLE obrladd_2.t
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
|  partition by range(`a`) subpartition by range(`b`)
| (partition `p0` values less than (10) (
| subpartition `s0` values less than (5),
| subpartition `s1` values less than (MAXVALUE)),
| partition `p1` values less than (20) (
| subpartition `p1sp0` values less than (MAXVALUE)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrladd_2
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| s0
< row
: bytes
| p0
: bytes
| s1
< row
: bytes
| p1
: bytes
| p1sp0
> exec
| DROP DATABASE IF EXISTS obrladd_2
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrladd_3
< ok 0 0
> exec
| CREATE DATABASE obrladd_3
< ok 1 0
> exec
| USE obrladd_3
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES IN (1) (SUBPARTITION s0 VALUES LESS THAN (5)))
< ok 0 0
> query
| SHOW CREATE TABLE obrladd_3.t
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
|  partition by list(`a`) subpartition by range(`b`)
| (partition `p0` values in (1) (
| subpartition `s0` values less than (5)))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES IN (3) (SUBPARTITION s3 VALUES LESS THAN (7)))
< ok 0 0
> query
| SHOW CREATE TABLE obrladd_3.t
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
|  partition by list(`a`) subpartition by range(`b`)
| (partition `p0` values in (1) (
| subpartition `s0` values less than (5)),
| partition `p2` values in (3) (
| subpartition `s3` values less than (7)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrladd_3
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| s0
< row
: bytes
| p2
: bytes
| s3
> exec
| DROP DATABASE IF EXISTS obrladd_3
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrladd_4
< ok 0 0
> exec
| CREATE DATABASE obrladd_4
< ok 1 0
> exec
| USE obrladd_4
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES IN (1)), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SHOW CREATE TABLE obrladd_4.t
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
|  partition by range(`a`) subpartition by list(`b`)
| (partition `p0` values less than (10) (
| subpartition `s0` values in (1)),
| partition `p1` values less than (20) (
| subpartition `P8192` values in (DEFAULT)))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))
< ok 0 0
> query
| SHOW CREATE TABLE obrladd_4.t
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
|  partition by range(`a`) subpartition by list(`b`)
| (partition `p0` values less than (10) (
| subpartition `s0` values in (1)),
| partition `p1` values less than (20) (
| subpartition `P8192` values in (DEFAULT)),
| partition `p2` values less than (30) (
| subpartition `p2sp0` values in (DEFAULT)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrladd_4
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| s0
< row
: bytes
| p1
: bytes
| P8192
< row
: bytes
| p2
: bytes
| p2sp0
> exec
| DROP DATABASE IF EXISTS obrladd_4
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrladd_5
< ok 0 0
> exec
| CREATE DATABASE obrladd_5
< ok 1 0
> exec
| USE obrladd_5
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE COLUMNS(b) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obrladd_5.t
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
|  partition by range(`a`) subpartition by range columns(`b`)
| (partition `p0` values less than (10) (
| subpartition `P8192` values less than (MAXVALUE)))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SHOW CREATE TABLE obrladd_5.t
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
|  partition by range(`a`) subpartition by range columns(`b`)
| (partition `p0` values less than (10) (
| subpartition `P8192` values less than (MAXVALUE)),
| partition `p1` values less than (20) (
| subpartition `p1sp0` values less than (MAXVALUE)))
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obrladd_5
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| P8192
< row
: bytes
| p1
: bytes
| p1sp0
> exec
| DROP DATABASE IF EXISTS obrladd_5
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlkind_0
< ok 0 0
> exec
| CREATE DATABASE obrlkind_0
< ok 1 0
> exec
| USE obrlkind_0
< ok 0 0
> exec
| CREATE TABLE seed (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obrlkind_0.seed
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| seed
: bytes
| CREATE TABLE `seed` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES IN (5)))
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS obrlkind_0
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlkind_1
< ok 0 0
> exec
| CREATE DATABASE obrlkind_1
< ok 1 0
> exec
| USE obrlkind_1
< ok 0 0
> exec
| CREATE TABLE seed (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obrlkind_1.seed
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| seed
: bytes
| CREATE TABLE `seed` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0))
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS obrlkind_1
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlkind_2
< ok 0 0
> exec
| CREATE DATABASE obrlkind_2
< ok 1 0
> exec
| USE obrlkind_2
< ok 0 0
> exec
| CREATE TABLE seed (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obrlkind_2.seed
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| seed
: bytes
| CREATE TABLE `seed` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5)))
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS obrlkind_2
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlkind_3
< ok 0 0
> exec
| CREATE DATABASE obrlkind_3
< ok 1 0
> exec
| USE obrlkind_3
< ok 0 0
> exec
| CREATE TABLE seed (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obrlkind_3.seed
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| seed
: bytes
| CREATE TABLE `seed` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (3)))
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS obrlkind_3
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlkind_4
< ok 0 0
> exec
| CREATE DATABASE obrlkind_4
< ok 1 0
> exec
| USE obrlkind_4
< ok 0 0
> exec
| CREATE TABLE seed (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obrlkind_4.seed
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| seed
: bytes
| CREATE TABLE `seed` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0), PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION s1 VALUES LESS THAN (5)))
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS obrlkind_4
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlkind_5
< ok 0 0
> exec
| CREATE DATABASE obrlkind_5
< ok 1 0
> exec
| USE obrlkind_5
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5)))
< ok 0 0
> query
| SHOW CREATE TABLE obrlkind_5.t
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
|  partition by range(`a`) subpartition by range(`b`)
| (partition `p0` values less than (10) (
| subpartition `s0` values less than (5)))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION z0 VALUES IN (5)))
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS obrlkind_5
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obrlkind_6
< ok 0 0
> exec
| CREATE DATABASE obrlkind_6
< ok 1 0
> exec
| USE obrlkind_6
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0))
< ok 0 0
> query
| SHOW CREATE TABLE obrlkind_6.t
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
|  partition by range(`a`) subpartition by hash(`b`)
| (partition `p0` values less than (10) (
| subpartition `s0`))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION z0 VALUES LESS THAN (5)))
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS obrlkind_6
< ok 1 0
