-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ReplacingAGeneratedDependency default#1
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
| ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL
< error
| 1901 HY000 Function or expression 'a' cannot be used in the GENERATED ALWAYS AS clause of `b`
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
|   `a` double GENERATED ALWAYS AS (1) VIRTUAL,
|   `b` double GENERATED ALWAYS AS (`a`) STORED
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
| CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) VIRTUAL, KEY (b))
< ok 0 0
> exec
| ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL
< error
| 1901 HY000 Function or expression 'a' cannot be used in the GENERATED ALWAYS AS clause of `b`
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
|   `a` double GENERATED ALWAYS AS (1) VIRTUAL,
|   `b` double GENERATED ALWAYS AS (`a`) VIRTUAL,
|   KEY `b` (`b`)
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
| CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) VIRTUAL, c DOUBLE AS (b) STORED)
< ok 0 0
> exec
| ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL
< error
| 1901 HY000 Function or expression 'b' cannot be used in the GENERATED ALWAYS AS clause of `c`
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
|   `a` double GENERATED ALWAYS AS (1) VIRTUAL,
|   `b` double GENERATED ALWAYS AS (`a`) VIRTUAL,
|   `c` double GENERATED ALWAYS AS (`b`) STORED
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
| CREATE TABLE t (a DOUBLE AS (1) STORED)
< ok 0 0
> exec
| ALTER TABLE t MODIFY a DOUBLE AS (RAND()) STORED
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the GENERATED ALWAYS AS clause of `a`
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
|   `a` double GENERATED ALWAYS AS (1) STORED
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
| CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL
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
|   `a` double GENERATED ALWAYS AS (rand()) VIRTUAL,
|   `b` double GENERATED ALWAYS AS (`a`) VIRTUAL
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
| CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) STORED)
< ok 0 0
> exec
| ALTER TABLE t MODIFY a DOUBLE AS (2) VIRTUAL
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
|   `a` double GENERATED ALWAYS AS (2) VIRTUAL,
|   `b` double GENERATED ALWAYS AS (`a`) STORED
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
| CREATE TABLE t (a DOUBLE AS (1) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL
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
|   `a` double GENERATED ALWAYS AS (rand()) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
