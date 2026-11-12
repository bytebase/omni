-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSpotCheck_CatalogVerification default#1
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
| child_ibfk_5
< row
: bytes
| fk_2
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
| fk_2
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
> exec
| USE sc
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
: null
: bytes
|
< row
: bytes
| ts2
: null
: bytes
|
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
| CREATE DATABASE sc_cs CHARACTER SET latin1 COLLATE latin1_bin
< ok 0 0
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
| latin1_bin
> exec
| DROP DATABASE IF EXISTS sc_cs
< ok 0 0
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
| NO ACTION
: bytes
| NO ACTION
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
|   `a` int DEFAULT NULL,
|   KEY `fk_1` (`a`),
|   CONSTRAINT `fk_1` FOREIGN KEY (`a`) REFERENCES `p5` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
> exec
| USE sc
< ok 0 0
> query
| SELECT LENGTH(NOW()), LENGTH(NOW(6)), LENGTH(CURRENT_TIMESTAMP)
< result
< column "LENGTH(NOW())" "BIGINT"
< column "LENGTH(NOW(6))" "BIGINT"
< column "LENGTH(CURRENT_TIMESTAMP)" "BIGINT"
< row
: int64 19
: int64 26
: int64 19
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
|   `id` int NOT NULL AUTO_INCREMENT,
|   `v` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `id` int NOT NULL AUTO_INCREMENT,
|   `v` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=30001
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tps5 (ts DATETIME(6) DEFAULT NOW())
< error
| 1067 42000 Invalid default value for 'ts'
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< ok 0 0
> query
|
| 			SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tps7' AND CONSTRAINT_TYPE='FOREIGN KEY'
| 			ORDER BY CONSTRAINT_NAME
| 		
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| fk_2
< row
: bytes
| tps7_ibfk_1
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< ok 0 0
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
| latin1_bin
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
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin
> exec
| DROP DATABASE IF EXISTS sc_cs2
< ok 0 0
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
|   `x` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
: null
< row
: bytes
| p1
: null
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
| fk_1
: bytes
| a
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
| Compact
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
| SELECT `sc`.`base10`.`id` AS `id` FROM `sc`.`base10`
: bytes
| CASCADED
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v10` (`id`) AS SELECT `sc`.`base10`.`id` AS `id` FROM `sc`.`base10`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tt11 (a INT)
< ok 0 0
> exec
| CREATE TRIGGER trg11 BEFORE INSERT ON tt11 FOR EACH ROW SET NEW.a = NEW.a
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 14 near "TRIGGER trg11 BEFORE INSERT ON tt11 FOR EACH ROW SET NEW.a = NEW.a" 
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
> exec
| USE sc
< ok 0 0
> exec
| CREATE TABLE tchk14 (a INT, CONSTRAINT chk14 CHECK (a > 0))
< ok 0 0
> query
|
| 			SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
| 			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tchk14' AND CONSTRAINT_TYPE='CHECK'
| 		
< result
< column "CONSTRAINT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
|   `x` int DEFAULT NULL,
|   `y` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
: null
: bytes
| YES
> exec
| DROP DATABASE IF EXISTS sc
< ok 0 0
> exec
| CREATE DATABASE sc
< ok 0 0
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
< column "NUMERIC_PRECISION" "BIGINT"
< column "NUMERIC_SCALE" "BIGINT"
< row
: bytes
| x
: bytes
| decimal(10,0)
: int64 10
: int64 0
< row
: bytes
| y
: bytes
| decimal(8,0)
: int64 8
: int64 0
< row
: bytes
| z
: bytes
| decimal(10,0)
: int64 10
: int64 0
