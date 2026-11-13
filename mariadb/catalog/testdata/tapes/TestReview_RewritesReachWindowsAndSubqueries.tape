-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_RewritesReachWindowsAndSubqueries default#1
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
| CREATE VIEW v_rw0 AS SELECT SUM(a) OVER w AS x FROM t WINDOW w AS (PARTITION BY (a = 1) SOUNDS LIKE '0')
< ok 0 0
> query
| SHOW CREATE VIEW v_rw0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw0` AS select sum(`t`.`a`) over ( partition by soundex(`t`.`a` = 1) = soundex('0')) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw0 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_rw0
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw0` AS select sum(`t`.`a`) over ( partition by soundex(`t`.`a` = 1) = soundex('0')) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_rw0 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_rw1 AS SELECT SUM(a) OVER (PARTITION BY (a = 1) SOUNDS LIKE '0') AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_rw1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw1` AS select sum(`t`.`a`) over ( partition by soundex(`t`.`a` = 1) = soundex('0')) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw1 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_rw1
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw1` AS select sum(`t`.`a`) over ( partition by soundex(`t`.`a` = 1) = soundex('0')) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_rw1 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_rw2 AS SELECT SUM(a) OVER (ORDER BY (a = 1) SOUNDS LIKE '0') AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_rw2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw2` AS select sum(`t`.`a`) over ( order by soundex(`t`.`a` = 1) = soundex('0')) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw2 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_rw2
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw2` AS select sum(`t`.`a`) over ( order by soundex(`t`.`a` = 1) = soundex('0')) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_rw2 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_rw3 AS SELECT SUM(a) OVER (PARTITION BY NOT (a > 0)) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_rw3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw3` AS select sum(`t`.`a`) over ( partition by `t`.`a` <= 0) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw3 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_rw3
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw3` AS select sum(`t`.`a`) over ( partition by `t`.`a` <= 0) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_rw3 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_rw4 AS SELECT SUM(a) OVER (PARTITION BY a AND b) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_rw4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw4` AS select sum(`t`.`a`) over ( partition by `t`.`a` <> 0 and `t`.`b` <> 0) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw4 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_rw4
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw4` AS select sum(`t`.`a`) over ( partition by `t`.`a` <> 0 and `t`.`b` <> 0) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_rw4 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_rw5 AS SELECT SUM(a) OVER (ORDER BY NOT (a > 0)) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_rw5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw5` AS select sum(`t`.`a`) over ( order by `t`.`a` <= 0) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw5 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_rw5
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw5` AS select sum(`t`.`a`) over ( order by `t`.`a` <= 0) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_rw5 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_rw6 AS SELECT SUM(a) OVER (PARTITION BY NOT (a > 0 AND b > 0)) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_rw6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw6` AS select sum(`t`.`a`) over ( partition by `t`.`a` <= 0 or `t`.`b` <= 0) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw6 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_rw6
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw6` AS select sum(`t`.`a`) over ( partition by `t`.`a` <= 0 or `t`.`b` <= 0) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_rw6 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_rw7 AS SELECT d.x FROM (SELECT (a = 1) SOUNDS LIKE '0' AS x FROM t) d
< ok 0 0
> query
| SHOW CREATE VIEW v_rw7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw7` AS select `d`.`x` AS `x` from (select soundex(`t`.`a` = 1) = soundex('0') AS `x` from `t`) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw7 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 1
> exec
| DROP VIEW v_rw7
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw7` AS select `d`.`x` AS `x` from (select soundex(`t`.`a` = 1) = soundex('0') AS `x` from `t`) `d`
< ok 0 0
> query
| SELECT * FROM v_rw7 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 1
> exec
| CREATE VIEW v_rw8 AS SELECT d.x FROM (SELECT NOT (a > 0) AS x FROM t) d
< ok 0 0
> query
| SHOW CREATE VIEW v_rw8
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw8
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw8` AS select `d`.`x` AS `x` from (select `t`.`a` <= 0 AS `x` from `t`) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw8 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 0
< row
: int64 0
> exec
| DROP VIEW v_rw8
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw8` AS select `d`.`x` AS `x` from (select `t`.`a` <= 0 AS `x` from `t`) `d`
< ok 0 0
> query
| SELECT * FROM v_rw8 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 0
< row
: int64 0
> exec
| CREATE VIEW v_rw9 AS SELECT (SELECT (a = 1) SOUNDS LIKE '0' FROM t LIMIT 1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_rw9
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw9
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw9` AS select (select soundex(`t`.`a` = 1) = soundex('0') from `t` limit 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw9 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 1
> exec
| DROP VIEW v_rw9
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw9` AS select (select soundex(`t`.`a` = 1) = soundex('0') from `t` limit 1) AS `x`
< ok 0 0
> query
| SELECT * FROM v_rw9 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 1
> exec
| CREATE VIEW v_rw10 AS SELECT a FROM t WHERE a IN (SELECT (b = 1) SOUNDS LIKE '0' FROM t)
< ok 0 0
> query
| SHOW CREATE VIEW v_rw10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw10` AS select `t`.`a` AS `a` from `t` where `t`.`a` in (select soundex(`t`.`b` = 1) = soundex('0') from `t`)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw10 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 1
> exec
| DROP VIEW v_rw10
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw10` AS select `t`.`a` AS `a` from `t` where `t`.`a` in (select soundex(`t`.`b` = 1) = soundex('0') from `t`)
< ok 0 0
> query
| SELECT * FROM v_rw10 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 1
> exec
| CREATE VIEW v_rw11 AS SELECT a FROM t WHERE EXISTS (SELECT 1 FROM t t2 WHERE (t2.a = 1) SOUNDS LIKE '0')
< ok 0 0
> query
| SHOW CREATE VIEW v_rw11
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw11
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw11` AS select `t`.`a` AS `a` from `t` where exists(select 1 from `t` `t2` where soundex(`t2`.`a` = 1) = soundex('0') limit 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw11 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_rw11
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw11` AS select `t`.`a` AS `a` from `t` where exists(select 1 from `t` `t2` where soundex(`t2`.`a` = 1) = soundex('0') limit 1)
< ok 0 0
> query
| SELECT * FROM v_rw11 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_rw12 AS SELECT a FROM t WHERE a = ANY (SELECT (b = 1) SOUNDS LIKE '0' FROM t)
< ok 0 0
> query
| SHOW CREATE VIEW v_rw12
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw12
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw12` AS select `t`.`a` AS `a` from `t` where `t`.`a` in (select soundex(`t`.`b` = 1) = soundex('0') from `t`)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw12 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 1
> exec
| DROP VIEW v_rw12
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw12` AS select `t`.`a` AS `a` from `t` where `t`.`a` in (select soundex(`t`.`b` = 1) = soundex('0') from `t`)
< ok 0 0
> query
| SELECT * FROM v_rw12 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 1
> exec
| CREATE VIEW v_rw13 AS WITH cte AS (SELECT (a = 1) SOUNDS LIKE '0' AS x FROM t) SELECT cte.x FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v_rw13
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw13
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw13` AS with cte as (select soundex(`t`.`a` = 1) = soundex('0') AS `x` from `t`)select `cte`.`x` AS `x` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw13 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 1
> exec
| DROP VIEW v_rw13
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw13` AS with cte as (select soundex(`t`.`a` = 1) = soundex('0') AS `x` from `t`)select `cte`.`x` AS `x` from `cte`
< ok 0 0
> query
| SELECT * FROM v_rw13 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 1
> exec
| CREATE VIEW v_rw14 AS SELECT t1.a FROM t t1 JOIN t t2 ON (t1.a = 1) SOUNDS LIKE '0'
< ok 0 0
> query
| SHOW CREATE VIEW v_rw14
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw14
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw14` AS select `t1`.`a` AS `a` from (`t` `t1` join `t` `t2` on(soundex(`t1`.`a` = 1) = soundex('0')))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw14 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 1
< row
: int64 1
< row
: int64 2
< row
: int64 2
> exec
| DROP VIEW v_rw14
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw14` AS select `t1`.`a` AS `a` from (`t` `t1` join `t` `t2` on(soundex(`t1`.`a` = 1) = soundex('0')))
< ok 0 0
> query
| SELECT * FROM v_rw14 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 1
< row
: int64 1
< row
: int64 2
< row
: int64 2
> exec
| CREATE VIEW v_rw15 AS SELECT t1.a FROM t t1 JOIN t t2 ON NOT (t1.a > 0)
< ok 0 0
> query
| SHOW CREATE VIEW v_rw15
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw15
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw15` AS select `t1`.`a` AS `a` from (`t` `t1` join `t` `t2` on(`t1`.`a` <= 0))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw15 ORDER BY 1
< result
< column "a" "INT"
> exec
| DROP VIEW v_rw15
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw15` AS select `t1`.`a` AS `a` from (`t` `t1` join `t` `t2` on(`t1`.`a` <= 0))
< ok 0 0
> query
| SELECT * FROM v_rw15 ORDER BY 1
< result
< column "a" "INT"
> exec
| CREATE VIEW v_rw16 AS SELECT d.x FROM t t1 JOIN (SELECT NOT (a > 0) AS x FROM t) d ON d.x = t1.a
< ok 0 0
> query
| SHOW CREATE VIEW v_rw16
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw16
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw16` AS select `d`.`x` AS `x` from (`t` `t1` join (select `t`.`a` <= 0 AS `x` from `t`) `d` on(`d`.`x` = `t1`.`a`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw16 ORDER BY 1
< result
< column "x" "INT"
> exec
| DROP VIEW v_rw16
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw16` AS select `d`.`x` AS `x` from (`t` `t1` join (select `t`.`a` <= 0 AS `x` from `t`) `d` on(`d`.`x` = `t1`.`a`))
< ok 0 0
> query
| SELECT * FROM v_rw16 ORDER BY 1
< result
< column "x" "INT"
> exec
| CREATE VIEW v_rw17 AS SELECT d.x FROM (SELECT SUM(a) OVER (PARTITION BY NOT (a > 0)) AS x FROM t) d
< ok 0 0
> query
| SHOW CREATE VIEW v_rw17
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rw17
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw17` AS select `d`.`x` AS `x` from (select sum(`t`.`a`) over ( partition by `t`.`a` <= 0) AS `x` from `t`) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_rw17 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_rw17
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rw17` AS select `d`.`x` AS `x` from (select sum(`t`.`a`) over ( partition by `t`.`a` <= 0) AS `x` from `t`) `d`
< ok 0 0
> query
| SELECT * FROM v_rw17 ORDER BY 1
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
