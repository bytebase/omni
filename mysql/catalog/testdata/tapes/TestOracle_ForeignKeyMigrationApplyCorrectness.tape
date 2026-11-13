-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_ForeignKeyMigrationApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS fkapply_f_add_fk_autoindex
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_add_fk_autoindex DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_add_fk_autoindex
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_add_fk_autoindex.p
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
| SHOW CREATE TABLE fkapply_f_add_fk_autoindex.c
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
| DROP DATABASE IF EXISTS fkapply_t_add_fk_autoindex
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_add_fk_autoindex DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_add_fk_autoindex
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_add_fk_autoindex.p
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
| SHOW CREATE TABLE fkapply_t_add_fk_autoindex.c
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
> query
| SHOW CREATE TABLE fkapply.p
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
> exec
| DROP DATABASE IF EXISTS fkapply_f_add_fk_reuse_user_index
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_add_fk_reuse_user_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_add_fk_reuse_user_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_add_fk_reuse_user_index.p
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
| SHOW CREATE TABLE fkapply_f_add_fk_reuse_user_index.c
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
| DROP DATABASE IF EXISTS fkapply_t_add_fk_reuse_user_index
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_add_fk_reuse_user_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_add_fk_reuse_user_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_add_fk_reuse_user_index.p
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
| SHOW CREATE TABLE fkapply_t_add_fk_reuse_user_index.c
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
> query
| SHOW CREATE TABLE fkapply.p
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
> exec
| DROP DATABASE IF EXISTS fkapply_f_add_fk_cascade
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_add_fk_cascade DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_add_fk_cascade
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_add_fk_cascade.p
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
| SHOW CREATE TABLE fkapply_f_add_fk_cascade.c
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
| DROP DATABASE IF EXISTS fkapply_t_add_fk_cascade
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_add_fk_cascade DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_add_fk_cascade
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE CASCADE ON UPDATE SET NULL)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_add_fk_cascade.p
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
| SHOW CREATE TABLE fkapply_t_add_fk_cascade.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`) ON DELETE CASCADE ON UPDATE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkapply.p
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
> exec
| DROP DATABASE IF EXISTS fkapply_f_add_fk_unnamed
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_add_fk_unnamed DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_add_fk_unnamed
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_add_fk_unnamed.p
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
| SHOW CREATE TABLE fkapply_f_add_fk_unnamed.c
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
| DROP DATABASE IF EXISTS fkapply_t_add_fk_unnamed
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_add_fk_unnamed DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_add_fk_unnamed
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_add_fk_unnamed.p
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
| SHOW CREATE TABLE fkapply_t_add_fk_unnamed.c
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
> query
| SHOW CREATE TABLE fkapply.p
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
> exec
| DROP DATABASE IF EXISTS fkapply_f_add_fk_composite
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_add_fk_composite DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_add_fk_composite
< ok 0 0
> exec
| CREATE TABLE p (a INT NOT NULL, b INT NOT NULL, c INT, PRIMARY KEY (a,b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_add_fk_composite.p
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
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkapply_f_add_fk_composite.c
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS fkapply_t_add_fk_composite
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_add_fk_composite DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_add_fk_composite
< ok 0 0
> exec
| CREATE TABLE p (a INT NOT NULL, b INT NOT NULL, c INT, PRIMARY KEY (a,b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (x,y) REFERENCES p(a,b))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_add_fk_composite.p
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
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkapply_t_add_fk_composite.c
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
> query
| SHOW CREATE TABLE fkapply.p
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
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS fkapply_f_drop_fk_and_backing_index
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_drop_fk_and_backing_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_drop_fk_and_backing_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_drop_fk_and_backing_index.p
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
| SHOW CREATE TABLE fkapply_f_drop_fk_and_backing_index.c
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
| DROP DATABASE IF EXISTS fkapply_t_drop_fk_and_backing_index
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_drop_fk_and_backing_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_drop_fk_and_backing_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_drop_fk_and_backing_index.p
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
| SHOW CREATE TABLE fkapply_t_drop_fk_and_backing_index.c
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
> query
| SHOW CREATE TABLE fkapply.p
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
> exec
| DROP DATABASE IF EXISTS fkapply_f_drop_fk_keep_user_index
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_drop_fk_keep_user_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_drop_fk_keep_user_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_drop_fk_keep_user_index.p
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
| SHOW CREATE TABLE fkapply_f_drop_fk_keep_user_index.c
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
| DROP DATABASE IF EXISTS fkapply_t_drop_fk_keep_user_index
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_drop_fk_keep_user_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_drop_fk_keep_user_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_drop_fk_keep_user_index.p
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
| SHOW CREATE TABLE fkapply_t_drop_fk_keep_user_index.c
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
> query
| SHOW CREATE TABLE fkapply.p
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
> exec
| DROP DATABASE IF EXISTS fkapply_f_drop_fk_unnamed
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_drop_fk_unnamed DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_drop_fk_unnamed
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_drop_fk_unnamed.p
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
| SHOW CREATE TABLE fkapply_f_drop_fk_unnamed.c
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
| DROP DATABASE IF EXISTS fkapply_t_drop_fk_unnamed
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_drop_fk_unnamed DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_drop_fk_unnamed
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_drop_fk_unnamed.p
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
| SHOW CREATE TABLE fkapply_t_drop_fk_unnamed.c
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
> query
| SHOW CREATE TABLE fkapply.p
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
> exec
| DROP DATABASE IF EXISTS fkapply_f_drop_fk_two_covering_indexes
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_drop_fk_two_covering_indexes DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_drop_fk_two_covering_indexes
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY pid (pid), KEY pid_2 (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_drop_fk_two_covering_indexes.p
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
| SHOW CREATE TABLE fkapply_f_drop_fk_two_covering_indexes.c
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
|   KEY `pid_2` (`pid`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS fkapply_t_drop_fk_two_covering_indexes
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_drop_fk_two_covering_indexes DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_drop_fk_two_covering_indexes
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY pid (pid))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_drop_fk_two_covering_indexes.p
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
| SHOW CREATE TABLE fkapply_t_drop_fk_two_covering_indexes.c
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
|   KEY `pid` (`pid`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkapply.p
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
> exec
| DROP DATABASE IF EXISTS fkapply_f_modify_fk_action
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_modify_fk_action DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_modify_fk_action
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE CASCADE)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_modify_fk_action.p
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
| SHOW CREATE TABLE fkapply_f_modify_fk_action.c
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
| DROP DATABASE IF EXISTS fkapply_t_modify_fk_action
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_modify_fk_action DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_modify_fk_action
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE SET NULL)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_modify_fk_action.p
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
| SHOW CREATE TABLE fkapply_t_modify_fk_action.c
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
> query
| SHOW CREATE TABLE fkapply.p
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
> exec
| DROP DATABASE IF EXISTS fkapply_f_modify_fk_refcol
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_modify_fk_refcol DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_modify_fk_refcol
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(a))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_modify_fk_refcol.p
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
| SHOW CREATE TABLE fkapply_f_modify_fk_refcol.c
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
| DROP DATABASE IF EXISTS fkapply_t_modify_fk_refcol
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_modify_fk_refcol DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_modify_fk_refcol
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(b))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_modify_fk_refcol.p
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
| SHOW CREATE TABLE fkapply_t_modify_fk_refcol.c
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
> query
| SHOW CREATE TABLE fkapply.p
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
> exec
| DROP DATABASE IF EXISTS fkapply_f_modify_fk_column
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_modify_fk_column DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_modify_fk_column
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (x) REFERENCES p(a))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_modify_fk_column.p
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
| SHOW CREATE TABLE fkapply_f_modify_fk_column.c
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
| DROP DATABASE IF EXISTS fkapply_t_modify_fk_column
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_modify_fk_column DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_modify_fk_column
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (y) REFERENCES p(b))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_modify_fk_column.p
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
| SHOW CREATE TABLE fkapply_t_modify_fk_column.c
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
> query
| SHOW CREATE TABLE fkapply.p
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
> exec
| DROP DATABASE IF EXISTS fkapply_f_self_referential
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_self_referential DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_self_referential
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, parent_id INT)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_self_referential.p
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
| SHOW CREATE TABLE fkapply_f_self_referential.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS fkapply_t_self_referential
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_self_referential DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_self_referential
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, parent_id INT, CONSTRAINT fk_self FOREIGN KEY (parent_id) REFERENCES c(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_self_referential.p
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
| SHOW CREATE TABLE fkapply_t_self_referential.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_self` (`parent_id`),
|   CONSTRAINT `fk_self` FOREIGN KEY (`parent_id`) REFERENCES `c` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkapply.p
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
> exec
| DROP DATABASE IF EXISTS fkapply_f_add_second_fk
< ok 0 0
> exec
| CREATE DATABASE fkapply_f_add_second_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_f_add_second_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pa INT, pb INT, CONSTRAINT fka FOREIGN KEY (pa) REFERENCES p(a))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_f_add_second_fk.p
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
| SHOW CREATE TABLE fkapply_f_add_second_fk.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `pa` int DEFAULT NULL,
|   `pb` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fka` (`pa`),
|   CONSTRAINT `fka` FOREIGN KEY (`pa`) REFERENCES `p` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS fkapply_t_add_second_fk
< ok 0 0
> exec
| CREATE DATABASE fkapply_t_add_second_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply_t_add_second_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pa INT, pb INT, CONSTRAINT fka FOREIGN KEY (pa) REFERENCES p(a), CONSTRAINT fkb FOREIGN KEY (pb) REFERENCES p(b))
< ok 0 0
> query
| SHOW CREATE TABLE fkapply_t_add_second_fk.p
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
| SHOW CREATE TABLE fkapply_t_add_second_fk.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `pa` int DEFAULT NULL,
|   `pb` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fka` (`pa`),
|   KEY `fkb` (`pb`),
|   CONSTRAINT `fka` FOREIGN KEY (`pa`) REFERENCES `p` (`a`),
|   CONSTRAINT `fkb` FOREIGN KEY (`pb`) REFERENCES `p` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkapply.p
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

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/add-fk-autoindex default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 0 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` ADD CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
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

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/add-fk-cascade default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` ADD CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`) ON DELETE CASCADE ON UPDATE SET NULL
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`) ON DELETE CASCADE ON UPDATE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/add-fk-composite default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (a INT NOT NULL, b INT NOT NULL, c INT, PRIMARY KEY (a,b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT)
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` ADD CONSTRAINT `fk` FOREIGN KEY (`x`, `y`) REFERENCES `p` (`a`, `b`)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
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

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/add-fk-reuse-user-index default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid))
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` ADD CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
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

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/add-fk-unnamed default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` ADD CONSTRAINT `c_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
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

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/add-second-fk default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pa INT, pb INT, CONSTRAINT fka FOREIGN KEY (pa) REFERENCES p(a))
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` ADD CONSTRAINT `fkb` FOREIGN KEY (`pb`) REFERENCES `p` (`b`)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `pa` int DEFAULT NULL,
|   `pb` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fka` (`pa`),
|   KEY `fkb` (`pb`),
|   CONSTRAINT `fka` FOREIGN KEY (`pa`) REFERENCES `p` (`a`),
|   CONSTRAINT `fkb` FOREIGN KEY (`pb`) REFERENCES `p` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/drop-fk-and-backing-index default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` DROP INDEX `fk`
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
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

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/drop-fk-keep-user-index default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
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

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/drop-fk-two-covering-indexes default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY pid (pid), KEY pid_2 (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` DROP INDEX `pid_2`
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
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
|   KEY `pid` (`pid`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/drop-fk-unnamed default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` DROP FOREIGN KEY `c_ibfk_1`
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` DROP INDEX `pid`
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
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

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/modify-fk-action default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE CASCADE)
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` ADD CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`) ON DELETE SET NULL
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
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

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/modify-fk-column default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (x) REFERENCES p(a))
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` DROP INDEX `fk`
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` ADD CONSTRAINT `fk` FOREIGN KEY (`y`) REFERENCES `p` (`b`)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
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

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/modify-fk-refcol default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(a))
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` ADD CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`b`)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
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

== TestOracle_ForeignKeyMigrationApplyCorrectness/8.0/self-referential default#1
> exec
| DROP DATABASE IF EXISTS fkapply
< ok 2 0
> exec
| CREATE DATABASE fkapply DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkapply
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, parent_id INT)
< ok 0 0
> exec
| ALTER TABLE `fkapply`.`c` ADD CONSTRAINT `fk_self` FOREIGN KEY (`parent_id`) REFERENCES `c` (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE fkapply.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_self` (`parent_id`),
|   CONSTRAINT `fk_self` FOREIGN KEY (`parent_id`) REFERENCES `c` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
