-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_simple` AS select trim(`t`.`a`) AS `TRIM(a)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_leading` AS select trim(leading 'x' from `t`.`a`) AS `TRIM(LEADING 'x' FROM a)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_trailing` AS select trim(trailing 'x' from `t`.`a`) AS `TRIM(TRAILING 'x' FROM a)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_both` AS select trim(both 'x' from `t`.`a`) AS `TRIM(BOTH 'x' FROM a)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
