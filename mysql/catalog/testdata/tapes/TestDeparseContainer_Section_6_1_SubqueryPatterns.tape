-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestDeparseContainer_Section_6_1_SubqueryPatterns default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_scalar_subquery_in_select
< ok 0 0
> exec
| CREATE VIEW v_scalar_subquery_in_select AS SELECT (SELECT MAX(a) FROM t) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_scalar_subquery_in_select
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_scalar_subquery_in_select
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_scalar_subquery_in_select` AS select (select max(`t`.`a`) from `t`) AS `(SELECT MAX(a) FROM t)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_in_subquery
< ok 0 0
> exec
| CREATE VIEW v_in_subquery AS SELECT * FROM t WHERE a IN (SELECT a FROM t WHERE a > 0)
< ok 0 0
> query
| SHOW CREATE VIEW v_in_subquery
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_in_subquery
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_in_subquery` AS select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c` from `t` where `t`.`a` in (select `t`.`a` from `t` where (`t`.`a` > 0))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_exists_subquery
< ok 0 0
> exec
| CREATE VIEW v_exists_subquery AS SELECT * FROM t WHERE EXISTS (SELECT 1 FROM t WHERE a > 0)
< ok 0 0
> query
| SHOW CREATE VIEW v_exists_subquery
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_exists_subquery
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_exists_subquery` AS select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c` from `t` where exists(select 1 from `t` where (`t`.`a` > 0))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_correlated_subquery
< ok 0 0
> exec
| CREATE VIEW v_correlated_subquery AS SELECT a, (SELECT COUNT(*) FROM t t2 WHERE t2.a = t1.a) FROM t t1
< ok 0 0
> query
| SHOW CREATE VIEW v_correlated_subquery
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_correlated_subquery
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_correlated_subquery` AS select `t1`.`a` AS `a`,(select count(0) from `t` `t2` where (`t2`.`a` = `t1`.`a`)) AS `(SELECT COUNT(*) FROM t t2 WHERE t2.a = t1.a)` from `t` `t1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_nested_subqueries_2_levels
< ok 0 0
> exec
| CREATE VIEW v_nested_subqueries_2_levels AS SELECT * FROM t WHERE a IN (SELECT a FROM t WHERE b IN (SELECT b FROM t WHERE c > 0))
< ok 0 0
> query
| SHOW CREATE VIEW v_nested_subqueries_2_levels
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nested_subqueries_2_levels
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nested_subqueries_2_levels` AS select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c` from `t` where `t`.`a` in (select `t`.`a` from `t` where `t`.`b` in (select `t`.`b` from `t` where (`t`.`c` > 0)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
