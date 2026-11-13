-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| 1051 42S02 Unknown table 'catalog_shared.t_err_noexist'
> exec
| DROP TABLE IF EXISTS t_err_noexist2
< ok 0 0
> exec
| ALTER TABLE t_err_noexist2 ADD COLUMN x INT
< error
| 1146 42S02 Table 'catalog_shared.t_err_noexist2' doesn't exist
