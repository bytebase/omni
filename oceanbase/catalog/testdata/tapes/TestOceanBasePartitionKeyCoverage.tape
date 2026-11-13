-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOceanBasePartitionKeyCoverage default#1
> exec
| DROP DATABASE IF EXISTS pkc_0
< ok 0 0
> exec
| CREATE DATABASE pkc_0
< ok 1 0
> exec
| USE pkc_0
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_1
< ok 0 0
> exec
| CREATE DATABASE pkc_1
< ok 1 0
> exec
| USE pkc_1
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_2
< ok 0 0
> exec
| CREATE DATABASE pkc_2
< ok 1 0
> exec
| USE pkc_2
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_2
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_3
< ok 0 0
> exec
| CREATE DATABASE pkc_3
< ok 1 0
> exec
| USE pkc_3
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST(a) (PARTITION p0 VALUES IN (1, 2))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_3
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_4
< ok 0 0
> exec
| CREATE DATABASE pkc_4
< ok 1 0
> exec
| USE pkc_4
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN (10))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_4
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_5
< ok 0 0
> exec
| CREATE DATABASE pkc_5
< ok 1 0
> exec
| USE pkc_5
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST COLUMNS(a) (PARTITION p0 VALUES IN (1))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_5
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_6
< ok 0 0
> exec
| CREATE DATABASE pkc_6
< ok 1 0
> exec
| USE pkc_6
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY RANGE COLUMNS(a, d) (PARTITION p0 VALUES LESS THAN (10, 10))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_6
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_7
< ok 0 0
> exec
| CREATE DATABASE pkc_7
< ok 1 0
> exec
| USE pkc_7
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY KEY(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_7
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_8
< ok 0 0
> exec
| CREATE DATABASE pkc_8
< ok 1 0
> exec
| USE pkc_8
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, PRIMARY KEY (a), UNIQUE KEY (b)) PARTITION BY KEY() PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_8
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_9
< ok 0 0
> exec
| CREATE DATABASE pkc_9
< ok 1 0
> exec
| USE pkc_9
< ok 0 0
> exec
| CREATE TABLE c (d DATE, b INT, UNIQUE KEY (b)) PARTITION BY HASH(YEAR(d)) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_9
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_10
< ok 0 0
> exec
| CREATE DATABASE pkc_10
< ok 1 0
> exec
| USE pkc_10
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY HASH(a + d) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_10
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_11
< ok 0 0
> exec
| CREATE DATABASE pkc_11
< ok 1 0
> exec
| USE pkc_11
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (a)) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_11
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_12
< ok 0 0
> exec
| CREATE DATABASE pkc_12
< ok 1 0
> exec
| USE pkc_12
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (a)) PARTITION BY RANGE(a) SUBPARTITION BY KEY(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_12
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_13
< ok 0 0
> exec
| CREATE DATABASE pkc_13
< ok 1 0
> exec
| USE pkc_13
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (d), PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_13
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_14
< ok 0 0
> exec
| CREATE DATABASE pkc_14
< ok 1 0
> exec
| USE pkc_14
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (d), PRIMARY KEY (b, a)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_14
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_15
< ok 0 0
> exec
| CREATE DATABASE pkc_15
< ok 1 0
> exec
| USE pkc_15
< ok 0 0
> exec
| CREATE TABLE c (a INT PRIMARY KEY, b INT) PARTITION BY HASH(b) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_15
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_16
< ok 0 0
> exec
| CREATE DATABASE pkc_16
< ok 1 0
> exec
| USE pkc_16
< ok 0 0
> exec
| CREATE TABLE c (a INT UNIQUE, b INT) PARTITION BY HASH(b) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_16
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_17
< ok 0 0
> exec
| CREATE DATABASE pkc_17
< ok 1 0
> exec
| USE pkc_17
< ok 0 0
> exec
| CREATE TABLE c (a INT, b SERIAL) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_17
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_18
< ok 0 0
> exec
| CREATE DATABASE pkc_18
< ok 1 0
> exec
| USE pkc_18
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY ((b + 1))) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_18
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_19
< ok 0 0
> exec
| CREATE DATABASE pkc_19
< ok 1 0
> exec
| USE pkc_19
< ok 0 0
> exec
| CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(5), b)) PARTITION BY KEY(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_19
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_20
< ok 0 0
> exec
| CREATE DATABASE pkc_20
< ok 1 0
> exec
| USE pkc_20
< ok 0 0
> exec
| CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(20))) PARTITION BY KEY(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_20
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_21
< ok 0 0
> exec
| CREATE DATABASE pkc_21
< ok 1 0
> exec
| USE pkc_21
< ok 0 0
> exec
| CREATE TABLE c (a VARCHAR(20), b VARCHAR(20), UNIQUE KEY (a, b(3))) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_21
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_22
< ok 0 0
> exec
| CREATE DATABASE pkc_22
< ok 1 0
> exec
| USE pkc_22
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (A, b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_22
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_23
< ok 0 0
> exec
| CREATE DATABASE pkc_23
< ok 1 0
> exec
| USE pkc_23
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b) GLOBAL) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_23
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_24
< ok 0 0
> exec
| CREATE DATABASE pkc_24
< ok 1 0
> exec
| USE pkc_24
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b) LOCAL) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_24
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_25
< ok 0 0
> exec
| CREATE DATABASE pkc_25
< ok 1 0
> exec
| USE pkc_25
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, PRIMARY KEY (b), UNIQUE KEY (a) GLOBAL) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_25
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_26
< ok 0 0
> exec
| CREATE DATABASE pkc_26
< ok 1 0
> exec
| USE pkc_26
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (d) GLOBAL, UNIQUE KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_26
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_27
< ok 0 0
> exec
| CREATE DATABASE pkc_27
< ok 1 0
> exec
| USE pkc_27
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_27
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_28
< ok 0 0
> exec
| CREATE DATABASE pkc_28
< ok 1 0
> exec
| USE pkc_28
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE KEY u (a, b)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_28
< ok 3 0
> exec
| DROP DATABASE IF EXISTS pkc_29
< ok 0 0
> exec
| CREATE DATABASE pkc_29
< ok 1 0
> exec
| USE pkc_29
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE u (b) GLOBAL
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_29
< ok 3 0
> exec
| DROP DATABASE IF EXISTS pkc_30
< ok 0 0
> exec
| CREATE DATABASE pkc_30
< ok 1 0
> exec
| USE pkc_30
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE u (b) LOCAL
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_30
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_31
< ok 0 0
> exec
| CREATE DATABASE pkc_31
< ok 1 0
> exec
| USE pkc_31
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT uc UNIQUE (c)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_31
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_32
< ok 0 0
> exec
| CREATE DATABASE pkc_32
< ok 1 0
> exec
| USE pkc_32
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT uc UNIQUE (c) GLOBAL
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_32
< ok 3 0
> exec
| DROP DATABASE IF EXISTS pkc_33
< ok 0 0
> exec
| CREATE DATABASE pkc_33
< ok 1 0
> exec
| USE pkc_33
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (b)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_33
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_34
< ok 0 0
> exec
| CREATE DATABASE pkc_34
< ok 1 0
> exec
| USE pkc_34
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (b) GLOBAL
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_34
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_35
< ok 0 0
> exec
| CREATE DATABASE pkc_35
< ok 1 0
> exec
| USE pkc_35
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (c)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_35
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_36
< ok 0 0
> exec
| CREATE DATABASE pkc_36
< ok 1 0
> exec
| USE pkc_36
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (a, b)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_36
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_37
< ok 0 0
> exec
| CREATE DATABASE pkc_37
< ok 1 0
> exec
| USE pkc_37
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN d INT UNIQUE
< error
| 1503 HY000 A UNIQUE INDEX/PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_37
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_38
< ok 0 0
> exec
| CREATE DATABASE pkc_38
< ok 1 0
> exec
| USE pkc_38
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b), ADD COLUMN f INT
< error
| 1503 HY000 A UNIQUE INDEX/PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_38
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_39
< ok 0 0
> exec
| CREATE DATABASE pkc_39
< ok 1 0
> exec
| USE pkc_39
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN f INT, ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX/PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_39
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_40
< ok 0 0
> exec
| CREATE DATABASE pkc_40
< ok 1 0
> exec
| USE pkc_40
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (c), ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_40
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_41
< ok 0 0
> exec
| CREATE DATABASE pkc_41
< ok 1 0
> exec
| USE pkc_41
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE u10 (b) GLOBAL, ADD UNIQUE u11 (c)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_41
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_42
< ok 0 0
> exec
| CREATE DATABASE pkc_42
< ok 1 0
> exec
| USE pkc_42
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE ((b + 1))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_42
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_43
< ok 0 0
> exec
| CREATE DATABASE pkc_43
< ok 1 0
> exec
| USE pkc_43
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b), ADD COLUMN a INT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS pkc_43
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_44
< ok 0 0
> exec
| CREATE DATABASE pkc_44
< ok 1 0
> exec
| USE pkc_44
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b), DROP INDEX nope
< error
| 1091 42000 Can't DROP 'nope'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS pkc_44
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_45
< ok 0 0
> exec
| CREATE DATABASE pkc_45
< ok 1 0
> exec
| USE pkc_45
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ui ON t (c)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_45
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_46
< ok 0 0
> exec
| CREATE DATABASE pkc_46
< ok 1 0
> exec
| USE pkc_46
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ui ON t (c, a)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_46
< ok 3 0
> exec
| DROP DATABASE IF EXISTS pkc_47
< ok 0 0
> exec
| CREATE DATABASE pkc_47
< ok 1 0
> exec
| USE pkc_47
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ui ON t (c) GLOBAL
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_47
< ok 3 0
> exec
| DROP DATABASE IF EXISTS pkc_48
< ok 0 0
> exec
| CREATE DATABASE pkc_48
< ok 1 0
> exec
| USE pkc_48
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ui ON t (c) LOCAL
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_48
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_49
< ok 0 0
> exec
| CREATE DATABASE pkc_49
< ok 1 0
> exec
| USE pkc_49
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ue ON t ((c + 2))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_49
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_50
< ok 0 0
> exec
| CREATE DATABASE pkc_50
< ok 1 0
> exec
| USE pkc_50
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE INDEX ic ON t (c)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_50
< ok 3 0
> exec
| DROP DATABASE IF EXISTS pkc_51
< ok 0 0
> exec
| CREATE DATABASE pkc_51
< ok 1 0
> exec
| USE pkc_51
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_51
< ok 1 0
> exec
| DROP DATABASE IF EXISTS pkc_52
< ok 0 0
> exec
| CREATE DATABASE pkc_52
< ok 1 0
> exec
| USE pkc_52
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS(a) (PARTITION p0 VALUES IN (1))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_52
< ok 1 0
> exec
| DROP DATABASE IF EXISTS pkc_53
< ok 0 0
> exec
| CREATE DATABASE pkc_53
< ok 1 0
> exec
| USE pkc_53
< ok 0 0
> exec
| CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION p0 VALUES LESS THAN (2000))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_53
< ok 1 0
> exec
| DROP DATABASE IF EXISTS pkc_54
< ok 0 0
> exec
| CREATE DATABASE pkc_54
< ok 1 0
> exec
| USE pkc_54
< ok 0 0
> exec
| CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION p0 VALUES LESS THAN (2000))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b, d)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_54
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_55
< ok 0 0
> exec
| CREATE DATABASE pkc_55
< ok 1 0
> exec
| USE pkc_55
< ok 0 0
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT) PARTITION BY KEY() PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_55
< ok 1 0
> exec
| DROP DATABASE IF EXISTS pkc_56
< ok 0 0
> exec
| CREATE DATABASE pkc_56
< ok 1 0
> exec
| USE pkc_56
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_56
< ok 1 0
> exec
| DROP DATABASE IF EXISTS pkc_57
< ok 0 0
> exec
| CREATE DATABASE pkc_57
< ok 1 0
> exec
| USE pkc_57
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_57
< ok 1 0
> exec
| DROP DATABASE IF EXISTS pkc_58
< ok 0 0
> exec
| CREATE DATABASE pkc_58
< ok 1 0
> exec
| USE pkc_58
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a, b)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_58
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_59
< ok 0 0
> exec
| CREATE DATABASE pkc_59
< ok 1 0
> exec
| USE pkc_59
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a(5))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_59
< ok 1 0
> exec
| DROP DATABASE IF EXISTS pkc_60
< ok 0 0
> exec
| CREATE DATABASE pkc_60
< ok 1 0
> exec
| USE pkc_60
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a(20))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_60
< ok 1 0
> exec
| DROP DATABASE IF EXISTS pkc_61
< ok 0 0
> exec
| CREATE DATABASE pkc_61
< ok 1 0
> exec
| USE pkc_61
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, UNIQUE KEY (b))
< ok 0 0
> exec
| ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX/PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_61
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_62
< ok 0 0
> exec
| CREATE DATABASE pkc_62
< ok 1 0
> exec
| USE pkc_62
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, PRIMARY KEY (b))
< ok 0 0
> exec
| ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_62
< ok 1 0
> exec
| DROP DATABASE IF EXISTS pkc_63
< ok 0 0
> exec
| CREATE DATABASE pkc_63
< ok 1 0
> exec
| USE pkc_63
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, UNIQUE KEY (b) GLOBAL)
< ok 0 0
> exec
| ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_63
< ok 2 0
> exec
| DROP DATABASE IF EXISTS pkc_64
< ok 0 0
> exec
| CREATE DATABASE pkc_64
< ok 1 0
> exec
| USE pkc_64
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
| 1503 HY000 A UNIQUE INDEX/PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS pkc_64
< ok 3 0
> exec
| DROP DATABASE IF EXISTS pkc_65
< ok 0 0
> exec
| CREATE DATABASE pkc_65
< ok 1 0
> exec
| USE pkc_65
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, UNIQUE KEY (b)) PARTITION BY HASH(b) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t CHANGE b b2 INT
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b2, a)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pkc_65
< ok 3 0
> exec
| DROP DATABASE IF EXISTS pkc_66
< ok 0 0
> exec
| CREATE DATABASE pkc_66
< ok 1 0
> exec
| USE pkc_66
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
> exec
| DROP DATABASE IF EXISTS pkc_66
< ok 2 0
