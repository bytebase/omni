-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_ForeignKeyMultiTableApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS fkmulti_f_drop_column_and_its_fk
< ok 0 0
> exec
| CREATE DATABASE fkmulti_f_drop_column_and_its_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_f_drop_column_and_its_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_f_drop_column_and_its_fk' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_f_drop_column_and_its_fk.c
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
| SHOW CREATE TABLE fkmulti_f_drop_column_and_its_fk.p
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
| DROP DATABASE IF EXISTS fkmulti_t_drop_column_and_its_fk
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_column_and_its_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_column_and_its_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_column_and_its_fk' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_drop_column_and_its_fk.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_t_drop_column_and_its_fk.p
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
| DROP DATABASE IF EXISTS fkmulti_f_drop_referenced_table
< ok 0 0
> exec
| CREATE DATABASE fkmulti_f_drop_referenced_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_f_drop_referenced_table
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_f_drop_referenced_table' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_f_drop_referenced_table.c
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
| SHOW CREATE TABLE fkmulti_f_drop_referenced_table.p
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
| DROP DATABASE IF EXISTS fkmulti_t_drop_referenced_table
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_referenced_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_referenced_table
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_referenced_table' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
> query
| SHOW CREATE TABLE fkmulti_t_drop_referenced_table.c
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
| DROP DATABASE IF EXISTS fkmulti_f_drop_two_fks_sharing_index
< ok 0 0
> exec
| CREATE DATABASE fkmulti_f_drop_two_fks_sharing_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_f_drop_two_fks_sharing_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fka FOREIGN KEY (pid) REFERENCES p(id), CONSTRAINT fkb FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_f_drop_two_fks_sharing_index' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_f_drop_two_fks_sharing_index.c
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
|   KEY `fkb` (`pid`),
|   CONSTRAINT `fka` FOREIGN KEY (`pid`) REFERENCES `p` (`id`),
|   CONSTRAINT `fkb` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_f_drop_two_fks_sharing_index.p
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
| DROP DATABASE IF EXISTS fkmulti_t_drop_two_fks_sharing_index
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_two_fks_sharing_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_two_fks_sharing_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_two_fks_sharing_index' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_drop_two_fks_sharing_index.c
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
| SHOW CREATE TABLE fkmulti_t_drop_two_fks_sharing_index.p
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
| DROP DATABASE IF EXISTS fkmulti_f_drop_one_of_two_fks_sharing_index
< ok 0 0
> exec
| CREATE DATABASE fkmulti_f_drop_one_of_two_fks_sharing_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_f_drop_one_of_two_fks_sharing_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fka FOREIGN KEY (pid) REFERENCES p(id), CONSTRAINT fkb FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_f_drop_one_of_two_fks_sharing_index' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_f_drop_one_of_two_fks_sharing_index.c
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
|   KEY `fkb` (`pid`),
|   CONSTRAINT `fka` FOREIGN KEY (`pid`) REFERENCES `p` (`id`),
|   CONSTRAINT `fkb` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_f_drop_one_of_two_fks_sharing_index.p
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
| DROP DATABASE IF EXISTS fkmulti_t_drop_one_of_two_fks_sharing_index
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_one_of_two_fks_sharing_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_one_of_two_fks_sharing_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fkb FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_one_of_two_fks_sharing_index' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_drop_one_of_two_fks_sharing_index.c
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
|   KEY `fkb` (`pid`),
|   CONSTRAINT `fkb` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_t_drop_one_of_two_fks_sharing_index.p
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
| DROP DATABASE IF EXISTS fkmulti_f_drop_fk_reuse_name_for_other_index
< ok 0 0
> exec
| CREATE DATABASE fkmulti_f_drop_fk_reuse_name_for_other_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_f_drop_fk_reuse_name_for_other_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, other_col INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_f_drop_fk_reuse_name_for_other_index' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_f_drop_fk_reuse_name_for_other_index.c
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
|   `other_col` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk` (`pid`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_f_drop_fk_reuse_name_for_other_index.p
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
| DROP DATABASE IF EXISTS fkmulti_t_drop_fk_reuse_name_for_other_index
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_fk_reuse_name_for_other_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_fk_reuse_name_for_other_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, other_col INT, KEY fk (other_col))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_fk_reuse_name_for_other_index' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_drop_fk_reuse_name_for_other_index.c
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
|   `other_col` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk` (`other_col`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_t_drop_fk_reuse_name_for_other_index.p
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
| DROP DATABASE IF EXISTS fkmulti_f_drop_parent_and_child_name_order
< ok 0 0
> exec
| CREATE DATABASE fkmulti_f_drop_parent_and_child_name_order DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_f_drop_parent_and_child_name_order
< ok 0 0
> exec
| CREATE TABLE a_parent (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE z_child (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES a_parent(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_f_drop_parent_and_child_name_order' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| a_parent
< row
: bytes
| z_child
> query
| SHOW CREATE TABLE fkmulti_f_drop_parent_and_child_name_order.a_parent
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a_parent
: bytes
| CREATE TABLE `a_parent` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_f_drop_parent_and_child_name_order.z_child
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| z_child
: bytes
| CREATE TABLE `z_child` (
|   `id` int NOT NULL,
|   `pid` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk` (`pid`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `a_parent` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS fkmulti_t_drop_parent_and_child_name_order
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_parent_and_child_name_order DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_parent_and_child_name_order
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_parent_and_child_name_order' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS fkmulti_f_drop_two_unnamed_fks_sharing_index
< ok 0 0
> exec
| CREATE DATABASE fkmulti_f_drop_two_unnamed_fks_sharing_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_f_drop_two_unnamed_fks_sharing_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id), FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_f_drop_two_unnamed_fks_sharing_index' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_f_drop_two_unnamed_fks_sharing_index.c
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
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `p` (`id`),
|   CONSTRAINT `c_ibfk_2` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE fkmulti_f_drop_two_unnamed_fks_sharing_index.p
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
| DROP DATABASE IF EXISTS fkmulti_t_drop_two_unnamed_fks_sharing_index
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_two_unnamed_fks_sharing_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_two_unnamed_fks_sharing_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_two_unnamed_fks_sharing_index' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_drop_two_unnamed_fks_sharing_index.c
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
| SHOW CREATE TABLE fkmulti_t_drop_two_unnamed_fks_sharing_index.p
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

== TestOracle_ForeignKeyMultiTableApplyCorrectness/8.0/drop-column-and-its-fk default#1
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
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP INDEX `fk`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP COLUMN `pid`
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ForeignKeyMultiTableApplyCorrectness/8.0/drop-fk-reuse-name-for-other-index default#1
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
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, other_col INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP INDEX `fk`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` ADD KEY `fk` (`other_col`)
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
|   `pid` int DEFAULT NULL,
|   `other_col` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk` (`other_col`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ForeignKeyMultiTableApplyCorrectness/8.0/drop-one-of-two-fks-sharing-index default#1
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
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fka FOREIGN KEY (pid) REFERENCES p(id), CONSTRAINT fkb FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fka`
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
|   `pid` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fkb` (`pid`),
|   CONSTRAINT `fkb` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ForeignKeyMultiTableApplyCorrectness/8.0/drop-parent-and-child-name-order default#1
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
| CREATE TABLE a_parent (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE z_child (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES a_parent(id))
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`z_child` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| DROP TABLE `fkmulti`.`a_parent`
< ok 0 0
> exec
| DROP TABLE `fkmulti`.`z_child`
< ok 0 0

== TestOracle_ForeignKeyMultiTableApplyCorrectness/8.0/drop-referenced-table default#1
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
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| DROP TABLE `fkmulti`.`p`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP INDEX `fk`
< ok 0 0
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
|   `pid` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ForeignKeyMultiTableApplyCorrectness/8.0/drop-two-fks-sharing-index default#1
> exec
| DROP DATABASE IF EXISTS fkmulti
< ok 1 0
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
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fka FOREIGN KEY (pid) REFERENCES p(id), CONSTRAINT fkb FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fka`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fkb`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP INDEX `fkb`
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
|   `pid` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ForeignKeyMultiTableApplyCorrectness/8.0/drop-two-unnamed-fks-sharing-index default#1
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
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id), FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `c_ibfk_1`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `c_ibfk_2`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP INDEX `pid`
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
|   `pid` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
