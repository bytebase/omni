-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_EventIdempotence default#1

== TestOracle_EventIdempotence default#2
> exec
| DROP DATABASE IF EXISTS ev_idem_every_default
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_default DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_default
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_default.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_every_explicit_enable
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_explicit_enable DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_explicit_enable
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 SECOND ENABLE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_explicit_enable.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 SECOND STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_every_disable
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_disable DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_disable
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 5 MINUTE DISABLE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_disable.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 5 MINUTE STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_every_disable_on_slave
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_disable_on_slave DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_disable_on_slave
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 30 MINUTE ON COMPLETION NOT PRESERVE DISABLE ON SLAVE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_disable_on_slave.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 30 MINUTE STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE ON SLAVE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_every_preserve
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_preserve DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_preserve
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY ON COMPLETION PRESERVE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_preserve.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 DAY STARTS '2026-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_every_explicit_starts
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_explicit_starts DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_explicit_starts
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 5 MINUTE STARTS '2025-03-01 12:00:00' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_explicit_starts.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 5 MINUTE STARTS '2025-03-01 12:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_every_starts_ends
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_starts_ends DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_starts_ends
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2025-01-01 00:00:00' ENDS '2030-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE COMMENT 'c' DO SET @q = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_starts_ends.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 DAY STARTS '2025-01-01 00:00:00' ENDS '2030-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE COMMENT 'c' DO SET @q = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_every_only_ends
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_only_ends DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_only_ends
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENDS '2030-01-01 00:00:00' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_only_ends.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ENDS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_every_interval_quoted_hm
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_interval_quoted_hm DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_interval_quoted_hm
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY '1:30' HOUR_MINUTE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_interval_quoted_hm.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY '1:30' HOUR_MINUTE STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_every_interval_quoted_leadingzero
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_interval_quoted_leadingzero DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_interval_quoted_leadingzero
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY '02:03:04' HOUR_SECOND DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_interval_quoted_leadingzero.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY '2:3:4' HOUR_SECOND STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_every_interval_quoted_daysecond
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_interval_quoted_daysecond DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_interval_quoted_daysecond
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY '1 12:30:00' DAY_SECOND DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_interval_quoted_daysecond.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY '1 12:30:0' DAY_SECOND STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_every_interval_single_quoted
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_interval_single_quoted DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_interval_single_quoted
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY '05' MINUTE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_interval_single_quoted.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 5 MINUTE STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_every_comment
< ok 0 0
> exec
| CREATE DATABASE ev_idem_every_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_every_comment
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'it''s a test' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_every_comment.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'it''s a test' DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_at_literal
< ok 0 0
> exec
| CREATE DATABASE ev_idem_at_literal DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_at_literal
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE DO SET @y = 2
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_at_literal.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO SET @y = 2
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ev_idem_at_preserve_disable
< ok 0 0
> exec
| CREATE DATABASE ev_idem_at_preserve_disable DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ev_idem_at_preserve_disable
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE AT '2030-06-01 00:00:00' ON COMPLETION PRESERVE DISABLE DO SET @y = 2
< ok 0 0
> query
| SHOW CREATE EVENT ev_idem_at_preserve_disable.e
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE AT '2030-06-01 00:00:00' ON COMPLETION PRESERVE DISABLE DO SET @y = 2
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
