-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainerPartitionKeyCoverage default#1
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE(a) (PARTITION `p0` VALUES LESS THAN (10), PARTITION `p1` VALUES LESS THAN MAXVALUE)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST(a) (PARTITION `p0` VALUES IN (1, 2))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE COLUMNS(`a`) (PARTITION `p0` VALUES LESS THAN (10))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST COLUMNS(`a`) (PARTITION `p0` VALUES IN (1))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY RANGE COLUMNS(`a`,`d`) (PARTITION `p0` VALUES LESS THAN (10, 10))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY KEY(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, PRIMARY KEY (a), UNIQUE KEY (b)) PARTITION BY KEY() PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
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
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY HASH(a + d) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LINEAR HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (a)) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION `p0` VALUES LESS THAN (10))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (a)) PARTITION BY RANGE(a) SUBPARTITION BY KEY(b) SUBPARTITIONS 2 (PARTITION `p0` VALUES LESS THAN (10))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (d), PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (d), PRIMARY KEY (b, a)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT PRIMARY KEY, b INT) PARTITION BY HASH(b) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT UNIQUE, b INT) PARTITION BY HASH(b) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b SERIAL) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ub (b), UNIQUE KEY ua (a)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT NOT NULL, b INT, UNIQUE KEY ub (b), UNIQUE KEY ua (a)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(5), b)) PARTITION BY KEY(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(19))) PARTITION BY KEY(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
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
| CREATE TABLE c (a VARCHAR(20), b VARCHAR(20), UNIQUE KEY (a, b(3))) PARTITION BY KEY(a) PARTITIONS 2
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
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
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
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (b)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (c)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
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
| CREATE TABLE t (a INT, b INT NOT NULL) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN d INT UNIQUE
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN d INT PRIMARY KEY
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b), ADD COLUMN f INT
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN f INT, ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (c), ADD PRIMARY KEY (b)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (c), ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b), ADD COLUMN a INT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE u (b), ADD UNIQUE KEY u (c)
< error
| 1061 42000 Duplicate key name 'u'
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b), DROP INDEX nope
< error
| 1091 42000 Can't DROP INDEX `nope`; check that it exists
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ui ON t (c)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
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
| CREATE TABLE t (a INT, b INT NOT NULL) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ub ON t (b)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
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
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) (PARTITION `p0` VALUES LESS THAN (10), PARTITION `p1` VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS(`a`) (PARTITION `p0` VALUES IN (1))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION `p0` VALUES LESS THAN (2000))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION `p0` VALUES LESS THAN (2000))
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
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY (a)) PARTITION BY KEY() PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION `p0` VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a(5))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
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
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, PRIMARY KEY (b))
< ok 0 0
> exec
| ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
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
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, UNIQUE KEY (b)) PARTITION BY HASH(b) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t RENAME COLUMN b TO b2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b2, a)
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, UNIQUE KEY (b)) PARTITION BY HASH(b) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t CHANGE b b2 INT
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
