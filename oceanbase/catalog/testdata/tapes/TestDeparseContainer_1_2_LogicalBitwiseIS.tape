-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDeparseContainer_1_2_LogicalBitwiseIS default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_bitwise
< ok 0 0
> exec
| CREATE VIEW v_bitwise AS SELECT a | b, a & b, a ^ b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bitwise
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bitwise
: bytes
| CREATE VIEW `v_bitwise` AS select (`test`.`t`.`a` | `test`.`t`.`b`) AS `a | b`,(`test`.`t`.`a` & `test`.`t`.`b`) AS `a & b`,(`test`.`t`.`a` ^ `test`.`t`.`b`) AS `a ^ b` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_shifts
< ok 0 0
> exec
| CREATE VIEW v_shifts AS SELECT a << b, a >> b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_shifts
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_shifts
: bytes
| CREATE VIEW `v_shifts` AS select (`test`.`t`.`a` << `test`.`t`.`b`) AS `a << b`,(`test`.`t`.`a` >> `test`.`t`.`b`) AS `a >> b` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_bitnot
< ok 0 0
> exec
| CREATE VIEW v_bitnot AS SELECT ~a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bitnot
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bitnot
: bytes
| CREATE VIEW `v_bitnot` AS select ~(`test`.`t`.`a`) AS `~a` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_isnull
< ok 0 0
> exec
| CREATE VIEW v_isnull AS SELECT a IS NULL, a IS NOT NULL FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_isnull
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_isnull
: bytes
| CREATE VIEW `v_isnull` AS select (`test`.`t`.`a` is null) AS `a IS NULL`,(`test`.`t`.`a` is not null) AS `a IS NOT NULL` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_in
< ok 0 0
> exec
| CREATE VIEW v_in AS SELECT a IN (1,2,3), a NOT IN (1,2,3) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_in
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_in
: bytes
| CREATE VIEW `v_in` AS select (`test`.`t`.`a` in (1,2,3)) AS `a IN (1,2,3)`,(`test`.`t`.`a` not in (1,2,3)) AS `a NOT IN (1,2,3)` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_like_esc
< ok 0 0
> exec
| CREATE VIEW v_like_esc AS SELECT a LIKE 'x' ESCAPE '\\' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_like_esc
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_like_esc
: bytes
| CREATE VIEW `v_like_esc` AS select (`test`.`t`.`a` like 'x' escape '\\') AS `a LIKE 'x' ESCAPE '\\'` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
