-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_UnionOrderByKeepsTheOldAlias default#1
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
| CREATE VIEW uv0 (x) AS SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW uv0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| uv0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `uv0` AS select `t`.`a` AS `x` from `t` union select `t`.`b` AS `b` from `t` order by `z`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM uv0
< error
| 1356 HY000 View 'testdb.uv0' references invalid table(s) or column(s) or function(s) or definer/invoker of view lack rights to use them
> exec
| CREATE VIEW uv1 (x) AS SELECT a AS z FROM t UNION ALL SELECT b FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW uv1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| uv1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `uv1` AS select `t`.`a` AS `x` from `t` union all select `t`.`b` AS `b` from `t` order by `z`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM uv1
< error
| 1356 HY000 View 'testdb.uv1' references invalid table(s) or column(s) or function(s) or definer/invoker of view lack rights to use them
> exec
| CREATE VIEW uv2 (x) AS SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY z DESC
< ok 0 0
> query
| SHOW CREATE VIEW uv2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| uv2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `uv2` AS select `t`.`a` AS `x` from `t` union select `t`.`b` AS `b` from `t` order by `z` desc
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM uv2
< error
| 1356 HY000 View 'testdb.uv2' references invalid table(s) or column(s) or function(s) or definer/invoker of view lack rights to use them
> exec
| CREATE VIEW uv3 (x) AS SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY z + 1
< ok 0 0
> query
| SHOW CREATE VIEW uv3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| uv3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `uv3` AS select `t`.`a` AS `x` from `t` union select `t`.`b` AS `b` from `t` order by `z` + 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM uv3
< error
| 1356 HY000 View 'testdb.uv3' references invalid table(s) or column(s) or function(s) or definer/invoker of view lack rights to use them
> exec
| CREATE VIEW uv4 (x, y) AS SELECT a AS z, b AS w FROM t UNION SELECT b, a FROM t ORDER BY w, z
< ok 0 0
> query
| SHOW CREATE VIEW uv4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| uv4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `uv4` AS select `t`.`a` AS `x`,`t`.`b` AS `y` from `t` union select `t`.`b` AS `b`,`t`.`a` AS `a` from `t` order by `w`,`z`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM uv4
< error
| 1356 HY000 View 'testdb.uv4' references invalid table(s) or column(s) or function(s) or definer/invoker of view lack rights to use them
> exec
| CREATE VIEW uv5 (x) AS SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY 1
< ok 0 0
> query
| SHOW CREATE VIEW uv5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| uv5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `uv5` AS select `t`.`a` AS `x` from `t` union select `t`.`b` AS `b` from `t` order by 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM uv5
< ok 0 0
> exec
| CREATE VIEW uv6  AS SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW uv6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| uv6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `uv6` AS select `t`.`a` AS `z` from `t` union select `t`.`b` AS `b` from `t` order by `z`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM uv6
< ok 0 0
> exec
| CREATE VIEW uv7 (x) AS SELECT a AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW uv7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| uv7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `uv7` AS select `t`.`a` AS `x` from `t` order by `t`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM uv7
< ok 0 0
> exec
| CREATE VIEW uv_new (x) AS SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY x
< error
| 1054 42S22 Unknown column 'x' in 'ORDER BY'
