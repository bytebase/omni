-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_VersionDispatchBareCharset default#1
> exec
| DROP DATABASE IF EXISTS vd_bare
< ok 0 0
> exec
| CREATE DATABASE vd_bare
< ok 1 0
> exec
| USE vd_bare
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(50)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE vd_bare.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_VersionDispatchBareCharset/8.0 default#1
> exec
| DROP DATABASE IF EXISTS vdb
< ok 0 0
> exec
| CREATE DATABASE vdb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vdb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(50)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `vdb`.`t` ADD COLUMN `isbn` varchar(20)
< ok 0 0
> query
| SHOW CREATE TABLE vdb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   `isbn` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
