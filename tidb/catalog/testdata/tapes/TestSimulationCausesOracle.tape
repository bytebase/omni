-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSimulationCausesOracle default#1
> exec
| ALTER TABLE a NOCACHE
< error
| 1146 42S02 Table 'test.a' doesn't exist
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
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD PARTITION (PARTITION p1 VALUES LESS THAN (10))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD PARTITION PARTITIONS 2
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a DROP PARTITION p0
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a DROP PARTITION IF EXISTS p0
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a TRUNCATE PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a COALESCE PARTITION 1
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| CREATE TABLE s (b INT)
< ok 0 0
> exec
| ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE s
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ANALYZE PARTITION p0
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a PARTITION p0 ATTRIBUTES = DEFAULT
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a LAST PARTITION LESS THAN (10)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a FIRST PARTITION LESS THAN (20)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a REMOVE PARTITIONING
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a NOCACHE
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
| CREATE VIEW v AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| ALTER TABLE a NOCACHE
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
| CREATE VIEW v AS SELECT t.*
< error
| 1051 42S02 Unknown table 't'
> exec
| ALTER TABLE a NOCACHE
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
| CREATE VIEW v AS SELECT t.* FROM a t
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a NOCACHE
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
| CREATE TABLE z
< error
| 1113 42000 A table must have at least 1 column
> exec
| ALTER TABLE a NOCACHE
< error
| 1146 42S02 Table 'test.a' doesn't exist
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
| ALTER TABLE a NOCACHE
< error
| 1146 42S02 Table 'test.a' doesn't exist
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
| ALTER TABLE a NOCACHE
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
| CREATE VIEW v AS WITH c AS (SELECT *) SELECT 1
< error
| 1051 42S02 Unknown table ''
> exec
| ALTER TABLE a NOCACHE
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
| CREATE VIEW v AS SELECT * FROM DUAL
< error
| 1051 42S02 Unknown table ''
> exec
| ALTER TABLE a NOCACHE
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
| CREATE TABLE `dual` (x INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT * FROM `dual`
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
< error
| 1146 42S02 Table 'test.a' doesn't exist
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
| ALTER TABLE a NOCACHE
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
| CREATE VIEW v AS SELECT missing.* FROM a
< error
| 1051 42S02 Unknown table 'missing'
> exec
| ALTER TABLE a NOCACHE
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
| CREATE VIEW v AS SELECT other.a.* FROM a
< error
| 1051 42S02 Unknown table 'a'
> exec
| ALTER TABLE a NOCACHE
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
| CREATE VIEW v AS SELECT t.* FROM a AS t
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a NOCACHE
< error
| 1146 42S02 Table 'test.a' doesn't exist
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
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e TINYINT DEFAULT(CURRENT_TIMESTAMP)
< error
| 1067 42000 Invalid default value for 'e'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e TINYINT DEFAULT CURRENT_TIMESTAMP(1)
< error
| 1067 42000 Invalid default value for 'e'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a MODIFY b INT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'b'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e DATETIME DEFAULT NOW()
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a NOCACHE
< error
| 1146 42S02 Table 'test.a' doesn't exist
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
| ALTER TABLE a NOCACHE
< error
| 1146 42S02 Table 'test.a' doesn't exist
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
| ALTER TABLE a NOCACHE
< error
| 1146 42S02 Table 'test.a' doesn't exist
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
| ALTER TABLE a NOCACHE
< error
| 1146 42S02 Table 'test.a' doesn't exist
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
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD dt DATETIME
< ok 0 0
> exec
| ALTER TABLE a ALTER dt SET DEFAULT (CURRENT_TIMESTAMP)
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ALTER b SET DEFAULT (LOCALTIME)
< error
| 3770 HY000 Default value expression of column 'b' contains a disallowed function: `LOCALTIME`.
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD dt DATETIME(1)
< ok 0 0
> exec
| ALTER TABLE a ALTER dt SET DEFAULT (CURRENT_TIMESTAMP)
< error
| 1067 42000 Invalid default value for 'dt'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD dt DATETIME(1)
< ok 0 0
> exec
| ALTER TABLE a ALTER dt SET DEFAULT (CURRENT_TIMESTAMP(1))
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a PARTITION BY HASH(a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a PARTITION BY KEY(a)
< error
| 1488 HY000 Field in list of fields for partition function not found in table
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD c2 INT
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(c2) INTERVAL(1) FIRST PARTITION LESS THAN (c) LAST PARTITION LESS THAN (20)
< error
| 1054 42S22 Unknown column 'c' in 'expression'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a PARTITION BY RANGE(a) INTERVAL(a) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (20)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e INT AS (a)
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e INT AS (c + 1)
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e INT AS (c) STORED
< error
| 3106 HY000 'Adding generated stored column through ALTER TABLE' is not supported for generated columns.
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a TTL a + INTERVAL 1 DAY
< error
| 1054 42S22 Unknown column 'a' in 'TTL config'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a TTL c + INTERVAL 1 DAY
< error
| 8148 HY000 Field 'c' is of a not supported type for TTL config, expect DATETIME, DATE or TIMESTAMP
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD dt DATETIME
< ok 0 0
> exec
| ALTER TABLE a TTL dt + INTERVAL 1 DAY
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| DROP VIEW a
< error
| 1347 HY000 'test.a' is not VIEW
> exec
| ALTER TABLE a NOCACHE
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
| DROP VIEW IF EXISTS a
< error
| 1347 HY000 'test.a' is not VIEW
> exec
| ALTER TABLE a NOCACHE
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
| DROP VIEW IF EXISTS nope
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e INT DEFAULT NEXT VALUE FOR a
< error
| 1347 HY000 'test.a' is not SEQUENCE
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD c INT DEFAULT NEXT VALUE FOR a
< error
| 1347 HY000 'test.a' is not SEQUENCE
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a RENAME KEY c TO q
< error
| 1176 42000 Key 'c' doesn't exist in table 'a'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ALTER INDEX c VISIBLE
< error
| 1176 42000 Key 'c' doesn't exist in table 'a'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ALTER INDEX c INVISIBLE
< error
| 1176 42000 Key 'c' doesn't exist in table 'a'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e TINYINT AUTO_RANDOM
< error
| 8216 HY000 Invalid auto random: unsupported add column 'e' constraint AUTO_RANDOM when altering 'test.a'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD c TINYINT AUTO_RANDOM
< error
| 8216 HY000 Invalid auto random: unsupported add column 'c' constraint AUTO_RANDOM when altering 'test.a'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a AUTO_RANDOM_BASE 1
< error
| 8216 HY000 Invalid auto random: alter auto_random_base of a non auto_random table
> exec
| ALTER TABLE a NOCACHE
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
| CREATE TABLE a (b BIGINT PRIMARY KEY AUTO_RANDOM)
< ok 0 0
> exec
| ALTER TABLE a AUTO_RANDOM_BASE 2
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD FOREIGN KEY(c) REFERENCES a
< error
| 1239 42000 Incorrect foreign key definition for 'fk_1': Key reference and table reference don't match
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD FOREIGN KEY(c) REFERENCES a(b)
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e TINYTEXT UNICODE
< error
| 1115 42000 Unknown character set: 'ucs2'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE nope ADD e TINYTEXT UNICODE
< error
| 1115 42000 Unknown character set: 'ucs2'
> exec
| ALTER TABLE a NOCACHE
< error
| 1146 42S02 Table 'test.a' doesn't exist
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
| ALTER TABLE a PLACEMENT POLICY 'binary'
< error
| 8239 HY000 Unknown placement policy 'binary'
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e TINYINT ON UPDATE CURRENT_TIMESTAMP
< error
| 1294 HY000 Invalid ON UPDATE clause for 'e' column
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e ENUM('binary', 'binary')
< error
| 1291 HY000 Column 'e' has duplicated value 'binary' in ENUM
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e ENUM('x', 'y')
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e TINYINT DEFAULT a()
< error
| 3770 HY000 Default value expression of column 'e' contains a disallowed function: `a`.
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a ADD e DATETIME DEFAULT NOW()
< ok 0 0
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a CACHE
< ok 0 0
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /* CACHED ON */
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a PARTITION BY KEY()
< ok 0 0
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`b`) PARTITIONS 1
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a PARTITION BY LINEAR KEY()
< ok 0 0
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`b`) PARTITIONS 1
> exec
| ALTER TABLE a NOCACHE
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
| ALTER TABLE a PARTITION BY HASH(b)
< ok 0 0
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY HASH (`b`) PARTITIONS 1
> exec
| ALTER TABLE a NOCACHE
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
