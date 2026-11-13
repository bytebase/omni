-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_PartitionKeyCoverage/add_column_primary_key default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN d INT PRIMARY KEY
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_column_then_unique default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN f INT, ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_column_unique default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN d INT UNIQUE
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_constraint_primary_key default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (c)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_constraint_unique default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT uc UNIQUE (c)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_functional_unique default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE ((b + 1))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_prefix_unique default#1
> exec
| CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a(5))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_prefix_unique_of_the_whole_column default#1
> exec
| CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a(20))
< ok 0 0

== TestContainer_PartitionKeyCoverage/add_primary_key default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (b)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_primary_key_covered default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (a, b)
< ok 0 0

== TestContainer_PartitionKeyCoverage/add_two_uniques default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (c), ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_unique default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_unique_after_remove_partitioning default#1
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t REMOVE PARTITIONING
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< ok 0 0

== TestContainer_PartitionKeyCoverage/add_unique_covered default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE KEY u (a, b)
< ok 0 0

== TestContainer_PartitionKeyCoverage/add_unique_on_a_NOT_NULL_column_is_promoted default#1
> exec
| CREATE TABLE t (a INT, b INT NOT NULL) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_unique_then_column default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b), ADD COLUMN f INT
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_unique_to_key default#1
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_unique_to_key_over_a_promoted_key default#1
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY (a)) PARTITION BY KEY() PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_unique_to_key_over_the_primary_key default#1
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT) PARTITION BY KEY() PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_unique_to_list_columns default#1
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS(a) (PARTITION p0 VALUES IN (1))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_unique_to_range default#1
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_unique_to_range_expression default#1
> exec
| CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION p0 VALUES LESS THAN (2000))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/add_unique_to_range_expression_covered default#1
> exec
| CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION p0 VALUES LESS THAN (2000))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b, d)
< ok 0 0

== TestContainer_PartitionKeyCoverage/add_unique_to_subpartitioned default#1
> exec
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_a_NOT_NULL_unique_key_is_promoted_past_a_nullable_one default#1
> exec
| CREATE TABLE c (a INT, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_a_nullable_unique_key_is_not_promoted default#1
> exec
| CREATE TABLE c (a INT NOT NULL, b INT, UNIQUE KEY ub (b), UNIQUE KEY ua (a)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_column_primary_key default#1
> exec
| CREATE TABLE c (a INT PRIMARY KEY, b INT) PARTITION BY HASH(b) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_column_unique default#1
> exec
| CREATE TABLE c (a INT UNIQUE, b INT) PARTITION BY HASH(b) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_covered default#1
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (A, b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0

== TestContainer_PartitionKeyCoverage/create_first_NOT_NULL_unique_key_is_promoted_and_covers default#1
> exec
| CREATE TABLE c (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_first_NOT_NULL_unique_key_is_promoted_and_leaves_a_column_out default#1
> exec
| CREATE TABLE c (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ub (b), UNIQUE KEY ua (a)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_functional_unique default#1
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY ((a + 1))) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_functional_unique_index default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ue ON t ((c + 2))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_hash_expression default#1
> exec
| CREATE TABLE c (d DATE, b INT, UNIQUE KEY (b)) PARTITION BY HASH(YEAR(d)) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_hash_of_two_columns default#1
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY HASH(a + d) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_hash_primary default#1
> exec
| CREATE TABLE c (a INT, b INT, PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_hash_unique default#1
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_hash_with_a_generated_invisible_primary_key default#1
> exec
| SET SESSION sql_generate_invisible_primary_key = ON
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1235 42000 This version of MySQL doesn't yet support 'generating invisible primary key for the partitioned tables'

== TestContainer_PartitionKeyCoverage/create_key default#1
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY KEY(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_key_over_a_promoted_unique_key default#1
> exec
| CREATE TABLE c (a INT NOT NULL, b INT, UNIQUE KEY (a)) PARTITION BY KEY() PARTITIONS 2
< ok 0 0

== TestContainer_PartitionKeyCoverage/create_key_over_the_primary_key default#1
> exec
| CREATE TABLE c (a INT, b INT, PRIMARY KEY (a), UNIQUE KEY (b)) PARTITION BY KEY() PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_linear_hash default#1
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LINEAR HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_list default#1
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST(a) (PARTITION p0 VALUES IN (1, 2))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_list_columns default#1
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST COLUMNS(a) (PARTITION p0 VALUES IN (1))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_plain_index default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE INDEX ic ON t (c)
< ok 0 0

== TestContainer_PartitionKeyCoverage/create_prefix default#1
> exec
| CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(5), b)) PARTITION BY KEY(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_prefix_of_the_whole_column default#1
> exec
| CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(20))) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0

== TestContainer_PartitionKeyCoverage/create_prefix_on_another_column default#1
> exec
| CREATE TABLE c (a VARCHAR(20), b VARCHAR(20), UNIQUE KEY (a, b(3))) PARTITION BY KEY(a) PARTITIONS 2
< ok 0 0

== TestContainer_PartitionKeyCoverage/create_prefix_short_by_one default#1
> exec
| CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(19))) PARTITION BY KEY(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_primary_key_covers default#1
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (d), PRIMARY KEY (b, a)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_primary_key_named_first default#1
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (d), PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_range default#1
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_range_columns default#1
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN (10))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_range_columns_pair default#1
> exec
| CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY RANGE COLUMNS(a, d) (PARTITION p0 VALUES LESS THAN (10, 10))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_serial_is_promoted default#1
> exec
| CREATE TABLE c (a INT, b SERIAL) PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_subpartition_hash default#1
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (a)) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_subpartition_key default#1
> exec
| CREATE TABLE c (a INT, b INT, UNIQUE KEY (a)) PARTITION BY RANGE(a) SUBPARTITION BY KEY(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_unique_index default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ui ON t (c)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_unique_index_covered default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ui ON t (c, a)
< ok 0 0

== TestContainer_PartitionKeyCoverage/create_unique_index_on_a_NOT_NULL_column_is_promoted default#1
> exec
| CREATE TABLE t (a INT, b INT NOT NULL) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| CREATE UNIQUE INDEX ub ON t (b)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/create_with_a_generated_invisible_primary_key default#1
> exec
| SET SESSION sql_generate_invisible_primary_key = ON
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10))
< error
| 1235 42000 This version of MySQL doesn't yet support 'generating invisible primary key for the partitioned tables'

== TestContainer_PartitionKeyCoverage/create_with_a_primary_key_when_invisible_keys_are_generated default#1
> exec
| SET SESSION sql_generate_invisible_primary_key = ON
< ok 0 0
> exec
| CREATE TABLE c (a INT, b INT, PRIMARY KEY (a)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0

== TestContainer_PartitionKeyCoverage/duplicate_column_comes_first default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b), ADD COLUMN a INT
< error
| 1060 42S21 Duplicate column name 'a'

== TestContainer_PartitionKeyCoverage/duplicate_key_name_comes_first default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE u (b), ADD UNIQUE KEY u (c)
< error
| 1061 42000 Duplicate key name 'u'

== TestContainer_PartitionKeyCoverage/missing_drop_comes_first default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (b), DROP INDEX nope
< error
| 1091 42000 Can't DROP 'nope'; check that column/key exists

== TestContainer_PartitionKeyCoverage/primary_key_named_before_unique default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (c), ADD PRIMARY KEY (b)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/repartition default#1
> exec
| CREATE TABLE t (a INT, b INT, UNIQUE KEY (b))
< ok 0 0
> exec
| ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/repartition_again default#1
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
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).

== TestContainer_PartitionKeyCoverage/repartition_over_a_primary_key default#1
> exec
| CREATE TABLE t (a INT, b INT, PRIMARY KEY (b))
< ok 0 0
> exec
| ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).
