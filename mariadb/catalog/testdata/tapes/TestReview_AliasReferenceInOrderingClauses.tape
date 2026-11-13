-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_AliasReferenceInOrderingClauses default#1
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
| CREATE VIEW v_al0 (x) AS SELECT a AS z FROM t ORDER BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al0` AS select `t`.`a` AS `x` from `t` order by `t`.`a` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al0
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al0
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al0` AS select `t`.`a` AS `x` from `t` order by `t`.`a` + 0
< ok 0 0
> query
| SELECT * FROM v_al0
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al1 AS SELECT a + 1 AS z FROM t ORDER BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al1` AS select `t`.`a` + 1 AS `z` from `t` order by `z` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al1
< result
< column "z" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| DROP VIEW v_al1
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al1` AS select `t`.`a` + 1 AS `z` from `t` order by `z` + 0
< ok 0 0
> query
| SELECT * FROM v_al1
< result
< column "z" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| CREATE VIEW v_al2 AS SELECT 2 AS z FROM t ORDER BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al2` AS select 2 AS `z` from `t` order by `z` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al2
< result
< column "z" "INT"
< row
: int64 2
< row
: int64 2
> exec
| DROP VIEW v_al2
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al2` AS select 2 AS `z` from `t` order by `z` + 0
< ok 0 0
> query
| SELECT * FROM v_al2
< result
< column "z" "INT"
< row
: int64 2
< row
: int64 2
> exec
| CREATE VIEW v_al3 AS SELECT ABS(a) AS z FROM t ORDER BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al3` AS select abs(`t`.`a`) AS `z` from `t` order by `z` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al3
< result
< column "z" "BIGINT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al3
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al3` AS select abs(`t`.`a`) AS `z` from `t` order by `z` + 0
< ok 0 0
> query
| SELECT * FROM v_al3
< result
< column "z" "BIGINT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al4 AS SELECT t.a AS z FROM t ORDER BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al4` AS select `t`.`a` AS `z` from `t` order by `t`.`a` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al4
< result
< column "z" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al4
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al4` AS select `t`.`a` AS `z` from `t` order by `t`.`a` + 0
< ok 0 0
> query
| SELECT * FROM v_al4
< result
< column "z" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al5 AS SELECT a AS z FROM t ORDER BY -z
< ok 0 0
> query
| SHOW CREATE VIEW v_al5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al5` AS select `t`.`a` AS `z` from `t` order by -`t`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al5
< result
< column "z" "INT"
< row
: int64 2
< row
: int64 1
> exec
| DROP VIEW v_al5
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al5` AS select `t`.`a` AS `z` from `t` order by -`t`.`a`
< ok 0 0
> query
| SELECT * FROM v_al5
< result
< column "z" "INT"
< row
: int64 2
< row
: int64 1
> exec
| CREATE VIEW v_al6 AS SELECT a + 1 AS z, COUNT(*) AS n FROM t GROUP BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al6` AS select `t`.`a` + 1 AS `z`,count(0) AS `n` from `t` group by `t`.`a` + 1 + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al6
< result
< column "z" "BIGINT"
< column "n" "BIGINT"
< row
: int64 2
: int64 1
< row
: int64 3
: int64 1
> exec
| DROP VIEW v_al6
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al6` AS select `t`.`a` + 1 AS `z`,count(0) AS `n` from `t` group by `t`.`a` + 1 + 0
< ok 0 0
> query
| SELECT * FROM v_al6
< result
< column "z" "BIGINT"
< column "n" "BIGINT"
< row
: int64 2
: int64 1
< row
: int64 3
: int64 1
> exec
| CREATE VIEW v_al7 AS SELECT 2.5 AS z, COUNT(*) AS n FROM t GROUP BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al7` AS select 2.5 AS `z`,count(0) AS `n` from `t` group by 2.5 + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al7
< result
< column "z" "DECIMAL"
< column "n" "BIGINT"
< row
: bytes
| 2.5
: int64 2
> exec
| DROP VIEW v_al7
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al7` AS select 2.5 AS `z`,count(0) AS `n` from `t` group by 2.5 + 0
< ok 0 0
> query
| SELECT * FROM v_al7
< result
< column "z" "DECIMAL"
< column "n" "BIGINT"
< row
: bytes
| 2.5
: int64 2
> exec
| CREATE VIEW v_al8 AS SELECT a AS z, COUNT(*) AS n FROM t GROUP BY ABS(z)
< ok 0 0
> query
| SHOW CREATE VIEW v_al8
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al8
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al8` AS select `t`.`a` AS `z`,count(0) AS `n` from `t` group by abs(`t`.`a`)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al8
< result
< column "z" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_al8
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al8` AS select `t`.`a` AS `z`,count(0) AS `n` from `t` group by abs(`t`.`a`)
< ok 0 0
> query
| SELECT * FROM v_al8
< result
< column "z" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_al9 AS SELECT 2 AS z, COUNT(*) AS n FROM t GROUP BY ABS(z)
< ok 0 0
> query
| SHOW CREATE VIEW v_al9
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al9
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al9` AS select 2 AS `z`,count(0) AS `n` from `t` group by abs(2)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al9
< result
< column "z" "INT"
< column "n" "BIGINT"
< row
: int64 2
: int64 2
> exec
| DROP VIEW v_al9
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al9` AS select 2 AS `z`,count(0) AS `n` from `t` group by abs(2)
< ok 0 0
> query
| SELECT * FROM v_al9
< result
< column "z" "INT"
< column "n" "BIGINT"
< row
: int64 2
: int64 2
> exec
| CREATE VIEW v_al10 AS SELECT a AS z FROM t ORDER BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al10` AS select `t`.`a` AS `z` from `t` order by `t`.`a` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al10
< result
< column "z" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al10
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al10` AS select `t`.`a` AS `z` from `t` order by `t`.`a` + 0
< ok 0 0
> query
| SELECT * FROM v_al10
< result
< column "z" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al11 (z) AS SELECT a AS z FROM t ORDER BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al11
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al11
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al11` AS select `t`.`a` AS `z` from `t` order by `t`.`a` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al11
< result
< column "z" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al11
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al11` AS select `t`.`a` AS `z` from `t` order by `t`.`a` + 0
< ok 0 0
> query
| SELECT * FROM v_al11
< result
< column "z" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al12 (x) AS SELECT a AS z FROM t ORDER BY ABS(z)
< ok 0 0
> query
| SHOW CREATE VIEW v_al12
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al12
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al12` AS select `t`.`a` AS `x` from `t` order by abs(`t`.`a`)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al12
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al12
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al12` AS select `t`.`a` AS `x` from `t` order by abs(`t`.`a`)
< ok 0 0
> query
| SELECT * FROM v_al12
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al13 AS SELECT a AS z FROM t ORDER BY ABS(z)
< ok 0 0
> query
| SHOW CREATE VIEW v_al13
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al13
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al13` AS select `t`.`a` AS `z` from `t` order by abs(`t`.`a`)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al13
< result
< column "z" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al13
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al13` AS select `t`.`a` AS `z` from `t` order by abs(`t`.`a`)
< ok 0 0
> query
| SELECT * FROM v_al13
< result
< column "z" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al14 (x) AS SELECT a AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW v_al14
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al14
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al14` AS select `t`.`a` AS `x` from `t` order by `t`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al14
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al14
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al14` AS select `t`.`a` AS `x` from `t` order by `t`.`a`
< ok 0 0
> query
| SELECT * FROM v_al14
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al15 (x, y) AS SELECT a AS z, COUNT(*) AS n FROM t GROUP BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al15
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al15
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al15` AS select `t`.`a` AS `x`,count(0) AS `y` from `t` group by `t`.`a` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al15
< result
< column "x" "INT"
< column "y" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_al15
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al15` AS select `t`.`a` AS `x`,count(0) AS `y` from `t` group by `t`.`a` + 0
< ok 0 0
> query
| SELECT * FROM v_al15
< result
< column "x" "INT"
< column "y" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_al16 AS SELECT a AS z, COUNT(*) AS n FROM t GROUP BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al16
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al16
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al16` AS select `t`.`a` AS `z`,count(0) AS `n` from `t` group by `t`.`a` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al16
< result
< column "z" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_al16
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al16` AS select `t`.`a` AS `z`,count(0) AS `n` from `t` group by `t`.`a` + 0
< ok 0 0
> query
| SELECT * FROM v_al16
< result
< column "z" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_al17 (x) AS SELECT a AS z FROM t HAVING z > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al17
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al17
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al17` AS select `t`.`a` AS `x` from `t` having `x` > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al17
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al17
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al17` AS select `t`.`a` AS `x` from `t` having `x` > 0
< ok 0 0
> query
| SELECT * FROM v_al17
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al18 AS SELECT a AS z FROM t HAVING z > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al18
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al18
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al18` AS select `t`.`a` AS `z` from `t` having `z` > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al18
< result
< column "z" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al18
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al18` AS select `t`.`a` AS `z` from `t` having `z` > 0
< ok 0 0
> query
| SELECT * FROM v_al18
< result
< column "z" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al19 (x) AS SELECT 2 AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW v_al19
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al19
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al19` AS select 2 AS `x` from `t` order by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al19
< result
< column "x" "INT"
< row
: int64 2
< row
: int64 2
> exec
| DROP VIEW v_al19
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al19` AS select 2 AS `x` from `t` order by ''
< ok 0 0
> query
| SELECT * FROM v_al19
< result
< column "x" "INT"
< row
: int64 2
< row
: int64 2
> exec
| CREATE VIEW v_al20 (x) AS SELECT a AS z FROM t HAVING EXISTS(SELECT 1 FROM (SELECT 1 AS z HAVING z=1) AS d)
< ok 0 0
> query
| SHOW CREATE VIEW v_al20
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al20
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al20` AS select `t`.`a` AS `x` from `t` having exists(select 1 from (select 1 AS `z` having `z` = 1) `d` limit 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al20
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al20
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al20` AS select `t`.`a` AS `x` from `t` having exists(select 1 from (select 1 AS `z` having `z` = 1) `d` limit 1)
< ok 0 0
> query
| SELECT * FROM v_al20
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al21 (x) AS SELECT a AS z FROM t HAVING z > 0 AND EXISTS(SELECT 1 FROM (SELECT 1 AS z HAVING z=1) AS d)
< ok 0 0
> query
| SHOW CREATE VIEW v_al21
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al21
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al21` AS select `t`.`a` AS `x` from `t` having `x` > 0 and exists(select 1 from (select 1 AS `z` having `z` = 1) `d` limit 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al21
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al21
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al21` AS select `t`.`a` AS `x` from `t` having `x` > 0 and exists(select 1 from (select 1 AS `z` having `z` = 1) `d` limit 1)
< ok 0 0
> query
| SELECT * FROM v_al21
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al22 (x) AS SELECT a AS z FROM t HAVING (SELECT MAX(b) FROM t t2 HAVING MAX(b) > 0) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_al22
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al22
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al22` AS select `t`.`a` AS `x` from `t` having (select max(`t2`.`b`) from `t` `t2` having max(`t2`.`b`) > 0) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al22
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al22
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al22` AS select `t`.`a` AS `x` from `t` having (select max(`t2`.`b`) from `t` `t2` having max(`t2`.`b`) > 0) > 0
< ok 0 0
> query
| SELECT * FROM v_al22
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al23 AS WITH c(x) AS (SELECT a AS z FROM t HAVING z > 0) SELECT c.x FROM c
< ok 0 0
> query
| SHOW CREATE VIEW v_al23
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al23
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al23` AS with c(`x`) as (select `t`.`a` AS `x` from `t` having `x` > 0)select `c`.`x` AS `x` from `c`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al23
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al23
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al23` AS with c(`x`) as (select `t`.`a` AS `x` from `t` having `x` > 0)select `c`.`x` AS `x` from `c`
< ok 0 0
> query
| SELECT * FROM v_al23
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al24 AS WITH c(x) AS (SELECT a AS z FROM t ORDER BY z + 0) SELECT c.x FROM c
< ok 0 0
> query
| SHOW CREATE VIEW v_al24
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al24
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al24` AS with c(`x`) as (select `t`.`a` AS `x` from `t` order by `t`.`a` + 0)select `c`.`x` AS `x` from `c`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al24
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al24
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al24` AS with c(`x`) as (select `t`.`a` AS `x` from `t` order by `t`.`a` + 0)select `c`.`x` AS `x` from `c`
< ok 0 0
> query
| SELECT * FROM v_al24
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_al25 AS WITH c(x, y) AS (SELECT a AS z, COUNT(*) AS n FROM t GROUP BY z + 0) SELECT c.x FROM c
< ok 0 0
> query
| SHOW CREATE VIEW v_al25
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_al25
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al25` AS with c(`x`,`y`) as (select `t`.`a` AS `x`,count(0) AS `y` from `t` group by `t`.`a` + 0)select `c`.`x` AS `x` from `c`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_al25
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_al25
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_al25` AS with c(`x`,`y`) as (select `t`.`a` AS `x`,count(0) AS `y` from `t` group by `t`.`a` + 0)select `c`.`x` AS `x` from `c`
< ok 0 0
> query
| SELECT * FROM v_al25
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
