-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
