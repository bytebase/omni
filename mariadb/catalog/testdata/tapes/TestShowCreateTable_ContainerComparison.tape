-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestShowCreateTable_ContainerComparison default#1
> exec
| DROP TABLE IF EXISTS t_types
< ok 0 0
> exec
| CREATE TABLE t_types (a INT, b VARCHAR(100), c TEXT, d DECIMAL(10,2), e DATETIME)
< ok 0 0
> query
| SHOW CREATE TABLE t_types
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_types
: bytes
| CREATE TABLE `t_types` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(100) DEFAULT NULL,
|   `c` text DEFAULT NULL,
|   `d` decimal(10,2) DEFAULT NULL,
|   `e` datetime DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_defaults
< ok 0 0
> exec
| CREATE TABLE t_defaults (id INT NOT NULL, name VARCHAR(50) DEFAULT 'test', active TINYINT(1) DEFAULT 1)
< ok 0 0
> query
| SHOW CREATE TABLE t_defaults
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_defaults
: bytes
| CREATE TABLE `t_defaults` (
|   `id` int(11) NOT NULL,
|   `name` varchar(50) DEFAULT 'test',
|   `active` tinyint(1) DEFAULT 1
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_auto
< ok 0 0
> exec
| CREATE TABLE t_auto (id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE t_auto
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_auto
: bytes
| CREATE TABLE `t_auto` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_multi_pk
< ok 0 0
> exec
| CREATE TABLE t_multi_pk (a INT NOT NULL, b INT NOT NULL, c VARCHAR(10), PRIMARY KEY (a, b))
< ok 0 0
> query
| SHOW CREATE TABLE t_multi_pk
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_multi_pk
: bytes
| CREATE TABLE `t_multi_pk` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   `c` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_unique
< ok 0 0
> exec
| CREATE TABLE t_unique (id INT, email VARCHAR(255), UNIQUE KEY idx_email (email))
< ok 0 0
> query
| SHOW CREATE TABLE t_unique
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_unique
: bytes
| CREATE TABLE `t_unique` (
|   `id` int(11) DEFAULT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   UNIQUE KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
