-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestViewNamesAndJoinColumnsOracle default#1
> exec
| DROP DATABASE IF EXISTS omni_view_names
< ok 0 0
> exec
| CREATE DATABASE omni_view_names
< ok 0 0
> exec
| USE omni_view_names
< ok 0 0
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
| CREATE VIEW a04 AS SELECT (SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1), (SELECT DISTINCT a FROM n1 LIMIT 1) FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a04
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1)
< row
: bytes
| (SELECT DISTINCT a FROM n1 LIMIT 1)
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a04` (`(SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1)`, `(SELECT DISTINCT a FROM n1 LIMIT 1)`) AS SELECT (SELECT `a` AS `a` FROM `omni_view_names`.`n1` GROUP BY `a` HAVING `a`>0 ORDER BY `a` LIMIT 1) AS `(SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1)`,(SELECT DISTINCT `a` AS `a` FROM `omni_view_names`.`n1` LIMIT 1) AS `(SELECT DISTINCT a FROM n1 LIMIT 1)` FROM `omni_view_names`.`n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a05 AS SELECT a+1, a  +  1, count(*), cast(a as char) FROM n1 GROUP BY a
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a05
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a05` (`a+1`, `a  +  1`, `count(*)`, `cast(a as char)`) AS SELECT `a`+1 AS `a+1`,`a`+1 AS `a  +  1`,COUNT(1) AS `count(*)`,CAST(`a` AS CHAR(11)) AS `cast(a as char)` FROM `omni_view_names`.`n1` GROUP BY `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a06 AS SELECT (a), ((b)), (1), 0004, - 2, x'41', B'10', true, date '2020-01-01', ('e'), (null) FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a06
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (a)
< row
: bytes
| ((b))
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
| TRUE
< row
: bytes
| date '2020-01-01'
< row
: bytes
| e
< row
: bytes
| NULL
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a06` (`(a)`, `((b))`, `1`, `0004`, `- 2`, `x'41'`, `B'10'`, `TRUE`, `date '2020-01-01'`, `e`, `NULL`) AS SELECT (`a`) AS `(a)`,((`b`)) AS `((b))`,(1) AS `1`,4 AS `0004`,-2 AS `- 2`,x'41' AS `x'41'`,b'10' AS `B'10'`,TRUE AS `TRUE`,DATE '2020-01-01' AS `date '2020-01-01'`,(_UTF8MB4'e') AS `e`,(NULL) AS `NULL` FROM `omni_view_names`.`n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a07 AS SELECT a /* c */ + 1, a/*!50100 +2*/ FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a07
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a /* c  + 1
< row
: bytes
| a +2
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a07` (`a /* c  + 1`, `a +2`) AS SELECT `a`+1 AS `a /* c  + 1`,`a`+2 AS `a +2` FROM `omni_view_names`.`n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a08 AS SELECT concat('éééééééééééééééééééééééééééééééééééééééééééééééééé') FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a08
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| name_exp_1
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a08` (`name_exp_1`) AS SELECT CONCAT(_UTF8MB4'éééééééééééééééééééééééééééééééééééééééééééééééééé') AS `concat('éééééééééééééééééééééééééééééééééééééééééééééééééé')` FROM `omni_view_names`.`n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a09 AS SELECT a, concat('eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee') FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a09
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
< row
: bytes
| name_exp_2
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a09` (`a`, `name_exp_2`) AS SELECT `a` AS `a`,CONCAT(_UTF8MB4'eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee') AS `concat('eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee')` FROM `omni_view_names`.`n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a10 AS SELECT * FROM (SELECT a+1, b FROM n1) d
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a10
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a+1
< row
: bytes
| b
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a10` (`a+1`, `b`) AS SELECT `d`.`a+1` AS `a+1`,`omni_view_names`.`d`.`b` AS `b` FROM (SELECT `a`+1 AS `a+1`,`b` AS `b` FROM `omni_view_names`.`n1`) AS `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a11 AS SELECT (select   1), (SELECT a FROM n1 LIMIT 1) + 0 FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a11
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| (select   1)
< row
: bytes
| (SELECT a FROM n1 LIMIT 1) + 0
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a11` (`(select   1)`, `(SELECT a FROM n1 LIMIT 1) + 0`) AS SELECT (SELECT 1 AS `1`) AS `(select   1)`,(SELECT `a` AS `a` FROM `omni_view_names`.`n1` LIMIT 1)+0 AS `(SELECT a FROM n1 LIMIT 1) + 0` FROM `omni_view_names`.`n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a12 AS SELECT a IN (SELECT a FROM n1), a in (1,2) FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a12
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a IN (SELECT a FROM n1)
< row
: bytes
| a in (1,2)
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a12` (`a IN (SELECT a FROM n1)`, `a in (1,2)`) AS SELECT `a` IN (SELECT `a` AS `a` FROM `omni_view_names`.`n1`) AS `a IN (SELECT a FROM n1)`,`a` IN (1,2) AS `a in (1,2)` FROM `omni_view_names`.`n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a13 AS SELECT a + 1 AS `x + y`, count(*) AS `count(z)` FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a13
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| x + y
< row
: bytes
| count(z)
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a13` (`x + y`, `count(z)`) AS SELECT `a`+1 AS `x + y`,COUNT(1) AS `count(z)` FROM `omni_view_names`.`n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a14 AS SELECT EXISTS (SELECT 1 UNION SELECT 2), a IN (SELECT 1 UNION SELECT 2) FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a14
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| EXISTS (SELECT 1 UNION SELECT 2)
< row
: bytes
| a IN (SELECT 1 UNION SELECT 2)
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a14` (`EXISTS (SELECT 1 UNION SELECT 2)`, `a IN (SELECT 1 UNION SELECT 2)`) AS SELECT EXISTS (SELECT 1 AS `1` UNION SELECT 2 AS `2`) AS `EXISTS (SELECT 1 UNION SELECT 2)`,`a` IN (SELECT 1 AS `1` UNION SELECT 2 AS `2`) AS `a IN (SELECT 1 UNION SELECT 2)` FROM `omni_view_names`.`n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a15 AS SELECT * FROM (SELECT `a`+1 FROM n1) d
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a15
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| `a`+1
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a15` (```a``+1`) AS SELECT `d`.```a``+1` AS ```a``+1` FROM (SELECT `a`+1 AS ```a``+1` FROM `omni_view_names`.`n1`) AS `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a16 AS SELECT * FROM (SELECT a AS `x``y` FROM n1) d
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a16
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| x`y
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a16` (`x``y`) AS SELECT `omni_view_names`.`d`.`x``y` AS `x``y` FROM (SELECT `a` AS `x``y` FROM `omni_view_names`.`n1`) AS `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a17 AS SELECT a */*c*/3, a + '*/', concat('/**/!x', b), a /*!50100 */*/*c*/3, a + 1 -- t */
| , b FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a17
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a *c3
< row
: bytes
| a + ''
< row
: bytes
| concat('/*!x', b)
< row
: bytes
| a  *c3
< row
: bytes
| a + 1 -- t 
< row
: bytes
| b
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a17` (`a *c3`, `a + ''`, `concat('/*!x', b)`, `a  *c3`, `a + 1 -- t `, `b`) AS SELECT `a`*3 AS `a *c3`,`a`+_UTF8MB4'*/' AS `a + ''`,CONCAT(_UTF8MB4'/**/!x', `b`) AS `concat('/*!x', b)`,`a`*3 AS `a  *c3`,`a`+1 AS `a + 1 -- t `,`b` AS `b` FROM `omni_view_names`.`n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW a18 AS SELECT concat('xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx*/*/*/*/') FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| a18
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| concat('xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx')
> query
| SHOW CREATE VIEW a18
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a18
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a18` (`concat('xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx')`) AS SELECT CONCAT(_UTF8MB4'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx*/*/*/*/') AS `concat('xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx')` FROM `omni_view_names`.`n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b01 AS SELECT * FROM (j1 JOIN j2 USING (k1)) JOIN j3 USING (k2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b01
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b01` (`k2`, `k1`, `a`, `b`, `k3`, `c`) AS SELECT `omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j1`.`k1` AS `k1`,`omni_view_names`.`j1`.`a` AS `a`,`omni_view_names`.`j2`.`b` AS `b`,`omni_view_names`.`j3`.`k3` AS `k3`,`omni_view_names`.`j3`.`c` AS `c` FROM (`omni_view_names`.`j1` JOIN `omni_view_names`.`j2` USING (`k1`)) JOIN `omni_view_names`.`j3` USING (`k2`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b02 AS SELECT * FROM j3 JOIN (j1 JOIN j2 USING (k1)) USING (k2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b02
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b02` (`k2`, `k3`, `c`, `k1`, `a`, `b`) AS SELECT `omni_view_names`.`j3`.`k2` AS `k2`,`omni_view_names`.`j3`.`k3` AS `k3`,`omni_view_names`.`j3`.`c` AS `c`,`omni_view_names`.`j1`.`k1` AS `k1`,`omni_view_names`.`j1`.`a` AS `a`,`omni_view_names`.`j2`.`b` AS `b` FROM `omni_view_names`.`j3` JOIN (`omni_view_names`.`j1` JOIN `omni_view_names`.`j2` USING (`k1`)) USING (`k2`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b03 AS SELECT j1.a, k2, c FROM (j1 JOIN j2 ON j1.k1 = j2.k1) JOIN j3 USING (k2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b03
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b03` (`a`, `k2`, `c`) AS SELECT `j1`.`a` AS `a`,`k2` AS `k2`,`c` AS `c` FROM (`omni_view_names`.`j1` JOIN `omni_view_names`.`j2` ON `j1`.`k1`=`j2`.`k1`) JOIN `omni_view_names`.`j3` USING (`k2`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b04 AS SELECT * FROM (j1 NATURAL JOIN j2) NATURAL JOIN j3
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b04
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b04` (`k2`, `k1`, `a`, `b`, `k3`, `c`) AS SELECT `omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j1`.`k1` AS `k1`,`omni_view_names`.`j1`.`a` AS `a`,`omni_view_names`.`j2`.`b` AS `b`,`omni_view_names`.`j3`.`k3` AS `k3`,`omni_view_names`.`j3`.`c` AS `c` FROM (`omni_view_names`.`j1` NATURAL JOIN `omni_view_names`.`j2`) NATURAL JOIN `omni_view_names`.`j3`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b05 AS SELECT * FROM j3 NATURAL JOIN (j1 NATURAL JOIN j2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b05
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b05` (`k2`, `k3`, `c`, `k1`, `a`, `b`) AS SELECT `omni_view_names`.`j3`.`k2` AS `k2`,`omni_view_names`.`j3`.`k3` AS `k3`,`omni_view_names`.`j3`.`c` AS `c`,`omni_view_names`.`j1`.`k1` AS `k1`,`omni_view_names`.`j1`.`a` AS `a`,`omni_view_names`.`j2`.`b` AS `b` FROM `omni_view_names`.`j3` NATURAL JOIN (`omni_view_names`.`j1` NATURAL JOIN `omni_view_names`.`j2`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b06 AS SELECT * FROM j1 RIGHT JOIN j2 USING (k1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b06
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b06` (`k1`, `k2`, `b`, `a`) AS SELECT `omni_view_names`.`j2`.`k1` AS `k1`,`omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j2`.`b` AS `b`,`omni_view_names`.`j1`.`a` AS `a` FROM `omni_view_names`.`j1` RIGHT JOIN `omni_view_names`.`j2` USING (`k1`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b07 AS SELECT k1 FROM j1 RIGHT JOIN j2 USING (k1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b07
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| k1
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b07` (`k1`) AS SELECT `k1` AS `k1` FROM `omni_view_names`.`j1` RIGHT JOIN `omni_view_names`.`j2` USING (`k1`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b08 AS SELECT * FROM (j1 RIGHT JOIN j2 USING (k1)) JOIN j3 USING (k2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b08
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b08` (`k2`, `k1`, `b`, `a`, `k3`, `c`) AS SELECT `omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j2`.`k1` AS `k1`,`omni_view_names`.`j2`.`b` AS `b`,`omni_view_names`.`j1`.`a` AS `a`,`omni_view_names`.`j3`.`k3` AS `k3`,`omni_view_names`.`j3`.`c` AS `c` FROM (`omni_view_names`.`j1` RIGHT JOIN `omni_view_names`.`j2` USING (`k1`)) JOIN `omni_view_names`.`j3` USING (`k2`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b09 AS SELECT * FROM (j1 JOIN j2 USING (k1)) RIGHT JOIN j3 USING (k2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b09
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b09` (`k2`, `k3`, `c`, `k1`, `a`, `b`) AS SELECT `omni_view_names`.`j3`.`k2` AS `k2`,`omni_view_names`.`j3`.`k3` AS `k3`,`omni_view_names`.`j3`.`c` AS `c`,`omni_view_names`.`j1`.`k1` AS `k1`,`omni_view_names`.`j1`.`a` AS `a`,`omni_view_names`.`j2`.`b` AS `b` FROM (`omni_view_names`.`j1` JOIN `omni_view_names`.`j2` USING (`k1`)) RIGHT JOIN `omni_view_names`.`j3` USING (`k2`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b10 AS SELECT * FROM j1 NATURAL RIGHT JOIN j2
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b10
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b10` (`k1`, `k2`, `b`, `a`) AS SELECT `omni_view_names`.`j2`.`k1` AS `k1`,`omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j2`.`b` AS `b`,`omni_view_names`.`j1`.`a` AS `a` FROM `omni_view_names`.`j1` NATURAL RIGHT JOIN `omni_view_names`.`j2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b11 AS SELECT * FROM ((j1 JOIN j2 USING (k1)) JOIN j3 USING (k2)) JOIN j4 USING (k3)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b11
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b11` (`k3`, `k2`, `k1`, `a`, `b`, `c`, `d`) AS SELECT `omni_view_names`.`j3`.`k3` AS `k3`,`omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j1`.`k1` AS `k1`,`omni_view_names`.`j1`.`a` AS `a`,`omni_view_names`.`j2`.`b` AS `b`,`omni_view_names`.`j3`.`c` AS `c`,`omni_view_names`.`j4`.`d` AS `d` FROM ((`omni_view_names`.`j1` JOIN `omni_view_names`.`j2` USING (`k1`)) JOIN `omni_view_names`.`j3` USING (`k2`)) JOIN `omni_view_names`.`j4` USING (`k3`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b12 AS SELECT * FROM (SELECT k1, a FROM j1) d JOIN j2 USING (k1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b12
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b12` (`k1`, `a`, `k2`, `b`) AS SELECT `omni_view_names`.`d`.`k1` AS `k1`,`omni_view_names`.`d`.`a` AS `a`,`omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j2`.`b` AS `b` FROM (SELECT `k1` AS `k1`,`a` AS `a` FROM `omni_view_names`.`j1`) AS `d` JOIN `omni_view_names`.`j2` USING (`k1`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b13 AS SELECT * FROM j2 JOIN (SELECT k2, c FROM j3) d USING (k2)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b13
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b13` (`k2`, `k1`, `b`, `c`) AS SELECT `omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j2`.`k1` AS `k1`,`omni_view_names`.`j2`.`b` AS `b`,`omni_view_names`.`d`.`c` AS `c` FROM `omni_view_names`.`j2` JOIN (SELECT `k2` AS `k2`,`c` AS `c` FROM `omni_view_names`.`j3`) AS `d` USING (`k2`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b14 AS SELECT * FROM (SELECT k1, a FROM j1) d NATURAL JOIN j2
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b14
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b14` (`k1`, `a`, `k2`, `b`) AS SELECT `omni_view_names`.`d`.`k1` AS `k1`,`omni_view_names`.`d`.`a` AS `a`,`omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j2`.`b` AS `b` FROM (SELECT `k1` AS `k1`,`a` AS `a` FROM `omni_view_names`.`j1`) AS `d` NATURAL JOIN `omni_view_names`.`j2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b15 AS SELECT * FROM j5 JOIN (SELECT k1, k2 FROM j2) d USING (k2, k1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b15
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b15` (`k1`, `k2`, `a`, `b`) AS SELECT `omni_view_names`.`j5`.`k1` AS `k1`,`omni_view_names`.`j5`.`k2` AS `k2`,`omni_view_names`.`j5`.`a` AS `a`,`omni_view_names`.`j5`.`b` AS `b` FROM `omni_view_names`.`j5` JOIN (SELECT `k1` AS `k1`,`k2` AS `k2` FROM `omni_view_names`.`j2`) AS `d` USING (`k2`,`k1`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b16 AS SELECT * FROM j1 RIGHT JOIN j2 USING (k1) WHERE k1 > 0 ORDER BY k1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b16
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b16` (`k1`, `k2`, `b`, `a`) AS SELECT `omni_view_names`.`j2`.`k1` AS `k1`,`omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j2`.`b` AS `b`,`omni_view_names`.`j1`.`a` AS `a` FROM `omni_view_names`.`j1` RIGHT JOIN `omni_view_names`.`j2` USING (`k1`) WHERE `k1`>0 ORDER BY `k1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b17 AS SELECT * FROM u JOIN ((q RIGHT JOIN p ON 1=1) JOIN s ON 1=1) USING (k)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b17
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
< row
: bytes
| s1
< row
: bytes
| s2
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b17` (`k`, `u1`, `q1`, `q2`, `p1`, `s1`, `s2`) AS SELECT `omni_view_names`.`u`.`k` AS `k`,`omni_view_names`.`u`.`u1` AS `u1`,`omni_view_names`.`q`.`q1` AS `q1`,`omni_view_names`.`q`.`q2` AS `q2`,`omni_view_names`.`p`.`p1` AS `p1`,`omni_view_names`.`s`.`s1` AS `s1`,`omni_view_names`.`s`.`s2` AS `s2` FROM `omni_view_names`.`u` JOIN ((`omni_view_names`.`q` RIGHT JOIN `omni_view_names`.`p` ON 1=1) JOIN `omni_view_names`.`s` ON 1=1) USING (`k`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b18 AS SELECT * FROM (p JOIN r USING (k)) RIGHT JOIN q ON 1=1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b18
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b18` (`k`, `p1`, `r1`, `q1`, `q2`) AS SELECT `omni_view_names`.`p`.`k` AS `k`,`omni_view_names`.`p`.`p1` AS `p1`,`omni_view_names`.`r`.`r1` AS `r1`,`omni_view_names`.`q`.`q1` AS `q1`,`omni_view_names`.`q`.`q2` AS `q2` FROM (`omni_view_names`.`p` JOIN `omni_view_names`.`r` USING (`k`)) RIGHT JOIN `omni_view_names`.`q` ON 1=1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b19 AS SELECT * FROM j1 JOIN j2 USING (K1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b19
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b19` (`k1`, `a`, `k2`, `b`) AS SELECT `omni_view_names`.`j1`.`k1` AS `k1`,`omni_view_names`.`j1`.`a` AS `a`,`omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j2`.`b` AS `b` FROM `omni_view_names`.`j1` JOIN `omni_view_names`.`j2` USING (`K1`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b20 AS SELECT * FROM j4, j1 JOIN j2 USING (k1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b20
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b20` (`k3`, `d`, `k1`, `a`, `k2`, `b`) AS SELECT `omni_view_names`.`j4`.`k3` AS `k3`,`omni_view_names`.`j4`.`d` AS `d`,`omni_view_names`.`j1`.`k1` AS `k1`,`omni_view_names`.`j1`.`a` AS `a`,`omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j2`.`b` AS `b` FROM (`omni_view_names`.`j4`) JOIN (`omni_view_names`.`j1` JOIN `omni_view_names`.`j2` USING (`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b21 AS SELECT * FROM j4 JOIN j1 ON j4.d = j1.a JOIN j2 USING (k1)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b21
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b21` (`k1`, `k3`, `d`, `a`, `k2`, `b`) AS SELECT `omni_view_names`.`j1`.`k1` AS `k1`,`omni_view_names`.`j4`.`k3` AS `k3`,`omni_view_names`.`j4`.`d` AS `d`,`omni_view_names`.`j1`.`a` AS `a`,`omni_view_names`.`j2`.`k2` AS `k2`,`omni_view_names`.`j2`.`b` AS `b` FROM (`omni_view_names`.`j4` JOIN `omni_view_names`.`j1` ON `j4`.`d`=`j1`.`a`) JOIN `omni_view_names`.`j2` USING (`k1`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b22 AS SELECT * FROM u JOIN (q RIGHT JOIN (s RIGHT JOIN p ON 1=1) ON 1=1) USING (k)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b22
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
| s1
< row
: bytes
| s2
< row
: bytes
| p1
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b22` (`k`, `u1`, `q1`, `q2`, `s1`, `s2`, `p1`) AS SELECT `omni_view_names`.`u`.`k` AS `k`,`omni_view_names`.`u`.`u1` AS `u1`,`omni_view_names`.`q`.`q1` AS `q1`,`omni_view_names`.`q`.`q2` AS `q2`,`omni_view_names`.`s`.`s1` AS `s1`,`omni_view_names`.`s`.`s2` AS `s2`,`omni_view_names`.`p`.`p1` AS `p1` FROM `omni_view_names`.`u` JOIN (`omni_view_names`.`q` RIGHT JOIN (`omni_view_names`.`s` RIGHT JOIN `omni_view_names`.`p` ON 1=1) ON 1=1) USING (`k`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b23 AS SELECT * FROM p NATURAL RIGHT JOIN q
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b23
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b23` (`q1`, `q2`, `k`, `p1`) AS SELECT `omni_view_names`.`q`.`q1` AS `q1`,`omni_view_names`.`q`.`q2` AS `q2`,`omni_view_names`.`p`.`k` AS `k`,`omni_view_names`.`p`.`p1` AS `p1` FROM `omni_view_names`.`p` NATURAL RIGHT JOIN `omni_view_names`.`q`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b24 AS SELECT * FROM u NATURAL JOIN (p NATURAL RIGHT JOIN q)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b24
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b24` (`k`, `u1`, `q1`, `q2`, `p1`) AS SELECT `omni_view_names`.`u`.`k` AS `k`,`omni_view_names`.`u`.`u1` AS `u1`,`omni_view_names`.`q`.`q1` AS `q1`,`omni_view_names`.`q`.`q2` AS `q2`,`omni_view_names`.`p`.`p1` AS `p1` FROM `omni_view_names`.`u` NATURAL JOIN (`omni_view_names`.`p` NATURAL RIGHT JOIN `omni_view_names`.`q`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW b25 AS SELECT * FROM (q NATURAL LEFT JOIN p) NATURAL JOIN u
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| b25
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b25` (`k`, `q1`, `q2`, `p1`, `u1`) AS SELECT `omni_view_names`.`p`.`k` AS `k`,`omni_view_names`.`q`.`q1` AS `q1`,`omni_view_names`.`q`.`q2` AS `q2`,`omni_view_names`.`p`.`p1` AS `p1`,`omni_view_names`.`u`.`u1` AS `u1` FROM (`omni_view_names`.`q` NATURAL LEFT JOIN `omni_view_names`.`p`) NATURAL JOIN `omni_view_names`.`u`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS omni_view_names
< ok 0 0
