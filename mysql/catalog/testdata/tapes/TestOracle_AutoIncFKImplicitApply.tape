-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_AutoIncFKImplicitApply default#1

== TestOracle_AutoIncFKImplicitApply/8.0/autoinc-fk-implicit-only default#1
> exec
| DROP DATABASE IF EXISTS aifk_autoinc_fk_implicit_only
< ok 0 0
> exec
| CREATE DATABASE aifk_autoinc_fk_implicit_only DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE aifk_autoinc_fk_implicit_only
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL AUTO_INCREMENT, x INT, FOREIGN KEY (id) REFERENCES p(id)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE aifk_autoinc_fk_implicit_only.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE aifk_autoinc_fk_implicit_only.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `x` int DEFAULT NULL,
|   KEY `id` (`id`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`id`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS aifk_autoinc_fk_implicit_only
< ok 2 0
> exec
| CREATE DATABASE aifk_autoinc_fk_implicit_only DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE aifk_autoinc_fk_implicit_only
< ok 0 0
> exec
| CREATE TABLE `aifk_autoinc_fk_implicit_only`.`c` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `x` int DEFAULT NULL,
|   KEY `id` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| CREATE TABLE `aifk_autoinc_fk_implicit_only`.`p` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `aifk_autoinc_fk_implicit_only`.`c` ADD CONSTRAINT `c_ibfk_1` FOREIGN KEY (`id`) REFERENCES `p` (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE aifk_autoinc_fk_implicit_only.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE aifk_autoinc_fk_implicit_only.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `x` int DEFAULT NULL,
|   KEY `id` (`id`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`id`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
