-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestScenario_C10 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
| NONE
: bytes
| YES
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
| NONE
: bytes
| YES
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH CHECK OPTION
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
| YES
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
| YES
> exec
| DROP DATABASE IF EXISTS testdb
< ok 3 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_agg` AS select count(0) AS `COUNT(*)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_merge_bad` AS select distinct `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 3 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
< ok 3 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
| YES
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
< ok 4 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
| NO
< row
: bytes
| b
: bytes
| YES
