-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`b`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
