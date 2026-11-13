-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_3_2_TableErrors default#1
> exec
| DROP TABLE IF EXISTS t_err_dup
< ok 0 0
> exec
| CREATE TABLE t_err_dup (id INT)
< ok 0 0
> exec ^
< error
| 1050 42S01 Table 't_err_dup' already exists
> exec
| DROP TABLE IF EXISTS t_err_dup
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_err_noexist
< ok 0 0
> exec
| DROP TABLE t_err_noexist
< error
| 1051 42S02 Unknown table 'test.t_err_noexist'
> exec
| DROP TABLE IF EXISTS t_err_noexist2
< ok 0 0
> exec
| ALTER TABLE t_err_noexist2 ADD COLUMN x INT
< error
| 1146 42S02 Table 'test.t_err_noexist2' doesn't exist
