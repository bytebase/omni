-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_2_7_TruncateTable default#1
> exec
| DROP TABLE IF EXISTS t_trunc1
< ok 0 0
> exec
| CREATE TABLE t_trunc1 (id INT PRIMARY KEY, name VARCHAR(100))
< ok 0 0
> exec
| TRUNCATE TABLE t_trunc1
< ok 0 0
> query
| SHOW CREATE TABLE t_trunc1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_trunc1
: bytes
| CREATE TABLE `t_trunc1` (
|   `id` int(11) NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_trunc_ai
< ok 0 0
> exec
| CREATE TABLE t_trunc_ai (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100)) AUTO_INCREMENT=1000
< ok 0 0
> query
| SHOW CREATE TABLE t_trunc_ai
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_trunc_ai
: bytes
| CREATE TABLE `t_trunc_ai` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| TRUNCATE TABLE t_trunc_ai
< ok 0 0
> query
| SHOW CREATE TABLE t_trunc_ai
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_trunc_ai
: bytes
| CREATE TABLE `t_trunc_ai` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_trunc_noexist
< ok 0 0
> exec
| TRUNCATE TABLE t_trunc_noexist
< error
| 1146 42S02 Table 'catalog_shared.t_trunc_noexist' doesn't exist
