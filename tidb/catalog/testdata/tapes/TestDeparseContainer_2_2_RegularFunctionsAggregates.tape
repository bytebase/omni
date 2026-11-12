-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparseContainer_2_2_RegularFunctionsAggregates default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_str_funcs
< ok 0 0
> exec
| CREATE VIEW v_str_funcs AS SELECT CONCAT(a, b), UPPER(a), LOWER(a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_str_funcs
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_str_funcs
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_str_funcs` (`CONCAT(a, b)`, `UPPER(a)`, `LOWER(a)`) AS SELECT CONCAT(`a`, `b`) AS `CONCAT(a, b)`,UPPER(`a`) AS `UPPER(a)`,LOWER(`a`) AS `LOWER(a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_null_funcs
< ok 0 0
> exec
| CREATE VIEW v_null_funcs AS SELECT IFNULL(a, 0), COALESCE(a, b, 0), NULLIF(a, 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_null_funcs
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_null_funcs
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_null_funcs` (`IFNULL(a, 0)`, `COALESCE(a, b, 0)`, `NULLIF(a, 0)`) AS SELECT IFNULL(`a`, 0) AS `IFNULL(a, 0)`,COALESCE(`a`, `b`, 0) AS `COALESCE(a, b, 0)`,NULLIF(`a`, 0) AS `NULLIF(a, 0)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_if_func
< ok 0 0
> exec
| CREATE VIEW v_if_func AS SELECT IF(a > 0, 'yes', 'no') FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_if_func
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_if_func
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_if_func` (`IF(a > 0, 'yes', 'no')`) AS SELECT IF(`a`>0, _UTF8MB4'yes', _UTF8MB4'no') AS `IF(a > 0, 'yes', 'no')` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_num_funcs
< ok 0 0
> exec
| CREATE VIEW v_num_funcs AS SELECT ABS(a), GREATEST(a, b), LEAST(a, b) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_num_funcs
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_num_funcs
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_num_funcs` (`ABS(a)`, `GREATEST(a, b)`, `LEAST(a, b)`) AS SELECT ABS(`a`) AS `ABS(a)`,GREATEST(`a`, `b`) AS `GREATEST(a, b)`,LEAST(`a`, `b`) AS `LEAST(a, b)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_agg_funcs
< ok 0 0
> exec
| CREATE VIEW v_agg_funcs AS SELECT SUM(a), AVG(a), MAX(a), MIN(a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_agg_funcs
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_agg_funcs
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_agg_funcs` (`SUM(a)`, `AVG(a)`, `MAX(a)`, `MIN(a)`) AS SELECT SUM(`a`) AS `SUM(a)`,AVG(`a`) AS `AVG(a)`,MAX(`a`) AS `MAX(a)`,MIN(`a`) AS `MIN(a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_nested_funcs
< ok 0 0
> exec
| CREATE VIEW v_nested_funcs AS SELECT CONCAT(UPPER(TRIM(a)), LOWER(b)) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_nested_funcs
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nested_funcs
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nested_funcs` (`CONCAT(UPPER(TRIM(a)), LOWER(b))`) AS SELECT CONCAT(UPPER(TRIM(`a`)), LOWER(`b`)) AS `CONCAT(UPPER(TRIM(a)), LOWER(b))` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_multi_agg
< ok 0 0
> exec
| CREATE VIEW v_multi_agg AS SELECT COUNT(*), SUM(a), AVG(b), MAX(c) FROM t GROUP BY a
< ok 0 0
> query
| SHOW CREATE VIEW v_multi_agg
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_multi_agg
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_multi_agg` (`COUNT(*)`, `SUM(a)`, `AVG(b)`, `MAX(c)`) AS SELECT COUNT(1) AS `COUNT(*)`,SUM(`a`) AS `SUM(a)`,AVG(`b`) AS `AVG(b)`,MAX(`c`) AS `MAX(c)` FROM `test`.`t` GROUP BY `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
