-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_UnaryBinaryInStoredExpressions default#1
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
| CREATE TABLE t (a INT, CHECK (BINARY 'a' = 'A'))
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
|   CONSTRAINT `CONSTRAINT_1` CHECK (cast('a' as char charset binary) = 'A')
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `CONSTRAINT_1` CHECK (cast('a' as char charset binary) = 'A')
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   CONSTRAINT `CONSTRAINT_1` CHECK (cast('a' as char charset binary) = 'A')
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
| CREATE TABLE t (a INT DEFAULT (2 * BINARY (1 + 2)))
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
|   `a` int(11) DEFAULT (2 * cast(1 + 2 as char charset binary))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT (2 * cast(1 + 2 as char charset binary))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `a` int(11) DEFAULT (2 * cast(1 + 2 as char charset binary))
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
| CREATE TABLE t (a VARCHAR(5), b INT AS (BINARY (1 + 2)) VIRTUAL)
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
|   `a` varchar(5) DEFAULT NULL,
|   `b` int(11) GENERATED ALWAYS AS (cast(1 + 2 as char charset binary)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varchar(5) DEFAULT NULL,
|   `b` int(11) GENERATED ALWAYS AS (cast(1 + 2 as char charset binary)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `a` varchar(5) DEFAULT NULL,
|   `b` int(11) GENERATED ALWAYS AS (cast(1 + 2 as char charset binary)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
