-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestDeparse_Section_4_1_Container default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_not_gt
< ok 0 0
> exec
| CREATE VIEW v_not_gt AS SELECT NOT (a > 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_not_gt
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_not_gt
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_not_gt` AS select `t`.`a` <= 0 AS `NOT (a > 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_not_lt
< ok 0 0
> exec
| CREATE VIEW v_not_lt AS SELECT NOT (a < 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_not_lt
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_not_lt
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_not_lt` AS select `t`.`a` >= 0 AS `NOT (a < 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_not_ge
< ok 0 0
> exec
| CREATE VIEW v_not_ge AS SELECT NOT (a >= 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_not_ge
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_not_ge
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_not_ge` AS select `t`.`a` < 0 AS `NOT (a >= 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_not_le
< ok 0 0
> exec
| CREATE VIEW v_not_le AS SELECT NOT (a <= 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_not_le
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_not_le
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_not_le` AS select `t`.`a` > 0 AS `NOT (a <= 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_not_eq
< ok 0 0
> exec
| CREATE VIEW v_not_eq AS SELECT NOT (a = 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_not_eq
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_not_eq
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_not_eq` AS select `t`.`a` <> 0 AS `NOT (a = 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_not_ne
< ok 0 0
> exec
| CREATE VIEW v_not_ne AS SELECT NOT (a <> 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_not_ne
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_not_ne
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_not_ne` AS select `t`.`a` = 0 AS `NOT (a <> 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_not_col
< ok 0 0
> exec
| CREATE VIEW v_not_col AS SELECT NOT a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_not_col
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_not_col
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_not_col` AS select `t`.`a` = 0 AS `NOT a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_not_add
< ok 0 0
> exec
| CREATE VIEW v_not_add AS SELECT NOT (a + 1) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_not_add
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_not_add
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_not_add` AS select !(`t`.`a` + 1) AS `NOT (a + 1)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_bang_col
< ok 0 0
> exec
| CREATE VIEW v_bang_col AS SELECT !a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bang_col
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bang_col
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bang_col` AS select `t`.`a` = 0 AS `!a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
