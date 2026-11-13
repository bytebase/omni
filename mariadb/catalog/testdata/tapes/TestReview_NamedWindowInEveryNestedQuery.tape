-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_NamedWindowInEveryNestedQuery default#1
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
| INSERT INTO t VALUES (1,2),(2,1)
< ok 2 0
> exec
| CREATE VIEW v_nw0 AS SELECT d.x FROM (SELECT SUM(a) OVER w AS x FROM t WINDOW w AS (ORDER BY a)) d
< ok 0 0
> query
| SHOW CREATE VIEW v_nw0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nw0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw0` AS select `d`.`x` AS `x` from (select sum(`t`.`a`) over ( order by `t`.`a`) AS `x` from `t`) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_nw0
< result
< column "x" "DECIMAL"
< row
: bytes
| 1
< row
: bytes
| 3
> exec
| DROP VIEW v_nw0
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw0` AS select `d`.`x` AS `x` from (select sum(`t`.`a`) over ( order by `t`.`a`) AS `x` from `t`) `d`
< ok 0 0
> query
| SELECT * FROM v_nw0
< result
< column "x" "DECIMAL"
< row
: bytes
| 1
< row
: bytes
| 3
> exec
| CREATE VIEW v_nw1 AS WITH c AS (SELECT SUM(a) OVER w AS x FROM t WINDOW w AS (ORDER BY a)) SELECT c.x FROM c
< ok 0 0
> query
| SHOW CREATE VIEW v_nw1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nw1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw1` AS with c as (select sum(`t`.`a`) over ( order by `t`.`a`) AS `x` from `t`)select `c`.`x` AS `x` from `c`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_nw1
< result
< column "x" "DECIMAL"
< row
: bytes
| 1
< row
: bytes
| 3
> exec
| DROP VIEW v_nw1
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw1` AS with c as (select sum(`t`.`a`) over ( order by `t`.`a`) AS `x` from `t`)select `c`.`x` AS `x` from `c`
< ok 0 0
> query
| SELECT * FROM v_nw1
< result
< column "x" "DECIMAL"
< row
: bytes
| 1
< row
: bytes
| 3
> exec
| CREATE VIEW v_nw2 AS SELECT (SELECT SUM(a) OVER w FROM t WINDOW w AS (ORDER BY a) LIMIT 1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_nw2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nw2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw2` AS select (select sum(`t`.`a`) over ( order by `t`.`a`) from `t` limit 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_nw2
< result
< column "x" "DECIMAL"
< row
: bytes
| 1
> exec
| DROP VIEW v_nw2
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw2` AS select (select sum(`t`.`a`) over ( order by `t`.`a`) from `t` limit 1) AS `x`
< ok 0 0
> query
| SELECT * FROM v_nw2
< result
< column "x" "DECIMAL"
< row
: bytes
| 1
> exec
| CREATE VIEW v_nw3 AS SELECT a FROM t WHERE EXISTS (SELECT SUM(a) OVER w FROM t t2 WINDOW w AS (ORDER BY a))
< ok 0 0
> query
| SHOW CREATE VIEW v_nw3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nw3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw3` AS select `t`.`a` AS `a` from `t` where exists(select sum(`t2`.`a`) over ( order by `t2`.`a`) from `t` `t2` limit 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_nw3
< result
< column "a" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_nw3
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw3` AS select `t`.`a` AS `a` from `t` where exists(select sum(`t2`.`a`) over ( order by `t2`.`a`) from `t` `t2` limit 1)
< ok 0 0
> query
| SELECT * FROM v_nw3
< result
< column "a" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_nw4 AS SELECT a FROM t WHERE a IN (SELECT SUM(a) OVER w FROM t t2 WINDOW w AS (ORDER BY a))
< ok 0 0
> query
| SHOW CREATE VIEW v_nw4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nw4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw4` AS select `t`.`a` AS `a` from `t` where `t`.`a` in (select sum(`t2`.`a`) over ( order by `t2`.`a`) from `t` `t2`)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_nw4
< result
< column "a" "INT"
> exec
| DROP VIEW v_nw4
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw4` AS select `t`.`a` AS `a` from `t` where `t`.`a` in (select sum(`t2`.`a`) over ( order by `t2`.`a`) from `t` `t2`)
< ok 0 0
> query
| SELECT * FROM v_nw4
< result
< column "a" "INT"
> exec
| CREATE VIEW v_nw5 AS SELECT 1 AS x UNION SELECT SUM(a) OVER w FROM t WINDOW w AS (ORDER BY a)
< ok 0 0
> query
| SHOW CREATE VIEW v_nw5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nw5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw5` AS select 1 AS `x` union select sum(`t`.`a`) over ( order by `t`.`a`) AS `SUM(a) OVER w` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_nw5
< result
< column "x" "DECIMAL"
< row
: bytes
| 1
< row
: bytes
| 3
> exec
| DROP VIEW v_nw5
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw5` AS select 1 AS `x` union select sum(`t`.`a`) over ( order by `t`.`a`) AS `SUM(a) OVER w` from `t`
< ok 0 0
> query
| SELECT * FROM v_nw5
< result
< column "x" "DECIMAL"
< row
: bytes
| 1
< row
: bytes
| 3
> exec
| CREATE VIEW v_nw6 AS SELECT t1.a FROM t t1 JOIN (SELECT SUM(a) OVER w AS x FROM t WINDOW w AS (ORDER BY a)) d ON d.x = t1.a
< ok 0 0
> query
| SHOW CREATE VIEW v_nw6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nw6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw6` AS select `t1`.`a` AS `a` from (`t` `t1` join (select sum(`t`.`a`) over ( order by `t`.`a`) AS `x` from `t`) `d` on(`d`.`x` = `t1`.`a`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_nw6
< result
< column "a" "INT"
< row
: int64 1
> exec
| DROP VIEW v_nw6
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw6` AS select `t1`.`a` AS `a` from (`t` `t1` join (select sum(`t`.`a`) over ( order by `t`.`a`) AS `x` from `t`) `d` on(`d`.`x` = `t1`.`a`))
< ok 0 0
> query
| SELECT * FROM v_nw6
< result
< column "a" "INT"
< row
: int64 1
> exec
| CREATE VIEW v_nw7 AS SELECT SUM(a) OVER (PARTITION BY (SELECT SUM(b) OVER w FROM t WINDOW w AS (ORDER BY b) LIMIT 1)) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_nw7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nw7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw7` AS select sum(`t`.`a`) over ( partition by (select sum(`t`.`b`) over ( order by `t`.`b`) from `t` limit 1)) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_nw7
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_nw7
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nw7` AS select sum(`t`.`a`) over ( partition by (select sum(`t`.`b`) over ( order by `t`.`b`) from `t` limit 1)) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_nw7
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
