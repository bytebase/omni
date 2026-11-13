-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_left_add
: bytes
| CREATE VIEW `v_left_add` AS select ((`test`.`t`.`a` + `test`.`t`.`b`) + `test`.`t`.`c`) AS `a + b + c` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_mul_add
: bytes
| CREATE VIEW `v_mul_add` AS select ((`test`.`t`.`a` * `test`.`t`.`b`) + `test`.`t`.`c`) AS `a * b + c` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_add_mul
: bytes
| CREATE VIEW `v_add_mul` AS select (`test`.`t`.`a` + (`test`.`t`.`b` * `test`.`t`.`c`)) AS `a + b * c` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_paren_mul
: bytes
| CREATE VIEW `v_paren_mul` AS select ((`test`.`t`.`a` + `test`.`t`.`b`) * `test`.`t`.`c`) AS `(a + b) * c` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_mixed
: bytes
| CREATE VIEW `v_mixed` AS select (((`test`.`t`.`a` > 0) and (`test`.`t`.`b` < 10)) or (`test`.`t`.`c` = 5)) AS `a > 0 AND b < 10 OR c = 5` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
