-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ChangingAGeneratedDependencyName default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) STORED)
< ok 0 0
> exec
| ALTER TABLE t CHANGE a z DOUBLE AS (RAND()) VIRTUAL
< error
| 1901 HY000 Function or expression 'z' cannot be used in the GENERATED ALWAYS AS clause of `b`
