-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
< column "Created" "TIMESTAMP"
< row
: bytes
| tr_before_insert
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `tr_before_insert` BEFORE INSERT ON `t_trigger` FOR EACH ROW SET NEW.val = NEW.val + 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: time 2026-01-01T00:00:00Z
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
< column "Created" "TIMESTAMP"
< row
: bytes
| tr_after_update
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `tr_after_update` AFTER UPDATE ON `t_trigger` FOR EACH ROW SET @updated = @updated + 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: time 2026-01-01T00:00:00Z
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
< column "Created" "TIMESTAMP"
< row
: bytes
| tr_before_delete
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `tr_before_delete` BEFORE DELETE ON `t_trigger` FOR EACH ROW SET @deleted = OLD.id
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: time 2026-01-01T00:00:00Z
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
< column "Created" "TIMESTAMP"
< row
: bytes
| tr_show_test
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `tr_show_test` AFTER INSERT ON `t_trigger` FOR EACH ROW SET @last_insert = NEW.id
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: time 2026-01-01T00:00:00Z
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
< column "Created" "TIMESTAMP"
< row
: bytes
| tr_multi_1
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `tr_multi_1` BEFORE INSERT ON `t_trigger` FOR EACH ROW SET NEW.val = NEW.val + 10
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: time 2026-01-01T00:00:00Z
> query
| SHOW CREATE TRIGGER tr_multi_2
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
| tr_multi_2
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `tr_multi_2` BEFORE INSERT ON `t_trigger` FOR EACH ROW SET NEW.val = NEW.val + 20
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: time 2026-01-01T00:00:00Z
