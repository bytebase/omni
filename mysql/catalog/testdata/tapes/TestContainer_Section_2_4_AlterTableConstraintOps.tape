-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_2_4_AlterTableConstraintOps default#1
> exec
| DROP TABLE IF EXISTS t_parent_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_add_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_con
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_ne
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_enf
< ok 0 0
> exec
| CREATE TABLE t_parent_fk (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t_child_fk (id INT PRIMARY KEY, parent_id INT)
< ok 0 0
> exec
| ALTER TABLE t_child_fk ADD CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES t_parent_fk(id)
< ok 0 0
> query
| SHOW CREATE TABLE t_child_fk
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_child_fk
: bytes
| CREATE TABLE `t_child_fk` (
|   `id` int NOT NULL,
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_parent` (`parent_id`),
|   CONSTRAINT `fk_parent` FOREIGN KEY (`parent_id`) REFERENCES `t_parent_fk` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_parent_fk
< error
| 3730 HY000 Cannot drop table 't_parent_fk' referenced by a foreign key constraint 'fk_parent' on table 't_child_fk'.
> exec
| DROP TABLE IF EXISTS t_child_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_add_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_con
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_ne
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_enf
< ok 0 0
> exec
| CREATE TABLE t_add_chk (id INT PRIMARY KEY, age INT)
< ok 0 0
> exec
| ALTER TABLE t_add_chk ADD CHECK (age > 0)
< ok 0 0
> query
| SHOW CREATE TABLE t_add_chk
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_add_chk
: bytes
| CREATE TABLE `t_add_chk` (
|   `id` int NOT NULL,
|   `age` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_add_chk_chk_1` CHECK ((`age` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_parent_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_add_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_con
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_ne
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_enf
< ok 0 0
> exec
| CREATE TABLE t_parent_dfk (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t_child_dfk (id INT PRIMARY KEY, parent_id INT, CONSTRAINT fk_drop FOREIGN KEY (parent_id) REFERENCES t_parent_dfk(id))
< ok 0 0
> exec
| ALTER TABLE t_child_dfk DROP FOREIGN KEY fk_drop
< ok 0 0
> query
| SHOW CREATE TABLE t_child_dfk
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_child_dfk
: bytes
| CREATE TABLE `t_child_dfk` (
|   `id` int NOT NULL,
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_drop` (`parent_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_parent_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_add_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_con
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_ne
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_enf
< ok 0 0
> exec
| CREATE TABLE t_drop_chk (id INT PRIMARY KEY, age INT, CONSTRAINT chk_age CHECK (age > 0))
< ok 0 0
> exec
| ALTER TABLE t_drop_chk DROP CHECK chk_age
< ok 0 0
> query
| SHOW CREATE TABLE t_drop_chk
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_drop_chk
: bytes
| CREATE TABLE `t_drop_chk` (
|   `id` int NOT NULL,
|   `age` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_parent_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_add_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_con
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_ne
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_enf
< ok 0 0
> exec
| CREATE TABLE t_drop_con (id INT PRIMARY KEY, val INT, CONSTRAINT chk_val CHECK (val >= 0))
< ok 0 0
> exec
| ALTER TABLE t_drop_con DROP CONSTRAINT chk_val
< ok 0 0
> query
| SHOW CREATE TABLE t_drop_con
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_drop_con
: bytes
| CREATE TABLE `t_drop_con` (
|   `id` int NOT NULL,
|   `val` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_parent_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_add_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_con
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_ne
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_enf
< ok 0 0
> exec
| CREATE TABLE t_chk_ne (id INT PRIMARY KEY, score INT, CONSTRAINT chk_score CHECK (score >= 0))
< ok 0 0
> exec
| ALTER TABLE t_chk_ne ALTER CHECK chk_score NOT ENFORCED
< ok 0 0
> query
| SHOW CREATE TABLE t_chk_ne
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chk_ne
: bytes
| CREATE TABLE `t_chk_ne` (
|   `id` int NOT NULL,
|   `score` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `chk_score` CHECK ((`score` >= 0)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_parent_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_child_dfk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_add_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_drop_con
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_ne
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_enf
< ok 0 0
> exec
| CREATE TABLE t_chk_enf (id INT PRIMARY KEY, score INT, CONSTRAINT chk_score2 CHECK (score >= 0) /*!80016 NOT ENFORCED */)
< ok 0 0
> exec
| ALTER TABLE t_chk_enf ALTER CHECK chk_score2 ENFORCED
< ok 0 0
> query
| SHOW CREATE TABLE t_chk_enf
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chk_enf
: bytes
| CREATE TABLE `t_chk_enf` (
|   `id` int NOT NULL,
|   `score` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `chk_score2` CHECK ((`score` >= 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
