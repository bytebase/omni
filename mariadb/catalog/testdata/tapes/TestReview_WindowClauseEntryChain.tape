-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_WindowClauseEntryChain default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW wc1 AS SELECT SUM(b) OVER w2 AS x FROM t WINDOW w1 AS (PARTITION BY a), w2 AS (w1 ORDER BY b)
< ok 0 0
> query
| SHOW CREATE VIEW wc1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| wc1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `wc1` AS select sum(`t`.`b`) over ( partition by `t`.`a` order by `t`.`b`) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW wc2 AS SELECT SUM(b) OVER w3 AS x FROM t WINDOW w1 AS (PARTITION BY a), w2 AS (w1), w3 AS (w2 ORDER BY b)
< ok 0 0
> query
| SHOW CREATE VIEW wc2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| wc2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `wc2` AS select sum(`t`.`b`) over ( partition by `t`.`a` order by `t`.`b`) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW wc3 AS SELECT SUM(b) OVER (w2) AS x FROM t WINDOW w1 AS (PARTITION BY a), w2 AS (w1 ORDER BY b)
< ok 0 0
> query
| SHOW CREATE VIEW wc3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| wc3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `wc3` AS select sum(`t`.`b`) over ( partition by `t`.`a` order by `t`.`b`) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW wc4 AS SELECT SUM(b) OVER (w2 ROWS UNBOUNDED PRECEDING) AS x FROM t WINDOW w1 AS (PARTITION BY a), w2 AS (w1 ORDER BY b)
< ok 0 0
> query
| SHOW CREATE VIEW wc4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| wc4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `wc4` AS select sum(`t`.`b`) over ( partition by `t`.`a` order by `t`.`b` rows between  unbounded  preceding  and  current row ) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
