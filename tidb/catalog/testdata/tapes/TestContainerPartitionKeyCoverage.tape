-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainerPartitionKeyCoverage default#1
> exec
| DROP DATABASE IF EXISTS omni_partition_coverage
< ok 0 0
> exec
| CREATE DATABASE omni_partition_coverage
< ok 0 0
> exec
| USE omni_partition_coverage
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A CLUSTERED INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, PRIMARY KEY (b) CLUSTERED) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A CLUSTERED INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, PRIMARY KEY (b) NONCLUSTERED) PARTITION BY HASH(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'PRIMARY', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b VARCHAR(10), PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A CLUSTERED INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST(a) (PARTITION p0 VALUES IN (1, 2))
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN (10))
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST COLUMNS(a) (PARTITION p0 VALUES IN (1))
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY RANGE COLUMNS(a, d) (PARTITION p0 VALUES LESS THAN (10, 10))
< error
| 8264 HY000 Global Index is needed for index 'a', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY KEY(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, PRIMARY KEY (a), UNIQUE KEY (b)) PARTITION BY KEY() PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT NOT NULL, b INT, UNIQUE KEY (a)) PARTITION BY KEY() PARTITIONS 2
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (d DATE, b INT, UNIQUE KEY (b)) PARTITION BY HASH(YEAR(d)) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY HASH(a + d) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'a', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LINEAR HASH(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (a)) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, PRIMARY KEY (b), UNIQUE KEY (d)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'd', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b VARCHAR(10), d INT, PRIMARY KEY (b), UNIQUE KEY (d)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A CLUSTERED INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b VARCHAR(10), d INT, UNIQUE KEY (d), PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'd', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, PRIMARY KEY (b) NONCLUSTERED, UNIQUE KEY ud (d)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'PRIMARY', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT UNIQUE, d INT, UNIQUE KEY ud (d)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'ud', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b VARCHAR(10) PRIMARY KEY, d INT, UNIQUE KEY ud (d)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'ud', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT PRIMARY KEY, b INT) PARTITION BY HASH(b) PARTITIONS 2
< error
| 1503 HY000 A CLUSTERED INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT UNIQUE, b INT) PARTITION BY HASH(b) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'a', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, KEY b (b), UNIQUE KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'b_2', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b SERIAL) PARTITION BY HASH(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(5), b)) PARTITION BY KEY(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'a', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(19))) PARTITION BY KEY(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'a', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(20))) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (A, b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 8264 HY000 Global Index is needed for index 'b_2', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE KEY u (a, b)
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT uc UNIQUE (c)
< error
| 8264 HY000 Global Index is needed for index 'uc', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (b)
< error
| 8264 HY000 Global Index is needed for index 'PRIMARY', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (b) NONCLUSTERED
< error
| 8264 HY000 Global Index is needed for index 'PRIMARY', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (c)
< error
| 8264 HY000 Global Index is needed for index 'PRIMARY', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (a, b)
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b), ADD COLUMN f INT
< error
| 8264 HY000 Global Index is needed for index 'b_2', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN f INT, ADD UNIQUE (b)
< error
| 8264 HY000 Global Index is needed for index 'b_2', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (c), ADD PRIMARY KEY (b)
< error
| 8264 HY000 Global Index is needed for index 'c', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (b), ADD UNIQUE (c)
< error
| 8264 HY000 Global Index is needed for index 'PRIMARY', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (c), ADD UNIQUE (b)
< error
| 8264 HY000 Global Index is needed for index 'c', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (c), ADD COLUMN a INT
< error
| 8264 HY000 Global Index is needed for index 'c', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN a INT, ADD UNIQUE (c)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (c), DROP INDEX nope
< error
| 8264 HY000 Global Index is needed for index 'c', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ui ON t (c)
< error
| 8264 HY000 Global Index is needed for index 'ui', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ui ON t (c, a)
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ue ON t ((c + 2))
< error
| 8264 HY000 Global Index is needed for index 'ue', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE INDEX ic ON t (c)
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS(a) (PARTITION p0 VALUES IN (1))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION p0 VALUES LESS THAN (2000))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION p0 VALUES LESS THAN (2000))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b, d)
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT) PARTITION BY KEY() PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY (a)) PARTITION BY KEY() PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a)
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a(5))
< error
| 8264 HY000 Global Index is needed for index 'a', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a(20))
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t REMOVE PARTITIONING
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, UNIQUE KEY (b))
< ok 0 0
> exec
| ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'b', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, PRIMARY KEY (b))
< ok 0 0
> exec
| ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A CLUSTERED INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (b)
< ok 0 0
> exec
| ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'PRIMARY', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, UNIQUE KEY (a, b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t PARTITION BY HASH(b) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b, c)
< ok 0 0
> exec
| ALTER TABLE t PARTITION BY HASH(c) PARTITIONS 2
< error
| 8264 HY000 Global Index is needed for index 'a', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption
> exec
| DROP DATABASE IF EXISTS omni_partition_coverage
< ok 0 0
