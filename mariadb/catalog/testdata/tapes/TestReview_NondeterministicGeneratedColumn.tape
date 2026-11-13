-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_NondeterministicGeneratedColumn default#1
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
| CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL)
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
| CREATE TABLE t (a DOUBLE AS (RAND()) STORED)
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the GENERATED ALWAYS AS clause of `a`
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
| CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, KEY(a))
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the GENERATED ALWAYS AS clause of `a`
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
| CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL UNIQUE)
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the GENERATED ALWAYS AS clause of `a`
