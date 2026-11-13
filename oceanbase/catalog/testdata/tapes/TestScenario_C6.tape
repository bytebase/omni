-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestScenario_C6 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) PARTITION BY HASH(id)
< ok 0 0
> query
| SELECT PARTITION_NAME FROM information_schema.PARTITIONS
|         WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|         ORDER BY PARTITION_ORDINAL_POSITION
< result
< column "PARTITION_NAME" "VARCHAR"
< row
: bytes
| p0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT, d DATE)
|             PARTITION BY RANGE(YEAR(d))
|             SUBPARTITION BY HASH(id)
|             (PARTITION p0 VALUES LESS THAN (2000),
|              PARTITION p1 VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SELECT PARTITION_NAME, SUBPARTITION_NAME
|             FROM information_schema.PARTITIONS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|             ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
< result
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
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
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) PARTITION BY KEY(id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key(`id`)
| (partition `p0`,
| partition `p1`,
| partition `p2`,
| partition `p3`)
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT) PARTITION BY KEY() PARTITIONS 4
< ok 0 0
> query
| SELECT COALESCE(PARTITION_EXPRESSION,'')
|             FROM information_schema.PARTITIONS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|             LIMIT 1
< result
< column "COALESCE(PARTITION_EXPRESSION,'')" "TEXT"
< row
: bytes
|
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) PARTITION BY LINEAR HASH(id) PARTITIONS 4
< error
| 1064 42000 …near 'LINEAR HASH(id) PARTITIONS 4' at line 1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) PARTITION BY RANGE(id) PARTITIONS 4
< error
| 1064 42000 Wrong number of partitions defined, mismatch with previous setting
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) PARTITION BY RANGE(id)
|             (PARTITION p0 VALUES LESS THAN MAXVALUE,
|              PARTITION p1 VALUES LESS THAN (100))
< error
| 1493 HY000 VALUES LESS THAN value must be strictly increasing for each partition
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (c INT) PARTITION BY LIST(c)
|              (PARTITION p0 VALUES IN (1,2), PARTITION p1 VALUES IN (3,4))
< ok 0 0
> query
| SELECT PARTITION_NAME FROM information_schema.PARTITIONS
|         WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|         ORDER BY PARTITION_ORDINAL_POSITION
< result
< column "PARTITION_NAME" "VARCHAR"
< row
: bytes
| p0
< row
: bytes
| p1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (c INT) PARTITION BY LIST(c)
|              (PARTITION p0 VALUES IN (1,2), PARTITION pd VALUES IN (DEFAULT))
< ok 0 0
> query
| SELECT PARTITION_NAME FROM information_schema.PARTITIONS
|         WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|         ORDER BY PARTITION_ORDINAL_POSITION
< result
< column "PARTITION_NAME" "VARCHAR"
< row
: bytes
| p0
< row
: bytes
| pd
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10), b VARCHAR(10))
|             PARTITION BY RANGE(CONCAT(a,b))
|             (PARTITION p0 VALUES LESS THAN ('m'),
|              PARTITION p1 VALUES LESS THAN MAXVALUE)
< error
| 1564 HY000 This partition function is not allowed
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t1 (ts TIMESTAMP) PARTITION BY RANGE(ts)
|             (PARTITION p0 VALUES LESS THAN (100))
< error
| 1659 HY000 Field 'ts' is of a not allowed type for this type of partitioning
> exec
| CREATE TABLE t2 (ts TIMESTAMP NOT NULL)
|             PARTITION BY RANGE(UNIX_TIMESTAMP(ts))
|             (PARTITION p0 VALUES LESS THAN (100),
|              PARTITION p1 VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, UNIQUE KEY (a))
|             PARTITION BY HASH(b) PARTITIONS 4
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT, d DATE)
|             ENGINE=InnoDB
|             PARTITION BY RANGE(YEAR(d))
|             SUBPARTITION BY HASH(id) SUBPARTITIONS 2
|             (PARTITION p0 VALUES LESS THAN (2000)
|                 (SUBPARTITION s0, SUBPARTITION s1),
|              PARTITION p1 VALUES LESS THAN MAXVALUE
|                 (SUBPARTITION s2, SUBPARTITION s3))
< ok 0 0
> query
| SELECT SUBPARTITION_NAME
|             FROM information_schema.PARTITIONS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|             ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
< result
< column "SUBPARTITION_NAME" "VARCHAR"
< row
: bytes
| s0
< row
: bytes
| s1
< row
: bytes
| s2
< row
: bytes
| s3
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 6
< ok 0 0
> exec
| ALTER TABLE t COALESCE PARTITION 2
< error
| 1064 42000 …near 'COALESCE PARTITION 2' at line 1
> query
| SELECT PARTITION_NAME FROM information_schema.PARTITIONS
|         WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|         ORDER BY PARTITION_ORDINAL_POSITION
< result
< column "PARTITION_NAME" "VARCHAR"
< row
: bytes
| p0
< row
: bytes
| p1
< row
: bytes
| p2
< row
: bytes
| p3
< row
: bytes
| p4
< row
: bytes
| p5
