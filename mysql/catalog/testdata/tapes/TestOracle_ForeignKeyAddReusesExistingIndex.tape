-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_ForeignKeyAddReusesExistingIndex default#1
> exec
| DROP DATABASE IF EXISTS fkreuse_f
< ok 0 0
> exec
| CREATE DATABASE fkreuse_f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkreuse_f
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid))
< ok 0 0
> query
| SHOW CREATE TABLE fkreuse_f.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`),
|   UNIQUE KEY `ub` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkreuse_f.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `pid` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `my_idx` (`pid`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS fkreuse_t
< ok 0 0
> exec
| CREATE DATABASE fkreuse_t DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkreuse_t
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkreuse_t.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`),
|   UNIQUE KEY `ub` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkreuse_t.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `pid` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `my_idx` (`pid`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
