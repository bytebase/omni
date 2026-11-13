-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestScenario_C14 default#1
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
| CREATE TABLE t (a INT, CHECK (a > 0))
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='information_schema' AND TABLE_NAME='TABLE_CONSTRAINTS'
|                AND COLUMN_NAME='ENFORCED'
< result
< column "COLUMN_NAME" "VARCHAR"
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
|   CONSTRAINT `CONSTRAINT_1` CHECK (`a` > 0)
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
| CREATE TABLE t (a INT, CONSTRAINT c_pos CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE t ALTER CHECK c_pos NOT ENFORCED
< error
| 1064 42000 …near 'CHECK c_pos NOT ENFORCED' at line 1
> exec
| ALTER TABLE t ALTER CHECK c_pos ENFORCED
< error
| 1064 42000 …near 'CHECK c_pos ENFORCED' at line 1
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
| CREATE TABLE t (
|     a INT,
|     g INT AS (a * 2) STORED,
|     CONSTRAINT c_g_nonneg CHECK (g >= 0)
| )
< ok 0 0
> query
| SELECT CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME='c_g_nonneg'
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| c_g_nonneg
: bytes
| `g` >= 0
> exec
| INSERT INTO t (a) VALUES (5)
< ok 1 0
> exec
| INSERT INTO t (a) VALUES (-1)
< error
| 4025 23000 CONSTRAINT `c_g_nonneg` failed for `testdb`.`t`
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
| CREATE TABLE t1 (a INT, CHECK (CHAR_LENGTH(CAST(a AS CHAR)) < 10))
< ok 0 0
> exec
| CREATE TABLE other (id INT)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, CHECK (a IN (SELECT id FROM other)))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec
| CREATE TABLE t3 (a INT, CHECK (a < NOW()))
< error
| 1901 HY000 Function or expression 'current_timestamp()' cannot be used in the CHECK clause of `CONSTRAINT_1`
> exec
| CREATE TABLE t4 (a INT, CHECK (a < @max))
< error
| 1901 HY000 Function or expression '@max' cannot be used in the CHECK clause of `CONSTRAINT_1`
> exec
| CREATE TABLE t5 (a INT, CHECK (a < RAND()))
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the CHECK clause of `CONSTRAINT_1`
