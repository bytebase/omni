-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ConstantOrderingItems default#1
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
| CREATE VIEW v_co0 AS SELECT a, COUNT(*) AS n FROM t GROUP BY TRUE
< ok 0 0
> query
| SHOW CREATE VIEW v_co0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co0` AS select `t`.`a` AS `a`,count(0) AS `n` from `t` group by 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co0
< result
< column "a" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_co0
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co0` AS select `t`.`a` AS `a`,count(0) AS `n` from `t` group by 1
< ok 0 0
> query
| SELECT * FROM v_co0
< result
< column "a" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_co1 AS SELECT a, COUNT(*) AS n FROM t GROUP BY 1
< ok 0 0
> query
| SHOW CREATE VIEW v_co1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co1` AS select `t`.`a` AS `a`,count(0) AS `n` from `t` group by 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co1
< result
< column "a" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_co1
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co1` AS select `t`.`a` AS `a`,count(0) AS `n` from `t` group by 1
< ok 0 0
> query
| SELECT * FROM v_co1
< result
< column "a" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_co2 AS SELECT a, b FROM t ORDER BY TRUE
< ok 0 0
> query
| SHOW CREATE VIEW v_co2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co2` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co2
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_co2
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co2` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 1
< ok 0 0
> query
| SELECT * FROM v_co2
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_co3 AS SELECT a, b FROM t ORDER BY (TRUE)
< ok 0 0
> query
| SHOW CREATE VIEW v_co3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co3` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co3
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_co3
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co3` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 1
< ok 0 0
> query
| SELECT * FROM v_co3
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_co4 AS SELECT a, b FROM t ORDER BY +2
< ok 0 0
> query
| SHOW CREATE VIEW v_co4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co4` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 2
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co4
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 2
: int64 1
< row
: int64 1
: int64 2
> exec
| DROP VIEW v_co4
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co4` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 2
< ok 0 0
> query
| SELECT * FROM v_co4
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 2
: int64 1
< row
: int64 1
: int64 2
> exec
| CREATE VIEW v_co5 AS SELECT a, b FROM t ORDER BY TRUE, a DESC
< ok 0 0
> query
| SHOW CREATE VIEW v_co5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co5` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 1,`t`.`a` desc
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co5
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_co5
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co5` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 1,`t`.`a` desc
< ok 0 0
> query
| SELECT * FROM v_co5
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_co6 AS SELECT a, b FROM t GROUP BY a, TRUE
< ok 0 0
> query
| SHOW CREATE VIEW v_co6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co6` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` group by `t`.`a`,1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co6
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_co6
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co6` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` group by `t`.`a`,1
< ok 0 0
> query
| SELECT * FROM v_co6
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_co7 AS SELECT a, b FROM t ORDER BY 2.5
< ok 0 0
> query
| SHOW CREATE VIEW v_co7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co7` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 2.5
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co7
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_co7
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co7` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 2.5
< ok 0 0
> query
| SELECT * FROM v_co7
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_co8 AS SELECT a, b FROM t ORDER BY NULL
< ok 0 0
> query
| SHOW CREATE VIEW v_co8
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co8
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co8` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by NULL
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co8
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_co8
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co8` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by NULL
< ok 0 0
> query
| SELECT * FROM v_co8
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_co9 AS SELECT a, b FROM t ORDER BY 1e2
< ok 0 0
> query
| SHOW CREATE VIEW v_co9
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co9
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co9` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 1e2
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co9
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_co9
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co9` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 1e2
< ok 0 0
> query
| SELECT * FROM v_co9
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_co10 AS SELECT a, b FROM t ORDER BY X'31'
< ok 0 0
> query
| SHOW CREATE VIEW v_co10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co10` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by X'31'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co10
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_co10
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co10` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by X'31'
< ok 0 0
> query
| SELECT * FROM v_co10
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_co11 AS SELECT a, b FROM t ORDER BY b'1'
< ok 0 0
> query
| SHOW CREATE VIEW v_co11
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co11
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co11` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 0x01
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co11
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_co11
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co11` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 0x01
< ok 0 0
> query
| SELECT * FROM v_co11
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_co12 AS SELECT a, b FROM t ORDER BY 0x31
< ok 0 0
> query
| SHOW CREATE VIEW v_co12
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co12
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co12` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 0x31
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co12
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_co12
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co12` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` order by 0x31
< ok 0 0
> query
| SELECT * FROM v_co12
< result
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_co13 AS SELECT SUM(a) OVER (ORDER BY TRUE) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co13
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co13
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co13` AS select sum(`t`.`a`) over ( order by '') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co13
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_co13
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co13` AS select sum(`t`.`a`) over ( order by '') AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co13
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_co14 AS SELECT SUM(a) OVER (ORDER BY 1) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co14
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co14
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co14` AS select sum(`t`.`a`) over ( order by '') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co14
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_co14
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co14` AS select sum(`t`.`a`) over ( order by '') AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co14
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_co15 AS SELECT SUM(a) OVER (ORDER BY -1) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co15
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co15
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co15` AS select sum(`t`.`a`) over ( order by '') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co15
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_co15
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co15` AS select sum(`t`.`a`) over ( order by '') AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co15
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_co16 AS SELECT SUM(a) OVER (ORDER BY 1 DESC) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co16
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co16
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co16` AS select sum(`t`.`a`) over ( order by '' desc) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co16
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_co16
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co16` AS select sum(`t`.`a`) over ( order by '' desc) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co16
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_co17 AS SELECT SUM(a) OVER (ORDER BY 1, a) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co17
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co17
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co17` AS select sum(`t`.`a`) over ( order by '',`t`.`a`) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co17
< result
< column "x" "DECIMAL"
< row
: bytes
| 1
< row
: bytes
| 3
> exec
| DROP VIEW v_co17
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co17` AS select sum(`t`.`a`) over ( order by '',`t`.`a`) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co17
< result
< column "x" "DECIMAL"
< row
: bytes
| 1
< row
: bytes
| 3
> exec
| CREATE VIEW v_co18 AS SELECT SUM(a) OVER (PARTITION BY TRUE) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co18
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co18
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co18` AS select sum(`t`.`a`) over ( partition by '') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co18
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_co18
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co18` AS select sum(`t`.`a`) over ( partition by '') AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co18
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_co19 AS SELECT SUM(a) OVER (PARTITION BY 1) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co19
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co19
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co19` AS select sum(`t`.`a`) over ( partition by '') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co19
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_co19
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co19` AS select sum(`t`.`a`) over ( partition by '') AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co19
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_co20 AS SELECT SUM(a) OVER (ORDER BY 2.5) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co20
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co20
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co20` AS select sum(`t`.`a`) over ( order by 2.5) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co20
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_co20
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co20` AS select sum(`t`.`a`) over ( order by 2.5) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co20
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_co21 AS SELECT SUM(a) OVER (ORDER BY NULL) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co21
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co21
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co21` AS select sum(`t`.`a`) over ( order by NULL) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co21
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_co21
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co21` AS select sum(`t`.`a`) over ( order by NULL) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co21
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_co22 AS SELECT SUM(a) OVER (ORDER BY 1e2) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co22
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co22
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co22` AS select sum(`t`.`a`) over ( order by 1e2) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co22
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_co22
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co22` AS select sum(`t`.`a`) over ( order by 1e2) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co22
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_co23 AS SELECT SUM(a) OVER (ORDER BY X'31') AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co23
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co23
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co23` AS select sum(`t`.`a`) over ( order by X'31') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co23
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_co23
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co23` AS select sum(`t`.`a`) over ( order by X'31') AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co23
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_co24 AS SELECT SUM(a) OVER w AS x FROM t WINDOW w AS (ORDER BY TRUE)
< ok 0 0
> query
| SHOW CREATE VIEW v_co24
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co24
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co24` AS select sum(`t`.`a`) over ( order by '') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co24
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_co24
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co24` AS select sum(`t`.`a`) over ( order by '') AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co24
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_co25 AS SELECT SUM(a) OVER w AS x FROM t WINDOW w AS (PARTITION BY 1)
< ok 0 0
> query
| SHOW CREATE VIEW v_co25
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co25
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co25` AS select sum(`t`.`a`) over ( partition by '') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co25
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| DROP VIEW v_co25
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co25` AS select sum(`t`.`a`) over ( partition by '') AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co25
< result
< column "x" "DECIMAL"
< row
: bytes
| 3
< row
: bytes
| 3
> exec
| CREATE VIEW v_co26 AS SELECT GROUP_CONCAT(a ORDER BY TRUE) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co26
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co26
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co26` AS select group_concat(`t`.`a` order by 1 ASC separator ',') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co26
< result
< column "x" "TEXT"
< row
: bytes
| 1,2
> exec
| DROP VIEW v_co26
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co26` AS select group_concat(`t`.`a` order by 1 ASC separator ',') AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co26
< result
< column "x" "TEXT"
< row
: bytes
| 1,2
> exec
| CREATE VIEW v_co27 AS SELECT GROUP_CONCAT(a ORDER BY 2.5) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_co27
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co27
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co27` AS select group_concat(`t`.`a` order by 2.5 ASC separator ',') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co27
< result
< column "x" "TEXT"
< row
: bytes
| 2,1
> exec
| DROP VIEW v_co27
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co27` AS select group_concat(`t`.`a` order by 2.5 ASC separator ',') AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_co27
< result
< column "x" "TEXT"
< row
: bytes
| 2,1
> exec
| CREATE VIEW v_co28 AS SELECT 2 AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW v_co28
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_co28
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co28` AS select 2 AS `z` from `t` order by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_co28
< result
< column "z" "INT"
< row
: int64 2
< row
: int64 2
> exec
| DROP VIEW v_co28
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_co28` AS select 2 AS `z` from `t` order by ''
< ok 0 0
> query
| SELECT * FROM v_co28
< result
< column "z" "INT"
< row
: int64 2
< row
: int64 2
