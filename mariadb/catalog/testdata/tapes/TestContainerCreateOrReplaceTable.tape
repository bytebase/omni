-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainerCreateOrReplaceTable/ALTER_dropping_an_index_keeps_a_CHECK_of_its_name default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 1 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE ca (a INT, KEY i (a), CONSTRAINT i CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE ca DROP INDEX i
< ok 0 0
> query
| SHOW CREATE TABLE `orp`.`ca`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ca
: bytes
| CREATE TABLE `ca` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `i` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/a_definition_rejected_before_the_drop_keeps_the_table default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 0 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE k1 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE k1 (a TIME(7))
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 6
> exec
| CREATE TABLE k2 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE k2 (a DECIMAL(5,6))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a')
> exec
| CREATE TABLE k3 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE k3 (a INT NOT NULL DEFAULT NULL)
< error
| 1067 42000 Invalid default value for 'a'
> exec
| CREATE TABLE k4 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE k4 (a INT) CHARSET latin1 COLLATE utf8mb4_bin
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> query
| SHOW CREATE TABLE `orp`.`k1`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| k1
: bytes
| CREATE TABLE `k1` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `orp`.`k2`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| k2
: bytes
| CREATE TABLE `k2` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `orp`.`k3`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| k3
: bytes
| CREATE TABLE `k3` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `orp`.`k4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| k4
: bytes
| CREATE TABLE `k4` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/a_definition_that_fails_after_the_drop_leaves_no_table default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 0 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE k1 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE k1 (a INT, a INT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| CREATE TABLE k2 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE k2 (a INT PRIMARY KEY, b INT PRIMARY KEY)
< error
| 1068 42000 Multiple primary key defined
> exec
| CREATE TABLE k3 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE k3 (a INT AUTO_INCREMENT PRIMARY KEY, b INT AUTO_INCREMENT)
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| CREATE TABLE k4 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE k4 (a INT, KEY i (a), KEY i (a))
< error
| 1061 42000 Duplicate key name 'i'
> query
| SHOW CREATE TABLE `orp`.`k1`
< error
| 1146 42S02 Table 'orp.k1' doesn't exist
> query
| SHOW CREATE TABLE `orp`.`k2`
< error
| 1146 42S02 Table 'orp.k2' doesn't exist
> query
| SHOW CREATE TABLE `orp`.`k3`
< error
| 1146 42S02 Table 'orp.k3' doesn't exist
> query
| SHOW CREATE TABLE `orp`.`k4`
< error
| 1146 42S02 Table 'orp.k4' doesn't exist
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/a_definition_without_columns_is_1113 default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 6 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE r (id INT)
< ok 0 0
> exec
| CREATE TRIGGER r_bi BEFORE INSERT ON r FOR EACH ROW SET NEW.id = 1
< ok 0 0
> exec
| CREATE OR REPLACE TABLE r
< error
| 1113 42000 A table must have at least 1 column
> exec
| CREATE TABLE r2 (id INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE r2 (INDEX (id))
< error
| 1113 42000 A table must have at least 1 column
> exec
| CREATE TABLE r3
< error
| 1113 42000 A table must have at least 1 column
> exec
| CREATE TABLE r4 (CHECK (1 > 0))
< error
| 1113 42000 A table must have at least 1 column
> query
| SHOW CREATE TABLE `orp`.`r`
< error
| 1146 42S02 Table 'orp.r' doesn't exist
> query
| SHOW CREATE TABLE `orp`.`r2`
< error
| 1146 42S02 Table 'orp.r2' doesn't exist
> query
| SHOW CREATE TABLE `orp`.`r3`
< error
| 1146 42S02 Table 'orp.r3' doesn't exist
> query
| SHOW CREATE TABLE `orp`.`r4`
< error
| 1146 42S02 Table 'orp.r4' doesn't exist
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/accepts_CREATE_TABLE_AS_SELECT_that_reads_another_table_of_the_name default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 1 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE src (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE orp2.src (q INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE src AS SELECT q FROM orp2.src
< ok 0 0
> exec
| CREATE TABLE src2 (id INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE src2 AS WITH src2 AS (SELECT 1 AS a) SELECT * FROM src2
< ok 1 0
> query
| SHOW CREATE TABLE `orp2`.`src`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| src
: bytes
| CREATE TABLE `src` (
|   `q` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/accepts_a_SELECT_whose_every_reference_of_the_name_reads_a_visible_CTE default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 1 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE s1 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE s1 AS WITH RECURSIVE s1 AS (SELECT 1 AS a UNION SELECT a + 1 FROM s1 WHERE a < 3) SELECT * FROM s1
< ok 3 0
> exec
| CREATE TABLE s2 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE s2 AS WITH s2 AS (SELECT 1 AS a), c1 AS (SELECT * FROM s2) SELECT * FROM c1
< ok 1 0
> exec
| CREATE TABLE s3 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE s3 AS SELECT * FROM (WITH s3 AS (SELECT 1 AS a) SELECT * FROM s3) d
< ok 1 0
> exec
| CREATE TABLE s4 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE s4 AS WITH s4 AS (SELECT 1 AS a) SELECT * FROM s4 WHERE EXISTS (SELECT 1 FROM s4)
< ok 1 0
> exec
| CREATE TABLE s5 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE s5 AS WITH s5 AS (SELECT 1 AS a) SELECT * FROM s5 UNION SELECT a FROM s5
< ok 1 0
> exec
| CREATE TABLE s6 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE s6 AS WITH RECURSIVE c1 AS (SELECT * FROM s6), s6 AS (SELECT 1 AS a) SELECT * FROM c1
< ok 1 0
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/creates_a_table_that_does_not_exist default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 1 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE OR REPLACE TABLE n (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE `orp`.`n`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| n
: bytes
| CREATE TABLE `n` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/creates_a_temporary_table_that_does_not_exist default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 0 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE OR REPLACE TEMPORARY TABLE t2 (x INT)
< ok 0 0
> query
| SHOW CREATE TABLE `orp`.`t2`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t2
: bytes
| CREATE TEMPORARY TABLE `t2` (
|   `x` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/dropping_an_index_keeps_a_CHECK_of_its_name default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 1 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE cd (a INT, KEY i (a), CONSTRAINT i CHECK (a > 0))
< ok 0 0
> exec
| DROP INDEX i ON cd
< ok 0 0
> query
| SHOW CREATE TABLE `orp`.`cd`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| cd
: bytes
| CREATE TABLE `cd` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `i` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/dropping_an_index_keeps_a_foreign_key_of_its_name default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 1 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE fp (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE fi (a INT, b INT, KEY j (a), KEY i (b), CONSTRAINT i FOREIGN KEY (a) REFERENCES fp (id))
< ok 0 0
> exec
| DROP INDEX i ON fi
< ok 0 0
> query
| SHOW CREATE TABLE `orp`.`fi`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| fi
: bytes
| CREATE TABLE `fi` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   KEY `j` (`a`),
|   CONSTRAINT `i` FOREIGN KEY (`a`) REFERENCES `fp` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/drops_the_triggers_of_the_table_it_replaces,_even_when_the_replacement_fails default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 4 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE tg (a INT)
< ok 0 0
> exec
| CREATE TRIGGER tg_bi BEFORE INSERT ON tg FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TABLE tg2 (a INT)
< ok 0 0
> exec
| CREATE TRIGGER tg2_bi BEFORE INSERT ON tg2 FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TABLE keep (a INT)
< ok 0 0
> exec
| CREATE TRIGGER keep_bi BEFORE INSERT ON keep FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE OR REPLACE TABLE tg (a INT, b INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE tg2 (a INT, a INT)
< error
| 1060 42S21 Duplicate column name 'a'
> query
| SHOW CREATE TABLE `orp`.`tg`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tg
: bytes
| CREATE TABLE `tg` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `orp`.`tg2`
< error
| 1146 42S02 Table 'orp.tg2' doesn't exist
> query
| SHOW CREATE TABLE `orp`.`keep`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| keep
: bytes
| CREATE TABLE `keep` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| keep_bi

== TestContainerCreateOrReplaceTable/keeps_the_table_when_the_LIKE_source_is_missing_or_a_view default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 1 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE dst (a INT)
< ok 0 0
> exec
| CREATE VIEW sv AS SELECT 1 AS a
< ok 0 0
> exec
| CREATE OR REPLACE TABLE dst LIKE nosuch
< error
| 1146 42S02 Table 'orp.nosuch' doesn't exist
> exec
| CREATE OR REPLACE TABLE dst LIKE nosuchdb.x
< error
| 1146 42S02 Table 'nosuchdb.x' doesn't exist
> exec
| CREATE OR REPLACE TABLE dst LIKE sv
< error
| 1347 HY000 'orp.sv' is not of type 'BASE TABLE'
> exec
| CREATE TABLE fresh LIKE sv
< error
| 1347 HY000 'orp.sv' is not of type 'BASE TABLE'
> exec
| CREATE TABLE fresh LIKE nosuchdb.x
< error
| 1146 42S02 Table 'nosuchdb.x' doesn't exist
> query
| SHOW CREATE TABLE `orp`.`dst`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| dst
: bytes
| CREATE TABLE `dst` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `orp`.`fresh`
< error
| 1146 42S02 Table 'orp.fresh' doesn't exist
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/refuses_CREATE_TABLE_AS_SELECT_that_reads_the_table_it_replaces default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 3 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE src (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE src AS SELECT * FROM src
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> exec
| CREATE OR REPLACE TABLE src AS SELECT y.id FROM src AS y
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> exec
| CREATE OR REPLACE TABLE src AS SELECT * FROM orp.src
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> exec
| CREATE OR REPLACE TABLE src AS SELECT 1 AS a FROM dual WHERE EXISTS (SELECT 1 FROM src)
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> exec
| CREATE OR REPLACE TABLE src AS SELECT 1 AS a FROM (SELECT id FROM src) d
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> exec
| CREATE OR REPLACE TABLE src AS SELECT (SELECT MAX(id) FROM src) AS m
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> exec
| CREATE OR REPLACE TABLE src AS SELECT 1 AS a UNION SELECT id FROM src
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> exec
| CREATE OR REPLACE TABLE src AS WITH c AS (SELECT id FROM src) SELECT * FROM c
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> exec
| CREATE OR REPLACE TABLE src (x INT) AS SELECT * FROM src
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> query
| SHOW CREATE TABLE `orp`.`src`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| src
: bytes
| CREATE TABLE `src` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/refuses_LIKE_of_the_table_itself,_existing_or_not default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 2 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE src (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE src LIKE src
< error
| 1066 42000 Not unique table/alias: 'src'
> exec
| CREATE OR REPLACE TABLE orp.src LIKE src
< error
| 1066 42000 Not unique table/alias: 'src'
> exec
| CREATE TABLE src LIKE src
< error
| 1066 42000 Not unique table/alias: 'src'
> exec
| CREATE TABLE IF NOT EXISTS src LIKE src
< error
| 1066 42000 Not unique table/alias: 'src'
> exec
| CREATE TABLE nx LIKE nx
< error
| 1066 42000 Not unique table/alias: 'nx'
> exec
| CREATE OR REPLACE TABLE nx LIKE nx
< error
| 1066 42000 Not unique table/alias: 'nx'
> exec
| CREATE TEMPORARY TABLE tm (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TEMPORARY TABLE tm LIKE tm
< error
| 1066 42000 Not unique table/alias: 'tm'
> query
| SHOW CREATE TABLE `orp`.`src`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| src
: bytes
| CREATE TABLE `src` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `orp`.`nx`
< error
| 1146 42S02 Table 'orp.nx' doesn't exist
> query
| SHOW CREATE TABLE `orp`.`tm`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tm
: bytes
| CREATE TEMPORARY TABLE `tm` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/refuses_a_SELECT_that_reads_the_table_outside_the_CTE_of_its_name default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 2 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 1 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE src (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE src AS SELECT * FROM src WHERE EXISTS (WITH src AS (SELECT 1) SELECT * FROM src)
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> exec
| CREATE OR REPLACE TABLE src AS WITH src AS (SELECT * FROM src) SELECT * FROM src
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> exec
| CREATE OR REPLACE TABLE src AS WITH c1 AS (SELECT * FROM src), src AS (SELECT 1 AS a) SELECT * FROM c1
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> exec
| CREATE OR REPLACE TABLE src AS SELECT (WITH src AS (SELECT 1 AS a) SELECT a FROM src) AS x, (SELECT COUNT(*) FROM src) AS y
< error
| 1093 HY000 Table 'src' is specified twice, both as a target for 'CREATE' and as a separate source for data
> query
| SHOW CREATE TABLE `orp`.`src`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| src
: bytes
| CREATE TABLE `src` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/refuses_a_foreign_key_parent_replaced_by_LIKE_or_CREATE_TABLE_AS_SELECT default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 2 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE src (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE par (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE chi (pid INT, CONSTRAINT fk_c FOREIGN KEY (pid) REFERENCES par (id))
< ok 0 0
> exec
| CREATE OR REPLACE TABLE par LIKE src
< error
| 1451 23000 Cannot delete or update a parent row: a foreign key constraint fails
> exec
| CREATE OR REPLACE TABLE par LIKE nosuch
< error
| 1146 42S02 Table 'orp.nosuch' doesn't exist
> exec
| CREATE OR REPLACE TABLE par AS SELECT 1 AS id
< error
| 1451 23000 Cannot delete or update a parent row: a foreign key constraint fails
> exec
| CREATE OR REPLACE TABLE par AS SELECT * FROM par
< error
| 1093 HY000 Table 'par' is specified twice, both as a target for 'CREATE' and as a separate source for data
> query
| SHOW CREATE TABLE `orp`.`par`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| par
: bytes
| CREATE TABLE `par` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `orp`.`chi`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| chi
: bytes
| CREATE TABLE `chi` (
|   `pid` int(11) DEFAULT NULL,
|   KEY `fk_c` (`pid`),
|   CONSTRAINT `fk_c` FOREIGN KEY (`pid`) REFERENCES `par` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/refuses_a_temporary_CREATE_TABLE_AS_SELECT_that_reads_the_table_it_replaces default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 1 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TEMPORARY TABLE tm (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TEMPORARY TABLE tm AS SELECT x.a FROM tm AS x
< error
| 1137 HY000 Can't reopen table: 'x'
> exec
| CREATE OR REPLACE TEMPORARY TABLE tm AS SELECT 1 AS a FROM dual WHERE EXISTS (SELECT 1 FROM tm)
< error
| 1137 HY000 Can't reopen table: 'tm'
> query
| SHOW CREATE TABLE `orp`.`tm`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tm
: bytes
| CREATE TEMPORARY TABLE `tm` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/refuses_a_view_of_the_name default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 1 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE VIEW vw AS SELECT 1 AS a
< ok 0 0
> exec
| CREATE OR REPLACE TABLE vw (id INT)
< error
| 1965 42S02 'orp.vw' is a view
> exec
| CREATE OR REPLACE TABLE vw AS SELECT 1 AS a
< error
| 1965 42S02 'orp.vw' is a view
> exec
| CREATE OR REPLACE TABLE vw LIKE nosuch
< error
| 1146 42S02 Table 'orp.nosuch' doesn't exist
> exec
| CREATE OR REPLACE TABLE vw LIKE vw
< error
| 1066 42000 Not unique table/alias: 'vw'
> exec
| CREATE TABLE vw (id INT)
< error
| 1050 42S01 Table 'vw' already exists
> query
| SHOW CREATE TABLE `orp`.`vw`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| vw
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `vw` AS select 1 AS `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/refuses_an_unknown_database default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 2 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE OR REPLACE TABLE nosuchdb.x (a INT)
< error
| 1049 42000 Unknown database 'nosuchdb'
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/refuses_to_drop_a_foreign_key_parent default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 1 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE par (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE chi (pid INT, CONSTRAINT fk_c FOREIGN KEY (pid) REFERENCES par (id))
< ok 0 0
> exec
| CREATE OR REPLACE TABLE par (id INT PRIMARY KEY, x INT)
< error
| 1451 23000 Cannot delete or update a parent row: a foreign key constraint fails
> query
| SHOW CREATE TABLE `orp`.`par`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| par
: bytes
| CREATE TABLE `par` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `orp`.`chi`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| chi
: bytes
| CREATE TABLE `chi` (
|   `pid` int(11) DEFAULT NULL,
|   KEY `fk_c` (`pid`),
|   CONSTRAINT `fk_c` FOREIGN KEY (`pid`) REFERENCES `par` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/refuses_to_drop_a_foreign_key_parent_of_another_database,_until_the_child_goes default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 2 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE par (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE orp2.chi (pid INT, CONSTRAINT fk_x FOREIGN KEY (pid) REFERENCES orp.par (id))
< ok 0 0
> exec
| CREATE OR REPLACE TABLE par (id INT PRIMARY KEY)
< error
| 1451 23000 Cannot delete or update a parent row: a foreign key constraint fails
> exec
| CREATE OR REPLACE TABLE orp2.chi (pid INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE par (id INT)
< ok 0 0
> query
| SHOW CREATE TABLE `orp`.`par`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| par
: bytes
| CREATE TABLE `par` (
|   `id` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `orp2`.`chi`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| chi
: bytes
| CREATE TABLE `chi` (
|   `pid` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/replaces_a_foreign_key_parent_with_foreign_key_checks_off default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 1 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 1 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE par (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE chi (pid INT, CONSTRAINT fk_c FOREIGN KEY (pid) REFERENCES par (id))
< ok 0 0
> exec
| CREATE OR REPLACE TABLE par (id INT PRIMARY KEY, y INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE par (z INT)
< ok 0 0
> query
| SHOW CREATE TABLE `orp`.`par`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| par
: bytes
| CREATE TABLE `par` (
|   `z` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `orp`.`chi`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| chi
: bytes
| CREATE TABLE `chi` (
|   `pid` int(11) DEFAULT NULL,
|   KEY `fk_c` (`pid`),
|   CONSTRAINT `fk_c` FOREIGN KEY (`pid`) REFERENCES `par` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/replaces_a_table default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 2 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE r (id INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE r (id INT, v INT)
< ok 0 0
> query
| SHOW CREATE TABLE `orp`.`r`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| r
: bytes
| CREATE TABLE `r` (
|   `id` int(11) DEFAULT NULL,
|   `v` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/replaces_a_table_whose_foreign_key_references_itself default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 2 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE s (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_s FOREIGN KEY (pid) REFERENCES s (id))
< ok 0 0
> exec
| CREATE OR REPLACE TABLE s (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE `orp`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TABLE `s` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/replaces_a_table_with_CREATE_TABLE_AS_SELECT_of_declared_columns default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 0 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE src (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE src (x INT) SELECT 1 AS x
< ok 1 0
> query
| SHOW CREATE TABLE `orp`.`src`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| src
: bytes
| CREATE TABLE `src` (
|   `x` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/replaces_a_table_with_LIKE default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 0 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE src (id INT PRIMARY KEY, s VARCHAR(5))
< ok 0 0
> exec
| CREATE TABLE dst (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE dst LIKE src
< ok 0 0
> query
| SHOW CREATE TABLE `orp`.`src`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| src
: bytes
| CREATE TABLE `src` (
|   `id` int(11) NOT NULL,
|   `s` varchar(5) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `orp`.`dst`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| dst
: bytes
| CREATE TABLE `dst` (
|   `id` int(11) NOT NULL,
|   `s` varchar(5) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/replaces_a_temporary_table default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 1 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TEMPORARY TABLE t1 (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE TEMPORARY TABLE t1 (a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE `orp`.`t1`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t1
: bytes
| CREATE TEMPORARY TABLE `t1` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"

== TestContainerCreateOrReplaceTable/replacing_an_index_keeps_a_CHECK_of_its_name default#1
> exec
| DROP DATABASE IF EXISTS orp
< ok 0 0
> exec
| DROP DATABASE IF EXISTS orp2
< ok 0 0
> exec
| CREATE DATABASE orp
< ok 1 0
> exec
| CREATE DATABASE orp2
< ok 1 0
> exec
| USE orp
< ok 0 0
> exec
| CREATE TABLE ci (a INT, KEY i (a), CONSTRAINT i CHECK (a > 0))
< ok 0 0
> exec
| CREATE OR REPLACE INDEX i ON ci (a)
< ok 0 0
> query
| SHOW CREATE TABLE `orp`.`ci`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ci
: bytes
| CREATE TABLE `ci` (
|   `a` int(11) DEFAULT NULL,
|   KEY `i` (`a`),
|   CONSTRAINT `i` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME
< result
< column "TRIGGER_NAME" "VARCHAR"
