-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
|   `id` int DEFAULT NULL,
|   `parent_id` int DEFAULT NULL,
|   KEY `fk_1` (`parent_id`),
|   CONSTRAINT `fk_1` FOREIGN KEY (`parent_id`) REFERENCES `nonexistent_parent` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
| 1824 HY000 Failed to open the referenced table 'nonexistent_on_parent'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET CHARACTER SET utf8mb4
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION'
< ok 0 0
