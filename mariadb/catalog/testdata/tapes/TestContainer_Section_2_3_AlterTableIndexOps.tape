-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
|   `id` int(11) NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL,
|   `body` text DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `idx_body` (`body`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL,
|   `name` varchar(100) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_new` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
