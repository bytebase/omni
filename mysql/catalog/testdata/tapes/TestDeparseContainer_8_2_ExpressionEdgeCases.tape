-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_cast_with_expression` AS select cast((`t`.`a` + `t`.`b`) as char charset utf8mb4) AS `CAST(a + b AS CHAR)` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_collate_expression` AS select (`tc`.`a` collate utf8mb4_unicode_ci) AS `a COLLATE utf8mb4_unicode_ci` from `tc`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_interval_arithmetic` AS select (`t`.`a` + interval 1 day) AS `INTERVAL 1 DAY + a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_82_sounds_like
< ok 0 0
> exec
| CREATE VIEW v_82_sounds_like AS SELECT a SOUNDS LIKE b FROM tc
< ok 0 0
> query
| SHOW CREATE VIEW v_82_sounds_like
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_sounds_like
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_sounds_like` AS select (soundex(`tc`.`a`) = soundex(`tc`.`b`)) AS `a SOUNDS LIKE b` from `tc`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_unary_operators` AS select -(`t`.`a`) AS `-a`,`t`.`a` AS `a` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_long_expression_auto_alias` AS select (case when (`t`.`a` > 0) then concat(`t`.`a`,`t`.`b`,`t`.`c`) when (`t`.`a` < 0) then concat(`t`.`c`,`t`.`b`,`t`.`a`) else NULL end) AS `Name_exp_1` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_reserved_word_column` AS select `t_rw`.`select` AS `select` from `t_rw`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_multiple_rewrites` AS select (`t`.`a` + `t`.`b`) AS `a + b`,(`t`.`a` <= 0) AS `NOT (a > 0)`,cast(`t`.`a` as char charset utf8mb4) AS `CAST(a AS CHAR)`,count(0) AS `COUNT(*)`,regexp_like(`t`.`a`,'x') AS `a REGEXP 'x'` from `t` group by `t`.`a`,`t`.`b` having (count(0) > 1)
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_all_logical_operators` AS select ((0 <> `t`.`a`) and (0 <> `t`.`b`)) AS `a AND b`,((0 <> `t`.`a`) or (0 <> `t`.`b`)) AS `a OR b`,(0 = `t`.`a`) AS `NOT a`,((0 <> `t`.`a`) xor (0 <> `t`.`b`)) AS `a XOR b`,(0 = `t`.`a`) AS `!a` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_function_boolean_context` AS select ((0 <> ifnull(`t`.`a`,0)) and (0 <> coalesce(`t`.`b`,0))) AS `IFNULL(a, 0) AND COALESCE(b, 0)` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_82_mixed_boolean_precedence` AS select (((`t`.`a` > 0) and (0 <> (`t`.`b` + 1))) or (`t`.`c` is null)) AS `(a > 0) AND (b + 1) OR (c IS NULL)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
