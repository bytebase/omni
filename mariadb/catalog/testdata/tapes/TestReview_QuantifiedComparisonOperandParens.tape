-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_QuantifiedComparisonOperandParens default#1
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
| CREATE VIEW qc0 AS SELECT (1 OR 0) = ANY (SELECT 0) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc0` AS select (1 or 0) in (select 0) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc0
< result
< column "x" "INT"
< row
: int64 0
> query
| SELECT * FROM (select (1 or 0) in (select 0) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 0
> exec
| CREATE VIEW qc1 AS SELECT (1 AND 1) = ANY (SELECT 1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc1` AS select (1 and 1) in (select 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc1
< result
< column "x" "INT"
< row
: int64 1
> query
| SELECT * FROM (select (1 and 1) in (select 1) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 1
> exec
| CREATE VIEW qc2 AS SELECT (1 XOR 0) = ANY (SELECT 1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc2` AS select (1 xor 0) in (select 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc2
< result
< column "x" "INT"
< row
: int64 1
> query
| SELECT * FROM (select (1 xor 0) in (select 1) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 1
> exec
| CREATE VIEW qc3 AS SELECT (1 = 1) = ANY (SELECT 1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc3` AS select (1 = 1) in (select 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc3
< result
< column "x" "INT"
< row
: int64 1
> query
| SELECT * FROM (select (1 = 1) in (select 1) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 1
> exec
| CREATE VIEW qc4 AS SELECT (1 BETWEEN 0 AND 2) = ANY (SELECT 1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc4` AS select (1 between 0 and 2) in (select 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc4
< result
< column "x" "INT"
< row
: int64 1
> query
| SELECT * FROM (select (1 between 0 and 2) in (select 1) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 1
> exec
| CREATE VIEW qc5 AS SELECT (1 OR 0) = SOME (SELECT 0) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc5` AS select (1 or 0) in (select 0) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc5
< result
< column "x" "INT"
< row
: int64 0
> query
| SELECT * FROM (select (1 or 0) in (select 0) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 0
> exec
| CREATE VIEW qc6 AS SELECT (1 OR 0) <> ALL (SELECT 1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc6` AS select !((1 or 0) in (select 1)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc6
< result
< column "x" "INT"
< row
: int64 0
> query
| SELECT * FROM (select !((1 or 0) in (select 1)) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 0
> exec
| CREATE VIEW qc7 AS SELECT 1 + 1 = ANY (SELECT 2) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc7` AS select 1 + 1 in (select 2) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc7
< result
< column "x" "INT"
< row
: int64 1
> query
| SELECT * FROM (select 1 + 1 in (select 2) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 1
> exec
| CREATE VIEW qc8 AS SELECT (NOT 1) = ANY (SELECT 0) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc8
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc8
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc8` AS select !1 in (select 0) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc8
< result
< column "x" "INT"
< row
: int64 1
> query
| SELECT * FROM (select !1 in (select 0) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 1
> exec
| CREATE VIEW qc9 AS SELECT (1 LIKE '1') = ANY (SELECT 1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc9
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc9
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc9` AS select 1 like '1' in (select 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc9
< result
< column "x" "INT"
< row
: int64 1
> query
| SELECT * FROM (select 1 like '1' in (select 1) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 1
> exec
| CREATE VIEW qc10 AS SELECT 1 = ANY (SELECT 1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc10` AS select 1 in (select 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc10
< result
< column "x" "INT"
< row
: int64 1
> query
| SELECT * FROM (select 1 in (select 1) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 1
> exec
| CREATE VIEW qc11 AS SELECT (1 OR 0) > ALL (SELECT 1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc11
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc11
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc11` AS select 1 or 0 > all (select 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc11
< result
< column "x" "INT"
< row
: int64 1
> query
| SELECT * FROM (select 1 or 0 > all (select 1) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 1
> exec
| CREATE VIEW qc12 AS SELECT (1 OR 0) < ANY (SELECT 1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW qc12
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| qc12
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `qc12` AS select 1 or 0 < any (select 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT x FROM qc12
< result
< column "x" "INT"
< row
: int64 1
> query
| SELECT * FROM (select 1 or 0 < any (select 1) AS `x`) AS q
< result
< column "x" "INT"
< row
: int64 1
