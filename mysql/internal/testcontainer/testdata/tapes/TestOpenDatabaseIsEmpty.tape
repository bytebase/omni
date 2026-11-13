-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOpenDatabaseIsEmpty default#1
> query
| SELECT DATABASE(), (SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE())
< result
< column "DATABASE()" "VARCHAR"
< column "(SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE())" "BIGINT"
< row
: bytes
| omni_t_empty
: int64 0
