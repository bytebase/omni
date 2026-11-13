-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_4_5_ViewsDeep default#1
> exec
| DROP VIEW IF EXISTS v_alter_test
< ok 0 0
> exec
| CREATE VIEW v_alter_test AS SELECT 1 AS a
< ok 0 0
> exec
| ALTER VIEW v_alter_test AS SELECT 2 AS b, 3 AS c
< ok 0 0
> query
| SHOW CREATE VIEW v_alter_test
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_alter_test
: bytes
| CREATE VIEW `v_alter_test` AS select 2 AS `b`,3 AS `c`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_alter_noexist
< ok 0 0
> exec
| ALTER VIEW v_alter_noexist AS SELECT 1
< error
| 1146 42S02 Table 'test.v_alter_noexist' doesn't exist
> exec
| DROP VIEW IF EXISTS v_dep_test
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_dep_base
< ok 0 0
> exec
| CREATE TABLE t_dep_base (id INT, val VARCHAR(100))
< ok 0 0
> exec
| CREATE VIEW v_dep_test AS SELECT id, val FROM t_dep_base
< ok 0 0
> exec
| DROP TABLE t_dep_base
< ok 0 0
> query
| SHOW CREATE VIEW v_dep_test
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_dep_test
: bytes
| CREATE VIEW `v_dep_test` AS select `test`.`t_dep_base`.`id` AS `id`,`test`.`t_dep_base`.`val` AS `val` from `test`.`t_dep_base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_dep_test
< ok 0 0
