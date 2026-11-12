-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_3_2_TableErrors default#1
> exec
| DROP TABLE IF EXISTS t_err_dup
< ok 0 0
> exec
| CREATE TABLE t_err_dup (id INT)
< ok 0 0
> exec
| CREATE TABLE t_err_dup (id INT)
< error
| 1050 42S01 Table 'test.t_err_dup' already exists
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
