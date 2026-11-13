-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSharedContainerResetsStateBetweenUses default#1
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
| SET SESSION foreign_key_checks = 0
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0

== TestSharedContainerResetsStateBetweenUses default#2
> query
| SELECT DATABASE(), @@foreign_key_checks, @@session.sql_mode
< result
< column "DATABASE()" "VARCHAR"
< column "@@foreign_key_checks" "BIGINT"
< column "@@session.sql_mode" "VARCHAR"
< row
: bytes
| test
: int64 1
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
> query
| SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = 'dirty_shared_state'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
