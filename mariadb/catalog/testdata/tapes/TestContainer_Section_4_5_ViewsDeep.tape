-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_alter_noexist
< ok 0 0
> exec
| ALTER VIEW v_alter_noexist AS SELECT 1
< error
| 1146 42S02 Table 'catalog_shared.v_alter_noexist' doesn't exist
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_dep_test` AS select `catalog_shared`.`t_dep_base`.`id` AS `id`,`catalog_shared`.`t_dep_base`.`val` AS `val` from `t_dep_base`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| utf8mb4_uca1400_ai_ci
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
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_col_aliases` AS select 1 AS `x`,2 AS `y`,3 AS `z`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
