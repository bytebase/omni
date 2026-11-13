-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_cc` AS select `t`.`a` <> 0 and `t`.`b` <> 0 AS `a AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_ac` AS select `t`.`a` + 1 and `t`.`b` <> 0 AS `(a + 1) AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_ca` AS select `t`.`a` > 0 and `t`.`b` + 1 AS `(a > 0) AND (b + 1)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_cmpx2` AS select `t`.`a` > 0 and `t`.`b` > 0 AS `(a > 0) AND (b > 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_abs` AS select abs(`t`.`a`) and `t`.`b` <> 0 AS `ABS(a) AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_case` AS select case when `t`.`a` > 0 then 1 else 0 end and `t`.`b` <> 0 AS `CASE WHEN a > 0 THEN 1 ELSE 0 END AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_if` AS select if(`t`.`a` > 0,1,0) and `t`.`b` <> 0 AS `IF(a > 0, 1, 0) AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_subq` AS select (select max(`t`.`a`) from `t`) and `t`.`b` <> 0 AS `(SELECT MAX(a) FROM t) AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_str` AS select 'hello' and 1 AS `'hello' AND 1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_ifnull` AS select ifnull(`t`.`a`,0) and `t`.`b` <> 0 AS `IFNULL(a, 0) AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_coal` AS select coalesce(`t`.`a`,`t`.`b`) and 1 AS `COALESCE(a, b) AND 1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_nullif` AS select nullif(`t`.`a`,0) and `t`.`b` <> 0 AS `NULLIF(a, 0) AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_great` AS select greatest(`t`.`a`,`t`.`b`) and 1 AS `GREATEST(a, b) AND 1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_least` AS select least(`t`.`a`,`t`.`b`) and 1 AS `LEAST(a, b) AND 1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
