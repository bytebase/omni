-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDeparseContainer_Section_6_1_SubqueryPatterns default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_scalar_subquery_in_select
< ok 0 0
> exec
| CREATE VIEW v_scalar_subquery_in_select AS SELECT (SELECT MAX(a) FROM t) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_scalar_subquery_in_select
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_scalar_subquery_in_select
: bytes
| CREATE VIEW `v_scalar_subquery_in_select` AS select (select max(`test`.`t`.`a`) from `test`.`t`) AS `(SELECT MAX(a) FROM t)` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
