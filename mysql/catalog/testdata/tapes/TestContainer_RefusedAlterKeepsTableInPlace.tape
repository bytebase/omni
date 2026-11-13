-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_RefusedAlterKeepsTableInPlace/ALTER_TABLE_d1.t_ADD_COLUMN_a_INT,_RENAME_TO_d2.u default#1
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
| SELECT CONCAT(table_schema, '.', table_name) FROM information_schema.tables WHERE table_schema IN ('d1', 'd2')
< result
< column "CONCAT(table_schema, '.', table_name)" "VARCHAR"
< row
: bytes
| d1.t

== TestContainer_RefusedAlterKeepsTableInPlace/ALTER_TABLE_d1.t_ADD_UNIQUE_(b),_RENAME_TO_d1.v default#1
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
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).
> query
| SELECT CONCAT(table_schema, '.', table_name) FROM information_schema.tables WHERE table_schema IN ('d1', 'd2')
< result
< column "CONCAT(table_schema, '.', table_name)" "VARCHAR"
< row
: bytes
| d1.t

== TestContainer_RefusedAlterKeepsTableInPlace/ALTER_TABLE_d1.t_ADD_UNIQUE_(b),_RENAME_TO_d2.t default#1
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
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).
> query
| SELECT CONCAT(table_schema, '.', table_name) FROM information_schema.tables WHERE table_schema IN ('d1', 'd2')
< result
< column "CONCAT(table_schema, '.', table_name)" "VARCHAR"
< row
: bytes
| d1.t

== TestContainer_RefusedAlterKeepsTableInPlace/ALTER_TABLE_d1.t_RENAME_TO_d2.u,_ADD_COLUMN_a_INT default#1
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
| SELECT CONCAT(table_schema, '.', table_name) FROM information_schema.tables WHERE table_schema IN ('d1', 'd2')
< result
< column "CONCAT(table_schema, '.', table_name)" "VARCHAR"
< row
: bytes
| d1.t

== TestContainer_RefusedAlterKeepsTableInPlace/ALTER_TABLE_d1.t_RENAME_TO_d2.u,_ADD_UNIQUE_(b) default#1
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
| 1503 HY000 A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered).
> query
| SELECT CONCAT(table_schema, '.', table_name) FROM information_schema.tables WHERE table_schema IN ('d1', 'd2')
< result
< column "CONCAT(table_schema, '.', table_name)" "VARCHAR"
< row
: bytes
| d1.t
