-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDeparseContainer_2_3_SpecialFunctions default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS tv (a VARCHAR(50), b VARCHAR(50), c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_trim_simple
< ok 0 0
> exec
| CREATE VIEW v_trim_simple AS SELECT TRIM(a) FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_simple
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_simple
: bytes
| CREATE VIEW `v_trim_simple` AS select trim(`test`.`tv`.`a`) AS `TRIM(a)` from `test`.`tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_trim_leading
< ok 0 0
> exec
| CREATE VIEW v_trim_leading AS SELECT TRIM(LEADING 'x' FROM a) FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_leading
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_leading
: bytes
| CREATE VIEW `v_trim_leading` AS select trim(leading 'x' from `test`.`tv`.`a`) AS `TRIM(LEADING 'x' FROM a)` from `test`.`tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_trim_trailing
< ok 0 0
> exec
| CREATE VIEW v_trim_trailing AS SELECT TRIM(TRAILING 'x' FROM a) FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_trailing
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_trailing
: bytes
| CREATE VIEW `v_trim_trailing` AS select trim(trailing 'x' from `test`.`tv`.`a`) AS `TRIM(TRAILING 'x' FROM a)` from `test`.`tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_trim_both
< ok 0 0
> exec
| CREATE VIEW v_trim_both AS SELECT TRIM(BOTH 'x' FROM a) FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_both
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_both
: bytes
| CREATE VIEW `v_trim_both` AS select trim(both 'x' from `test`.`tv`.`a`) AS `TRIM(BOTH 'x' FROM a)` from `test`.`tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_simple_case
< ok 0 0
> exec
| CREATE VIEW v_simple_case AS SELECT CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_simple_case
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_simple_case
: bytes
| CREATE VIEW `v_simple_case` AS select (case `test`.`t`.`a` when 1 then 'one' when 2 then 'two' else 'other' end) AS `CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
