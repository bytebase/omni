-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
