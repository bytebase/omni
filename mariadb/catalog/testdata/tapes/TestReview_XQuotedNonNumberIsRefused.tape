-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_XQuotedNonNumberIsRefused default#1
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
| CREATE TABLE t (a INT DEFAULT (X'41'))
< error
| 1067 42000 Invalid default value for 'a'
> exec
| CREATE TABLE t (a INT DEFAULT (X''))
< error
| 1067 42000 Invalid default value for 'a'
> exec
| CREATE TABLE t (a YEAR DEFAULT (X'41'))
< error
| 1067 42000 Invalid default value for 'a'
