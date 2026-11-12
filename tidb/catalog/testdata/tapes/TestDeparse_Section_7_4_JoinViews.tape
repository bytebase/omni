-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparse_Section_7_4_JoinViews default#1
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
| DROP VIEW IF EXISTS v_inner_join
< ok 0 0
> exec
| CREATE VIEW v_inner_join AS SELECT t1.a, t2.b FROM t1 JOIN t2 ON t1.a = t2.a
< ok 0 0
> query
| SHOW CREATE VIEW v_inner_join
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_inner_join
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_inner_join` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM `test`.`t1` JOIN `test`.`t2` ON `t1`.`a`=`t2`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_left_join
< ok 0 0
> exec
| CREATE VIEW v_left_join AS SELECT t1.a, t2.b FROM t1 LEFT JOIN t2 ON t1.a = t2.a
< ok 0 0
> query
| SHOW CREATE VIEW v_left_join
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_left_join
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_left_join` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM `test`.`t1` LEFT JOIN `test`.`t2` ON `t1`.`a`=`t2`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_multi_table
< ok 0 0
> exec
| CREATE VIEW v_multi_table AS SELECT t1.a FROM t1, t2 WHERE t1.a = t2.a
< ok 0 0
> query
| SHOW CREATE VIEW v_multi_table
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_multi_table
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_multi_table` (`a`) AS SELECT `t1`.`a` AS `a` FROM (`test`.`t1`) JOIN `test`.`t2` WHERE `t1`.`a`=`t2`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_subquery_from
< ok 0 0
> exec
| CREATE VIEW v_subquery_from AS SELECT d.x FROM (SELECT a AS x FROM t) d
< ok 0 0
> query
| SHOW CREATE VIEW v_subquery_from
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_subquery_from
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_subquery_from` (`x`) AS SELECT `d`.`x` AS `x` FROM (SELECT `a` AS `x` FROM `test`.`t`) AS `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
