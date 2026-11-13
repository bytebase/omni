-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_5_1_UseStatement default#1
> exec
| DROP DATABASE IF EXISTS use_test_db
< ok 0 0
> exec
| CREATE DATABASE use_test_db
< ok 1 0
> exec
| USE use_test_db
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_use_check
< ok 0 0
> exec
| CREATE TABLE t_use_check (id INT PRIMARY KEY)
< ok 0 0
> query
| SELECT DATABASE()
< result
< column "DATABASE()" "VARCHAR"
< row
: bytes
| use_test_db
> query
| SELECT TABLE_SCHEMA FROM information_schema.TABLES WHERE TABLE_NAME = 't_use_check'
< result
< column "TABLE_SCHEMA" "CHAR"
< row
: bytes
| use_test_db
> exec
| DROP DATABASE IF EXISTS use_test_db
< ok 1 0
> exec
| DROP DATABASE IF EXISTS use_nonexistent_db
< ok 0 0
> exec
| USE use_nonexistent_db
< error
| 1049 42000 Unknown database 'use_nonexistent_db'
