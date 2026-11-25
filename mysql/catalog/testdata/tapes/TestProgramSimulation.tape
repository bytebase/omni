-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestProgramSimulation default#1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER PROCEDURE utf8mb4;
< error
| 1305 42000 PROCEDURE omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP EVENT utf8mb4;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP FUNCTION utf8mb4;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP PROCEDURE utf8mb4;
< error
| 1305 42000 PROCEDURE omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TRIGGER utf8mb4;
< error
| 1360 HY000 Trigger does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT IF NOT EXISTS utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ENABLE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DISABLE ON SLAVE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DISABLE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE AT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ENABLE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE ON SLAVE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY STARTS utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY ENDS utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ON COMPLETION PRESERVE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON COMPLETION PRESERVE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON COMPLETION NOT PRESERVE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 COMMENT 's';
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO COMMIT;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : LOOP COMMIT; END LOOP;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO LOOP COMMIT; END LOOP;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 COMMENT 's';
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() DETERMINISTIC COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
|     DETERMINISTIC
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 LANGUAGE SQL;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 NO SQL;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 CONTAINS SQL;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 READS SQL DATA;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 MODIFIES SQL DATA;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 SQL SECURITY DEFINER;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() COMMENT 's' COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
|     COMMENT 's'
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() NOT DETERMINISTIC COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 SQL SECURITY INVOKER;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(IN utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(IN utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(OUT utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(OUT utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(INOUT utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(INOUT utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO LOOP COMMIT; COMMIT; END LOOP;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 INT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CONDITION FOR 1; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR 1 COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CURSOR FOR SELECT utf8mb4; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CURSOR FOR SELECT utf8mb4 zzz; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE CONTINUE HANDLER FOR 1 COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR SQLWARNING COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR NOT FOUND COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 INT DEFAULT utf8mb4; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() IF utf8mb4 THEN COMMIT; END IF;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| IF utf8mb4 THEN COMMIT; END IF
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() CASE WHEN utf8mb4 THEN COMMIT; END CASE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| CASE WHEN utf8mb4 THEN COMMIT; END CASE
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() utf8mb4 : BEGIN END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| utf8mb4 : BEGIN END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() BEGIN END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| BEGIN END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() utf8mb4 : LOOP COMMIT; END LOOP;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| utf8mb4 : LOOP COMMIT; END LOOP
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() LOOP COMMIT; END LOOP;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| LOOP COMMIT; END LOOP
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; ELSEIF utf8mb4 THEN COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; ELSE COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE utf8mb4 WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE utf8mb4 WHEN utf8mb4 THEN COMMIT; WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; ELSE COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END utf8mb4;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO WHILE utf8mb4 DO COMMIT; END WHILE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO REPEAT COMMIT; UNTIL utf8mb4 END REPEAT;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ON COMPLETION PRESERVE;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 RENAME TO utf8mb4;
< error
| 1551 HY000 Same old and new event name
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP FUNCTION utf8mb4.utf8mb4;
< error
| 1305 42000 FUNCTION utf8mb4.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS utf8mb4.utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP EVENT IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 PROCEDURE utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`utf8mb4`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE IF NOT EXISTS utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER PROCEDURE utf8mb4;
< error
| 1305 42000 PROCEDURE omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP EVENT utf8mb4;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP FUNCTION utf8mb4;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP PROCEDURE utf8mb4;
< error
| 1305 42000 PROCEDURE omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TRIGGER utf8mb4;
< error
| 1360 HY000 Trigger does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT IF NOT EXISTS utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ENABLE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DISABLE ON SLAVE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DISABLE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE AT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ENABLE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE ON SLAVE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY STARTS utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY ENDS utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ON COMPLETION PRESERVE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON COMPLETION PRESERVE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON COMPLETION NOT PRESERVE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 COMMENT 's';
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO COMMIT;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : LOOP COMMIT; END LOOP;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO LOOP COMMIT; END LOOP;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 COMMENT 's';
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() DETERMINISTIC COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
|     DETERMINISTIC
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 LANGUAGE SQL;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 NO SQL;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 CONTAINS SQL;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 READS SQL DATA;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 MODIFIES SQL DATA;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 SQL SECURITY DEFINER;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() COMMENT 's' COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
|     COMMENT 's'
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() NOT DETERMINISTIC COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 SQL SECURITY INVOKER;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(IN utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(IN utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(OUT utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(OUT utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(INOUT utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(INOUT utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO LOOP COMMIT; COMMIT; END LOOP;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 INT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CONDITION FOR 1; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR 1 COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CURSOR FOR SELECT utf8mb4; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CURSOR FOR SELECT utf8mb4 zzz; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE CONTINUE HANDLER FOR 1 COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR SQLWARNING COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR NOT FOUND COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 INT DEFAULT utf8mb4; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() IF utf8mb4 THEN COMMIT; END IF;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| IF utf8mb4 THEN COMMIT; END IF
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() CASE WHEN utf8mb4 THEN COMMIT; END CASE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| CASE WHEN utf8mb4 THEN COMMIT; END CASE
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() utf8mb4 : BEGIN END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| utf8mb4 : BEGIN END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() BEGIN END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| BEGIN END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() utf8mb4 : LOOP COMMIT; END LOOP;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| utf8mb4 : LOOP COMMIT; END LOOP
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() LOOP COMMIT; END LOOP;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| LOOP COMMIT; END LOOP
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; ELSEIF utf8mb4 THEN COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; ELSE COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE utf8mb4 WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE utf8mb4 WHEN utf8mb4 THEN COMMIT; WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; ELSE COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END utf8mb4;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO WHILE utf8mb4 DO COMMIT; END WHILE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO REPEAT COMMIT; UNTIL utf8mb4 END REPEAT;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ON COMPLETION PRESERVE;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT utf8mb4 RENAME TO utf8mb4;
< error
| 1551 HY000 Same old and new event name
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP FUNCTION utf8mb4.utf8mb4;
< error
| 1305 42000 FUNCTION utf8mb4.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS utf8mb4.utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP EVENT IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 PROCEDURE utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`utf8mb4`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE PROCEDURE IF NOT EXISTS utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER PROCEDURE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| DROP EVENT utf8mb4;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| DROP FUNCTION utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| DROP PROCEDURE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| DROP TRIGGER utf8mb4;
< error
| 1360 HY000 Trigger does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE EVENT IF NOT EXISTS utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ENABLE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DISABLE ON SLAVE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DISABLE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE AT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ENABLE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE ON SLAVE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY STARTS utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY ENDS utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ON COMPLETION PRESERVE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON COMPLETION PRESERVE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON COMPLETION NOT PRESERVE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 COMMENT 's';
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO COMMIT;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : LOOP COMMIT; END LOOP;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO LOOP COMMIT; END LOOP;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
|     COMMENT 's'
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() COMMIT;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() DETERMINISTIC COMMIT;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 LANGUAGE SQL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 NO SQL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     NO SQL
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 CONTAINS SQL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 READS SQL DATA;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     READS SQL DATA
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 MODIFIES SQL DATA;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     MODIFIES SQL DATA
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 SQL SECURITY DEFINER;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() COMMENT 's' COMMIT;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() NOT DETERMINISTIC COMMIT;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 SQL SECURITY INVOKER;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
|     SQL SECURITY INVOKER
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(utf8mb4 INT) COMMIT;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(IN utf8mb4 INT) COMMIT;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(OUT utf8mb4 INT) COMMIT;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(INOUT utf8mb4 INT) COMMIT;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO LOOP COMMIT; COMMIT; END LOOP;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 INT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CONDITION FOR 1; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR 1 COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CURSOR FOR SELECT utf8mb4; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CURSOR FOR SELECT utf8mb4 zzz; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE CONTINUE HANDLER FOR 1 COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR SQLWARNING COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR NOT FOUND COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 INT DEFAULT utf8mb4; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() IF utf8mb4 THEN COMMIT; END IF;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() CASE WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() utf8mb4 : BEGIN END;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() BEGIN END;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() utf8mb4 : LOOP COMMIT; END LOOP;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() LOOP COMMIT; END LOOP;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; ELSEIF utf8mb4 THEN COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; ELSE COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE utf8mb4 WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE utf8mb4 WHEN utf8mb4 THEN COMMIT; WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; ELSE COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END utf8mb4;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO WHILE utf8mb4 DO COMMIT; END WHILE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO REPEAT COMMIT; UNTIL utf8mb4 END REPEAT;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ON COMPLETION PRESERVE;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 RENAME TO utf8mb4;
< error
| 1551 HY000 Same old and new event name
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| DROP FUNCTION utf8mb4.utf8mb4;
< error
| 1305 42000 FUNCTION utf8mb4.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS utf8mb4.utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| DROP EVENT IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 PROCEDURE utf8mb4() COMMIT;
< error
| 1304 42000 PROCEDURE utf8mb4 already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE PROCEDURE utf8mb4() COMMIT;
| CREATE FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE PROCEDURE IF NOT EXISTS utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| FUNCTION
: bytes
| utf8mb4
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE FUNCTION `utf8mb4`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `utf8mb4`() RETURNS int
|     DETERMINISTIC
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT;
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER PROCEDURE utf8mb4;
< error
| 1305 42000 PROCEDURE omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| DROP EVENT utf8mb4;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| DROP FUNCTION utf8mb4;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| DROP PROCEDURE utf8mb4;
< error
| 1305 42000 PROCEDURE omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| DROP TRIGGER utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE EVENT IF NOT EXISTS utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ENABLE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DISABLE ON SLAVE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DISABLE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE AT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ENABLE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE ON SLAVE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY STARTS utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY ENDS utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ON COMPLETION PRESERVE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON COMPLETION PRESERVE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON COMPLETION NOT PRESERVE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 COMMENT 's';
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO COMMIT;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : LOOP COMMIT; END LOOP;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO LOOP COMMIT; END LOOP;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 COMMENT 's';
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() DETERMINISTIC COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
|     DETERMINISTIC
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 LANGUAGE SQL;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 NO SQL;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 CONTAINS SQL;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 READS SQL DATA;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 MODIFIES SQL DATA;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 SQL SECURITY DEFINER;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() COMMENT 's' COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
|     COMMENT 's'
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() NOT DETERMINISTIC COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 SQL SECURITY INVOKER;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(IN utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(IN utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(OUT utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(OUT utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(INOUT utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(INOUT utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO LOOP COMMIT; COMMIT; END LOOP;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 INT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CONDITION FOR 1; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR 1 COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CURSOR FOR SELECT utf8mb4; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CURSOR FOR SELECT utf8mb4 zzz; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE CONTINUE HANDLER FOR 1 COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR SQLWARNING COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR NOT FOUND COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION COMMIT; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 INT DEFAULT utf8mb4; END;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() IF utf8mb4 THEN COMMIT; END IF;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| IF utf8mb4 THEN COMMIT; END IF
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() CASE WHEN utf8mb4 THEN COMMIT; END CASE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| CASE WHEN utf8mb4 THEN COMMIT; END CASE
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() utf8mb4 : BEGIN END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| utf8mb4 : BEGIN END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() BEGIN END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| BEGIN END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() utf8mb4 : LOOP COMMIT; END LOOP;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| utf8mb4 : LOOP COMMIT; END LOOP
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() LOOP COMMIT; END LOOP;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| LOOP COMMIT; END LOOP
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; ELSEIF utf8mb4 THEN COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; ELSE COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; COMMIT; END IF;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE utf8mb4 WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE utf8mb4 WHEN utf8mb4 THEN COMMIT; WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; WHEN utf8mb4 THEN COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; ELSE COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; COMMIT; END CASE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END utf8mb4;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO WHILE utf8mb4 DO COMMIT; END WHILE;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO REPEAT COMMIT; UNTIL utf8mb4 END REPEAT;
< error
| 1539 HY000 Unknown event 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ON COMPLETION PRESERVE;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| ALTER EVENT utf8mb4 RENAME TO utf8mb4;
< error
| 1551 HY000 Same old and new event name
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| DROP FUNCTION utf8mb4.utf8mb4;
< error
| 1305 42000 FUNCTION utf8mb4.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS utf8mb4.utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| DROP EVENT IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 PROCEDURE utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`utf8mb4`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT);
| CREATE TRIGGER utf8mb4 BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1
< ok 0 0
> exec
| CREATE PROCEDURE IF NOT EXISTS utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
< row
: bytes
| TRIGGER
: bytes
| utf8mb4
> query
| SHOW CREATE TRIGGER `utf8mb4`
< result
< column "Trigger" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< column "Created" "TIMESTAMP" precision 0 scale 0
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `utf8mb4` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.c = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 2026-01-01 00:00:00.00
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER PROCEDURE utf8mb4;
< error
| 1305 42000 PROCEDURE omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP EVENT utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP FUNCTION utf8mb4;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP PROCEDURE utf8mb4;
< error
| 1305 42000 PROCEDURE omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP TRIGGER utf8mb4;
< error
| 1360 HY000 Trigger does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE EVENT IF NOT EXISTS utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ENABLE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DISABLE ON SLAVE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DISABLE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE AT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ENABLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE ON SLAVE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE ON SLAVE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY STARTS utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY utf8mb4 DAY ENDS utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ON COMPLETION PRESERVE DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON COMPLETION PRESERVE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON COMPLETION NOT PRESERVE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE COMMENT 's' DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; END IF;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO IF utf8mb4 THEN COMMIT; END IF
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; END CASE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO CASE WHEN utf8mb4 THEN COMMIT; END CASE
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO utf8mb4 : BEGIN END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : LOOP COMMIT; END LOOP;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO utf8mb4 : LOOP COMMIT; END LOOP
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO LOOP COMMIT; END LOOP;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO LOOP COMMIT; END LOOP
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 COMMENT 's';
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() DETERMINISTIC COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
|     DETERMINISTIC
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 LANGUAGE SQL;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 NO SQL;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 CONTAINS SQL;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 READS SQL DATA;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 MODIFIES SQL DATA;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 SQL SECURITY DEFINER;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() COMMENT 's' COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
|     COMMENT 's'
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() NOT DETERMINISTIC COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER FUNCTION utf8mb4 SQL SECURITY INVOKER;
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(IN utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(IN utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(OUT utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(OUT utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4(INOUT utf8mb4 INT) COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`(INOUT utf8mb4 INT)
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN COMMIT; END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN COMMIT; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO LOOP COMMIT; COMMIT; END LOOP;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO LOOP COMMIT; COMMIT; END LOOP
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 INT; END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN DECLARE utf8mb4 INT; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CONDITION FOR 1; END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN DECLARE utf8mb4 CONDITION FOR 1; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR 1 COMMIT; END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN DECLARE EXIT HANDLER FOR 1 COMMIT; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CURSOR FOR SELECT utf8mb4; END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN DECLARE utf8mb4 CURSOR FOR SELECT utf8mb4; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 CURSOR FOR SELECT utf8mb4 zzz; END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN DECLARE utf8mb4 CURSOR FOR SELECT utf8mb4 zzz; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE CONTINUE HANDLER FOR 1 COMMIT; END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN DECLARE CONTINUE HANDLER FOR 1 COMMIT; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR SQLWARNING COMMIT; END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN DECLARE EXIT HANDLER FOR SQLWARNING COMMIT; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR NOT FOUND COMMIT; END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN DECLARE EXIT HANDLER FOR NOT FOUND COMMIT; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION COMMIT; END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION COMMIT; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO BEGIN DECLARE utf8mb4 INT DEFAULT utf8mb4; END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN DECLARE utf8mb4 INT DEFAULT utf8mb4; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() IF utf8mb4 THEN COMMIT; END IF;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| IF utf8mb4 THEN COMMIT; END IF
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() CASE WHEN utf8mb4 THEN COMMIT; END CASE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| CASE WHEN utf8mb4 THEN COMMIT; END CASE
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() utf8mb4 : BEGIN END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| utf8mb4 : BEGIN END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() BEGIN END;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| BEGIN END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() utf8mb4 : LOOP COMMIT; END LOOP;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| utf8mb4 : LOOP COMMIT; END LOOP
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE utf8mb4() LOOP COMMIT; END LOOP;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| LOOP COMMIT; END LOOP
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; ELSEIF utf8mb4 THEN COMMIT; END IF;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO IF utf8mb4 THEN COMMIT; ELSEIF utf8mb4 THEN COMMIT; END IF
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; ELSE COMMIT; END IF;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO IF utf8mb4 THEN COMMIT; ELSE COMMIT; END IF
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO IF utf8mb4 THEN COMMIT; COMMIT; END IF;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO IF utf8mb4 THEN COMMIT; COMMIT; END IF
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE utf8mb4 WHEN utf8mb4 THEN COMMIT; END CASE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO CASE utf8mb4 WHEN utf8mb4 THEN COMMIT; END CASE
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE utf8mb4 WHEN utf8mb4 THEN COMMIT; WHEN utf8mb4 THEN COMMIT; END CASE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO CASE utf8mb4 WHEN utf8mb4 THEN COMMIT; WHEN utf8mb4 THEN COMMIT; END CASE
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; WHEN utf8mb4 THEN COMMIT; END CASE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO CASE WHEN utf8mb4 THEN COMMIT; WHEN utf8mb4 THEN COMMIT; END CASE
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; ELSE COMMIT; END CASE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO CASE WHEN utf8mb4 THEN COMMIT; ELSE COMMIT; END CASE
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO CASE WHEN utf8mb4 THEN COMMIT; COMMIT; END CASE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO CASE WHEN utf8mb4 THEN COMMIT; COMMIT; END CASE
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO utf8mb4 : BEGIN END utf8mb4
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO utf8mb4 : BEGIN END
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO utf8mb4 : BEGIN END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO WHILE utf8mb4 DO COMMIT; END WHILE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO WHILE utf8mb4 DO COMMIT; END WHILE
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO REPEAT COMMIT; UNTIL utf8mb4 END REPEAT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO REPEAT COMMIT; UNTIL utf8mb4 END REPEAT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE AT utf8mb4 ON COMPLETION PRESERVE;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 RENAME TO utf8mb4;
< error
| 1551 HY000 Same old and new event name
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP FUNCTION utf8mb4.utf8mb4;
< error
| 1305 42000 FUNCTION utf8mb4.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS utf8mb4.utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP EVENT IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 PROCEDURE utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`utf8mb4`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE PROCEDURE IF NOT EXISTS utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
< row
: bytes
| PROCEDURE
: bytes
| utf8mb4
> query
| SHOW CREATE PROCEDURE `utf8mb4`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `utf8mb4`()
| COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 RENAME TO cleanup
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| cleanup
> query
| SHOW CREATE EVENT `cleanup`
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
| cleanup
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `cleanup` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE AT '2031-01-01 00:00:00'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE AT '2031-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY 2 HOUR STARTS '2030-06-01 00:00:00' ENDS '2031-01-01 00:00:00'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 2 HOUR STARTS '2030-06-01 00:00:00' ENDS '2031-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE ON SLAVE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE ON SLAVE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DISABLE ON REPLICA
< error
| 1064 42000 …near 'REPLICA' at line 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON COMPLETION NOT PRESERVE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 COMMENT 'nightly'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE COMMENT 'nightly' DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 DO SELECT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 ON SCHEDULE EVERY 1 HOUR STARTS '2031-01-01 00:00:00' ENDS '2030-01-01 00:00:00'
< error
| 1543 HY000 ENDS is either invalid or before STARTS
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT utf8mb4 RENAME TO utf8mb4
< error
| 1551 HY000 Same old and new event name
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| ALTER EVENT nope DISABLE
< error
| 1539 HY000 Unknown event 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE EVENT utf8mb4 ON SCHEDULE AT '2030-01-01 00:00:00' DO COMMIT
< error
| 1537 HY000 Event 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| CREATE EVENT IF NOT EXISTS utf8mb4 ON SCHEDULE AT '2030-01-01 00:00:00' DO COMMIT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP EVENT utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP EVENT IF EXISTS nope
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| utf8mb4
> query
| SHOW CREATE EVENT `utf8mb4`
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
| utf8mb4
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `utf8mb4` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE EVENT utf8mb4 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT
< ok 0 0
> exec
| DROP EVENT nope
< error
| 1539 HY000 Unknown event 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT nightly ON SCHEDULE EVERY 1 HOUR STARTS '2030-01-01 00:00:00' ENDS '2031-01-01 00:00:00' DO DELETE FROM t WHERE c < 0
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| nightly
> query
| SHOW CREATE EVENT `nightly`
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
| nightly
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `nightly` ON SCHEDULE EVERY 1 HOUR STARTS '2030-01-01 00:00:00' ENDS '2031-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO DELETE FROM t WHERE c < 0
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT nightly ON SCHEDULE EVERY 1 HOUR STARTS '2031-01-01 00:00:00' ENDS '2030-01-01 00:00:00' DO COMMIT
< error
| 1543 HY000 ENDS is either invalid or before STARTS
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT nightly ON SCHEDULE AT '2030-01-01 00:00:00' + INTERVAL 1 DAY ON COMPLETION PRESERVE DISABLE COMMENT 'later' DO COMMIT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| nightly
> query
| SHOW CREATE EVENT `nightly`
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
| nightly
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `nightly` ON SCHEDULE AT '2030-01-02 00:00:00' ON COMPLETION PRESERVE DISABLE COMMENT 'later' DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT nightly ON SCHEDULE AT '2000-01-01 00:00:00' DO COMMIT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT nightly ON SCHEDULE AT '2000-01-01 00:00:00' ON COMPLETION PRESERVE DO COMMIT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| nightly
> query
| SHOW CREATE EVENT `nightly`
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
| nightly
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `nightly` ON SCHEDULE AT '2000-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT nightly ON SCHEDULE EVERY 0 DAY STARTS '2030-01-01 00:00:00' DO COMMIT
< error
| 1542 HY000 INTERVAL is either not positive or too big
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT nightly ON SCHEDULE EVERY '1:30' HOUR_MINUTE STARTS '2030-01-01 00:00:00' DO COMMIT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| nightly
> query
| SHOW CREATE EVENT `nightly`
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
| nightly
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `nightly` ON SCHEDULE EVERY '1:30' HOUR_MINUTE STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT nightly ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO BEGIN DELETE FROM t; INSERT INTO t VALUES (1, 1, 1); END
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| nightly
> query
| SHOW CREATE EVENT `nightly`
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
| nightly
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `nightly` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO BEGIN DELETE FROM t; INSERT INTO t VALUES (1, 1, 1); END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT nightly ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DISABLE ON SLAVE DO COMMIT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM" "string"
< column "ROUTINE_NAME" "VARCHAR" "string"
> query
| SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME
< result
< column "TRIGGER" "VARCHAR" "string"
< column "TRIGGER_NAME" "VARCHAR" "string"
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR" "string"
< column "EVENT_NAME" "VARCHAR" "string"
< row
: bytes
| EVENT
: bytes
| nightly
> query
| SHOW CREATE EVENT `nightly`
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
| nightly
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `nightly` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE ON SLAVE DO COMMIT
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT nightly ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DISABLE ON REPLICA DO COMMIT
< error
| 1064 42000 …near 'REPLICA DO COMMIT' at line 1
