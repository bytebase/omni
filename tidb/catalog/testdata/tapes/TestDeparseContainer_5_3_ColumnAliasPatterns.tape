-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_alias_as_space` (`col1`, `col2`) AS SELECT `a` AS `col1`,`b` AS `col2` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_expr_alias` (`sum_col`) AS SELECT `a`+`b` AS `sum_col` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_lit_auto` (`1`) AS SELECT 1 AS `1` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_star` (`a`, `b`, `c`) AS SELECT `test`.`t`.`a` AS `a`,`test`.`t`.`b` AS `b`,`test`.`t`.`c` AS `c` FROM `test`.`t`
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
