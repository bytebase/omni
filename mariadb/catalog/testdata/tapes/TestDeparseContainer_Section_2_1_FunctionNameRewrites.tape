-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestDeparseContainer_Section_2_1_FunctionNameRewrites default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_substr
< ok 0 0
> exec
| CREATE VIEW v_substr AS SELECT SUBSTRING('abc', 1, 2) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_substr
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_substr
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_substr` AS select substr('abc',1,2) AS `SUBSTRING('abc', 1, 2)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_cur_ts
< ok 0 0
> exec
| CREATE VIEW v_cur_ts AS SELECT CURRENT_TIMESTAMP FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cur_ts
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cur_ts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cur_ts` AS select current_timestamp() AS `CURRENT_TIMESTAMP` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_cur_ts_fn
< ok 0 0
> exec
| CREATE VIEW v_cur_ts_fn AS SELECT CURRENT_TIMESTAMP() FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cur_ts_fn
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cur_ts_fn
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cur_ts_fn` AS select current_timestamp() AS `CURRENT_TIMESTAMP()` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_cur_date
< ok 0 0
> exec
| CREATE VIEW v_cur_date AS SELECT CURRENT_DATE FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cur_date
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cur_date
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cur_date` AS select curdate() AS `CURRENT_DATE` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_cur_time
< ok 0 0
> exec
| CREATE VIEW v_cur_time AS SELECT CURRENT_TIME FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cur_time
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cur_time
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cur_time` AS select curtime() AS `CURRENT_TIME` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_cur_user
< ok 0 0
> exec
| CREATE VIEW v_cur_user AS SELECT CURRENT_USER FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cur_user
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cur_user
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cur_user` AS select current_user() AS `CURRENT_USER` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_now
< ok 0 0
> exec
| CREATE VIEW v_now AS SELECT NOW() FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_now
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_now
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_now` AS select current_timestamp() AS `NOW()` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_count_star
< ok 0 0
> exec
| CREATE VIEW v_count_star AS SELECT COUNT(*) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_count_star
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_count_star
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_count_star` AS select count(0) AS `COUNT(*)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_count_dist
< ok 0 0
> exec
| CREATE VIEW v_count_dist AS SELECT COUNT(DISTINCT a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_count_dist
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_count_dist
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_count_dist` AS select count(distinct `t`.`a`) AS `COUNT(DISTINCT a)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
