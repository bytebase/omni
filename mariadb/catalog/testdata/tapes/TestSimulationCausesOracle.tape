-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestSimulationCausesOracle default#1
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a
< error
| 1347 HY000 'omni_testsimulationcausesoracle.a' is not of type 'VIEW'
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a
< error
| 1146 42S02 Table 'omni_testsimulationcausesoracle.a' doesn't exist
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW a
< error
| 4092 42S02 Unknown VIEW: 'omni_testsimulationcausesoracle.a'
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW a
< error
| 4092 42S02 Unknown VIEW: 'omni_testsimulationcausesoracle.a'
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a
< error
| 1347 HY000 'omni_testsimulationcausesoracle.a' is not of type 'VIEW'
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY idx_d(missing)
< error
| 1072 42000 Key column 'missing' doesn't exist in table
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(missing)
< error
| 1072 42000 Key column 'missing' doesn't exist in table
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a
< error
| 1128 HY000 Function 'a' is not defined
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REMOVE PARTITIONING
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE REMOVE PARTITIONING
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REBUILD PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a OPTIMIZE PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ANALYZE PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHECK PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REPAIR PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY HASH(b) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE a REMOVE PARTITIONING
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a)
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c BLOB)
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT KEY(c)
< error
| 1283 HY000 Column 'c' cannot be part of FULLTEXT index
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARBINARY(10))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX i ON a(c)
< error
| 1283 HY000 Column 'c' cannot be part of FULLTEXT index
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c TEXT)
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(c)
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a
< error
| 1005 HY000 Can't create table `omni_testsimulationcausesoracle`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a
< error
| 1005 HY000 Can't create table `omni_testsimulationcausesoracle`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE z (q INT UNIQUE)
< ok 0 0
> exec
| ALTER TABLE a ADD q INT REFERENCES z
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `q` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `q` (`q`),
|   CONSTRAINT `a_ibfk_1` FOREIGN KEY (`q`) REFERENCES `z` (`q`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD q INT REFERENCES a(b)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `q` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `q` (`q`),
|   CONSTRAINT `a_ibfk_1` FOREIGN KEY (`q`) REFERENCES `a` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE z (q INT UNIQUE)
< ok 0 0
> exec
| ALTER TABLE a ADD q INT REFERENCES z(q, q)
< error
| 1239 42000 Incorrect foreign key definition for 'foreign key without name': Key reference and table reference don't match
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY k (c))
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a ADD q INT CONSTRAINT k REFERENCES p(id)
< error
| 1061 42000 Duplicate key name 'k'
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX i ON a(a) WITH PARSER p
< error
| 1128 HY000 Function 'p' is not defined
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE x (a INT, FULLTEXT KEY x (a) WITH PARSER p)
< error
| 1128 HY000 Function 'p' is not defined
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY HASH(b) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p ENGINE=zzz)
< error
| 1286 42000 Unknown storage engine 'zzz'
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE x (a INT) PARTITION BY HASH(a) (PARTITION p ENGINE=zzz)
< error
| 1286 42000 Unknown storage engine 'zzz'
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE x (a INT) ENGINE zzz
< error
| 1286 42000 Unknown storage engine 'zzz'
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE zzz
< error
| 1286 42000 Unknown storage engine 'zzz'
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE x ENGINE zzz
< error
| 1286 42000 Unknown storage engine 'zzz'
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| ALTER TABLE x ENGINE zzz
< error
| 1286 42000 Unknown storage engine 'zzz'
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE zzz
< error
| 1286 42000 Unknown storage engine 'zzz'
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE x (a INT) ENGINE=Aria
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE=MyISAM
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
