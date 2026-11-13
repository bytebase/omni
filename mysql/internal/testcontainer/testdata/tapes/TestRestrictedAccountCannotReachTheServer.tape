-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestRestrictedAccountCannotReachTheServer default#1
> exec
| CREATE TABLE kept (a INT)
< ok 0 0
> query
| SELECT COUNT(*) FROM kept
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0

== TestRestrictedAccountCannotReachTheServer default#2
> exec
| SHUTDOWN
< error
| 1227 42000 Access denied; you need (at least one of) the SHUTDOWN privilege(s) for this operation
> exec
| RESTART
< error
| 1227 42000 Access denied; you need (at least one of) the SHUTDOWN privilege(s) for this operation
> exec
| SET GLOBAL sql_mode = 'ANSI_QUOTES'
< error
| 1227 42000 Access denied; you need (at least one of) the SUPER or SYSTEM_VARIABLES_ADMIN privilege(s) for this operation
> exec
| SET PERSIST max_connections = 1
< error
| 1227 42000 Access denied; you need (at least one of) the SUPER or SYSTEM_VARIABLES_ADMIN privilege(s) for this operation
> exec
| ALTER USER root IDENTIFIED BY 'gone'
< error
| 1227 42000 Access denied; you need (at least one of) the CREATE USER privilege(s) for this operation
> exec
| SET PASSWORD FOR root = 'gone'
< error
| 1044 42000 Access denied for user 'omni_test'@'%' to database 'mysql'
> exec
| CREATE USER intruder
< error
| 1227 42000 Access denied; you need (at least one of) the CREATE USER privilege(s) for this operation
> exec
| DROP USER root
< error
| 1227 42000 Access denied; you need (at least one of) the CREATE USER privilege(s) for this operation
> exec
| GRANT ALL ON *.* TO omni_test
< error
| 1045 28000 Access denied for user 'omni_test'@'%' (using password: YES)
> exec
| DROP DATABASE mysql
< error
| 1044 42000 Access denied for user 'omni_test'@'%' to database 'mysql'
> exec
| CREATE DATABASE outside_the_prefix
< error
| 1044 42000 Access denied for user 'omni_test'@'%' to database 'outside_the_prefix'
> exec
| FLUSH TABLES WITH READ LOCK
< error
| 1227 42000 Access denied; you need (at least one of) the RELOAD or FLUSH_TABLES privilege(s) for this operation
> exec
| FLUSH PRIVILEGES
< error
| 1227 42000 Access denied; you need (at least one of) the RELOAD privilege(s) for this operation
> exec
| LOCK INSTANCE FOR BACKUP
< error
| 1227 42000 Access denied; you need (at least one of) the BACKUP_ADMIN privilege(s) for this operation
> exec
| CREATE RESOURCE GROUP rg TYPE = USER
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN privilege(s) for this operation
> exec
| RESET MASTER
< error
| 1227 42000 Access denied; you need (at least one of) the RELOAD or FLUSH_TABLES privilege(s) for this operation
> exec
| PURGE BINARY LOGS BEFORE NOW()
< error
| 1227 42000 Access denied; you need (at least one of) the SUPER or BINLOG_ADMIN privilege(s) for this operation
> exec
| INSTALL PLUGIN p SONAME 'p.so'
< error
| 1142 42000 INSERT command denied to user 'omni_test'@'172.17.0.1' for table 'plugin'
> exec
| KILL 1
< error
| 1094 HY000 Unknown thread id: 1
> query
| SELECT CURRENT_USER()
< result
< column "CURRENT_USER()" "VARCHAR"
< row
: bytes
| omni_test@%
