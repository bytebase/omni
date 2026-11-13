-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
| select `d`.`k1` AS `k1` from ((select `j1`.`k1` AS `k1` from `j1`) `d` join `j2` on((`d`.`k1` = `j2`.`k1`)))
< ok 0 0
> exec ^
< ok 0 0
> exec
| select d.k1 from ((select j1.k1 from j1) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) as d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d (x) join j2 on d.x = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) as d (x) join j2 on d.x = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d left join j2 on d.k1 = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d left outer join j2 on d.k1 = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d right join j2 on d.k1 = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d inner join j2 on d.k1 = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d cross join j2)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d join j2)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d join j2 using (k1))
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d natural join j2)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d natural left join j2)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d straight_join j2)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d straight_join j2 on d.k1 = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d, j2)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d, (select k1 from j2) e)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from (((select k1 from j1) d, j3) join j2 on 1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from (((select k1 from j1) d join j2 on d.k1 = j2.k1) join j3 on j3.k1 = d.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((((select k1 from j1) d join j2 on d.k1 = j2.k1) join j3 on j3.k1 = d.k1))
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d join (j2 join j3 on j2.k1 = j3.k1) on d.k1 = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) d join ((select k1 from j2) e join j3 on e.k1 = j3.k1) on d.k1 = e.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from (j1 join ((select 1) d join j2 on 1) on 1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from j1 join ((select 1) d join j2 on 1) on 1
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select 1) d join j2 on 1) join j3 on 1
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select 1) d join j2 on 1), j3
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from (((select k1 from j1)) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((((select k1 from j1))) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1 union select k1 from j2) d join j3 on d.k1 = j3.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from (((select k1 from j1) union (select k1 from j2)) d join j3 on d.k1 = j3.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from (((select k1 from j1) union (select k1 from j2) order by 1 limit 1) d join j3 on d.k1 = j3.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from (((select k1 from j1) except (select k1 from j2)) d join j3 on 1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from (((select k1 from j1) intersect (select k1 from j2)) d join j3 on 1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((with c as (select 1 k1) select k1 from c) d join j2 on 1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((table j1) d join j2 on 1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((values row(1)) d join j2 on 1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from (lateral (select 1) d join j2 on 1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from (j1 join lateral (select j1.k1) d on 1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from (j1, lateral (select j1.k1) d)
< ok 0 0
> exec ^
< ok 0 0
> exec
| update ((select k1 from j1) d join j2 on d.k1 = j2.k1) set j2.b = 1
< ok 0 0
> exec ^
< ok 0 0
> exec
| delete j2 from ((select k1 from j1) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| insert into j3 select d.k1, 1 from ((select k1 from j1) d join j2 on d.k1 = j2.k1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| with c as (select 1 k1) select * from ((select k1 from c) d join j2 on 1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select (select 1 from ((select 1) d join j2 on 1) limit 1)
< ok 0 0
> exec ^
< ok 0 0
> exec
| select 1 from j1 where exists (select 1 from ((select 1) d join j2 on 1))
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1)) d
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from (((select k1 from j1))) d
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) union (select k1 from j2)) d
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) union select k1 from j2) d
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) order by 1 limit 1) d
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select 1) for update) d
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select 1) lock in share mode) d
< ok 0 0
> exec ^
< ok 0 0
> exec
| select * from ((select k1 from j1) join j2 on 1)
< error
| 1248 42000 Every derived table must have its own alias
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| select * from ((select k1 from j1) join j2 on 1)
< error
| 1248 42000 Every derived table must have its own alias
> exec ^
: string
| select * from ((select k1 from j1) join j2 on 1)
| ; SELECT FROM
< error
| 1248 42000 Every derived table must have its own alias
> exec
| select * from ((select k1 from j1) join j2 on 1)
< error
| 1248 42000 Every derived table must have its own alias
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| select * from ((select k1 from j1) join j2 on 1)
< error
| 1248 42000 Every derived table must have its own alias
> exec ^
: string
| select * from ((select k1 from j1) join j2 on 1)
| ; SELECT FROM
< error
| 1248 42000 Every derived table must have its own alias
> exec
| select * from ((select k1 from j1) d)
< error
| 1064 42000 …near ')' at line 1
> exec ^
< error
| 1064 42000 …near ')' at line 1
> exec
| select * from (((select k1 from j1) d))
< error
| 1064 42000 …near '))' at line 1
> exec ^
< error
| 1064 42000 …near '))' at line 1
> exec
| select * from ((select k1 from j1) as d)
< error
| 1064 42000 …near ')' at line 1
> exec ^
< error
| 1064 42000 …near ')' at line 1
> exec
| select * from ((select k1 from j1) d) join j2 on 1
< error
| 1064 42000 …near ') join j2 on 1' at line 1
> exec ^
< error
| 1064 42000 …near ') join j2 on 1' at line 1
> exec
| select * from (((select 1) d) join j2 on 1)
< error
| 1064 42000 …near ') join j2 on 1)' at line 1
> exec ^
< error
| 1064 42000 …near ') join j2 on 1)' at line 1
> exec
| select * from (((select 1) d), j2)
< error
| 1064 42000 …near '), j2)' at line 1
> exec ^
< error
| 1064 42000 …near '), j2)' at line 1
> exec
| select * from j1, (lateral (select j1.k1) d)
< error
| 1064 42000 …near ')' at line 1
> exec ^
< error
| 1064 42000 …near ')' at line 1
> exec
| select * from (json_table('[1]', '$[*]' columns (x int path '$')) jt)
< error
| 1064 42000 …near ')' at line 1
> exec ^
< error
| 1064 42000 …near ')' at line 1
> exec
| select * from ((select k1 from j1) d join j2 on 1
< error
| 1064 42000 …near '' at line 1
> exec ^
< error
| 1064 42000 …near '' at line 1
> exec
| select * from ((select k1 from j1) d join)
< error
| 1064 42000 …near ')' at line 1
> exec ^
< error
| 1064 42000 …near ')' at line 1
> exec
| select * from ((select k1 from j1) d d2 join j2 on 1)
< error
| 1064 42000 …near 'd2 join j2 on 1)' at line 1
> exec ^
< error
| 1064 42000 …near 'd2 join j2 on 1)' at line 1
> exec
| select * from ((select k1 from j1) d join j2 on 1) d2
< error
| 1064 42000 …near 'd2' at line 1
> exec ^
< error
| 1064 42000 …near 'd2' at line 1
> exec
| select * from ((select k1 from j1) d (x join j2 on 1)
< error
| 1064 42000 …near 'join j2 on 1)' at line 1
> exec ^
< error
| 1064 42000 …near 'join j2 on 1)' at line 1
> exec
| select * from ((select k1 from j1) d use index (i) join j2 on 1)
< error
| 1064 42000 …near 'use index (i) join j2 on 1)' at line 1
> exec ^
< error
| 1064 42000 …near 'use index (i) join j2 on 1)' at line 1
> exec
| select * from ((select 1) union (select 2) d join j2 on 1)
< error
| 1064 42000 …near 'd join j2 on 1)' at line 1
> exec ^
< error
| 1064 42000 …near 'd join j2 on 1)' at line 1
> exec
| select * from ((select 1) order by 1 d join j2 on 1)
< error
| 1064 42000 …near 'd join j2 on 1)' at line 1
> exec ^
< error
| 1064 42000 …near 'd join j2 on 1)' at line 1
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
| select `d`.`k1` AS `k1`,`omni_t_mysql_parser`.`j2`.`b` AS `b` from ((select `omni_t_mysql_parser`.`j1`.`k1` AS `k1` from `omni_t_mysql_parser`.`j1`) `d` join `omni_t_mysql_parser`.`j2` on((`d`.`k1` = `omni_t_mysql_parser`.`j2`.`k1`)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `v1` AS select `d`.`k1` AS `k1`,`j2`.`b` AS `b` from ((select `j1`.`k1` AS `k1` from `j1`) `d` join `j2` on((`d`.`k1` = `j2`.`k1`)))
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
| select `d`.`k1` AS `k1` from (((select `omni_t_mysql_parser`.`j1`.`k1` AS `k1` from `omni_t_mysql_parser`.`j1`) `d` join `omni_t_mysql_parser`.`j2` on((`d`.`k1` = `omni_t_mysql_parser`.`j2`.`k1`))) join `omni_t_mysql_parser`.`j3` on((`omni_t_mysql_parser`.`j3`.`k1` = `d`.`k1`)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `v2` AS select `d`.`k1` AS `k1` from (((select `j1`.`k1` AS `k1` from `j1`) `d` join `j2` on((`d`.`k1` = `j2`.`k1`))) join `j3` on((`j3`.`k1` = `d`.`k1`)))
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
| select `d`.`k1` AS `k1` from ((select `omni_t_mysql_parser`.`j1`.`k1` AS `k1` from `omni_t_mysql_parser`.`j1` union select `omni_t_mysql_parser`.`j2`.`k1` AS `k1` from `omni_t_mysql_parser`.`j2`) `d` left join `omni_t_mysql_parser`.`j3` on((`omni_t_mysql_parser`.`j3`.`k1` = `d`.`k1`)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `v3` AS select `d`.`k1` AS `k1` from ((select `j1`.`k1` AS `k1` from `j1` union select `j2`.`k1` AS `k1` from `j2`) `d` left join `j3` on((`j3`.`k1` = `d`.`k1`)))
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
| select `d`.`x` AS `x` from ((select `omni_t_mysql_parser`.`j1`.`k1` AS `k1` from `omni_t_mysql_parser`.`j1`) `d` (`x`) join `omni_t_mysql_parser`.`j2` on((`d`.`x` = `omni_t_mysql_parser`.`j2`.`k1`)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `v4` AS select `d`.`x` AS `x` from ((select `j1`.`k1` AS `k1` from `j1`) `d` (`x`) join `j2` on((`d`.`x` = `j2`.`k1`)))
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
| select `omni_t_mysql_parser`.`j1`.`k1` AS `k1`,`d`.`y` AS `y` from `omni_t_mysql_parser`.`j1` join lateral (select (`omni_t_mysql_parser`.`j1`.`k1` + 1) AS `y`) `d`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `v5` AS select `j1`.`k1` AS `k1`,`d`.`y` AS `y` from (`j1` join lateral (select (`j1`.`k1` + 1) AS `y`) `d`)
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
| select `d`.`k1` AS `k1` from (select `omni_t_mysql_parser`.`j1`.`k1` AS `k1` from `omni_t_mysql_parser`.`j1`) `d` join `omni_t_mysql_parser`.`j2`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `v6` AS select `d`.`k1` AS `k1` from ((select `j1`.`k1` AS `k1` from `j1`) `d` join `j2`)
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
| select ((select `omni_t_mysql_parser`.`n1`.`a` from `omni_t_mysql_parser`.`n1` limit 1) + 0) AS `x`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `v7` AS select ((select `n1`.`a` from `n1` limit 1) + 0) AS `x`
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
| select `d`.`k1` AS `k1` from (`omni_t_mysql_parser`.`j2` join ((select `omni_t_mysql_parser`.`j1`.`k1` AS `k1` from `omni_t_mysql_parser`.`j1`) `d` join `omni_t_mysql_parser`.`j3` on((`d`.`k1` = `omni_t_mysql_parser`.`j3`.`k1`))) on((`omni_t_mysql_parser`.`j2`.`k1` = `d`.`k1`)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `v8` AS select `d`.`k1` AS `k1` from (`j2` join ((select `j1`.`k1` AS `k1` from `j1`) `d` join `j3` on((`d`.`k1` = `j3`.`k1`))) on((`j2`.`k1` = `d`.`k1`)))
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
| select `d`.`k1` AS `k1` from ((select `omni_t_mysql_parser`.`j1`.`k1` AS `k1` from `omni_t_mysql_parser`.`j1`) `d` join (select `omni_t_mysql_parser`.`j2`.`k1` AS `k1` from `omni_t_mysql_parser`.`j2`) `e` on((`d`.`k1` = `e`.`k1`)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `v9` AS select `d`.`k1` AS `k1` from ((select `j1`.`k1` AS `k1` from `j1`) `d` join (select `j2`.`k1` AS `k1` from `j2`) `e` on((`d`.`k1` = `e`.`k1`)))
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
| select `omni_t_mysql_parser`.`n1`.`a` AS `a` from `omni_t_mysql_parser`.`n1` where (((select `omni_t_mysql_parser`.`n1`.`a` from `omni_t_mysql_parser`.`n1` limit 1) = 1) and (`omni_t_mysql_parser`.`n1`.`a` in (((select `omni_t_mysql_parser`.`n1`.`a` from `omni_t_mysql_parser`.`n1` limit 1) + 0),2)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `v10` AS select `n1`.`a` AS `a` from `n1` where (((select `n1`.`a` from `n1` limit 1) = 1) and (`n1`.`a` in (((select `n1`.`a` from `n1` limit 1) + 0),2)))
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
| select ((select `omni_t_mysql_parser`.`n1`.`a` from `omni_t_mysql_parser`.`n1` limit 1) is null) AS `x`,if(((select count(0) from `omni_t_mysql_parser`.`n1`) = 0),'NO','YES') AS `y`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `v11` AS select ((select `n1`.`a` from `n1` limit 1) is null) AS `x`,if(((select count(0) from `n1`) = 0),'NO','YES') AS `y`
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
| select `d`.`k1` AS `k1` from ((select `omni_t_mysql_parser`.`j1`.`k1` AS `k1` from `omni_t_mysql_parser`.`j1` union select `omni_t_mysql_parser`.`j2`.`k1` AS `k1` from `omni_t_mysql_parser`.`j2`) `d` join `omni_t_mysql_parser`.`j3` on((`d`.`k1` = `omni_t_mysql_parser`.`j3`.`k1`)))
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `v12` AS select `d`.`k1` AS `k1` from ((select `j1`.`k1` AS `k1` from `j1` union select `j2`.`k1` AS `k1` from `j2`) `d` join `j3` on((`d`.`k1` = `j3`.`k1`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
