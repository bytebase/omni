-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_GroupByPrefersTheSourceColumn default#1
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
| INSERT INTO t VALUES (1,10),(1,20)
< ok 2 0
> exec
| CREATE VIEW gc0 AS SELECT a AS b, COUNT(*) AS n FROM t GROUP BY b
< ok 0 0
> query
| SHOW CREATE VIEW gc0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| gc0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `gc0` AS select `t`.`a` AS `b`,count(0) AS `n` from `t` group by `t`.`b`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM gc0 ORDER BY 1
< result
< column "b" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 1
: int64 1
> query
| SELECT * FROM (select `t`.`a` AS `b`,count(0) AS `n` from `t` group by `t`.`b`) AS q ORDER BY 1
< result
< column "b" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 1
: int64 1
> exec
| CREATE VIEW gc1 AS SELECT a AS b, COUNT(*) AS n FROM t GROUP BY b ORDER BY b
< ok 0 0
> query
| SHOW CREATE VIEW gc1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| gc1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `gc1` AS select `t`.`a` AS `b`,count(0) AS `n` from `t` group by `t`.`b` order by `t`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM gc1 ORDER BY 1
< result
< column "b" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 1
: int64 1
> query
| SELECT * FROM (select `t`.`a` AS `b`,count(0) AS `n` from `t` group by `t`.`b` order by `t`.`a`) AS q ORDER BY 1
< result
< column "b" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 1
: int64 1
> exec
| CREATE VIEW gc2 AS SELECT 2 AS b, COUNT(*) AS n FROM t GROUP BY b
< ok 0 0
> query
| SHOW CREATE VIEW gc2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| gc2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `gc2` AS select 2 AS `b`,count(0) AS `n` from `t` group by `t`.`b`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM gc2 ORDER BY 1
< result
< column "b" "INT"
< column "n" "BIGINT"
< row
: int64 2
: int64 1
< row
: int64 2
: int64 1
> query
| SELECT * FROM (select 2 AS `b`,count(0) AS `n` from `t` group by `t`.`b`) AS q ORDER BY 1
< result
< column "b" "INT"
< column "n" "BIGINT"
< row
: int64 2
: int64 1
< row
: int64 2
: int64 1
> exec
| CREATE VIEW gc3 AS SELECT a + 1 AS b, COUNT(*) AS n FROM t GROUP BY b
< ok 0 0
> query
| SHOW CREATE VIEW gc3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| gc3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `gc3` AS select `t`.`a` + 1 AS `b`,count(0) AS `n` from `t` group by `t`.`b`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM gc3 ORDER BY 1
< result
< column "b" "BIGINT"
< column "n" "BIGINT"
< row
: int64 2
: int64 1
< row
: int64 2
: int64 1
> query
| SELECT * FROM (select `t`.`a` + 1 AS `b`,count(0) AS `n` from `t` group by `t`.`b`) AS q ORDER BY 1
< result
< column "b" "BIGINT"
< column "n" "BIGINT"
< row
: int64 2
: int64 1
< row
: int64 2
: int64 1
> exec
| CREATE VIEW gc4 AS SELECT a AS b, COUNT(*) AS n FROM t GROUP BY t.b
< ok 0 0
> query
| SHOW CREATE VIEW gc4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| gc4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `gc4` AS select `t`.`a` AS `b`,count(0) AS `n` from `t` group by `t`.`b`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM gc4 ORDER BY 1
< result
< column "b" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 1
: int64 1
> query
| SELECT * FROM (select `t`.`a` AS `b`,count(0) AS `n` from `t` group by `t`.`b`) AS q ORDER BY 1
< result
< column "b" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 1
< row
: int64 1
: int64 1
> exec
| CREATE VIEW gc5 AS SELECT a AS z, COUNT(*) AS n FROM t GROUP BY z
< ok 0 0
> query
| SHOW CREATE VIEW gc5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| gc5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `gc5` AS select `t`.`a` AS `z`,count(0) AS `n` from `t` group by `t`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM gc5 ORDER BY 1
< result
< column "z" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 2
> query
| SELECT * FROM (select `t`.`a` AS `z`,count(0) AS `n` from `t` group by `t`.`a`) AS q ORDER BY 1
< result
< column "z" "INT"
< column "n" "BIGINT"
< row
: int64 1
: int64 2
> exec
| CREATE VIEW gc6 AS SELECT a + 1 AS z, COUNT(*) AS n FROM t GROUP BY z
< ok 0 0
> query
| SHOW CREATE VIEW gc6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| gc6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `gc6` AS select `t`.`a` + 1 AS `z`,count(0) AS `n` from `t` group by `t`.`a` + 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM gc6 ORDER BY 1
< result
< column "z" "BIGINT"
< column "n" "BIGINT"
< row
: int64 2
: int64 2
> query
| SELECT * FROM (select `t`.`a` + 1 AS `z`,count(0) AS `n` from `t` group by `t`.`a` + 1) AS q ORDER BY 1
< result
< column "z" "BIGINT"
< column "n" "BIGINT"
< row
: int64 2
: int64 2
> exec
| CREATE VIEW gc7 AS SELECT 2 AS z, COUNT(*) AS n FROM t GROUP BY z
< ok 0 0
> query
| SHOW CREATE VIEW gc7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| gc7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `gc7` AS select 2 AS `z`,count(0) AS `n` from `t` group by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM gc7 ORDER BY 1
< result
< column "z" "INT"
< column "n" "BIGINT"
< row
: int64 2
: int64 2
> query
| SELECT * FROM (select 2 AS `z`,count(0) AS `n` from `t` group by '') AS q ORDER BY 1
< result
< column "z" "INT"
< column "n" "BIGINT"
< row
: int64 2
: int64 2
> exec
| CREATE VIEW gc8 AS SELECT a AS b FROM t ORDER BY b
< ok 0 0
> query
| SHOW CREATE VIEW gc8
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| gc8
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `gc8` AS select `t`.`a` AS `b` from `t` order by `t`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM gc8 ORDER BY 1
< result
< column "b" "INT"
< row
: int64 1
< row
: int64 1
> query
| SELECT * FROM (select `t`.`a` AS `b` from `t` order by `t`.`a`) AS q ORDER BY 1
< result
< column "b" "INT"
< row
: int64 1
< row
: int64 1
> exec
| CREATE VIEW gc9 AS SELECT 2 AS b FROM t ORDER BY b
< ok 0 0
> query
| SHOW CREATE VIEW gc9
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| gc9
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `gc9` AS select 2 AS `b` from `t` order by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM gc9 ORDER BY 1
< result
< column "b" "INT"
< row
: int64 2
< row
: int64 2
> query
| SELECT * FROM (select 2 AS `b` from `t` order by '') AS q ORDER BY 1
< result
< column "b" "INT"
< row
: int64 2
< row
: int64 2
