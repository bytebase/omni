-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparse_Section_7_3_ExpressionViews default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_arithmetic_expr
< ok 0 0
> exec
| CREATE VIEW v_arithmetic_expr AS SELECT a + b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_arithmetic_expr
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_arithmetic_expr
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_arithmetic_expr` (`a + b`) AS SELECT `a`+`b` AS `a + b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_function_call
< ok 0 0
> exec
| CREATE VIEW v_function_call AS SELECT CONCAT(a, b) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_function_call
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_function_call
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_function_call` (`CONCAT(a, b)`) AS SELECT CONCAT(`a`, `b`) AS `CONCAT(a, b)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_case_expr
< ok 0 0
> exec
| CREATE VIEW v_case_expr AS SELECT CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_case_expr
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_case_expr
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_case_expr` (`CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END`) AS SELECT CASE WHEN `a`>0 THEN _UTF8MB4'pos' ELSE _UTF8MB4'neg' END AS `CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cast_expr
< ok 0 0
> exec
| CREATE VIEW v_cast_expr AS SELECT CAST(a AS CHAR) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cast_expr
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_expr
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_expr` (`CAST(a AS CHAR)`) AS SELECT CAST(`a` AS CHAR(11)) AS `CAST(a AS CHAR)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_aggregate_expr
< ok 0 0
> exec
| CREATE VIEW v_aggregate_expr AS SELECT COUNT(*), SUM(a) FROM t GROUP BY a HAVING SUM(a) > 10
< ok 0 0
> query
| SHOW CREATE VIEW v_aggregate_expr
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_aggregate_expr
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_aggregate_expr` (`COUNT(*)`, `SUM(a)`) AS SELECT COUNT(1) AS `COUNT(*)`,SUM(`a`) AS `SUM(a)` FROM `test`.`t` GROUP BY `a` HAVING SUM(`test`.`t`.`a`)>10
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_national_string
< ok 0 0
> exec
| CREATE VIEW v_national_string AS SELECT N'nat', n'x' AS y, N'a' 'b', N'' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_national_string
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_national_string
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_national_string` (`nat`, `y`, `b`, ``) AS SELECT _UTF8'nat' AS `nat`,_UTF8'x' AS `y`,_UTF8'a' AS `b`,_UTF8'' FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_trim_direction
< ok 0 0
> exec
| CREATE VIEW v_trim_direction AS SELECT TRIM(BOTH FROM c), TRIM(LEADING FROM c), TRIM(TRAILING FROM c), TRIM(BOTH 'x' FROM c) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_direction
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_direction
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_direction` (`TRIM(BOTH FROM c)`, `TRIM(LEADING FROM c)`, `TRIM(TRAILING FROM c)`, `TRIM(BOTH 'x' FROM c)`) AS SELECT TRIM(BOTH _UTF8MB4' ' FROM `c`) AS `TRIM(BOTH FROM c)`,TRIM(LEADING _UTF8MB4' ' FROM `c`) AS `TRIM(LEADING FROM c)`,TRIM(TRAILING _UTF8MB4' ' FROM `c`) AS `TRIM(TRAILING FROM c)`,TRIM(BOTH _UTF8MB4'x' FROM `c`) AS `TRIM(BOTH 'x' FROM c)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_char_using
< ok 0 0
> exec
| CREATE VIEW v_char_using AS SELECT CHAR(65 USING utf8mb4), CHAR(65, 66 USING 'Latin1'), CHAR(65 USING utf8mb3), CHAR(65) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_char_using
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_char_using
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_char_using` (`CHAR(65 USING utf8mb4)`, `CHAR(65, 66 USING 'Latin1')`, `CHAR(65 USING utf8mb3)`, `CHAR(65)`) AS SELECT CHAR_FUNC(65, 'utf8mb4') AS `CHAR(65 USING utf8mb4)`,CHAR_FUNC(65, 66, 'latin1') AS `CHAR(65, 66 USING 'Latin1')`,CHAR_FUNC(65, 'utf8') AS `CHAR(65 USING utf8mb3)`,CHAR_FUNC(65, NULL) AS `CHAR(65)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_row_shorthand
< ok 0 0
> exec
| CREATE VIEW v_row_shorthand AS SELECT (a, b) IN ((1, 2)), (a, b) = (1, 2), ROW(a, b) < (1, 2) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_row_shorthand
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_row_shorthand
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_row_shorthand` (`(a, b) IN ((1, 2))`, `(a, b) = (1, 2)`, `ROW(a, b) < (1, 2)`) AS SELECT ROW(`a`,`b`) IN (ROW(1,2)) AS `(a, b) IN ((1, 2))`,ROW(`a`,`b`)=ROW(1,2) AS `(a, b) = (1, 2)`,ROW(`a`,`b`)<ROW(1,2) AS `ROW(a, b) < (1, 2)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_low_priority
< ok 0 0
> exec
| CREATE VIEW v_low_priority AS SELECT LOW_PRIORITY a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_low_priority
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_low_priority
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_low_priority` (`a`) AS SELECT LOW_PRIORITY `a` AS `a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_high_priority
< ok 0 0
> exec
| CREATE VIEW v_high_priority AS SELECT HIGH_PRIORITY a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_high_priority
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_high_priority
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_high_priority` (`a`) AS SELECT HIGH_PRIORITY `a` AS `a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_delayed
< ok 0 0
> exec
| CREATE VIEW v_delayed AS SELECT DELAYED a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_delayed
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_delayed
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_delayed` (`a`) AS SELECT DELAYED `a` AS `a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_last_priority_wins
< ok 0 0
> exec
| CREATE VIEW v_last_priority_wins AS SELECT HIGH_PRIORITY LOW_PRIORITY a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_last_priority_wins
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_last_priority_wins
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_last_priority_wins` (`a`) AS SELECT LOW_PRIORITY `a` AS `a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_sql_cache
< ok 0 0
> exec
| CREATE VIEW v_sql_cache AS SELECT SQL_CACHE a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_sql_cache
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_sql_cache
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_sql_cache` (`a`) AS SELECT `a` AS `a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_sql_cache_and_no_cache
< ok 0 0
> exec
| CREATE VIEW v_sql_cache_and_no_cache AS SELECT SQL_CACHE SQL_NO_CACHE a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_sql_cache_and_no_cache
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_sql_cache_and_no_cache
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_sql_cache_and_no_cache` (`a`) AS SELECT SQL_NO_CACHE `a` AS `a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
