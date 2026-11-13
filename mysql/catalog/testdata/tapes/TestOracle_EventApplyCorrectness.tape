-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_EventApplyCorrectness default#1

== TestOracle_EventApplyCorrectness/8.0/alter-add-comment default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_add_comment
< ok 0 0
> exec
| CREATE DATABASE evto_alter_add_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_add_comment
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'new' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_add_comment.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'new' DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS evfrom_alter_add_comment
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_add_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_add_comment
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_add_comment.e
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
| DROP DATABASE IF EXISTS eva_alter_add_comment
< ok 0 0
> exec
| CREATE DATABASE eva_alter_add_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_add_comment
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'new' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_add_comment.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'new' DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci

== TestOracle_EventApplyCorrectness/8.0/alter-add-ends default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_add_ends
< ok 0 0
> exec
| CREATE DATABASE evto_alter_add_ends DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_add_ends
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENDS '2030-01-01 00:00:00' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_add_ends.e
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
| DROP DATABASE IF EXISTS evfrom_alter_add_ends
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_add_ends DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_add_ends
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_add_ends.e
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
| DROP DATABASE IF EXISTS eva_alter_add_ends
< ok 0 0
> exec
| CREATE DATABASE eva_alter_add_ends DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_add_ends
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ENDS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_add_ends.e
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

== TestOracle_EventApplyCorrectness/8.0/alter-at-to-every default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_at_to_every
< ok 0 0
> exec
| CREATE DATABASE evto_alter_at_to_every DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_at_to_every
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_at_to_every.e
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
| DROP DATABASE IF EXISTS evfrom_alter_at_to_every
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_at_to_every DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_at_to_every
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_at_to_every.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS eva_alter_at_to_every
< ok 0 0
> exec
| CREATE DATABASE eva_alter_at_to_every DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_at_to_every
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_at_to_every.e
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

== TestOracle_EventApplyCorrectness/8.0/alter-body default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_body
< ok 0 0
> exec
| CREATE DATABASE evto_alter_body DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_body
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 2
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_body.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 2
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS evfrom_alter_body
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_body DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_body
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_body.e
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
| DROP DATABASE IF EXISTS eva_alter_body
< ok 0 0
> exec
| CREATE DATABASE eva_alter_body DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_body
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 2
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_body.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 2
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci

== TestOracle_EventApplyCorrectness/8.0/alter-comment default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_comment
< ok 0 0
> exec
| CREATE DATABASE evto_alter_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_comment
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'b' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_comment.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'b' DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS evfrom_alter_comment
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_comment
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'a' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_comment.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'a' DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS eva_alter_comment
< ok 0 0
> exec
| CREATE DATABASE eva_alter_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_comment
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'a' DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'b' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_comment.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'b' DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci

== TestOracle_EventApplyCorrectness/8.0/alter-completion default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_completion
< ok 0 0
> exec
| CREATE DATABASE evto_alter_completion DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_completion
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ON COMPLETION PRESERVE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_completion.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS evfrom_alter_completion
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_completion DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_completion
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_completion.e
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
| DROP DATABASE IF EXISTS eva_alter_completion
< ok 0 0
> exec
| CREATE DATABASE eva_alter_completion DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_completion
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_completion.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci

== TestOracle_EventApplyCorrectness/8.0/alter-completion-reset default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_completion_reset
< ok 0 0
> exec
| CREATE DATABASE evto_alter_completion_reset DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_completion_reset
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_completion_reset.e
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
| DROP DATABASE IF EXISTS evfrom_alter_completion_reset
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_completion_reset DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_completion_reset
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ON COMPLETION PRESERVE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_completion_reset.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS eva_alter_completion_reset
< ok 0 0
> exec
| CREATE DATABASE eva_alter_completion_reset DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_completion_reset
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ON COMPLETION PRESERVE DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_completion_reset.e
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

== TestOracle_EventApplyCorrectness/8.0/alter-every-to-at default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_every_to_at
< ok 0 0
> exec
| CREATE DATABASE evto_alter_every_to_at DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_every_to_at
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_every_to_at.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS evfrom_alter_every_to_at
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_every_to_at DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_every_to_at
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_every_to_at.e
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
| DROP DATABASE IF EXISTS eva_alter_every_to_at
< ok 0 0
> exec
| CREATE DATABASE eva_alter_every_to_at DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_every_to_at
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_every_to_at.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci

== TestOracle_EventApplyCorrectness/8.0/alter-remove-comment default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_remove_comment
< ok 0 0
> exec
| CREATE DATABASE evto_alter_remove_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_remove_comment
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_remove_comment.e
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
| DROP DATABASE IF EXISTS evfrom_alter_remove_comment
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_remove_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_remove_comment
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'old' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_remove_comment.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'old' DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS eva_alter_remove_comment
< ok 0 0
> exec
| CREATE DATABASE eva_alter_remove_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_remove_comment
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'old' DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_remove_comment.e
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

== TestOracle_EventApplyCorrectness/8.0/alter-remove-ends default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_remove_ends
< ok 0 0
> exec
| CREATE DATABASE evto_alter_remove_ends DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_remove_ends
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_remove_ends.e
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
| DROP DATABASE IF EXISTS evfrom_alter_remove_ends
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_remove_ends DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_remove_ends
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENDS '2030-01-01 00:00:00' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_remove_ends.e
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
| DROP DATABASE IF EXISTS eva_alter_remove_ends
< ok 0 0
> exec
| CREATE DATABASE eva_alter_remove_ends DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_remove_ends
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENDS '2030-01-01 00:00:00' DO SET @x = 1
< ok 0 0
> exec
| DROP EVENT IF EXISTS `e`
< ok 0 0
> exec
| CREATE EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_remove_ends.e
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

== TestOracle_EventApplyCorrectness/8.0/alter-schedule-interval default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_schedule_interval
< ok 0 0
> exec
| CREATE DATABASE evto_alter_schedule_interval DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_schedule_interval
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 2 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_schedule_interval.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 2 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS evfrom_alter_schedule_interval
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_schedule_interval DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_schedule_interval
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_schedule_interval.e
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
| DROP DATABASE IF EXISTS eva_alter_schedule_interval
< ok 0 0
> exec
| CREATE DATABASE eva_alter_schedule_interval DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_schedule_interval
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 2 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_schedule_interval.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 2 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci

== TestOracle_EventApplyCorrectness/8.0/alter-schedule-unit default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_schedule_unit
< ok 0 0
> exec
| CREATE DATABASE evto_alter_schedule_unit DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_schedule_unit
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_schedule_unit.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 DAY STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS evfrom_alter_schedule_unit
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_schedule_unit DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_schedule_unit
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_schedule_unit.e
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
| DROP DATABASE IF EXISTS eva_alter_schedule_unit
< ok 0 0
> exec
| CREATE DATABASE eva_alter_schedule_unit DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_schedule_unit
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 1 DAY STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_schedule_unit.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 DAY STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci

== TestOracle_EventApplyCorrectness/8.0/alter-status-disable default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_status_disable
< ok 0 0
> exec
| CREATE DATABASE evto_alter_status_disable DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_status_disable
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_status_disable.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS evfrom_alter_status_disable
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_status_disable DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_status_disable
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENABLE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_status_disable.e
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
| DROP DATABASE IF EXISTS eva_alter_status_disable
< ok 0 0
> exec
| CREATE DATABASE eva_alter_status_disable DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_status_disable
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENABLE DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_status_disable.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci

== TestOracle_EventApplyCorrectness/8.0/alter-status-reenable default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_status_reenable
< ok 0 0
> exec
| CREATE DATABASE evto_alter_status_reenable DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_status_reenable
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENABLE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_status_reenable.e
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
| DROP DATABASE IF EXISTS evfrom_alter_status_reenable
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_status_reenable DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_status_reenable
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_status_reenable.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS eva_alter_status_reenable
< ok 0 0
> exec
| CREATE DATABASE eva_alter_status_reenable DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_status_reenable
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_status_reenable.e
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

== TestOracle_EventApplyCorrectness/8.0/alter-status-slave default#1
> exec
| DROP DATABASE IF EXISTS evto_alter_status_slave
< ok 0 0
> exec
| CREATE DATABASE evto_alter_status_slave DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_alter_status_slave
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE ON SLAVE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_alter_status_slave.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE ON SLAVE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS evfrom_alter_status_slave
< ok 0 0
> exec
| CREATE DATABASE evfrom_alter_status_slave DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_alter_status_slave
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_alter_status_slave.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS eva_alter_status_slave
< ok 0 0
> exec
| CREATE DATABASE eva_alter_status_slave DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_alter_status_slave
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x = 1
< ok 0 0
> exec
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE ON SLAVE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_alter_status_slave.e
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
| CREATE DEFINER=`root`@`%` EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE ON SLAVE DO SET @x = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci

== TestOracle_EventApplyCorrectness/8.0/create-at default#1
> exec
| DROP DATABASE IF EXISTS evto_create_at
< ok 0 0
> exec
| CREATE DATABASE evto_create_at DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_create_at
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE DO SET @y = 2
< ok 0 0
> query
| SHOW CREATE EVENT evto_create_at.e
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
| DROP DATABASE IF EXISTS eva_create_at
< ok 0 0
> exec
| CREATE DATABASE eva_create_at DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_create_at
< ok 0 0
> exec
| CREATE EVENT `e` ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE COMMENT '' DO SET @y = 2
< ok 0 0
> query
| SHOW CREATE EVENT eva_create_at.e
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

== TestOracle_EventApplyCorrectness/8.0/create-every default#1
> exec
| DROP DATABASE IF EXISTS evto_create_every
< ok 0 0
> exec
| CREATE DATABASE evto_create_every DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_create_every
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_create_every.e
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
| DROP DATABASE IF EXISTS eva_create_every
< ok 0 0
> exec
| CREATE DATABASE eva_create_every DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_create_every
< ok 0 0
> exec
| CREATE EVENT `e` ON SCHEDULE EVERY 1 HOUR STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_create_every.e
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

== TestOracle_EventApplyCorrectness/8.0/create-full default#1
> exec
| DROP DATABASE IF EXISTS evto_create_full
< ok 0 0
> exec
| CREATE DATABASE evto_create_full DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evto_create_full
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2025-01-01 00:00:00' ENDS '2030-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE COMMENT 'c' DO SET @q = 1
< ok 0 0
> query
| SHOW CREATE EVENT evto_create_full.e
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
| DROP DATABASE IF EXISTS eva_create_full
< ok 0 0
> exec
| CREATE DATABASE eva_create_full DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_create_full
< ok 0 0
> exec
| CREATE EVENT `e` ON SCHEDULE EVERY 1 DAY STARTS '2025-01-01 00:00:00' ENDS '2030-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE COMMENT 'c' DO SET @q = 1
< ok 0 0
> query
| SHOW CREATE EVENT eva_create_full.e
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

== TestOracle_EventApplyCorrectness/8.0/drop default#1
> exec
| DROP DATABASE IF EXISTS evfrom_drop
< ok 0 0
> exec
| CREATE DATABASE evfrom_drop DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE evfrom_drop
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> query
| SHOW CREATE EVENT evfrom_drop.e
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
| DROP DATABASE IF EXISTS eva_drop
< ok 0 0
> exec
| CREATE DATABASE eva_drop DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE eva_drop
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1
< ok 0 0
> exec
| DROP EVENT IF EXISTS `e`
< ok 0 0
> query
| SHOW CREATE EVENT eva_drop.e
< error
| 1539 HY000 Unknown event 'e'
