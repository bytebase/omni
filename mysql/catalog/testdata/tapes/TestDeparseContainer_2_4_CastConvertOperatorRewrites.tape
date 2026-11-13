-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestDeparseContainer_2_4_CastConvertOperatorRewrites default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS tj (a JSON, b INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_cast_char
< ok 0 0
> exec
| CREATE VIEW v_cast_char AS SELECT CAST(a AS CHAR) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cast_char
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_char
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_char` AS select cast(`t`.`a` as char charset utf8mb4) AS `CAST(a AS CHAR)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cast_char10
< ok 0 0
> exec
| CREATE VIEW v_cast_char10 AS SELECT CAST(a AS CHAR(10)) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cast_char10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_char10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_char10` AS select cast(`t`.`a` as char(10) charset utf8mb4) AS `CAST(a AS CHAR(10))` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cast_binary
< ok 0 0
> exec
| CREATE VIEW v_cast_binary AS SELECT CAST(a AS BINARY) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cast_binary
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_binary
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_binary` AS select cast(`t`.`a` as char charset binary) AS `CAST(a AS BINARY)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cast_su
< ok 0 0
> exec
| CREATE VIEW v_cast_su AS SELECT CAST(a AS SIGNED), CAST(a AS UNSIGNED) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cast_su
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_su
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_su` AS select cast(`t`.`a` as signed) AS `CAST(a AS SIGNED)`,cast(`t`.`a` as unsigned) AS `CAST(a AS UNSIGNED)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cast_dec
< ok 0 0
> exec
| CREATE VIEW v_cast_dec AS SELECT CAST(a AS DECIMAL(10,2)) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cast_dec
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_dec
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_dec` AS select cast(`t`.`a` as decimal(10,2)) AS `CAST(a AS DECIMAL(10,2))` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cast_ddj
< ok 0 0
> exec
| CREATE VIEW v_cast_ddj AS SELECT CAST(a AS DATE), CAST(a AS DATETIME), CAST(a AS JSON) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cast_ddj
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_ddj
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_ddj` AS select cast(`t`.`a` as date) AS `CAST(a AS DATE)`,cast(`t`.`a` as datetime) AS `CAST(a AS DATETIME)`,cast(`t`.`a` as json) AS `CAST(a AS JSON)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_conv_char
< ok 0 0
> exec
| CREATE VIEW v_conv_char AS SELECT CONVERT(a, CHAR) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_conv_char
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_conv_char
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_conv_char` AS select cast(`t`.`a` as char charset utf8mb4) AS `CONVERT(a, CHAR)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_conv_using
< ok 0 0
> exec
| CREATE VIEW v_conv_using AS SELECT CONVERT(a USING utf8mb4) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_conv_using
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_conv_using
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_conv_using` AS select convert(`t`.`a` using utf8mb4) AS `CONVERT(a USING utf8mb4)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_regexp
< ok 0 0
> exec
| CREATE VIEW v_regexp AS SELECT a REGEXP 'pattern' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_regexp
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_regexp
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_regexp` AS select regexp_like(`t`.`a`,'pattern') AS `a REGEXP 'pattern'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_not_regexp
< ok 0 0
> exec
| CREATE VIEW v_not_regexp AS SELECT a NOT REGEXP 'pattern' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_not_regexp
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_not_regexp
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_not_regexp` AS select (not(regexp_like(`t`.`a`,'pattern'))) AS `a NOT REGEXP 'pattern'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_json_ext
< ok 0 0
> exec
| CREATE VIEW v_json_ext AS SELECT a->'$.key' FROM tj
< ok 0 0
> query
| SHOW CREATE VIEW v_json_ext
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_json_ext
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_json_ext` AS select json_extract(`tj`.`a`,'$.key') AS `a->'$.key'` from `tj`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_json_unq
< ok 0 0
> exec
| CREATE VIEW v_json_unq AS SELECT a->>'$.key' FROM tj
< ok 0 0
> query
| SHOW CREATE VIEW v_json_unq
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_json_unq
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_json_unq` AS select json_unquote(json_extract(`tj`.`a`,'$.key')) AS `a->>'$.key'` from `tj`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
