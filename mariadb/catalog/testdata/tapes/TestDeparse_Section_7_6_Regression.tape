-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestDeparse_Section_7_6_Regression default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_col_alias
< ok 0 0
> exec
| CREATE VIEW v_col_alias(x, y) AS SELECT a, b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_col_alias
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_col_alias
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_col_alias` AS select `t`.`a` AS `x`,`t`.`b` AS `y` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_reg_simple_select_column
< ok 0 0
> exec
| CREATE VIEW v_reg_simple_select_column AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_reg_simple_select_column
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_reg_simple_select_column
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_reg_simple_select_column` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_reg_simple_select_where
< ok 0 0
> exec
| CREATE VIEW v_reg_simple_select_where AS SELECT a FROM t WHERE a > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_reg_simple_select_where
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_reg_simple_select_where
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_reg_simple_select_where` AS select `t`.`a` AS `a` from `t` where `t`.`a` > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_reg_complex_union
< ok 0 0
> exec
| CREATE VIEW v_reg_complex_union AS SELECT a FROM t UNION SELECT b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_reg_complex_union
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_reg_complex_union
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_reg_complex_union` AS select `t`.`a` AS `a` from `t` union select `t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_reg_complex_boolean
< ok 0 0
> exec
| CREATE VIEW v_reg_complex_boolean AS SELECT a AND b, a OR b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_reg_complex_boolean
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_reg_complex_boolean
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_reg_complex_boolean` AS select `t`.`a` <> 0 and `t`.`b` <> 0 AS `a AND b`,`t`.`a` <> 0 or `t`.`b` <> 0 AS `a OR b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
