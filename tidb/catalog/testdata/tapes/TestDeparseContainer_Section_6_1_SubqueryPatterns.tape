-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_scalar_subquery_in_select` (`(SELECT MAX(a) FROM t)`) AS SELECT (SELECT MAX(`a`) AS `MAX(a)` FROM `test`.`t`) AS `(SELECT MAX(a) FROM t)` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_in_subquery` (`a`, `b`, `c`) AS SELECT `test`.`t`.`a` AS `a`,`test`.`t`.`b` AS `b`,`test`.`t`.`c` AS `c` FROM `test`.`t` WHERE `a` IN (SELECT `a` AS `a` FROM `test`.`t` WHERE `a`>0)
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_exists_subquery` (`a`, `b`, `c`) AS SELECT `test`.`t`.`a` AS `a`,`test`.`t`.`b` AS `b`,`test`.`t`.`c` AS `c` FROM `test`.`t` WHERE EXISTS (SELECT 1 AS `1` FROM `test`.`t` WHERE `a`>0)
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_correlated_subquery` (`a`, `(SELECT COUNT(*) FROM t t2 WHERE t2.a = t1.a)`) AS SELECT `a` AS `a`,(SELECT COUNT(1) AS `COUNT(*)` FROM `test`.`t` AS `t2` WHERE `t2`.`a`=`t1`.`a`) AS `(SELECT COUNT(*) FROM t t2 WHERE t2.a = t1.a)` FROM `test`.`t` AS `t1`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nested_subqueries_2_levels` (`a`, `b`, `c`) AS SELECT `test`.`t`.`a` AS `a`,`test`.`t`.`b` AS `b`,`test`.`t`.`c` AS `c` FROM `test`.`t` WHERE `a` IN (SELECT `a` AS `a` FROM `test`.`t` WHERE `b` IN (SELECT `b` AS `b` FROM `test`.`t` WHERE `c`>0))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
