-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestParenSubqueryGroupOracle default#1
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
| select `d`.`k1` AS `k1` from ((select `j1`.`k1` AS `k1` from `j1`) `d` join `j2` on(`d`.`k1` = `j2`.`k1`))
< ok 0 0
> exec
| select d.k1 from ((select j1.k1 from j1) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) as d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d (x) join j2 on d.x = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d left join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d right join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d inner join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d cross join j2)
< ok 0 0
> exec
| select * from ((select k1 from j1) d join j2 using (k1))
< ok 0 0
> exec
| select * from ((select k1 from j1) d natural join j2)
< ok 0 0
> exec
| select * from ((select k1 from j1) d straight_join j2)
< ok 0 0
> exec
| select * from ((select k1 from j1) d, j2)
< ok 0 0
> exec
| select * from ((select k1 from j1) d, (select k1 from j2) e)
< ok 0 0
> exec
| select * from (((select k1 from j1) d join j2 on d.k1 = j2.k1) join j3 on j3.k1 = d.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d join ((select k1 from j2) e join j3 on e.k1 = j3.k1) on d.k1 = e.k1)
< ok 0 0
> exec
| select * from (((select k1 from j1)) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from (((select k1 from j1) union (select k1 from j2)) d join j3 on d.k1 = j3.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1 union select k1 from j2) d join j3 on d.k1 = j3.k1)
< ok 0 0
> exec
| select * from ((with c as (select 1 k1) select k1 from c) d join j2 on 1)
< ok 0 0
> exec
| select * from j1 join ((select 1) d join j2 on 1) on 1
< ok 0 0
> exec
| update ((select k1 from j1) d join j2 on d.k1 = j2.k1) set j2.b = 1
< ok 0 0
> exec
| delete j2 from ((select k1 from j1) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d)
< ok 0 0
> exec
| select * from (((select k1 from j1) d))
< ok 0 0
> exec
| select * from (((select k1 from j1) d) join j2 on 1)
< ok 0 0
> exec
| select * from (((select k1 from j1) d), j2)
< ok 0 0
> exec
| select * from (j2 join ((select 1) d) on 1)
< ok 0 0
> exec
| select * from ((select k1 from j1)) d
< ok 0 0
> exec
| select * from ((select k1 from j1) union (select k1 from j2)) d
< ok 0 0
> exec
| select * from ((select k1 from j1) limit 1) d
< ok 0 0
> exec
| select ((select `n1`.`a` from `n1` limit 1) + 0) AS `x`
< ok 0 0
> exec
| select ((select 1) + 0)
< ok 0 0
> exec
| select ((select a from n1 limit 1) + 0) from n1
< ok 0 0
> exec
| select ((select 1) = 1)
< ok 0 0
> exec
| select ((select 1) in (1))
< ok 0 0
> exec
| select ((select 1) not in (1, 2))
< ok 0 0
> exec
| select ((select 1) is null)
< ok 0 0
> exec
| select ((select 1) is not true)
< ok 0 0
> exec
| select ((select 1) between 0 and 2)
< ok 0 0
> exec
| select ((select 1) like 'a')
< ok 0 0
> exec
| select ((select 1) and 1)
< ok 0 0
> exec
| select ((select 1) <=> 1)
< ok 0 0
> exec
| select ((select 1) * 2 + 1)
< ok 0 0
> exec
| select ((select 1) - (select 2))
< ok 0 0
> exec
| select ((select 1) collate utf8mb4_bin)
< ok 0 0
> exec
| select ((select 1) = any (select 1))
< ok 0 0
> exec
| select ((select 1), 2) = (1, 2)
< ok 0 0
> exec
| select (((select 1) + 0))
< ok 0 0
> exec
| select (((select 1)) + 0)
< ok 0 0
> exec
| select ((((select 1) + 0) + 0) + 0)
< ok 0 0
> exec
| select (((select 1) limit 1) + 0)
< ok 0 0
> exec
| select ((with c as (select 1 x) select x from c) = 1)
< ok 0 0
> exec
| select abs(((select 1) + 0))
< ok 0 0
> exec
| select if(((select count(0) from n1) = 0), 'NO', 'YES')
< ok 0 0
> exec
| select coalesce(((select a from n1 limit 1) + 0), 1)
< ok 0 0
> exec
| select 1 in ((select 1) + 0, 2)
< ok 0 0
> exec
| select 1 in (((select 1) + 0), 2)
< ok 0 0
> exec
| select 1 in (((select 1) = 1))
< ok 0 0
> exec
| select * from n1 where a in ((select 1) + 0)
< ok 0 0
> exec
| select * from n1 where ((select a from n1 limit 1) = 1) and ((select 1) + 0 = 1)
< ok 0 0
> exec
| select * from n1 where (a, a) in (((select 1), 1))
< ok 0 0
> exec
| select ((select 1))
< ok 0 0
> exec
| select ((select 1)) + 0
< ok 0 0
> exec
| select ((select 1) union (select 1))
< ok 0 0
> exec
| select 1 in ((select 1) union (select 2))
< ok 0 0
> exec
| select ((select 1) union (select 2) limit 1)
< ok 0 0
> exec
| select * from ((select k1 from j1) d join j2 on 1
< error
| 1064 42000 …near '' at line 1
> exec
| select * from ((select k1 from j1) d join)
< error
| 1064 42000 …near ')' at line 1
> exec
| select * from ((select k1 from j1) d d2 join j2 on 1)
< error
| 1064 42000 …near 'd2 join j2 on 1)' at line 1
> exec
| select * from ((select k1 from j1) d join j2 on 1) d2
< error
| 1064 42000 …near 'd2' at line 1
> exec
| select * from ((select 1) union (select 2) d join j2 on 1)
< error
| 1064 42000 …near 'd join j2 on 1)' at line 1
> exec
| select * from ((select 1) d union (select 2))
< error
| 1064 42000 …near 'union (select 2))' at line 1
> exec
| select ((select 1) + 0
< error
| 1064 42000 …near '' at line 1
> exec
| select ((select 1) +)
< error
| 1064 42000 …near ')' at line 1
> exec
| select ((select 1) + 0 from n1)
< error
| 1064 42000 …near 'from n1)' at line 1
> exec
| select ((select 1) d)
< error
| 1064 42000 …near 'd)' at line 1
> exec
| select ((select 1) union (select 2) = 1)
< error
| 1064 42000 …near '= 1)' at line 1
> exec
| select (((select 1) + 0) union (select 1))
< error
| 1064 42000 …near 'union (select 1))' at line 1
> exec
| select exists ((select 1) = 1)
< error
| 1064 42000 …near '= 1)' at line 1
> exec
| CREATE VIEW v1 AS SELECT * FROM (SELECT k1 FROM j1) d JOIN j2 USING (k1)
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?
: string
| v1
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1`,`test`.`j2`.`b` AS `b` from ((select `test`.`j1`.`k1` AS `k1` from `test`.`j1`) `d` join `test`.`j2` on(`d`.`k1` = `test`.`j2`.`k1`))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v1` AS select `d`.`k1` AS `k1`,`j2`.`b` AS `b` from ((select `j1`.`k1` AS `k1` from `j1`) `d` join `j2` on(`d`.`k1` = `j2`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v2 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d JOIN j2 ON d.k1 = j2.k1 JOIN j3 ON j3.k1 = d.k1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?
: string
| v2
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1` from (((select `test`.`j1`.`k1` AS `k1` from `test`.`j1`) `d` join `test`.`j2` on(`d`.`k1` = `test`.`j2`.`k1`)) join `test`.`j3` on(`test`.`j3`.`k1` = `d`.`k1`))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v2` AS select `d`.`k1` AS `k1` from (((select `j1`.`k1` AS `k1` from `j1`) `d` join `j2` on(`d`.`k1` = `j2`.`k1`)) join `j3` on(`j3`.`k1` = `d`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v3 AS SELECT d.k1 FROM (SELECT k1 FROM j1 UNION SELECT k1 FROM j2) d LEFT JOIN j3 ON j3.k1 = d.k1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?
: string
| v3
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1` from ((select `test`.`j1`.`k1` AS `k1` from `test`.`j1` union select `test`.`j2`.`k1` AS `k1` from `test`.`j2`) `d` left join `test`.`j3` on(`test`.`j3`.`k1` = `d`.`k1`))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v3` AS select `d`.`k1` AS `k1` from ((select `j1`.`k1` AS `k1` from `j1` union select `j2`.`k1` AS `k1` from `j2`) `d` left join `j3` on(`j3`.`k1` = `d`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v4 AS SELECT d.x FROM (SELECT k1 FROM j1) d (x) JOIN j2 ON d.x = j2.k1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?
: string
| v4
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`x` AS `x` from ((select `test`.`j1`.`k1` AS `k1` from `test`.`j1`) `d`(`x`) join `test`.`j2` on(`d`.`x` = `test`.`j2`.`k1`))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v4` AS select `d`.`x` AS `x` from ((select `j1`.`k1` AS `x` from `j1`) `d`(`x`) join `j2` on(`d`.`x` = `j2`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v6 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d, j2
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?
: string
| v6
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1` from (select `test`.`j1`.`k1` AS `k1` from `test`.`j1`) `d` join `test`.`j2`
> query
| SHOW CREATE VIEW v6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v6` AS select `d`.`k1` AS `k1` from ((select `j1`.`k1` AS `k1` from `j1`) `d` join `j2`)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v7 AS SELECT ((SELECT a FROM n1 LIMIT 1) + 0) AS x FROM n1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?
: string
| v7
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select (select `test`.`n1`.`a` from `test`.`n1` limit 1) + 0 AS `x` from `test`.`n1`
> query
| SHOW CREATE VIEW v7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v7` AS select (select `n1`.`a` from `n1` limit 1) + 0 AS `x` from `n1`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v8 AS SELECT d.k1 FROM j2 JOIN ((SELECT k1 FROM j1) d JOIN j3 USING (k1)) ON j2.k1 = d.k1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?
: string
| v8
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1` from (`test`.`j2` join ((select `test`.`j1`.`k1` AS `k1` from `test`.`j1`) `d` join `test`.`j3` on(`d`.`k1` = `test`.`j3`.`k1`)) on(`test`.`j2`.`k1` = `d`.`k1`))
> query
| SHOW CREATE VIEW v8
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v8
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v8` AS select `d`.`k1` AS `k1` from (`j2` join ((select `j1`.`k1` AS `k1` from `j1`) `d` join `j3` on(`d`.`k1` = `j3`.`k1`)) on(`j2`.`k1` = `d`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v9 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d JOIN (SELECT k1 FROM j2) e ON d.k1 = e.k1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?
: string
| v9
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1` from ((select `test`.`j1`.`k1` AS `k1` from `test`.`j1`) `d` join (select `test`.`j2`.`k1` AS `k1` from `test`.`j2`) `e` on(`d`.`k1` = `e`.`k1`))
> query
| SHOW CREATE VIEW v9
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v9
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v9` AS select `d`.`k1` AS `k1` from ((select `j1`.`k1` AS `k1` from `j1`) `d` join (select `j2`.`k1` AS `k1` from `j2`) `e` on(`d`.`k1` = `e`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v10 AS SELECT a FROM n1 WHERE ((SELECT a FROM n1 LIMIT 1) = 1) AND a IN ((SELECT a FROM n1 LIMIT 1) + 0, 2)
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?
: string
| v10
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `test`.`n1`.`a` AS `a` from `test`.`n1` where (select `test`.`n1`.`a` from `test`.`n1` limit 1) = 1 and `test`.`n1`.`a` in ((select `test`.`n1`.`a` from `test`.`n1` limit 1) + 0,2)
> query
| SHOW CREATE VIEW v10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v10` AS select `n1`.`a` AS `a` from `n1` where (select `n1`.`a` from `n1` limit 1) = 1 and `n1`.`a` in ((select `n1`.`a` from `n1` limit 1) + 0,2)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v11 AS SELECT ((SELECT a FROM n1 LIMIT 1) IS NULL) AS x, if(((SELECT count(0) FROM n1) = 0), 'NO', 'YES') AS y
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?
: string
| v11
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select (select `test`.`n1`.`a` from `test`.`n1` limit 1) is null AS `x`,if((select count(0) from `test`.`n1`) = 0,'NO','YES') AS `y`
> query
| SHOW CREATE VIEW v11
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v11
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v11` AS select (select `n1`.`a` from `n1` limit 1) is null AS `x`,if((select count(0) from `n1`) = 0,'NO','YES') AS `y`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v12 AS SELECT d.k1 FROM ((SELECT k1 FROM j1) UNION (SELECT k1 FROM j2)) d JOIN j3 ON d.k1 = j3.k1
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?
: string
| v12
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `d`.`k1` AS `k1` from (((select `test`.`j1`.`k1` AS `k1` from `test`.`j1`) union (select `test`.`j2`.`k1` AS `k1` from `test`.`j2`)) `d` join `test`.`j3` on(`d`.`k1` = `test`.`j3`.`k1`))
> query
| SHOW CREATE VIEW v12
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v12
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v12` AS select `d`.`k1` AS `k1` from (((select `j1`.`k1` AS `k1` from `j1`) union (select `j2`.`k1` AS `k1` from `j2`)) `d` join `j3` on(`d`.`k1` = `j3`.`k1`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
