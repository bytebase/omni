-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ReplacedGeneratedColumnKeepsItsIndex default#1
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
| CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, KEY (a))
< ok 0 0
> exec
| ALTER TABLE t CHANGE a b DOUBLE AS (RAND()) VIRTUAL
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the GENERATED ALWAYS AS clause of `b`
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
| CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, KEY (a))
< ok 0 0
> exec
| ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL
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
| CREATE TABLE t (x INT, a DOUBLE AS (1) VIRTUAL, KEY (x, a))
< ok 0 0
> exec
| ALTER TABLE t CHANGE a b DOUBLE AS (RAND()) VIRTUAL
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the GENERATED ALWAYS AS clause of `b`
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
| CREATE TABLE t (a DOUBLE AS (1) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE t CHANGE a b DOUBLE AS (RAND()) VIRTUAL
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
|   `b` double GENERATED ALWAYS AS (rand()) VIRTUAL
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
| CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, KEY (a))
< ok 0 0
> exec
| ALTER TABLE t CHANGE a b DOUBLE AS (2) VIRTUAL
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
|   `b` double GENERATED ALWAYS AS (2) VIRTUAL,
|   KEY `a` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
