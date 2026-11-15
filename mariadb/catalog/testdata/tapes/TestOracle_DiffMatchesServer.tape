-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestOracle_DiffMatchesServer default#1
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b BIGINT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` bigint(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT(11), b BIGINT(20))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` bigint(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a BOOLEAN DEFAULT TRUE)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` tinyint(1) DEFAULT 1
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a TINYINT(1) DEFAULT 1)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` tinyint(1) DEFAULT 1
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT DEFAULT 1, d DECIMAL(10,2) DEFAULT 1.5)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT 1,
|   `d` decimal(10,2) DEFAULT 1.50
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT DEFAULT '1', d DECIMAL(10,2) DEFAULT '1.50')
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT 1,
|   `d` decimal(10,2) DEFAULT 1.50
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b TEXT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` text DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT DEFAULT NULL, b TEXT DEFAULT NULL)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` text DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a TIMESTAMP)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` timestamp NULL DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a TIMESTAMP NULL DEFAULT NULL)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` timestamp NULL DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10) CHARACTER SET utf8mb4) DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10) CHARACTER SET utf8)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10) CHARACTER SET utf8mb3)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT c CHECK (a > 0))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `c` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT c CHECK (((a > 0))))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `c` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT) ENGINE=innodb
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME DEFAULT NOW())
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT current_timestamp()
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT current_timestamp()
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> exec
| CREATE TABLE t (a INT DEFAULT NEXTVAL(s))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT nextval(`omni_testoracle_diffmatchesserver`.`s`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> exec
| CREATE TABLE t (a INT DEFAULT (NEXT VALUE FOR s))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT nextval(`omni_testoracle_diffmatchesserver`.`s`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE b (id INT)
< ok 0 0
> exec
| CREATE VIEW v (x) AS SELECT id FROM b
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| b
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| b
: bytes
| CREATE TABLE `b` (
|   `id` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `b`.`id` AS `x` from `b`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE b (id INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id AS x FROM b
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| b
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| b
: bytes
| CREATE TABLE `b` (
|   `id` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `b`.`id` AS `x` from `b`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a TIME(2) DEFAULT '1:2:3')
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` time(2) DEFAULT '01:02:03.00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a TIME(2) DEFAULT '01:02:03.00')
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` time(2) DEFAULT '01:02:03.00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT(5) ZEROFILL DEFAULT 3)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(5) unsigned zerofill DEFAULT 00003
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT(5) UNSIGNED ZEROFILL DEFAULT '00003')
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(5) unsigned zerofill DEFAULT 00003
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (j JSON)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `j` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`j`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (j JSON DEFAULT NULL)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `j` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`j`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a TEXT, UNIQUE KEY u (a))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` text DEFAULT NULL,
|   UNIQUE KEY `u` (`a`) USING HASH
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a TEXT, UNIQUE KEY u (a) USING HASH)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` text DEFAULT NULL,
|   UNIQUE KEY `u` (`a`) USING HASH
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, KEY ka (a), KEY kb (b))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   KEY `ka` (`a`),
|   KEY `kb` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, KEY kb (b), KEY ka (a))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   KEY `kb` (`b`),
|   KEY `ka` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE SEQUENCE s START WITH 1 INCREMENT BY 1 MINVALUE 1 CACHE 1000 NOCYCLE
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT) WITH SYSTEM VERSIONING
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT WITH SYSTEM VERSIONING)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT COMMENT 'x')
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL COMMENT 'x'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a BIGINT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` bigint(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT DEFAULT 1)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT 1
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT DEFAULT 2)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT 2
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT INVISIBLE)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL INVISIBLE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, UNIQUE KEY k (a))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   UNIQUE KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a) IGNORED)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a DESC))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`a` DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a) USING HASH)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`a`) USING HASH
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT CHECK (a > 0))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT CHECK (a > 1))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL CHECK (`a` > 1)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT c CHECK (a > 0))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `c` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT c CHECK (a > 1))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `c` CHECK (`a` > 1)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `f` (`a`),
|   CONSTRAINT `f` FOREIGN KEY (`a`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id) ON DELETE CASCADE)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `f` (`a`),
|   CONSTRAINT `f` FOREIGN KEY (`a`) REFERENCES `p` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT) COMMENT 'x'
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='x'
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT) ENGINE=MyISAM
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT) STATS_PERSISTENT=0
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci STATS_PERSISTENT=0
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10)) COLLATE utf8mb4_bin
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10)) COLLATE utf8mb4_unicode_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT) WITH SYSTEM VERSIONING
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT) WITH SYSTEM VERSIONING
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT WITHOUT SYSTEM VERSIONING) WITH SYSTEM VERSIONING
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL WITHOUT SYSTEM VERSIONING
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE SEQUENCE s INCREMENT BY 2
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 2 cache 1000 nocycle ENGINE=InnoDB
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE SEQUENCE s COMMENT 'x'
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB COMMENT='x'
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE b (id INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM b
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| b
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| b
: bytes
| CREATE TABLE `b` (
|   `id` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `b`.`id` AS `id` from `b`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE b (id INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id + 1 AS id FROM b
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| b
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| b
: bytes
| CREATE TABLE `b` (
|   `id` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `b`.`id` + 1 AS `id` from `b`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE b (id INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM b
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| b
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| b
: bytes
| CREATE TABLE `b` (
|   `id` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `b`.`id` AS `id` from `b`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE b (id INT)
< ok 0 0
> exec
| CREATE ALGORITHM=MERGE VIEW v AS SELECT id FROM b
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| b
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| b
: bytes
| CREATE TABLE `b` (
|   `id` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=MERGE DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `b`.`id` AS `id` from `b`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 2
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY HASH (`a`)
| PARTITIONS 2
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 4
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY HASH (`a`)
| PARTITIONS 4
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `f` (`a`),
|   CONSTRAINT `f` FOREIGN KEY (`a`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id) ON DELETE RESTRICT ON UPDATE RESTRICT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `f` (`a`),
|   CONSTRAINT `f` FOREIGN KEY (`a`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `f` (`a`),
|   CONSTRAINT `f` FOREIGN KEY (`a`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id) ON DELETE NO ACTION)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `f` (`a`),
|   CONSTRAINT `f` FOREIGN KEY (`a`) REFERENCES `p` (`id`) ON DELETE NO ACTION
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 1
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| e
> query
| SHOW CREATE EVENT `e`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| e
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 24 HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| e
> query
| SHOW CREATE EVENT `e`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| e
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 24 HOUR STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 1
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| e
> query
| SHOW CREATE EVENT `e`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| e
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DISABLE DO SELECT 1
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| e
> query
| SHOW CREATE EVENT `e`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| e
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT) COMMENT 'counter AUTO_INCREMENT=10'
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='counter AUTO_INCREMENT=10'
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT) COMMENT 'counter AUTO_INCREMENT=20'
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='counter AUTO_INCREMENT=20'
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY LENGTH('01') HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| e
> query
| SHOW CREATE EVENT `e`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| e
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 2 HOUR STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY LENGTH('1') HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| e
> query
| SHOW CREATE EVENT `e`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| e
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 1
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| e
> query
| SHOW CREATE EVENT `e`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| e
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 2
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| e
> query
| SHOW CREATE EVENT `e`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| e
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 2
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
< row
: bytes
| tr
: bytes
| BEFORE
| INSERT
| t
| 1
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
| SET NEW.a = 1
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 2
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
< row
: bytes
| tr
: bytes
| BEFORE
| INSERT
| t
| 1
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
| SET NEW.a = 2
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE FUNCTION f(x INT) RETURNS INT DETERMINISTIC RETURN x
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| f
> query
| SHOW CREATE FUNCTION `f`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f`(x INT) RETURNS int(11)
|     DETERMINISTIC
| RETURN x
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE FUNCTION f(x INT) RETURNS INT DETERMINISTIC RETURN x + 1
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| f
> query
| SHOW CREATE FUNCTION `f`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f`(x INT) RETURNS int(11)
|     DETERMINISTIC
| RETURN x + 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE FUNCTION f(x INT) RETURNS INT DETERMINISTIC RETURN x
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| f
> query
| SHOW CREATE FUNCTION `f`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f`(x INT) RETURNS int(11)
|     DETERMINISTIC
| RETURN x
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE FUNCTION f(x INT) RETURNS INT NOT DETERMINISTIC RETURN x
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| f
> query
| SHOW CREATE FUNCTION `f`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f`(x INT) RETURNS int(11)
| RETURN x
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, g INT AS (a + 1) VIRTUAL, KEY k (g))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `g` int(11) GENERATED ALWAYS AS (`a` + 1) VIRTUAL,
|   KEY `k` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, g INT AS (a + 1) STORED, KEY k (g))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `g` int(11) GENERATED ALWAYS AS (`a` + 1) STORED,
|   KEY `k` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, g1 INT AS (a + 1) VIRTUAL, g2 INT AS (g1 + 1) VIRTUAL)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `g1` int(11) GENERATED ALWAYS AS (`a` + 1) VIRTUAL,
|   `g2` int(11) GENERATED ALWAYS AS (`g1` + 1) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE x (i INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `x`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| x
: bytes
| CREATE TABLE `x` (
|   `i` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE x (i INT)
< ok 0 0
> exec
| CREATE TABLE t (id INT, s DATE NOT NULL, e DATE NOT NULL, PERIOD FOR p (s, e), PRIMARY KEY (id, p WITHOUT OVERLAPS))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL,
|   PERIOD FOR `p` (`s`, `e`),
|   PRIMARY KEY (`id`,`p` WITHOUT OVERLAPS)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `x`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| x
: bytes
| CREATE TABLE `x` (
|   `i` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (c INT, d INT)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, KEY k (a), CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   KEY `k` (`a`),
|   CONSTRAINT `f` FOREIGN KEY (`a`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, KEY k (a, b), CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   KEY `k` (`a`,`b`),
|   CONSTRAINT `f` FOREIGN KEY (`a`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (s DATE NOT NULL, e DATE NOT NULL)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `s` date NOT NULL,
|   `e` date NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_diffmatchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_diffmatchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (s DATE NOT NULL, e DATE NOT NULL, PERIOD FOR p (s, e))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_diffmatchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `s` date NOT NULL,
|   `e` date NOT NULL,
|   PERIOD FOR `p` (`s`, `e`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'
< result
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT TRIGGER_NAME, CONCAT_WS('\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT)" "LONGTEXT"
> query
| SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()
< result
< column "EVENT_NAME" "VARCHAR"
