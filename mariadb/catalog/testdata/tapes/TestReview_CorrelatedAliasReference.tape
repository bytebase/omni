-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_CorrelatedAliasReference default#1
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
| CREATE VIEW v_cr0 (x) AS SELECT a AS z FROM t HAVING (SELECT z) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr0` AS select `t`.`a` AS `x` from `t` having (select `x`) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr0
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr0
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr0` AS select `t`.`a` AS `x` from `t` having (select `x`) > 0
< ok 0 0
> query
| SELECT * FROM v_cr0
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr1 AS SELECT a AS z FROM t HAVING (SELECT z) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr1` AS select `t`.`a` AS `z` from `t` having (select `z`) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr1
< result
< column "z" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr1
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr1` AS select `t`.`a` AS `z` from `t` having (select `z`) > 0
< ok 0 0
> query
| SELECT * FROM v_cr1
< result
< column "z" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr2 (x) AS SELECT a AS z FROM t HAVING (SELECT z + 0) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr2` AS select `t`.`a` AS `x` from `t` having (select `x` + 0) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr2
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr2
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr2` AS select `t`.`a` AS `x` from `t` having (select `x` + 0) > 0
< ok 0 0
> query
| SELECT * FROM v_cr2
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr3 (x) AS SELECT a AS z FROM t HAVING (SELECT MAX(z)) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr3` AS select `t`.`a` AS `x` from `t` having (select max(`x`)) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr3
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr3
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr3` AS select `t`.`a` AS `x` from `t` having (select max(`x`)) > 0
< ok 0 0
> query
| SELECT * FROM v_cr3
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr4 (x) AS SELECT a AS z FROM t HAVING (SELECT (SELECT z)) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr4` AS select `t`.`a` AS `x` from `t` having (select (select `x`)) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr4
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr4
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr4` AS select `t`.`a` AS `x` from `t` having (select (select `x`)) > 0
< ok 0 0
> query
| SELECT * FROM v_cr4
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr5 (x) AS SELECT a AS z FROM t HAVING (SELECT COUNT(*) FROM t t2 WHERE t2.a = z) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr5` AS select `t`.`a` AS `x` from `t` having (select count(0) from `t` `t2` where `t2`.`a` = `x`) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr5
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr5
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr5` AS select `t`.`a` AS `x` from `t` having (select count(0) from `t` `t2` where `t2`.`a` = `x`) > 0
< ok 0 0
> query
| SELECT * FROM v_cr5
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr6 (x) AS SELECT a AS z FROM t HAVING (SELECT 1 FROM t t2 GROUP BY z) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr6` AS select `t`.`a` AS `x` from `t` having (select 1 from `t` `t2` group by `x`) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr6
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr6
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr6` AS select `t`.`a` AS `x` from `t` having (select 1 from `t` `t2` group by `x`) > 0
< ok 0 0
> query
| SELECT * FROM v_cr6
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr7 (x) AS SELECT a AS z FROM t HAVING (SELECT 1 FROM t t2 ORDER BY z LIMIT 1) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr7` AS select `t`.`a` AS `x` from `t` having (select 1 from `t` `t2` order by `x` limit 1) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr7
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr7
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr7` AS select `t`.`a` AS `x` from `t` having (select 1 from `t` `t2` order by `x` limit 1) > 0
< ok 0 0
> query
| SELECT * FROM v_cr7
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr8 (x) AS SELECT a AS z FROM t HAVING (SELECT z UNION SELECT z) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr8
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr8
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr8` AS select `t`.`a` AS `x` from `t` having (select `x` union select `x`) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr8
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr8
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr8` AS select `t`.`a` AS `x` from `t` having (select `x` union select `x`) > 0
< ok 0 0
> query
| SELECT * FROM v_cr8
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr9 (x, y) AS SELECT a AS z, b AS w FROM t HAVING (SELECT z + w) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr9
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr9
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr9` AS select `t`.`a` AS `x`,`t`.`b` AS `y` from `t` having (select `x` + `y`) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr9
< result
< column "x" "INT"
< column "y" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| DROP VIEW v_cr9
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr9` AS select `t`.`a` AS `x`,`t`.`b` AS `y` from `t` having (select `x` + `y`) > 0
< ok 0 0
> query
| SELECT * FROM v_cr9
< result
< column "x" "INT"
< column "y" "INT"
< row
: int64 1
: int64 2
< row
: int64 2
: int64 1
> exec
| CREATE VIEW v_cr10 (x) AS SELECT a AS z FROM t HAVING (SELECT z HAVING z = 1) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr10` AS select `t`.`a` AS `x` from `t` having (select `x` having `x` = 1) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr10
< result
< column "x" "INT"
< row
: int64 1
> exec
| DROP VIEW v_cr10
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr10` AS select `t`.`a` AS `x` from `t` having (select `x` having `x` = 1) > 0
< ok 0 0
> query
| SELECT * FROM v_cr10
< result
< column "x" "INT"
< row
: int64 1
> exec
| CREATE VIEW v_cr11 (x) AS SELECT a AS z FROM t HAVING EXISTS(SELECT 1 FROM (SELECT 1 AS z HAVING z=1) AS d)
< ok 0 0
> query
| SHOW CREATE VIEW v_cr11
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr11
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr11` AS select `t`.`a` AS `x` from `t` having exists(select 1 from (select 1 AS `z` having `z` = 1) `d` limit 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr11
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr11
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr11` AS select `t`.`a` AS `x` from `t` having exists(select 1 from (select 1 AS `z` having `z` = 1) `d` limit 1)
< ok 0 0
> query
| SELECT * FROM v_cr11
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr12 (x) AS SELECT a AS z FROM t HAVING (SELECT 2 AS z HAVING z = 2) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr12
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr12
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr12` AS select `t`.`a` AS `x` from `t` having (select 2 AS `z` having `z` = 2) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr12
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr12
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr12` AS select `t`.`a` AS `x` from `t` having (select 2 AS `z` having `z` = 2) > 0
< ok 0 0
> query
| SELECT * FROM v_cr12
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr13 (x) AS SELECT a AS z FROM t HAVING (SELECT 1 AS z FROM t t2 GROUP BY z) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr13
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr13
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr13` AS select `t`.`a` AS `x` from `t` having (select 1 AS `z` from `t` `t2` group by '') > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr13
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr13
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr13` AS select `t`.`a` AS `x` from `t` having (select 1 AS `z` from `t` `t2` group by '') > 0
< ok 0 0
> query
| SELECT * FROM v_cr13
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr14 (x) AS SELECT a AS z FROM t HAVING (SELECT 1 AS z FROM t t2 ORDER BY z LIMIT 1) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr14
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr14
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr14` AS select `t`.`a` AS `x` from `t` having (select 1 AS `z` from `t` `t2` order by '' limit 1) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr14
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr14
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr14` AS select `t`.`a` AS `x` from `t` having (select 1 AS `z` from `t` `t2` order by '' limit 1) > 0
< ok 0 0
> query
| SELECT * FROM v_cr14
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr15 (x) AS SELECT a AS z FROM t HAVING (SELECT 1 AS z UNION SELECT 2 ORDER BY z DESC LIMIT 1) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr15
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr15
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr15` AS select `t`.`a` AS `x` from `t` having (select 1 AS `z` union select 2 order by `z` desc limit 1) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr15
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr15
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr15` AS select `t`.`a` AS `x` from `t` having (select 1 AS `z` union select 2 order by `z` desc limit 1) > 0
< ok 0 0
> query
| SELECT * FROM v_cr15
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr16 (x) AS SELECT a AS z FROM t HAVING (SELECT 1 AS z UNION SELECT 2 ORDER BY 1 LIMIT 1) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr16
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr16
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr16` AS select `t`.`a` AS `x` from `t` having (select 1 AS `z` union select 2 order by 1 limit 1) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr16
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr16
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr16` AS select `t`.`a` AS `x` from `t` having (select 1 AS `z` union select 2 order by 1 limit 1) > 0
< ok 0 0
> query
| SELECT * FROM v_cr16
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr17 (x) AS SELECT a AS z FROM t HAVING (SELECT 1 AS z FROM t t2 UNION SELECT 2 ORDER BY z + 0 LIMIT 1) > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr17
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr17
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr17` AS select `t`.`a` AS `x` from `t` having (select 1 AS `z` from `t` `t2` union select 2 order by `z` + 0 limit 1) > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr17
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr17
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr17` AS select `t`.`a` AS `x` from `t` having (select 1 AS `z` from `t` `t2` union select 2 order by `z` + 0 limit 1) > 0
< ok 0 0
> query
| SELECT * FROM v_cr17
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr18 (x) AS SELECT a + 1 AS z FROM t ORDER BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr18
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr18
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr18` AS select `t`.`a` + 1 AS `x` from `t` order by `x` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr18
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| DROP VIEW v_cr18
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr18` AS select `t`.`a` + 1 AS `x` from `t` order by `x` + 0
< ok 0 0
> query
| SELECT * FROM v_cr18
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| CREATE VIEW v_cr19 AS SELECT a + 1 AS z FROM t ORDER BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr19
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr19
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr19` AS select `t`.`a` + 1 AS `z` from `t` order by `z` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr19
< result
< column "z" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| DROP VIEW v_cr19
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr19` AS select `t`.`a` + 1 AS `z` from `t` order by `z` + 0
< ok 0 0
> query
| SELECT * FROM v_cr19
< result
< column "z" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| CREATE VIEW v_cr20 (x) AS SELECT a + 1 AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW v_cr20
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr20
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr20` AS select `t`.`a` + 1 AS `x` from `t` order by `t`.`a` + 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr20
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| DROP VIEW v_cr20
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr20` AS select `t`.`a` + 1 AS `x` from `t` order by `t`.`a` + 1
< ok 0 0
> query
| SELECT * FROM v_cr20
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| CREATE VIEW v_cr21 (x) AS SELECT a + 1 AS z FROM t ORDER BY ABS(z)
< ok 0 0
> query
| SHOW CREATE VIEW v_cr21
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr21
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr21` AS select `t`.`a` + 1 AS `x` from `t` order by abs(`x`)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr21
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| DROP VIEW v_cr21
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr21` AS select `t`.`a` + 1 AS `x` from `t` order by abs(`x`)
< ok 0 0
> query
| SELECT * FROM v_cr21
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| CREATE VIEW v_cr22 (x) AS SELECT 2 AS z FROM t ORDER BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr22
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr22
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr22` AS select 2 AS `x` from `t` order by `x` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr22
< result
< column "x" "INT"
< row
: int64 2
< row
: int64 2
> exec
| DROP VIEW v_cr22
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr22` AS select 2 AS `x` from `t` order by `x` + 0
< ok 0 0
> query
| SELECT * FROM v_cr22
< result
< column "x" "INT"
< row
: int64 2
< row
: int64 2
> exec
| CREATE VIEW v_cr23 (x, y) AS SELECT a + 1 AS z, COUNT(*) AS n FROM t GROUP BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr23
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr23
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr23` AS select `t`.`a` + 1 AS `x`,count(0) AS `y` from `t` group by `t`.`a` + 1 + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr23
< result
< column "x" "BIGINT"
< column "y" "BIGINT"
< row
: int64 2
: int64 1
< row
: int64 3
: int64 1
> exec
| DROP VIEW v_cr23
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr23` AS select `t`.`a` + 1 AS `x`,count(0) AS `y` from `t` group by `t`.`a` + 1 + 0
< ok 0 0
> query
| SELECT * FROM v_cr23
< result
< column "x" "BIGINT"
< column "y" "BIGINT"
< row
: int64 2
: int64 1
< row
: int64 3
: int64 1
> exec
| CREATE VIEW v_cr24 (x) AS SELECT a + 1 AS z FROM t HAVING z > 0 ORDER BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr24
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr24
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr24` AS select `t`.`a` + 1 AS `x` from `t` having `x` > 0 order by `x` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr24
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| DROP VIEW v_cr24
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr24` AS select `t`.`a` + 1 AS `x` from `t` having `x` > 0 order by `x` + 0
< ok 0 0
> query
| SELECT * FROM v_cr24
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| CREATE VIEW v_cr25 (x) AS SELECT a + 1 AS z FROM t ORDER BY (SELECT z) + 0
< ok 0 0
> query
| SHOW CREATE VIEW v_cr25
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr25
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr25` AS select `t`.`a` + 1 AS `x` from `t` order by (select `x`) + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr25
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| DROP VIEW v_cr25
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr25` AS select `t`.`a` + 1 AS `x` from `t` order by (select `x`) + 0
< ok 0 0
> query
| SELECT * FROM v_cr25
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| CREATE VIEW v_cr26 AS WITH c(x) AS (SELECT a AS z FROM t HAVING (SELECT z) > 0) SELECT c.x FROM c
< ok 0 0
> query
| SHOW CREATE VIEW v_cr26
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr26
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr26` AS with c(`x`) as (select `t`.`a` AS `x` from `t` having (select `x`) > 0)select `c`.`x` AS `x` from `c`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr26
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_cr26
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr26` AS with c(`x`) as (select `t`.`a` AS `x` from `t` having (select `x`) > 0)select `c`.`x` AS `x` from `c`
< ok 0 0
> query
| SELECT * FROM v_cr26
< result
< column "x" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_cr27 AS WITH c(x) AS (SELECT a + 1 AS z FROM t ORDER BY z + 0) SELECT c.x FROM c
< ok 0 0
> query
| SHOW CREATE VIEW v_cr27
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cr27
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr27` AS with c(`x`) as (select `t`.`a` + 1 AS `x` from `t` order by `x` + 0)select `c`.`x` AS `x` from `c`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_cr27
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3
> exec
| DROP VIEW v_cr27
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cr27` AS with c(`x`) as (select `t`.`a` + 1 AS `x` from `t` order by `x` + 0)select `c`.`x` AS `x` from `c`
< ok 0 0
> query
| SELECT * FROM v_cr27
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3
