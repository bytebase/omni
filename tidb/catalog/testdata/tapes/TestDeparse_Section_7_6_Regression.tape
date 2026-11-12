-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_col_alias` (`x`, `y`) AS SELECT `a` AS `a`,`b` AS `b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_reg_simple_select_column` (`a`) AS SELECT `a` AS `a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_reg_simple_select_where` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` WHERE `a`>0
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_reg_complex_union` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` UNION SELECT `b` AS `b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_reg_complex_boolean` (`a AND b`, `a OR b`) AS SELECT `a` AND `b` AS `a AND b`,`a` OR `b` AS `a OR b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
