-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_RejectedAlterLeavesDependenciesAlone default#1
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
| ALTER TABLE t RENAME COLUMN a TO z, ADD COLUMN b INT
< error
| 1060 42S21 Duplicate column name 'b'
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
| ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL
< error
| 1901 HY000 Function or expression 'a' cannot be used in the GENERATED ALWAYS AS clause of `b`
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
| ALTER TABLE t ADD COLUMN b INT, RENAME COLUMN a TO z
< error
| 1060 42S21 Duplicate column name 'b'
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
| ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL
< error
| 1901 HY000 Function or expression 'a' cannot be used in the GENERATED ALWAYS AS clause of `b`
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
| CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, b DOUBLE AS (a) VIRTUAL)
< ok 0 0
> exec
| ALTER TABLE t RENAME COLUMN a TO z, ADD COLUMN b INT
< error
| 1060 42S21 Duplicate column name 'b'
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
| CREATE INDEX i ON t(b)
< error
| 1901 HY000 Function or expression 'a' cannot be used in the GENERATED ALWAYS AS clause of `b`
