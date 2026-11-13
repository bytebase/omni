-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestCurrentDatabaseDeleteFormsOracle default#1
> exec
| USE test
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t (a INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t1 (a INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t2 (a INT)
< ok 0 0
> exec
| USE test
< ok 0 0
> exec
| PREPARE omni_delete FROM 'DELETE t FROM test.t'
< ok 0 0
> exec
| USE test
< ok 0 0
> exec
| PREPARE omni_delete FROM 'DELETE test.t FROM t'
< ok 0 0
> exec
| USE test
< ok 0 0
> exec
| PREPARE omni_delete FROM 'DELETE test.t1 FROM t1'
< ok 0 0
> exec
| USE test
< ok 0 0
> exec
| PREPARE omni_delete FROM 'DELETE t1 FROM test.t1'
< ok 0 0
> exec
| USE test
< ok 0 0
> exec
| PREPARE omni_delete FROM 'DELETE test.t.* FROM t'
< ok 0 0
> exec
| USE test
< ok 0 0
> exec
| PREPARE omni_delete FROM 'DELETE test.t1.*, t2 FROM t1, t2'
< ok 0 0
> exec
| USE test
< ok 0 0
> exec
| PREPARE omni_delete FROM 'DELETE FROM test.t2 USING t2'
< ok 0 0
