-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBasePartitionNames/a_table_with_no_partitions default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a DROP PARTITION p
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a TRUNCATE PARTITION p
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a REORGANIZE PARTITION p INTO (PARTITION q VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a EXCHANGE PARTITION p WITH TABLE a WITHOUT VALIDATION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible

== TestOracle_OceanBasePartitionNames/names_in_a_partitioning default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY HASH(b + zz)
< error
| 1054 42S22 Unknown column 'zz' in 'partition function'
> exec
| ALTER TABLE a PARTITION BY RANGE COLUMNS(zz)
< error
| 1054 42S22 Unknown column 'zz' in 'partition function'
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY HASH(zz)
< error
| 1054 42S22 Unknown column 'zz' in 'partition function'
> exec
| ALTER TABLE a PARTITION BY HASH(c)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| ALTER TABLE a PARTITION BY HASH(c)(PARTITION p(SUBPARTITION s))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION p(SUBPARTITION s VALUES IN(zz)))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| ALTER TABLE a PARTITION BY KEY(zz)(PARTITION p(SUBPARTITION s))
< error
| 1054 42S22 Unknown column 'zz' in 'partition function'
> exec
| ALTER TABLE a PARTITION BY LIST(zz)(PARTITION p VALUES IN(c))
< error
| 1054 42S22 Unknown column 'zz' in 'partition function'
> exec
| ALTER TABLE a PARTITION BY LIST(b)(PARTITION p VALUES IN(c))
< error
| 1054 42S22 Unknown column 'c' in 'partition function'
> exec
| ALTER TABLE a PARTITION BY RANGE(b)(PARTITION p VALUES LESS THAN(c))
< error
| 1054 42S22 Unknown column 'c' in 'partition function'
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY) PARTITION BY HASH(zz)
< error
| 1054 42S22 Unknown column 'zz' in 'partition function'
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY) PARTITION BY LIST(b)(PARTITION p VALUES IN(c))
< error
| 1054 42S22 Unknown column 'c' in 'partition function'
