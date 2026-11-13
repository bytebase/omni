-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_DropCompositePKColumn default#1
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
| CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b))
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN a
< error
| 1072 42000 Key column 'a' doesn't exist in table
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
| CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b))
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN a, DROP PRIMARY KEY
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
| CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b))
< ok 0 0
> exec
| ALTER TABLE t DROP PRIMARY KEY, DROP COLUMN a
< ok 0 0
