-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2_vov
: bytes
| CREATE VIEW `v2_vov` AS select `test`.`v1_vov`.`a` AS `a` from `test`.`v1_vov`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ten_plus_columns
: bytes
| CREATE VIEW `v_ten_plus_columns` AS select `test`.`t`.`a` AS `a`,`test`.`t`.`b` AS `b`,`test`.`t`.`c` AS `c`,(`test`.`t`.`a` + 1) AS `a + 1`,(`test`.`t`.`b` + 1) AS `b + 1`,(`test`.`t`.`c` + 1) AS `c + 1`,(`test`.`t`.`a` * `test`.`t`.`b`) AS `a * b`,(`test`.`t`.`b` * `test`.`t`.`c`) AS `b * c`,(`test`.`t`.`a` * `test`.`t`.`c`) AS `a * c`,((`test`.`t`.`a` + `test`.`t`.`b`) + `test`.`t`.`c`) AS `a + b + c` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_reserved_word_aliases
: bytes
| CREATE VIEW `v_reserved_word_aliases` AS select `test`.`t`.`a` AS `select`,`test`.`t`.`b` AS `from`,`test`.`t`.`c` AS `where` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_between_column_bounds
: bytes
| CREATE VIEW `v_between_column_bounds` AS select (`test`.`t`.`a` between `test`.`t`.`b` and `test`.`t`.`c`) AS `a BETWEEN b AND c` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
