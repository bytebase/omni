-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestViewNamesAndJoinColumnsOracle default#1
> exec
| CREATE TABLE n1 (a INT, b INT)
< ok 0 0
> exec
| CREATE TABLE j1 (k1 INT, a INT)
< ok 0 0
> exec
| CREATE TABLE j2 (k1 INT, k2 INT, b INT)
< ok 0 0
> exec
| CREATE TABLE j3 (k2 INT, k3 INT, c INT)
< ok 0 0
> exec
| CREATE TABLE j4 (k3 INT, d INT)
< ok 0 0
> exec
| CREATE TABLE j5 (a INT, k1 INT, b INT, k2 INT)
< ok 0 0
> exec
| CREATE TABLE p (k INT, p1 INT)
< ok 0 0
> exec
| CREATE TABLE q (q1 INT, q2 INT)
< ok 0 0
> exec
| CREATE TABLE r (r1 INT, k INT)
< ok 0 0
> exec
| CREATE TABLE s (s1 INT, s2 INT)
< ok 0 0
> exec
| CREATE TABLE u (k INT, u1 INT)
< ok 0 0
> exec
| CREATE VIEW a01 AS SELECT (WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2) FROM n1
< ok 0 0
> query
| SHOW CREATE VIEW a01
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a01
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a01` AS select (with `c` as (select 1 AS `x`) select `c`.`x` from `c` union select 2) AS `(WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2)` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a01' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2)
> exec
| CREATE VIEW a02 AS select (with c as (select 1 as x) select x from   c union select 2) from n1
< ok 0 0
> query
| SHOW CREATE VIEW a02
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a02
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a02` AS select (with `c` as (select 1 AS `x`) select `c`.`x` from `c` union select 2) AS `(with c as (select 1 as x) select x from   c union select 2)` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a02' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (with c as (select 1 as x) select x from   c union select 2)
> exec
| CREATE VIEW a03 AS SELECT (SELECT 1 UNION SELECT 2) FROM n1
< ok 0 0
> query
| SHOW CREATE VIEW a03
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a03
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a03` AS select (select 1 union select 2) AS `(SELECT 1 UNION SELECT 2)` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a03' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (SELECT 1 UNION SELECT 2)
> exec
| CREATE VIEW a04 AS SELECT (SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1), (SELECT DISTINCT a FROM n1 LIMIT 1) FROM n1
< ok 0 0
> query
| SHOW CREATE VIEW a04
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a04
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a04` AS select (select `n1`.`a` from `n1` group by `n1`.`a` having (`n1`.`a` > 0) order by `n1`.`a` limit 1) AS `(SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1)`,(select distinct `n1`.`a` from `n1` limit 1) AS `(SELECT DISTINCT a FROM n1 LIMIT 1)` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a04' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1)
< row
: bytes
| (SELECT DISTINCT a FROM n1 LIMIT 1)
> exec
| CREATE VIEW a05 AS SELECT a+1, a  +  1, count(*), cast(a as char) FROM n1 GROUP BY a
< ok 0 0
> query
| SHOW CREATE VIEW a05
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a05
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a05` AS select (`n1`.`a` + 1) AS `a+1`,(`n1`.`a` + 1) AS `a  +  1`,count(0) AS `count(*)`,cast(`n1`.`a` as char charset utf8mb4) AS `cast(a as char)` from `n1` group by `n1`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a05' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a+1
< row
: bytes
| a  +  1
< row
: bytes
| count(*)
< row
: bytes
| cast(a as char)
> exec
| CREATE VIEW a06 AS SELECT (a), ((b)), (1), 0004, - 2, x'41', B'10', true, date '2020-01-01', ('e'), (null) FROM n1
< ok 0 0
> query
| SHOW CREATE VIEW a06
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a06
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a06` AS select `n1`.`a` AS `a`,`n1`.`b` AS `b`,1 AS `1`,4 AS `0004`,-(2) AS `- 2`,0x41 AS `x'41'`,0x02 AS `B'10'`,true AS `true`,DATE'2020-01-01' AS `date '2020-01-01'`,'e' AS `e`,NULL AS `NULL` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a06' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| 1
< row
: bytes
| 0004
< row
: bytes
| - 2
< row
: bytes
| x'41'
< row
: bytes
| B'10'
< row
: bytes
| true
< row
: bytes
| date '2020-01-01'
< row
: bytes
| e
< row
: bytes
| NULL
> exec
| CREATE VIEW a07 AS SELECT a /* c */ + 1, a/*!50100 +2*/ FROM n1
< ok 0 0
> query
| SHOW CREATE VIEW a07
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a07
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a07` AS select (`n1`.`a` + 1) AS `a /* c */ + 1`,(`n1`.`a` + 2) AS `a +2` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a07' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a /* c */ + 1
< row
: bytes
| a +2
> exec
| CREATE VIEW a08 AS SELECT concat('éééééééééééééééééééééééééééééééééééééééééééééééééé') FROM n1
< ok 0 0
> query
| SHOW CREATE VIEW a08
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a08
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a08` AS select concat('éééééééééééééééééééééééééééééééééééééééééééééééééé') AS `concat('éééééééééééééééééééééééééééééééééééééééééééééééééé')` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a08' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| concat('éééééééééééééééééééééééééééééééééééééééééééééééééé')
> exec
| CREATE VIEW a09 AS SELECT concat('eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee') FROM n1
< ok 0 0
> query
| SHOW CREATE VIEW a09
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a09
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a09` AS select concat('eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee') AS `Name_exp_1` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a09' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| Name_exp_1
> exec
| CREATE VIEW a10 AS SELECT * FROM (SELECT a+1, b FROM n1) d
< ok 0 0
> query
| SHOW CREATE VIEW a10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a10` AS select `d`.`a+1` AS `a+1`,`d`.`b` AS `b` from (select (`n1`.`a` + 1) AS `a+1`,`n1`.`b` AS `b` from `n1`) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a10' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a+1
< row
: bytes
| b
> exec
| CREATE VIEW a11 AS SELECT (select   1), (SELECT a FROM n1 LIMIT 1) + 0 FROM n1
< ok 0 0
> query
| SHOW CREATE VIEW a11
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a11
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a11` AS select (select 1) AS `(select   1)`,((select `n1`.`a` from `n1` limit 1) + 0) AS `(SELECT a FROM n1 LIMIT 1) + 0` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a11' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (select   1)
< row
: bytes
| (SELECT a FROM n1 LIMIT 1) + 0
> exec
| CREATE VIEW a12 AS SELECT a IN (SELECT a FROM n1), a in (1,2) FROM n1
< ok 0 0
> query
| SHOW CREATE VIEW a12
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a12
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a12` AS select `n1`.`a` in (select `n1`.`a` from `n1`) AS `a IN (SELECT a FROM n1)`,(`n1`.`a` in (1,2)) AS `a in (1,2)` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a12' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a IN (SELECT a FROM n1)
< row
: bytes
| a in (1,2)
> exec
| CREATE VIEW a13 AS SELECT a + 1 AS `x + y`, count(*) AS `count(z)` FROM n1
< ok 0 0
> query
| SHOW CREATE VIEW a13
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a13
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a13` AS select (`n1`.`a` + 1) AS `x + y`,count(0) AS `count(z)` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a13' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| x + y
< row
: bytes
| count(z)
> exec
| CREATE VIEW a14 AS SELECT EXISTS (SELECT 1 UNION SELECT 2), a IN (SELECT 1 UNION SELECT 2) FROM n1
< ok 0 0
> query
| SHOW CREATE VIEW a14
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a14
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a14` AS select exists(select 1 union select 2) AS `EXISTS (SELECT 1 UNION SELECT 2)`,`n1`.`a` in (select 1 union select 2) AS `a IN (SELECT 1 UNION SELECT 2)` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a14' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| EXISTS (SELECT 1 UNION SELECT 2)
< row
: bytes
| a IN (SELECT 1 UNION SELECT 2)
> exec
| CREATE VIEW a15 AS SELECT * FROM (SELECT `a`+1 FROM n1) d
< ok 0 0
> query
| SHOW CREATE VIEW a15
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a15
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a15` AS select `d`.```a``+1` AS ```a``+1` from (select (`n1`.`a` + 1) AS ```a``+1` from `n1`) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a15' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| `a`+1
> exec
| CREATE VIEW a16 AS SELECT * FROM (SELECT a AS `x``y` FROM n1) d
< ok 0 0
> query
| SHOW CREATE VIEW a16
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a16
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a16` AS select `d`.`x``y` AS `x``y` from (select `n1`.`a` AS `x``y` from `n1`) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a16' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| x`y
> exec
| CREATE VIEW a17 AS SELECT a /*! + 1 /*! + 2 */ */*c*/3 FROM n1
< ok 0 0
> query
| SHOW CREATE VIEW a17
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a17
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a17` AS select ((`n1`.`a` + 1) + (2 * 3)) AS `a  + 1  + 2  */*c*/3` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'a17' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a  + 1  + 2  */*c*/3
> exec
| CREATE VIEW b01 AS SELECT * FROM (j1 JOIN j2 USING (k1)) JOIN j3 USING (k2)
< ok 0 0
> query
| SHOW CREATE VIEW b01
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b01
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b01` AS select `j2`.`k2` AS `k2`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`b` AS `b`,`j3`.`k3` AS `k3`,`j3`.`c` AS `c` from ((`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) join `j3` on((`j2`.`k2` = `j3`.`k2`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b01' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| k3
< row
: bytes
| c
> exec
| CREATE VIEW b02 AS SELECT * FROM j3 JOIN (j1 JOIN j2 USING (k1)) USING (k2)
< ok 0 0
> query
| SHOW CREATE VIEW b02
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b02
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b02` AS select `j3`.`k2` AS `k2`,`j3`.`k3` AS `k3`,`j3`.`c` AS `c`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`b` AS `b` from (`j3` join (`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) on((`j3`.`k2` = `j2`.`k2`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b02' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k3
< row
: bytes
| c
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| b
> exec
| CREATE VIEW b03 AS SELECT j1.a, k2, c FROM (j1 JOIN j2 ON j1.k1 = j2.k1) JOIN j3 USING (k2)
< ok 0 0
> query
| SHOW CREATE VIEW b03
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b03
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b03` AS select `j1`.`a` AS `a`,`j2`.`k2` AS `k2`,`j3`.`c` AS `c` from ((`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) join `j3` on((`j2`.`k2` = `j3`.`k2`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b03' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
< row
: bytes
| k2
< row
: bytes
| c
> exec
| CREATE VIEW b04 AS SELECT * FROM (j1 NATURAL JOIN j2) NATURAL JOIN j3
< ok 0 0
> query
| SHOW CREATE VIEW b04
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b04
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b04` AS select `j2`.`k2` AS `k2`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`b` AS `b`,`j3`.`k3` AS `k3`,`j3`.`c` AS `c` from ((`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) join `j3` on((`j2`.`k2` = `j3`.`k2`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b04' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| k3
< row
: bytes
| c
> exec
| CREATE VIEW b05 AS SELECT * FROM j3 NATURAL JOIN (j1 NATURAL JOIN j2)
< ok 0 0
> query
| SHOW CREATE VIEW b05
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b05
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b05` AS select `j3`.`k2` AS `k2`,`j3`.`k3` AS `k3`,`j3`.`c` AS `c`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`b` AS `b` from (`j3` join (`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) on((`j3`.`k2` = `j2`.`k2`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b05' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k3
< row
: bytes
| c
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| b
> exec
| CREATE VIEW b06 AS SELECT * FROM j1 RIGHT JOIN j2 USING (k1)
< ok 0 0
> query
| SHOW CREATE VIEW b06
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b06
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b06` AS select `j2`.`k1` AS `k1`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b`,`j1`.`a` AS `a` from (`j2` left join `j1` on((`j2`.`k1` = `j1`.`k1`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b06' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| k2
< row
: bytes
| b
< row
: bytes
| a
> exec
| CREATE VIEW b07 AS SELECT k1 FROM j1 RIGHT JOIN j2 USING (k1)
< ok 0 0
> query
| SHOW CREATE VIEW b07
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b07
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b07` AS select `j2`.`k1` AS `k1` from (`j2` left join `j1` on((`j2`.`k1` = `j1`.`k1`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b07' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
> exec
| CREATE VIEW b08 AS SELECT * FROM (j1 RIGHT JOIN j2 USING (k1)) JOIN j3 USING (k2)
< ok 0 0
> query
| SHOW CREATE VIEW b08
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b08
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b08` AS select `j2`.`k2` AS `k2`,`j2`.`k1` AS `k1`,`j2`.`b` AS `b`,`j1`.`a` AS `a`,`j3`.`k3` AS `k3`,`j3`.`c` AS `c` from ((`j2` left join `j1` on((`j2`.`k1` = `j1`.`k1`))) join `j3` on((`j2`.`k2` = `j3`.`k2`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b08' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k1
< row
: bytes
| b
< row
: bytes
| a
< row
: bytes
| k3
< row
: bytes
| c
> exec
| CREATE VIEW b09 AS SELECT * FROM (j1 JOIN j2 USING (k1)) RIGHT JOIN j3 USING (k2)
< ok 0 0
> query
| SHOW CREATE VIEW b09
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b09
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b09` AS select `j3`.`k2` AS `k2`,`j3`.`k3` AS `k3`,`j3`.`c` AS `c`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`b` AS `b` from (`j3` left join (`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) on((`j3`.`k2` = `j2`.`k2`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b09' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k3
< row
: bytes
| c
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| b
> exec
| CREATE VIEW b10 AS SELECT * FROM j1 NATURAL RIGHT JOIN j2
< ok 0 0
> query
| SHOW CREATE VIEW b10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b10` AS select `j2`.`k1` AS `k1`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b`,`j1`.`a` AS `a` from (`j2` left join `j1` on((`j2`.`k1` = `j1`.`k1`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b10' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| k2
< row
: bytes
| b
< row
: bytes
| a
> exec
| CREATE VIEW b11 AS SELECT * FROM ((j1 JOIN j2 USING (k1)) JOIN j3 USING (k2)) JOIN j4 USING (k3)
< ok 0 0
> query
| SHOW CREATE VIEW b11
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b11
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b11` AS select `j3`.`k3` AS `k3`,`j2`.`k2` AS `k2`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`b` AS `b`,`j3`.`c` AS `c`,`j4`.`d` AS `d` from (((`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) join `j3` on((`j2`.`k2` = `j3`.`k2`))) join `j4` on((`j3`.`k3` = `j4`.`k3`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b11' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k3
< row
: bytes
| k2
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| c
< row
: bytes
| d
> exec
| CREATE VIEW b12 AS SELECT * FROM (SELECT k1, a FROM j1) d JOIN j2 USING (k1)
< ok 0 0
> query
| SHOW CREATE VIEW b12
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b12
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b12` AS select `d`.`k1` AS `k1`,`d`.`a` AS `a`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b` from ((select `j1`.`k1` AS `k1`,`j1`.`a` AS `a` from `j1`) `d` join `j2` on((`d`.`k1` = `j2`.`k1`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b12' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| k2
< row
: bytes
| b
> exec
| CREATE VIEW b13 AS SELECT * FROM j2 JOIN (SELECT k2, c FROM j3) d USING (k2)
< ok 0 0
> query
| SHOW CREATE VIEW b13
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b13
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b13` AS select `j2`.`k2` AS `k2`,`j2`.`k1` AS `k1`,`j2`.`b` AS `b`,`d`.`c` AS `c` from (`j2` join (select `j3`.`k2` AS `k2`,`j3`.`c` AS `c` from `j3`) `d` on((`j2`.`k2` = `d`.`k2`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b13' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k2
< row
: bytes
| k1
< row
: bytes
| b
< row
: bytes
| c
> exec
| CREATE VIEW b14 AS SELECT * FROM (SELECT k1, a FROM j1) d NATURAL JOIN j2
< ok 0 0
> query
| SHOW CREATE VIEW b14
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b14
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b14` AS select `d`.`k1` AS `k1`,`d`.`a` AS `a`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b` from ((select `j1`.`k1` AS `k1`,`j1`.`a` AS `a` from `j1`) `d` join `j2` on((`d`.`k1` = `j2`.`k1`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b14' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| k2
< row
: bytes
| b
> exec
| CREATE VIEW b15 AS SELECT * FROM j5 JOIN (SELECT k1, k2 FROM j2) d USING (k2, k1)
< ok 0 0
> query
| SHOW CREATE VIEW b15
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b15
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b15` AS select `j5`.`k1` AS `k1`,`j5`.`k2` AS `k2`,`j5`.`a` AS `a`,`j5`.`b` AS `b` from (`j5` join (select `j2`.`k1` AS `k1`,`j2`.`k2` AS `k2` from `j2`) `d` on(((`j5`.`k1` = `d`.`k1`) and (`j5`.`k2` = `d`.`k2`))))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b15' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| k2
< row
: bytes
| a
< row
: bytes
| b
> exec
| CREATE VIEW b16 AS SELECT * FROM j1 RIGHT JOIN j2 USING (k1) WHERE k1 > 0 ORDER BY k1
< ok 0 0
> query
| SHOW CREATE VIEW b16
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b16
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b16` AS select `j2`.`k1` AS `k1`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b`,`j1`.`a` AS `a` from (`j2` left join `j1` on((`j2`.`k1` = `j1`.`k1`))) where (`j2`.`k1` > 0) order by `j2`.`k1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b16' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| k2
< row
: bytes
| b
< row
: bytes
| a
> exec
| CREATE VIEW b17 AS SELECT * FROM u JOIN ((q RIGHT JOIN p ON 1=1) JOIN s ON 1=1) USING (k)
< ok 0 0
> query
| SHOW CREATE VIEW b17
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b17
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b17` AS select `u`.`k` AS `k`,`u`.`u1` AS `u1`,`p`.`p1` AS `p1`,`q`.`q1` AS `q1`,`q`.`q2` AS `q2`,`s`.`s1` AS `s1`,`s`.`s2` AS `s2` from (`u` join ((`p` left join `q` on((1 = 1))) join `s` on((1 = 1))) on((`u`.`k` = `p`.`k`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b17' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k
< row
: bytes
| u1
< row
: bytes
| p1
< row
: bytes
| q1
< row
: bytes
| q2
< row
: bytes
| s1
< row
: bytes
| s2
> exec
| CREATE VIEW b18 AS SELECT * FROM (p JOIN r USING (k)) RIGHT JOIN q ON 1=1
< ok 0 0
> query
| SHOW CREATE VIEW b18
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b18
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b18` AS select `p`.`k` AS `k`,`p`.`p1` AS `p1`,`r`.`r1` AS `r1`,`q`.`q1` AS `q1`,`q`.`q2` AS `q2` from (`q` left join (`p` join `r` on((`p`.`k` = `r`.`k`))) on((1 = 1)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b18' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k
< row
: bytes
| p1
< row
: bytes
| r1
< row
: bytes
| q1
< row
: bytes
| q2
> exec
| CREATE VIEW b19 AS SELECT * FROM j1 JOIN j2 USING (K1)
< ok 0 0
> query
| SHOW CREATE VIEW b19
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b19
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b19` AS select `j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b` from (`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b19' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| k2
< row
: bytes
| b
> exec
| CREATE VIEW b20 AS SELECT * FROM j4, j1 JOIN j2 USING (k1)
< ok 0 0
> query
| SHOW CREATE VIEW b20
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b20
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b20` AS select `j4`.`k3` AS `k3`,`j4`.`d` AS `d`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b` from (`j4` join (`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b20' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k3
< row
: bytes
| d
< row
: bytes
| k1
< row
: bytes
| a
< row
: bytes
| k2
< row
: bytes
| b
> exec
| CREATE VIEW b21 AS SELECT * FROM j4 JOIN j1 ON j4.d = j1.a JOIN j2 USING (k1)
< ok 0 0
> query
| SHOW CREATE VIEW b21
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b21
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b21` AS select `j1`.`k1` AS `k1`,`j4`.`k3` AS `k3`,`j4`.`d` AS `d`,`j1`.`a` AS `a`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b` from ((`j4` join `j1` on((`j4`.`d` = `j1`.`a`))) join `j2` on((`j1`.`k1` = `j2`.`k1`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b21' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
< row
: bytes
| k3
< row
: bytes
| d
< row
: bytes
| a
< row
: bytes
| k2
< row
: bytes
| b
> exec
| CREATE VIEW b22 AS SELECT * FROM u JOIN (q RIGHT JOIN (s RIGHT JOIN p ON 1=1) ON 1=1) USING (k)
< ok 0 0
> query
| SHOW CREATE VIEW b22
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b22
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b22` AS select `u`.`k` AS `k`,`u`.`u1` AS `u1`,`p`.`p1` AS `p1`,`s`.`s1` AS `s1`,`s`.`s2` AS `s2`,`q`.`q1` AS `q1`,`q`.`q2` AS `q2` from (`u` join ((`p` left join `s` on((1 = 1))) left join `q` on((1 = 1))) on((`u`.`k` = `p`.`k`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b22' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k
< row
: bytes
| u1
< row
: bytes
| p1
< row
: bytes
| s1
< row
: bytes
| s2
< row
: bytes
| q1
< row
: bytes
| q2
> exec
| CREATE VIEW b23 AS SELECT * FROM p NATURAL RIGHT JOIN q
< ok 0 0
> query
| SHOW CREATE VIEW b23
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b23
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b23` AS select `q`.`q1` AS `q1`,`q`.`q2` AS `q2`,`p`.`k` AS `k`,`p`.`p1` AS `p1` from (`q` left join `p` on(true))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b23' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| q1
< row
: bytes
| q2
< row
: bytes
| k
< row
: bytes
| p1
> exec
| CREATE VIEW b24 AS SELECT * FROM u NATURAL JOIN (p NATURAL RIGHT JOIN q)
< ok 0 0
> query
| SHOW CREATE VIEW b24
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b24
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b24` AS select `u`.`k` AS `k`,`u`.`u1` AS `u1`,`q`.`q1` AS `q1`,`q`.`q2` AS `q2`,`p`.`p1` AS `p1` from (`u` join (`q` left join `p` on(true)) on((`u`.`k` = `p`.`k`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b24' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k
< row
: bytes
| u1
< row
: bytes
| q1
< row
: bytes
| q2
< row
: bytes
| p1
> exec
| CREATE VIEW b25 AS SELECT * FROM (q NATURAL LEFT JOIN p) NATURAL JOIN u
< ok 0 0
> query
| SHOW CREATE VIEW b25
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b25
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b25` AS select `p`.`k` AS `k`,`q`.`q1` AS `q1`,`q`.`q2` AS `q2`,`p`.`p1` AS `p1`,`u`.`u1` AS `u1` from ((`q` left join `p` on(true)) join `u` on((`p`.`k` = `u`.`k`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'b25' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k
< row
: bytes
| q1
< row
: bytes
| q2
< row
: bytes
| p1
< row
: bytes
| u1
