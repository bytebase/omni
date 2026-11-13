-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_CheckNameScope default#1
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
| CREATE TABLE t1 (a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, CONSTRAINT ck CHECK (a < 9))
< ok 0 0
> exec
| ALTER TABLE t1 ADD CONSTRAINT ck CHECK (a < 10)
< error
| 1826 HY000 Duplicate CHECK constraint name 'ck'
