-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDeparse_Section_7_6_Regression default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_reg_simple_select_column
< ok 0 0
> exec
| CREATE VIEW v_reg_simple_select_column AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_reg_simple_select_column
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_reg_simple_select_column
: bytes
| CREATE VIEW `v_reg_simple_select_column` AS select `test`.`t`.`a` AS `a` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_reg_simple_select_where
< ok 0 0
> exec
| CREATE VIEW v_reg_simple_select_where AS SELECT a FROM t WHERE a > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_reg_simple_select_where
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_reg_simple_select_where
: bytes
| CREATE VIEW `v_reg_simple_select_where` AS select `test`.`t`.`a` AS `a` from `test`.`t` where (`test`.`t`.`a` > 0)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
