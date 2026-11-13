-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestDeparseContainer_2_3_SpecialFunctions default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS tv (a VARCHAR(50), b VARCHAR(50), c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_trim_simple
< ok 0 0
> exec
| CREATE VIEW v_trim_simple AS SELECT TRIM(a) FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_simple
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_simple
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_simple` AS select trim(`tv`.`a`) AS `TRIM(a)` from `tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_trim_leading
< ok 0 0
> exec
| CREATE VIEW v_trim_leading AS SELECT TRIM(LEADING 'x' FROM a) FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_leading
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_leading
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_leading` AS select trim(leading 'x' from `tv`.`a`) AS `TRIM(LEADING 'x' FROM a)` from `tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_trim_trailing
< ok 0 0
> exec
| CREATE VIEW v_trim_trailing AS SELECT TRIM(TRAILING 'x' FROM a) FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_trailing
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_trailing
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_trailing` AS select trim(trailing 'x' from `tv`.`a`) AS `TRIM(TRAILING 'x' FROM a)` from `tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_trim_both
< ok 0 0
> exec
| CREATE VIEW v_trim_both AS SELECT TRIM(BOTH 'x' FROM a) FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_both
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_both
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_both` AS select trim(both 'x' from `tv`.`a`) AS `TRIM(BOTH 'x' FROM a)` from `tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_gc_basic
< ok 0 0
> exec
| CREATE VIEW v_gc_basic AS SELECT GROUP_CONCAT(a ORDER BY a SEPARATOR ',') FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_gc_basic
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_gc_basic
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_gc_basic` AS select group_concat(`t`.`a` order by `t`.`a` ASC separator ',') AS `GROUP_CONCAT(a ORDER BY a SEPARATOR ',')` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_gc_full
< ok 0 0
> exec
| CREATE VIEW v_gc_full AS SELECT GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';') FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_gc_full
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_gc_full
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_gc_full` AS select group_concat(distinct `t`.`a` order by `t`.`a` DESC separator ';') AS `GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';')` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_simple_case
< ok 0 0
> exec
| CREATE VIEW v_simple_case AS SELECT CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_simple_case
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_simple_case
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_simple_case` AS select (case `t`.`a` when 1 then 'one' when 2 then 'two' else 'other' end) AS `CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
