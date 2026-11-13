-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_2_3_AlterTableIndexOps default#1
> exec
| DROP TABLE IF EXISTS t_add_idx
< ok 0 0
> exec
| CREATE TABLE t_add_idx (id INT PRIMARY KEY, name VARCHAR(100))
< ok 0 0
> exec
| ALTER TABLE t_add_idx ADD INDEX idx_name (name)
< ok 0 0
> query
| SHOW CREATE TABLE t_add_idx
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_add_idx
: bytes
| CREATE TABLE `t_add_idx` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_add_uniq
< ok 0 0
> exec
| CREATE TABLE t_add_uniq (id INT PRIMARY KEY, email VARCHAR(255))
< ok 0 0
> exec
| ALTER TABLE t_add_uniq ADD UNIQUE INDEX idx_email (email)
< ok 0 0
> query
| SHOW CREATE TABLE t_add_uniq
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_add_uniq
: bytes
| CREATE TABLE `t_add_uniq` (
|   `id` int NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_add_ft
< ok 0 0
> exec
| CREATE TABLE t_add_ft (id INT PRIMARY KEY, body TEXT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE t_add_ft ADD FULLTEXT INDEX idx_body (body)
< ok 0 0
> query
| SHOW CREATE TABLE t_add_ft
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_add_ft
: bytes
| CREATE TABLE `t_add_ft` (
|   `id` int NOT NULL,
|   `body` text,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `idx_body` (`body`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_add_pk
< ok 0 0
> exec
| CREATE TABLE t_add_pk (id INT NOT NULL, name VARCHAR(100))
< ok 0 0
> exec
| ALTER TABLE t_add_pk ADD PRIMARY KEY (id)
< ok 0 0
> query
| SHOW CREATE TABLE t_add_pk
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_add_pk
: bytes
| CREATE TABLE `t_add_pk` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_drop_idx
< ok 0 0
> exec
| CREATE TABLE t_drop_idx (id INT PRIMARY KEY, name VARCHAR(100), KEY idx_name (name))
< ok 0 0
> exec
| ALTER TABLE t_drop_idx DROP INDEX idx_name
< ok 0 0
> query
| SHOW CREATE TABLE t_drop_idx
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_drop_idx
: bytes
| CREATE TABLE `t_drop_idx` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_drop_pk
< ok 0 0
> exec
| CREATE TABLE t_drop_pk (id INT NOT NULL, name VARCHAR(100), PRIMARY KEY (id))
< ok 0 0
> exec
| ALTER TABLE t_drop_pk DROP PRIMARY KEY
< ok 0 0
> query
| SHOW CREATE TABLE t_drop_pk
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_drop_pk
: bytes
| CREATE TABLE `t_drop_pk` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_ren_idx
< ok 0 0
> exec
| CREATE TABLE t_ren_idx (id INT PRIMARY KEY, name VARCHAR(100), KEY idx_old (name))
< ok 0 0
> exec
| ALTER TABLE t_ren_idx RENAME INDEX idx_old TO idx_new
< ok 0 0
> query
| SHOW CREATE TABLE t_ren_idx
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ren_idx
: bytes
| CREATE TABLE `t_ren_idx` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_new` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_idx_invis
< ok 0 0
> exec
| CREATE TABLE t_idx_invis (id INT PRIMARY KEY, name VARCHAR(100), KEY idx_name (name))
< ok 0 0
> exec
| ALTER TABLE t_idx_invis ALTER INDEX idx_name INVISIBLE
< ok 0 0
> query
| SHOW CREATE TABLE t_idx_invis
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx_invis
: bytes
| CREATE TABLE `t_idx_invis` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_name` (`name`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_idx_vis
< ok 0 0
> exec
| CREATE TABLE t_idx_vis (id INT PRIMARY KEY, name VARCHAR(100), KEY idx_name (name) INVISIBLE)
< ok 0 0
> exec
| ALTER TABLE t_idx_vis ALTER INDEX idx_name VISIBLE
< ok 0 0
> query
| SHOW CREATE TABLE t_idx_vis
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx_vis
: bytes
| CREATE TABLE `t_idx_vis` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_dup_pk
< ok 0 0
> exec
| CREATE TABLE t_dup_pk (id INT PRIMARY KEY, val INT NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t_dup_pk ADD PRIMARY KEY (val)
< error
| 1068 42000 Multiple primary key defined
