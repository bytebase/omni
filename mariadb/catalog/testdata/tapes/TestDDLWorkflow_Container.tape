-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   `age` int(11) DEFAULT 0,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   `age` int(11) DEFAULT 0,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `idx_email` (`email`),
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(200) NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   `age` int(11) DEFAULT 0,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `idx_email` (`email`),
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(200) NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   `age` int(11) DEFAULT 0,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `user_id` int(11) NOT NULL,
|   `amount` decimal(10,2) DEFAULT NULL,
|   `created_at` timestamp NULL DEFAULT current_timestamp(),
|   PRIMARY KEY (`id`),
|   KEY `idx_user` (`user_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(200) NOT NULL,
|   `email_address` varchar(255) DEFAULT NULL,
|   `age` int(11) DEFAULT 0,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `idx_email` (`email_address`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(200) NOT NULL,
|   `email_address` varchar(255) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `idx_email` (`email_address`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
