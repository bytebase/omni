-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparseContainer_Section_5_1_SelectClauses default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_all_clauses
< ok 0 0
> exec
| CREATE VIEW v_all_clauses AS SELECT a FROM t WHERE a > 1 GROUP BY a HAVING COUNT(*) > 1 ORDER BY a LIMIT 10
< ok 0 0
> query
| SHOW CREATE VIEW v_all_clauses
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_all_clauses
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_all_clauses` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` WHERE `a`>1 GROUP BY `a` HAVING COUNT(1)>1 ORDER BY `a` LIMIT 10
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_alias_orderby
< ok 0 0
> exec
| CREATE VIEW v_alias_orderby AS SELECT a, COUNT(*) cnt FROM t GROUP BY a HAVING COUNT(*) > 1 ORDER BY cnt DESC
< ok 0 0
> query
| SHOW CREATE VIEW v_alias_orderby
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_alias_orderby
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_alias_orderby` (`a`, `cnt`) AS SELECT `a` AS `a`,COUNT(1) AS `cnt` FROM `test`.`t` GROUP BY `a` HAVING COUNT(1)>1 ORDER BY `cnt` DESC
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_distinct_orderby
< ok 0 0
> exec
| CREATE VIEW v_distinct_orderby AS SELECT DISTINCT a FROM t ORDER BY a DESC
< ok 0 0
> query
| SHOW CREATE VIEW v_distinct_orderby
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_distinct_orderby
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_distinct_orderby` (`a`) AS SELECT DISTINCT `a` AS `a` FROM `test`.`t` ORDER BY `a` DESC
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_multi_orderby
< ok 0 0
> exec
| CREATE VIEW v_multi_orderby AS SELECT a FROM t ORDER BY a, b DESC
< ok 0 0
> query
| SHOW CREATE VIEW v_multi_orderby
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_multi_orderby
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_multi_orderby` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` ORDER BY `a`,`b` DESC
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_limit_offset
< ok 0 0
> exec
| CREATE VIEW v_limit_offset AS SELECT a FROM t LIMIT 10 OFFSET 5
< ok 0 0
> query
| SHOW CREATE VIEW v_limit_offset
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_limit_offset
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_limit_offset` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` LIMIT 5,10
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_expr_groupby
< ok 0 0
> exec
| CREATE VIEW v_expr_groupby AS SELECT a + b FROM t GROUP BY a + b
< ok 0 0
> query
| SHOW CREATE VIEW v_expr_groupby
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_expr_groupby
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_expr_groupby` (`a + b`) AS SELECT `a`+`b` AS `a + b` FROM `test`.`t` GROUP BY `a`+`b`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
