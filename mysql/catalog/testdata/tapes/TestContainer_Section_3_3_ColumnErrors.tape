-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
| 1054 42S22 Unknown column 'nonexistent' in 't_err_nocol'
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
