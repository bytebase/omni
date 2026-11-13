-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseViewOptions default#1
> exec
| DROP DATABASE IF EXISTS obview_options
< ok 0 0
> exec
| CREATE DATABASE obview_options
< ok 1 0
> exec
| USE obview_options
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v_plain AS SELECT a FROM t
< ok 0 0
> exec
| CREATE ALGORITHM=MERGE VIEW v_merge AS SELECT a FROM t
< ok 0 0
> exec
| CREATE ALGORITHM=TEMPTABLE VIEW v_temp AS SELECT a FROM t
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW v_inv AS SELECT a FROM t
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW v_def AS SELECT a FROM t
< ok 0 0
> exec
| CREATE VIEW v_check AS SELECT a FROM t WHERE a > 0 WITH CHECK OPTION
< ok 0 0
> query
| SHOW CREATE VIEW v_plain
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_plain
: bytes
| CREATE VIEW `v_plain` AS select `obview_options`.`t`.`a` AS `a` from `obview_options`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT SECURITY_TYPE, CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'obview_options' AND TABLE_NAME = ?
: string
| v_plain
< result
< column "SECURITY_TYPE" "VARCHAR"
< column "CHECK_OPTION" "VARCHAR"
< row
: bytes
| NONE
: bytes
| NONE
> query
| SHOW CREATE VIEW v_merge
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_merge
: bytes
| CREATE VIEW `v_merge` AS select `obview_options`.`t`.`a` AS `a` from `obview_options`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT SECURITY_TYPE, CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'obview_options' AND TABLE_NAME = ?
: string
| v_merge
< result
< column "SECURITY_TYPE" "VARCHAR"
< column "CHECK_OPTION" "VARCHAR"
< row
: bytes
| NONE
: bytes
| NONE
> query
| SHOW CREATE VIEW v_temp
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_temp
: bytes
| CREATE VIEW `v_temp` AS select `obview_options`.`t`.`a` AS `a` from `obview_options`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT SECURITY_TYPE, CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'obview_options' AND TABLE_NAME = ?
: string
| v_temp
< result
< column "SECURITY_TYPE" "VARCHAR"
< column "CHECK_OPTION" "VARCHAR"
< row
: bytes
| NONE
: bytes
| NONE
> query
| SHOW CREATE VIEW v_inv
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_inv
: bytes
| CREATE VIEW `v_inv` AS select `obview_options`.`t`.`a` AS `a` from `obview_options`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT SECURITY_TYPE, CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'obview_options' AND TABLE_NAME = ?
: string
| v_inv
< result
< column "SECURITY_TYPE" "VARCHAR"
< column "CHECK_OPTION" "VARCHAR"
< row
: bytes
| NONE
: bytes
| NONE
> query
| SHOW CREATE VIEW v_def
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_def
: bytes
| CREATE VIEW `v_def` AS select `obview_options`.`t`.`a` AS `a` from `obview_options`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT SECURITY_TYPE, CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'obview_options' AND TABLE_NAME = ?
: string
| v_def
< result
< column "SECURITY_TYPE" "VARCHAR"
< column "CHECK_OPTION" "VARCHAR"
< row
: bytes
| NONE
: bytes
| NONE
> query
| SHOW CREATE VIEW v_check
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_check
: bytes
| CREATE VIEW `v_check` AS select `obview_options`.`t`.`a` AS `a` from `obview_options`.`t` where (`obview_options`.`t`.`a` > 0) WITH CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT SECURITY_TYPE, CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'obview_options' AND TABLE_NAME = ?
: string
| v_check
< result
< column "SECURITY_TYPE" "VARCHAR"
< column "CHECK_OPTION" "VARCHAR"
< row
: bytes
| NONE
: bytes
| CASCADED
