-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_char
: bytes
| CREATE VIEW `v_cast_char` AS select cast(`test`.`t`.`a` as char charset utf8mb4) AS `CAST(a AS CHAR)` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_char10
: bytes
| CREATE VIEW `v_cast_char10` AS select cast(`test`.`t`.`a` as char(10) charset utf8mb4) AS `CAST(a AS CHAR(10))` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_binary
: bytes
| CREATE VIEW `v_cast_binary` AS select cast(`test`.`t`.`a` as char charset binary) AS `CAST(a AS BINARY)` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_su
: bytes
| CREATE VIEW `v_cast_su` AS select cast(`test`.`t`.`a` as signed) AS `CAST(a AS SIGNED)`,cast(`test`.`t`.`a` as unsigned) AS `CAST(a AS UNSIGNED)` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_conv_char
: bytes
| CREATE VIEW `v_conv_char` AS select cast(`test`.`t`.`a` as char charset utf8mb4) AS `CONVERT(a, CHAR)` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
