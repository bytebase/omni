-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_3_3_ColumnErrors default#1
> exec
| DROP TABLE IF EXISTS t_err_nocol
< ok 0 0
> exec
| CREATE TABLE t_err_nocol (id INT, name VARCHAR(50))
< ok 0 0
> exec
| ALTER TABLE t_err_nocol MODIFY COLUMN name VARCHAR(50) AFTER nonexistent
< error
| 1054 42S22 Unknown column 'nonexistent' in 'generated column function'
> exec
| DROP TABLE IF EXISTS t_err_nocol
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_err_dupcol
< ok 0 0
> exec
| CREATE TABLE t_err_dupcol (a INT, a VARCHAR(10))
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP TABLE IF EXISTS t_err_multipk
< ok 0 0
> exec
| CREATE TABLE t_err_multipk (a INT, b INT, PRIMARY KEY (a), PRIMARY KEY (b))
< error
| 1068 42000 Multiple primary key defined
