-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestDeparseContainer_3_3_ComplexPrecedence default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_left_add
< ok 0 0
> exec
| CREATE VIEW v_left_add AS SELECT a + b + c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_left_add
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_left_add
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_left_add` AS select `t`.`a` + `t`.`b` + `t`.`c` AS `a + b + c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_mul_add
< ok 0 0
> exec
| CREATE VIEW v_mul_add AS SELECT a * b + c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_mul_add
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_mul_add
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_mul_add` AS select `t`.`a` * `t`.`b` + `t`.`c` AS `a * b + c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_add_mul
< ok 0 0
> exec
| CREATE VIEW v_add_mul AS SELECT a + b * c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_add_mul
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_add_mul
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_add_mul` AS select `t`.`a` + `t`.`b` * `t`.`c` AS `a + b * c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_paren_mul
< ok 0 0
> exec
| CREATE VIEW v_paren_mul AS SELECT (a + b) * c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_paren_mul
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_paren_mul
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_paren_mul` AS select (`t`.`a` + `t`.`b`) * `t`.`c` AS `(a + b) * c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_or_and
< ok 0 0
> exec
| CREATE VIEW v_or_and AS SELECT a OR b AND c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_or_and
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_or_and
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_or_and` AS select `t`.`a` <> 0 or `t`.`b` <> 0 and `t`.`c` <> 0 AS `a OR b AND c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_paren_or
< ok 0 0
> exec
| CREATE VIEW v_paren_or AS SELECT (a OR b) AND c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_paren_or
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_paren_or
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_paren_or` AS select (`t`.`a` <> 0 or `t`.`b` <> 0) and `t`.`c` <> 0 AS `(a OR b) AND c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_mixed
< ok 0 0
> exec
| CREATE VIEW v_mixed AS SELECT a > 0 AND b < 10 OR c = 5 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_mixed
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_mixed
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_mixed` AS select `t`.`a` > 0 and `t`.`b` < 10 or `t`.`c` = 5 AS `a > 0 AND b < 10 OR c = 5` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
