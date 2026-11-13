-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDeparseContainer_4_2_MultiTableDerived default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t2 (a INT, b INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t3 (b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_derived
< ok 0 0
> exec
| CREATE VIEW v_derived AS SELECT d.x FROM (SELECT a AS x FROM t) d
< ok 0 0
> query
| SHOW CREATE VIEW v_derived
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_derived
: bytes
| CREATE VIEW `v_derived` AS select `d`.`x` AS `x` from (select `test`.`t`.`a` AS `x` from `test`.`t`) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_derived_where
< ok 0 0
> exec
| CREATE VIEW v_derived_where AS SELECT d.x FROM (SELECT a AS x FROM t WHERE a > 0) AS d WHERE d.x < 10
< ok 0 0
> query
| SHOW CREATE VIEW v_derived_where
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_derived_where
: bytes
| CREATE VIEW `v_derived_where` AS select `d`.`x` AS `x` from (select `test`.`t`.`a` AS `x` from `test`.`t` where (`test`.`t`.`a` > 0)) `d` where (`d`.`x` < 10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_alias_as
< ok 0 0
> exec
| CREATE VIEW v_alias_as AS SELECT x.a FROM t AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_alias_as
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_alias_as
: bytes
| CREATE VIEW `v_alias_as` AS select `x`.`a` AS `a` from `test`.`t` `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_alias_noas
< ok 0 0
> exec
| CREATE VIEW v_alias_noas AS SELECT x.a FROM t x
< ok 0 0
> query
| SHOW CREATE VIEW v_alias_noas
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_alias_noas
: bytes
| CREATE VIEW `v_alias_noas` AS select `x`.`a` AS `a` from `test`.`t` `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
