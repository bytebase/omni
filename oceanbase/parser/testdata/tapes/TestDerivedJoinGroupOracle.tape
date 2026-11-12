-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDerivedJoinGroupOracle default#1
> exec
| CREATE TABLE j1 (k1 INT, a INT, KEY i (k1))
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
| CREATE VIEW v1 AS SELECT * FROM (SELECT k1 FROM j1) d JOIN j2 USING (k1)
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v1
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1`,`omni_t_oceanbase_parser`.`j2`.`b` AS `b` from ((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) `d` join `omni_t_oceanbase_parser`.`j2` on (`d`.`k1` = `omni_t_oceanbase_parser`.`j2`.`k1`))
> query
| SHOW CREATE VIEW v1
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE VIEW `v1` AS select `d`.`k1` AS `k1`,`omni_t_oceanbase_parser`.`j2`.`b` AS `b` from ((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) `d` join `omni_t_oceanbase_parser`.`j2` on (`d`.`k1` = `omni_t_oceanbase_parser`.`j2`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v2 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d JOIN j2 ON d.k1 = j2.k1 JOIN j3 ON j3.k1 = d.k1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v2
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1` from (((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) `d` join `omni_t_oceanbase_parser`.`j2` on (`d`.`k1` = `omni_t_oceanbase_parser`.`j2`.`k1`)) join `omni_t_oceanbase_parser`.`j3` on (`omni_t_oceanbase_parser`.`j3`.`k1` = `d`.`k1`))
> query
| SHOW CREATE VIEW v2
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2
: bytes
| CREATE VIEW `v2` AS select `d`.`k1` AS `k1` from (((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) `d` join `omni_t_oceanbase_parser`.`j2` on (`d`.`k1` = `omni_t_oceanbase_parser`.`j2`.`k1`)) join `omni_t_oceanbase_parser`.`j3` on (`omni_t_oceanbase_parser`.`j3`.`k1` = `d`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v3 AS SELECT d.k1 FROM (SELECT k1 FROM j1 UNION SELECT k1 FROM j2) d LEFT JOIN j3 ON j3.k1 = d.k1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v3
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1` from (((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) union (select `omni_t_oceanbase_parser`.`j2`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j2`)) `d` left join `omni_t_oceanbase_parser`.`j3` on (`omni_t_oceanbase_parser`.`j3`.`k1` = `d`.`k1`))
> query
| SHOW CREATE VIEW v3
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v3
: bytes
| CREATE VIEW `v3` AS select `d`.`k1` AS `k1` from (((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) union (select `omni_t_oceanbase_parser`.`j2`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j2`)) `d` left join `omni_t_oceanbase_parser`.`j3` on (`omni_t_oceanbase_parser`.`j3`.`k1` = `d`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v4 AS SELECT d.x FROM (SELECT k1 FROM j1) d (x) JOIN j2 ON d.x = j2.k1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v4
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`x` AS `x` from ((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `x` from `omni_t_oceanbase_parser`.`j1`) `d` join `omni_t_oceanbase_parser`.`j2` on (`d`.`x` = `omni_t_oceanbase_parser`.`j2`.`k1`))
> query
| SHOW CREATE VIEW v4
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v4
: bytes
| CREATE VIEW `v4` AS select `d`.`x` AS `x` from ((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `x` from `omni_t_oceanbase_parser`.`j1`) `d` join `omni_t_oceanbase_parser`.`j2` on (`d`.`x` = `omni_t_oceanbase_parser`.`j2`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v5 AS SELECT j1.k1, d.y FROM j1, LATERAL (SELECT j1.k1 + 1 AS y) d
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v5
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1`,`d`.`y` AS `y` from `omni_t_oceanbase_parser`.`j1`,lateral (select (`omni_t_oceanbase_parser`.`j1`.`k1` + 1) AS `y`) `d`
> query
| SHOW CREATE VIEW v5
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v5
: bytes
| CREATE VIEW `v5` AS select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1`,`d`.`y` AS `y` from `omni_t_oceanbase_parser`.`j1`,lateral (select (`omni_t_oceanbase_parser`.`j1`.`k1` + 1) AS `y`) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v6 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d, j2
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v6
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1` from (select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) `d`,`omni_t_oceanbase_parser`.`j2`
> query
| SHOW CREATE VIEW v6
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v6
: bytes
| CREATE VIEW `v6` AS select `d`.`k1` AS `k1` from (select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) `d`,`omni_t_oceanbase_parser`.`j2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v7 AS SELECT (SELECT a FROM n1 LIMIT 1) + 0 AS x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v7
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select ((select `omni_t_oceanbase_parser`.`n1`.`a` from `omni_t_oceanbase_parser`.`n1` limit 1) + 0) AS `x`
> query
| SHOW CREATE VIEW v7
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v7
: bytes
| CREATE VIEW `v7` AS select ((select `omni_t_oceanbase_parser`.`n1`.`a` from `omni_t_oceanbase_parser`.`n1` limit 1) + 0) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v8 AS SELECT d.k1 FROM j2 JOIN ((SELECT k1 FROM j1) d JOIN j3 USING (k1)) ON j2.k1 = d.k1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v8
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1` from (`omni_t_oceanbase_parser`.`j2` join ((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) `d` join `omni_t_oceanbase_parser`.`j3` on (`d`.`k1` = `omni_t_oceanbase_parser`.`j3`.`k1`)) on (`omni_t_oceanbase_parser`.`j2`.`k1` = `d`.`k1`))
> query
| SHOW CREATE VIEW v8
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v8
: bytes
| CREATE VIEW `v8` AS select `d`.`k1` AS `k1` from (`omni_t_oceanbase_parser`.`j2` join ((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) `d` join `omni_t_oceanbase_parser`.`j3` on (`d`.`k1` = `omni_t_oceanbase_parser`.`j3`.`k1`)) on (`omni_t_oceanbase_parser`.`j2`.`k1` = `d`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v9 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d JOIN (SELECT k1 FROM j2) e ON d.k1 = e.k1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v9
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1` from ((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) `d` join (select `omni_t_oceanbase_parser`.`j2`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j2`) `e` on (`d`.`k1` = `e`.`k1`))
> query
| SHOW CREATE VIEW v9
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v9
: bytes
| CREATE VIEW `v9` AS select `d`.`k1` AS `k1` from ((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) `d` join (select `omni_t_oceanbase_parser`.`j2`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j2`) `e` on (`d`.`k1` = `e`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v10 AS SELECT a FROM n1 WHERE (SELECT a FROM n1 LIMIT 1) = 1 AND a IN ((SELECT a FROM n1 LIMIT 1) + 0, 2)
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v10
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `omni_t_oceanbase_parser`.`n1`.`a` AS `a` from `omni_t_oceanbase_parser`.`n1` where ((select `omni_t_oceanbase_parser`.`n1`.`a` from `omni_t_oceanbase_parser`.`n1` limit 1) = 1) and (`omni_t_oceanbase_parser`.`n1`.`a` in (((select `omni_t_oceanbase_parser`.`n1`.`a` from `omni_t_oceanbase_parser`.`n1` limit 1) + 0),2))
> query
| SHOW CREATE VIEW v10
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v10
: bytes
| CREATE VIEW `v10` AS select `omni_t_oceanbase_parser`.`n1`.`a` AS `a` from `omni_t_oceanbase_parser`.`n1` where ((select `omni_t_oceanbase_parser`.`n1`.`a` from `omni_t_oceanbase_parser`.`n1` limit 1) = 1) and (`omni_t_oceanbase_parser`.`n1`.`a` in (((select `omni_t_oceanbase_parser`.`n1`.`a` from `omni_t_oceanbase_parser`.`n1` limit 1) + 0),2))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v11 AS SELECT ((SELECT a FROM n1 LIMIT 1) IS NULL) AS x, if((SELECT count(0) FROM n1) = 0, 'NO', 'YES') AS y
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v11
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select ((select `omni_t_oceanbase_parser`.`n1`.`a` from `omni_t_oceanbase_parser`.`n1` limit 1) is null) AS `x`,(case when ((select count(0) from `omni_t_oceanbase_parser`.`n1`) = 0) then 'NO' else 'YES' end) AS `y`
> query
| SHOW CREATE VIEW v11
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v11
: bytes
| CREATE VIEW `v11` AS select ((select `omni_t_oceanbase_parser`.`n1`.`a` from `omni_t_oceanbase_parser`.`n1` limit 1) is null) AS `x`,(case when ((select count(0) from `omni_t_oceanbase_parser`.`n1`) = 0) then 'NO' else 'YES' end) AS `y`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v12 AS SELECT d.k1 FROM ((SELECT k1 FROM j1) UNION (SELECT k1 FROM j2)) d JOIN j3 ON d.k1 = j3.k1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v12
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1` from (((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) union (select `omni_t_oceanbase_parser`.`j2`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j2`)) `d` join `omni_t_oceanbase_parser`.`j3` on (`d`.`k1` = `omni_t_oceanbase_parser`.`j3`.`k1`))
> query
| SHOW CREATE VIEW v12
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v12
: bytes
| CREATE VIEW `v12` AS select `d`.`k1` AS `k1` from (((select `omni_t_oceanbase_parser`.`j1`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j1`) union (select `omni_t_oceanbase_parser`.`j2`.`k1` AS `k1` from `omni_t_oceanbase_parser`.`j2`)) `d` join `omni_t_oceanbase_parser`.`j3` on (`d`.`k1` = `omni_t_oceanbase_parser`.`j3`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
