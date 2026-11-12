-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_char` (`CAST(a AS CHAR)`) AS SELECT CAST(`a` AS CHAR(11)) AS `CAST(a AS CHAR)` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_char10` (`CAST(a AS CHAR(10))`) AS SELECT CAST(`a` AS CHAR(10)) AS `CAST(a AS CHAR(10))` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_binary` (`CAST(a AS BINARY)`) AS SELECT CAST(`a` AS BINARY(11)) AS `CAST(a AS BINARY)` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_su` (`CAST(a AS SIGNED)`, `CAST(a AS UNSIGNED)`) AS SELECT CAST(`a` AS SIGNED) AS `CAST(a AS SIGNED)`,CAST(`a` AS UNSIGNED) AS `CAST(a AS UNSIGNED)` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_dec` (`CAST(a AS DECIMAL(10,2))`) AS SELECT CAST(`a` AS DECIMAL(10, 2)) AS `CAST(a AS DECIMAL(10,2))` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_ddj` (`CAST(a AS DATE)`, `CAST(a AS DATETIME)`, `CAST(a AS JSON)`) AS SELECT CAST(`a` AS DATE) AS `CAST(a AS DATE)`,CAST(`a` AS DATETIME) AS `CAST(a AS DATETIME)`,CAST(`a` AS JSON) AS `CAST(a AS JSON)` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_conv_char` (`CONVERT(a, CHAR)`) AS SELECT CONVERT(`a`, CHAR(11)) AS `CONVERT(a, CHAR)` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_conv_using` (`CONVERT(a USING utf8mb4)`) AS SELECT CONVERT(`a` USING 'utf8mb4') AS `CONVERT(a USING utf8mb4)` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_regexp` (`a REGEXP 'pattern'`) AS SELECT `a` REGEXP _UTF8MB4'pattern' AS `a REGEXP 'pattern'` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_not_regexp` (`a NOT REGEXP 'pattern'`) AS SELECT `a` NOT REGEXP _UTF8MB4'pattern' AS `a NOT REGEXP 'pattern'` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_json_ext` (`a->'$.key'`) AS SELECT JSON_EXTRACT(`a`, _UTF8MB4'$.key') AS `a->'$.key'` FROM `test`.`tj`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_json_unq` (`a->>'$.key'`) AS SELECT JSON_UNQUOTE(JSON_EXTRACT(`a`, _UTF8MB4'$.key')) AS `a->>'$.key'` FROM `test`.`tj`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
