-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
| 1062 23000 Duplicated primary key
