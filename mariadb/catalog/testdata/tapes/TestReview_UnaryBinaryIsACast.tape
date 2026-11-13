-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_UnaryBinaryIsACast default#1
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
| CREATE TABLE t (a INT, b VARCHAR(10))
< ok 0 0
> exec
| INSERT INTO t VALUES (1,'x'),(2,'y')
< ok 2 0
> exec
| CREATE VIEW v_bn0 AS SELECT 2 * BINARY (1 + 2) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_bn0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bn0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn0` AS select 2 * cast(1 + 2 as char charset binary) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_bn0
< result
< column "x" "DOUBLE"
< row
: float64 6
> exec
| DROP VIEW v_bn0
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn0` AS select 2 * cast(1 + 2 as char charset binary) AS `x`
< ok 0 0
> query
| SELECT * FROM v_bn0
< result
< column "x" "DOUBLE"
< row
: float64 6
> exec
| CREATE VIEW v_bn1 AS SELECT BINARY (1 + 2) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_bn1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bn1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn1` AS select cast(1 + 2 as char charset binary) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_bn1
< result
< column "x" "VARBINARY"
< row
: bytes
| 3
> exec
| DROP VIEW v_bn1
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn1` AS select cast(1 + 2 as char charset binary) AS `x`
< ok 0 0
> query
| SELECT * FROM v_bn1
< result
< column "x" "VARBINARY"
< row
: bytes
| 3
> exec
| CREATE VIEW v_bn2 AS SELECT 2 * BINARY 3 AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_bn2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bn2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn2` AS select 2 * cast(3 as char charset binary) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_bn2
< result
< column "x" "DOUBLE"
< row
: float64 6
> exec
| DROP VIEW v_bn2
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn2` AS select 2 * cast(3 as char charset binary) AS `x`
< ok 0 0
> query
| SELECT * FROM v_bn2
< result
< column "x" "DOUBLE"
< row
: float64 6
> exec
| CREATE VIEW v_bn3 AS SELECT BINARY 'a' = 'A' AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_bn3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bn3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn3` AS select cast('a' as char charset binary) = 'A' AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_bn3
< result
< column "x" "INT"
< row
: int64 0
> exec
| DROP VIEW v_bn3
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn3` AS select cast('a' as char charset binary) = 'A' AS `x`
< ok 0 0
> query
| SELECT * FROM v_bn3
< result
< column "x" "INT"
< row
: int64 0
> exec
| CREATE VIEW v_bn4 AS SELECT -BINARY (1 + 2) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_bn4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bn4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn4` AS select -cast(1 + 2 as char charset binary) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_bn4
< result
< column "x" "DOUBLE"
< row
: float64 -3
> exec
| DROP VIEW v_bn4
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn4` AS select -cast(1 + 2 as char charset binary) AS `x`
< ok 0 0
> query
| SELECT * FROM v_bn4
< result
< column "x" "DOUBLE"
< row
: float64 -3
> exec
| CREATE VIEW v_bn5 AS SELECT CONCAT(BINARY (1+2), 'x') AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_bn5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bn5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn5` AS select concat(cast(1 + 2 as char charset binary),'x') AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_bn5
< result
< column "x" "VARBINARY"
< row
: bytes
| 3x
> exec
| DROP VIEW v_bn5
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn5` AS select concat(cast(1 + 2 as char charset binary),'x') AS `x`
< ok 0 0
> query
| SELECT * FROM v_bn5
< result
< column "x" "VARBINARY"
< row
: bytes
| 3x
> exec
| CREATE VIEW v_bn6 AS SELECT 2 * BINARY (SELECT 3) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_bn6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bn6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn6` AS select 2 * cast((select 3) as char charset binary) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_bn6
< result
< column "x" "DOUBLE"
< row
: float64 6
> exec
| DROP VIEW v_bn6
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn6` AS select 2 * cast((select 3) as char charset binary) AS `x`
< ok 0 0
> query
| SELECT * FROM v_bn6
< result
< column "x" "DOUBLE"
< row
: float64 6
> exec
| CREATE VIEW v_bn7 AS SELECT BINARY b AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bn7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bn7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn7` AS select cast(`t`.`b` as char charset binary) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_bn7
< result
< column "x" "VARBINARY"
< row
: bytes
| x
< row
: bytes
| y
> exec
| DROP VIEW v_bn7
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn7` AS select cast(`t`.`b` as char charset binary) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_bn7
< result
< column "x" "VARBINARY"
< row
: bytes
| x
< row
: bytes
| y
> exec
| CREATE VIEW v_bn8 AS SELECT 2 * CAST((1 + 2) AS BINARY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_bn8
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bn8
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn8` AS select 2 * cast(1 + 2 as char charset binary) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_bn8
< result
< column "x" "DOUBLE"
< row
: float64 6
> exec
| DROP VIEW v_bn8
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn8` AS select 2 * cast(1 + 2 as char charset binary) AS `x`
< ok 0 0
> query
| SELECT * FROM v_bn8
< result
< column "x" "DOUBLE"
< row
: float64 6
> exec
| CREATE VIEW v_bn9 AS SELECT CAST(b AS BINARY) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bn9
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bn9
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn9` AS select cast(`t`.`b` as char charset binary) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_bn9
< result
< column "x" "VARBINARY"
< row
: bytes
| x
< row
: bytes
| y
> exec
| DROP VIEW v_bn9
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bn9` AS select cast(`t`.`b` as char charset binary) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_bn9
< result
< column "x" "VARBINARY"
< row
: bytes
| x
< row
: bytes
| y
