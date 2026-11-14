-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_all_clauses` AS select `t`.`a` AS `a` from `t` where `t`.`a` > 1 group by `t`.`a` having count(0) > 1 order by `t`.`a` limit 10
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_alias_orderby` AS select `t`.`a` AS `a`,count(0) AS `cnt` from `t` group by `t`.`a` having count(0) > 1 order by count(0) desc
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_distinct_orderby` AS select distinct `t`.`a` AS `a` from `t` order by `t`.`a` desc
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_multi_orderby` AS select `t`.`a` AS `a` from `t` order by `t`.`a`,`t`.`b` desc
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_limit_offset` AS select `t`.`a` AS `a` from `t` limit 5,10
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_fetch_only
< ok 0 0
> exec
| CREATE VIEW v_fetch_only AS SELECT a FROM t ORDER BY a OFFSET 1 ROWS FETCH NEXT 3 ROWS ONLY
< ok 0 0
> query
| SHOW CREATE VIEW v_fetch_only
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_fetch_only
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_fetch_only` AS select `t`.`a` AS `a` from `t` order by `t`.`a` limit 1,3
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_fetch_row
< ok 0 0
> exec
| CREATE VIEW v_fetch_row AS SELECT a FROM t FETCH FIRST ROW ONLY
< ok 0 0
> query
| SHOW CREATE VIEW v_fetch_row
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_fetch_row
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_fetch_row` AS select `t`.`a` AS `a` from `t` limit 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_fetch_ties
< ok 0 0
> exec
| CREATE VIEW v_fetch_ties AS SELECT a FROM t ORDER BY a OFFSET 1 ROW FETCH FIRST 2 ROWS WITH TIES
< ok 0 0
> query
| SHOW CREATE VIEW v_fetch_ties
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_fetch_ties
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_fetch_ties` AS select `t`.`a` AS `a` from `t` order by `t`.`a` offset 1 rows  fetch first 2 rows with ties
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_offset_rows
< ok 0 0
> exec
| CREATE VIEW v_offset_rows AS SELECT a FROM t OFFSET 2 ROWS
< ok 0 0
> query
| SHOW CREATE VIEW v_offset_rows
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_offset_rows
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_offset_rows` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_select_unique
< ok 0 0
> exec
| CREATE VIEW v_select_unique AS SELECT UNIQUE a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_select_unique
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_select_unique
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_unique` AS select distinct `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_expr_groupby` AS select `t`.`a` + `t`.`b` AS `a + b` from `t` group by `t`.`a` + `t`.`b`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
