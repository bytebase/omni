-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_select_constant
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_constant` (`1`) AS SELECT 1 AS `1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_select_column
< ok 0 0
> exec
| CREATE VIEW v_select_column AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_select_column
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_select_column
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_column` (`a`) AS SELECT `a` AS `a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_select_alias
< ok 0 0
> exec
| CREATE VIEW v_select_alias AS SELECT a AS col1 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_select_alias
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_select_alias
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_alias` (`col1`) AS SELECT `a` AS `col1` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_select_multi_columns
< ok 0 0
> exec
| CREATE VIEW v_select_multi_columns AS SELECT a, b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_select_multi_columns
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_select_multi_columns
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_multi_columns` (`a`, `b`) AS SELECT `a` AS `a`,`b` AS `b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_select_where
< ok 0 0
> exec
| CREATE VIEW v_select_where AS SELECT a FROM t WHERE a > 0
< ok 0 0
> query
| SHOW CREATE VIEW v_select_where
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_select_where
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_where` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` WHERE `a`>0
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_select_orderby_limit
< ok 0 0
> exec
| CREATE VIEW v_select_orderby_limit AS SELECT a FROM t ORDER BY a LIMIT 10
< ok 0 0
> query
| SHOW CREATE VIEW v_select_orderby_limit
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_select_orderby_limit
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_orderby_limit` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` ORDER BY `a` LIMIT 10
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
