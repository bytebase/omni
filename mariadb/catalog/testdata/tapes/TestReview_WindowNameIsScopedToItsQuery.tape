-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_WindowNameIsScopedToItsQuery default#1
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
| CREATE VIEW w_shadowed AS SELECT (SELECT SUM(b) OVER w FROM t WINDOW w AS (ORDER BY b) LIMIT 1) AS x FROM t WINDOW w AS (ORDER BY a)
< ok 0 0
> query
| SHOW CREATE VIEW w_shadowed
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| w_shadowed
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `w_shadowed` AS select (select sum(`t`.`b`) over ( order by `t`.`b`) from `t` limit 1) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW w_subquery AS SELECT (SELECT SUM(b) OVER w FROM t WINDOW w AS (ORDER BY b) LIMIT 1) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW w_subquery
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| w_subquery
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `w_subquery` AS select (select sum(`t`.`b`) over ( order by `t`.`b`) from `t` limit 1) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW w_derived AS SELECT x FROM (SELECT ROW_NUMBER() OVER w AS x FROM t WINDOW w AS (ORDER BY a)) d
< ok 0 0
> query
| SHOW CREATE VIEW w_derived
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| w_derived
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `w_derived` AS select `d`.`x` AS `x` from (select row_number() over ( order by `t`.`a`) AS `x` from `t`) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
