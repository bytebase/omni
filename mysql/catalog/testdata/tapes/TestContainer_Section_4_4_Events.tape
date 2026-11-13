-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_4_4_Events default#1
> exec
| DROP EVENT IF EXISTS ev_test_every
< ok 0 0
> exec
| CREATE EVENT ev_test_every ON SCHEDULE EVERY 1 HOUR DO SELECT 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_test_every
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
| ev_test_every
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_test_every` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP EVENT IF EXISTS ev_test_at
< ok 0 0
> exec
| CREATE EVENT ev_test_at ON SCHEDULE AT '2035-01-01 00:00:00' DO SELECT 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_test_at
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
| ev_test_at
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_test_at` ON SCHEDULE AT '2035-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP EVENT IF EXISTS ev_test_opts
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_trigger
< ok 0 0
> exec
| CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))
< ok 0 0
> exec
| CREATE EVENT ev_test_opts ON SCHEDULE EVERY 1 DAY ON COMPLETION PRESERVE ENABLE COMMENT 'daily cleanup' DO DELETE FROM t_trigger WHERE id < 0
< ok 0 0
> query
| SHOW CREATE EVENT ev_test_opts
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
| ev_test_opts
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_test_opts` ON SCHEDULE EVERY 1 DAY STARTS '2026-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE COMMENT 'daily cleanup' DO DELETE FROM t_trigger WHERE id < 0
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP EVENT IF EXISTS ev_test_disabled
< ok 0 0
> exec
| CREATE EVENT ev_test_disabled ON SCHEDULE EVERY 1 HOUR DISABLE DO SELECT 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_test_disabled
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
| ev_test_disabled
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_test_disabled` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP EVENT IF EXISTS ev_alter_test
< ok 0 0
> exec
| CREATE EVENT ev_alter_test ON SCHEDULE EVERY 1 HOUR DO SELECT 1
< ok 0 0
> exec
| ALTER EVENT ev_alter_test DISABLE
< ok 0 0
> query
| SHOW CREATE EVENT ev_alter_test
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
| ev_alter_test
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_alter_test` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP EVENT IF EXISTS ev_rename_old
< ok 0 0
> exec
| DROP EVENT IF EXISTS ev_rename_new
< ok 0 0
> exec
| CREATE EVENT ev_rename_old ON SCHEDULE EVERY 1 HOUR DO SELECT 1
< ok 0 0
> exec
| ALTER EVENT ev_rename_old RENAME TO ev_rename_new
< ok 0 0
> query
| SHOW CREATE EVENT ev_rename_new
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
| ev_rename_new
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_rename_new` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP EVENT IF EXISTS ev_drop_test
< ok 0 0
> exec
| CREATE EVENT ev_drop_test ON SCHEDULE EVERY 1 HOUR DO SELECT 1
< ok 0 0
> exec
| DROP EVENT ev_drop_test
< ok 0 0
> exec
| DROP EVENT IF EXISTS ev_nonexistent
< ok 0 0
> exec
| DROP EVENT IF EXISTS ev_no_exist
< ok 0 0
> exec
| DROP EVENT ev_no_exist
< error
| 1539 HY000 Unknown event 'ev_no_exist'
> exec
| DROP EVENT IF EXISTS ev_show_test
< ok 0 0
> exec
| CREATE EVENT ev_show_test ON SCHEDULE EVERY 1 HOUR COMMENT 'test event' DO SELECT 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_show_test
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
| ev_show_test
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_show_test` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'test event' DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
