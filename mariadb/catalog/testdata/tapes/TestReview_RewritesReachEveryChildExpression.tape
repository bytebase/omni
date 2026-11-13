-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_RewritesReachEveryChildExpression default#1
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
| CREATE VIEW v_ch0 AS SELECT ROW((2=1) SOUNDS LIKE '0',1) = ROW(1,1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_ch0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ch0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch0` AS select (soundex(2 = 1) = soundex('0'),1) = (1,1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_ch0 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 1
> exec
| DROP VIEW v_ch0
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch0` AS select (soundex(2 = 1) = soundex('0'),1) = (1,1) AS `x`
< ok 0 0
> query
| SELECT * FROM v_ch0 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 1
> exec
| CREATE VIEW v_ch1 AS SELECT ROW(NOT (a > 0),1) = ROW(1,1) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_ch1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ch1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch1` AS select (`t`.`a` <= 0,1) = (1,1) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_ch1 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 0
< row
: int64 0
> exec
| DROP VIEW v_ch1
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch1` AS select (`t`.`a` <= 0,1) = (1,1) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_ch1 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 0
< row
: int64 0
> exec
| CREATE VIEW v_ch2 AS SELECT (a,b) = ((a = 1) SOUNDS LIKE '0', 1) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_ch2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ch2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch2` AS select (`t`.`a`,`t`.`b`) = (soundex(`t`.`a` = 1) = soundex('0'),1) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_ch2 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 0
< row
: int64 0
> exec
| DROP VIEW v_ch2
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch2` AS select (`t`.`a`,`t`.`b`) = (soundex(`t`.`a` = 1) = soundex('0'),1) AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_ch2 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 0
< row
: int64 0
> exec
| CREATE VIEW v_ch3 AS SELECT '2020-01-01' + INTERVAL ((2=1) SOUNDS LIKE '0') DAY AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_ch3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ch3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch3` AS select '2020-01-01' + interval soundex(2 = 1) = soundex('0') day AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_ch3 ORDER BY 1
< result
< column "x" "VARCHAR"
< row
: bytes
| 2020-01-02
> exec
| DROP VIEW v_ch3
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch3` AS select '2020-01-01' + interval soundex(2 = 1) = soundex('0') day AS `x`
< ok 0 0
> query
| SELECT * FROM v_ch3 ORDER BY 1
< result
< column "x" "VARCHAR"
< row
: bytes
| 2020-01-02
> exec
| CREATE VIEW v_ch4 AS SELECT '2020-01-01' + INTERVAL (NOT (a > 0)) DAY AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_ch4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ch4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch4` AS select '2020-01-01' + interval `t`.`a` <= 0 day AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_ch4 ORDER BY 1
< result
< column "x" "CHAR"
< row
: bytes
| 2020-01-01
< row
: bytes
| 2020-01-01
> exec
| DROP VIEW v_ch4
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch4` AS select '2020-01-01' + interval `t`.`a` <= 0 day AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_ch4 ORDER BY 1
< result
< column "x" "CHAR"
< row
: bytes
| 2020-01-01
< row
: bytes
| 2020-01-01
> exec
| CREATE VIEW v_ch5 AS SELECT EXTRACT(DAY FROM ((2=1) SOUNDS LIKE '0')) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_ch5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ch5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch5` AS select extract(day from soundex(2 = 1) = soundex('0')) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_ch5 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 0
> exec
| DROP VIEW v_ch5
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch5` AS select extract(day from soundex(2 = 1) = soundex('0')) AS `x`
< ok 0 0
> query
| SELECT * FROM v_ch5 ORDER BY 1
< result
< column "x" "INT"
< row
: int64 0
> exec
| CREATE VIEW v_ch6 AS SELECT GROUP_CONCAT(a ORDER BY (a = 1) SOUNDS LIKE '0') AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_ch6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ch6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch6` AS select group_concat(`t`.`a` order by soundex(`t`.`a` = 1) = soundex('0') ASC separator ',') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_ch6 ORDER BY 1
< result
< column "x" "TEXT"
< row
: bytes
| 2,1
> exec
| DROP VIEW v_ch6
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch6` AS select group_concat(`t`.`a` order by soundex(`t`.`a` = 1) = soundex('0') ASC separator ',') AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_ch6 ORDER BY 1
< result
< column "x" "TEXT"
< row
: bytes
| 2,1
> exec
| CREATE VIEW v_ch7 AS SELECT GROUP_CONCAT(a ORDER BY NOT (a > 0)) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_ch7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ch7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch7` AS select group_concat(`t`.`a` order by `t`.`a` <= 0 ASC separator ',') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_ch7 ORDER BY 1
< result
< column "x" "TEXT"
< row
: bytes
| 2,1
> exec
| DROP VIEW v_ch7
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ch7` AS select group_concat(`t`.`a` order by `t`.`a` <= 0 ASC separator ',') AS `x` from `t`
< ok 0 0
> query
| SELECT * FROM v_ch7 ORDER BY 1
< result
< column "x" "TEXT"
< row
: bytes
| 2,1
