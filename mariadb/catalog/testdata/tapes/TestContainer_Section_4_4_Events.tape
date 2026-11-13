-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_test_every` ON SCHEDULE EVERY 1 HOUR STARTS '2026-11-13 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
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
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_test_at` ON SCHEDULE AT '2035-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
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
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_test_opts` ON SCHEDULE EVERY 1 DAY STARTS '2026-11-13 00:00:00' ON COMPLETION PRESERVE ENABLE COMMENT 'daily cleanup' DO DELETE FROM t_trigger WHERE id < 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
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
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_test_disabled` ON SCHEDULE EVERY 1 HOUR STARTS '2026-11-13 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
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
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_alter_test` ON SCHEDULE EVERY 1 HOUR STARTS '2026-11-13 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
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
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_rename_new` ON SCHEDULE EVERY 1 HOUR STARTS '2026-11-13 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
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
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_show_test` ON SCHEDULE EVERY 1 HOUR STARTS '2026-11-13 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'test event' DO SELECT 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
