-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDeparse_Section_7_3_ExpressionViews default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_arithmetic_expr
< ok 0 0
> exec
| CREATE VIEW v_arithmetic_expr AS SELECT a + b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_arithmetic_expr
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_arithmetic_expr
: bytes
| CREATE VIEW `v_arithmetic_expr` AS select (`test`.`t`.`a` + `test`.`t`.`b`) AS `a + b` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_case_expr
< ok 0 0
> exec
| CREATE VIEW v_case_expr AS SELECT CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_case_expr
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_case_expr
: bytes
| CREATE VIEW `v_case_expr` AS select (case when (`test`.`t`.`a` > 0) then 'pos' else 'neg' end) AS `CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cast_expr
< ok 0 0
> exec
| CREATE VIEW v_cast_expr AS SELECT CAST(a AS CHAR) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cast_expr
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_expr
: bytes
| CREATE VIEW `v_cast_expr` AS select cast(`test`.`t`.`a` as char charset utf8mb4) AS `CAST(a AS CHAR)` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
