-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_basic_lit` AS select 1 AS `1`,1.5 AS `1.5`,'hello' AS `hello`,NULL AS `NULL`,1 AS `TRUE`,0 AS `FALSE` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_hex_bit` AS select 0xff AS `0xFF`,X'ff' AS `X'FF'`,0x0a AS `0b1010`,0x0a AS `b'1010'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_charset` AS select _utf8mb4'hello' AS `hello`,_latin1'world' AS `world` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_empty_str` AS select '' AS `Name_exp_1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_esc_quotes` AS select 'it\'s' AS `it's` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_esc_bslash` AS select 'back\\slash' AS `back\slash` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_temporal` AS select DATE'2024-01-01' AS `DATE '2024-01-01'`,TIME'12:00:00' AS `TIME '12:00:00'`,TIMESTAMP'2024-01-01 12:00:00' AS `TIMESTAMP '2024-01-01 12:00:00'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_func_args` AS select concat(`t`.`a`,`t`.`b`,`t`.`c`) AS `CONCAT(a, b, c)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_in_list` AS select `t`.`a` in (1,2,3) AS `a IN (1, 2, 3)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
