-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTableSampleViewTiDBDiffers default#1
> exec
| DROP DATABASE IF EXISTS omni_tablesample_view
< ok 0 0
> exec
| CREATE DATABASE omni_tablesample_view
< ok 0 0
> exec
| USE omni_tablesample_view
< ok 0 0
> exec
| CREATE TABLE tbl (a INT, b INT, PRIMARY KEY (a))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM tbl TABLESAMPLE REGIONS() REPEATABLE (1) WHERE a > 1
< ok 0 0
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`a`) AS SELECT `a` AS `a` FROM `omni_tablesample_view`.`tbl` TABLESAMPLE REGION () REPEATABLE(1) WHERE `a`>1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| SELECT * FROM v
< error
| 1105 HY000 line 1 column 71 near "REGION () REPEATABLE(1) WHERE `a`>1" 
> exec
| SELECT `a` AS `a` FROM `omni_tablesample_view`.`tbl` TABLESAMPLE REGIONS () REPEATABLE(1) WHERE `a`>1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_tablesample_view
< ok 0 0
