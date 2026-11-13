-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_CheckDiffNonEmptyCorrect default#1
> exec
| DROP DATABASE IF EXISTS chk_nf
< ok 0 0
> exec
| CREATE DATABASE chk_nf
< ok 1 0
> exec
| USE chk_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE chk_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_nt
< ok 0 0
> exec
| CREATE DATABASE chk_nt
< ok 1 0
> exec
| USE chk_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_nf
< ok 1 0
> exec
| CREATE DATABASE chk_nf
< ok 1 0
> exec
| USE chk_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_nt
< ok 1 0
> exec
| CREATE DATABASE chk_nt
< ok 1 0
> exec
| USE chk_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE chk_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_nf
< ok 1 0
> exec
| CREATE DATABASE chk_nf
< ok 1 0
> exec
| USE chk_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_nt
< ok 1 0
> exec
| CREATE DATABASE chk_nt
< ok 1 0
> exec
| USE chk_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 10))
< ok 0 0
> query
| SHOW CREATE TABLE chk_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck` CHECK ((`a` > 10))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_nf
< ok 1 0
> exec
| CREATE DATABASE chk_nf
< ok 1 0
> exec
| USE chk_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b))
< ok 0 0
> query
| SHOW CREATE TABLE chk_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck` CHECK ((`a` < `b`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_nt
< ok 1 0
> exec
| CREATE DATABASE chk_nt
< ok 1 0
> exec
| USE chk_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE chk_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
