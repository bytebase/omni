-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparseContainer_5_2_SetOperations default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_union
< ok 0 0
> exec
| CREATE VIEW v_union AS SELECT a FROM t UNION SELECT b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_union
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_union
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_union` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` UNION SELECT `b` AS `b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_union_all
< ok 0 0
> exec
| CREATE VIEW v_union_all AS SELECT a FROM t UNION ALL SELECT b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_union_all
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_union_all
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_union_all` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` UNION ALL SELECT `b` AS `b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_multi_union
< ok 0 0
> exec
| CREATE VIEW v_multi_union AS SELECT a FROM t UNION SELECT b FROM t UNION SELECT c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_multi_union
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_multi_union
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_multi_union` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` UNION SELECT `b` AS `b` FROM `test`.`t` UNION SELECT `c` AS `c` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_intersect
< ok 0 0
> exec
| CREATE VIEW v_intersect AS SELECT a FROM t INTERSECT SELECT b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_intersect
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_intersect
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_intersect` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` INTERSECT SELECT `b` AS `b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_except
< ok 0 0
> exec
| CREATE VIEW v_except AS SELECT a FROM t EXCEPT SELECT b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_except
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_except
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_except` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` EXCEPT SELECT `b` AS `b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_union_ordlim
< ok 0 0
> exec
| CREATE VIEW v_union_ordlim AS SELECT a FROM t UNION SELECT b FROM t ORDER BY 1 LIMIT 5
< ok 0 0
> query
| SHOW CREATE VIEW v_union_ordlim
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_union_ordlim
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_union_ordlim` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` UNION SELECT `b` AS `b` FROM `test`.`t` ORDER BY 1 LIMIT 5
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
