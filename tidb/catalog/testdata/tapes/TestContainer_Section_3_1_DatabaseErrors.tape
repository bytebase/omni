-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_3_1_DatabaseErrors default#1
> exec
| DROP DATABASE IF EXISTS db_err_dup
< ok 0 0
> exec
| CREATE DATABASE db_err_dup
< ok 0 0
> exec
| CREATE DATABASE db_err_dup
< error
| 1007 HY000 Can't create database 'db_err_dup'; database exists
> exec
| DROP DATABASE IF EXISTS db_err_dup
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_err_unknown_xyz
< ok 0 0
> exec
| USE db_err_unknown_xyz
< error
| 1049 42000 Unknown database 'db_err_unknown_xyz'
> exec
| DROP DATABASE db_err_unknown_xyz
< error
| 1008 HY000 Can't drop database 'db_err_unknown_xyz'; database doesn't exist
