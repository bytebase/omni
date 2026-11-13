-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestDeparseContainer_5_3_ColumnAliasPatterns default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t2 (a INT, b INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_alias_as_space
< ok 0 0
> exec
| CREATE VIEW v_alias_as_space AS SELECT a AS col1, b col2 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_alias_as_space
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_alias_as_space
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_alias_as_space` AS select `t`.`a` AS `col1`,`t`.`b` AS `col2` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_expr_alias
< ok 0 0
> exec
| CREATE VIEW v_expr_alias AS SELECT a + b AS sum_col FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_expr_alias
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_expr_alias
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_expr_alias` AS select (`t`.`a` + `t`.`b`) AS `sum_col` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_lit_auto
< ok 0 0
> exec
| CREATE VIEW v_lit_auto AS SELECT 1 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_lit_auto
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_lit_auto
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_lit_auto` AS select 1 AS `1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_star
< ok 0 0
> exec
| CREATE VIEW v_star AS SELECT * FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_star
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_star
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_star` AS select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_same_name_join
< ok 0 0
> exec
| CREATE VIEW v_same_name_join AS SELECT t1.a, t2.a FROM t1 JOIN t2 ON t1.a = t2.a
< error
| 1060 42S21 Duplicate column name 'a'
