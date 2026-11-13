-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestParenSubqueryOperandOracle default#1
> exec
| DROP DATABASE IF EXISTS omni_paren_operand
< ok 0 0
> exec
| CREATE DATABASE omni_paren_operand
< ok 0 0
> exec
| USE omni_paren_operand
< ok 0 0
> exec
| CREATE TABLE j1 (k1 INT, a INT)
< ok 0 0
> exec
| CREATE TABLE j2 (k1 INT, b INT)
< ok 0 0
> exec
| CREATE TABLE j3 (k1 INT, c INT)
< ok 0 0
> exec
| CREATE TABLE n1 (a INT)
< ok 0 0
> exec
| select (((select 1) + 0))
< ok 0 0
> exec
| select (((select 1) = 0))
< ok 0 0
> exec
| select ((((select 1) = 0)))
< ok 0 0
> exec
| select ((((select 1) = 0)) or (((select 2)) = 2))
< ok 0 0
> exec
| select (((select count(0) from n1) = 0))
< ok 0 0
> exec
| select 1 in (((select 1) = 1))
< ok 0 0
> exec
| select 1 in ((((select 1)) = 1))
< ok 0 0
> exec
| select if((((select count(0) from n1) = 0)), 'NO', 'YES')
< ok 0 0
> exec
| select ((select 1) + 0)
< ok 0 0
> exec
| select (((select 1)) + 0)
< ok 0 0
> exec
| select (((select 1)))
< ok 0 0
> exec
| select (((select 1) union (select 2)))
< error
| 1242 21000 Subquery returns more than 1 row
> exec
| select 1 in (((select 1)))
< ok 0 0
> exec
| select 1 in (((select 1) union (select 2)))
< ok 0 0
> exec
| SELECT `d`.`k1` AS `k1` FROM ((SELECT `k1` AS `k1` FROM `test`.`j1`) AS `d` JOIN `test`.`j2` ON `d`.`k1`=`j2`.`k1`) JOIN `test`.`j3` ON `j3`.`k1`=`d`.`k1`
< error
| 1146 42S02 Table 'test.j1' doesn't exist
> exec
| SELECT `d`.`k1` AS `k1` FROM `test`.`j2` JOIN ((SELECT `k1` AS `k1` FROM `test`.`j1`) AS `d` JOIN `test`.`j3` USING (`k1`)) ON `j2`.`k1`=`d`.`k1`
< error
| 1146 42S02 Table 'test.j2' doesn't exist
> exec
| SELECT `d`.`k1` AS `k1` FROM ((SELECT `k1` AS `k1` FROM `test`.`j1`) UNION (SELECT `k1` AS `k1` FROM `test`.`j2`)) AS `d` JOIN `test`.`j3` ON `d`.`k1`=`j3`.`k1`
< error
| 1146 42S02 Table 'test.j1' doesn't exist
> exec
| select (((select 1) + 0)
< error
| 1064 42000 …line 1 column 24 near "" 
> exec
| select (((select 1) +))
< error
| 1064 42000 …line 1 column 22 near "))" 
> exec
| select (((select 1) d))
< error
| 1064 42000 …line 1 column 21 near "d))" 
> exec
| select (((select 1) union (select 2) + 0))
< error
| 1064 42000 …line 1 column 38 near "+ 0))" 
> exec
| CREATE OR REPLACE VIEW v1 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d JOIN j2 ON d.k1 = j2.k1 JOIN j3 ON j3.k1 = d.k1
< ok 0 0
> query
| SHOW CREATE VIEW v1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v1` (`k1`) AS SELECT `d`.`k1` AS `k1` FROM ((SELECT `k1` AS `k1` FROM `omni_paren_operand`.`j1`) AS `d` JOIN `omni_paren_operand`.`j2` ON `d`.`k1`=`j2`.`k1`) JOIN `omni_paren_operand`.`j3` ON `j3`.`k1`=`d`.`k1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v2 AS SELECT d.k1 FROM j2 JOIN ((SELECT k1 FROM j1) d JOIN j3 USING (k1)) ON j2.k1 = d.k1
< ok 0 0
> query
| SHOW CREATE VIEW v2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v2` (`k1`) AS SELECT `d`.`k1` AS `k1` FROM `omni_paren_operand`.`j2` JOIN ((SELECT `k1` AS `k1` FROM `omni_paren_operand`.`j1`) AS `d` JOIN `omni_paren_operand`.`j3` USING (`k1`)) ON `j2`.`k1`=`d`.`k1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v3 AS SELECT d.k1 FROM ((SELECT k1 FROM j1) UNION (SELECT k1 FROM j2)) d JOIN j3 ON d.k1 = j3.k1
< ok 0 0
> query
| SHOW CREATE VIEW v3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v3` (`k1`) AS SELECT `d`.`k1` AS `k1` FROM ((SELECT `k1` AS `k1` FROM `omni_paren_operand`.`j1`) UNION (SELECT `k1` AS `k1` FROM `omni_paren_operand`.`j2`)) AS `d` JOIN `omni_paren_operand`.`j3` ON `d`.`k1`=`j3`.`k1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v4 AS SELECT (((SELECT a FROM n1 LIMIT 1) + 0)) AS x, if((((SELECT count(0) FROM n1) = 0)), 'NO', 'YES') AS y
< ok 0 0
> query
| SHOW CREATE VIEW v4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v4` (`x`, `y`) AS SELECT (((SELECT `a` AS `a` FROM `omni_paren_operand`.`n1` LIMIT 1)+0)) AS `x`,IF((((SELECT COUNT(0) AS `count(0)` FROM `omni_paren_operand`.`n1`)=0)), _UTF8MB4'NO', _UTF8MB4'YES') AS `y`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v5 AS SELECT a FROM n1 WHERE a IN ((((SELECT a FROM n1 LIMIT 1) + 0)), 2)
< ok 0 0
> query
| SHOW CREATE VIEW v5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v5` (`a`) AS SELECT `a` AS `a` FROM `omni_paren_operand`.`n1` WHERE `a` IN ((((SELECT `a` AS `a` FROM `omni_paren_operand`.`n1` LIMIT 1)+0)),2)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS omni_paren_operand
< ok 0 0
