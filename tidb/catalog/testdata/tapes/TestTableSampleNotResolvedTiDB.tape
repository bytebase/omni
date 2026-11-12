-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTableSampleNotResolvedTiDB default#1
> exec
| DROP DATABASE IF EXISTS omni_tablesample_resolve
< ok 0 0
> exec
| CREATE DATABASE omni_tablesample_resolve
< ok 0 0
> exec
| USE omni_tablesample_resolve
< ok 0 0
> exec
| CREATE TABLE tbl (a INT, b INT, PRIMARY KEY (a))
< ok 0 0
> exec
| SELECT a FROM tbl TABLESAMPLE REGIONS(nosuch)
< ok 0 0
> exec
| select `tbl`.`a` from `omni_tablesample_resolve`.`tbl` tablesample regions (`nosuch`)
< ok 0 0
> exec
| SELECT a FROM tbl TABLESAMPLE REGIONS(nosuch.col)
< ok 0 0
> exec
| select `tbl`.`a` from `omni_tablesample_resolve`.`tbl` tablesample regions (`nosuch`.`col`)
< ok 0 0
> exec
| SELECT a FROM tbl TABLESAMPLE REGIONS(nosuchfunc(1))
< ok 0 0
> exec
| select `tbl`.`a` from `omni_tablesample_resolve`.`tbl` tablesample regions (nosuchfunc(1))
< ok 0 0
> exec
| SELECT a FROM tbl TABLESAMPLE REGIONS((SELECT nosuch FROM nosuch))
< ok 0 0
> exec
| select `tbl`.`a` from `omni_tablesample_resolve`.`tbl` tablesample regions ((select `nosuch` from `nosuch`))
< ok 0 0
> exec
| SELECT a FROM tbl TABLESAMPLE REGIONS() REPEATABLE (nosuch)
< ok 0 0
> exec
| select `tbl`.`a` from `omni_tablesample_resolve`.`tbl` tablesample regions () repeatable (`nosuch`)
< ok 0 0
> exec
| SELECT a FROM tbl TABLESAMPLE REGIONS() REPEATABLE (1) WHERE a > 1
< ok 0 0
> exec
| select `tbl`.`a` from `omni_tablesample_resolve`.`tbl` tablesample regions () repeatable (1) where (`tbl`.`a` > 1)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_tablesample_resolve
< ok 0 0
