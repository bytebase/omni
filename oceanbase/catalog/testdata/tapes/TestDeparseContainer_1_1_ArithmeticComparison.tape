-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_add
: bytes
| CREATE VIEW `v_add` AS select (`test`.`t`.`a` + `test`.`t`.`b`) AS `a + b` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_sub
: bytes
| CREATE VIEW `v_sub` AS select (`test`.`t`.`a` - `test`.`t`.`b`) AS `a - b` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_mul
: bytes
| CREATE VIEW `v_mul` AS select (`test`.`t`.`a` * `test`.`t`.`b`) AS `a * b` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_div
: bytes
| CREATE VIEW `v_div` AS select (`test`.`t`.`a` / `test`.`t`.`b`) AS `a / b` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_mod
: bytes
| CREATE VIEW `v_mod` AS select (`test`.`t`.`a` % `test`.`t`.`b`) AS `a MOD b` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_eq
: bytes
| CREATE VIEW `v_eq` AS select (`test`.`t`.`a` = `test`.`t`.`b`) AS `a = b` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_neq
: bytes
| CREATE VIEW `v_neq` AS select (`test`.`t`.`a` <> `test`.`t`.`b`) AS `a != b` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_neq2
: bytes
| CREATE VIEW `v_neq2` AS select (`test`.`t`.`a` <> `test`.`t`.`b`) AS `a <> b` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cmp
: bytes
| CREATE VIEW `v_cmp` AS select (`test`.`t`.`a` > `test`.`t`.`b`) AS `a > b`,(`test`.`t`.`a` < `test`.`t`.`b`) AS `a < b`,(`test`.`t`.`a` >= `test`.`t`.`b`) AS `a >= b`,(`test`.`t`.`a` <= `test`.`t`.`b`) AS `a <= b` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nseq
: bytes
| CREATE VIEW `v_nseq` AS select (`test`.`t`.`a` <=> `test`.`t`.`b`) AS `a <=> b` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
