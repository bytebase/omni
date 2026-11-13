-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_IndexFKColumnChangeNoIndexOps default#1
> exec
| DROP DATABASE IF EXISTS idxfkc_a
< ok 0 0
> exec
| CREATE DATABASE idxfkc_a DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfkc_a
< ok 0 0
> exec
| CREATE TABLE p (a INT NOT NULL PRIMARY KEY, b INT NOT NULL, UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (x) REFERENCES p(a))
< ok 0 0
> query
| SHOW CREATE TABLE idxfkc_a.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`a`),
|   UNIQUE KEY `ub` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE idxfkc_a.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk` (`x`),
|   CONSTRAINT `fk` FOREIGN KEY (`x`) REFERENCES `p` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idxfkc_b
< ok 0 0
> exec
| CREATE DATABASE idxfkc_b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfkc_b
< ok 0 0
> exec
| CREATE TABLE p (a INT NOT NULL PRIMARY KEY, b INT NOT NULL, UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (y) REFERENCES p(b))
< ok 0 0
> query
| SHOW CREATE TABLE idxfkc_b.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`a`),
|   UNIQUE KEY `ub` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE idxfkc_b.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk` (`y`),
|   CONSTRAINT `fk` FOREIGN KEY (`y`) REFERENCES `p` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
