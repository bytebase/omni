-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDeparseContainer_3_1_BooleanContextWrapping default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_bool_cmpx2
< ok 0 0
> exec
| CREATE VIEW v_bool_cmpx2 AS SELECT (a > 0) AND (b > 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bool_cmpx2
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bool_cmpx2
: bytes
| CREATE VIEW `v_bool_cmpx2` AS select ((`test`.`t`.`a` > 0) and (`test`.`t`.`b` > 0)) AS `(a > 0) AND (b > 0)` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
