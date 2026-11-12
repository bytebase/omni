-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_add` (`a + b`) AS SELECT `a`+`b` AS `a + b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_sub` (`a - b`) AS SELECT `a`-`b` AS `a - b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_mul` (`a * b`) AS SELECT `a`*`b` AS `a * b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_div` (`a / b`) AS SELECT `a`/`b` AS `a / b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_intdiv` (`a DIV b`) AS SELECT `a` DIV `b` AS `a DIV b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_mod` (`a MOD b`) AS SELECT `a`%`b` AS `a MOD b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_eq` (`a = b`) AS SELECT `a`=`b` AS `a = b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_neq` (`a != b`) AS SELECT `a`!=`b` AS `a != b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_neq2` (`a <> b`) AS SELECT `a`!=`b` AS `a <> b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cmp` (`a > b`, `a < b`, `a >= b`, `a <= b`) AS SELECT `a`>`b` AS `a > b`,`a`<`b` AS `a < b`,`a`>=`b` AS `a >= b`,`a`<=`b` AS `a <= b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nseq` (`a <=> b`) AS SELECT `a`<=>`b` AS `a <=> b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
