-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=1000
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
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_trunc_noexist
< ok 0 0
> exec
| TRUNCATE TABLE t_trunc_noexist
< error
| 1146 42S02 Table 'test.t_trunc_noexist' doesn't exist
