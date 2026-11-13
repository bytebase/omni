-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_InnoDBHashRoundTrip default#1
> exec
| DROP DATABASE IF EXISTS hash_innodb_unique_hash
< ok 0 0
> exec
| CREATE DATABASE hash_innodb_unique_hash
< ok 1 0
> exec
| USE hash_innodb_unique_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), UNIQUE KEY uk (name) USING HASH) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE hash_innodb_unique_hash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_innodb_key_hash
< ok 0 0
> exec
| CREATE DATABASE hash_innodb_key_hash
< ok 1 0
> exec
| USE hash_innodb_key_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), KEY k (name) USING HASH) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE hash_innodb_key_hash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_innodb_default_engine_hash
< ok 0 0
> exec
| CREATE DATABASE hash_innodb_default_engine_hash
< ok 1 0
> exec
| USE hash_innodb_default_engine_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), KEY k (name) USING HASH)
< ok 0 0
> query
| SHOW CREATE TABLE hash_innodb_default_engine_hash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_innodb_key_btree
< ok 0 0
> exec
| CREATE DATABASE hash_innodb_key_btree
< ok 1 0
> exec
| USE hash_innodb_key_btree
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), KEY k (name) USING BTREE) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE hash_innodb_key_btree.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_memory_unique_hash
< ok 0 0
> exec
| CREATE DATABASE hash_memory_unique_hash
< ok 1 0
> exec
| USE hash_memory_unique_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), UNIQUE KEY uk (name) USING HASH) ENGINE=MEMORY
< ok 0 0
> query
| SHOW CREATE TABLE hash_memory_unique_hash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`name`) USING HASH
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_memory_key_hash
< ok 0 0
> exec
| CREATE DATABASE hash_memory_key_hash
< ok 1 0
> exec
| USE hash_memory_key_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), KEY k (name) USING HASH) ENGINE=MEMORY
< ok 0 0
> query
| SHOW CREATE TABLE hash_memory_key_hash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`) USING HASH
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_alter_add_unique_hash_innodb
< ok 0 0
> exec
| CREATE DATABASE hash_alter_add_unique_hash_innodb
< ok 1 0
> exec
| USE hash_alter_add_unique_hash_innodb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=InnoDB;ALTER TABLE t ADD UNIQUE uk (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_alter_add_unique_hash_innodb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_alter_add_unique_hash_memory
< ok 0 0
> exec
| CREATE DATABASE hash_alter_add_unique_hash_memory
< ok 1 0
> exec
| USE hash_alter_add_unique_hash_memory
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=MEMORY;ALTER TABLE t ADD UNIQUE uk (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_alter_add_unique_hash_memory.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`name`) USING HASH
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_alter_add_key_hash_innodb
< ok 0 0
> exec
| CREATE DATABASE hash_alter_add_key_hash_innodb
< ok 1 0
> exec
| USE hash_alter_add_key_hash_innodb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=InnoDB;ALTER TABLE t ADD KEY k (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_alter_add_key_hash_innodb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_alter_add_key_hash_memory
< ok 0 0
> exec
| CREATE DATABASE hash_alter_add_key_hash_memory
< ok 1 0
> exec
| USE hash_alter_add_key_hash_memory
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=MEMORY;ALTER TABLE t ADD KEY k (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_alter_add_key_hash_memory.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`) USING HASH
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_create_index_hash_innodb
< ok 0 0
> exec
| CREATE DATABASE hash_create_index_hash_innodb
< ok 1 0
> exec
| USE hash_create_index_hash_innodb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=InnoDB;CREATE INDEX k ON t (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_create_index_hash_innodb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_create_index_hash_memory
< ok 0 0
> exec
| CREATE DATABASE hash_create_index_hash_memory
< ok 1 0
> exec
| USE hash_create_index_hash_memory
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=MEMORY;CREATE INDEX k ON t (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_create_index_hash_memory.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`) USING HASH
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_create_unique_index_hash_innodb
< ok 0 0
> exec
| CREATE DATABASE hash_create_unique_index_hash_innodb
< ok 1 0
> exec
| USE hash_create_unique_index_hash_innodb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=InnoDB;CREATE UNIQUE INDEX uk ON t (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_create_unique_index_hash_innodb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS hash_create_unique_index_hash_memory
< ok 0 0
> exec
| CREATE DATABASE hash_create_unique_index_hash_memory
< ok 1 0
> exec
| USE hash_create_unique_index_hash_memory
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=MEMORY;CREATE UNIQUE INDEX uk ON t (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_create_unique_index_hash_memory.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`name`) USING HASH
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
