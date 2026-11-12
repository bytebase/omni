-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScenario_C10 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SELECT DEFINER, SECURITY_TYPE, CHECK_OPTION, IS_UPDATABLE
|              FROM information_schema.VIEWS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v'
< result
< column "DEFINER" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< column "CHECK_OPTION" "VARCHAR"
< column "IS_UPDATABLE" "VARCHAR"
< row
: bytes
| root@%
: bytes
| DEFINER
: bytes
| CASCADED
: bytes
| NO
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`a`) AS SELECT `a` AS `a` FROM `testdb`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SELECT DEFINER, SECURITY_TYPE, CHECK_OPTION, IS_UPDATABLE
|              FROM information_schema.VIEWS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v'
< result
< column "DEFINER" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< column "CHECK_OPTION" "VARCHAR"
< column "IS_UPDATABLE" "VARCHAR"
< row
: bytes
| root@%
: bytes
| DEFINER
: bytes
| CASCADED
: bytes
| NO
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION
< ok 0 0
> exec
| CREATE VIEW v2 AS SELECT a FROM t WHERE a > 0 WITH LOCAL CHECK OPTION
< ok 0 0
> query
| SELECT DEFINER, SECURITY_TYPE, CHECK_OPTION, IS_UPDATABLE
|              FROM information_schema.VIEWS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v1'
< result
< column "DEFINER" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< column "CHECK_OPTION" "VARCHAR"
< column "IS_UPDATABLE" "VARCHAR"
< row
: bytes
| root@%
: bytes
| DEFINER
: bytes
| CASCADED
: bytes
| NO
> query
| SELECT DEFINER, SECURITY_TYPE, CHECK_OPTION, IS_UPDATABLE
|              FROM information_schema.VIEWS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v2'
< result
< column "DEFINER" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< column "CHECK_OPTION" "VARCHAR"
< column "IS_UPDATABLE" "VARCHAR"
< row
: bytes
| root@%
: bytes
| DEFINER
: bytes
| LOCAL
: bytes
| NO
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED VIEW v_agg AS SELECT COUNT(*) FROM t
< ok 0 0
> exec
| CREATE ALGORITHM=MERGE VIEW v_merge_bad AS SELECT DISTINCT a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_agg
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_agg
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_agg` (`COUNT(*)`) AS SELECT COUNT(1) AS `COUNT(*)` FROM `testdb`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW v_merge_bad
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_merge_bad
: bytes
| CREATE ALGORITHM=MERGE DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_merge_bad` (`a`) AS SELECT DISTINCT `a` AS `a` FROM `testdb`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SELECT SECURITY_TYPE FROM information_schema.VIEWS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v'
< result
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| DEFINER
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v_auto AS SELECT a, a+1, COUNT(*) FROM t GROUP BY a
< ok 0 0
> exec
| CREATE VIEW v_list (x,y,z) AS SELECT a, a+1, COUNT(*) FROM t GROUP BY a
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_auto'
|              ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
< row
: bytes
| a+1
< row
: bytes
| COUNT(*)
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_list'
|              ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| x
< row
: bytes
| y
< row
: bytes
| z
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v_ok AS SELECT a FROM t
< ok 0 0
> exec
| CREATE VIEW v_distinct AS SELECT DISTINCT a FROM t
< ok 0 0
> exec
| CREATE ALGORITHM=TEMPTABLE VIEW v_temp AS SELECT a FROM t
< ok 0 0
> query
| SELECT IS_UPDATABLE FROM information_schema.VIEWS
|                  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_ok'
< result
< column "IS_UPDATABLE" "VARCHAR"
< row
: bytes
| NO
> query
| SELECT IS_UPDATABLE FROM information_schema.VIEWS
|                  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_distinct'
< result
< column "IS_UPDATABLE" "VARCHAR"
< row
: bytes
| NO
> query
| SELECT IS_UPDATABLE FROM information_schema.VIEWS
|                  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_temp'
< result
< column "IS_UPDATABLE" "VARCHAR"
< row
: bytes
| NO
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t1 (id INT NOT NULL, a INT NOT NULL)
< ok 0 0
> exec
| CREATE TABLE t2 (id INT NOT NULL, b INT NOT NULL)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t1.a, t2.b FROM t1 LEFT JOIN t2 ON t1.id = t2.id
< ok 0 0
> query
| SELECT COLUMN_NAME, IS_NULLABLE FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v'
|              ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| a
: bytes
| YES
< row
: bytes
| b
: bytes
| YES
