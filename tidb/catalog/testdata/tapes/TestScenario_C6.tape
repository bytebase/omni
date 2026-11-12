-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScenario_C6 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
: null
< row
: bytes
| p1
: null
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) ENGINE=InnoDB PARTITION BY HASH(id) PARTITIONS 2
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY HASH (`id`) PARTITIONS 2
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`id`) PARTITIONS 4
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
| `id`
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) PARTITION BY LINEAR HASH(id) PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY HASH (`id`) PARTITIONS 4
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) PARTITION BY RANGE(id) PARTITIONS 4
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) PARTITION BY RANGE(id)
|             (PARTITION p0 VALUES LESS THAN MAXVALUE,
|              PARTITION p1 VALUES LESS THAN (100))
< error
| 1481 HY000 MAXVALUE can only be used in last partition definition
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
< row
: bytes
| 8.0.11-TiDB-v8.5.5
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
> query
| SELECT COALESCE(PARTITION_DESCRIPTION,'')
|             FROM information_schema.PARTITIONS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND PARTITION_NAME='pd'
< result
< column "COALESCE(PARTITION_DESCRIPTION,'')" "TEXT"
< row
: bytes
| DEFAULT
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, UNIQUE KEY (a))
|             PARTITION BY HASH(b) PARTITIONS 4
< error
| 8264 HY000 Global Index is needed for index 'a', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) PARTITION BY HASH(id)
|              (PARTITION p0 COMMENT='first' ENGINE=InnoDB,
|               PARTITION p1 COMMENT='second' ENGINE=InnoDB)
< ok 0 0
> query
| SELECT PARTITION_NAME, PARTITION_COMMENT
|             FROM information_schema.PARTITIONS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|             ORDER BY PARTITION_ORDINAL_POSITION
< result
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_COMMENT" "VARCHAR"
< row
: bytes
| p0
: bytes
| first
< row
: bytes
| p1
: bytes
| second
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT, d DATE)
|             ENGINE=InnoDB
|             PARTITION BY RANGE(YEAR(d))
|             SUBPARTITION BY HASH(id) SUBPARTITIONS 2
|             (PARTITION p0 VALUES LESS THAN (2000)
|                 (SUBPARTITION s0 COMMENT='sa', SUBPARTITION s1 COMMENT='sb'),
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
: null
< row
: null
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 3
< ok 0 0
> exec
| ALTER TABLE t ADD PARTITION PARTITIONS 2
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
< row
: bytes
| p2
< row
: bytes
| p3
< row
: bytes
| p4
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 6
< ok 0 0
> exec
| ALTER TABLE t COALESCE PARTITION 2
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
< row
: bytes
| p2
< row
: bytes
| p3
