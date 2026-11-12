-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_1_11_PrimaryKey default#1
> exec
| DROP TABLE IF EXISTS t_pk1
< ok 0 0
> exec
| CREATE TABLE t_pk1 (id INT NOT NULL, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE t_pk1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_pk1
: bytes
| CREATE TABLE `t_pk1` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_pk2
< ok 0 0
> exec
| CREATE TABLE t_pk2 (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a, b))
< ok 0 0
> query
| SHOW CREATE TABLE t_pk2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_pk2
: bytes
| CREATE TABLE `t_pk2` (
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`a`,`b`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_pk3
< ok 0 0
> exec
| CREATE TABLE t_pk3 (id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, name VARCHAR(100), PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE t_pk3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_pk3
: bytes
| CREATE TABLE `t_pk3` (
|   `id` bigint unsigned NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_pk4
< ok 0 0
> exec
| CREATE TABLE t_pk4 (c INT NOT NULL, b INT NOT NULL, a INT NOT NULL, PRIMARY KEY (b, a))
< ok 0 0
> query
| SHOW CREATE TABLE t_pk4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_pk4
: bytes
| CREATE TABLE `t_pk4` (
|   `c` int NOT NULL,
|   `b` int NOT NULL,
|   `a` int NOT NULL,
|   PRIMARY KEY (`b`,`a`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_pk5
< ok 0 0
> exec
| CREATE TABLE t_pk5 (id INT, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE t_pk5
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_pk5
: bytes
| CREATE TABLE `t_pk5` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_pk6
< ok 0 0
> exec
| CREATE TABLE t_pk6 (id INT NOT NULL, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE t_pk6
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_pk6
: bytes
| CREATE TABLE `t_pk6` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
