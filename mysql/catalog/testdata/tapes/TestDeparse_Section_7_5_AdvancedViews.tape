-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_union_view` AS select `t`.`a` AS `a` from `t` union select `t`.`b` AS `b` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_view` AS with `cte` as (select `t`.`a` AS `a` from `t`) select `cte`.`a` AS `a` from `cte`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_window_func_view` AS select `t`.`a` AS `a`,row_number() OVER (ORDER BY `t`.`a` )  AS `ROW_NUMBER() OVER (ORDER BY a)` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nested_subquery_view` AS select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c` from `t` where `t`.`a` in (select `t`.`a` from `t` where (`t`.`a` > 0))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_boolean_expr_view` AS select ((0 <> `t`.`a`) and (0 <> `t`.`b`)) AS `a AND b`,((0 <> `t`.`a`) or (0 <> `t`.`b`)) AS `a OR b` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_combined_rewrite_view` AS select (`t`.`a` + `t`.`b`) AS `a + b`,(`t`.`a` <= 0) AS `NOT (a > 0)`,cast(`t`.`a` as char charset utf8mb4) AS `CAST(a AS CHAR)`,count(0) AS `COUNT(*)` from `t` group by `t`.`a`,`t`.`b` having (count(0) > 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
