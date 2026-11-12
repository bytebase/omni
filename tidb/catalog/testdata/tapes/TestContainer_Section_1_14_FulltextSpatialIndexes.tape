-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_1_14_FulltextSpatialIndexes default#1
> exec
| DROP TABLE IF EXISTS t_ft1
< ok 0 0
> exec
| CREATE TABLE t_ft1 (id INT PRIMARY KEY, body TEXT, FULLTEXT KEY `ft_idx` (body))
< ok 0 0
> query
| SHOW CREATE TABLE t_ft1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ft1
: bytes
| CREATE TABLE `t_ft1` (
|   `id` int NOT NULL,
|   `body` text DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_ft2
< ok 0 0
> exec
| CREATE TABLE t_ft2 (id INT PRIMARY KEY, title VARCHAR(200), body TEXT, FULLTEXT KEY `ft_multi` (title, body))
< ok 0 0
> query
| SHOW CREATE TABLE t_ft2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ft2
: bytes
| CREATE TABLE `t_ft2` (
|   `id` int NOT NULL,
|   `title` varchar(200) DEFAULT NULL,
|   `body` text DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_ft3
< ok 0 0
> exec
| CREATE TABLE t_ft3 (id INT PRIMARY KEY, body TEXT, FULLTEXT KEY (body))
< ok 0 0
> query
| SHOW CREATE TABLE t_ft3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ft3
: bytes
| CREATE TABLE `t_ft3` (
|   `id` int NOT NULL,
|   `body` text DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_sp1
< ok 0 0
> exec
| CREATE TABLE t_sp1 (id INT PRIMARY KEY, geo_col GEOMETRY NOT NULL, SPATIAL KEY `sp_idx` (geo_col))
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 56 near "GEOMETRY NOT NULL, SPATIAL KEY `sp_idx` (geo_col))" 
