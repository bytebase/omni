-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDeparse_Section_7_2_SimpleViews default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_select_constant
< ok 0 0
> exec
| CREATE VIEW v_select_constant AS SELECT 1
< ok 0 0
> query
| SHOW CREATE VIEW v_select_constant
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_select_constant
: bytes
| CREATE VIEW `v_select_constant` AS select 1 AS `1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
