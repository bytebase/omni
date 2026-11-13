-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDeparseContainer_2_2_RegularFunctionsAggregates default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_agg_funcs
< ok 0 0
> exec
| CREATE VIEW v_agg_funcs AS SELECT SUM(a), AVG(a), MAX(a), MIN(a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_agg_funcs
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_agg_funcs
: bytes
| CREATE VIEW `v_agg_funcs` AS select sum(`test`.`t`.`a`) AS `SUM(a)`,avg(`test`.`t`.`a`) AS `AVG(a)`,max(`test`.`t`.`a`) AS `MAX(a)`,min(`test`.`t`.`a`) AS `MIN(a)` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
