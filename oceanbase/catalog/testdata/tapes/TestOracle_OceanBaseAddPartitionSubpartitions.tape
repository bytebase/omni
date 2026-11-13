-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseAddPartitionSubpartitions default#1
> exec
| DROP DATABASE IF EXISTS obreorg
< ok 0 0
> exec
| CREATE DATABASE obreorg
< ok 1 0
> exec
| USE obreorg
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b)) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SHOW CREATE TABLE obreorg.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`a`, `b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(`a`) subpartition by hash(`b`) subpartition template (
| subpartition `s0`,
| subpartition `s1`)
| (partition `p0` values less than (10),
| partition `p1` values less than (20))
> exec
| ALTER TABLE t REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (15), PARTITION p1b VALUES LESS THAN (20))
< error
| 1235 0A000 split subpartition is not supported
> exec
| ALTER TABLE t REORGANIZE PARTITION p0, p1 INTO (PARTITION p01 VALUES LESS THAN (20))
< error
| 1235 0A000 split subpartition is not supported
> exec
| DROP DATABASE IF EXISTS obreorg
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_0
< ok 0 0
> exec
| CREATE DATABASE obaddpart_0
< ok 1 0
> exec
| USE obaddpart_0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_0.t
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
|  partition by range(`a`) subpartition by hash(`b`) subpartition template (
| subpartition `s0`,
| subpartition `s1`)
| (partition `p0` values less than (10))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_0
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
> query
| SHOW CREATE TABLE obaddpart_0.t
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
|  partition by range(`a`) subpartition by hash(`b`) subpartition template (
| subpartition `s0`,
| subpartition `s1`)
| (partition `p0` values less than (10),
| partition `p1` values less than (20))
> exec
| DROP DATABASE IF EXISTS obaddpart_0
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_1
< ok 0 0
> exec
| CREATE DATABASE obaddpart_1
< ok 1 0
> exec
| USE obaddpart_1
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_1.t
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
|  partition by range(`a`) subpartition by hash(`b`) subpartition template (
| subpartition `s0`,
| subpartition `s1`)
| (partition `p0` values less than (10))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20), PARTITION p2 VALUES LESS THAN (30))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_1
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
> query
| SHOW CREATE TABLE obaddpart_1.t
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
|  partition by range(`a`) subpartition by hash(`b`) subpartition template (
| subpartition `s0`,
| subpartition `s1`)
| (partition `p0` values less than (10),
| partition `p1` values less than (20),
| partition `p2` values less than (30))
> exec
| DROP DATABASE IF EXISTS obaddpart_1
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_2
< ok 0 0
> exec
| CREATE DATABASE obaddpart_2
< ok 1 0
> exec
| USE obaddpart_2
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_2.t
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
|  partition by range(`a`) subpartition by hash(`b`) subpartition template (
| subpartition `s0`,
| subpartition `s1`)
| (partition `p0` values less than (10))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION x0, SUBPARTITION x1, SUBPARTITION x2))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_2
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
| x0
< row
: bytes
| p1
: bytes
| x1
< row
: bytes
| p1
: bytes
| x2
> query
| SHOW CREATE TABLE obaddpart_2.t
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
| subpartition `p0ss0`,
| subpartition `p0ss1`),
| partition `p1` values less than (20) (
| subpartition `x0`,
| subpartition `x1`,
| subpartition `x2`))
> exec
| DROP DATABASE IF EXISTS obaddpart_2
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_3
< ok 0 0
> exec
| CREATE DATABASE obaddpart_3
< ok 1 0
> exec
| USE obaddpart_3
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_3.t
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
|  partition by range(`a`) subpartition by hash(`b`) subpartition template (
| subpartition `s0`,
| subpartition `s1`)
| (partition `p0` values less than (10))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION x0))
< ok 0 0
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_3
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
| x0
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
> query
| SHOW CREATE TABLE obaddpart_3.t
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
| subpartition `p0ss0`,
| subpartition `p0ss1`),
| partition `p1` values less than (20) (
| subpartition `x0`),
| partition `p2` values less than (30) (
| subpartition `p2ss0`,
| subpartition `p2ss1`))
> exec
| DROP DATABASE IF EXISTS obaddpart_3
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_4
< ok 0 0
> exec
| CREATE DATABASE obaddpart_4
< ok 1 0
> exec
| USE obaddpart_4
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST(a) SUBPARTITION BY KEY(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES IN (1, 2))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_4.t
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
|  partition by list(`a`) subpartition by key(`b`) subpartition template (
| subpartition `s0`,
| subpartition `s1`)
| (partition `p0` values in (1,2))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES IN (3))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_4
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
> query
| SHOW CREATE TABLE obaddpart_4.t
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
|  partition by list(`a`) subpartition by key(`b`) subpartition template (
| subpartition `s0`,
| subpartition `s1`)
| (partition `p0` values in (1,2),
| partition `p1` values in (3))
> exec
| DROP DATABASE IF EXISTS obaddpart_4
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_5
< ok 0 0
> exec
| CREATE DATABASE obaddpart_5
< ok 1 0
> exec
| USE obaddpart_5
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY KEY(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1)
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_5.t
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
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_5
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
> query
| SHOW CREATE TABLE obaddpart_5.t
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
> exec
| DROP DATABASE IF EXISTS obaddpart_5
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_6
< ok 0 0
> exec
| CREATE DATABASE obaddpart_6
< ok 1 0
> exec
| USE obaddpart_6
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_6.t
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
|  partition by range(`a`) subpartition by hash(`b`) subpartition template (
| subpartition `p0`,
| subpartition `p1`)
| (partition `p0` values less than (10))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_6
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0sp0
< row
: bytes
| p0
: bytes
| p0sp1
< row
: bytes
| p1
: bytes
| p1sp0
< row
: bytes
| p1
: bytes
| p1sp1
> query
| SHOW CREATE TABLE obaddpart_6.t
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
|  partition by range(`a`) subpartition by hash(`b`) subpartition template (
| subpartition `p0`,
| subpartition `p1`)
| (partition `p0` values less than (10),
| partition `p1` values less than (20))
> exec
| DROP DATABASE IF EXISTS obaddpart_6
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_7
< ok 0 0
> exec
| CREATE DATABASE obaddpart_7
< ok 1 0
> exec
| USE obaddpart_7
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_7.t
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
|  partition by range(`a`) subpartition by hash(`b`) subpartition template (
| subpartition `p0`,
| subpartition `p1`)
| (partition `p0` values less than (10))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION y0, SUBPARTITION y1, SUBPARTITION y2))
< ok 0 0
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_7
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0sp0
< row
: bytes
| p0
: bytes
| p0sp1
< row
: bytes
| p1
: bytes
| y0
< row
: bytes
| p1
: bytes
| y1
< row
: bytes
| p1
: bytes
| y2
< row
: bytes
| p2
: bytes
| p2sp0
< row
: bytes
| p2
: bytes
| p2sp1
> query
| SHOW CREATE TABLE obaddpart_7.t
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
| subpartition `p0sp0`,
| subpartition `p0sp1`),
| partition `p1` values less than (20) (
| subpartition `y0`,
| subpartition `y1`,
| subpartition `y2`),
| partition `p2` values less than (30) (
| subpartition `p2sp0`,
| subpartition `p2sp1`))
> exec
| DROP DATABASE IF EXISTS obaddpart_7
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_8
< ok 0 0
> exec
| CREATE DATABASE obaddpart_8
< ok 1 0
> exec
| USE obaddpart_8
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0, SUBPARTITION s1))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_8.t
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
| subpartition `s0`,
| subpartition `s1`))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30) (SUBPARTITION s2, SUBPARTITION s3))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_8
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
< row
: bytes
| p2
: bytes
| s2
< row
: bytes
| p2
: bytes
| s3
> query
| SHOW CREATE TABLE obaddpart_8.t
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
| subpartition `s0`,
| subpartition `s1`),
| partition `p1` values less than (20) (
| subpartition `p1sp0`),
| partition `p2` values less than (30) (
| subpartition `s2`,
| subpartition `s3`))
> exec
| DROP DATABASE IF EXISTS obaddpart_8
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_9
< ok 0 0
> exec
| CREATE DATABASE obaddpart_9
< ok 1 0
> exec
| USE obaddpart_9
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0, SUBPARTITION s1), PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION s2))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_9.t
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
| subpartition `s0`,
| subpartition `s1`),
| partition `p1` values less than (20) (
| subpartition `s2`))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_9
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
< row
: bytes
| p2
: bytes
| p2sp0
> query
| SHOW CREATE TABLE obaddpart_9.t
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
| subpartition `s0`,
| subpartition `s1`),
| partition `p1` values less than (20) (
| subpartition `s2`),
| partition `p2` values less than (30) (
| subpartition `p2sp0`))
> exec
| DROP DATABASE IF EXISTS obaddpart_9
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_10
< ok 0 0
> exec
| CREATE DATABASE obaddpart_10
< ok 1 0
> exec
| USE obaddpart_10
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_10.t
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
| subpartition `s0`),
| partition `p1` values less than (20) (
| subpartition `p1sp0`))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_10
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
| p1sp0
< row
: bytes
| p2
: bytes
| p2sp0
> query
| SHOW CREATE TABLE obaddpart_10.t
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
| subpartition `s0`),
| partition `p1` values less than (20) (
| subpartition `p1sp0`),
| partition `p2` values less than (30) (
| subpartition `p2sp0`))
> exec
| DROP DATABASE IF EXISTS obaddpart_10
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_11
< ok 0 0
> exec
| CREATE DATABASE obaddpart_11
< ok 1 0
> exec
| USE obaddpart_11
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0, SUBPARTITION s1, SUBPARTITION s2), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_11.t
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
| subpartition `s0`,
| subpartition `s1`,
| subpartition `s2`),
| partition `p1` values less than (20) (
| subpartition `p1sp0`,
| subpartition `p1sp1`))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_11
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
| p0
: bytes
| s2
< row
: bytes
| p1
: bytes
| p1sp0
< row
: bytes
| p1
: bytes
| p1sp1
< row
: bytes
| p2
: bytes
| p2sp0
> query
| SHOW CREATE TABLE obaddpart_11.t
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
| subpartition `s0`,
| subpartition `s1`,
| subpartition `s2`),
| partition `p1` values less than (20) (
| subpartition `p1sp0`,
| subpartition `p1sp1`),
| partition `p2` values less than (30) (
| subpartition `p2sp0`))
> exec
| DROP DATABASE IF EXISTS obaddpart_11
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_12
< ok 0 0
> exec
| CREATE DATABASE obaddpart_12
< ok 1 0
> exec
| USE obaddpart_12
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_12.t
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
|  partition by range(`a`) subpartition by hash(`b`) subpartition template (
| subpartition `p0`)
| (partition `p0` values less than (10))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_12
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0sp0
< row
: bytes
| p1
: bytes
| p1sp0
> query
| SHOW CREATE TABLE obaddpart_12.t
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
|  partition by range(`a`) subpartition by hash(`b`) subpartition template (
| subpartition `p0`)
| (partition `p0` values less than (10),
| partition `p1` values less than (20))
> exec
| DROP DATABASE IF EXISTS obaddpart_12
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obaddpart_13
< ok 0 0
> exec
| CREATE DATABASE obaddpart_13
< ok 1 0
> exec
| USE obaddpart_13
< ok 0 0
> exec
| CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SHOW CREATE TABLE obaddpart_13.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(`a`)
| (partition `p0` values less than (10))
> exec
| ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
: string
| obaddpart_13
< result
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(SUBPARTITION_NAME, '')" "VARCHAR"
< row
: bytes
| p0
: bytes
|
< row
: bytes
| p1
: bytes
|
> query
| SHOW CREATE TABLE obaddpart_13.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range(`a`)
| (partition `p0` values less than (10),
| partition `p1` values less than (20))
> exec
| DROP DATABASE IF EXISTS obaddpart_13
< ok 1 0
