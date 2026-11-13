-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_2_1_CreateTableVariants default#1
> exec
| DROP TABLE IF EXISTS t_ine
< ok 0 0
> exec
| CREATE TABLE t_ine (id INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t_ine (id INT)
< ok 0 0
> query
| SHOW CREATE TABLE t_ine
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ine
: bytes
| CREATE TABLE `t_ine` (
|   `id` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TEMPORARY TABLE IF EXISTS t_temp
< ok 0 0
> exec
| CREATE TEMPORARY TABLE t_temp (id INT, name VARCHAR(50))
< ok 0 0
> query
| SHOW CREATE TABLE t_temp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_temp
: bytes
| CREATE TEMPORARY TABLE `t_temp` (
|   `id` int(11) DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_like_dst
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_like_src
< ok 0 0
> exec
| CREATE TABLE t_like_src (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL DEFAULT '', score DECIMAL(10,2))
< ok 0 0
> exec
| CREATE TABLE t_like_dst LIKE t_like_src
< ok 0 0
> query
| SHOW CREATE TABLE t_like_dst
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_like_dst
: bytes
| CREATE TABLE `t_like_dst` (
|   `id` int(11) NOT NULL,
|   `name` varchar(100) NOT NULL DEFAULT '',
|   `score` decimal(10,2) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_view_conflict
< ok 0 0
> exec
| DROP VIEW IF EXISTS t_view_conflict
< ok 0 0
> exec
| CREATE VIEW t_view_conflict AS SELECT 1 AS a
< ok 0 0
> exec
| CREATE TABLE t_view_conflict (id INT)
< error
| 1050 42S01 Table 't_view_conflict' already exists
> exec
| DROP TABLE IF EXISTS `select`
< ok 0 0
> exec
| CREATE TABLE `select` (`from` INT, `where` VARCHAR(50))
< ok 0 0
> query
| SHOW CREATE TABLE `select`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| select
: bytes
| CREATE TABLE `select` (
|   `from` int(11) DEFAULT NULL,
|   `where` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
