-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparseContainer_Section_6_2_CTEPatterns default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_simple_cte
< ok 0 0
> exec
| CREATE VIEW v_simple_cte AS WITH cte AS (SELECT a FROM t) SELECT * FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v_simple_cte
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_simple_cte
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_simple_cte` (`a`) AS WITH `cte` AS (SELECT `a` AS `a` FROM `test`.`t`) SELECT `cte`.`a` AS `a` FROM `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cte_with_column_list
< ok 0 0
> exec
| CREATE VIEW v_cte_with_column_list AS WITH cte(x) AS (SELECT a FROM t) SELECT x FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v_cte_with_column_list
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cte_with_column_list
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_with_column_list` (`x`) AS WITH `cte` (`x`) AS (SELECT `a` AS `a` FROM `test`.`t`) SELECT `x` AS `x` FROM `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_recursive_cte
< ok 0 0
> exec
| CREATE VIEW v_recursive_cte AS WITH RECURSIVE cte AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM cte WHERE n < 10) SELECT * FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v_recursive_cte
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_recursive_cte
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_recursive_cte` (`n`) AS WITH RECURSIVE `cte` AS (SELECT 1 AS `n` UNION ALL SELECT `n`+1 AS `n + 1` FROM `cte` WHERE `n`<10) SELECT `cte`.`n` AS `n` FROM `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_multiple_ctes
< ok 0 0
> exec
| CREATE VIEW v_multiple_ctes AS WITH c1 AS (SELECT a FROM t), c2 AS (SELECT b FROM t) SELECT c1.a, c2.b FROM c1, c2
< ok 0 0
> query
| SHOW CREATE VIEW v_multiple_ctes
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_multiple_ctes
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_multiple_ctes` (`a`, `b`) AS WITH `c1` AS (SELECT `a` AS `a` FROM `test`.`t`), `c2` AS (SELECT `b` AS `b` FROM `test`.`t`) SELECT `c1`.`a` AS `a`,`c2`.`b` AS `b` FROM (`c1`) JOIN `c2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cte_used_in_union
< ok 0 0
> exec
| CREATE VIEW v_cte_used_in_union AS WITH cte AS (SELECT a FROM t) SELECT * FROM cte UNION SELECT * FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v_cte_used_in_union
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cte_used_in_union
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_used_in_union` (`a`) AS WITH `cte` AS (SELECT `a` AS `a` FROM `test`.`t`) SELECT `cte`.`a` AS `a` FROM `cte` UNION SELECT `cte`.`a` AS `a` FROM `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_recursive_cte_multi_col
< ok 0 0
> exec
| CREATE VIEW v_recursive_cte_multi_col AS WITH RECURSIVE cte AS (SELECT 1 AS n, 'a' AS s UNION ALL SELECT n + 1, CONCAT(s, 'a') FROM cte WHERE n < 5) SELECT * FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v_recursive_cte_multi_col
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_recursive_cte_multi_col
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_recursive_cte_multi_col` (`n`, `s`) AS WITH RECURSIVE `cte` AS (SELECT 1 AS `n`,_UTF8MB4'a' AS `s` UNION ALL SELECT `n`+1 AS `n + 1`,CONCAT(`s`, _UTF8MB4'a') AS `CONCAT(s, 'a')` FROM `cte` WHERE `n`<5) SELECT `cte`.`n` AS `n`,`cte`.`s` AS `s` FROM `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_recursive_cte_join_base
< ok 0 0
> exec
| CREATE VIEW v_recursive_cte_join_base AS WITH RECURSIVE cte AS (SELECT a, 1 AS depth FROM t WHERE a = 1 UNION ALL SELECT t.a, cte.depth + 1 FROM t JOIN cte ON t.b = cte.a WHERE cte.depth < 5) SELECT * FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v_recursive_cte_join_base
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_recursive_cte_join_base
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_recursive_cte_join_base` (`a`, `depth`) AS WITH RECURSIVE `cte` AS (SELECT `a` AS `a`,1 AS `depth` FROM `test`.`t` WHERE `a`=1 UNION ALL SELECT `t`.`a` AS `a`,`cte`.`depth`+1 AS `cte.depth + 1` FROM `test`.`t` JOIN `cte` ON `t`.`b`=`cte`.`a` WHERE `cte`.`depth`<5) SELECT `cte`.`a` AS `a`,`cte`.`depth` AS `depth` FROM `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cte_references_cte
< ok 0 0
> exec
| CREATE VIEW v_cte_references_cte AS WITH c1 AS (SELECT a FROM t), c2 AS (SELECT * FROM c1 WHERE a > 0) SELECT * FROM c2
< ok 0 0
> query
| SHOW CREATE VIEW v_cte_references_cte
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cte_references_cte
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_references_cte` (`a`) AS WITH `c1` AS (SELECT `a` AS `a` FROM `test`.`t`), `c2` AS (SELECT `c1`.`a` AS `a` FROM `c1` WHERE `a`>0) SELECT `c2`.`a` AS `a` FROM `c2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cte_with_aggregation
< ok 0 0
> exec
| CREATE VIEW v_cte_with_aggregation AS WITH cte AS (SELECT a, COUNT(*) AS cnt FROM t GROUP BY a) SELECT * FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v_cte_with_aggregation
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cte_with_aggregation
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_with_aggregation` (`a`, `cnt`) AS WITH `cte` AS (SELECT `a` AS `a`,COUNT(1) AS `cnt` FROM `test`.`t` GROUP BY `a`) SELECT `cte`.`a` AS `a`,`cte`.`cnt` AS `cnt` FROM `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cte_with_inner_union
< ok 0 0
> exec
| CREATE VIEW v_cte_with_inner_union AS WITH cte AS (SELECT a FROM t UNION SELECT b FROM t) SELECT * FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v_cte_with_inner_union
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cte_with_inner_union
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_with_inner_union` (`a`) AS WITH `cte` AS (SELECT `a` AS `a` FROM `test`.`t` UNION SELECT `b` AS `b` FROM `test`.`t`) SELECT `cte`.`a` AS `a` FROM `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cte_with_complex_expr
< ok 0 0
> exec
| CREATE VIEW v_cte_with_complex_expr AS WITH cte AS (SELECT a + b AS sum_val, CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END AS label FROM t) SELECT * FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v_cte_with_complex_expr
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cte_with_complex_expr
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_with_complex_expr` (`sum_val`, `label`) AS WITH `cte` AS (SELECT `a`+`b` AS `sum_val`,CASE WHEN `a`>0 THEN _UTF8MB4'pos' ELSE _UTF8MB4'neg' END AS `label` FROM `test`.`t`) SELECT `cte`.`sum_val` AS `sum_val`,`cte`.`label` AS `label` FROM `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
