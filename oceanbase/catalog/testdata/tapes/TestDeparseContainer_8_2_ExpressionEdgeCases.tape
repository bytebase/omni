-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_cast_with_expression
: bytes
| CREATE VIEW `v_82_cast_with_expression` AS select cast((`test`.`t`.`a` + `test`.`t`.`b`) as char charset utf8mb4) AS `CAST(a + b AS CHAR)` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_interval_arithmetic
: bytes
| CREATE VIEW `v_82_interval_arithmetic` AS select (`test`.`t`.`a` + interval 1 day) AS `INTERVAL 1 DAY + a` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_sounds_like
: bytes
| CREATE VIEW `v_82_sounds_like` AS select (soundex(`test`.`tc`.`a`) = soundex(`test`.`tc`.`b`)) AS `a SOUNDS LIKE b` from `test`.`tc`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_unary_operators
: bytes
| CREATE VIEW `v_82_unary_operators` AS select -(`test`.`t`.`a`) AS `-a`,`test`.`t`.`a` AS `a` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_82_reserved_word_column
: bytes
| CREATE VIEW `v_82_reserved_word_column` AS select `test`.`t_rw`.`select` AS `select` from `test`.`t_rw`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
