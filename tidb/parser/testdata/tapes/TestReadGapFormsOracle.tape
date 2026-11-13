-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestReadGapFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| SELECT 'a' ILIKE 'b'
< ok 0 0
> exec ^
: string
| SELECT 'a' NOT ILIKE 'b'
< ok 0 0
> exec ^
: string
| SELECT 'a' ! ILIKE 'b'
< error
| 1064 42000 …line 1 column 12 near "! ILIKE 'b'" 
> exec ^
: string
| SELECT 'a' ILIKE 'b' ESCAPE 'x'
< ok 0 0
> exec ^
: string
| SELECT 'a' ILIKE 'b' ESCAPE 'xy'
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| SELECT 'a' ILIKE 'b' ESCAPE 'xy'
| ; SELECT FROM
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| SELECT 'a' ILIKE 'b' ESCAPE ''
< ok 0 0
> exec ^
: string
| SELECT 'a' ILIKE 'b' ESCAPE concat('c')
< error
| 1064 42000 …line 1 column 34 near "concat('c')" 
> exec ^
: string
| SELECT 'a' ilike 'b' ilike 'c'
< error
| 1064 42000 …line 1 column 26 near "ilike 'c'" 
> exec ^
: string
| SELECT 'a' ILIKE 'b' LIKE 'c'
< error
| 1064 42000 …line 1 column 25 near "LIKE 'c'" 
> exec ^
: string
| SELECT 'a' ILIKE 'b' + 1
< error
| 1064 42000 …line 1 column 22 near "+ 1" 
> exec ^
: string
| SELECT 'a' ILIKE ('b')
< ok 0 0
> exec ^
: string
| SELECT 'a' ILIKE
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| SELECT ILIKE 'a'
< error
| 1064 42000 …line 1 column 12 near "ILIKE 'a'" 
> exec ^
: string
| SELECT 'a' ILIKE ALL (SELECT 1)
< error
| 1064 42000 …line 1 column 20 near "ALL (SELECT 1)" 
> exec ^
: string
| SELECT 'a' ILIKE 'b' IS NULL
< ok 0 0
> exec ^
: string
| SELECT NOT 'a' ILIKE 'b'
< ok 0 0
> exec ^
: string
| SELECT 1 FROM t WHERE a ILIKE '%x%' OR b NOT ILIKE 'y'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a ILIKE '%x%' OR b NOT ILIKE 'y'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a ILIKE b FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a ILIKE b FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT ilike FROM t
< error
| 1064 42000 …line 1 column 12 near "ilike FROM t" 
> exec ^
: string
| SELECT `ilike` FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT `ilike` FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.ilike FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.ilike FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT ilike(1)
< error
| 1064 42000 …line 1 column 12 near "ilike(1)" 
> exec ^
: string
| SELECT `ilike`(1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ilike'
> exec ^
: string
| SELECT `ilike`(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE ilike (a INT)
< error
| 1064 42000 …line 1 column 18 near "ilike (a INT)" 
> exec ^
: string
| CREATE TABLE `ilike` (a INT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (ilike INT)
< error
| 1064 42000 …line 1 column 21 near "ilike INT)" 
> exec ^
: string
| SELECT 1 AS ilike
< error
| 1064 42000 …line 1 column 17 near "ilike" 
> exec ^
: string
| SELECT 1 ilike
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| SELECT 1 AS `ilike`
< ok 0 0
> exec ^
: string
| SHOW TABLES ILIKE 'a'
< error
| 1064 42000 …line 1 column 17 near "ILIKE 'a'" 
> exec ^
: string
| SHOW DATABASES ILIKE 'a'
< error
| 1064 42000 …line 1 column 20 near "ILIKE 'a'" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) NULL PARTITION
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) MAXVALUE PARTITION
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) NULL PARTITION MAXVALUE PARTITION
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) MAXVALUE PARTITION NULL PARTITION
< error
| 1064 42000 …line 1 column 141 near "NULL PARTITION" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) NULL PARTITION
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) MAXVALUE PARTITION
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1)
< error
| 1064 42000 …line 1 column 88 near "" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) LAST PARTITION LESS THAN (5)
< error
| 1064 42000 …line 1 column 63 near "LAST PARTITION LESS THAN (5)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) LAST PARTITION LESS THAN (5) FIRST PARTITION LESS THAN (1)
< error
| 1064 42000 …line 1 column 63 near "LAST PARTITION LESS THAN (5) FIRST PARTITION LESS THAN (1)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN 1 LAST PARTITION LESS THAN 5
< error
| 1064 42000 …line 1 column 86 near "1 LAST PARTITION LESS THAN 5" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1,2) LAST PARTITION LESS THAN (5)
< error
| 1064 42000 …line 1 column 88 near ",2) LAST PARTITION LESS THAN (5)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (MAXVALUE) LAST PARTITION LESS THAN (5)
< error
| 1064 42000 …line 1 column 94 near "MAXVALUE) LAST PARTITION LESS THAN (5)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL ()
< error
| 1064 42000 …line 1 column 57 near ")" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL 1
< error
| 1064 42000 …line 1 column 56 near "1" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL
< error
| 1064 42000 …line 1 column 54 near "" 
> exec ^
: string
| CREATE TABLE t (a DATETIME) PARTITION BY RANGE (TO_DAYS(a)) INTERVAL (1 DAY) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 MONTH) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2021-01-01')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 FOO)
< error
| 1064 42000 …line 1 column 74 near "FOO)" 
> exec ^
: string
| CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 YEAR_MONTH)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME) PARTITION BY RANGE COLUMNS (a) INTERVAL (1 SQL_TSI_DAY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE FIELDS (a) INTERVAL (1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1+1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1,2)
< error
| 1064 42000 …line 1 column 58 near ",2)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) PARTITIONS 3
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) (PARTITION p VALUES LESS THAN (1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) (PARTITION p VALUES LESS THAN (1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) INTERVAL (1) SUBPARTITION BY HASH (b)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) SUBPARTITION BY HASH (b)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) UPDATE INDEXES (i GLOBAL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5) UPDATE INDEXES (i GLOBAL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) INTERVAL (1)
< error
| 1064 42000 …line 1 column 53 near "INTERVAL (1)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST COLUMNS (a) INTERVAL (1)
< error
| 1064 42000 …line 1 column 61 near "INTERVAL (1)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) INTERVAL (1)
< error
| 1064 42000 …line 1 column 53 near "INTERVAL (1)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY (a) INTERVAL (1)
< error
| 1064 42000 …line 1 column 52 near "INTERVAL (1)" 
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE (a) INTERVAL (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE (a) INTERVAL (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE (a) INTERVAL (1) FIRST PARTITION LESS THAN (1) LAST PARTITION LESS THAN (5)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE COLUMNS (a) INTERVAL (1 MONTH) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2021-01-01') NULL PARTITION MAXVALUE PARTITION
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE COLUMNS (a) INTERVAL (1 MONTH) FIRST PARTITION LESS THAN ('2020-01-01') LAST PARTITION LESS THAN ('2021-01-01') NULL PARTITION MAXVALUE PARTITION
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME
| ; SELECT FROM
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY)
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY)
| ; SELECT FROM
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x CURRENT, PARTITION y HISTORY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y HISTORY, PARTITION z CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x, PARTITION y CURRENT)
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x, PARTITION y CURRENT)
| ; SELECT FROM
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y)
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y)
| ; SELECT FROM
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x VALUES LESS THAN (1), PARTITION y CURRENT)
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x VALUES LESS THAN (1), PARTITION y CURRENT)
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x VALUES IN (1), PARTITION y CURRENT)
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x VALUES IN (1), PARTITION y CURRENT)
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x DEFAULT, PARTITION y CURRENT)
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x DEFAULT, PARTITION y CURRENT)
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL 1 DAY (PARTITION x HISTORY, PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL 1+1 DAY (PARTITION x HISTORY, PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL 1 FOO (PARTITION x HISTORY, PARTITION y CURRENT)
< error
| 1064 42000 …line 1 column 62 near "FOO (PARTITION x HISTORY, PARTITION y CURRENT)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL 1 (PARTITION x HISTORY, PARTITION y CURRENT)
< error
| 1064 42000 …line 1 column 60 near "(PARTITION x HISTORY, PARTITION y CURRENT)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL (PARTITION x HISTORY, PARTITION y CURRENT)
< error
| 1064 42000 …line 1 column 67 near "PARTITION x HISTORY, PARTITION y CURRENT)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME INTERVAL (1) DAY (PARTITION x HISTORY, PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME LIMIT 5 (PARTITION x HISTORY, PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME LIMIT 0 (PARTITION x HISTORY, PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME LIMIT (PARTITION x HISTORY, PARTITION y CURRENT)
< error
| 1064 42000 …line 1 column 55 near "(PARTITION x HISTORY, PARTITION y CURRENT)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME LIMIT 'a' (PARTITION x HISTORY, PARTITION y CURRENT)
< error
| 1064 42000 …line 1 column 57 near "'a' (PARTITION x HISTORY, PARTITION y CURRENT)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME LIMIT 5, 6 (PARTITION x HISTORY, PARTITION y CURRENT)
< error
| 1064 42000 …line 1 column 56 near ", 6 (PARTITION x HISTORY, PARTITION y CURRENT)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME PARTITIONS 2 (PARTITION x HISTORY, PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME PARTITIONS 3 (PARTITION x HISTORY, PARTITION y CURRENT)
< error
| 1484 HY000 Wrong number of partitions defined, mismatch with previous setting
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME PARTITIONS 3 (PARTITION x HISTORY, PARTITION y CURRENT)
| ; SELECT FROM
< error
| 1484 HY000 Wrong number of partitions defined, mismatch with previous setting
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME PARTITIONS 2
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME PARTITIONS 2
| ; SELECT FROM
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY SYSTEM_TIME SUBPARTITION BY HASH (b) (PARTITION x HISTORY, PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY (SUBPARTITION s), PARTITION y CURRENT (SUBPARTITION t))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY (SUBPARTITION s), PARTITION y CURRENT (SUBPARTITION t))
| ; SELECT FROM
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY COMMENT 'c', PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT) UPDATE INDEXES (i GLOBAL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LINEAR SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT)
< error
| 1064 42000 …line 1 column 54 near "SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT)" 
> exec ^
: string
| ALTER TABLE t PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION BY SYSTEM_TIME LIMIT 5 (PARTITION x HISTORY, PARTITION y CURRENT)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t PARTITION BY SYSTEM_TIME LIMIT 5 (PARTITION x HISTORY, PARTITION y CURRENT)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION x HISTORY)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION x HISTORY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION x CURRENT)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION x CURRENT)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t REORGANIZE PARTITION a INTO (PARTITION x HISTORY)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t REORGANIZE PARTITION a INTO (PARTITION x HISTORY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION x HISTORY)
< error
| 4113 HY000 Wrong partitioning type, expected type: `SYSTEM_TIME`
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION x HISTORY)
| ; SELECT FROM
< error
| 4113 HY000 Wrong partitioning type, expected type: `SYSTEM_TIME`
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION x CURRENT)
< error
| 4113 HY000 Wrong partitioning type, expected type: `SYSTEM_TIME`
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION x CURRENT)
| ; SELECT FROM
< error
| 4113 HY000 Wrong partitioning type, expected type: `SYSTEM_TIME`
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION x HISTORY)
< error
| 4113 HY000 Wrong partitioning type, expected type: `SYSTEM_TIME`
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION x HISTORY)
| ; SELECT FROM
< error
| 4113 HY000 Wrong partitioning type, expected type: `SYSTEM_TIME`
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY (a) (PARTITION x CURRENT)
< error
| 4113 HY000 Wrong partitioning type, expected type: `SYSTEM_TIME`
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY (a) (PARTITION x CURRENT)
| ; SELECT FROM
< error
| 4113 HY000 Wrong partitioning type, expected type: `SYSTEM_TIME`
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'a'
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'a'" 
> exec ^
: string
| CREATE TABLESPACE ts
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts" 
> exec ^
: string
| CREATE UNDO TABLESPACE ts ADD DATAFILE 'a.ibu'
< error
| 1064 42000 …line 1 column 11 near "UNDO TABLESPACE ts ADD DATAFILE 'a.ibu'" 
> exec ^
: string
| ALTER TABLESPACE ts ADD DATAFILE 'a'
< error
| 1064 42000 …line 1 column 16 near "TABLESPACE ts ADD DATAFILE 'a'" 
> exec ^
: string
| ALTER TABLESPACE ts RENAME TO b
< error
| 1064 42000 …line 1 column 16 near "TABLESPACE ts RENAME TO b" 
> exec ^
: string
| DROP TABLESPACE ts
< error
| 1064 42000 …line 1 column 15 near "TABLESPACE ts" 
> exec ^
: string
| DROP UNDO TABLESPACE ts
< error
| 1064 42000 …line 1 column 9 near "UNDO TABLESPACE ts" 
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'u'
< error
| 1064 42000 …line 1 column 14 near "LOGFILE GROUP lg ADD UNDOFILE 'u'" 
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'u'
< error
| 1064 42000 …line 1 column 13 near "LOGFILE GROUP lg ADD UNDOFILE 'u'" 
> exec ^
: string
| DROP LOGFILE GROUP lg
< error
| 1064 42000 …line 1 column 12 near "LOGFILE GROUP lg" 
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (USER 'r')
< error
| 1064 42000 …line 1 column 13 near "SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (USER 'r')" 
> exec ^
: string
| ALTER SERVER s OPTIONS (USER 'r')
< error
| 1064 42000 …line 1 column 12 near "SERVER s OPTIONS (USER 'r')" 
> exec ^
: string
| DROP SERVER s
< error
| 1064 42000 …line 1 column 11 near "SERVER s" 
> exec ^
: string
| CREATE TABLE t (a INT) TABLESPACE ts
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLESPACE = ts
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLESPACE ts STORAGE DISK
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p TABLESPACE ts)
< ok 0 0
> exec ^
: string
| ALTER TABLE t TABLESPACE ts
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t TABLESPACE ts
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t DISCARD TABLESPACE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t DISCARD TABLESPACE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t IMPORT TABLESPACE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t IMPORT TABLESPACE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| GRANT CREATE TABLESPACE ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT CREATE TABLESPACE ON d.* TO u
< ok 0 0
> exec ^
: string
| SHOW CREATE TABLESPACE ts
< error
| 1064 42000 …line 1 column 22 near "TABLESPACE ts" 
> exec ^
: string
| FLUSH TABLES t FOR EXPORT
< error
| 1064 42000 …line 1 column 18 near "FOR EXPORT" 
> exec ^
: string
| SELECT tablespace FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT tablespace FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE tablespace (a INT)
< ok 0 0
> exec ^
: string
| SELECT * FROM information_schema.tablespaces
< ok 0 0
