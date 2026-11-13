-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_4_3_Triggers default#1
> exec
| DROP TABLE IF EXISTS t_trigger
< ok 0 0
> exec
| CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr_before_insert
< ok 0 0
> exec
| CREATE TRIGGER tr_before_insert BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = NEW.val + 1
< ok 0 0
> query
| SHOW CREATE TRIGGER tr_before_insert
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr_before_insert
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER tr_before_insert BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = NEW.val + 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
> exec
| DROP TRIGGER IF EXISTS tr_after_update
< ok 0 0
> exec
| CREATE TRIGGER tr_after_update AFTER UPDATE ON t_trigger FOR EACH ROW SET @updated = @updated + 1
< ok 0 0
> query
| SHOW CREATE TRIGGER tr_after_update
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr_after_update
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER tr_after_update AFTER UPDATE ON t_trigger FOR EACH ROW SET @updated = @updated + 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
> exec
| DROP TRIGGER IF EXISTS tr_before_delete
< ok 0 0
> exec
| CREATE TRIGGER tr_before_delete BEFORE DELETE ON t_trigger FOR EACH ROW SET @deleted = OLD.id
< ok 0 0
> query
| SHOW CREATE TRIGGER tr_before_delete
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr_before_delete
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER tr_before_delete BEFORE DELETE ON t_trigger FOR EACH ROW SET @deleted = OLD.id
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
> exec
| DROP TRIGGER IF EXISTS tr_drop_test
< ok 0 0
> exec
| CREATE TRIGGER tr_drop_test BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = 0
< ok 0 0
> exec
| DROP TRIGGER tr_drop_test
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr_nonexistent_drop
< ok 0 0
> exec ^
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr_no_exist
< ok 0 0
> exec
| DROP TRIGGER tr_no_exist
< error
| 1360 HY000 Trigger does not exist
> exec
| DROP TRIGGER IF EXISTS tr_show_test
< ok 0 0
> exec
| CREATE TRIGGER tr_show_test AFTER INSERT ON t_trigger FOR EACH ROW SET @last_insert = NEW.id
< ok 0 0
> query
| SHOW CREATE TRIGGER tr_show_test
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr_show_test
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER tr_show_test AFTER INSERT ON t_trigger FOR EACH ROW SET @last_insert = NEW.id
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
> exec
| DROP TRIGGER IF EXISTS tr_multi_2
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr_multi_1
< ok 0 0
> exec
| CREATE TRIGGER tr_multi_1 BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = NEW.val + 10
< ok 0 0
> exec
| CREATE TRIGGER tr_multi_2 BEFORE INSERT ON t_trigger FOR EACH ROW FOLLOWS tr_multi_1 SET NEW.val = NEW.val + 20
< ok 0 0
> query
| SHOW CREATE TRIGGER tr_multi_1
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr_multi_1
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER tr_multi_1 BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = NEW.val + 10
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
> query
| SHOW CREATE TRIGGER tr_multi_2
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr_multi_2
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER tr_multi_2 BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = NEW.val + 20
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
