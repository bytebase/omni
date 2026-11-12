-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAsOfNotInQueryScopeTiDB default#1
> exec
| DROP DATABASE IF EXISTS omni_as_of_scope
< ok 0 0
> exec
| CREATE DATABASE omni_as_of_scope
< ok 0 0
> exec
| USE omni_as_of_scope
< ok 0 0
> exec
| CREATE TABLE tbl (a INT, b INT, PRIMARY KEY (a))
< ok 0 0
> exec
| SELECT a FROM tbl AS OF TIMESTAMP a
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| SELECT a FROM tbl AS OF TIMESTAMP tbl.a
< error
| 1054 42S22 Unknown column 'tbl.a' in 'expression'
> exec
| SELECT a FROM tbl AS OF TIMESTAMP nosuch
< error
| 1054 42S22 Unknown column 'nosuch' in 'expression'
> exec
| DROP DATABASE IF EXISTS omni_as_of_scope
< ok 0 0
