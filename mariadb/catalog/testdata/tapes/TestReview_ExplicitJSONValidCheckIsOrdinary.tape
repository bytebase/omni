-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ExplicitJSONValidCheckIsOrdinary default#1
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
| CREATE TABLE t (j LONGTEXT, CONSTRAINT j CHECK (json_valid(j)))
< ok 0 0
> exec
| ALTER TABLE t DROP CONSTRAINT j
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (j LONGTEXT, CONSTRAINT ck CHECK (json_valid(j)))
< ok 0 0
> exec
| ALTER TABLE t DROP CONSTRAINT ck
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (j JSON, CONSTRAINT ck CHECK (json_valid(j)))
< ok 0 0
> exec
| ALTER TABLE t DROP CONSTRAINT ck
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (j JSON)
< ok 0 0
> exec
| ALTER TABLE t DROP CONSTRAINT j
< error
| 1091 42000 Can't DROP CONSTRAINT `j`; check that it exists
