-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDeparseContainer_Section_5_1_SelectClauses default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_limit_offset
< ok 0 0
> exec
| CREATE VIEW v_limit_offset AS SELECT a FROM t LIMIT 10 OFFSET 5
< ok 0 0
> query
| SHOW CREATE VIEW v_limit_offset
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_limit_offset
: bytes
| CREATE VIEW `v_limit_offset` AS select `test`.`t`.`a` AS `a` from `test`.`t` limit 5,10
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_expr_groupby
< ok 0 0
> exec
| CREATE VIEW v_expr_groupby AS SELECT a + b FROM t GROUP BY a + b
< ok 0 0
> query
| SHOW CREATE VIEW v_expr_groupby
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_expr_groupby
: bytes
| CREATE VIEW `v_expr_groupby` AS select (`test`.`t`.`a` + `test`.`t`.`b`) AS `a + b` from `test`.`t` group by (`test`.`t`.`a` + `test`.`t`.`b`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
