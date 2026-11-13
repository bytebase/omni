-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_alter_test
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_alter_test` AS select 2 AS `b`,3 AS `c`
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
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_dep_test
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_dep_test` AS select `t_dep_base`.`id` AS `id`,`t_dep_base`.`val` AS `val` from `t_dep_base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_dep_test
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_show_basic
< ok 0 0
> exec
| CREATE VIEW v_show_basic AS SELECT 1 AS a, 2 AS b
< ok 0 0
> query
| SHOW CREATE VIEW v_show_basic
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_show_basic
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_show_basic` AS select 1 AS `a`,2 AS `b`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_show_opts
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW v_show_opts AS SELECT 1 AS x
< ok 0 0
> query
| SHOW CREATE VIEW v_show_opts
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_show_opts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `v_show_opts` AS select 1 AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_col_aliases
< ok 0 0
> exec
| CREATE VIEW v_col_aliases (x, y, z) AS SELECT 1, 2, 3
< ok 0 0
> query
| SHOW CREATE VIEW v_col_aliases
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_col_aliases
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_col_aliases` (`x`,`y`,`z`) AS select 1 AS `1`,2 AS `2`,3 AS `3`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
