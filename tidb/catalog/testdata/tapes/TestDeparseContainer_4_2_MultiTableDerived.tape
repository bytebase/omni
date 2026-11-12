-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparseContainer_4_2_MultiTableDerived default#1
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
| CREATE TABLE IF NOT EXISTS t3 (b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_3join
< ok 0 0
> exec
| CREATE VIEW v_3join AS SELECT t1.a, t2.b, t3.c FROM t1 JOIN t2 ON t1.a = t2.a JOIN t3 ON t2.b = t3.b
< ok 0 0
> query
| SHOW CREATE VIEW v_3join
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_3join
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_3join` (`a`, `b`, `c`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b`,`t3`.`c` AS `c` FROM (`test`.`t1` JOIN `test`.`t2` ON `t1`.`a`=`t2`.`a`) JOIN `test`.`t3` ON `t2`.`b`=`t3`.`b`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_chain_left
< ok 0 0
> exec
| CREATE VIEW v_chain_left AS SELECT t1.a FROM t1 LEFT JOIN t2 ON t1.a = t2.a LEFT JOIN t3 ON t1.a = t3.b
< ok 0 0
> query
| SHOW CREATE VIEW v_chain_left
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_chain_left
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_chain_left` (`a`) AS SELECT `t1`.`a` AS `a` FROM (`test`.`t1` LEFT JOIN `test`.`t2` ON `t1`.`a`=`t2`.`a`) LEFT JOIN `test`.`t3` ON `t1`.`a`=`t3`.`b`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_derived
< ok 0 0
> exec
| CREATE VIEW v_derived AS SELECT d.x FROM (SELECT a AS x FROM t) d
< ok 0 0
> query
| SHOW CREATE VIEW v_derived
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_derived
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_derived` (`x`) AS SELECT `d`.`x` AS `x` FROM (SELECT `a` AS `x` FROM `test`.`t`) AS `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_derived_where
< ok 0 0
> exec
| CREATE VIEW v_derived_where AS SELECT d.x FROM (SELECT a AS x FROM t WHERE a > 0) AS d WHERE d.x < 10
< ok 0 0
> query
| SHOW CREATE VIEW v_derived_where
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_derived_where
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_derived_where` (`x`) AS SELECT `d`.`x` AS `x` FROM (SELECT `a` AS `x` FROM `test`.`t` WHERE `a`>0) AS `d` WHERE `d`.`x`<10
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_alias_as
< ok 0 0
> exec
| CREATE VIEW v_alias_as AS SELECT x.a FROM t AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_alias_as
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_alias_as
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_alias_as` (`a`) AS SELECT `x`.`a` AS `a` FROM `test`.`t` AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_alias_noas
< ok 0 0
> exec
| CREATE VIEW v_alias_noas AS SELECT x.a FROM t x
< ok 0 0
> query
| SHOW CREATE VIEW v_alias_noas
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_alias_noas
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_alias_noas` (`a`) AS SELECT `x`.`a` AS `a` FROM `test`.`t` AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
