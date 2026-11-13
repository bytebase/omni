-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_IndexDiffFKImplicit default#1
> exec
| DROP DATABASE IF EXISTS idxfk_fk_named_implicit
< ok 0 0
> exec
| CREATE DATABASE idxfk_fk_named_implicit DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfk_fk_named_implicit
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE idxfk_fk_named_implicit.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE idxfk_fk_named_implicit.c
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
|   KEY `fk` (`pid`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idxfk_fk_unnamed_implicit
< ok 0 0
> exec
| CREATE DATABASE idxfk_fk_unnamed_implicit DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfk_fk_unnamed_implicit
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE idxfk_fk_unnamed_implicit.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE idxfk_fk_unnamed_implicit.c
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
|   KEY `pid` (`pid`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idxfk_fk_explicit_covering
< ok 0 0
> exec
| CREATE DATABASE idxfk_fk_explicit_covering DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfk_fk_explicit_covering
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE idxfk_fk_explicit_covering.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE idxfk_fk_explicit_covering.c
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
> exec
| DROP DATABASE IF EXISTS idxfk_fk_composite
< ok 0 0
> exec
| CREATE DATABASE idxfk_fk_composite DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfk_fk_composite
< ok 0 0
> exec
| CREATE TABLE p (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (x,y) REFERENCES p(a,b))
< ok 0 0
> query
| SHOW CREATE TABLE idxfk_fk_composite.p
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
|   PRIMARY KEY (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE idxfk_fk_composite.c
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
|   KEY `fk` (`x`,`y`),
|   CONSTRAINT `fk` FOREIGN KEY (`x`, `y`) REFERENCES `p` (`a`, `b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idxfk_fk_backing_name_collision
< ok 0 0
> exec
| CREATE DATABASE idxfk_fk_backing_name_collision DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfk_fk_backing_name_collision
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, a INT, pid INT, KEY pid (a), FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE idxfk_fk_backing_name_collision.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE idxfk_fk_backing_name_collision.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `pid` (`a`),
|   KEY `pid_2` (`pid`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idxfk_fk_user_named_ibfk
< ok 0 0
> exec
| CREATE DATABASE idxfk_fk_user_named_ibfk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfk_fk_user_named_ibfk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT c_ibfk_1 FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE idxfk_fk_user_named_ibfk.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE idxfk_fk_user_named_ibfk.c
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
|   KEY `c_ibfk_1` (`pid`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
