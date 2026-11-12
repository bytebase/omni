-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAsOfViewTiDBDiffers default#1
> exec
| DROP DATABASE IF EXISTS omni_as_of_view
< ok 0 0
> exec
| CREATE DATABASE omni_as_of_view
< ok 0 0
> exec
| USE omni_as_of_view
< ok 0 0
> exec
| CREATE TABLE tbl (a INT, b INT, PRIMARY KEY (a))
< ok 0 0
> exec
| DO SLEEP(3)
< ok 0 0
> exec
| SELECT a FROM tbl AS OF TIMESTAMP NOW() - INTERVAL 1 SECOND WHERE a > 1
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM tbl AS OF TIMESTAMP NOW() - INTERVAL 1 SECOND WHERE a > 1
< error
| 1356 HY000 View 'omni_as_of_view.v' references invalid table(s) or column(s) or function(s) or definer/invoker of view lack rights to use them
> exec
| SELECT `a` AS `a` FROM `omni_as_of_view`.`tbl` AS OF TIMESTAMP DATE_SUB(NOW(), INTERVAL 1 SECOND) WHERE `a`>1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_as_of_view
< ok 0 0
