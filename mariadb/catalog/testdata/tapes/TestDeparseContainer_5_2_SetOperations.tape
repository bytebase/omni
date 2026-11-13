-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_union` AS select `t`.`a` AS `a` from `t` union select `t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_union_all` AS select `t`.`a` AS `a` from `t` union all select `t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_multi_union` AS select `t`.`a` AS `a` from `t` union select `t`.`b` AS `b` from `t` union select `t`.`c` AS `c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_intersect` AS select `t`.`a` AS `a` from `t` intersect select `t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_except` AS select `t`.`a` AS `a` from `t` except select `t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_union_ordlim` AS select `t`.`a` AS `a` from `t` union select `t`.`b` AS `b` from `t` order by 1 limit 5
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
