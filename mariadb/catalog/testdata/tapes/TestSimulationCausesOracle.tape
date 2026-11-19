-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestSimulationCausesOracle default#1
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (g INT, id INT AUTO_INCREMENT, KEY k (g, id)) ENGINE=MEMORY
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, KEY k (b))
< ok 0 0
> exec
| CREATE INDEX i USING RTREE ON a (missing)
< error
| 1072 42000 Key column 'missing' doesn't exist in table
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| i
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, KEY k (b))
< ok 0 0
> exec
| CREATE INDEX i USING RTREE ON a (b)
< error
| 1210 HY000 Incorrect arguments to RTREE INDEX
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, KEY k (`missing name`))
< error
| 1072 42000 Key column 'missing name' doesn't exist in table
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT) ENGINE='DEFAULT'
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
|   `b` int(11) DEFAULT NULL
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
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE a (e ENUM('s', 'S')) COLLATE utf8mb4_bin
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
|   `e` enum('s','S') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (e ENUM('a', 'a '))
< error
| 1291 HY000 Column 'e' has duplicated value 'a' in ENUM
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (e ENUM('a ', 'b '))
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
|   `e` enum('a','b') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ROW_FORMAT=COMPACT
< error
| 1814 HY000 Tablespace has been discarded for table `a`
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD e INT
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD KEY IF NOT EXISTS idx_d (d) x = 1
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT) WITH SYSTEM VERSIONING
< ok 0 0
> exec
| ALTER TABLE a DROP SYSTEM VERSIONING PARTITION BY SYSTEM_TIME
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
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (e ENUM('I', 'i') COLLATE utf8mb4_turkish_ci)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (e ENUM('i', 'İ') COLLATE utf8mb4_turkish_ci)
< error
| 1291 HY000 Column 'e' has duplicated value 'i' in ENUM
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b TEXT, FULLTEXT KEY f (missing) WITH PARSER unavailable)
< error
| 1128 HY000 Function 'unavailable' is not defined
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b TEXT, FULLTEXT KEY f (missing))
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, KEY k (b) x = 1) ENGINE='DEFAULT'
< error
| 1911 HY000 Unknown option 'x'
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD KEY IF NOT EXISTS j (d) x = 1
< error
| 1911 HY000 Unknown option 'x'
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a FORCE
< error
| 1814 HY000 Tablespace has been discarded for table `a`
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD INDEX i (d)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| TRUNCATE TABLE a
< ok 0 0
> exec
| ALTER TABLE a IMPORT TABLESPACE
< error
| 1813 HY000 Tablespace for table 'omni_testsimulationcausesoracle/a' exists. Please DISCARD the tablespace before IMPORT
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD UNIQUE KEY IF NOT EXISTS idx_d (d, missing_period WITHOUT OVERLAPS)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ALGORITHM=COPY, ADD e INT
< error
| 1814 HY000 Tablespace has been discarded for table `a`
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ALGORITHM=INPLACE, ADD e INT
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (e ENUM('s', 'S') CHARACTER SET latin1) COLLATE utf8mb4_bin
< error
| 1291 HY000 Column 'e' has duplicated value 's' in ENUM
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, PRIMARY KEY (b(1), b))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b TEXT, PRIMARY KEY (b))
< error
| 1170 42000 BLOB/TEXT column 'b' used in key specification without a key length
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD e INT CHECK (e > f), ADD f INT CHECK (f < e)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD e INT CHECK (e > 0)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE INDEX i ON a (c)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE UNIQUE INDEX u ON a (c)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD COLUMN IF NOT EXISTS c ENUM('x', 'x'), ADD e INT AS (1)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD COLUMN IF NOT EXISTS c INT DEFAULT missing
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a RENAME KEY zzz TO yyy
< error
| 1176 42000 Key 'zzz' doesn't exist in table 'a'
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT)
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a, a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD (KEY k (c) x = 1)
< error
| 1911 HY000 Unknown option 'x'
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ALTER COLUMN c SET DEFAULT missing
< error
| 1054 42S22 Unknown column 'missing' in 'DEFAULT'
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ALTER COLUMN c SET DEFAULT e, ADD COLUMN e INT
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
|   `c` int(11) DEFAULT `e`,
|   `d` int(11) DEFAULT NULL,
|   `e` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, g INT AS (missing))
< error
| 1054 42S22 Unknown column 'missing' in 'GENERATED ALWAYS AS'
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD g INT AS (missing)
< error
| 1054 42S22 Unknown column 'missing' in 'GENERATED ALWAYS AS'
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD COLUMN IF NOT EXISTS c INT AS (e), ADD e INT AS (1)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD e INT AS (c), ADD COLUMN IF NOT EXISTS c INT AS (1)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT)
< ok 0 0
> exec
| ALTER TABLE a ADD e ENUM('s', 'S'), COLLATE utf8mb4_bin
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
|   `b` int(11) DEFAULT NULL,
|   `e` enum('s','S') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (e ENUM('s', 'S')) COLLATE utf8mb4_bin
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4
< error
| 1291 HY000 Column 'e' has duplicated value 's' in ENUM
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
|   `e` enum('s','S') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query ^
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `e` enum('s','S') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a MODIFY missing ENUM('s', 's')
< error
| 1054 42S22 Unknown column 'missing' in 'a'
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE a (e ENUM('s', 'S'))
< error
| 1291 HY000 Column 'e' has duplicated value 's' in ENUM
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE UNIQUE INDEX u ON a (missing, missing WITHOUT OVERLAPS)
< error
| 4156 HY000 Period `missing` is not found in table
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, UNIQUE (missing, missing WITHOUT OVERLAPS))
< error
| 4156 HY000 Period `missing` is not found in table
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD CHECK (e > 0), ADD e INT
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
|   `e` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   CONSTRAINT `CONSTRAINT_1` CHECK (`e` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TEMPORARY TABLE a (b INT)
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< error
| 3007 HY000 Cannot DISCARD/IMPORT tablespace associated with temporary table
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TEMPORARY TABLE a (b INT)
< ok 0 0
> exec
| ALTER TABLE a IMPORT TABLESPACE
< error
| 3007 HY000 Cannot DISCARD/IMPORT tablespace associated with temporary table
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query ^
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TEMPORARY TABLE a (b INT) ENGINE=MyISAM
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
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
| CREATE TEMPORARY TABLE `a` (
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
| CREATE TEMPORARY TABLE `a` (
|   `b` int(11) DEFAULT NULL
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT) PARTITION BY HASH (b) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< error
| 1031 HY000 Storage engine InnoDB of the table `omni_testsimulationcausesoracle`.`a` doesn't have this option
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY HASH (`b`)
| PARTITIONS 2
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY HASH (`b`)
| PARTITIONS 2
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD e INT CHECK (c > 0), RENAME COLUMN c TO z
< error
| 1054 42S22 Unknown column 'c' in 'CHECK'
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a ADD e INT CHECK (z > 0), RENAME COLUMN c TO z
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
|   `z` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `e` int(11) DEFAULT NULL CHECK (`z` > 0),
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b VARCHAR(3)) COLLATE utf8mb4_bin
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET DEFAULT
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
|   `b` varchar(3) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (e ENUM('s', 'S')) COLLATE utf8mb4_bin
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET DEFAULT
< error
| 1291 HY000 Column 'e' has duplicated value 's' in ENUM
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT) COLLATE utf8mb4_bin
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET DEFAULT, ADD e ENUM('s', 'S')
< error
| 1291 HY000 Column 'e' has duplicated value 's' in ENUM
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT) COLLATE utf8mb4_bin
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET DEFAULT, ADD e ENUM('s', 't')
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
|   `b` int(11) DEFAULT NULL,
|   `e` enum('s','t') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (e ENUM('s', 'S')) COLLATE utf8mb4_bin
< ok 0 0
> exec
| ALTER TABLE a COLLATE DEFAULT
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = 'IGNORE_BAD_TABLE_OPTIONS'
< ok 0 0
> exec
| CREATE TABLE a (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 DISTANCE=cosine x=1)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| ALTER TABLE a ADD c INT
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
|   `v` vector(3) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   VECTOR KEY `i` (`v`) `M`=6 `DISTANCE`=cosine /* `x`=1 */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE INDEX idx_d ON a (c) x = 1
< error
| 1061 42000 Duplicate key name 'idx_d'
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE INDEX i ON a (c) x = 1
< error
| 1911 HY000 Unknown option 'x'
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE INDEX idx_d ON a (missing) x = 1
< error
| 1072 42000 Key column 'missing' doesn't exist in table
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER DATABASE CHARACTER SET latin1 COLLATE latin1_bin
< ok 1 0
> exec
| CREATE TABLE a (b VARCHAR(10)) CHARSET DEFAULT
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
|   `b` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER DATABASE CHARACTER SET binary
< ok 1 0
> exec
| CREATE TABLE a (e ENUM('s', 'S') CHARACTER SET utf8mb4 COLLATE utf8mb4_bin)
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET DEFAULT
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
|   `e` enum('s','S') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=binary
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (e ENUM('s', 'S') COLLATE DEFAULT) COLLATE utf8mb4_bin
< error
| 1291 HY000 Column 'e' has duplicated value 's' in ENUM
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b VARCHAR(10) COLLATE DEFAULT) COLLATE utf8mb4_bin
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
|   `b` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (e ENUM('s', 'S') NOT NULL COLLATE utf8mb4_bin)
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
|   `e` enum('s','S') CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b VARCHAR(10))
< ok 0 0
> exec
| ALTER TABLE a COLLATE latin1_bin, CHARACTER SET latin1
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
|   `b` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (e ENUM('s', 'S')) COLLATE utf8mb4_bin
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4, MODIFY e ENUM('s', 't')
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
|   `e` enum('s','t') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, e ENUM('s', 'S')) COLLATE utf8mb4_bin
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4, DROP COLUMN e
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT) COLLATE utf8mb4_bin
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4, ADD e ENUM('s', 'S') COLLATE utf8mb4_bin
< error
| 1291 HY000 Column 'e' has duplicated value 's' in ENUM
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT)
< ok 0 0
> exec
| ALTER TABLE a ADD e ENUM('s', 'S'), CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_bin
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
|   `b` int(11) DEFAULT NULL,
|   `e` enum('s','S') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b VARCHAR(10)) CHARSET latin1
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4, COLLATE utf8mb4_bin
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
|   `b` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b VARCHAR(10), e ENUM('s', 'S') COLLATE utf8mb4_bin)
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET binary
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
|   `b` varbinary(10) DEFAULT NULL,
|   `e` enum('s','S') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=binary
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, c INT)
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN c, ADD COLUMN IF NOT EXISTS c INT DEFAULT missing
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
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, c INT)
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN c TO z, ADD COLUMN IF NOT EXISTS c INT DEFAULT missing
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
|   `b` int(11) DEFAULT NULL,
|   `z` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT)
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN c INT, ADD COLUMN IF NOT EXISTS c INT DEFAULT missing
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
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, KEY k (b))
< ok 0 0
> exec
| ALTER TABLE a DROP KEY k, ADD KEY IF NOT EXISTS k (b) x = 1
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
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT)
< ok 0 0
> exec
| ALTER TABLE a ADD KEY k (b), ADD KEY IF NOT EXISTS k (b) x = 1
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
|   `b` int(11) DEFAULT NULL,
|   KEY `k` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| ALTER TABLE a PARTITION BY HASH (b) PARTITIONS 2
< error
| 1814 HY000 Tablespace has been discarded for table `a`
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b VARCHAR(10) PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET latin1
< error
| 1814 HY000 Tablespace has been discarded for table `a`
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
|   `b` varchar(10) NOT NULL,
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
|   `b` varchar(10) NOT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b VARCHAR(10) PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (id INT AUTO_INCREMENT, c INT, KEY k (id))
< ok 0 0
> exec
| DROP INDEX k ON a
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `c` int(11) DEFAULT NULL,
|   KEY `k` (`id`)
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `c` int(11) DEFAULT NULL,
|   KEY `k` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (id INT AUTO_INCREMENT, c INT, KEY k (id))
< ok 0 0
> exec
| CREATE OR REPLACE INDEX k ON a (c)
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `c` int(11) DEFAULT NULL,
|   KEY `k` (`id`)
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `c` int(11) DEFAULT NULL,
|   KEY `k` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (id INT AUTO_INCREMENT, c INT, KEY k (id))
< ok 0 0
> exec
| CREATE OR REPLACE INDEX k ON a (id, c)
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `c` int(11) DEFAULT NULL,
|   KEY `k` (`id`,`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE VECTOR INDEX idx_d ON a (b)
< error
| 1210 HY000 Incorrect arguments to VECTOR INDEX
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
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (v VECTOR(3) NOT NULL, VECTOR INDEX idx_d (v))
< ok 0 0
> exec
| CREATE VECTOR INDEX idx_d ON a (v)
< error
| 1061 42000 Duplicate key name 'idx_d'
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
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `idx_d` (`v`)
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
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `idx_d` (`v`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (g GEOMETRY, KEY k (g(5)))
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
|   `g` geometry DEFAULT NULL,
|   KEY `k` (`g`(5))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE p (x INT, y INT, KEY (x, y))
< ok 0 0
> exec
| CREATE TABLE a (b INT, FOREIGN KEY (b, b) REFERENCES p(x, y))
< error
| 1060 42S21 Duplicate column name 'b'
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT)
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY (b, b) REFERENCES z(x, y)
< error
| 1060 42S21 Duplicate column name 'b'
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
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT x = 1) ENGINE='DEFAULT'
< error
| 1911 HY000 Unknown option 'x'
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| SET sql_mode = 'IGNORE_BAD_TABLE_OPTIONS'
< ok 0 0
> exec
| CREATE TABLE a (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 `m``x`=1)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| ALTER TABLE a ADD c INT
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
|   `v` vector(3) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   VECTOR KEY `i` (`v`) `M`=6 /* `m``x`=1 */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, CHECK (a.missing > 0))
< error
| 1054 42S22 Unknown column 'missing' in 'CHECK'
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, CHECK (p.b > 0))
< error
| 1901 HY000 Function or expression 'p.b' cannot be used in the CHECK clause of `CONSTRAINT_1`
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT DEFAULT (p.b))
< error
| 1901 HY000 Function or expression 'p.b' cannot be used in the DEFAULT clause of `b`
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE SEQUENCE id
< ok 0 0
> exec
| CREATE TABLE a (id INT DEFAULT NEXTVAL(id))
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
|   `id` int(11) DEFAULT nextval(`omni_testsimulationcausesoracle`.`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (d INT DEFAULT missing, KEY k (nosuch))
< error
| 1072 42000 Key column 'nosuch' doesn't exist in table
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT, CHECK (missing > 0), KEY k (nosuch))
< error
| 1072 42000 Key column 'nosuch' doesn't exist in table
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?
: string
| a
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT) WITH SYSTEM VERSIONING
< ok 0 0
> exec
| ALTER TABLE a ADD PERIOD FOR SYSTEM_TIME(b, b)
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
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e))
< ok 0 0
> exec
| ALTER TABLE a DROP PERIOD FOR app_time, ADD UNIQUE u (missing, app_time WITHOUT OVERLAPS)
< error
| 4156 HY000 Period `app_time` is not found in table
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
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL,
|   PERIOD FOR `app_time` (`s`, `e`)
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
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL,
|   PERIOD FOR `app_time` (`s`, `e`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (b INT COLLATE utf8mb4_bin)
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
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (j JSON COLLATE latin1_swedish_ci)
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
|   `j` longtext CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL CHECK (json_valid(`j`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (id INT AUTO_INCREMENT NOT NULL, c INT NOT NULL, UNIQUE KEY k (id), UNIQUE KEY u (c) IGNORED)
< ok 0 0
> exec
| CREATE OR REPLACE INDEX k ON a (c)
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `c` int(11) NOT NULL,
|   UNIQUE KEY `k` (`id`),
|   UNIQUE KEY `u` (`c`) IGNORED
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `c` int(11) NOT NULL,
|   UNIQUE KEY `k` (`id`),
|   UNIQUE KEY `u` (`c`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (s VARCHAR(10)) CHARSET=latin1
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET DEFAULT COLLATE DEFAULT
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
|   `s` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
| CREATE TABLE a (s VARCHAR(10)) CHARSET DEFAULT COLLATE DEFAULT
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
|   `s` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
< ok 1 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
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
| DROP SEQUENCE IF EXISTS id
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
