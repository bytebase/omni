-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestDeparseContainer_8_1_ViewOfViewComplexStructures default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v2_vov
< ok 0 0
> exec
| DROP VIEW IF EXISTS v1_vov
< ok 0 0
> exec
| CREATE VIEW v1_vov AS SELECT a FROM t
< ok 0 0
> exec
| CREATE VIEW v2_vov AS SELECT * FROM v1_vov
< ok 0 0
> query
| SHOW CREATE VIEW v2_vov
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2_vov
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v2_vov` AS select `v1_vov`.`a` AS `a` from `v1_vov`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_ten_plus_columns
< ok 0 0
> exec
| CREATE VIEW v_ten_plus_columns AS SELECT a, b, c, a + 1, b + 1, c + 1, a * b, b * c, a * c, a + b + c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_ten_plus_columns
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ten_plus_columns
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ten_plus_columns` AS select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c`,`t`.`a` + 1 AS `a + 1`,`t`.`b` + 1 AS `b + 1`,`t`.`c` + 1 AS `c + 1`,`t`.`a` * `t`.`b` AS `a * b`,`t`.`b` * `t`.`c` AS `b * c`,`t`.`a` * `t`.`c` AS `a * c`,`t`.`a` + `t`.`b` + `t`.`c` AS `a + b + c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_reserved_word_aliases
< ok 0 0
> exec
| CREATE VIEW v_reserved_word_aliases AS SELECT a AS `select`, b AS `from`, c AS `where` FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_reserved_word_aliases
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_reserved_word_aliases
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_reserved_word_aliases` AS select `t`.`a` AS `select`,`t`.`b` AS `from`,`t`.`c` AS `where` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_case_without_else
< ok 0 0
> exec
| CREATE VIEW v_case_without_else AS SELECT CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' END FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_case_without_else
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_case_without_else
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_case_without_else` AS select case when `t`.`a` > 0 then 'pos' when `t`.`a` < 0 then 'neg' end AS `CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' END` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_between_column_bounds
< ok 0 0
> exec
| CREATE VIEW v_between_column_bounds AS SELECT a BETWEEN b AND c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_between_column_bounds
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_between_column_bounds
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_between_column_bounds` AS select `t`.`a` between `t`.`b` and `t`.`c` AS `a BETWEEN b AND c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
