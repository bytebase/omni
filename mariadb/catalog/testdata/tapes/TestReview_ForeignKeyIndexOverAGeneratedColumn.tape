-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ForeignKeyIndexOverAGeneratedColumn default#1
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
| CREATE TABLE p (a DOUBLE PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES p(a)
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (a DOUBLE PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, CONSTRAINT fk FOREIGN KEY (a) REFERENCES p(a))
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the GENERATED ALWAYS AS clause of `a`
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
| CREATE TABLE p (a DOUBLE PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (r DOUBLE AS (RAND()) VIRTUAL, a DOUBLE AS (r) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES p(a)
< error
| 1901 HY000 Function or expression 'r' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (a DOUBLE PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (x DOUBLE, a DOUBLE AS (x) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES p(a)
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `x` double DEFAULT NULL,
|   `a` double GENERATED ALWAYS AS (`x`) VIRTUAL,
|   KEY `fk` (`a`),
|   CONSTRAINT `fk` FOREIGN KEY (`a`) REFERENCES `p` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
