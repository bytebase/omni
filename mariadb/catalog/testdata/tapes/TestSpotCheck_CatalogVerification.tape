-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestSpotCheck_CatalogVerification default#1
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE child (
| 				a INT,
| 				b INT,
| 				CONSTRAINT child_ibfk_5 FOREIGN KEY (a) REFERENCES parent(id),
| 				FOREIGN KEY (b) REFERENCES parent(id)
| 			)
< ok 0 0
> query
|
| 			SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='child' AND CONSTRAINT_TYPE='FOREIGN KEY'
| 			ORDER BY CONSTRAINT_NAME
| 		
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| child_ibfk_1
< row
: bytes
| child_ibfk_5
> exec
| CREATE TABLE parent2 (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE child2 (a INT, b INT)
< ok 0 0
> exec
| ALTER TABLE child2 ADD CONSTRAINT child2_ibfk_5 FOREIGN KEY (a) REFERENCES parent2(id)
< ok 0 0
> exec
| ALTER TABLE child2 ADD FOREIGN KEY (b) REFERENCES parent2(id)
< ok 0 0
> query
|
| 			SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='child2' AND CONSTRAINT_TYPE='FOREIGN KEY'
| 			ORDER BY CONSTRAINT_NAME
| 		
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| child2_ibfk_5
< row
: bytes
| child2_ibfk_6
> exec
| DROP DATABASE IF EXISTS sc
< ok 4 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE pt (id INT) PARTITION BY HASH(id) PARTITIONS 4
< ok 0 0
> query
|
| 			SELECT PARTITION_NAME FROM information_schema.PARTITIONS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='pt'
| 			ORDER BY PARTITION_ORDINAL_POSITION
| 		
< result
< column "PARTITION_NAME" "VARCHAR"
< row
: bytes
| p0
< row
: bytes
| p1
< row
: bytes
| p2
< row
: bytes
| p3
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| SET SESSION explicit_defaults_for_timestamp=OFF
< ok 0 0
> exec
| SET SESSION sql_mode=''
< ok 0 0
> exec
| CREATE TABLE ts (
| 				ts1 TIMESTAMP NOT NULL,
| 				ts2 TIMESTAMP NOT NULL
| 			)
< ok 0 0
> query
|
| 			SELECT COLUMN_NAME, COLUMN_DEFAULT, EXTRA FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='ts'
| 			ORDER BY ORDINAL_POSITION
| 		
< result
< column "COLUMN_NAME" "VARCHAR"
< column "COLUMN_DEFAULT" "TEXT"
< column "EXTRA" "VARCHAR"
< row
: bytes
| ts1
: bytes
| current_timestamp()
: bytes
| on update current_timestamp()
< row
: bytes
| ts2
: bytes
| '0000-00-00 00:00:00'
: bytes
|
> exec
| SET SESSION sql_mode=DEFAULT
< ok 0 0
> exec
| SET SESSION explicit_defaults_for_timestamp=ON
< ok 0 0
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE pk (id INT, PRIMARY KEY(id))
< ok 0 0
> query
|
| 			SELECT IS_NULLABLE FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='pk' AND COLUMN_NAME='id'
| 		
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| NO
> exec
| DROP DATABASE IF EXISTS sc_cs
< ok 0 0
> exec
| CREATE DATABASE sc_cs CHARACTER SET latin1 COLLATE latin1_swedish_ci
< ok 1 0
> exec
| CREATE TABLE sc_cs.t (c VARCHAR(10))
< ok 0 0
> query
|
| 			SELECT TABLE_COLLATION FROM information_schema.TABLES
| 			WHERE TABLE_SCHEMA='sc_cs' AND TABLE_NAME='t'
| 		
< result
< column "TABLE_COLLATION" "VARCHAR"
< row
: bytes
| latin1_swedish_ci
> exec
| DROP DATABASE IF EXISTS sc_cs
< ok 1 0
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE p5 (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c5 (
| 				a INT,
| 				FOREIGN KEY (a) REFERENCES p5(id)
| 			)
< ok 0 0
> query
|
| 			SELECT DELETE_RULE, UPDATE_RULE FROM information_schema.REFERENTIAL_CONSTRAINTS
| 			WHERE CONSTRAINT_SCHEMA='sc' AND TABLE_NAME='c5'
| 		
< result
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< row
: bytes
| RESTRICT
: bytes
| RESTRICT
> query
| SHOW CREATE TABLE c5
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c5
: bytes
| CREATE TABLE `c5` (
|   `a` int(11) DEFAULT NULL,
|   KEY `a` (`a`),
|   CONSTRAINT `c5_ibfk_1` FOREIGN KEY (`a`) REFERENCES `p5` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS sc
< ok 2 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE base (id INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT * FROM base
< ok 0 0
> query
|
| 			SELECT SECURITY_TYPE FROM information_schema.VIEWS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='v'
| 		
< result
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| DEFINER
> exec
| DROP DATABASE IF EXISTS sc
< ok 2 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> query
| SELECT LENGTH(NOW()), LENGTH(NOW(6)), LENGTH(CURRENT_TIMESTAMP)
< result
< column "LENGTH(NOW())" "INT"
< column "LENGTH(NOW(6))" "INT"
< column "LENGTH(CURRENT_TIMESTAMP)" "INT"
< row
: int64 19
: int64 26
: int64 19
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE ai (id INT AUTO_INCREMENT PRIMARY KEY, v INT)
< ok 0 0
> query
| SHOW CREATE TABLE ai
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ai
: bytes
| CREATE TABLE `ai` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `v` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| INSERT INTO ai (v) VALUES (10),(20),(30)
< ok 3 1
> query
| SHOW CREATE TABLE ai
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ai
: bytes
| CREATE TABLE `ai` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `v` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tchk (
| 				a INT,
| 				CONSTRAINT tchk_chk_5 CHECK (a > 0),
| 				b INT,
| 				CHECK (b < 100)
| 			)
< ok 0 0
> query
|
| 			SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tchk' AND CONSTRAINT_TYPE='CHECK'
| 			ORDER BY CONSTRAINT_NAME
| 		
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| CONSTRAINT_1
< row
: bytes
| tchk_chk_5
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tchk2 (
| 				a INT,
| 				b INT,
| 				CONSTRAINT tchk2_chk_20 CHECK (a > 0)
| 			)
< ok 0 0
> exec
| ALTER TABLE tchk2 ADD CHECK (b > 0)
< ok 0 0
> query
|
| 			SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tchk2' AND CONSTRAINT_TYPE='CHECK'
| 			ORDER BY CONSTRAINT_NAME
| 		
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| CONSTRAINT_1
< row
: bytes
| tchk2_chk_20
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tps5 (ts DATETIME(6) DEFAULT NOW())
< ok 0 0
> query
|
| 			SELECT COLUMN_DEFAULT FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tps5' AND COLUMN_NAME='ts'
| 		
< result
< column "COLUMN_DEFAULT" "TEXT"
< row
: bytes
| current_timestamp(6)
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE p7 (id INT PRIMARY KEY)
< ok 0 0
> exec
|
| 			CREATE TABLE tps7 (
| 				a INT,
| 				CONSTRAINT tps7_ibfk_1 FOREIGN KEY (a) REFERENCES p7(id),
| 				b INT,
| 				FOREIGN KEY (b) REFERENCES p7(id)
| 			)
| 		
< error
| 1005 HY000 Can't create table `sc`.`tps7` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE cc (a INT, CHECK (a > 0), b INT, CHECK (b < 100))
< ok 0 0
> query
|
| 			SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='cc' AND CONSTRAINT_TYPE='CHECK'
| 			ORDER BY CONSTRAINT_NAME
| 		
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| CONSTRAINT_1
< row
: bytes
| CONSTRAINT_2
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE t2 (x REAL)
< ok 0 0
> query
|
| 			SELECT DATA_TYPE, COLUMN_TYPE FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='t2' AND COLUMN_NAME='x'
| 		
< result
< column "DATA_TYPE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| double
: bytes
| double
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tbool (flag BOOL)
< ok 0 0
> query
|
| 			SELECT DATA_TYPE, COLUMN_TYPE FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tbool' AND COLUMN_NAME='flag'
| 		
< result
< column "DATA_TYPE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| tinyint
: bytes
| tinyint(1)
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tai (id INT AUTO_INCREMENT PRIMARY KEY)
< ok 0 0
> query
|
| 			SELECT IS_NULLABLE FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tai' AND COLUMN_NAME='id'
| 		
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| NO
> exec
| DROP DATABASE IF EXISTS sc_cs2
< ok 0 0
> exec
| CREATE DATABASE sc_cs2 CHARACTER SET utf8mb4
< ok 1 0
> exec
| CREATE TABLE sc_cs2.t (c VARCHAR(10)) CHARSET latin1
< ok 0 0
> query
|
| 			SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='sc_cs2' AND TABLE_NAME='t' AND COLUMN_NAME='c'
| 		
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| latin1
: bytes
| latin1_swedish_ci
> query
| SHOW CREATE TABLE sc_cs2.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci
> exec
| DROP DATABASE IF EXISTS sc_cs2
< ok 1 0
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tnocs (x INT)
< ok 0 0
> query
| SHOW CREATE TABLE tnocs
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tnocs
: bytes
| CREATE TABLE `tnocs` (
|   `x` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE pm (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE cm (a INT, FOREIGN KEY (a) REFERENCES pm(id))
< ok 0 0
> query
|
| 			SELECT MATCH_OPTION FROM information_schema.REFERENTIAL_CONSTRAINTS
| 			WHERE CONSTRAINT_SCHEMA='sc' AND TABLE_NAME='cm'
| 		
< result
< column "MATCH_OPTION" "VARCHAR"
< row
: bytes
| NONE
> exec
| DROP DATABASE IF EXISTS sc
< ok 2 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE phd (id INT) PARTITION BY HASH(id)
< ok 0 0
> query
|
| 			SELECT COUNT(*) FROM information_schema.PARTITIONS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='phd' AND PARTITION_NAME IS NOT NULL
| 		
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 1
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
|
| 			CREATE TABLE psp (id INT, d INT)
| 			PARTITION BY RANGE(id)
| 			SUBPARTITION BY HASH(d) SUBPARTITIONS 3 (
| 				PARTITION p0 VALUES LESS THAN (10),
| 				PARTITION p1 VALUES LESS THAN (20)
| 			)
| 		
< ok 0 0
> query
|
| 			SELECT PARTITION_NAME, SUBPARTITION_NAME FROM information_schema.PARTITIONS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='psp'
| 			ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
| 		
< result
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< row
: bytes
| p0
: bytes
| p0sp0
< row
: bytes
| p0
: bytes
| p0sp1
< row
: bytes
| p0
: bytes
| p0sp2
< row
: bytes
| p1
: bytes
| p1sp0
< row
: bytes
| p1
: bytes
| p1sp1
< row
: bytes
| p1
: bytes
| p1sp2
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tidx (a INT, KEY(a))
< ok 0 0
> query
|
| 			SELECT INDEX_TYPE FROM information_schema.STATISTICS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tidx' AND INDEX_NAME='a'
| 		
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| BTREE
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE pfk (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE cfk (a INT, FOREIGN KEY (a) REFERENCES pfk(id))
< ok 0 0
> query
|
| 			SELECT INDEX_NAME, COLUMN_NAME FROM information_schema.STATISTICS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='cfk'
| 			ORDER BY INDEX_NAME, SEQ_IN_INDEX
| 		
< result
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| a
> exec
| DROP DATABASE IF EXISTS sc
< ok 2 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE teng (x INT)
< ok 0 0
> query
|
| 			SELECT ENGINE FROM information_schema.TABLES
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='teng'
| 		
< result
< column "ENGINE" "VARCHAR"
< row
: bytes
| InnoDB
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE trf (x INT)
< ok 0 0
> query
|
| 			SELECT ROW_FORMAT FROM information_schema.TABLES
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='trf'
| 		
< result
< column "ROW_FORMAT" "VARCHAR"
< row
: bytes
| Dynamic
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tgen (a INT, b INT AS (a + 1))
< ok 0 0
> query
|
| 			SELECT EXTRA FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tgen' AND COLUMN_NAME='b'
| 		
< result
< column "EXTRA" "VARCHAR"
< row
: bytes
| VIRTUAL GENERATED
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE base10 (id INT)
< ok 0 0
> exec
| CREATE VIEW v10 AS SELECT * FROM base10
< ok 0 0
> query
|
| 			SELECT VIEW_DEFINITION, CHECK_OPTION, DEFINER FROM information_schema.VIEWS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='v10'
| 		
< result
< column "VIEW_DEFINITION" "TEXT"
< column "CHECK_OPTION" "VARCHAR"
< column "DEFINER" "VARCHAR"
< row
: bytes
| select `sc`.`base10`.`id` AS `id` from `sc`.`base10`
: bytes
| NONE
: bytes
| root@%
> query
| SHOW CREATE VIEW v10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v10` AS select `base10`.`id` AS `id` from `base10`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS sc
< ok 2 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tt11 (a INT)
< ok 0 0
> exec
| CREATE TRIGGER trg11 BEFORE INSERT ON tt11 FOR EACH ROW SET NEW.a = NEW.a
< ok 0 0
> query
|
| 			SELECT DEFINER FROM information_schema.TRIGGERS
| 			WHERE TRIGGER_SCHEMA='sc' AND TRIGGER_NAME='trg11'
| 		
< result
< column "DEFINER" "VARCHAR"
< row
: bytes
| root@%
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tchk14 (a INT, CONSTRAINT chk14 CHECK (a > 0))
< ok 0 0
> query
|
| 			SELECT COLUMN_NAME FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='information_schema' AND TABLE_NAME='TABLE_CONSTRAINTS'
| 			  AND COLUMN_NAME='ENFORCED'
| 		
< result
< column "COLUMN_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tpos (a INT, b INT)
< ok 0 0
> exec
| ALTER TABLE tpos ADD COLUMN c INT
< ok 0 0
> query
|
| 			SELECT COLUMN_NAME FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tpos'
| 			ORDER BY ORDINAL_POSITION
| 		
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| c
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tnn (x INT, y INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE tnn
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tnn
: bytes
| CREATE TABLE `tnn` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE td (x INT DEFAULT NULL)
< ok 0 0
> query
|
| 			SELECT COLUMN_DEFAULT, IS_NULLABLE FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='td' AND COLUMN_NAME='x'
| 		
< result
< column "COLUMN_DEFAULT" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| NULL
: bytes
| YES
> exec
| DROP DATABASE IF EXISTS sc
< ok 1 0
> exec
| CREATE DATABASE sc
< ok 1 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE d (x DECIMAL, y DECIMAL(8), z NUMERIC)
< ok 0 0
> query
|
| 			SELECT COLUMN_NAME, COLUMN_TYPE, NUMERIC_PRECISION, NUMERIC_SCALE
| 			FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='d'
| 			ORDER BY ORDINAL_POSITION
| 		
< result
< column "COLUMN_NAME" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "NUMERIC_PRECISION" "UNSIGNED BIGINT"
< column "NUMERIC_SCALE" "UNSIGNED BIGINT"
< row
: bytes
| x
: bytes
| decimal(10,0)
: uint64 10
: uint64 0
< row
: bytes
| y
: bytes
| decimal(8,0)
: uint64 8
: uint64 0
< row
: bytes
| z
: bytes
| decimal(10,0)
: uint64 10
: uint64 0
