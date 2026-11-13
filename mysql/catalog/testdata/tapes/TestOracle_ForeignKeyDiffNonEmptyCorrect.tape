-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_ForeignKeyDiffNonEmptyCorrect default#1
> exec
| DROP DATABASE IF EXISTS fknec_f_add_fk
< ok 0 0
> exec
| CREATE DATABASE fknec_f_add_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_f_add_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SHOW CREATE TABLE fknec_f_add_fk.p
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
| SHOW CREATE TABLE fknec_f_add_fk.c
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS fknec_t_add_fk
< ok 0 0
> exec
| CREATE DATABASE fknec_t_add_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_t_add_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fknec_t_add_fk.p
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
| SHOW CREATE TABLE fknec_t_add_fk.c
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
| DROP DATABASE IF EXISTS fknec_f_drop_fk
< ok 0 0
> exec
| CREATE DATABASE fknec_f_drop_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_f_drop_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fknec_f_drop_fk.p
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
| SHOW CREATE TABLE fknec_f_drop_fk.c
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
| DROP DATABASE IF EXISTS fknec_t_drop_fk
< ok 0 0
> exec
| CREATE DATABASE fknec_t_drop_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_t_drop_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SHOW CREATE TABLE fknec_t_drop_fk.p
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
| SHOW CREATE TABLE fknec_t_drop_fk.c
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS fknec_f_modify_action
< ok 0 0
> exec
| CREATE DATABASE fknec_f_modify_action DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_f_modify_action
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE CASCADE)
< ok 0 0
> query
| SHOW CREATE TABLE fknec_f_modify_action.p
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
| SHOW CREATE TABLE fknec_f_modify_action.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS fknec_t_modify_action
< ok 0 0
> exec
| CREATE DATABASE fknec_t_modify_action DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_t_modify_action
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE SET NULL)
< ok 0 0
> query
| SHOW CREATE TABLE fknec_t_modify_action.p
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
| SHOW CREATE TABLE fknec_t_modify_action.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS fknec_f_modify_refcol
< ok 0 0
> exec
| CREATE DATABASE fknec_f_modify_refcol DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_f_modify_refcol
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(a))
< ok 0 0
> query
| SHOW CREATE TABLE fknec_f_modify_refcol.p
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
| SHOW CREATE TABLE fknec_f_modify_refcol.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS fknec_t_modify_refcol
< ok 0 0
> exec
| CREATE DATABASE fknec_t_modify_refcol DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_t_modify_refcol
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(b))
< ok 0 0
> query
| SHOW CREATE TABLE fknec_t_modify_refcol.p
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
| SHOW CREATE TABLE fknec_t_modify_refcol.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
