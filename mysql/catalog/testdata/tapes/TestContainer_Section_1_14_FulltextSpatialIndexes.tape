-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
|   `body` text,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ft_idx` (`body`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `body` text,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ft_multi` (`title`,`body`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `body` text,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `body` (`body`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_sp1
< ok 0 0
> exec
| CREATE TABLE t_sp1 (id INT PRIMARY KEY, geo_col GEOMETRY NOT NULL, SPATIAL KEY `sp_idx` (geo_col))
< ok 0 0
> query
| SHOW CREATE TABLE t_sp1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_sp1
: bytes
| CREATE TABLE `t_sp1` (
|   `id` int NOT NULL,
|   `geo_col` geometry NOT NULL,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sp_idx` (`geo_col`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
