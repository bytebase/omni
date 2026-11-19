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
| CREATE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
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
| CREATE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
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
| CREATE VIEW a AS SELECT b FROM a
< error
| 1050 42S01 Table 'a' already exists
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
| CREATE VIEW a AS SELECT ABS(a)
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
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
| CREATE VIEW a AS SELECT 1 UNION SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
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
| CREATE VIEW a AS SELECT a HAVING 1
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
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
| CREATE VIEW a AS SELECT 1 WHERE a
< error
| 1054 42S22 Unknown column 'a' in 'WHERE'
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
| CREATE VIEW a AS SELECT (SELECT 1 AS x HAVING x)
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
| CREATE TABLE a SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
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
| CREATE TABLE a SELECT a
< error
| 1050 42S01 Table 'a' already exists
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
| CREATE TABLE a (g INT)
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(g)
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
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
| CREATE TABLE a (g INT NOT NULL)
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(g)
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
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
| CREATE TABLE a (g POINT)
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(g)
< error
| 1252 42000 All parts of a SPATIAL index must be NOT NULL
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
| CREATE TABLE a (g POINT NOT NULL)
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(g)
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
| ALTER TABLE a ROW_FORMAT DEFAULT
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
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
| ALTER TABLE a TABLESPACE a
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
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
| ALTER TABLE a ADD a SERIAL
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
|   `a` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
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
| ALTER TABLE a ADD a CHAR BYTE
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
|   `a` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
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
| ALTER TABLE a ADD a CHAR COLLATE BINARY
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
|   `a` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
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
| ALTER TABLE a MAX_ROWS 9223372036854775808
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=4294967295
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
| ALTER TABLE a MAX_ROWS 2e3
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=2
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
| ALTER TABLE a MIN_ROWS 3e2
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MIN_ROWS=3
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
| ALTER TABLE a MAX_ROWS 1.5
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=1
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
| ALTER TABLE a MAX_ROWS 1e0
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=1
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
| CREATE TABLE a (b INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| ALTER TABLE a MODIFY c SERIAL
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
|   `c` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `c` (`c`)
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
| CREATE TABLE a (b INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| ALTER TABLE a CHANGE c y SERIAL
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
|   `y` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `y` (`y`)
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
| ALTER TABLE a MIN_ROWS 1.5
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MIN_ROWS=1
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
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a)
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
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
| ALTER TABLE a ADD a INT AS(a) STORED
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
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
| ALTER TABLE a ADD a INT DEFAULT a
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
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
| ALTER TABLE a ADD a INT DEFAULT(a)
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
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
| ALTER TABLE a ALTER a SET DEFAULT a
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
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
| CREATE TABLE a (b INT PRIMARY KEY, c INT AS(b + 1) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE a ADD d INT AS(c)
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
| CREATE TABLE a (b INT PRIMARY KEY, c INT AS(b + 1) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE a MODIFY c INT AS(c)
< error
| 4029 01000 Expression for field `c` is referring to uninitialized field `c`
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
| CREATE TABLE a (b INT PRIMARY KEY, c INT AS(b + 1) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE a CHANGE c c INT AS(c)
< error
| 4029 01000 Expression for field `c` is referring to uninitialized field `c`
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
| CREATE TABLE x (a INT AS(a))
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
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
| CREATE TABLE x (a INT DEFAULT a)
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
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
| CREATE TABLE x (a INT AS(b), b INT AS(a))
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `b`
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
| CREATE TABLE x (a INT AS(b), b INT)
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
| ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
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
| ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin zzz
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
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
| ALTER TABLE a ADD a CHAR ASCII COLLATE BINARY
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'latin1'
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
| ALTER TABLE a ADD a CHAR BYTE COLLATE utf8mb4_bin
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'binary'
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
| ALTER TABLE a ADD a NCHAR COLLATE utf8mb4_bin
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'utf8mb3'
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
| ALTER TABLE x ADD a CHAR ASCII COLLATE utf8mb4_bin
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
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
| ALTER TABLE a ADD a CHAR CHARACTER SET ascii COLLATE ascii_general_ci
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
| CREATE TABLE a (x CHAR ASCII COLLATE utf8mb4_bin)
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
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
| ALTER TABLE x ADD FULLTEXT(a) WITH PARSER a
< error
| 1128 HY000 Function 'a' is not defined
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| x
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
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
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a zzz
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
| ALTER TABLE a ROW_FORMAT FIXED
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
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
| ALTER TABLE a ROW_FORMAT PAGE
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
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
| ALTER TABLE a ROW_FORMAT FIXED ENGINE=MyISAM
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
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d)) ENGINE=MyISAM
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT PAGE
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
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d)) ENGINE=Aria
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT FIXED
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
| CREATE TABLE x (k INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE x ADD q INT AS(r), ADD r INT AS(q)
< error
| 4029 01000 Expression for field `q` is referring to uninitialized field `r`
> query
| SHOW CREATE TABLE x
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| x
: bytes
| CREATE TABLE `x` (
|   `k` int(11) NOT NULL,
|   PRIMARY KEY (`k`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| x
: bytes
| CREATE TABLE `x` (
|   `k` int(11) NOT NULL,
|   PRIMARY KEY (`k`)
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
| CREATE TABLE x (k INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE x ADD (q INT AS(r), r INT AS(q))
< error
| 4029 01000 Expression for field `q` is referring to uninitialized field `r`
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
| CREATE TABLE x (k INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE x ADD q INT AS(r), ADD r INT
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
| ALTER TABLE a ENGINE=InnoDB ROW_FORMAT FIXED
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
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
| ALTER TABLE x OPTIMIZE PARTITION ALL
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| x
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
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
| ALTER TABLE x ANALYZE PARTITION p1
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| x
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
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
| ALTER TABLE x CHECK PARTITION ALL
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| x
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
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
| ALTER TABLE x REPAIR PARTITION ALL
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| x
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
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
| ALTER TABLE a PARTITION BY SYSTEM_TIME
< error
| 4124 HY000 Table `a` is not system-versioned
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
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
| CREATE TABLE a(PERIOD FOR SYSTEM_TIME(a, a))
< error
| 4125 HY000 Wrong parameters for `a`: missing 'WITH SYSTEM VERSIONING'
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
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
| ALTER TABLE a ADD PERIOD FOR SYSTEM_TIME(a, a)
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
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
| ALTER TABLE a ADD CHECK(a)
< error
| 1054 42S22 Unknown column 'a' in 'CHECK'
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
| ALTER TABLE a ADD CONSTRAINT IF NOT EXISTS a CHECK(a)
< error
| 1054 42S22 Unknown column 'a' in 'CHECK'
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
| ALTER TABLE a ADD a INT DEFAULT VALUE
< error
| 1054 42S22 Unknown column 'VALUE' in 'DEFAULT'
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
| ALTER TABLE a ADD KEY(c, c)
< error
| 1060 42S21 Duplicate column name 'c'
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
| ALTER TABLE a ADD KEY(c(1))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
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
| ALTER TABLE a ADD KEY USING RTREE(c)
< error
| 1210 HY000 Incorrect arguments to RTREE INDEX
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
| CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY a (c))
< ok 0 0
> exec
| CREATE SPATIAL INDEX a ON a(b)
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
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
| CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY a (c))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(b)
< error
| 1283 HY000 Column 'b' cannot be part of FULLTEXT index
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
| ALTER TABLE a ADD KEY(c) x = 1
< error
| 1911 HY000 Unknown option 'x'
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
| ALTER TABLE a ADD a INT AUTO_INCREMENT
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
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
| ALTER TABLE a ADD a ENUM('s', 's')
< error
| 1291 HY000 Column 'a' has duplicated value 's' in ENUM
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
| ALTER TABLE a ADD UNIQUE(c, c WITHOUT OVERLAPS)
< error
| 4156 HY000 Period `c` is not found in table
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
| ALTER TABLE a IMPORT TABLESPACE
< error
| 1813 HY000 Tablespace for table 'omni_testsimulationcausesoracle/a' exists. Please DISCARD the tablespace before IMPORT
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
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
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a IMPORT TABLESPACE
< error
| 1815 HY000 Internal error: Drop all secondary indexes before importing table omni_testsimulationcausesoracle/a when .cfg file is missing.
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
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
| CREATE TABLE a (b INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a IMPORT TABLESPACE
< error
| 1030 HY000 Got error 194 "Tablespace is missing for a table" from storage engine InnoDB
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
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`b`)
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
| CREATE TABLE a (b INT) ENGINE=MyISAM
< ok 0 0
> exec
| ALTER TABLE a IMPORT TABLESPACE
< error
| 1031 HY000 Storage engine MyISAM of the table `omni_testsimulationcausesoracle`.`a` doesn't have this option
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
|   `b` int(11) DEFAULT NULL
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) DEFAULT NULL
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
| CREATE TABLE a (g INT, id INT AUTO_INCREMENT, KEY k (g, id)) ENGINE=MyISAM
< ok 0 0
> exec
| ALTER TABLE a ADD c INT
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
| ALTER TABLE a RENAME KEY nosuch TO idx_d
< error
| 1176 42000 Key 'nosuch' doesn't exist in table 'a'
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
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
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
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
| ALTER TABLE a ADD SYSTEM VERSIONING PARTITION BY SYSTEM_TIME
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
| CREATE TABLE a (b INT, KEY k (missing, missing))
< error
| 1072 42000 Key column 'missing' doesn't exist in table
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
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
| CREATE TABLE a (b INT, PRIMARY KEY (missing, missing))
< error
| 1072 42000 Key column 'missing' doesn't exist in table
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
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
| CREATE TABLE a (e ENUM('s', 'S') COLLATE utf8mb4_bin)
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
|   `e` enum('s','S') CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL
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
| CREATE TABLE a (b INT PRIMARY KEY, a INT, KEY idx_d (b))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY KEY
< error
| 1060 42S21 Duplicate column name 'a'
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
