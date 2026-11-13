-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_EventApplyCorrectness default#1
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_create_every
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| CREATE EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_create_every
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_create_at
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| ONE TIME
: bytes
| 2030-01-01 00:00:00
: bytes
|
: bytes
|
: bytes
|
: bytes
| ENABLED
: bytes
| PRESERVE
: bytes
|
: bytes
| SET @y = 2;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_create_at
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| ONE TIME
: bytes
| 2030-01-01 00:00:00
: bytes
|
: bytes
|
: bytes
|
: bytes
| ENABLED
: bytes
| PRESERVE
: bytes
|
: bytes
| SET @y = 2;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_create_full
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| DAY
: bytes
| 2030-01-01 00:00:00
: bytes
| DISABLED
: bytes
| PRESERVE
: bytes
| c
: bytes
| SET @q = 1;
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
| CREATE EVENT `e` ON SCHEDULE EVERY 1 DAY ENDS '2030-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE COMMENT 'c' DO SET @q = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_create_full
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| DAY
: bytes
| 2030-01-01 00:00:00
: bytes
| DISABLED
: bytes
| PRESERVE
: bytes
| c
: bytes
| SET @q = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_drop
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_drop
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_schedule_interval
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 2
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_schedule_interval
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| ALTER EVENT `e` ON SCHEDULE EVERY 2 HOUR ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_schedule_interval
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 2
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_schedule_unit
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| DAY
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_schedule_unit
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| ALTER EVENT `e` ON SCHEDULE EVERY 1 DAY ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_schedule_unit
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| DAY
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_every_to_at
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| ONE TIME
: bytes
| 2030-01-01 00:00:00
: bytes
|
: bytes
|
: bytes
|
: bytes
| ENABLED
: bytes
| PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_every_to_at
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_every_to_at
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| ONE TIME
: bytes
| 2030-01-01 00:00:00
: bytes
|
: bytes
|
: bytes
|
: bytes
| ENABLED
: bytes
| PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_status_disable
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| DISABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_status_disable
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE DISABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_status_disable
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| DISABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
< error
| 1064 42000 …near 'ON' at line 1
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_completion
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_completion
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_completion
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_comment
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
| b
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_comment
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
| a
: bytes
| SET @x = 1;
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
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE ENABLE COMMENT 'b' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_comment
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
| b
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_add_comment
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
| new
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_add_comment
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE ENABLE COMMENT 'new' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_add_comment
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
| new
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_body
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 2;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_body
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 2
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_body
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 2;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_add_ends
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
| 2030-01-01 00:00:00
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_add_ends
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR ENDS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_add_ends
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
| 2030-01-01 00:00:00
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_remove_comment
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_remove_comment
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
| old
: bytes
| SET @x = 1;
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
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_remove_comment
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
| old
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_remove_ends
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_remove_ends
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
| 2030-01-01 00:00:00
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| CREATE EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_remove_ends
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_status_reenable
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_status_reenable
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| DISABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_status_reenable
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_completion_reset
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_completion_reset
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_completion_reset
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evto_alter_at_to_every
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| evfrom_alter_at_to_every
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| ONE TIME
: bytes
| 2030-01-01 00:00:00
: bytes
|
: bytes
|
: bytes
|
: bytes
| ENABLED
: bytes
| PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x = 1
< ok 0 0
> query
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| eva_alter_at_to_every
: string
| e
< result
< column "EVENT_TYPE" "VARCHAR"
< column "IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), '')" "VARCHAR"
< column "STATUS" "VARCHAR"
< column "ON_COMPLETION" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "EVENT_DEFINITION" "TEXT"
< row
: bytes
| RECURRING
: bytes
|
: bytes
| 1
: bytes
| HOUR
: bytes
|
: bytes
| ENABLED
: bytes
| NOT PRESERVE
: bytes
|
: bytes
| SET @x = 1;
