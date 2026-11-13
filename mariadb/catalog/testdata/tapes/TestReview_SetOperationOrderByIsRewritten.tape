-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_SetOperationOrderByIsRewritten default#1
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
| CREATE VIEW v_so0 AS SELECT 2 AS a UNION SELECT 3 ORDER BY ((a=2) SOUNDS LIKE '0'), a DESC
< ok 0 0
> query
| SHOW CREATE VIEW v_so0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_so0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_so0` AS select 2 AS `a` union select 3 AS `3` order by soundex(`a` = 2) = soundex('0'),`a` desc
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_so0 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 2
< row
: int64 3
> exec
| DROP VIEW v_so0
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_so0` AS select 2 AS `a` union select 3 AS `3` order by soundex(`a` = 2) = soundex('0'),`a` desc
< ok 0 0
> query
| SELECT * FROM v_so0 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 2
< row
: int64 3
> exec
| CREATE VIEW v_so1 AS SELECT 2 AS a UNION SELECT 3 ORDER BY NOT (a > 2)
< ok 0 0
> query
| SHOW CREATE VIEW v_so1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_so1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_so1` AS select 2 AS `a` union select 3 AS `3` order by `a` <= 2
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_so1 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 2
< row
: int64 3
> exec
| DROP VIEW v_so1
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_so1` AS select 2 AS `a` union select 3 AS `3` order by `a` <= 2
< ok 0 0
> query
| SELECT * FROM v_so1 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 2
< row
: int64 3
> exec
| CREATE VIEW v_so2 AS SELECT a FROM t UNION SELECT b FROM t ORDER BY ((a=2) SOUNDS LIKE '0')
< ok 0 0
> query
| SHOW CREATE VIEW v_so2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_so2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_so2` AS select `t`.`a` AS `a` from `t` union select `t`.`b` AS `b` from `t` order by soundex(`a` = 2) = soundex('0')
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_so2 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_so2
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_so2` AS select `t`.`a` AS `a` from `t` union select `t`.`b` AS `b` from `t` order by soundex(`a` = 2) = soundex('0')
< ok 0 0
> query
| SELECT * FROM v_so2 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_so3 AS SELECT a FROM t UNION ALL SELECT b FROM t ORDER BY NOT (a > 2) LIMIT 5
< ok 0 0
> query
| SHOW CREATE VIEW v_so3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_so3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_so3` AS select `t`.`a` AS `a` from `t` union all select `t`.`b` AS `b` from `t` order by `a` <= 2 limit 5
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_so3 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 1
< row
: int64 1
< row
: int64 2
< row
: int64 2
> exec
| DROP VIEW v_so3
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_so3` AS select `t`.`a` AS `a` from `t` union all select `t`.`b` AS `b` from `t` order by `a` <= 2 limit 5
< ok 0 0
> query
| SELECT * FROM v_so3 ORDER BY 1
< result
< column "a" "INT"
< row
: int64 1
< row
: int64 1
< row
: int64 2
< row
: int64 2
