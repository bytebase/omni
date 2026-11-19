-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSimulationCausesOracle default#1
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REMOVE PARTITIONING
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p1 VALUES LESS THAN (10))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION PARTITIONS 2
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION p0
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION IF EXISTS p0
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COALESCE PARTITION 1
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT)
< ok 0 0
> exec
| ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE s
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ANALYZE PARTITION p0
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION p0 ATTRIBUTES = DEFAULT
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LAST PARTITION LESS THAN (10)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FIRST PARTITION LESS THAN (20)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE a REMOVE PARTITIONING
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE a ANALYZE PARTITION p0
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE a PARTITION p0 ATTRIBUTES = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY RANGE (b) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (20)
< ok 0 0
> exec
| ALTER TABLE a LAST PARTITION LESS THAN (1000)
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.*
< error
| 1051 42S02 Unknown table 't'
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.* FROM a t
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT count(*)
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE z
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE z (KEY (b))
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT *) d
< error
| 1051 42S02 Unknown table ''
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS WITH c AS (SELECT *) SELECT 1
< error
| 1051 42S02 Unknown table ''
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT * FROM DUAL
< error
| 1051 42S02 Unknown table ''
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE `dual` (x INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT * FROM `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS WITH c AS (SELECT *) SELECT * FROM (WITH c AS (SELECT 1 AS x) SELECT x FROM c) d
< error
| 1051 42S02 Unknown table ''
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT missing.* FROM a
< error
| 1051 42S02 Unknown table 'missing'
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT other.a.* FROM a
< error
| 1051 42S02 Unknown table 'a'
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.* FROM a AS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT 1 FROM DUAL
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE z (LIKE a)
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| ALTER TABLE nope ADD e TINYINT DEFAULT NOW()
< error
| 1067 42000 Invalid default value for 'e'
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD b TINYINT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'b'
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e TINYINT DEFAULT(CURRENT_TIMESTAMP)
< error
| 1067 42000 Invalid default value for 'e'
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e TINYINT DEFAULT CURRENT_TIMESTAMP(1)
< error
| 1067 42000 Invalid default value for 'e'
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY b INT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'b'
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e DATETIME DEFAULT NOW()
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE z (a INT DEFAULT NOW())
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE z (a DATETIME DEFAULT NOW())
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE z (a DATETIME DEFAULT(CURRENT_TIMESTAMP(1)))
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE z (a DATETIME(1) DEFAULT(CURRENT_TIMESTAMP(1)))
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER b SET DEFAULT (CURRENT_TIMESTAMP)
< error
| 1067 42000 Invalid default value for 'b'
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD dt DATETIME
< ok 0 0
> exec
| ALTER TABLE a ALTER dt SET DEFAULT (CURRENT_TIMESTAMP)
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER b SET DEFAULT (NOW())
< error
| 3770 HY000 Default value expression of column 'b' contains a disallowed function: `NOW`.
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER b SET DEFAULT (LOCALTIME)
< error
| 3770 HY000 Default value expression of column 'b' contains a disallowed function: `LOCALTIME`.
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD dt DATETIME(1)
< ok 0 0
> exec
| ALTER TABLE a ALTER dt SET DEFAULT (CURRENT_TIMESTAMP)
< error
| 1067 42000 Invalid default value for 'dt'
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD dt DATETIME(1)
< ok 0 0
> exec
| ALTER TABLE a ALTER dt SET DEFAULT (CURRENT_TIMESTAMP(1))
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
