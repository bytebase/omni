-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestOracle_DatabaseMetadataRestoresOnServer default#1
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_databasemetadatarestoresonserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_databasemetadatarestoresonserver`
< ok 1 0
> exec
| USE `omni_testoracle_databasemetadatarestoresonserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE SEQUENCE seq_order START WITH 100 INCREMENT BY 5 MINVALUE 1 MAXVALUE 9999 CACHE 20 CYCLE COMMENT='order seq'
< ok 0 0
> exec
| CREATE SEQUENCE seq_u32 AS INT UNSIGNED START WITH 1 NOCACHE NOCYCLE
< ok 0 0
> exec
| CREATE TABLE parent (
| 	id INT AUTO_INCREMENT PRIMARY KEY,
| 	code VARCHAR(20) NOT NULL DEFAULT 'abc',
| 	expr_str VARCHAR(20) NOT NULL DEFAULT '1 + 2',
| 	expr_col INT DEFAULT (1 + 2),
| 	fn_col DATE DEFAULT curdate(),
| 	bit_col BIT(3) DEFAULT b'101',
| 	nullable_col INT,
| 	ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	gen_stored VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_s')) STORED,
| 	gen_virt VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_v')) VIRTUAL,
| 	g POINT NOT NULL,
| 	secret INT INVISIBLE NOT NULL DEFAULT 0,
| 	SPATIAL INDEX sp_idx (g),
| 	INDEX idx_ignored (code) IGNORED COMMENT 'not used by optimizer',
| 	INDEX idx_prefix (code(10) DESC)
| ) ENGINE=InnoDB COMMENT='parent''s table'
< ok 0 0
> exec
| CREATE TABLE child (
| 	id INT PRIMARY KEY,
| 	parent_id INT,
| 	score INT,
| 	CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id) ON DELETE CASCADE,
| 	CONSTRAINT chk_score CHECK (score >= 0)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE part_tbl (
| 	id INT NOT NULL,
| 	yr INT NOT NULL,
| 	PRIMARY KEY (id, yr)
| ) ENGINE=InnoDB
| PARTITION BY RANGE (yr)
| SUBPARTITION BY HASH (id) SUBPARTITIONS 2 (
| 	PARTITION p0 VALUES LESS THAN (2025),
| 	PARTITION p1 VALUES LESS THAN MAXVALUE
| )
< ok 0 0
> exec
| CREATE TABLE versioned_explicit (
| 	id INT PRIMARY KEY,
| 	val VARCHAR(50),
| 	row_start TIMESTAMP(6) GENERATED ALWAYS AS ROW START,
| 	row_end TIMESTAMP(6) GENERATED ALWAYS AS ROW END,
| 	PERIOD FOR SYSTEM_TIME (row_start, row_end)
| ) ENGINE=InnoDB WITH SYSTEM VERSIONING
< ok 0 0
> exec
| CREATE VIEW v_cascaded (vid, vcode) AS SELECT id, code FROM parent WITH CASCADED CHECK OPTION
< ok 0 0
> exec
| CREATE VIEW v_local AS SELECT id, code FROM parent WITH LOCAL CHECK OPTION
< ok 0 0
> exec
| CREATE FUNCTION fn_add(a INT, b INT) RETURNS INT DETERMINISTIC COMMENT 'adds two ints'
| BEGIN
| 	RETURN a + b;
| END
< ok 0 0
> exec
| CREATE PROCEDURE proc_noop() COMMENT 'noop proc'
| BEGIN
| 	SELECT 1;
| END
< ok 0 0
> exec
| CREATE TRIGGER trg_parent AFTER INSERT ON parent FOR EACH ROW
| BEGIN
| 	INSERT INTO child(id) VALUES (NEW.id);
| END
< ok 0 0
> exec
| CREATE TRIGGER trg_z_first AFTER INSERT ON parent FOR EACH ROW PRECEDES trg_parent
| BEGIN
| 	INSERT INTO child(id) VALUES (NEW.id + 1);
| END
< ok 0 0
> exec
| CREATE EVENT ev_daily ON SCHEDULE EVERY 1 DAY COMMENT 'daily job' DO SELECT 1
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_databasemetadatarestoresonserver" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| parent
: bytes
| BASE TABLE
< row
: bytes
| part_tbl
: bytes
| BASE TABLE
< row
: bytes
| seq_order
: bytes
| SEQUENCE
< row
: bytes
| seq_u32
: bytes
| SEQUENCE
< row
: bytes
| v_cascaded
: bytes
| VIEW
< row
: bytes
| v_local
: bytes
| VIEW
< row
: bytes
| versioned_explicit
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `child`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| child
: bytes
| CREATE TABLE `child` (
|   `id` int(11) NOT NULL,
|   `parent_id` int(11) DEFAULT NULL,
|   `score` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_parent` (`parent_id`),
|   CONSTRAINT `fk_parent` FOREIGN KEY (`parent_id`) REFERENCES `parent` (`id`) ON DELETE CASCADE,
|   CONSTRAINT `chk_score` CHECK (`score` >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `parent`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| parent
: bytes
| CREATE TABLE `parent` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `code` varchar(20) NOT NULL DEFAULT 'abc',
|   `expr_str` varchar(20) NOT NULL DEFAULT '1 + 2',
|   `expr_col` int(11) DEFAULT (1 + 2),
|   `fn_col` date DEFAULT curdate(),
|   `bit_col` bit(3) DEFAULT b'101',
|   `nullable_col` int(11) DEFAULT NULL,
|   `ts` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
|   `gen_stored` varchar(30) GENERATED ALWAYS AS (concat(`code`,'_s')) STORED,
|   `gen_virt` varchar(30) GENERATED ALWAYS AS (concat(`code`,'_v')) VIRTUAL,
|   `g` point NOT NULL,
|   `secret` int(11) NOT NULL DEFAULT 0 INVISIBLE,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sp_idx` (`g`),
|   KEY `idx_ignored` (`code`) COMMENT 'not used by optimizer' IGNORED,
|   KEY `idx_prefix` (`code`(10) DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='parent''s table'
> query
| SHOW CREATE TABLE `part_tbl`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| part_tbl
: bytes
| CREATE TABLE `part_tbl` (
|   `id` int(11) NOT NULL,
|   `yr` int(11) NOT NULL,
|   PRIMARY KEY (`id`,`yr`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY RANGE (`yr`)
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION `p0` VALUES LESS THAN (2025) ENGINE = InnoDB,
|  PARTITION `p1` VALUES LESS THAN MAXVALUE ENGINE = InnoDB)
> query
| SHOW CREATE SEQUENCE `seq_order`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| seq_order
: bytes
| CREATE SEQUENCE `seq_order` start with 100 minvalue 1 maxvalue 9999 increment by 5 cache 20 cycle ENGINE=InnoDB COMMENT='order seq'
> query
| SHOW CREATE SEQUENCE `seq_u32`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| seq_u32
: bytes
| CREATE SEQUENCE `seq_u32` as int unsigned start with 1 minvalue 1 maxvalue 4294967294 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE VIEW `v_cascaded`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_cascaded
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cascaded` AS select `parent`.`id` AS `vid`,`parent`.`code` AS `vcode` from `parent` WITH CASCADED CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `v_local`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_local
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_local` AS select `parent`.`id` AS `id`,`parent`.`code` AS `code` from `parent` WITH LOCAL CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `versioned_explicit`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| versioned_explicit
: bytes
| CREATE TABLE `versioned_explicit` (
|   `id` int(11) NOT NULL,
|   `val` varchar(50) DEFAULT NULL,
|   `row_start` timestamp(6) GENERATED ALWAYS AS ROW START,
|   `row_end` timestamp(6) GENERATED ALWAYS AS ROW END,
|   PRIMARY KEY (`id`,`row_end`),
|   PERIOD FOR SYSTEM_TIME (`row_start`, `row_end`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| fn_add
> query
| SHOW CREATE FUNCTION `fn_add`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| fn_add
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_add`(a INT, b INT) RETURNS int(11)
|     DETERMINISTIC
|     COMMENT 'adds two ints'
| BEGIN
| 	RETURN a + b;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| proc_noop
> query
| SHOW CREATE PROCEDURE `proc_noop`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| proc_noop
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `proc_noop`()
|     COMMENT 'noop proc'
| BEGIN
| 	SELECT 1;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
< row
: bytes
| trg_z_first
: bytes
| AFTER
| INSERT
| parent
| 1
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
| BEGIN
| 	INSERT INTO child(id) VALUES (NEW.id + 1);
| END
< row
: bytes
| trg_parent
: bytes
| AFTER
| INSERT
| parent
| 2
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
| BEGIN
| 	INSERT INTO child(id) VALUES (NEW.id);
| END
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| ev_daily
> query
| SHOW CREATE EVENT `ev_daily`
< result
< column "Event" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "time_zone" "VARCHAR" "string"
< column "Create Event" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| ev_daily
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_daily` ON SCHEDULE EVERY 1 DAY STARTS '2026-11-13 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'daily job' DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_databasemetadatarestoresonserver`
< ok 8 0
> exec
| CREATE DATABASE `omni_testoracle_databasemetadatarestoresonserver`
< ok 1 0
> exec
| USE `omni_testoracle_databasemetadatarestoresonserver`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| SET @OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@@EXPLICIT_DEFAULTS_FOR_TIMESTAMP, EXPLICIT_DEFAULTS_FOR_TIMESTAMP=1
< ok 0 0
> exec
| CREATE VIEW `v_cascaded` AS SELECT
|   1 AS `vid`,
|   1 AS `vcode`
< ok 0 0
> exec
| CREATE VIEW `v_local` AS SELECT
|   1 AS `id`,
|   1 AS `code`
< ok 0 0
> exec
| CREATE SEQUENCE `seq_order` start with 100 minvalue 1 maxvalue 9999 increment by 5 cache 20 cycle COMMENT='order seq'
< ok 0 0
> exec
| CREATE SEQUENCE `seq_u32` as int unsigned start with 1 minvalue 1 maxvalue 4294967294 increment by 1 nocache nocycle
< ok 0 0
> exec
| CREATE TABLE `child` (
|   `id` int(11) NOT NULL,
|   `parent_id` int(11) DEFAULT NULL,
|   `score` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_parent` (`parent_id`),
|   CONSTRAINT `fk_parent` FOREIGN KEY (`parent_id`) REFERENCES `parent` (`id`) ON DELETE CASCADE,
|   CONSTRAINT `chk_score` CHECK (`score` >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| CREATE TABLE `parent` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `code` varchar(20) NOT NULL DEFAULT 'abc',
|   `expr_str` varchar(20) NOT NULL DEFAULT '1 + 2',
|   `expr_col` int(11) DEFAULT (1 + 2),
|   `fn_col` date DEFAULT curdate(),
|   `bit_col` bit(3) DEFAULT b'101',
|   `nullable_col` int(11) DEFAULT NULL,
|   `ts` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `gen_stored` varchar(30) GENERATED ALWAYS AS (concat(`code`,'_s')) STORED,
|   `gen_virt` varchar(30) GENERATED ALWAYS AS (concat(`code`,'_v')) VIRTUAL,
|   `g` point NOT NULL,
|   `secret` int(11) NOT NULL DEFAULT 0 INVISIBLE,
|   PRIMARY KEY (`id`),
|   KEY `idx_ignored` (`code`) COMMENT 'not used by optimizer' IGNORED,
|   KEY `idx_prefix` (`code`(10) DESC),
|   SPATIAL KEY `sp_idx` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='parent\'s table'
< ok 0 0
> exec
| CREATE TABLE `part_tbl` (
|   `id` int(11) NOT NULL,
|   `yr` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `yr`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
| PARTITION BY RANGE (`yr`)
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION `p0` VALUES LESS THAN (2025) ENGINE=InnoDB,
|  PARTITION `p1` VALUES LESS THAN MAXVALUE ENGINE=InnoDB)
< ok 0 0
> exec
| CREATE TABLE `versioned_explicit` (
|   `id` int(11) NOT NULL,
|   `val` varchar(50) DEFAULT NULL,
|   `row_start` timestamp(6) GENERATED ALWAYS AS ROW START NOT NULL,
|   `row_end` timestamp(6) GENERATED ALWAYS AS ROW END NOT NULL,
|   PERIOD FOR SYSTEM_TIME (`row_start`, `row_end`),
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
< ok 0 0
> exec
| SET @saved_cs_client = @@character_set_client, @saved_cs_results = @@character_set_results, @saved_col_connection = @@collation_connection, @saved_sql_mode = @@sql_mode
< ok 0 0
> exec
| CREATE FUNCTION `fn_add`(a INT, b INT) RETURNS int(11)
|     DETERMINISTIC
|     COMMENT 'adds two ints'
| BEGIN
| 	RETURN a + b;
| END
< ok 0 0
> exec
| SET character_set_client = @saved_cs_client, character_set_results = @saved_cs_results, collation_connection = @saved_col_connection, sql_mode = @saved_sql_mode
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_cascaded` AS select `parent`.`id` AS `vid`,`parent`.`code` AS `vcode` from `parent` WITH CASCADED CHECK OPTION
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_local` AS select `parent`.`id` AS `id`,`parent`.`code` AS `code` from `parent` WITH LOCAL CHECK OPTION
< ok 0 0
> exec
| CREATE PROCEDURE `proc_noop`()
|     COMMENT 'noop proc'
| BEGIN
| 	SELECT 1;
| END
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` EVENT `ev_daily` ON SCHEDULE EVERY 1 DAY ON COMPLETION NOT PRESERVE ENABLE COMMENT 'daily job' DO SELECT 1
< ok 0 0
> exec
| CREATE TRIGGER `trg_z_first`
| AFTER INSERT ON `parent` FOR EACH ROW
| BEGIN
| 	INSERT INTO child(id) VALUES (NEW.id + 1);
| END
< ok 0 0
> exec
| CREATE TRIGGER `trg_parent`
| AFTER INSERT ON `parent` FOR EACH ROW
| BEGIN
| 	INSERT INTO child(id) VALUES (NEW.id);
| END
< ok 0 0
> exec
| SET EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_databasemetadatarestoresonserver" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| parent
: bytes
| BASE TABLE
< row
: bytes
| part_tbl
: bytes
| BASE TABLE
< row
: bytes
| seq_order
: bytes
| SEQUENCE
< row
: bytes
| seq_u32
: bytes
| SEQUENCE
< row
: bytes
| v_cascaded
: bytes
| VIEW
< row
: bytes
| v_local
: bytes
| VIEW
< row
: bytes
| versioned_explicit
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `child`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| child
: bytes
| CREATE TABLE `child` (
|   `id` int(11) NOT NULL,
|   `parent_id` int(11) DEFAULT NULL,
|   `score` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_parent` (`parent_id`),
|   CONSTRAINT `fk_parent` FOREIGN KEY (`parent_id`) REFERENCES `parent` (`id`) ON DELETE CASCADE,
|   CONSTRAINT `chk_score` CHECK (`score` >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `parent`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| parent
: bytes
| CREATE TABLE `parent` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `code` varchar(20) NOT NULL DEFAULT 'abc',
|   `expr_str` varchar(20) NOT NULL DEFAULT '1 + 2',
|   `expr_col` int(11) DEFAULT (1 + 2),
|   `fn_col` date DEFAULT curdate(),
|   `bit_col` bit(3) DEFAULT b'101',
|   `nullable_col` int(11) DEFAULT NULL,
|   `ts` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
|   `gen_stored` varchar(30) GENERATED ALWAYS AS (concat(`code`,'_s')) STORED,
|   `gen_virt` varchar(30) GENERATED ALWAYS AS (concat(`code`,'_v')) VIRTUAL,
|   `g` point NOT NULL,
|   `secret` int(11) NOT NULL DEFAULT 0 INVISIBLE,
|   PRIMARY KEY (`id`),
|   KEY `idx_ignored` (`code`) COMMENT 'not used by optimizer' IGNORED,
|   KEY `idx_prefix` (`code`(10) DESC),
|   SPATIAL KEY `sp_idx` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='parent''s table'
> query
| SHOW CREATE TABLE `part_tbl`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| part_tbl
: bytes
| CREATE TABLE `part_tbl` (
|   `id` int(11) NOT NULL,
|   `yr` int(11) NOT NULL,
|   PRIMARY KEY (`id`,`yr`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY RANGE (`yr`)
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION `p0` VALUES LESS THAN (2025) ENGINE = InnoDB,
|  PARTITION `p1` VALUES LESS THAN MAXVALUE ENGINE = InnoDB)
> query
| SHOW CREATE SEQUENCE `seq_order`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| seq_order
: bytes
| CREATE SEQUENCE `seq_order` start with 100 minvalue 1 maxvalue 9999 increment by 5 cache 20 cycle ENGINE=InnoDB COMMENT='order seq'
> query
| SHOW CREATE SEQUENCE `seq_u32`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| seq_u32
: bytes
| CREATE SEQUENCE `seq_u32` as int unsigned start with 1 minvalue 1 maxvalue 4294967294 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE VIEW `v_cascaded`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_cascaded
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cascaded` AS select `parent`.`id` AS `vid`,`parent`.`code` AS `vcode` from `parent` WITH CASCADED CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `v_local`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_local
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_local` AS select `parent`.`id` AS `id`,`parent`.`code` AS `code` from `parent` WITH LOCAL CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `versioned_explicit`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| versioned_explicit
: bytes
| CREATE TABLE `versioned_explicit` (
|   `id` int(11) NOT NULL,
|   `val` varchar(50) DEFAULT NULL,
|   `row_start` timestamp(6) GENERATED ALWAYS AS ROW START,
|   `row_end` timestamp(6) GENERATED ALWAYS AS ROW END,
|   PRIMARY KEY (`id`,`row_end`),
|   PERIOD FOR SYSTEM_TIME (`row_start`, `row_end`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| fn_add
> query
| SHOW CREATE FUNCTION `fn_add`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| fn_add
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_add`(a INT, b INT) RETURNS int(11)
|     DETERMINISTIC
|     COMMENT 'adds two ints'
| BEGIN
| 	RETURN a + b;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| proc_noop
> query
| SHOW CREATE PROCEDURE `proc_noop`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| proc_noop
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `proc_noop`()
|     COMMENT 'noop proc'
| BEGIN
| 	SELECT 1;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
< row
: bytes
| trg_z_first
: bytes
| AFTER
| INSERT
| parent
| 1
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
| BEGIN
| 	INSERT INTO child(id) VALUES (NEW.id + 1);
| END
< row
: bytes
| trg_parent
: bytes
| AFTER
| INSERT
| parent
| 2
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
| BEGIN
| 	INSERT INTO child(id) VALUES (NEW.id);
| END
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| ev_daily
> query
| SHOW CREATE EVENT `ev_daily`
< result
< column "Event" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "time_zone" "VARCHAR" "string"
< column "Create Event" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| ev_daily
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_daily` ON SCHEDULE EVERY 1 DAY STARTS '2026-11-13 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'daily job' DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_databasemetadatarestoresonserver`
< ok 8 0
> exec
| CREATE DATABASE `omni_testoracle_databasemetadatarestoresonserver`
< ok 1 0
> exec
| USE `omni_testoracle_databasemetadatarestoresonserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE SEQUENCE `seq_order` start with 100 minvalue 1 maxvalue 9999 increment by 5 cache 20 cycle COMMENT='order seq'
< ok 0 0
> exec
| CREATE SEQUENCE `seq_u32` as int unsigned start with 1 minvalue 1 maxvalue 4294967294 increment by 1 nocache nocycle
< ok 0 0
> exec
| CREATE TABLE `child` (
|   `id` int(11) NOT NULL,
|   `parent_id` int(11) DEFAULT NULL,
|   `score` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_parent` (`parent_id`),
|   CONSTRAINT `fk_parent` FOREIGN KEY (`parent_id`) REFERENCES `parent` (`id`) ON DELETE CASCADE,
|   CONSTRAINT `chk_score` CHECK (`score` >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| CREATE TABLE `parent` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `code` varchar(20) NOT NULL DEFAULT 'abc',
|   `expr_str` varchar(20) NOT NULL DEFAULT '1 + 2',
|   `expr_col` int(11) DEFAULT (1 + 2),
|   `fn_col` date DEFAULT curdate(),
|   `bit_col` bit(3) DEFAULT b'101',
|   `nullable_col` int(11) DEFAULT NULL,
|   `ts` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `gen_stored` varchar(30) GENERATED ALWAYS AS (concat(`code`,'_s')) STORED,
|   `gen_virt` varchar(30) GENERATED ALWAYS AS (concat(`code`,'_v')) VIRTUAL,
|   `g` point NOT NULL,
|   `secret` int(11) NOT NULL DEFAULT 0 INVISIBLE,
|   PRIMARY KEY (`id`),
|   KEY `idx_ignored` (`code`) COMMENT 'not used by optimizer' IGNORED,
|   KEY `idx_prefix` (`code`(10) DESC),
|   SPATIAL KEY `sp_idx` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='parent\'s table'
< ok 0 0
> exec
| CREATE TABLE `part_tbl` (
|   `id` int(11) NOT NULL,
|   `yr` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `yr`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
| PARTITION BY RANGE (`yr`)
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION `p0` VALUES LESS THAN (2025) ENGINE=InnoDB,
|  PARTITION `p1` VALUES LESS THAN MAXVALUE ENGINE=InnoDB)
< ok 0 0
> exec
| CREATE TABLE `versioned_explicit` (
|   `id` int(11) NOT NULL,
|   `val` varchar(50) DEFAULT NULL,
|   `row_start` timestamp(6) GENERATED ALWAYS AS ROW START NOT NULL,
|   `row_end` timestamp(6) GENERATED ALWAYS AS ROW END NOT NULL,
|   PERIOD FOR SYSTEM_TIME (`row_start`, `row_end`),
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
< ok 0 0
> exec
| CREATE FUNCTION `fn_add`(a INT, b INT) RETURNS int(11)
|     DETERMINISTIC
|     COMMENT 'adds two ints'
| BEGIN
| 	RETURN a + b;
| END
< ok 0 0
> exec
| CREATE PROCEDURE `proc_noop`()
|     COMMENT 'noop proc'
| BEGIN
| 	SELECT 1;
| END
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_cascaded` AS
| select
|   `parent`.`id` AS `vid`,
|   `parent`.`code` AS `vcode`
| from `parent`
| WITH CASCADED CHECK OPTION
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_local` AS
| select
|   `parent`.`id` AS `id`,
|   `parent`.`code` AS `code`
| from `parent`
| WITH LOCAL CHECK OPTION
< ok 0 0
> exec
| CREATE TRIGGER `trg_parent` AFTER INSERT ON `parent` FOR EACH ROW
| BEGIN
| 	INSERT INTO child(id) VALUES (NEW.id);
| END
< ok 0 0
> exec
| CREATE TRIGGER `trg_z_first` AFTER INSERT ON `parent` FOR EACH ROW
| BEGIN
| 	INSERT INTO child(id) VALUES (NEW.id + 1);
| END
< ok 0 0
> exec
| CREATE EVENT `ev_daily` ON SCHEDULE EVERY 1 DAY ON COMPLETION NOT PRESERVE ENABLE COMMENT 'daily job' DO SELECT 1
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_databasemetadatarestoresonserver" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| parent
: bytes
| BASE TABLE
< row
: bytes
| part_tbl
: bytes
| BASE TABLE
< row
: bytes
| seq_order
: bytes
| SEQUENCE
< row
: bytes
| seq_u32
: bytes
| SEQUENCE
< row
: bytes
| v_cascaded
: bytes
| VIEW
< row
: bytes
| v_local
: bytes
| VIEW
< row
: bytes
| versioned_explicit
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `child`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| child
: bytes
| CREATE TABLE `child` (
|   `id` int(11) NOT NULL,
|   `parent_id` int(11) DEFAULT NULL,
|   `score` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_parent` (`parent_id`),
|   CONSTRAINT `fk_parent` FOREIGN KEY (`parent_id`) REFERENCES `parent` (`id`) ON DELETE CASCADE,
|   CONSTRAINT `chk_score` CHECK (`score` >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `parent`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| parent
: bytes
| CREATE TABLE `parent` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `code` varchar(20) NOT NULL DEFAULT 'abc',
|   `expr_str` varchar(20) NOT NULL DEFAULT '1 + 2',
|   `expr_col` int(11) DEFAULT (1 + 2),
|   `fn_col` date DEFAULT curdate(),
|   `bit_col` bit(3) DEFAULT b'101',
|   `nullable_col` int(11) DEFAULT NULL,
|   `ts` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
|   `gen_stored` varchar(30) GENERATED ALWAYS AS (concat(`code`,'_s')) STORED,
|   `gen_virt` varchar(30) GENERATED ALWAYS AS (concat(`code`,'_v')) VIRTUAL,
|   `g` point NOT NULL,
|   `secret` int(11) NOT NULL DEFAULT 0 INVISIBLE,
|   PRIMARY KEY (`id`),
|   KEY `idx_ignored` (`code`) COMMENT 'not used by optimizer' IGNORED,
|   KEY `idx_prefix` (`code`(10) DESC),
|   SPATIAL KEY `sp_idx` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='parent''s table'
> query
| SHOW CREATE TABLE `part_tbl`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| part_tbl
: bytes
| CREATE TABLE `part_tbl` (
|   `id` int(11) NOT NULL,
|   `yr` int(11) NOT NULL,
|   PRIMARY KEY (`id`,`yr`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY RANGE (`yr`)
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION `p0` VALUES LESS THAN (2025) ENGINE = InnoDB,
|  PARTITION `p1` VALUES LESS THAN MAXVALUE ENGINE = InnoDB)
> query
| SHOW CREATE SEQUENCE `seq_order`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| seq_order
: bytes
| CREATE SEQUENCE `seq_order` start with 100 minvalue 1 maxvalue 9999 increment by 5 cache 20 cycle ENGINE=InnoDB COMMENT='order seq'
> query
| SHOW CREATE SEQUENCE `seq_u32`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| seq_u32
: bytes
| CREATE SEQUENCE `seq_u32` as int unsigned start with 1 minvalue 1 maxvalue 4294967294 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE VIEW `v_cascaded`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_cascaded
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cascaded` AS select `parent`.`id` AS `vid`,`parent`.`code` AS `vcode` from `parent` WITH CASCADED CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `v_local`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_local
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_local` AS select `parent`.`id` AS `id`,`parent`.`code` AS `code` from `parent` WITH LOCAL CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `versioned_explicit`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| versioned_explicit
: bytes
| CREATE TABLE `versioned_explicit` (
|   `id` int(11) NOT NULL,
|   `val` varchar(50) DEFAULT NULL,
|   `row_start` timestamp(6) GENERATED ALWAYS AS ROW START,
|   `row_end` timestamp(6) GENERATED ALWAYS AS ROW END,
|   PRIMARY KEY (`id`,`row_end`),
|   PERIOD FOR SYSTEM_TIME (`row_start`, `row_end`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| fn_add
> query
| SHOW CREATE FUNCTION `fn_add`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| fn_add
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_add`(a INT, b INT) RETURNS int(11)
|     DETERMINISTIC
|     COMMENT 'adds two ints'
| BEGIN
| 	RETURN a + b;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| proc_noop
> query
| SHOW CREATE PROCEDURE `proc_noop`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| proc_noop
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `proc_noop`()
|     COMMENT 'noop proc'
| BEGIN
| 	SELECT 1;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
< row
: bytes
| trg_parent
: bytes
| AFTER
| INSERT
| parent
| 1
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
| BEGIN
| 	INSERT INTO child(id) VALUES (NEW.id);
| END
< row
: bytes
| trg_z_first
: bytes
| AFTER
| INSERT
| parent
| 2
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
| BEGIN
| 	INSERT INTO child(id) VALUES (NEW.id + 1);
| END
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| ev_daily
> query
| SHOW CREATE EVENT `ev_daily`
< result
< column "Event" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "time_zone" "VARCHAR" "string"
< column "Create Event" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| ev_daily
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_daily` ON SCHEDULE EVERY 1 DAY STARTS '2026-11-13 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'daily job' DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
