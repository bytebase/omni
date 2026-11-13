-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_simple_cte` AS with cte as (select `t`.`a` AS `a` from `t`)select `cte`.`a` AS `a` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_with_column_list` AS with cte(`x`) as (select `t`.`a` AS `x` from `t`)select `cte`.`x` AS `x` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_recursive_cte` AS with recursive cte as (select 1 AS `n` union all select `cte`.`n` + 1 AS `n + 1` from `cte` where `cte`.`n` < 10)select `cte`.`n` AS `n` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_multiple_ctes` AS with c1 as (select `t`.`a` AS `a` from `t`), c2 as (select `t`.`b` AS `b` from `t`)select `c1`.`a` AS `a`,`c2`.`b` AS `b` from (`c1` join `c2`)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_used_in_union` AS with cte as (select `t`.`a` AS `a` from `t`)select `cte`.`a` AS `a` from `cte` union select `cte`.`a` AS `a` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_recursive_cte_multi_col` AS with recursive cte as (select 1 AS `n`,'a' AS `s` union all select `cte`.`n` + 1 AS `n + 1`,concat(`cte`.`s`,'a') AS `CONCAT(s, 'a')` from `cte` where `cte`.`n` < 5)select `cte`.`n` AS `n`,`cte`.`s` AS `s` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_recursive_cte_join_base` AS with recursive cte as (select `t`.`a` AS `a`,1 AS `depth` from `t` where `t`.`a` = 1 union all select `t`.`a` AS `a`,`cte`.`depth` + 1 AS `cte.depth + 1` from (`t` join `cte` on(`t`.`b` = `cte`.`a`)) where `cte`.`depth` < 5)select `cte`.`a` AS `a`,`cte`.`depth` AS `depth` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_references_cte` AS with c1 as (select `t`.`a` AS `a` from `t`), c2 as (select `c1`.`a` AS `a` from `c1` where `c1`.`a` > 0)select `c2`.`a` AS `a` from `c2`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_with_aggregation` AS with cte as (select `t`.`a` AS `a`,count(0) AS `cnt` from `t` group by `t`.`a`)select `cte`.`a` AS `a`,`cte`.`cnt` AS `cnt` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_with_inner_union` AS with cte as (select `t`.`a` AS `a` from `t` union select `t`.`b` AS `b` from `t`)select `cte`.`a` AS `a` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cte_with_complex_expr` AS with cte as (select `t`.`a` + `t`.`b` AS `sum_val`,case when `t`.`a` > 0 then 'pos' else 'neg' end AS `label` from `t`)select `cte`.`sum_val` AS `sum_val`,`cte`.`label` AS `label` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
