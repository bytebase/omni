-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparseContainer_8_1_ViewOfViewComplexStructures default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v2_vov
< ok 0 0
> exec
| DROP VIEW IF EXISTS v1_vov
< ok 0 0
> exec
| CREATE VIEW v1_vov AS SELECT a FROM t
< ok 0 0
> exec
| CREATE VIEW v2_vov AS SELECT * FROM v1_vov
< ok 0 0
> query
| SHOW CREATE VIEW v2_vov
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2_vov
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v2_vov` (`a`) AS SELECT `test`.`v1_vov`.`a` AS `a` FROM `test`.`v1_vov`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_ten_plus_columns
< ok 0 0
> exec
| CREATE VIEW v_ten_plus_columns AS SELECT a, b, c, a + 1, b + 1, c + 1, a * b, b * c, a * c, a + b + c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_ten_plus_columns
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ten_plus_columns
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ten_plus_columns` (`a`, `b`, `c`, `a + 1`, `b + 1`, `c + 1`, `a * b`, `b * c`, `a * c`, `a + b + c`) AS SELECT `a` AS `a`,`b` AS `b`,`c` AS `c`,`a`+1 AS `a + 1`,`b`+1 AS `b + 1`,`c`+1 AS `c + 1`,`a`*`b` AS `a * b`,`b`*`c` AS `b * c`,`a`*`c` AS `a * c`,`a`+`b`+`c` AS `a + b + c` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_reserved_word_aliases
< ok 0 0
> exec
| CREATE VIEW v_reserved_word_aliases AS SELECT a AS `select`, b AS `from`, c AS `where` FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_reserved_word_aliases
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_reserved_word_aliases
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_reserved_word_aliases` (`select`, `from`, `where`) AS SELECT `a` AS `select`,`b` AS `from`,`c` AS `where` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_case_without_else
< ok 0 0
> exec
| CREATE VIEW v_case_without_else AS SELECT CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' END FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_case_without_else
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_case_without_else
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_case_without_else` (`CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' END`) AS SELECT CASE WHEN `a`>0 THEN _UTF8MB4'pos' WHEN `a`<0 THEN _UTF8MB4'neg' END AS `CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' END` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_between_column_bounds
< ok 0 0
> exec
| CREATE VIEW v_between_column_bounds AS SELECT a BETWEEN b AND c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_between_column_bounds
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_between_column_bounds
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_between_column_bounds` (`a BETWEEN b AND c`) AS SELECT `a` BETWEEN `b` AND `c` AS `a BETWEEN b AND c` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
