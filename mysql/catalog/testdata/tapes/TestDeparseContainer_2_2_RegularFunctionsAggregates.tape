-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_str_funcs` AS select concat(`t`.`a`,`t`.`b`) AS `CONCAT(a, b)`,upper(`t`.`a`) AS `UPPER(a)`,lower(`t`.`a`) AS `LOWER(a)` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_null_funcs` AS select ifnull(`t`.`a`,0) AS `IFNULL(a, 0)`,coalesce(`t`.`a`,`t`.`b`,0) AS `COALESCE(a, b, 0)`,nullif(`t`.`a`,0) AS `NULLIF(a, 0)` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_if_func` AS select if((`t`.`a` > 0),'yes','no') AS `IF(a > 0, 'yes', 'no')` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_num_funcs` AS select abs(`t`.`a`) AS `ABS(a)`,greatest(`t`.`a`,`t`.`b`) AS `GREATEST(a, b)`,least(`t`.`a`,`t`.`b`) AS `LEAST(a, b)` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_agg_funcs` AS select sum(`t`.`a`) AS `SUM(a)`,avg(`t`.`a`) AS `AVG(a)`,max(`t`.`a`) AS `MAX(a)`,min(`t`.`a`) AS `MIN(a)` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nested_funcs` AS select concat(upper(trim(`t`.`a`)),lower(`t`.`b`)) AS `CONCAT(UPPER(TRIM(a)), LOWER(b))` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_multi_agg` AS select count(0) AS `COUNT(*)`,sum(`t`.`a`) AS `SUM(a)`,avg(`t`.`b`) AS `AVG(b)`,max(`t`.`c`) AS `MAX(c)` from `t` group by `t`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
