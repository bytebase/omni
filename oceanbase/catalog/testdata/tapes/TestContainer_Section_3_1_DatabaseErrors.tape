-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
