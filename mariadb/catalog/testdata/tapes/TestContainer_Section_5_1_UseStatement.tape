-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_5_1_UseStatement default#1
> exec
| DROP DATABASE IF EXISTS use_test_db
< ok 0 0
> exec
| CREATE DATABASE use_test_db
< ok 1 0
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
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS use_test_db
< ok 1 0
> exec
| DROP DATABASE IF EXISTS use_nonexistent_db
< ok 0 0
> exec
| USE use_nonexistent_db
< error
| 1049 42000 Unknown database 'use_nonexistent_db'
