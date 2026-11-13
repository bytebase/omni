-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOceanBaseRefusedAlterKeepsTableInPlace default#1
> exec
| DROP DATABASE IF EXISTS d1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS d2
< ok 0 0
> exec
| CREATE DATABASE d1
< ok 1 0
> exec
| CREATE DATABASE d2
< ok 1 0
> exec
| CREATE TABLE d1.t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE d1.t ADD UNIQUE (b), RENAME TO d2.t
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> query
| SELECT GROUP_CONCAT(CONCAT(table_schema, '.', table_name)) FROM information_schema.tables WHERE table_schema IN ('d1', 'd2')
< result
< column "GROUP_CONCAT(CONCAT(table_schema, '.', table_name))" "VARCHAR"
< row
: bytes
| d1.t
> exec
| DROP DATABASE IF EXISTS d1
< ok 1 0
> exec
| DROP DATABASE IF EXISTS d2
< ok 0 0
> exec
| DROP DATABASE IF EXISTS d1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS d2
< ok 0 0
> exec
| CREATE DATABASE d1
< ok 1 0
> exec
| CREATE DATABASE d2
< ok 1 0
> exec
| CREATE TABLE d1.t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE d1.t RENAME TO d2.u, ADD UNIQUE (b)
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> query
| SELECT GROUP_CONCAT(CONCAT(table_schema, '.', table_name)) FROM information_schema.tables WHERE table_schema IN ('d1', 'd2')
< result
< column "GROUP_CONCAT(CONCAT(table_schema, '.', table_name))" "VARCHAR"
< row
: bytes
| d1.t
> exec
| DROP DATABASE IF EXISTS d1
< ok 1 0
> exec
| DROP DATABASE IF EXISTS d2
< ok 0 0
> exec
| DROP DATABASE IF EXISTS d1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS d2
< ok 0 0
> exec
| CREATE DATABASE d1
< ok 1 0
> exec
| CREATE DATABASE d2
< ok 1 0
> exec
| CREATE TABLE d1.t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE d1.t ADD UNIQUE (b), RENAME TO d1.v
< error
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function
> query
| SELECT GROUP_CONCAT(CONCAT(table_schema, '.', table_name)) FROM information_schema.tables WHERE table_schema IN ('d1', 'd2')
< result
< column "GROUP_CONCAT(CONCAT(table_schema, '.', table_name))" "VARCHAR"
< row
: bytes
| d1.t
> exec
| DROP DATABASE IF EXISTS d1
< ok 1 0
> exec
| DROP DATABASE IF EXISTS d2
< ok 0 0
> exec
| DROP DATABASE IF EXISTS d1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS d2
< ok 0 0
> exec
| CREATE DATABASE d1
< ok 1 0
> exec
| CREATE DATABASE d2
< ok 1 0
> exec
| CREATE TABLE d1.t (a INT, b INT)
< ok 0 0
> exec
| ALTER TABLE d1.t RENAME TO d2.u, ADD COLUMN a INT
< error
| 1060 42S21 Duplicate column name 'a'
> query
| SELECT GROUP_CONCAT(CONCAT(table_schema, '.', table_name)) FROM information_schema.tables WHERE table_schema IN ('d1', 'd2')
< result
< column "GROUP_CONCAT(CONCAT(table_schema, '.', table_name))" "VARCHAR"
< row
: bytes
| d1.t
> exec
| DROP DATABASE IF EXISTS d1
< ok 1 0
> exec
| DROP DATABASE IF EXISTS d2
< ok 0 0
> exec
| DROP DATABASE IF EXISTS d1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS d2
< ok 0 0
> exec
| CREATE DATABASE d1
< ok 1 0
> exec
| CREATE DATABASE d2
< ok 1 0
> exec
| CREATE TABLE d1.t (a INT, b INT)
< ok 0 0
> exec
| ALTER TABLE d1.t ADD COLUMN a INT, RENAME TO d2.u
< error
| 1060 42S21 Duplicate column name 'a'
> query
| SELECT GROUP_CONCAT(CONCAT(table_schema, '.', table_name)) FROM information_schema.tables WHERE table_schema IN ('d1', 'd2')
< result
< column "GROUP_CONCAT(CONCAT(table_schema, '.', table_name))" "VARCHAR"
< row
: bytes
| d1.t
> exec
| DROP DATABASE IF EXISTS d1
< ok 1 0
> exec
| DROP DATABASE IF EXISTS d2
< ok 0 0
