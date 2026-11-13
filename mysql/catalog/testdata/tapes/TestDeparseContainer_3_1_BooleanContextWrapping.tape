-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_cc` AS select ((0 <> `t`.`a`) and (0 <> `t`.`b`)) AS `a AND b` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_ac` AS select ((0 <> (`t`.`a` + 1)) and (0 <> `t`.`b`)) AS `(a + 1) AND b` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_ca` AS select ((`t`.`a` > 0) and (0 <> (`t`.`b` + 1))) AS `(a > 0) AND (b + 1)` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_cmpx2` AS select ((`t`.`a` > 0) and (`t`.`b` > 0)) AS `(a > 0) AND (b > 0)` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_abs` AS select ((0 <> abs(`t`.`a`)) and (0 <> `t`.`b`)) AS `ABS(a) AND b` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_case` AS select ((0 <> (case when (`t`.`a` > 0) then 1 else 0 end)) and (0 <> `t`.`b`)) AS `CASE WHEN a > 0 THEN 1 ELSE 0 END AND b` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_if` AS select ((0 <> if((`t`.`a` > 0),1,0)) and (0 <> `t`.`b`)) AS `IF(a > 0, 1, 0) AND b` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_subq` AS select ((0 <> (select max(`t`.`a`) from `t`)) and (0 <> `t`.`b`)) AS `(SELECT MAX(a) FROM t) AND b` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_str` AS select ((0 <> 'hello') and (0 <> 1)) AS `'hello' AND 1` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_ifnull` AS select ((0 <> ifnull(`t`.`a`,0)) and (0 <> `t`.`b`)) AS `IFNULL(a, 0) AND b` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_coal` AS select ((0 <> coalesce(`t`.`a`,`t`.`b`)) and (0 <> 1)) AS `COALESCE(a, b) AND 1` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_nullif` AS select ((0 <> nullif(`t`.`a`,0)) and (0 <> `t`.`b`)) AS `NULLIF(a, 0) AND b` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_great` AS select ((0 <> greatest(`t`.`a`,`t`.`b`)) and (0 <> 1)) AS `GREATEST(a, b) AND 1` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bool_least` AS select ((0 <> least(`t`.`a`,`t`.`b`)) and (0 <> 1)) AS `LEAST(a, b) AND 1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
