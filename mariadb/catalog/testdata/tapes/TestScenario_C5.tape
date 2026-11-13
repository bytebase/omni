-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestScenario_C5 default#1
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
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id))
< ok 0 0
> query
| SELECT DELETE_RULE, UPDATE_RULE, MATCH_OPTION
|             FROM information_schema.REFERENTIAL_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='c'
< result
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
< row
: bytes
| RESTRICT
: bytes
| RESTRICT
: bytes
| NONE
> query
| SHOW CREATE TABLE c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `a` int(11) DEFAULT NULL,
|   KEY `a` (`a`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`a`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE c (
|             a INT NOT NULL,
|             FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET NULL
|         )
< error
| 1005 HY000 Can't create table `testdb`.`c` (errno: 150 "Foreign key constraint is incorrectly formed")
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
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id))
< ok 0 0
> query
| SELECT MATCH_OPTION FROM information_schema.REFERENTIAL_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='c'
< result
< column "MATCH_OPTION" "VARCHAR"
< row
: bytes
| NONE
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id) ON DELETE CASCADE)
< ok 0 0
> query
| SELECT DELETE_RULE, UPDATE_RULE
|             FROM information_schema.REFERENTIAL_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='c'
< result
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< row
: bytes
| CASCADE
: bytes
| RESTRICT
> query
| SHOW CREATE TABLE c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `a` int(11) DEFAULT NULL,
|   KEY `a` (`a`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`a`) REFERENCES `p` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE c (
|             a INT DEFAULT 0,
|             FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET DEFAULT
|         )
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (id BIGINT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id))
< error
| 1005 HY000 Can't create table `testdb`.`c` (errno: 150 "Foreign key constraint is incorrectly formed")
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
| CREATE TABLE p (id INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE c (
|             a INT,
|             b INT GENERATED ALWAYS AS (a+1) VIRTUAL,
|             FOREIGN KEY (b) REFERENCES p(id)
|         )
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
|             a INT,
|             CONSTRAINT chk_pos CHECK (a > 0)
|         )
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
|   CONSTRAINT `chk_pos` CHECK (`a` > 0)
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
| CREATE TABLE t1 (a INT CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, CHECK (a > 0))
< ok 0 0
> query
| SELECT CONSTRAINT_NAME, CHECK_CLAUSE
|             FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='t1'
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| a
: bytes
| `a` > 0
> query
| SELECT CONSTRAINT_NAME, CHECK_CLAUSE
|             FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='t2'
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| CONSTRAINT_1
: bytes
| `a` > 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT CHECK (a > b), b INT)
< ok 0 0
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
| CREATE TABLE t (a INT, b INT, CHECK (a > b))
< ok 0 0
