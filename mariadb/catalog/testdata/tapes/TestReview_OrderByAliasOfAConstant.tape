-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_OrderByAliasOfAConstant default#1
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
| CREATE VIEW oa0 AS SELECT 2 AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa0` AS select 2 AS `z` from `t` order by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa0
< ok 0 0
> exec
| CREATE VIEW oa1 AS SELECT 0 AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa1` AS select 0 AS `z` from `t` order by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa1
< ok 0 0
> exec
| CREATE VIEW oa2 AS SELECT -2 AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa2` AS select -2 AS `z` from `t` order by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa2
< ok 0 0
> exec
| CREATE VIEW oa3 AS SELECT +2 AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa3` AS select 2 AS `z` from `t` order by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa3
< ok 0 0
> exec
| CREATE VIEW oa4 AS SELECT (2) AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa4` AS select 2 AS `z` from `t` order by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa4
< ok 0 0
> exec
| CREATE VIEW oa5 AS SELECT TRUE AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa5` AS select 1 AS `z` from `t` order by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa5
< ok 0 0
> exec
| CREATE VIEW oa6 AS SELECT 2 AS z, a FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa6` AS select 2 AS `z`,`t`.`a` AS `a` from `t` order by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa6
< ok 0 0
> exec
| CREATE VIEW oa7 AS SELECT a, 2 AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa7` AS select `t`.`a` AS `a`,2 AS `z` from `t` order by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa7
< ok 0 0
> exec
| CREATE VIEW oa8 AS SELECT 2 AS z FROM t ORDER BY z DESC
< ok 0 0
> query
| SHOW CREATE VIEW oa8
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa8
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa8` AS select 2 AS `z` from `t` order by '' desc
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa8
< ok 0 0
> exec
| CREATE VIEW oa9 AS SELECT 2 AS z FROM t GROUP BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa9
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa9
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa9` AS select 2 AS `z` from `t` group by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa9
< ok 0 0
> exec
| CREATE VIEW oa10 AS SELECT 2 AS z FROM t GROUP BY z ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa10` AS select 2 AS `z` from `t` group by '' order by ''
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa10
< ok 0 0
> exec
| CREATE VIEW oa11 AS SELECT 2.5 AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa11
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa11
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa11` AS select 2.5 AS `z` from `t` order by 2.5
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa11
< ok 0 0
> exec
| CREATE VIEW oa12 AS SELECT 1e2 AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa12
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa12
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa12` AS select 1e2 AS `z` from `t` order by 1e2
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa12
< ok 0 0
> exec
| CREATE VIEW oa13 AS SELECT 'x' AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa13
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa13
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa13` AS select 'x' AS `z` from `t` order by 'x'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa13
< ok 0 0
> exec
| CREATE VIEW oa14 AS SELECT NULL AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa14
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa14
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa14` AS select NULL AS `z` from `t` order by NULL
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa14
< ok 0 0
> exec
| CREATE VIEW oa15 AS SELECT a AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa15
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa15
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa15` AS select `t`.`a` AS `z` from `t` order by `t`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa15
< ok 0 0
> exec
| CREATE VIEW oa16 AS SELECT a AS z FROM t GROUP BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa16
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa16
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa16` AS select `t`.`a` AS `z` from `t` group by `t`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa16
< ok 0 0
> exec
| CREATE VIEW oa17 AS SELECT 1 + 1 AS z FROM t ORDER BY z
< ok 0 0
> query
| SHOW CREATE VIEW oa17
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa17
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa17` AS select 1 + 1 AS `z` from `t` order by 1 + 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa17
< ok 0 0
> exec
| CREATE VIEW oa18 AS SELECT 2 AS z FROM t ORDER BY z + 0
< ok 0 0
> query
| SHOW CREATE VIEW oa18
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa18
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa18` AS select 2 AS `z` from `t` order by `z` + 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa18
< ok 0 0
> exec
| CREATE VIEW oa19 AS SELECT 2 AS z FROM t ORDER BY 1
< ok 0 0
> query
| SHOW CREATE VIEW oa19
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa19
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa19` AS select 2 AS `z` from `t` order by 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa19
< ok 0 0
> exec
| CREATE VIEW oa20 AS SELECT a AS z, b FROM t ORDER BY z, b
< ok 0 0
> query
| SHOW CREATE VIEW oa20
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| oa20
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `oa20` AS select `t`.`a` AS `z`,`t`.`b` AS `b` from `t` order by `t`.`a`,`t`.`b`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| SELECT * FROM oa20
< ok 0 0
