-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparseContainer_3_1_BooleanContextWrapping default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_bool_cc
< ok 0 0
> exec
| CREATE VIEW v_bool_cc AS SELECT a AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_cc
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_cc
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_cc` (`a AND b`) AS SELECT `a` AND `b` AS `a AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_ac
< ok 0 0
> exec
| CREATE VIEW v_bool_ac AS SELECT (a + 1) AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_ac
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_ac
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_ac` (`(a + 1) AND b`) AS SELECT (`a`+1) AND `b` AS `(a + 1) AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_ca
< ok 0 0
> exec
| CREATE VIEW v_bool_ca AS SELECT (a > 0) AND (b + 1) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_ca
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_ca
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_ca` (`(a > 0) AND (b + 1)`) AS SELECT (`a`>0) AND (`b`+1) AS `(a > 0) AND (b + 1)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_cmpx2
< ok 0 0
> exec
| CREATE VIEW v_bool_cmpx2 AS SELECT (a > 0) AND (b > 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_cmpx2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_cmpx2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_cmpx2` (`(a > 0) AND (b > 0)`) AS SELECT (`a`>0) AND (`b`>0) AS `(a > 0) AND (b > 0)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_abs
< ok 0 0
> exec
| CREATE VIEW v_bool_abs AS SELECT ABS(a) AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_abs
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_abs
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_abs` (`ABS(a) AND b`) AS SELECT ABS(`a`) AND `b` AS `ABS(a) AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_case
< ok 0 0
> exec
| CREATE VIEW v_bool_case AS SELECT CASE WHEN a > 0 THEN 1 ELSE 0 END AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_case
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_case
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_case` (`CASE WHEN a > 0 THEN 1 ELSE 0 END AND b`) AS SELECT CASE WHEN `a`>0 THEN 1 ELSE 0 END AND `b` AS `CASE WHEN a > 0 THEN 1 ELSE 0 END AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_if
< ok 0 0
> exec
| CREATE VIEW v_bool_if AS SELECT IF(a > 0, 1, 0) AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_if
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_if
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_if` (`IF(a > 0, 1, 0) AND b`) AS SELECT IF(`a`>0, 1, 0) AND `b` AS `IF(a > 0, 1, 0) AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_subq
< ok 0 0
> exec
| CREATE VIEW v_bool_subq AS SELECT (SELECT MAX(a) FROM t) AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_subq
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_subq
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_subq` (`(SELECT MAX(a) FROM t) AND b`) AS SELECT (SELECT MAX(`a`) AS `MAX(a)` FROM `test`.`t`) AND `b` AS `(SELECT MAX(a) FROM t) AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_str
< ok 0 0
> exec
| CREATE VIEW v_bool_str AS SELECT 'hello' AND 1 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_str
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_str
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_str` (`'hello' AND 1`) AS SELECT _UTF8MB4'hello' AND 1 AS `'hello' AND 1` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_ifnull
< ok 0 0
> exec
| CREATE VIEW v_bool_ifnull AS SELECT IFNULL(a, 0) AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_ifnull
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_ifnull
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_ifnull` (`IFNULL(a, 0) AND b`) AS SELECT IFNULL(`a`, 0) AND `b` AS `IFNULL(a, 0) AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_coal
< ok 0 0
> exec
| CREATE VIEW v_bool_coal AS SELECT COALESCE(a, b) AND 1 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_coal
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_coal
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_coal` (`COALESCE(a, b) AND 1`) AS SELECT COALESCE(`a`, `b`) AND 1 AS `COALESCE(a, b) AND 1` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_nullif
< ok 0 0
> exec
| CREATE VIEW v_bool_nullif AS SELECT NULLIF(a, 0) AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_nullif
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_nullif
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_nullif` (`NULLIF(a, 0) AND b`) AS SELECT NULLIF(`a`, 0) AND `b` AS `NULLIF(a, 0) AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_great
< ok 0 0
> exec
| CREATE VIEW v_bool_great AS SELECT GREATEST(a, b) AND 1 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_great
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_great
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_great` (`GREATEST(a, b) AND 1`) AS SELECT GREATEST(`a`, `b`) AND 1 AS `GREATEST(a, b) AND 1` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bool_least
< ok 0 0
> exec
| CREATE VIEW v_bool_least AS SELECT LEAST(a, b) AND 1 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_least
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_least
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_least` (`LEAST(a, b) AND 1`) AS SELECT LEAST(`a`, `b`) AND 1 AS `LEAST(a, b) AND 1` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
