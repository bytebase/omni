-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparse_Section_7_5_AdvancedViews default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_union_view
< ok 0 0
> exec
| CREATE VIEW v_union_view AS SELECT a FROM t UNION SELECT b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_union_view
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_union_view
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_union_view` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` UNION SELECT `b` AS `b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cte_view
< ok 0 0
> exec
| CREATE VIEW v_cte_view AS WITH cte AS (SELECT a FROM t) SELECT * FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v_cte_view
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cte_view
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_view` (`a`) AS WITH `cte` AS (SELECT `a` AS `a` FROM `test`.`t`) SELECT `cte`.`a` AS `a` FROM `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_window_func_view
< ok 0 0
> exec
| CREATE VIEW v_window_func_view AS SELECT a, ROW_NUMBER() OVER (ORDER BY a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_window_func_view
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_window_func_view
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_window_func_view` (`a`, `ROW_NUMBER() OVER (ORDER BY a)`) AS SELECT `a` AS `a`,ROW_NUMBER() OVER (ORDER BY `a`) AS `ROW_NUMBER() OVER (ORDER BY a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_nested_subquery_view
< ok 0 0
> exec
| CREATE VIEW v_nested_subquery_view AS SELECT * FROM t WHERE a IN (SELECT a FROM t WHERE a > 0)
< ok 0 0
> query
| SHOW CREATE VIEW v_nested_subquery_view
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nested_subquery_view
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nested_subquery_view` (`a`, `b`, `c`) AS SELECT `test`.`t`.`a` AS `a`,`test`.`t`.`b` AS `b`,`test`.`t`.`c` AS `c` FROM `test`.`t` WHERE `a` IN (SELECT `a` AS `a` FROM `test`.`t` WHERE `a`>0)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_boolean_expr_view
< ok 0 0
> exec
| CREATE VIEW v_boolean_expr_view AS SELECT a AND b, a OR b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_boolean_expr_view
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_boolean_expr_view
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_boolean_expr_view` (`a AND b`, `a OR b`) AS SELECT `a` AND `b` AS `a AND b`,`a` OR `b` AS `a OR b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_combined_rewrite_view
< ok 0 0
> exec
| CREATE VIEW v_combined_rewrite_view AS SELECT a + b, NOT (a > 0), CAST(a AS CHAR), COUNT(*) FROM t GROUP BY a, b HAVING COUNT(*) > 1
< ok 0 0
> query
| SHOW CREATE VIEW v_combined_rewrite_view
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_combined_rewrite_view
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_combined_rewrite_view` (`a + b`, `NOT (a > 0)`, `CAST(a AS CHAR)`, `COUNT(*)`) AS SELECT `a`+`b` AS `a + b`,NOT (`a`>0) AS `NOT (a > 0)`,CAST(`a` AS CHAR(11)) AS `CAST(a AS CHAR)`,COUNT(1) AS `COUNT(*)` FROM `test`.`t` GROUP BY `a`,`b` HAVING COUNT(1)>1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
