-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ValuesRowsTakeTheRewrites default#1
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
| CREATE VIEW v_vs0 AS VALUES (((2=1) SOUNDS LIKE '0'))
< ok 0 0
> query
| SHOW CREATE VIEW v_vs0
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_vs0
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs0` AS values (soundex(2 = 1) = soundex('0'))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_vs0
< result
< column "soundex(2 = 1) = soundex('0')" "INT"
< row
: int64 1
> exec
| DROP VIEW v_vs0
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs0` AS values (soundex(2 = 1) = soundex('0'))
< ok 0 0
> query
| SELECT * FROM v_vs0
< result
< column "soundex(2 = 1) = soundex('0')" "INT"
< row
: int64 1
> exec
| CREATE VIEW v_vs1 AS VALUES (NOT (1 > 0))
< ok 0 0
> query
| SHOW CREATE VIEW v_vs1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_vs1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs1` AS values (1 <= 0)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_vs1
< result
< column "1 <= 0" "INT"
< row
: int64 0
> exec
| DROP VIEW v_vs1
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs1` AS values (1 <= 0)
< ok 0 0
> query
| SELECT * FROM v_vs1
< result
< column "1 <= 0" "INT"
< row
: int64 0
> exec
| CREATE VIEW v_vs2 AS VALUES (((2=1) SOUNDS LIKE '0')), (1)
< ok 0 0
> query
| SHOW CREATE VIEW v_vs2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_vs2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs2` AS values (soundex(2 = 1) = soundex('0')),(1)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_vs2
< result
< column "soundex(2 = 1) = soundex('0')" "INT"
< row
: int64 1
< row
: int64 1
> exec
| DROP VIEW v_vs2
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs2` AS values (soundex(2 = 1) = soundex('0')),(1)
< ok 0 0
> query
| SELECT * FROM v_vs2
< result
< column "soundex(2 = 1) = soundex('0')" "INT"
< row
: int64 1
< row
: int64 1
> exec
| CREATE VIEW v_vs3 AS VALUES (1 AND 2)
< ok 0 0
> query
| SHOW CREATE VIEW v_vs3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_vs3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs3` AS values (1 and 2)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_vs3
< result
< column "1 and 2" "INT"
< row
: int64 1
> exec
| DROP VIEW v_vs3
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs3` AS values (1 and 2)
< ok 0 0
> query
| SELECT * FROM v_vs3
< result
< column "1 and 2" "INT"
< row
: int64 1
> exec
| CREATE VIEW v_vs4 AS SELECT 1 UNION VALUES (((2=1) SOUNDS LIKE '0'))
< ok 0 0
> query
| SHOW CREATE VIEW v_vs4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_vs4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs4` AS select 1 AS `1` union values (soundex(2 = 1) = soundex('0'))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_vs4
< result
< column "1" "INT"
< row
: int64 1
> exec
| DROP VIEW v_vs4
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs4` AS select 1 AS `1` union values (soundex(2 = 1) = soundex('0'))
< ok 0 0
> query
| SELECT * FROM v_vs4
< result
< column "1" "INT"
< row
: int64 1
> exec
| CREATE VIEW v_vs5 AS VALUES (1),(2) ORDER BY 1
< ok 0 0
> query
| SHOW CREATE VIEW v_vs5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_vs5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs5` AS values (1),(2) order by 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_vs5
< result
< column "1" "INT"
< row
: int64 1
< row
: int64 2
> exec
| DROP VIEW v_vs5
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs5` AS values (1),(2) order by 1
< ok 0 0
> query
| SELECT * FROM v_vs5
< result
< column "1" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE VIEW v_vs6 AS VALUES (2),(1) ORDER BY 1 DESC
< ok 0 0
> query
| SHOW CREATE VIEW v_vs6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_vs6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs6` AS values (2),(1) order by 1 desc
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_vs6
< result
< column "2" "INT"
< row
: int64 2
< row
: int64 1
> exec
| DROP VIEW v_vs6
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs6` AS values (2),(1) order by 1 desc
< ok 0 0
> query
| SELECT * FROM v_vs6
< result
< column "2" "INT"
< row
: int64 2
< row
: int64 1
> exec
| CREATE VIEW v_vs7 AS VALUES (((2=1) SOUNDS LIKE '0')) ORDER BY 1
< ok 0 0
> query
| SHOW CREATE VIEW v_vs7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_vs7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs7` AS values (soundex(2 = 1) = soundex('0')) order by 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_vs7
< result
< column "soundex(2 = 1) = soundex('0')" "INT"
< row
: int64 1
> exec
| DROP VIEW v_vs7
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs7` AS values (soundex(2 = 1) = soundex('0')) order by 1
< ok 0 0
> query
| SELECT * FROM v_vs7
< result
< column "soundex(2 = 1) = soundex('0')" "INT"
< row
: int64 1
> exec
| CREATE VIEW v_vs8 AS VALUES (1),(2) LIMIT 1
< ok 0 0
> query
| SHOW CREATE VIEW v_vs8
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_vs8
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs8` AS values (1),(2) limit 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT * FROM v_vs8
< result
< column "1" "INT"
< row
: int64 1
> exec
| DROP VIEW v_vs8
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_vs8` AS values (1),(2) limit 1
< ok 0 0
> query
| SELECT * FROM v_vs8
< result
< column "1" "INT"
< row
: int64 1
