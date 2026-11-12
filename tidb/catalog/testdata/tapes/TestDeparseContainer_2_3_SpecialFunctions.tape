-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_simple` (`TRIM(a)`) AS SELECT TRIM(`a`) AS `TRIM(a)` FROM `test`.`tv`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_leading` (`TRIM(LEADING 'x' FROM a)`) AS SELECT TRIM(LEADING _UTF8MB4'x' FROM `a`) AS `TRIM(LEADING 'x' FROM a)` FROM `test`.`tv`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_trailing` (`TRIM(TRAILING 'x' FROM a)`) AS SELECT TRIM(TRAILING _UTF8MB4'x' FROM `a`) AS `TRIM(TRAILING 'x' FROM a)` FROM `test`.`tv`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_both` (`TRIM(BOTH 'x' FROM a)`) AS SELECT TRIM(BOTH _UTF8MB4'x' FROM `a`) AS `TRIM(BOTH 'x' FROM a)` FROM `test`.`tv`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_gc_basic` (`GROUP_CONCAT(a ORDER BY a SEPARATOR ',')`) AS SELECT GROUP_CONCAT(`a` ORDER BY `a` SEPARATOR ',') AS `GROUP_CONCAT(a ORDER BY a SEPARATOR ',')` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_gc_full` (`GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';')`) AS SELECT GROUP_CONCAT(DISTINCT `a` ORDER BY `a` DESC SEPARATOR ';') AS `GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';')` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_simple_case` (`CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END`) AS SELECT CASE `a` WHEN 1 THEN _UTF8MB4'one' WHEN 2 THEN _UTF8MB4'two' ELSE _UTF8MB4'other' END AS `CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
