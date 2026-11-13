-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_substr
: bytes
| CREATE VIEW `v_substr` AS select substr('abc',1,2) AS `SUBSTRING('abc', 1, 2)` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cur_ts
: bytes
| CREATE VIEW `v_cur_ts` AS select now() AS `CURRENT_TIMESTAMP` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cur_ts_fn
: bytes
| CREATE VIEW `v_cur_ts_fn` AS select now() AS `CURRENT_TIMESTAMP()` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cur_date
: bytes
| CREATE VIEW `v_cur_date` AS select curdate() AS `CURRENT_DATE` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cur_time
: bytes
| CREATE VIEW `v_cur_time` AS select curtime() AS `CURRENT_TIME` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cur_user
: bytes
| CREATE VIEW `v_cur_user` AS select current_user() AS `CURRENT_USER` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_now
: bytes
| CREATE VIEW `v_now` AS select now() AS `NOW()` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_count_dist
: bytes
| CREATE VIEW `v_count_dist` AS select count(distinct `test`.`t`.`a`) AS `COUNT(DISTINCT a)` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
