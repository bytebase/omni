-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_5_2_SetVariables default#1
> exec
| DROP TABLE IF EXISTS t_fkoff_child
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t_fkoff_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES nonexistent_parent(id))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE t_fkoff_child
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fkoff_child
: bytes
| CREATE TABLE `t_fkoff_child` (
|   `id` int(11) DEFAULT NULL,
|   `parent_id` int(11) DEFAULT NULL,
|   KEY `parent_id` (`parent_id`),
|   CONSTRAINT `t_fkoff_child_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `nonexistent_parent` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fkdrop_child
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fkdrop_parent
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE t_fkdrop_parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t_fkdrop_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES t_fkdrop_parent(id))
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP TABLE t_fkdrop_parent
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fkon_child
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE t_fkon_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES nonexistent_on_parent(id))
< error
| 1005 HY000 Can't create table `catalog_shared`.`t_fkon_child` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET CHARACTER SET utf8mb4
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION'
< ok 0 0
