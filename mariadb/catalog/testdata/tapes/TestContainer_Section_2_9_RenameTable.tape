-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_2_9_RenameTable default#1
> exec
| DROP TABLE IF EXISTS t_ren1
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ren2
< ok 0 0
> exec
| CREATE TABLE t_ren1 (id INT PRIMARY KEY, name VARCHAR(100))
< ok 0 0
> exec
| RENAME TABLE t_ren1 TO t_ren2
< ok 0 0
> query
| SHOW CREATE TABLE t_ren2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ren2
: bytes
| CREATE TABLE `t_ren2` (
|   `id` int(11) NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_ren_cross
< ok 0 0
> exec
| CREATE DATABASE IF NOT EXISTS test2
< ok 1 0
> exec
| DROP TABLE IF EXISTS test2.t_ren_cross
< ok 0 0
> exec
| CREATE TABLE t_ren_cross (id INT PRIMARY KEY, val VARCHAR(50) NOT NULL DEFAULT '')
< ok 0 0
> exec
| RENAME TABLE t_ren_cross TO test2.t_ren_cross
< ok 0 0
> query
| SHOW CREATE TABLE test2.t_ren_cross
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ren_cross
: bytes
| CREATE TABLE `t_ren_cross` (
|   `id` int(11) NOT NULL,
|   `val` varchar(50) NOT NULL DEFAULT '',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_mp1
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_mp2
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_mp3
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_mp4
< ok 0 0
> exec
| CREATE TABLE t_mp1 (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t_mp3 (val VARCHAR(100))
< ok 0 0
> exec
| RENAME TABLE t_mp1 TO t_mp2, t_mp3 TO t_mp4
< ok 0 0
> query
| SHOW CREATE TABLE t_mp2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_mp2
: bytes
| CREATE TABLE `t_mp2` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE t_mp4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_mp4
: bytes
| CREATE TABLE `t_mp4` (
|   `val` varchar(100) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_noexist_ren
< ok 0 0
> exec
| RENAME TABLE t_noexist_ren TO t_noexist_ren2
< error
| 1146 42S02 Table 'catalog_shared.t_noexist_ren' doesn't exist
> exec
| DROP TABLE IF EXISTS t_ren_exist1
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ren_exist2
< ok 0 0
> exec
| CREATE TABLE t_ren_exist1 (id INT)
< ok 0 0
> exec
| CREATE TABLE t_ren_exist2 (id INT)
< ok 0 0
> exec
| RENAME TABLE t_ren_exist1 TO t_ren_exist2
< error
| 1050 42S01 Table 't_ren_exist2' already exists
