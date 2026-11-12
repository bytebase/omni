-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_4_5_ViewsDeep default#1
> exec
| DROP VIEW IF EXISTS v_alter_test
< ok 0 0
> exec
| CREATE VIEW v_alter_test AS SELECT 1 AS a
< ok 0 0
> exec
| ALTER VIEW v_alter_test AS SELECT 2 AS b, 3 AS c
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 10 near "VIEW v_alter_test AS SELECT 2 AS b, 3 AS c" 
> exec
| CREATE OR REPLACE VIEW v_alter_test AS SELECT 2 AS b, 3 AS c
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_alter_test` (`b`, `c`) AS SELECT 2 AS `b`,3 AS `c`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_dep_test` (`id`, `val`) AS SELECT `id` AS `id`,`val` AS `val` FROM `test`.`t_dep_base`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_show_basic` (`a`, `b`) AS SELECT 1 AS `a`,2 AS `b`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `v_show_opts` (`x`) AS SELECT 1 AS `x`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_col_aliases` (`x`, `y`, `z`) AS SELECT 1 AS `1`,2 AS `2`,3 AS `3`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
