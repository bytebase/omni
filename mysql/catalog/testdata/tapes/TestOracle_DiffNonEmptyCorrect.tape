-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_DiffNonEmptyCorrect default#1
> exec
| DROP DATABASE IF EXISTS diff_nf
< ok 0 0
> exec
| CREATE DATABASE diff_nf
< ok 1 0
> exec
| USE diff_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS diff_nt
< ok 0 0
> exec
| CREATE DATABASE diff_nt
< ok 1 0
> exec
| USE diff_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE diff_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS diff_nf
< ok 1 0
> exec
| CREATE DATABASE diff_nf
< ok 1 0
> exec
| USE diff_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE diff_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS diff_nt
< ok 1 0
> exec
| CREATE DATABASE diff_nt
< ok 1 0
> exec
| USE diff_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS diff_nf
< ok 1 0
> exec
| CREATE DATABASE diff_nf
< ok 1 0
> exec
| USE diff_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS diff_nt
< ok 1 0
> exec
| CREATE DATABASE diff_nt
< ok 1 0
> exec
| USE diff_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v BIGINT)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` bigint DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS diff_nf
< ok 1 0
> exec
| CREATE DATABASE diff_nf
< ok 1 0
> exec
| USE diff_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT '0',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS diff_nt
< ok 1 0
> exec
| CREATE DATABASE diff_nt
< ok 1 0
> exec
| USE diff_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 1)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT '1',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS diff_nf
< ok 1 0
> exec
| CREATE DATABASE diff_nf
< ok 1 0
> exec
| USE diff_nf
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS diff_nt
< ok 1 0
> exec
| CREATE DATABASE diff_nt
< ok 1 0
> exec
| USE diff_nt
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE diff_nt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `v` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
