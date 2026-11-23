-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestViewTemporalPrecisionWrapOracle default#1
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(0) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `CURRENT_TIME(0)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(6) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime(6) AS `CURRENT_TIME(6)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(7) FROM a
< error
| 1426 42000 Too big precision specified for 'curtime'. Maximum is 6
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(65535) FROM a
< error
| 1426 42000 Too big precision specified for 'curtime'. Maximum is 6
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(65536) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `CURRENT_TIME(65536)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(2147483647) FROM a
< error
| 1426 42000 Too big precision specified for 'curtime'. Maximum is 6
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(2147483648) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `CURRENT_TIME(2147483648)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(2147483649) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime(1) AS `CURRENT_TIME(2147483649)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(2147483654) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime(6) AS `CURRENT_TIME(2147483654)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(2147483655) FROM a
< error
| 1426 42000 Too big precision specified for 'curtime'. Maximum is 6
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(9223372036854775808) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `CURRENT_TIME(9223372036854775808)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(18446744073709551615) FROM a
< error
| 1426 42000 Too big precision specified for 'curtime'. Maximum is 6
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(0x06) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime(6) AS `CURRENT_TIME(0x06)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(0x41) FROM a
< error
| 1426 42000 Too big precision specified for 'curtime'. Maximum is 6
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(0x10000) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `CURRENT_TIME(0x10000)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(0x10001) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime(1) AS `CURRENT_TIME(0x10001)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(0x8000000000000000) FROM a
< error
| 1426 42000 Too big precision specified for 'curtime'. Maximum is 6
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NOW(2147483648) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select current_timestamp() AS `NOW(2147483648)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NOW(2147483649) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select current_timestamp(1) AS `NOW(2147483649)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS _v
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SYSDATE(2147483648) FROM a
< ok 0 0
> query
| SHOW CREATE VIEW _v
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select sysdate() AS `SYSDATE(2147483648)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_wrap_ct
< ok 0 0
> exec
| CREATE TABLE t_wrap_ct (c DATETIME DEFAULT CURRENT_TIMESTAMP(2147483648))
< ok 0 0
> query
| SHOW CREATE TABLE t_wrap_ct
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_wrap_ct
: bytes
| CREATE TABLE `t_wrap_ct` (
|   `c` datetime DEFAULT current_timestamp()
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_bad_time
< ok 0 0
> exec
| CREATE TABLE t_bad_time (c TIME(65536))
< error
| 1426 42000 Too big precision specified for 'c'. Maximum is 6
