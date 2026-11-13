-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_EventIdempotence default#1
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| ev_idem_every_default
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| ev_idem_every_explicit_enable
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
| SECOND
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| ev_idem_every_disable
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
| 5
: bytes
| MINUTE
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
< error
| 1064 42000 …near 'ON' at line 1
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| ev_idem_every_preserve
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
| PRESERVE
: bytes
|
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| ev_idem_every_explicit_starts
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
| 5
: bytes
| MINUTE
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| ev_idem_every_starts_ends
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| ev_idem_every_only_ends
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
< error
| 1064 42000 …near ''1:30'' at line 1
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
< error
| 1064 42000 …near ''02:03:04'' at line 1
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
< error
| 1064 42000 …near ''1 12:30:00'' at line 1
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
< error
| 1064 42000 …near ''05'' at line 1
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| ev_idem_every_comment
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
| it's a test
: bytes
| SET @x = 1;
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| ev_idem_at_literal
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
| SELECT EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
| 			IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
| 			STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
| 		FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? AND EVENT_NAME = ?
: string
| ev_idem_at_preserve_disable
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
| 2030-06-01 00:00:00
: bytes
|
: bytes
|
: bytes
|
: bytes
| DISABLED
: bytes
| PRESERVE
: bytes
|
: bytes
| SET @y = 2;
