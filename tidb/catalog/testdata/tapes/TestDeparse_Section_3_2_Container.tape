-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparse_Section_3_2_Container default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a VARCHAR(50), b VARCHAR(50))
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_trim_simple
< ok 0 0
> exec
| CREATE VIEW v_trim_simple AS SELECT TRIM(a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_simple
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_simple
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_simple` (`TRIM(a)`) AS SELECT TRIM(`a`) AS `TRIM(a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_trim_leading
< ok 0 0
> exec
| CREATE VIEW v_trim_leading AS SELECT TRIM(LEADING 'x' FROM a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_leading
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_leading
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_leading` (`TRIM(LEADING 'x' FROM a)`) AS SELECT TRIM(LEADING _UTF8MB4'x' FROM `a`) AS `TRIM(LEADING 'x' FROM a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_trim_trailing
< ok 0 0
> exec
| CREATE VIEW v_trim_trailing AS SELECT TRIM(TRAILING 'x' FROM a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_trailing
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_trailing
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_trailing` (`TRIM(TRAILING 'x' FROM a)`) AS SELECT TRIM(TRAILING _UTF8MB4'x' FROM `a`) AS `TRIM(TRAILING 'x' FROM a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_trim_both
< ok 0 0
> exec
| CREATE VIEW v_trim_both AS SELECT TRIM(BOTH 'x' FROM a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_both
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_both
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_both` (`TRIM(BOTH 'x' FROM a)`) AS SELECT TRIM(BOTH _UTF8MB4'x' FROM `a`) AS `TRIM(BOTH 'x' FROM a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
