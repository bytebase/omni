-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_AutoIncGroupingFKImplicitBacking default#1
> exec
| DROP DATABASE IF EXISTS fkmulti_f_aig_fk_implicit_backing
< ok 0 0
> exec
| CREATE DATABASE fkmulti_f_aig_fk_implicit_backing DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_f_aig_fk_implicit_backing
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, x INT)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_f_aig_fk_implicit_backing' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_f_aig_fk_implicit_backing.c
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_f_aig_fk_implicit_backing.p
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
> exec
| DROP DATABASE IF EXISTS fkmulti_t_aig_fk_implicit_backing
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_aig_fk_implicit_backing DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_aig_fk_implicit_backing
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, x INT, seq INT NOT NULL AUTO_INCREMENT, CONSTRAINT fk_seq FOREIGN KEY (seq) REFERENCES p (id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_aig_fk_implicit_backing' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_aig_fk_implicit_backing.c
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
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `fk_seq` (`seq`),
|   CONSTRAINT `fk_seq` FOREIGN KEY (`seq`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_t_aig_fk_implicit_backing.p
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
> exec
| DROP DATABASE IF EXISTS fkmulti_f_aig_fk_implicit_backing_drop
< ok 0 0
> exec
| CREATE DATABASE fkmulti_f_aig_fk_implicit_backing_drop DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_f_aig_fk_implicit_backing_drop
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, x INT, seq INT NOT NULL AUTO_INCREMENT, CONSTRAINT fk_seq FOREIGN KEY (seq) REFERENCES p (id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_f_aig_fk_implicit_backing_drop' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_f_aig_fk_implicit_backing_drop.c
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
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `fk_seq` (`seq`),
|   CONSTRAINT `fk_seq` FOREIGN KEY (`seq`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_f_aig_fk_implicit_backing_drop.p
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
> exec
| DROP DATABASE IF EXISTS fkmulti_t_aig_fk_implicit_backing_drop
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_aig_fk_implicit_backing_drop DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_aig_fk_implicit_backing_drop
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, x INT)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_aig_fk_implicit_backing_drop' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_aig_fk_implicit_backing_drop.c
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_t_aig_fk_implicit_backing_drop.p
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
> exec
| DROP DATABASE IF EXISTS fkmulti_f_aig_fk_rename_in_place_to_ai
< ok 0 0
> exec
| CREATE DATABASE fkmulti_f_aig_fk_rename_in_place_to_ai DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_f_aig_fk_rename_in_place_to_ai
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, old_id INT, CONSTRAINT fk FOREIGN KEY (old_id) REFERENCES p (id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_f_aig_fk_rename_in_place_to_ai' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_f_aig_fk_rename_in_place_to_ai.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `old_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk` (`old_id`),
|   CONSTRAINT `fk` FOREIGN KEY (`old_id`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_f_aig_fk_rename_in_place_to_ai.p
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
> exec
| DROP DATABASE IF EXISTS fkmulti_t_aig_fk_rename_in_place_to_ai
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_aig_fk_rename_in_place_to_ai DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_aig_fk_rename_in_place_to_ai
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, old_id INT, seq INT NOT NULL AUTO_INCREMENT, CONSTRAINT fk FOREIGN KEY (seq) REFERENCES p (id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_aig_fk_rename_in_place_to_ai' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_aig_fk_rename_in_place_to_ai.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `old_id` int DEFAULT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `fk` (`seq`),
|   CONSTRAINT `fk` FOREIGN KEY (`seq`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_t_aig_fk_rename_in_place_to_ai.p
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

== TestOracle_AutoIncGroupingFKImplicitBacking/8.0/aig-fk-implicit-backing default#1
> exec
| DROP DATABASE IF EXISTS fkmulti
< ok 0 0
> exec
| CREATE DATABASE fkmulti DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, x INT)
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT, ADD KEY `fk_seq` (`seq`)
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` ADD CONSTRAINT `fk_seq` FOREIGN KEY (`seq`) REFERENCES `p` (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE fkmulti.p
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
| SHOW CREATE TABLE fkmulti.c
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
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `fk_seq` (`seq`),
|   CONSTRAINT `fk_seq` FOREIGN KEY (`seq`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingFKImplicitBacking/8.0/aig-fk-implicit-backing-drop default#1
> exec
| DROP DATABASE IF EXISTS fkmulti
< ok 2 0
> exec
| CREATE DATABASE fkmulti DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, x INT, seq INT NOT NULL AUTO_INCREMENT, CONSTRAINT fk_seq FOREIGN KEY (seq) REFERENCES p (id))
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fk_seq`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP INDEX `fk_seq`, DROP COLUMN `seq`
< ok 0 0
> query
| SHOW CREATE TABLE fkmulti.p
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
| SHOW CREATE TABLE fkmulti.c
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_AutoIncGroupingFKImplicitBacking/8.0/aig-fk-rename-in-place-to-ai default#1
> exec
| DROP DATABASE IF EXISTS fkmulti
< ok 2 0
> exec
| CREATE DATABASE fkmulti DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, old_id INT, CONSTRAINT fk FOREIGN KEY (old_id) REFERENCES p (id))
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP INDEX `fk`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT, ADD KEY `fk` (`seq`)
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` ADD CONSTRAINT `fk` FOREIGN KEY (`seq`) REFERENCES `p` (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE fkmulti.p
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
| SHOW CREATE TABLE fkmulti.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `old_id` int DEFAULT NULL,
|   `seq` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   KEY `fk` (`seq`),
|   CONSTRAINT `fk` FOREIGN KEY (`seq`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
