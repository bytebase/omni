-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_3_1_DatabaseErrors default#1
> exec
| DROP DATABASE IF EXISTS db_err_dup
< ok 0 0
> exec
| CREATE DATABASE db_err_dup
< ok 1 0
> exec ^
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
