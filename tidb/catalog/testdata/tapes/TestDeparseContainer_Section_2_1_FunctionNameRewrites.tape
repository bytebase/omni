-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_substr` (`SUBSTRING('abc', 1, 2)`) AS SELECT SUBSTRING(_UTF8MB4'abc', 1, 2) AS `SUBSTRING('abc', 1, 2)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cur_ts` (`CURRENT_TIMESTAMP`) AS SELECT CURRENT_TIMESTAMP() AS `CURRENT_TIMESTAMP` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cur_ts_fn` (`CURRENT_TIMESTAMP()`) AS SELECT CURRENT_TIMESTAMP() AS `CURRENT_TIMESTAMP()` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cur_date` (`CURRENT_DATE`) AS SELECT CURRENT_DATE() AS `CURRENT_DATE` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cur_time` (`CURRENT_TIME`) AS SELECT CURRENT_TIME() AS `CURRENT_TIME` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cur_user` (`CURRENT_USER`) AS SELECT CURRENT_USER() AS `CURRENT_USER` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_now` (`NOW()`) AS SELECT NOW() AS `NOW()` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_count_star` (`COUNT(*)`) AS SELECT COUNT(1) AS `COUNT(*)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_count_dist` (`COUNT(DISTINCT a)`) AS SELECT COUNT(DISTINCT `a`) AS `COUNT(DISTINCT a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
