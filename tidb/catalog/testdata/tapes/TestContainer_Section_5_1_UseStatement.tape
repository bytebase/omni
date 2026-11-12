-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_5_1_UseStatement default#1
> exec
| DROP DATABASE IF EXISTS use_test_db
< ok 0 0
> exec
| CREATE DATABASE use_test_db
< ok 0 0
> exec
| USE use_test_db
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_use_check
< ok 0 0
> exec
| CREATE TABLE t_use_check (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE t_use_check
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_use_check
: bytes
| CREATE TABLE `t_use_check` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS use_test_db
< ok 0 0
> exec
| DROP DATABASE IF EXISTS use_nonexistent_db
< ok 0 0
> exec
| USE use_nonexistent_db
< error
| 1049 42000 Unknown database 'use_nonexistent_db'
