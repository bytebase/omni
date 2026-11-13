-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDeparseContainer_1_3_LiteralsSpacing default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_esc_quotes
: bytes
| CREATE VIEW `v_esc_quotes` AS select 'it\'s' AS `it's` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_esc_bslash
: bytes
| CREATE VIEW `v_esc_bslash` AS select 'back\\slash' AS `back\slash` from `test`.`t`
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
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_in_list
: bytes
| CREATE VIEW `v_in_list` AS select (`test`.`t`.`a` in (1,2,3)) AS `a IN (1, 2, 3)` from `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
