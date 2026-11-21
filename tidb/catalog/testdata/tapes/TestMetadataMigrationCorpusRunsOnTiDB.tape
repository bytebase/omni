-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestMetadataMigrationCorpusRunsOnTiDB default#1
> exec
| DROP DATABASE IF EXISTS omni_corpus_a
< ok 0 0
> exec
| CREATE DATABASE omni_corpus_a
< ok 0 0
> exec
| USE omni_corpus_a
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_corpus_b
< ok 0 0
> exec
| CREATE DATABASE omni_corpus_b
< ok 0 0
> exec
| USE omni_corpus_b
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE users (
|     id INT NOT NULL AUTO_INCREMENT,
|     username VARCHAR(50) NOT NULL,
|     email VARCHAR(100) NOT NULL,
|     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_email (email),
|     INDEX idx_username (username)
| )
< ok 0 0
> exec
| CREATE TABLE posts (
|     id INT NOT NULL AUTO_INCREMENT,
|     user_id INT NOT NULL,
|     title VARCHAR(200) NOT NULL,
|     content TEXT,
|     published_at DATETIME,
|     PRIMARY KEY (id),
|     INDEX idx_user_id (user_id),
|     CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
| )
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_corpus_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| posts
: bytes
| BASE TABLE
< row
: bytes
| users
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `posts`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| posts
: bytes
| CREATE TABLE `posts` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `user_id` int NOT NULL,
|   `title` varchar(200) NOT NULL,
|   `content` text DEFAULT NULL,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_user_id` (`user_id`),
|   CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `users`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) NOT NULL,
|   `email` varchar(100) NOT NULL,
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `uk_email` (`email`),
|   KEY `idx_username` (`username`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| USE omni_corpus_a
< ok 0 0
> exec
| CREATE TABLE `omni_corpus_a`.`posts` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `user_id` int NOT NULL,
|   `title` varchar(200) NOT NULL,
|   `content` text DEFAULT NULL,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_user_id` (`user_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
< ok 0 0
> exec
| CREATE TABLE `omni_corpus_a`.`users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) NOT NULL,
|   `email` varchar(100) NOT NULL,
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_username` (`username`),
|   UNIQUE KEY `uk_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`posts` ADD CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_corpus_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| posts
: bytes
| BASE TABLE
< row
: bytes
| users
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `posts`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| posts
: bytes
| CREATE TABLE `posts` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `user_id` int NOT NULL,
|   `title` varchar(200) NOT NULL,
|   `content` text DEFAULT NULL,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_user_id` (`user_id`),
|   CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `users`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) NOT NULL,
|   `email` varchar(100) NOT NULL,
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_username` (`username`),
|   UNIQUE KEY `uk_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_corpus_a
< ok 0 0
> exec
| CREATE DATABASE omni_corpus_a
< ok 0 0
> exec
| USE omni_corpus_a
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE users (
|     id INT NOT NULL AUTO_INCREMENT,
|     username VARCHAR(50) NOT NULL,
|     email VARCHAR(100) NOT NULL,
|     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     is_active BOOLEAN DEFAULT true,
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_email (email),
|     INDEX idx_username (username),
|     INDEX idx_email_active (email, is_active)
| )
< ok 0 0
> exec
| CREATE TABLE posts (
|     id INT NOT NULL AUTO_INCREMENT,
|     user_id INT NOT NULL,
|     title VARCHAR(200) NOT NULL,
|     content TEXT,
|     published_at DATETIME,
|     PRIMARY KEY (id),
|     INDEX idx_user_id (user_id),
|     CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
|     CONSTRAINT chk_title_length CHECK (CHAR_LENGTH(title) > 0)
| )
< ok 0 0
> exec
| CREATE TABLE comments (
|     id INT NOT NULL AUTO_INCREMENT,
|     post_id INT NOT NULL,
|     user_id INT NOT NULL,
|     content TEXT NOT NULL,
|     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     INDEX idx_post_user (post_id, user_id),
|     CONSTRAINT fk_comment_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
|     CONSTRAINT fk_comment_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
| )
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_corpus_b
< ok 0 0
> exec
| CREATE DATABASE omni_corpus_b
< ok 0 0
> exec
| USE omni_corpus_b
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE users (
|     id INT NOT NULL AUTO_INCREMENT,
|     username VARCHAR(50) NOT NULL,
|     email VARCHAR(100) NOT NULL,
|     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_email (email),
|     INDEX idx_username (username)
| )
< ok 0 0
> exec
| CREATE TABLE posts (
|     id INT NOT NULL AUTO_INCREMENT,
|     user_id INT NOT NULL,
|     title VARCHAR(200) NOT NULL,
|     content TEXT,
|     published_at DATETIME,
|     PRIMARY KEY (id),
|     INDEX idx_user_id (user_id),
|     CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
| )
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_corpus_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| posts
: bytes
| BASE TABLE
< row
: bytes
| users
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `posts`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| posts
: bytes
| CREATE TABLE `posts` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `user_id` int NOT NULL,
|   `title` varchar(200) NOT NULL,
|   `content` text DEFAULT NULL,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_user_id` (`user_id`),
|   CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `users`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) NOT NULL,
|   `email` varchar(100) NOT NULL,
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `uk_email` (`email`),
|   KEY `idx_username` (`username`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| USE omni_corpus_a
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`comments` DROP FOREIGN KEY `fk_comment_post`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`comments` DROP FOREIGN KEY `fk_comment_user`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`users` DROP INDEX `idx_email_active`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`users` DROP COLUMN `is_active`
< ok 0 0
> exec
| DROP TABLE `omni_corpus_a`.`comments`
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_corpus_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| posts
: bytes
| BASE TABLE
< row
: bytes
| users
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `posts`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| posts
: bytes
| CREATE TABLE `posts` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `user_id` int NOT NULL,
|   `title` varchar(200) NOT NULL,
|   `content` text DEFAULT NULL,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_user_id` (`user_id`),
|   CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `users`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) NOT NULL,
|   `email` varchar(100) NOT NULL,
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `uk_email` (`email`),
|   KEY `idx_username` (`username`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_corpus_a
< ok 0 0
> exec
| CREATE DATABASE omni_corpus_a
< ok 0 0
> exec
| USE omni_corpus_a
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE sales (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_name VARCHAR(100) NOT NULL,
|     sale_amount DECIMAL(10, 2) NOT NULL,
|     sale_date DATE NOT NULL,
|     PRIMARY KEY (id),
|     INDEX idx_product_date (product_name, sale_date)
| )
< ok 0 0
> exec
| CREATE VIEW monthly_sales AS
| SELECT
|     YEAR(sale_date) AS year,
|     MONTH(sale_date) AS month,
|     SUM(sale_amount) AS total_sales
| FROM sales
| GROUP BY YEAR(sale_date), MONTH(sale_date)
< ok 0 0
> exec
| CREATE VIEW top_products AS
| SELECT
|     product_name,
|     COUNT(*) AS sale_count,
|     SUM(sale_amount) AS total_revenue
| FROM sales
| GROUP BY product_name
| ORDER BY SUM(sale_amount) DESC
< ok 0 0
> exec
| CREATE TABLE sales_reports (
|     id INT NOT NULL AUTO_INCREMENT,
|     report_type VARCHAR(50) NOT NULL,
|     report_data JSON,
|     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     INDEX idx_type (report_type)
| )
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_corpus_b
< ok 0 0
> exec
| CREATE DATABASE omni_corpus_b
< ok 0 0
> exec
| USE omni_corpus_b
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE sales (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_name VARCHAR(100) NOT NULL,
|     sale_amount DECIMAL(10, 2) NOT NULL,
|     sale_date DATE NOT NULL,
|     PRIMARY KEY (id),
|     INDEX idx_product_date (product_name, sale_date)
| )
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_corpus_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| sales
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `sales`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sales
: bytes
| CREATE TABLE `sales` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_name` varchar(100) NOT NULL,
|   `sale_amount` decimal(10,2) NOT NULL,
|   `sale_date` date NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_product_date` (`product_name`,`sale_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| USE omni_corpus_a
< ok 0 0
> exec
| DROP VIEW `omni_corpus_a`.`monthly_sales`
< ok 0 0
> exec
| DROP VIEW `omni_corpus_a`.`top_products`
< ok 0 0
> exec
| DROP TABLE `omni_corpus_a`.`sales_reports`
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_corpus_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| sales
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `sales`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sales
: bytes
| CREATE TABLE `sales` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_name` varchar(100) NOT NULL,
|   `sale_amount` decimal(10,2) NOT NULL,
|   `sale_date` date NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_product_date` (`product_name`,`sale_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_corpus_a
< ok 0 0
> exec
| CREATE DATABASE omni_corpus_a
< ok 0 0
> exec
| USE omni_corpus_a
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE users (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     email VARCHAR(150) NOT NULL,
|     age INT,
|     phone VARCHAR(20),
|     address TEXT,
|     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     INDEX idx_email (email),
|     INDEX idx_name (name),
|     INDEX idx_age (age),
|     CONSTRAINT chk_age_positive CHECK (age > 0),
|     CONSTRAINT uk_phone UNIQUE (phone)
| )
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_corpus_b
< ok 0 0
> exec
| CREATE DATABASE omni_corpus_b
< ok 0 0
> exec
| USE omni_corpus_b
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE users (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(50) NOT NULL,
|     email VARCHAR(100),
|     age INT,
|     PRIMARY KEY (id),
|     INDEX idx_email (email)
| )
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_corpus_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| users
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `users`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(50) NOT NULL,
|   `email` varchar(100) DEFAULT NULL,
|   `age` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| USE omni_corpus_a
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`users` DROP INDEX `idx_age`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`users` DROP INDEX `idx_name`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`users` DROP INDEX `uk_phone`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`users` DROP COLUMN `phone`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`users` DROP COLUMN `address`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`users` DROP COLUMN `created_at`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`users` MODIFY COLUMN `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`users` MODIFY COLUMN `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_corpus_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| users
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `users`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(50) NOT NULL,
|   `email` varchar(100) DEFAULT NULL,
|   `age` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_corpus_a
< ok 0 0
> exec
| CREATE DATABASE omni_corpus_a
< ok 0 0
> exec
| USE omni_corpus_a
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_corpus_b
< ok 0 0
> exec
| CREATE DATABASE omni_corpus_b
< ok 0 0
> exec
| USE omni_corpus_b
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (
|     a INT NOT NULL,
|     b INT GENERATED ALWAYS AS (a + 1) STORED,
|     c INT GENERATED ALWAYS AS (a * 2) VIRTUAL
| )
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_corpus_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL,
|   `b` int GENERATED ALWAYS AS (`a` + 1) STORED,
|   `c` int GENERATED ALWAYS AS (`a` * 2) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| USE omni_corpus_a
< ok 0 0
> exec
| CREATE TABLE `omni_corpus_a`.`t` (
|   `a` int NOT NULL,
|   `b` int GENERATED ALWAYS AS (`a` + 1) STORED,
|   `c` int GENERATED ALWAYS AS (`a` * 2) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_corpus_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL,
|   `b` int GENERATED ALWAYS AS (`a` + 1) STORED,
|   `c` int GENERATED ALWAYS AS (`a` * 2) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_corpus_a
< ok 0 0
> exec
| CREATE DATABASE omni_corpus_a
< ok 0 0
> exec
| USE omni_corpus_a
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE test_table (
|     id INT NOT NULL AUTO_INCREMENT,
|     simple_field VARCHAR(100),
|     PRIMARY KEY (id)
| )
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_corpus_b
< ok 0 0
> exec
| CREATE DATABASE omni_corpus_b
< ok 0 0
> exec
| USE omni_corpus_b
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE test_table (
|     id INT NOT NULL AUTO_INCREMENT,
|     simple_field VARCHAR(100),
|     field_with_quotes VARCHAR(100) COMMENT 'Field with "double" and ''single'' quotes',
|     field_with_symbols VARCHAR(100) COMMENT 'Field with symbols: @#$%^&*()_+-={}[]|:";''<>?,./ and more!',
|     field_with_unicode VARCHAR(100) COMMENT 'Unicode test: 你好世界 🌍 café naïve résumé Москва العالم 한국어',
|     field_multiline TEXT COMMENT 'Multi-line comment:
| Line 1 with regular text
| Line 2 with "quotes" and symbols @#$
| Line 3 with unicode: 测试 🚀 café
| Final line with mixed content',
|     field_with_sql VARCHAR(100) COMMENT 'Comment with SQL-like content: SELECT * FROM table WHERE id = ''123'' AND name LIKE "%test%"',
|     field_with_html VARCHAR(100) COMMENT 'HTML content: <div class="test">Hello & "World"</div> <!-- comment -->',
|     field_with_json VARCHAR(100) COMMENT 'JSON example: {"name": "test", "value": 123, "nested": {"key": "value with spaces"}}',
|     field_with_escape VARCHAR(100) COMMENT 'Escape sequences: \n \t \r " '' \\',
|     PRIMARY KEY (id)
| ) COMMENT = 'Table with "double quotes" and ''single quotes'' in comment';
|
| CREATE TABLE special_comments_table (
|     id BIGINT NOT NULL AUTO_RANDOM COMMENT 'ID with "quotes", symbols @#$, unicode 测试🌟, and
| multi-line
| content',
|     data_field JSON COMMENT 'Updated JSON field: {"new": "structure", "with": ["array", "of", "strings"], "escapes": "\n\t\r"}',
|     html_field TEXT,
|     sql_field VARCHAR(255) COMMENT 'SQL patterns: SELECT * FROM users WHERE name = ''O''Brien'' AND age > 21',
|     unicode_field VARCHAR(200) COMMENT '多语言支持: English, 中文, العربية, Русский, 日本語, 한국어, Français, Español',
|     symbols_field VARCHAR(100) COMMENT 'All symbols: !@#$%^&*()_+-={}[]|:";''<>?,./ plus tab    and newline
| test',
|     url_field VARCHAR(300) COMMENT 'URL with params: https://example.com/path?param1=value1&param2="quoted value"&param3=50%+discount',
|     regex_field VARCHAR(150) COMMENT 'Regex pattern: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$ for email validation',
|     PRIMARY KEY (id),
|     INDEX idx_unicode (unicode_field)
| ) COMMENT = 'Special characters test table with:
| - Double quotes: "test"
| - Single quotes: ''test''
| - Unicode: 特殊字符测试表 🎯
| - Symbols: @#$%^&*()
| - HTML: <div>content</div>
| - JSON: {"key": "value"}
| - URLs: https://example.com
| - Multi-line content with various encodings';
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_corpus_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| special_comments_table
: bytes
| BASE TABLE
< row
: bytes
| test_table
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `special_comments_table`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| special_comments_table
: bytes
| CREATE TABLE `special_comments_table` (
|   `id` bigint NOT NULL /*T![auto_rand] AUTO_RANDOM(5) */ COMMENT 'ID with "quotes", symbols @#$, unicode 测试🌟, and\nmulti-line\ncontent',
|   `data_field` json DEFAULT NULL COMMENT 'Updated JSON field: {"new": "structure", "with": ["array", "of", "strings"], "escapes": "\n	\r"}',
|   `html_field` text DEFAULT NULL,
|   `sql_field` varchar(255) DEFAULT NULL COMMENT 'SQL patterns: SELECT * FROM users WHERE name = ''O''Brien'' AND age > 21',
|   `unicode_field` varchar(200) DEFAULT NULL COMMENT '多语言支持: English, 中文, العربية, Русский, 日本語, 한국어, Français, Español',
|   `symbols_field` varchar(100) DEFAULT NULL COMMENT 'All symbols: !@#$%^&*()_+-={}[]|:";''<>?,./ plus tab    and newline\ntest',
|   `url_field` varchar(300) DEFAULT NULL COMMENT 'URL with params: https://example.com/path?param1=value1&param2="quoted value"&param3=50%+discount',
|   `regex_field` varchar(150) DEFAULT NULL COMMENT 'Regex pattern: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$ for email validation',
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_unicode` (`unicode_field`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Special characters test table with:\n- Double quotes: "test"\n- Single quotes: ''test''\n- Unicode: 特殊字符测试表 🎯\n- Symbols: @#$%^&*()\n- HTML: <div>content</div>\n- JSON: {"key": "value"}\n- URLs: https://example.com\n- Multi-line content with various encodings'
> query
| SHOW CREATE TABLE `test_table`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| test_table
: bytes
| CREATE TABLE `test_table` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `simple_field` varchar(100) DEFAULT NULL,
|   `field_with_quotes` varchar(100) DEFAULT NULL COMMENT 'Field with "double" and ''single'' quotes',
|   `field_with_symbols` varchar(100) DEFAULT NULL COMMENT 'Field with symbols: @#$%^&*()_+-={}[]|:";''<>?,./ and more!',
|   `field_with_unicode` varchar(100) DEFAULT NULL COMMENT 'Unicode test: 你好世界 🌍 café naïve résumé Москва العالم 한국어',
|   `field_multiline` text DEFAULT NULL COMMENT 'Multi-line comment:\nLine 1 with regular text\nLine 2 with "quotes" and symbols @#$\nLine 3 with unicode: 测试 🚀 café\nFinal line with mixed content',
|   `field_with_sql` varchar(100) DEFAULT NULL COMMENT 'Comment with SQL-like content: SELECT * FROM table WHERE id = ''123'' AND name LIKE "%test%"',
|   `field_with_html` varchar(100) DEFAULT NULL COMMENT 'HTML content: <div class="test">Hello & "World"</div> <!-- comment -->',
|   `field_with_json` varchar(100) DEFAULT NULL COMMENT 'JSON example: {"name": "test", "value": 123, "nested": {"key": "value with spaces"}}',
|   `field_with_escape` varchar(100) DEFAULT NULL COMMENT 'Escape sequences: \n 	 \r " '' \\',
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Table with "double quotes" and ''single quotes'' in comment'
> exec
| USE omni_corpus_a
< ok 0 0
> exec
| CREATE TABLE `omni_corpus_a`.`special_comments_table` (
|   `id` bigint NOT NULL /*T![auto_rand] AUTO_RANDOM(5) */ COMMENT 'ID with "quotes", symbols @#$, unicode 测试🌟, and\nmulti-line\ncontent',
|   `data_field` json DEFAULT NULL COMMENT 'Updated JSON field: {"new": "structure", "with": ["array", "of", "strings"], "escapes": "\n	\r"}',
|   `html_field` text DEFAULT NULL,
|   `sql_field` varchar(255) DEFAULT NULL COMMENT 'SQL patterns: SELECT * FROM users WHERE name = ''O''Brien'' AND age > 21',
|   `unicode_field` varchar(200) DEFAULT NULL COMMENT '多语言支持: English, 中文, العربية, Русский, 日本語, 한국어, Français, Español',
|   `symbols_field` varchar(100) DEFAULT NULL COMMENT 'All symbols: !@#$%^&*()_+-={}[]|:";''<>?,./ plus tab    and newline\ntest',
|   `url_field` varchar(300) DEFAULT NULL COMMENT 'URL with params: https://example.com/path?param1=value1&param2="quoted value"&param3=50%+discount',
|   `regex_field` varchar(150) DEFAULT NULL COMMENT 'Regex pattern: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$ for email validation',
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_unicode` (`unicode_field`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Special characters test table with:\n- Double quotes: "test"\n- Single quotes: ''test''\n- Unicode: 特殊字符测试表 🎯\n- Symbols: @#$%^&*()\n- HTML: <div>content</div>\n- JSON: {"key": "value"}\n- URLs: https://example.com\n- Multi-line content with various encodings'
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`test_table` COMMENT='Table with "double quotes" and ''single quotes'' in comment'
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`test_table` ADD COLUMN `field_with_quotes` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Field with "double" and ''single'' quotes' AFTER `simple_field`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`test_table` ADD COLUMN `field_with_symbols` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Field with symbols: @#$%^&*()_+-={}[]|:";''<>?,./ and more!' AFTER `field_with_quotes`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`test_table` ADD COLUMN `field_with_unicode` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Unicode test: 你好世界 🌍 café naïve résumé Москва العالم 한국어' AFTER `field_with_symbols`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`test_table` ADD COLUMN `field_multiline` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Multi-line comment:\nLine 1 with regular text\nLine 2 with "quotes" and symbols @#$\nLine 3 with unicode: 测试 🚀 café\nFinal line with mixed content' AFTER `field_with_unicode`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`test_table` ADD COLUMN `field_with_sql` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Comment with SQL-like content: SELECT * FROM table WHERE id = ''123'' AND name LIKE "%test%"' AFTER `field_multiline`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`test_table` ADD COLUMN `field_with_html` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'HTML content: <div class="test">Hello & "World"</div> <!-- comment -->' AFTER `field_with_sql`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`test_table` ADD COLUMN `field_with_json` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'JSON example: {"name": "test", "value": 123, "nested": {"key": "value with spaces"}}' AFTER `field_with_html`
< ok 0 0
> exec
| ALTER TABLE `omni_corpus_a`.`test_table` ADD COLUMN `field_with_escape` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Escape sequences: \n 	 \r " '' \\' AFTER `field_with_json`
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_corpus_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| special_comments_table
: bytes
| BASE TABLE
< row
: bytes
| test_table
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `special_comments_table`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| special_comments_table
: bytes
| CREATE TABLE `special_comments_table` (
|   `id` bigint NOT NULL /*T![auto_rand] AUTO_RANDOM(5) */ COMMENT 'ID with "quotes", symbols @#$, unicode 测试🌟, and\nmulti-line\ncontent',
|   `data_field` json DEFAULT NULL COMMENT 'Updated JSON field: {"new": "structure", "with": ["array", "of", "strings"], "escapes": "\n	\r"}',
|   `html_field` text DEFAULT NULL,
|   `sql_field` varchar(255) DEFAULT NULL COMMENT 'SQL patterns: SELECT * FROM users WHERE name = ''O''Brien'' AND age > 21',
|   `unicode_field` varchar(200) DEFAULT NULL COMMENT '多语言支持: English, 中文, العربية, Русский, 日本語, 한국어, Français, Español',
|   `symbols_field` varchar(100) DEFAULT NULL COMMENT 'All symbols: !@#$%^&*()_+-={}[]|:";''<>?,./ plus tab    and newline\ntest',
|   `url_field` varchar(300) DEFAULT NULL COMMENT 'URL with params: https://example.com/path?param1=value1&param2="quoted value"&param3=50%+discount',
|   `regex_field` varchar(150) DEFAULT NULL COMMENT 'Regex pattern: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$ for email validation',
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_unicode` (`unicode_field`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Special characters test table with:\n- Double quotes: "test"\n- Single quotes: ''test''\n- Unicode: 特殊字符测试表 🎯\n- Symbols: @#$%^&*()\n- HTML: <div>content</div>\n- JSON: {"key": "value"}\n- URLs: https://example.com\n- Multi-line content with various encodings'
> query
| SHOW CREATE TABLE `test_table`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| test_table
: bytes
| CREATE TABLE `test_table` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `simple_field` varchar(100) DEFAULT NULL,
|   `field_with_quotes` varchar(100) DEFAULT NULL COMMENT 'Field with "double" and ''single'' quotes',
|   `field_with_symbols` varchar(100) DEFAULT NULL COMMENT 'Field with symbols: @#$%^&*()_+-={}[]|:";''<>?,./ and more!',
|   `field_with_unicode` varchar(100) DEFAULT NULL COMMENT 'Unicode test: 你好世界 🌍 café naïve résumé Москва العالم 한국어',
|   `field_multiline` text DEFAULT NULL COMMENT 'Multi-line comment:\nLine 1 with regular text\nLine 2 with "quotes" and symbols @#$\nLine 3 with unicode: 测试 🚀 café\nFinal line with mixed content',
|   `field_with_sql` varchar(100) DEFAULT NULL COMMENT 'Comment with SQL-like content: SELECT * FROM table WHERE id = ''123'' AND name LIKE "%test%"',
|   `field_with_html` varchar(100) DEFAULT NULL COMMENT 'HTML content: <div class="test">Hello & "World"</div> <!-- comment -->',
|   `field_with_json` varchar(100) DEFAULT NULL COMMENT 'JSON example: {"name": "test", "value": 123, "nested": {"key": "value with spaces"}}',
|   `field_with_escape` varchar(100) DEFAULT NULL COMMENT 'Escape sequences: \n 	 \r " '' \\',
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Table with "double quotes" and ''single quotes'' in comment'
> exec
| USE test
< ok 0 0
