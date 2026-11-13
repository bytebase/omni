-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_IndexOnNondeterministicGeneratedColumn default#1
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
| CREATE INDEX idx ON t(a)
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
| CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX idx(a)
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
| CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE idx(a)
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
| CREATE TABLE t (b INT, a DOUBLE AS (RAND()) VIRTUAL)
< ok 0 0
> exec
| CREATE INDEX idx ON t(b,a)
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
| CREATE TABLE t (a DATETIME AS (NOW()) VIRTUAL)
< ok 0 0
> exec
| CREATE INDEX idx ON t(a)
< error
| 1901 HY000 Function or expression 'current_timestamp()' cannot be used in the GENERATED ALWAYS AS clause of `a`
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
| CREATE TABLE t (a CHAR(36) AS (UUID()) VIRTUAL)
< ok 0 0
> exec
| CREATE INDEX idx ON t(a)
< error
| 1901 HY000 Function or expression 'uuid()' cannot be used in the GENERATED ALWAYS AS clause of `a`
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
| CREATE TABLE t (a TEXT AS (UUID()) VIRTUAL)
< ok 0 0
> exec
| CREATE INDEX idx ON t(a)
< error
| 1901 HY000 Function or expression 'uuid()' cannot be used in the GENERATED ALWAYS AS clause of `a`
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
| CREATE TABLE t (b INT, a INT AS (b + 1) VIRTUAL)
< ok 0 0
> exec
| CREATE INDEX idx ON t(a)
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
|   `b` int(11) DEFAULT NULL,
|   `a` int(11) GENERATED ALWAYS AS (`b` + 1) VIRTUAL,
|   KEY `idx` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE INDEX idx ON t(a)
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
|   `a` int(11) DEFAULT NULL,
|   KEY `idx` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
