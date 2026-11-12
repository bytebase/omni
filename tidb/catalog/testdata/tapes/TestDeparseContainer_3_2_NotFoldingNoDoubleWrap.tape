-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_not_fold` (`NOT (a > 0)`) AS SELECT NOT (`a`>0) AS `NOT (a > 0)` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_not_nb` (`NOT (a + 1)`) AS SELECT NOT (`a`+1) AS `NOT (a + 1)` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_bang` (`!a`) AS SELECT !`a` AS `!a` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_cmp_ndw` (`(a = b) AND (a > 0)`) AS SELECT (`a`=`b`) AND (`a`>0) AS `(a = b) AND (a > 0)` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_pred_nw` (`(a IN (1, 2, 3)) AND (b BETWEEN 1 AND 10)`) AS SELECT (`a` IN (1,2,3)) AND (`b` BETWEEN 1 AND 10) AS `(a IN (1, 2, 3)) AND (b BETWEEN 1 AND 10)` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_islike_nw` (`(a IS NULL) AND (b LIKE 'x%')`) AS SELECT (`a` IS NULL) AND (`b` LIKE _UTF8MB4'x%') AS `(a IS NULL) AND (b LIKE 'x%')` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_s32_exists_nw` (`EXISTS(SELECT 1 FROM t WHERE a > 0) AND (b > 0)`) AS SELECT EXISTS (SELECT 1 AS `1` FROM `test`.`t` WHERE `a`>0) AND (`b`>0) AS `EXISTS(SELECT 1 FROM t WHERE a > 0) AND (b > 0)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
