-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparseContainer_1_3_LiteralsSpacing default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_basic_lit
< ok 0 0
> exec
| CREATE VIEW v_basic_lit AS SELECT 1, 1.5, 'hello', NULL, TRUE, FALSE FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_basic_lit
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_basic_lit
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_basic_lit` (`1`, `1.5`, `hello`, `NULL`, `TRUE`, `FALSE`) AS SELECT 1 AS `1`,1.5 AS `1.5`,_UTF8MB4'hello' AS `hello`,NULL AS `NULL`,TRUE AS `TRUE`,FALSE AS `FALSE` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_hex_bit
< ok 0 0
> exec
| CREATE VIEW v_hex_bit AS SELECT 0xFF, X'FF', 0b1010, b'1010' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_hex_bit
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_hex_bit
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_hex_bit` (`0xFF`, `X'FF'`, `0b1010`, `b'1010'`) AS SELECT x'ff' AS `0xFF`,x'ff' AS `X'FF'`,b'1010' AS `0b1010`,b'1010' AS `b'1010'` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_charset
< ok 0 0
> exec
| CREATE VIEW v_charset AS SELECT _utf8mb4'hello', _latin1'world' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_charset
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_charset
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_charset` (`hello`, `world`) AS SELECT _UTF8MB4'hello' AS `hello`,_LATIN1'world' AS `world` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_empty_str
< ok 0 0
> exec
| CREATE VIEW v_empty_str AS SELECT '' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_empty_str
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_empty_str
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_empty_str` (``) AS SELECT _UTF8MB4'' FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_esc_quotes
< ok 0 0
> exec
| CREATE VIEW v_esc_quotes AS SELECT 'it''s' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_esc_quotes
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_esc_quotes
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_esc_quotes` (`it's`) AS SELECT _UTF8MB4'it''s' AS `it's` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_esc_bslash
< ok 0 0
> exec
| CREATE VIEW v_esc_bslash AS SELECT 'back\\slash' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_esc_bslash
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_esc_bslash
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_esc_bslash` (`back\slash`) AS SELECT _UTF8MB4'back\\slash' AS `back\slash` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_temporal
< ok 0 0
> exec
| CREATE VIEW v_temporal AS SELECT DATE '2024-01-01', TIME '12:00:00', TIMESTAMP '2024-01-01 12:00:00' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_temporal
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_temporal
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_temporal` (`DATE '2024-01-01'`, `TIME '12:00:00'`, `TIMESTAMP '2024-01-01 12:00:00'`) AS SELECT DATE '2024-01-01' AS `DATE '2024-01-01'`,TIME '12:00:00' AS `TIME '12:00:00'`,TIMESTAMP '2024-01-01 12:00:00' AS `TIMESTAMP '2024-01-01 12:00:00'` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_func_args
< ok 0 0
> exec
| CREATE VIEW v_func_args AS SELECT CONCAT(a, b, c) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_func_args
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_func_args
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_func_args` (`CONCAT(a, b, c)`) AS SELECT CONCAT(`a`, `b`, `c`) AS `CONCAT(a, b, c)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_in_list
< ok 0 0
> exec
| CREATE VIEW v_in_list AS SELECT a IN (1, 2, 3) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_in_list
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_in_list
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_in_list` (`a IN (1, 2, 3)`) AS SELECT `a` IN (1,2,3) AS `a IN (1, 2, 3)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
