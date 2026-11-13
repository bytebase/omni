-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestDeparseContainer_3_2_NotFoldingNoDoubleWrap default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_s32_not_fold
< ok 0 0
> exec
| CREATE VIEW v_s32_not_fold AS SELECT NOT (a > 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_s32_not_fold
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_s32_not_fold
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_not_fold` AS select (`t`.`a` <= 0) AS `NOT (a > 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_s32_not_nb
< ok 0 0
> exec
| CREATE VIEW v_s32_not_nb AS SELECT NOT (a + 1) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_s32_not_nb
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_s32_not_nb
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_not_nb` AS select (0 = (`t`.`a` + 1)) AS `NOT (a + 1)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_s32_bang
< ok 0 0
> exec
| CREATE VIEW v_s32_bang AS SELECT !a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_s32_bang
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_s32_bang
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_bang` AS select (0 = `t`.`a`) AS `!a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_s32_cmp_ndw
< ok 0 0
> exec
| CREATE VIEW v_s32_cmp_ndw AS SELECT (a = b) AND (a > 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_s32_cmp_ndw
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_s32_cmp_ndw
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_cmp_ndw` AS select ((`t`.`a` = `t`.`b`) and (`t`.`a` > 0)) AS `(a = b) AND (a > 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_s32_pred_nw
< ok 0 0
> exec
| CREATE VIEW v_s32_pred_nw AS SELECT (a IN (1, 2, 3)) AND (b BETWEEN 1 AND 10) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_s32_pred_nw
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_s32_pred_nw
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_pred_nw` AS select ((`t`.`a` in (1,2,3)) and (`t`.`b` between 1 and 10)) AS `(a IN (1, 2, 3)) AND (b BETWEEN 1 AND 10)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_s32_islike_nw
< ok 0 0
> exec
| CREATE VIEW v_s32_islike_nw AS SELECT (a IS NULL) AND (b LIKE 'x%') FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_s32_islike_nw
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_s32_islike_nw
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_islike_nw` AS select ((`t`.`a` is null) and (`t`.`b` like 'x%')) AS `(a IS NULL) AND (b LIKE 'x%')` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_s32_exists_nw
< ok 0 0
> exec
| CREATE VIEW v_s32_exists_nw AS SELECT EXISTS(SELECT 1 FROM t WHERE a > 0) AND (b > 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_s32_exists_nw
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_s32_exists_nw
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_exists_nw` AS select (exists(select 1 from `t` where (`t`.`a` > 0)) and (`t`.`b` > 0)) AS `EXISTS(SELECT 1 FROM t WHERE a > 0) AND (b > 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
