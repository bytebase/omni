-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_6_2_ShowCreateOtherObjects default#1
> exec
| DROP VIEW IF EXISTS v_sc_basic
< ok 0 0
> exec
| CREATE VIEW v_sc_basic AS SELECT 1 AS col1, 2 AS col2
< ok 0 0
> query
| SHOW CREATE VIEW v_sc_basic
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_sc_basic
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_sc_basic` AS select 1 AS `col1`,2 AS `col2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP FUNCTION IF EXISTS fn_sc_multiply
< ok 0 0
> exec
| CREATE FUNCTION fn_sc_multiply(x INT, y INT) RETURNS INT DETERMINISTIC RETURN x * y
< ok 0 0
> query
| SHOW CREATE FUNCTION fn_sc_multiply
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_sc_multiply
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_sc_multiply`(x INT, y INT) RETURNS int
|     DETERMINISTIC
| RETURN x * y
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP PROCEDURE IF EXISTS sp_sc_swap
< ok 0 0
> exec
| CREATE PROCEDURE sp_sc_swap(INOUT a INT, INOUT b INT) BEGIN DECLARE tmp INT; SET tmp = a; SET a = b; SET b = tmp; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE sp_sc_swap
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| sp_sc_swap
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `sp_sc_swap`(INOUT a INT, INOUT b INT)
| BEGIN DECLARE tmp INT; SET tmp = a; SET a = b; SET b = tmp; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP TRIGGER IF EXISTS tr_sc_audit
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_sc_trigger
< ok 0 0
> exec
| CREATE TABLE t_sc_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, updated_at DATETIME)
< ok 0 0
> exec
| CREATE TRIGGER tr_sc_audit AFTER INSERT ON t_sc_trigger FOR EACH ROW SET @last_insert_id = NEW.id
< ok 0 0
> query
| SHOW CREATE TRIGGER tr_sc_audit
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "TIMESTAMP"
< row
: bytes
| tr_sc_audit
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `tr_sc_audit` AFTER INSERT ON `t_sc_trigger` FOR EACH ROW SET @last_insert_id = NEW.id
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: time 2026-01-01T00:00:00Z
> exec
| DROP EVENT IF EXISTS ev_sc_onetime
< ok 0 0
> exec
| CREATE EVENT ev_sc_onetime ON SCHEDULE AT '2035-06-15 12:00:00' ON COMPLETION PRESERVE COMMENT 'one-time event' DO SELECT 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_sc_onetime
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
| ev_sc_onetime
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_sc_onetime` ON SCHEDULE AT '2035-06-15 12:00:00' ON COMPLETION PRESERVE ENABLE COMMENT 'one-time event' DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
