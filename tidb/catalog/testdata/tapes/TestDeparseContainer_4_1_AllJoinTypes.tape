-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparseContainer_4_1_AllJoinTypes default#1
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
| DROP VIEW IF EXISTS v_right_join
< ok 0 0
> exec
| CREATE VIEW v_right_join AS SELECT t1.a, t2.b FROM t1 RIGHT JOIN t2 ON t1.a = t2.a
< ok 0 0
> query
| SHOW CREATE VIEW v_right_join
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_right_join
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_right_join` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM `test`.`t1` RIGHT JOIN `test`.`t2` ON `t1`.`a`=`t2`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cross_join
< ok 0 0
> exec
| CREATE VIEW v_cross_join AS SELECT t1.a, t2.b FROM t1 CROSS JOIN t2
< ok 0 0
> query
| SHOW CREATE VIEW v_cross_join
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cross_join
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cross_join` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM `test`.`t1` JOIN `test`.`t2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_natural_join
< ok 0 0
> exec
| CREATE VIEW v_natural_join AS SELECT * FROM t1 NATURAL JOIN t2
< ok 0 0
> query
| SHOW CREATE VIEW v_natural_join
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_natural_join
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_natural_join` (`a`, `b`) AS SELECT `test`.`t1`.`a` AS `a`,`test`.`t1`.`b` AS `b` FROM `test`.`t1` NATURAL JOIN `test`.`t2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_straight_join
< ok 0 0
> exec
| CREATE VIEW v_straight_join AS SELECT t1.a, t2.b FROM t1 STRAIGHT_JOIN t2 ON t1.a = t2.a
< ok 0 0
> query
| SHOW CREATE VIEW v_straight_join
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_straight_join
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_straight_join` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM `test`.`t1` STRAIGHT_JOIN `test`.`t2` ON `t1`.`a`=`t2`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_using
< ok 0 0
> exec
| CREATE VIEW v_using AS SELECT t1.a, t2.b FROM t1 JOIN t2 USING (a)
< ok 0 0
> query
| SHOW CREATE VIEW v_using
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_using
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_using` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM `test`.`t1` JOIN `test`.`t2` USING (`a`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_comma_join
< ok 0 0
> exec
| CREATE VIEW v_comma_join AS SELECT t1.a, t2.b FROM t1, t2 WHERE t1.a = t2.a
< ok 0 0
> query
| SHOW CREATE VIEW v_comma_join
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_comma_join
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_comma_join` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM (`test`.`t1`) JOIN `test`.`t2` WHERE `t1`.`a`=`t2`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
