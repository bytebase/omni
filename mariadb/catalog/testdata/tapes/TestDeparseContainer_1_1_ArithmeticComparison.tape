-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestDeparseContainer_1_1_ArithmeticComparison default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_add
< ok 0 0
> exec
| CREATE VIEW v_add AS SELECT a + b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_add
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_add
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_add` AS select `t`.`a` + `t`.`b` AS `a + b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_sub
< ok 0 0
> exec
| CREATE VIEW v_sub AS SELECT a - b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_sub
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_sub
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_sub` AS select `t`.`a` - `t`.`b` AS `a - b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_mul
< ok 0 0
> exec
| CREATE VIEW v_mul AS SELECT a * b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_mul
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_mul
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_mul` AS select `t`.`a` * `t`.`b` AS `a * b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_div
< ok 0 0
> exec
| CREATE VIEW v_div AS SELECT a / b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_div
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_div
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_div` AS select `t`.`a` / `t`.`b` AS `a / b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_intdiv
< ok 0 0
> exec
| CREATE VIEW v_intdiv AS SELECT a DIV b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_intdiv
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_intdiv
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_intdiv` AS select `t`.`a` DIV `t`.`b` AS `a DIV b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_mod
< ok 0 0
> exec
| CREATE VIEW v_mod AS SELECT a MOD b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_mod
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_mod
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_mod` AS select `t`.`a` MOD `t`.`b` AS `a MOD b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_eq
< ok 0 0
> exec
| CREATE VIEW v_eq AS SELECT a = b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_eq
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_eq
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_eq` AS select `t`.`a` = `t`.`b` AS `a = b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_neq
< ok 0 0
> exec
| CREATE VIEW v_neq AS SELECT a != b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_neq
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_neq
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_neq` AS select `t`.`a` <> `t`.`b` AS `a != b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_neq2
< ok 0 0
> exec
| CREATE VIEW v_neq2 AS SELECT a <> b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_neq2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_neq2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_neq2` AS select `t`.`a` <> `t`.`b` AS `a <> b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_cmp
< ok 0 0
> exec
| CREATE VIEW v_cmp AS SELECT a > b, a < b, a >= b, a <= b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cmp
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cmp
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cmp` AS select `t`.`a` > `t`.`b` AS `a > b`,`t`.`a` < `t`.`b` AS `a < b`,`t`.`a` >= `t`.`b` AS `a >= b`,`t`.`a` <= `t`.`b` AS `a <= b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_nseq
< ok 0 0
> exec
| CREATE VIEW v_nseq AS SELECT a <=> b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_nseq
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nseq
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nseq` AS select `t`.`a` <=> `t`.`b` AS `a <=> b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
