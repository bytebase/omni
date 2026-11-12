-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDDLWorkflow_Container default#1
> exec
| CREATE TABLE users (id INT NOT NULL AUTO_INCREMENT, name VARCHAR(100) NOT NULL, email VARCHAR(255), PRIMARY KEY (id), UNIQUE KEY idx_email (email)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE users
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| ALTER TABLE users ADD COLUMN age INT DEFAULT 0
< ok 0 0
> query
| SHOW CREATE TABLE users
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   `age` int DEFAULT '0',
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| CREATE INDEX idx_name ON users (name)
< ok 0 0
> query
| SHOW CREATE TABLE users
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   `age` int DEFAULT '0',
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `idx_email` (`email`),
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| ALTER TABLE users MODIFY COLUMN name VARCHAR(200) NOT NULL
< ok 0 0
> query
| SHOW CREATE TABLE users
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(200) NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   `age` int DEFAULT '0',
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `idx_email` (`email`),
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP INDEX idx_name ON users
< ok 0 0
> query
| SHOW CREATE TABLE users
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(200) NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   `age` int DEFAULT '0',
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| CREATE TABLE orders (id INT NOT NULL AUTO_INCREMENT, user_id INT NOT NULL, amount DECIMAL(10,2), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id), KEY idx_user (user_id)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE orders
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| orders
: bytes
| CREATE TABLE `orders` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `user_id` int NOT NULL,
|   `amount` decimal(10,2) DEFAULT NULL,
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_user` (`user_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| ALTER TABLE users CHANGE COLUMN email email_address VARCHAR(255)
< ok 0 0
> query
| SHOW CREATE TABLE users
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(200) NOT NULL,
|   `email_address` varchar(255) DEFAULT NULL,
|   `age` int DEFAULT '0',
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `idx_email` (`email_address`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| ALTER TABLE users DROP COLUMN age
< ok 0 0
> query
| SHOW CREATE TABLE users
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(200) NOT NULL,
|   `email_address` varchar(255) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `idx_email` (`email_address`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
