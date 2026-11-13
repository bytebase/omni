-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| DROP TABLE IF EXISTS t_drop_con
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
|   `id` int(11) NOT NULL,
|   `parent_id` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_parent` (`parent_id`),
|   CONSTRAINT `fk_parent` FOREIGN KEY (`parent_id`) REFERENCES `t_parent_fk` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_parent_fk
< error
| 1451 23000 Cannot delete or update a parent row: a foreign key constraint fails
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
| DROP TABLE IF EXISTS t_drop_con
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
|   `id` int(11) NOT NULL,
|   `age` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `CONSTRAINT_1` CHECK (`age` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
| DROP TABLE IF EXISTS t_drop_con
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
|   `id` int(11) NOT NULL,
|   `parent_id` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_drop` (`parent_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
| DROP TABLE IF EXISTS t_drop_con
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
|   `id` int(11) NOT NULL,
|   `val` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
