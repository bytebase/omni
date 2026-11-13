-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestStartContainerStartsEmpty default#1
> exec
| CREATE DATABASE dirty_shared_state
< ok 1 0
> exec
| USE dirty_shared_state
< ok 0 0
> exec
| CREATE TABLE dirty_table (id INT)
< ok 0 0
> exec
| CREATE USER dirty_user
< ok 0 0
> exec
| SET SESSION foreign_key_checks = 0
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0

== TestStartContainerStartsEmpty default#2
> query
| SELECT DATABASE(), @@foreign_key_checks, @@session.sql_mode
< result
< column "DATABASE()" "VARCHAR"
< column "@@foreign_key_checks" "BIGINT"
< column "@@session.sql_mode" "VARCHAR"
< row
: bytes
| catalog_shared
: int64 1
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
> query
| SELECT (SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = 'dirty_shared_state') + (SELECT COUNT(*) FROM mysql.global_priv WHERE User = 'dirty_user')
< result
< column "(SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = 'dirty_shared_state') + (SELECT COUNT(*) FROM mysql.global_priv WHERE User = 'dirty_user')" "BIGINT"
< row
: int64 0
