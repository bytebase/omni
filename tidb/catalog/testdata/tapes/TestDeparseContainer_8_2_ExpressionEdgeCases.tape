-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparseContainer_8_2_ExpressionEdgeCases default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t_rw (`select` INT, b INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS tc (a VARCHAR(50), b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_82_cast_with_expression
< ok 0 0
> exec
| CREATE VIEW v_82_cast_with_expression AS SELECT CAST(a + b AS CHAR) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_82_cast_with_expression
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_cast_with_expression
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_cast_with_expression` (`CAST(a + b AS CHAR)`) AS SELECT CAST(`a`+`b` AS CHAR(20)) AS `CAST(a + b AS CHAR)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_collate_expression
< ok 0 0
> exec
| CREATE VIEW v_82_collate_expression AS SELECT a COLLATE utf8mb4_unicode_ci FROM tc
< ok 0 0
> query
| SHOW CREATE VIEW v_82_collate_expression
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_collate_expression
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_collate_expression` (`a COLLATE utf8mb4_unicode_ci`) AS SELECT `a` COLLATE utf8mb4_unicode_ci AS `a COLLATE utf8mb4_unicode_ci` FROM `test`.`tc`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_interval_arithmetic
< ok 0 0
> exec
| CREATE VIEW v_82_interval_arithmetic AS SELECT INTERVAL 1 DAY + a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_82_interval_arithmetic
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_interval_arithmetic
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_interval_arithmetic` (`INTERVAL 1 DAY + a`) AS SELECT DATE_ADD(`a`, INTERVAL 1 DAY) AS `INTERVAL 1 DAY + a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_unary_operators
< ok 0 0
> exec
| CREATE VIEW v_82_unary_operators AS SELECT -a, +a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_82_unary_operators
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_unary_operators
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_unary_operators` (`-a`, `+a`) AS SELECT -`a` AS `-a`,+`a` AS `+a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_mod_regrouped
< ok 0 0
> exec
| CREATE VIEW v_82_mod_regrouped AS SELECT 2*MOD(a,3), MOD(a+1, 2) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_82_mod_regrouped
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_mod_regrouped
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_mod_regrouped` (`2*MOD(a,3)`, `MOD(a+1, 2)`) AS SELECT 2*`a`%3 AS `2*MOD(a,3)`,`a`+1%2 AS `MOD(a+1, 2)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_group_concat_backslash_separator
< ok 0 0
> exec
| CREATE VIEW v_82_group_concat_backslash_separator AS SELECT GROUP_CONCAT(a SEPARATOR '\\'), GROUP_CONCAT(a SEPARATOR 'x''y') FROM tc
< ok 0 0
> query
| SHOW CREATE VIEW v_82_group_concat_backslash_separator
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_group_concat_backslash_separator
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_group_concat_backslash_separator` (`GROUP_CONCAT(a SEPARATOR '\\')`, `GROUP_CONCAT(a SEPARATOR 'x''y')`) AS SELECT GROUP_CONCAT(`a` SEPARATOR '\\') AS `GROUP_CONCAT(a SEPARATOR '\\')`,GROUP_CONCAT(`a` SEPARATOR 'x''y') AS `GROUP_CONCAT(a SEPARATOR 'x''y')` FROM `test`.`tc`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_long_expression_auto_alias
< ok 0 0
> exec
| CREATE VIEW v_82_long_expression_auto_alias AS SELECT CASE WHEN a > 0 THEN CONCAT(a, b, c) WHEN a < 0 THEN CONCAT(c, b, a) ELSE NULL END FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_82_long_expression_auto_alias
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_long_expression_auto_alias
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_long_expression_auto_alias` (`name_exp_1`) AS SELECT CASE WHEN `a`>0 THEN CONCAT(`a`, `b`, `c`) WHEN `a`<0 THEN CONCAT(`c`, `b`, `a`) ELSE NULL END AS `CASE WHEN a > 0 THEN CONCAT(a, b, c) WHEN a < 0 THEN CONCAT(c, b, a) ELSE NULL END` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_reserved_word_column
< ok 0 0
> exec
| CREATE VIEW v_82_reserved_word_column AS SELECT `select` FROM t_rw
< ok 0 0
> query
| SHOW CREATE VIEW v_82_reserved_word_column
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_reserved_word_column
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_reserved_word_column` (`select`) AS SELECT `select` AS `select` FROM `test`.`t_rw`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_multiple_rewrites
< ok 0 0
> exec
| CREATE VIEW v_82_multiple_rewrites AS SELECT a + b, NOT (a > 0), CAST(a AS CHAR), COUNT(*), a REGEXP 'x' FROM t GROUP BY a, b HAVING COUNT(*) > 1
< ok 0 0
> query
| SHOW CREATE VIEW v_82_multiple_rewrites
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_multiple_rewrites
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_multiple_rewrites` (`a + b`, `NOT (a > 0)`, `CAST(a AS CHAR)`, `COUNT(*)`, `a REGEXP 'x'`) AS SELECT `a`+`b` AS `a + b`,NOT (`a`>0) AS `NOT (a > 0)`,CAST(`a` AS CHAR(11)) AS `CAST(a AS CHAR)`,COUNT(1) AS `COUNT(*)`,`a` REGEXP _UTF8MB4'x' AS `a REGEXP 'x'` FROM `test`.`t` GROUP BY `a`,`b` HAVING COUNT(1)>1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_all_logical_operators
< ok 0 0
> exec
| CREATE VIEW v_82_all_logical_operators AS SELECT a AND b, a OR b, NOT a, a XOR b, !a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_82_all_logical_operators
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_all_logical_operators
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_all_logical_operators` (`a AND b`, `a OR b`, `NOT a`, `a XOR b`, `!a`) AS SELECT `a` AND `b` AS `a AND b`,`a` OR `b` AS `a OR b`,NOT `a` AS `NOT a`,`a` XOR `b` AS `a XOR b`,!`a` AS `!a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_function_boolean_context
< ok 0 0
> exec
| CREATE VIEW v_82_function_boolean_context AS SELECT IFNULL(a, 0) AND COALESCE(b, 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_82_function_boolean_context
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_function_boolean_context
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_function_boolean_context` (`IFNULL(a, 0) AND COALESCE(b, 0)`) AS SELECT IFNULL(`a`, 0) AND COALESCE(`b`, 0) AS `IFNULL(a, 0) AND COALESCE(b, 0)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_mixed_boolean_precedence
< ok 0 0
> exec
| CREATE VIEW v_82_mixed_boolean_precedence AS SELECT (a > 0) AND (b + 1) OR (c IS NULL) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_82_mixed_boolean_precedence
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_mixed_boolean_precedence
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_mixed_boolean_precedence` (`(a > 0) AND (b + 1) OR (c IS NULL)`) AS SELECT (`a`>0) AND (`b`+1) OR (`c` IS NULL) AS `(a > 0) AND (b + 1) OR (c IS NULL)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_sounds_like
< ok 0 0
> exec
| CREATE VIEW v_82_sounds_like AS SELECT a SOUNDS LIKE b FROM tc
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 52 near "LIKE b FROM tc" 
