-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_ViewIdempotence default#1

== TestOracle_ViewIdempotence/8.0/algorithm-merge default#1
> exec
| DROP DATABASE IF EXISTS videm_algorithm_merge
< ok 0 0
> exec
| CREATE DATABASE videm_algorithm_merge DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_algorithm_merge
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE ALGORITHM=MERGE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_algorithm_merge.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_algorithm_merge.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=MERGE DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/algorithm-temptable default#1
> exec
| DROP DATABASE IF EXISTS videm_algorithm_temptable
< ok 0 0
> exec
| CREATE DATABASE videm_algorithm_temptable DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_algorithm_temptable
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE ALGORITHM=TEMPTABLE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_algorithm_temptable.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_algorithm_temptable.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=TEMPTABLE DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/check-cascaded default#1
> exec
| DROP DATABASE IF EXISTS videm_check_cascaded
< ok 0 0
> exec
| CREATE DATABASE videm_check_cascaded DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_check_cascaded
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION
< ok 0 0
> query
| SHOW CREATE TABLE videm_check_cascaded.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_check_cascaded.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t` where (`t`.`a` > 0) WITH CASCADED CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/check-local default#1
> exec
| DROP DATABASE IF EXISTS videm_check_local
< ok 0 0
> exec
| CREATE DATABASE videm_check_local DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_check_local
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a > 0 WITH LOCAL CHECK OPTION
< ok 0 0
> query
| SHOW CREATE TABLE videm_check_local.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_check_local.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t` where (`t`.`a` > 0) WITH LOCAL CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/cte default#1
> exec
| DROP DATABASE IF EXISTS videm_cte
< ok 0 0
> exec
| CREATE DATABASE videm_cte DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_cte
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS WITH cte AS (SELECT a, b FROM t) SELECT a FROM cte
< ok 0 0
> query
| SHOW CREATE TABLE videm_cte.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_cte.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS with `cte` as (select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`) select `cte`.`a` AS `a` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/date-arith-family default#1
> exec
| DROP DATABASE IF EXISTS videm_date_arith_family
< ok 0 0
> exec
| CREATE DATABASE videm_date_arith_family DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_date_arith_family
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT TIMESTAMPADD(HOUR, 1, a) AS c1, DATE_ADD(a, INTERVAL 1 DAY) AS c2, DATE_SUB(a, INTERVAL 30 MINUTE) AS c3, ADDDATE(a, INTERVAL 1 WEEK) AS c4, ADDDATE(a, 31) AS c5, SUBDATE(a, INTERVAL 1 MONTH) AS c6, SUBDATE(a, 31) AS c7 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_date_arith_family.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_date_arith_family.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (`t`.`a` + interval 1 hour) AS `c1`,(`t`.`a` + interval 1 day) AS `c2`,(`t`.`a` - interval 30 minute) AS `c3`,(`t`.`a` + interval 1 week) AS `c4`,(`t`.`a` + interval 31 day) AS `c5`,(`t`.`a` - interval 1 month) AS `c6`,(`t`.`a` - interval 31 day) AS `c7` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/derived-table-from default#1
> exec
| DROP DATABASE IF EXISTS videm_derived_table_from
< ok 0 0
> exec
| CREATE DATABASE videm_derived_table_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_derived_table_from
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT d.a FROM (SELECT a, b FROM t WHERE b > 0) d
< ok 0 0
> query
| SHOW CREATE TABLE videm_derived_table_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_derived_table_from.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `d`.`a` AS `a` from (select `t`.`a` AS `a`,`t`.`b` AS `b` from `t` where (`t`.`b` > 0)) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/derived-table-join default#1
> exec
| DROP DATABASE IF EXISTS videm_derived_table_join
< ok 0 0
> exec
| CREATE DATABASE videm_derived_table_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_derived_table_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a, d.y FROM t JOIN (SELECT x, y FROM u) d ON t.c = d.x
< ok 0 0
> query
| SHOW CREATE TABLE videm_derived_table_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_derived_table_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_derived_table_join.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a`,`d`.`y` AS `y` from (`t` join (select `u`.`x` AS `x`,`u`.`y` AS `y` from `u`) `d` on((`t`.`c` = `d`.`x`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/explicit-columns default#1
> exec
| DROP DATABASE IF EXISTS videm_explicit_columns
< ok 0 0
> exec
| CREATE DATABASE videm_explicit_columns DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_explicit_columns
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v (p, q) AS SELECT a, b FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_explicit_columns.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_explicit_columns.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`p`,`q`) AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/expr-and-func default#1
> exec
| DROP DATABASE IF EXISTS videm_expr_and_func
< ok 0 0
> exec
| CREATE DATABASE videm_expr_and_func DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_expr_and_func
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a + 1 AS s, UPPER(b) AS ub FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_expr_and_func.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_expr_and_func.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (`t`.`a` + 1) AS `s`,upper(`t`.`b`) AS `ub` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/extract-units default#1
> exec
| DROP DATABASE IF EXISTS videm_extract_units
< ok 0 0
> exec
| CREATE DATABASE videm_extract_units DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_extract_units
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT EXTRACT(SECOND FROM a) AS c1, extract(day from a) AS c2, EXTRACT(QUARTER FROM a) AS c3, EXTRACT(DAY_HOUR FROM a) AS c4, EXTRACT(YEAR_MONTH FROM a) AS c5, EXTRACT(HOUR_MICROSECOND FROM a) AS c6, EXTRACT(SQL_TSI_WEEK FROM a) AS c7 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_extract_units.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_extract_units.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select extract(second from `t`.`a`) AS `c1`,extract(day from `t`.`a`) AS `c2`,extract(quarter from `t`.`a`) AS `c3`,extract(day_hour from `t`.`a`) AS `c4`,extract(year_month from `t`.`a`) AS `c5`,extract(hour_microsecond from `t`.`a`) AS `c6`,extract(week from `t`.`a`) AS `c7` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/get-format-kinds default#1
> exec
| DROP DATABASE IF EXISTS videm_get_format_kinds
< ok 0 0
> exec
| CREATE DATABASE videm_get_format_kinds DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_get_format_kinds
< ok 0 0
> exec
| CREATE VIEW v AS SELECT GET_FORMAT(DATE, 'USA') AS c1, get_format(time, 'ISO') AS c2, GET_FORMAT(DATETIME, 'EUR') AS c3, GET_FORMAT(TIMESTAMP, 'JIS') AS c4
< ok 0 0
> query
| SHOW CREATE VIEW videm_get_format_kinds.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select get_format(DATE, 'USA') AS `c1`,get_format(TIME, 'ISO') AS `c2`,get_format(DATETIME, 'EUR') AS `c3`,get_format(DATETIME, 'JIS') AS `c4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/group-by-agg default#1
> exec
| DROP DATABASE IF EXISTS videm_group_by_agg
< ok 0 0
> exec
| CREATE DATABASE videm_group_by_agg DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_group_by_agg
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a, COUNT(*) AS n FROM t GROUP BY a
< ok 0 0
> query
| SHOW CREATE TABLE videm_group_by_agg.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_group_by_agg.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a`,count(0) AS `n` from `t` group by `t`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/interval-forms default#1
> exec
| DROP DATABASE IF EXISTS videm_interval_forms
< ok 0 0
> exec
| CREATE DATABASE videm_interval_forms DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_interval_forms
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME, n INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a + INTERVAL 1 DAY AS c1, a - INTERVAL 30 MINUTE AS c2, a + INTERVAL '1:1' DAY_HOUR AS c3, a + INTERVAL '1-1' YEAR_MONTH AS c4, a + INTERVAL (n + 1) DAY AS c5, ADDDATE(a, -5) AS c6, INTERVAL 1 DAY + a AS c7 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_interval_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL,
|   `n` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_interval_forms.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (`t`.`a` + interval 1 day) AS `c1`,(`t`.`a` - interval 30 minute) AS `c2`,(`t`.`a` + interval '1:1' day_hour) AS `c3`,(`t`.`a` + interval '1-1' year_month) AS `c4`,(`t`.`a` + interval (`t`.`n` + 1) day) AS `c5`,(`t`.`a` + interval -(5) day) AS `c6`,(`t`.`a` + interval 1 day) AS `c7` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/join default#1
> exec
| DROP DATABASE IF EXISTS videm_join
< ok 0 0
> exec
| CREATE DATABASE videm_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20), c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a, u.y FROM t JOIN u ON t.c = u.x
< ok 0 0
> query
| SHOW CREATE TABLE videm_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_join.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a`,`u`.`y` AS `y` from (`t` join `u` on((`t`.`c` = `u`.`x`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/join-3-table default#1
> exec
| DROP DATABASE IF EXISTS videm_join_3_table
< ok 0 0
> exec
| CREATE DATABASE videm_join_3_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join_3_table
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a, u.y, w.n FROM t JOIN u ON t.c = u.x JOIN w ON u.y = w.m
< ok 0 0
> query
| SHOW CREATE TABLE videm_join_3_table.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_3_table.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_3_table.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int DEFAULT NULL,
|   `n` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_join_3_table.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a`,`u`.`y` AS `y`,`w`.`n` AS `n` from ((`t` join `u` on((`t`.`c` = `u`.`x`))) join `w` on((`u`.`y` = `w`.`m`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/join-4-table default#1
> exec
| DROP DATABASE IF EXISTS videm_join_4_table
< ok 0 0
> exec
| CREATE DATABASE videm_join_4_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join_4_table
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE TABLE z (p INT, q INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a FROM t JOIN u ON t.c = u.x JOIN w ON u.y = w.m JOIN z ON w.n = z.p
< ok 0 0
> query
| SHOW CREATE TABLE videm_join_4_table.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_4_table.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_4_table.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int DEFAULT NULL,
|   `n` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_4_table.z
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| z
: bytes
| CREATE TABLE `z` (
|   `p` int DEFAULT NULL,
|   `q` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_join_4_table.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from (((`t` join `u` on((`t`.`c` = `u`.`x`))) join `w` on((`u`.`y` = `w`.`m`))) join `z` on((`w`.`n` = `z`.`p`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/join-aliased-compound-on default#1
> exec
| DROP DATABASE IF EXISTS videm_join_aliased_compound_on
< ok 0 0
> exec
| CREATE DATABASE videm_join_aliased_compound_on DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join_aliased_compound_on
< ok 0 0
> exec
| CREATE TABLE dept_emp (emp_no INT NOT NULL, dept_no CHAR(4) NOT NULL, from_date DATE NOT NULL, to_date DATE NOT NULL, PRIMARY KEY (emp_no, dept_no))
< ok 0 0
> exec
| CREATE TABLE departments (dept_no CHAR(4) NOT NULL, dept_name VARCHAR(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE VIEW dept_emp_latest_date AS SELECT emp_no, MAX(from_date) AS from_date, MAX(to_date) AS to_date FROM dept_emp GROUP BY emp_no
< ok 0 0
> exec
| CREATE VIEW current_dept_emp AS SELECT l.emp_no, d.dept_no, l.from_date, l.to_date FROM dept_emp d JOIN dept_emp_latest_date l ON d.emp_no = l.emp_no AND d.from_date = l.from_date JOIN departments dp ON d.dept_no = dp.dept_no
< ok 0 0
> query
| SHOW CREATE TABLE videm_join_aliased_compound_on.dept_emp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| dept_emp
: bytes
| CREATE TABLE `dept_emp` (
|   `emp_no` int NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   `to_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_aliased_compound_on.departments
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| departments
: bytes
| CREATE TABLE `departments` (
|   `dept_no` char(4) NOT NULL,
|   `dept_name` varchar(40) NOT NULL,
|   PRIMARY KEY (`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_join_aliased_compound_on.dept_emp_latest_date
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| dept_emp_latest_date
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `dept_emp_latest_date` AS select `dept_emp`.`emp_no` AS `emp_no`,max(`dept_emp`.`from_date`) AS `from_date`,max(`dept_emp`.`to_date`) AS `to_date` from `dept_emp` group by `dept_emp`.`emp_no`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW videm_join_aliased_compound_on.current_dept_emp
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| current_dept_emp
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `current_dept_emp` AS select `l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`l`.`from_date` AS `from_date`,`l`.`to_date` AS `to_date` from ((`dept_emp` `d` join `dept_emp_latest_date` `l` on(((`d`.`emp_no` = `l`.`emp_no`) and (`d`.`from_date` = `l`.`from_date`)))) join `departments` `dp` on((`d`.`dept_no` = `dp`.`dept_no`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/join-aliased-left-inner-compound-on default#1
> exec
| DROP DATABASE IF EXISTS videm_join_aliased_left_inner_compound_on
< ok 0 0
> exec
| CREATE DATABASE videm_join_aliased_left_inner_compound_on DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join_aliased_left_inner_compound_on
< ok 0 0
> exec
| CREATE TABLE a (id INT, k INT, v INT)
< ok 0 0
> exec
| CREATE TABLE b (id INT, k INT, w INT)
< ok 0 0
> exec
| CREATE TABLE c (id INT, k INT, z INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT x.v, y.w, q.z FROM a x LEFT JOIN b y ON x.id = y.id AND x.k = y.k INNER JOIN c q ON y.id = q.id AND y.k = q.k
< ok 0 0
> query
| SHOW CREATE TABLE videm_join_aliased_left_inner_compound_on.a
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `id` int DEFAULT NULL,
|   `k` int DEFAULT NULL,
|   `v` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_aliased_left_inner_compound_on.b
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| b
: bytes
| CREATE TABLE `b` (
|   `id` int DEFAULT NULL,
|   `k` int DEFAULT NULL,
|   `w` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_aliased_left_inner_compound_on.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int DEFAULT NULL,
|   `k` int DEFAULT NULL,
|   `z` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_join_aliased_left_inner_compound_on.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `x`.`v` AS `v`,`y`.`w` AS `w`,`q`.`z` AS `z` from ((`a` `x` left join `b` `y` on(((`x`.`id` = `y`.`id`) and (`x`.`k` = `y`.`k`)))) join `c` `q` on(((`y`.`id` = `q`.`id`) and (`y`.`k` = `q`.`k`))))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/join-group-by default#1
> exec
| DROP DATABASE IF EXISTS videm_join_group_by
< ok 0 0
> exec
| CREATE DATABASE videm_join_group_by DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join_group_by
< ok 0 0
> exec
| CREATE TABLE category (category_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_category (film_id INT, category_id INT)
< ok 0 0
> exec
| CREATE TABLE film (film_id INT, title VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_actor (actor_id INT, film_id INT)
< ok 0 0
> exec
| CREATE TABLE actor (actor_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT film.film_id AS FID, COUNT(*) AS n FROM category LEFT JOIN film_category ON category.category_id = film_category.category_id LEFT JOIN film ON film_category.film_id = film.film_id JOIN film_actor ON film.film_id = film_actor.film_id JOIN actor ON film_actor.actor_id = actor.actor_id GROUP BY film.film_id
< ok 0 0
> query
| SHOW CREATE TABLE videm_join_group_by.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `category_id` int DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_group_by.film_category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film_category
: bytes
| CREATE TABLE `film_category` (
|   `film_id` int DEFAULT NULL,
|   `category_id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_group_by.film
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film
: bytes
| CREATE TABLE `film` (
|   `film_id` int DEFAULT NULL,
|   `title` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_group_by.film_actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film_actor
: bytes
| CREATE TABLE `film_actor` (
|   `actor_id` int DEFAULT NULL,
|   `film_id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_group_by.actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| actor
: bytes
| CREATE TABLE `actor` (
|   `actor_id` int DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_join_group_by.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `film`.`film_id` AS `FID`,count(0) AS `n` from ((((`category` left join `film_category` on((`category`.`category_id` = `film_category`.`category_id`))) left join `film` on((`film_category`.`film_id` = `film`.`film_id`))) join `film_actor` on((`film`.`film_id` = `film_actor`.`film_id`))) join `actor` on((`film_actor`.`actor_id` = `actor`.`actor_id`))) group by `film`.`film_id`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/join-left-inner-mix default#1
> exec
| DROP DATABASE IF EXISTS videm_join_left_inner_mix
< ok 0 0
> exec
| CREATE DATABASE videm_join_left_inner_mix DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join_left_inner_mix
< ok 0 0
> exec
| CREATE TABLE category (category_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_category (film_id INT, category_id INT)
< ok 0 0
> exec
| CREATE TABLE film (film_id INT, title VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_actor (actor_id INT, film_id INT)
< ok 0 0
> exec
| CREATE TABLE actor (actor_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT film.film_id AS FID, category.name AS category, actor.name AS actor FROM category LEFT JOIN film_category ON category.category_id = film_category.category_id LEFT JOIN film ON film_category.film_id = film.film_id JOIN film_actor ON film.film_id = film_actor.film_id JOIN actor ON film_actor.actor_id = actor.actor_id
< ok 0 0
> query
| SHOW CREATE TABLE videm_join_left_inner_mix.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `category_id` int DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_left_inner_mix.film_category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film_category
: bytes
| CREATE TABLE `film_category` (
|   `film_id` int DEFAULT NULL,
|   `category_id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_left_inner_mix.film
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film
: bytes
| CREATE TABLE `film` (
|   `film_id` int DEFAULT NULL,
|   `title` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_left_inner_mix.film_actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film_actor
: bytes
| CREATE TABLE `film_actor` (
|   `actor_id` int DEFAULT NULL,
|   `film_id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_join_left_inner_mix.actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| actor
: bytes
| CREATE TABLE `actor` (
|   `actor_id` int DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_join_left_inner_mix.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `film`.`film_id` AS `FID`,`category`.`name` AS `category`,`actor`.`name` AS `actor` from ((((`category` left join `film_category` on((`category`.`category_id` = `film_category`.`category_id`))) left join `film` on((`film_category`.`film_id` = `film`.`film_id`))) join `film_actor` on((`film`.`film_id` = `film_actor`.`film_id`))) join `actor` on((`film_actor`.`actor_id` = `actor`.`actor_id`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/left-join default#1
> exec
| DROP DATABASE IF EXISTS videm_left_join
< ok 0 0
> exec
| CREATE DATABASE videm_left_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_left_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a, u.y FROM t LEFT JOIN u ON t.c = u.x
< ok 0 0
> query
| SHOW CREATE TABLE videm_left_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_left_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_left_join.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a`,`u`.`y` AS `y` from (`t` left join `u` on((`t`.`c` = `u`.`x`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/paren-exists-compare default#1
> exec
| DROP DATABASE IF EXISTS videm_paren_exists_compare
< ok 0 0
> exec
| CREATE DATABASE videm_paren_exists_compare DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_exists_compare
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT ((EXISTS(SELECT 1 FROM t)) = 1) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_exists_compare.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_paren_exists_compare.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (exists(select 1 from `t`) = 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/paren-exists-depth2 default#1
> exec
| DROP DATABASE IF EXISTS videm_paren_exists_depth2
< ok 0 0
> exec
| CREATE DATABASE videm_paren_exists_depth2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_exists_depth2
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT (EXISTS((SELECT 1 FROM t))) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_exists_depth2.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_paren_exists_depth2.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select exists(select 1 from `t`) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/paren-subquery-arith default#1
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_arith
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_arith DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_arith
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT ((SELECT max(a) FROM t) + 1) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_arith.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_paren_subquery_arith.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select ((select max(`t`.`a`) from `t`) + 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/paren-subquery-both-operands default#1
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_both_operands
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_both_operands DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_both_operands
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT ((SELECT min(a) FROM t) = (SELECT max(a) FROM t)) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_both_operands.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_paren_subquery_both_operands.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select ((select min(`t`.`a`) from `t`) = (select max(`t`.`a`) from `t`)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/paren-subquery-compare-if default#1
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_compare_if
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_compare_if DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_compare_if
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT if(((SELECT count(0) FROM t) = 0),'NO','YES') AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_compare_if.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_paren_subquery_compare_if.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select if(((select count(0) from `t`) = 0),'NO','YES') AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/paren-subquery-compare-if-depth2 default#1
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_compare_if_depth2
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_compare_if_depth2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_compare_if_depth2
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT if((((SELECT count(0) FROM t)) = 0),'NO','YES') AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_compare_if_depth2.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_paren_subquery_compare_if_depth2.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select if(((select count(0) from `t`) = 0),'NO','YES') AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/paren-subquery-in-depth2 default#1
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_in_depth2
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_in_depth2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_in_depth2
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT (1 IN ((SELECT max(a) FROM t))) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_in_depth2.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_paren_subquery_in_depth2.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select 1 in (select max(`t`.`a`) from `t`) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/paren-subquery-in-list default#1
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_in_list
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_in_list DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_in_list
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT (1 IN ((SELECT max(a) FROM t) + 1, 2)) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_in_list.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_paren_subquery_in_list.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (1 in (((select max(`t`.`a`) from `t`) + 1),2)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/paren-subquery-pure-depth2 default#1
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_pure_depth2
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_pure_depth2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_pure_depth2
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT ((SELECT count(0) FROM t)) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_pure_depth2.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_paren_subquery_pure_depth2.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (select count(0) from `t`) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/paren-subquery-row-expanded default#1
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_row_expanded
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_row_expanded DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_row_expanded
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT (((SELECT 1) = 1) AND (2 = 2)) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_row_expanded.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_paren_subquery_row_expanded.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (((select 1) = 1) and (2 = 2)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/paren-subquery-setop-compare default#1
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_setop_compare
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_setop_compare DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_setop_compare
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT ((SELECT 1 UNION SELECT 1) = 1) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_setop_compare.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_paren_subquery_setop_compare.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select ((select 1 union select 1) = 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/paren-subquery-where default#1
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_where
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_where DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_where
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE ((SELECT count(0) FROM t) > 0)
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_where.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_paren_subquery_where.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t` where ((select count(0) from `t`) > 0)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/simple default#1
> exec
| DROP DATABASE IF EXISTS videm_simple
< ok 0 0
> exec
| CREATE DATABASE videm_simple DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_simple
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20), c INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a, b FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_simple.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_simple.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/sql-security-invoker default#1
> exec
| DROP DATABASE IF EXISTS videm_sql_security_invoker
< ok 0 0
> exec
| CREATE DATABASE videm_sql_security_invoker DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_sql_security_invoker
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_sql_security_invoker.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_sql_security_invoker.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/star-expanded default#1
> exec
| DROP DATABASE IF EXISTS videm_star_expanded
< ok 0 0
> exec
| CREATE DATABASE videm_star_expanded DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_star_expanded
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20), c INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT * FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_star_expanded.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_star_expanded.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/substring-keyword-forms default#1
> exec
| DROP DATABASE IF EXISTS videm_substring_keyword_forms
< ok 0 0
> exec
| CREATE DATABASE videm_substring_keyword_forms DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_substring_keyword_forms
< ok 0 0
> exec
| CREATE TABLE t (s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT SUBSTRING(s FROM 2) AS c1, SUBSTRING(s FROM 2 FOR 3) AS c2, SUBSTRING(s, 2, 3) AS c3 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_substring_keyword_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `s` varchar(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_substring_keyword_forms.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select substr(`t`.`s`,2) AS `c1`,substr(`t`.`s`,2,3) AS `c2`,substr(`t`.`s`,2,3) AS `c3` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/trim-keyword-forms default#1
> exec
| DROP DATABASE IF EXISTS videm_trim_keyword_forms
< ok 0 0
> exec
| CREATE DATABASE videm_trim_keyword_forms DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_trim_keyword_forms
< ok 0 0
> exec
| CREATE TABLE t (s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT TRIM(BOTH 'x' FROM s) AS c1, TRIM(LEADING 'x' FROM s) AS c2, TRIM(TRAILING 'x' FROM s) AS c3, TRIM('x' FROM s) AS c4, TRIM(s) AS c5 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_trim_keyword_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `s` varchar(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_trim_keyword_forms.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select trim(both 'x' from `t`.`s`) AS `c1`,trim(leading 'x' from `t`.`s`) AS `c2`,trim(trailing 'x' from `t`.`s`) AS `c3`,trim('x' from `t`.`s`) AS `c4`,trim(`t`.`s`) AS `c5` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/tsdiff-all-units default#1
> exec
| DROP DATABASE IF EXISTS videm_tsdiff_all_units
< ok 0 0
> exec
| CREATE DATABASE videm_tsdiff_all_units DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_tsdiff_all_units
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME, b DATETIME)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT TIMESTAMPDIFF(MICROSECOND,a,b) AS c1, timestampdiff(second,a,b) AS c2, TIMESTAMPDIFF(MINUTE,a,b) AS c3, TIMESTAMPDIFF(HOUR,a,b) AS c4, TIMESTAMPDIFF(DAY,a,b) AS c5, TIMESTAMPDIFF(WEEK,a,b) AS c6, TIMESTAMPDIFF(MONTH,a,b) AS c7, TIMESTAMPDIFF(SQL_TSI_QUARTER,a,b) AS c8, TIMESTAMPDIFF(YEAR,a,b) AS c9 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_tsdiff_all_units.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL,
|   `b` datetime DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_tsdiff_all_units.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select timestampdiff(MICROSECOND,`t`.`a`,`t`.`b`) AS `c1`,timestampdiff(SECOND,`t`.`a`,`t`.`b`) AS `c2`,timestampdiff(MINUTE,`t`.`a`,`t`.`b`) AS `c3`,timestampdiff(HOUR,`t`.`a`,`t`.`b`) AS `c4`,timestampdiff(DAY,`t`.`a`,`t`.`b`) AS `c5`,timestampdiff(WEEK,`t`.`a`,`t`.`b`) AS `c6`,timestampdiff(MONTH,`t`.`a`,`t`.`b`) AS `c7`,timestampdiff(QUARTER,`t`.`a`,`t`.`b`) AS `c8`,timestampdiff(YEAR,`t`.`a`,`t`.`b`) AS `c9` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/tsdiff-column-now default#1
> exec
| DROP DATABASE IF EXISTS videm_tsdiff_column_now
< ok 0 0
> exec
| CREATE DATABASE videm_tsdiff_column_now DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_tsdiff_column_now
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id, TIMESTAMPDIFF(SECOND, created, NOW()) AS age_secs FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_tsdiff_column_now.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_tsdiff_column_now.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`id` AS `id`,timestampdiff(SECOND,`t`.`created`,now()) AS `age_secs` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/union default#1
> exec
| DROP DATABASE IF EXISTS videm_union
< ok 0 0
> exec
| CREATE DATABASE videm_union DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_union
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE TABLE u (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t UNION SELECT a FROM u
< ok 0 0
> query
| SHOW CREATE TABLE videm_union.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_union.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_union.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t` union select `u`.`a` AS `a` from `u`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/view-on-multi-join-view default#1
> exec
| DROP DATABASE IF EXISTS videm_view_on_multi_join_view
< ok 0 0
> exec
| CREATE DATABASE videm_view_on_multi_join_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_view_on_multi_join_view
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE VIEW base AS SELECT t.a AS a, u.y AS y, w.n AS n FROM t JOIN u ON t.c = u.x JOIN w ON u.y = w.m
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a, n FROM base WHERE a > 0
< ok 0 0
> query
| SHOW CREATE TABLE videm_view_on_multi_join_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_view_on_multi_join_view.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE videm_view_on_multi_join_view.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int DEFAULT NULL,
|   `n` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_view_on_multi_join_view.base
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `base` AS select `t`.`a` AS `a`,`u`.`y` AS `y`,`w`.`n` AS `n` from ((`t` join `u` on((`t`.`c` = `u`.`x`))) join `w` on((`u`.`y` = `w`.`m`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW videm_view_on_multi_join_view.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `base`.`a` AS `a`,`base`.`n` AS `n` from `base` where (`base`.`a` > 0)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/view-on-view default#1
> exec
| DROP DATABASE IF EXISTS videm_view_on_view
< ok 0 0
> exec
| CREATE DATABASE videm_view_on_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_view_on_view
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW base AS SELECT a, b FROM t
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM base
< ok 0 0
> query
| SHOW CREATE TABLE videm_view_on_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_view_on_view.base
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `base` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW videm_view_on_view.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `base`.`a` AS `a` from `base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/view-on-view-explicit-cols default#1
> exec
| DROP DATABASE IF EXISTS videm_view_on_view_explicit_cols
< ok 0 0
> exec
| CREATE DATABASE videm_view_on_view_explicit_cols DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_view_on_view_explicit_cols
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW base (x, y) AS SELECT a, b FROM t
< ok 0 0
> exec
| CREATE VIEW v AS SELECT x FROM base
< ok 0 0
> query
| SHOW CREATE TABLE videm_view_on_view_explicit_cols.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_view_on_view_explicit_cols.base
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `base` (`x`,`y`) AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW videm_view_on_view_explicit_cols.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `base`.`x` AS `x` from `base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/weight-string-forms default#1
> exec
| DROP DATABASE IF EXISTS videm_weight_string_forms
< ok 0 0
> exec
| CREATE DATABASE videm_weight_string_forms DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_weight_string_forms
< ok 0 0
> exec
| CREATE TABLE t (s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT WEIGHT_STRING(s) AS c1, WEIGHT_STRING(s AS CHAR(4)) AS c2, WEIGHT_STRING(s AS BINARY(4)) AS c3 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_weight_string_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `s` varchar(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_weight_string_forms.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select weight_string(`t`.`s`) AS `c1`,weight_string(`t`.`s` as char(4)) AS `c2`,weight_string(cast(`t`.`s` as char(4) charset binary)) AS `c3` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewIdempotence/8.0/where-expr default#1
> exec
| DROP DATABASE IF EXISTS videm_where_expr
< ok 0 0
> exec
| CREATE DATABASE videm_where_expr DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_where_expr
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a > 0
< ok 0 0
> query
| SHOW CREATE TABLE videm_where_expr.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW videm_where_expr.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t` where (`t`.`a` > 0)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
