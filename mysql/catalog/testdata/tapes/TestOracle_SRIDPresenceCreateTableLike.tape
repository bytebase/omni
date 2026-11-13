-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_SRIDPresenceCreateTableLike default#1
> exec
| DROP DATABASE IF EXISTS srid_like_clone_srid_zero
< ok 0 0
> exec
| CREATE DATABASE srid_like_clone_srid_zero
< ok 1 0
> exec
| USE srid_like_clone_srid_zero
< ok 0 0
> exec
| CREATE TABLE src (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0);
| CREATE TABLE clone LIKE src
< ok 0 0
> query
| SHOW CREATE TABLE srid_like_clone_srid_zero.clone
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| clone
: bytes
| CREATE TABLE `clone` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_like_clone_srid_4326
< ok 0 0
> exec
| CREATE DATABASE srid_like_clone_srid_4326
< ok 1 0
> exec
| USE srid_like_clone_srid_4326
< ok 0 0
> exec
| CREATE TABLE src (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326);
| CREATE TABLE clone LIKE src
< ok 0 0
> query
| SHOW CREATE TABLE srid_like_clone_srid_4326.clone
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| clone
: bytes
| CREATE TABLE `clone` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS srid_like_clone_srid_none
< ok 0 0
> exec
| CREATE DATABASE srid_like_clone_srid_none
< ok 1 0
> exec
| USE srid_like_clone_srid_none
< ok 0 0
> exec
| CREATE TABLE src (id INT PRIMARY KEY, g GEOMETRY NOT NULL);
| CREATE TABLE clone LIKE src
< ok 0 0
> query
| SHOW CREATE TABLE srid_like_clone_srid_none.clone
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| clone
: bytes
| CREATE TABLE `clone` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
