-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_AutoIncSupportingKeyIdempotence default#1
> exec
| DROP DATABASE IF EXISTS aidem_autoinc_nonpk_unique
< ok 0 0
> exec
| CREATE DATABASE aidem_autoinc_nonpk_unique
< ok 1 0
> exec
| USE aidem_autoinc_nonpk_unique
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, name VARCHAR(20), UNIQUE KEY uk (id)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE aidem_autoinc_nonpk_unique.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(20) DEFAULT NULL,
|   UNIQUE KEY `uk` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS aidem_autoinc_nonpk_unique_composite
< ok 0 0
> exec
| CREATE DATABASE aidem_autoinc_nonpk_unique_composite
< ok 1 0
> exec
| USE aidem_autoinc_nonpk_unique_composite
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk (id, a)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE aidem_autoinc_nonpk_unique_composite.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `a` int NOT NULL,
|   UNIQUE KEY `uk` (`id`,`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS aidem_autoinc_plain_key
< ok 0 0
> exec
| CREATE DATABASE aidem_autoinc_plain_key
< ok 1 0
> exec
| USE aidem_autoinc_plain_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT, KEY k (id)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE aidem_autoinc_plain_key.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `a` int DEFAULT NULL,
|   KEY `k` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS aidem_autoinc_unique_with_other_pk
< ok 0 0
> exec
| CREATE DATABASE aidem_autoinc_unique_with_other_pk
< ok 1 0
> exec
| USE aidem_autoinc_unique_with_other_pk
< ok 0 0
> exec
| CREATE TABLE t (pk INT NOT NULL, id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (pk), UNIQUE KEY uk (id)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE aidem_autoinc_unique_with_other_pk.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `pk` int NOT NULL,
|   `id` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`pk`),
|   UNIQUE KEY `uk` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS aidem_autoinc_primary_key_guard
< ok 0 0
> exec
| CREATE DATABASE aidem_autoinc_primary_key_guard
< ok 1 0
> exec
| USE aidem_autoinc_primary_key_guard
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR(20)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE aidem_autoinc_primary_key_guard.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
