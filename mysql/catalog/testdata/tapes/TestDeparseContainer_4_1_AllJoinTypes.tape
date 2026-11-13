-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_inner_join` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` join `t2` on((`t1`.`a` = `t2`.`a`)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_left_join` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` left join `t2` on((`t1`.`a` = `t2`.`a`)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_right_join` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t2` left join `t1` on((`t1`.`a` = `t2`.`a`)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cross_join` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` join `t2`)
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_natural_join` AS select `t1`.`a` AS `a`,`t1`.`b` AS `b` from (`t1` join `t2` on(((`t1`.`a` = `t2`.`a`) and (`t1`.`b` = `t2`.`b`))))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_straight_join` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` straight_join `t2` on((`t1`.`a` = `t2`.`a`)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_using` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` join `t2` on((`t1`.`a` = `t2`.`a`)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_comma_join` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` join `t2`) where (`t1`.`a` = `t2`.`a`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
