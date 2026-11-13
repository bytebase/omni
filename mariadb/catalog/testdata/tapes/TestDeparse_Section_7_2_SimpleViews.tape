-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_constant` AS select 1 AS `1`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_column` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_alias` AS select `t`.`a` AS `col1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_multi_columns` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_where` AS select `t`.`a` AS `a` from `t` where `t`.`a` > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select_orderby_limit` AS select `t`.`a` AS `a` from `t` order by `t`.`a` limit 10
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
