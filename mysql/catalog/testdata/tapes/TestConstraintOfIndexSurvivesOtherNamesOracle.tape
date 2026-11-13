-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestConstraintOfIndexSurvivesOtherNamesOracle default#1
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a), CONSTRAINT i CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE t DROP INDEX i
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
|   `a` int DEFAULT NULL,
|   CONSTRAINT `i` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a), CONSTRAINT i CHECK (a > 0))
< ok 0 0
> exec
| DROP INDEX i ON t
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
|   `a` int DEFAULT NULL,
|   CONSTRAINT `i` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, KEY j (a), KEY i (b), CONSTRAINT i FOREIGN KEY (a) REFERENCES p (id))
< ok 0 0
> exec
| DROP INDEX i ON t
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
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `j` (`a`),
|   CONSTRAINT `i` FOREIGN KEY (`a`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT i CHECK (a > 0), UNIQUE KEY i (a))
< ok 0 0
> exec
| DROP INDEX i ON t
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
|   `a` int DEFAULT NULL,
|   CONSTRAINT `i` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
