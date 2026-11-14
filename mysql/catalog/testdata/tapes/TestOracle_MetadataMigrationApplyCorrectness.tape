-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_MetadataMigrationApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS pa_mm_72b1de0e
< ok 0 0
> exec
| CREATE DATABASE pa_mm_72b1de0e
< ok 1 0
> exec
| USE pa_mm_72b1de0e
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE authors (
|   id int NOT NULL AUTO_INCREMENT,
|   name varchar(100) NOT NULL,
|   email varchar(100) DEFAULT NULL,
|   PRIMARY KEY (id),
|   UNIQUE KEY uk_email (email)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE books (
|   id int NOT NULL AUTO_INCREMENT,
|   title varchar(200) NOT NULL,
|   author_id int NOT NULL,
|   isbn varchar(20) DEFAULT NULL,
|   published_year int DEFAULT NULL,
|   price decimal(8,2) DEFAULT NULL,
|   PRIMARY KEY (id),
|   KEY idx_author (author_id),
|   KEY idx_year (published_year),
|   UNIQUE KEY uk_isbn (isbn),
|   CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES authors (id),
|   CONSTRAINT chk_price_positive CHECK (price > 0),
|   CONSTRAINT chk_year_valid CHECK ((published_year >= 1000) and (published_year <= 2100))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE products (
|   id int NOT NULL AUTO_INCREMENT,
|   name varchar(50) NOT NULL,
|   price decimal(8,2) NOT NULL,
|   description text DEFAULT NULL,
|   category varchar(30) DEFAULT NULL,
|   is_active tinyint(1) DEFAULT '1',
|   PRIMARY KEY (id),
|   KEY idx_category (category),
|   KEY idx_price (price)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_72b1de0e`.`books` DROP FOREIGN KEY `fk_author`
< ok 0 0
> exec
| DROP TABLE `pa_mm_72b1de0e`.`authors`
< ok 0 0
> exec
| DROP TABLE `pa_mm_72b1de0e`.`books`
< ok 0 0
> exec
| DROP TABLE `pa_mm_72b1de0e`.`products`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_72b1de0e_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_72b1de0e_to
< ok 1 0
> exec
| USE pa_mm_72b1de0e_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_72b1de0e
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_72b1de0e
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_72b1de0e
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_72b1de0e
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_72b1de0e
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_72b1de0e
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_72b1de0e_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_72b1de0e_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_72b1de0e_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_72b1de0e_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_72b1de0e_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_72b1de0e_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_85ed181e
< ok 0 0
> exec
| CREATE DATABASE pa_mm_85ed181e
< ok 1 0
> exec
| USE pa_mm_85ed181e
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE `pa_mm_85ed181e`.`posts` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `user_id` int NOT NULL,
|   `title` varchar(200) NOT NULL,
|   `content` text,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| CREATE TABLE `pa_mm_85ed181e`.`users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) NOT NULL,
|   `email` varchar(100) NOT NULL,
|   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| ALTER TABLE `pa_mm_85ed181e`.`posts` ADD KEY `idx_user_id` (`user_id`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_85ed181e`.`users` ADD KEY `idx_username` (`username`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_85ed181e`.`users` ADD UNIQUE KEY `uk_email` (`email`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_85ed181e`.`posts` ADD CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_85ed181e_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_85ed181e_to
< ok 1 0
> exec
| USE pa_mm_85ed181e_to
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_85ed181e
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_85ed181e
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
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
| SHOW CREATE TABLE `pa_mm_85ed181e`.`posts`
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
|   `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `content` text COLLATE utf8mb4_unicode_ci,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_user_id` (`user_id`),
|   CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SHOW CREATE TABLE `pa_mm_85ed181e`.`users`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`),
|   KEY `idx_username` (`username`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_85ed181e
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_85ed181e
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_85ed181e
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_85ed181e
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_85ed181e_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_85ed181e_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
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
| SHOW CREATE TABLE `pa_mm_85ed181e_to`.`posts`
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
|   `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `content` text COLLATE utf8mb4_unicode_ci,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_user_id` (`user_id`),
|   CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SHOW CREATE TABLE `pa_mm_85ed181e_to`.`users`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`),
|   KEY `idx_username` (`username`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_85ed181e_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_85ed181e_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_85ed181e_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_85ed181e_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_3fdeecb0
< ok 0 0
> exec
| CREATE DATABASE pa_mm_3fdeecb0
< ok 1 0
> exec
| USE pa_mm_3fdeecb0
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_3fdeecb0`.`comments` DROP FOREIGN KEY `fk_comment_post`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_3fdeecb0`.`comments` DROP FOREIGN KEY `fk_comment_user`
< ok 0 0
> exec
| DROP TABLE `pa_mm_3fdeecb0`.`comments`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_3fdeecb0`.`posts` DROP CHECK `chk_title_length`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_3fdeecb0`.`users` DROP INDEX `idx_email_active`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_3fdeecb0`.`users` DROP COLUMN `is_active`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_3fdeecb0_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_3fdeecb0_to
< ok 1 0
> exec
| USE pa_mm_3fdeecb0_to
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_3fdeecb0
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_3fdeecb0
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
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
| SHOW CREATE TABLE `pa_mm_3fdeecb0`.`posts`
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
|   `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `content` text COLLATE utf8mb4_unicode_ci,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_user_id` (`user_id`),
|   CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SHOW CREATE TABLE `pa_mm_3fdeecb0`.`users`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`),
|   KEY `idx_username` (`username`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_3fdeecb0
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_3fdeecb0
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_3fdeecb0
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_3fdeecb0
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_3fdeecb0_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_3fdeecb0_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
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
| SHOW CREATE TABLE `pa_mm_3fdeecb0_to`.`posts`
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
|   `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `content` text COLLATE utf8mb4_unicode_ci,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_user_id` (`user_id`),
|   CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SHOW CREATE TABLE `pa_mm_3fdeecb0_to`.`users`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`),
|   KEY `idx_username` (`username`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_3fdeecb0_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_3fdeecb0_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_3fdeecb0_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_3fdeecb0_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_e11cc33a
< ok 0 0
> exec
| CREATE DATABASE pa_mm_e11cc33a
< ok 1 0
> exec
| USE pa_mm_e11cc33a
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     price DECIMAL(10, 2) NOT NULL,
|     stock INT DEFAULT 0,
|     PRIMARY KEY (id),
|     INDEX idx_name (name)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE orders (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL,
|     total DECIMAL(10, 2),
|     order_date DATETIME DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE VIEW product_inventory AS
| SELECT
|     p.id,
|     p.name,
|     p.price,
|     p.stock,
|     COALESCE(SUM(o.quantity), 0) AS total_ordered
| FROM products p
| LEFT JOIN orders o ON p.id = o.product_id
| GROUP BY p.id, p.name, p.price, p.stock
< ok 0 0
> exec
| CREATE VIEW low_stock_products AS
| SELECT * FROM products WHERE stock < 10
< ok 0 0
> exec
| CREATE TRIGGER update_order_total
| BEFORE INSERT ON orders
| FOR EACH ROW
| BEGIN
|     DECLARE product_price DECIMAL(10, 2);
|     SELECT price INTO product_price FROM products WHERE id = NEW.product_id;
|     SET NEW.total = NEW.quantity * product_price;
| END
< ok 0 0
> exec
| CREATE PROCEDURE GetProductInventory(IN product_name VARCHAR(100))
| BEGIN
|     SELECT * FROM product_inventory
|     WHERE name LIKE CONCAT('%', product_name, '%');
| END
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP VIEW IF EXISTS `pa_mm_e11cc33a`.`low_stock_products`
< ok 0 0
> exec
| DROP VIEW IF EXISTS `pa_mm_e11cc33a`.`product_inventory`
< ok 0 0
> exec
| DROP PROCEDURE `pa_mm_e11cc33a`.`GetProductInventory`
< ok 0 0
> exec
| DROP TRIGGER `pa_mm_e11cc33a`.`update_order_total`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_e11cc33a_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_e11cc33a_to
< ok 1 0
> exec
| USE pa_mm_e11cc33a_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     price DECIMAL(10, 2) NOT NULL,
|     stock INT DEFAULT 0,
|     PRIMARY KEY (id),
|     INDEX idx_name (name)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE orders (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL,
|     total DECIMAL(10, 2),
|     order_date DATETIME DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_e11cc33a
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_e11cc33a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| orders
: bytes
| BASE TABLE
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_e11cc33a`.`orders`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| orders
: bytes
| CREATE TABLE `orders` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL,
|   `total` decimal(10,2) DEFAULT NULL,
|   `order_date` datetime DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   KEY `fk_product` (`product_id`),
|   CONSTRAINT `fk_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_e11cc33a`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `stock` int DEFAULT '0',
|   PRIMARY KEY (`id`),
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_e11cc33a
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_e11cc33a
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_e11cc33a
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_e11cc33a
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_e11cc33a_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_e11cc33a_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| orders
: bytes
| BASE TABLE
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_e11cc33a_to`.`orders`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| orders
: bytes
| CREATE TABLE `orders` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL,
|   `total` decimal(10,2) DEFAULT NULL,
|   `order_date` datetime DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   KEY `fk_product` (`product_id`),
|   CONSTRAINT `fk_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_e11cc33a_to`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `stock` int DEFAULT '0',
|   PRIMARY KEY (`id`),
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_e11cc33a_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_e11cc33a_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_e11cc33a_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_e11cc33a_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_78faf820
< ok 0 0
> exec
| CREATE DATABASE pa_mm_78faf820
< ok 1 0
> exec
| USE pa_mm_78faf820
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE sales (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL,
|     unit_price DECIMAL(10, 2) NOT NULL,
|     sale_date DATE NOT NULL,
|     PRIMARY KEY (id),
|     INDEX idx_product_date (product_id, sale_date)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE sales_summary (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     month_year VARCHAR(7) NOT NULL,
|     total_sales DECIMAL(10, 2),
|     tax_amount DECIMAL(10, 2),
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_product_month (product_id, month_year)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE FUNCTION CalculateTotalSales(p_product_id INT, p_start_date DATE, p_end_date DATE)
| RETURNS DECIMAL(10, 2)
| DETERMINISTIC
| READS SQL DATA
| BEGIN
|     DECLARE total DECIMAL(10, 2);
|
|     SELECT COALESCE(SUM(quantity * unit_price), 0) INTO total
|     FROM sales
|     WHERE product_id = p_product_id
|     AND sale_date BETWEEN p_start_date AND p_end_date;
|
|     RETURN total;
| END
< ok 0 0
> exec
| CREATE FUNCTION GetTaxAmount(amount DECIMAL(10, 2), tax_rate DECIMAL(5, 2))
| RETURNS DECIMAL(10, 2)
| DETERMINISTIC
| NO SQL
| BEGIN
|     RETURN amount * (tax_rate / 100);
| END
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE `pa_mm_78faf820`.`sales_summary`
< ok 0 0
> exec
| DROP FUNCTION `pa_mm_78faf820`.`CalculateTotalSales`
< ok 0 0
> exec
| DROP FUNCTION `pa_mm_78faf820`.`GetTaxAmount`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_78faf820_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_78faf820_to
< ok 1 0
> exec
| USE pa_mm_78faf820_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE sales (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL,
|     unit_price DECIMAL(10, 2) NOT NULL,
|     sale_date DATE NOT NULL,
|     PRIMARY KEY (id),
|     INDEX idx_product_date (product_id, sale_date)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_78faf820
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_78faf820
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| sales
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_78faf820`.`sales`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sales
: bytes
| CREATE TABLE `sales` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL,
|   `unit_price` decimal(10,2) NOT NULL,
|   `sale_date` date NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_product_date` (`product_id`,`sale_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_78faf820
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_78faf820
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_78faf820
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_78faf820
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_78faf820_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_78faf820_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| sales
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_78faf820_to`.`sales`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sales
: bytes
| CREATE TABLE `sales` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL,
|   `unit_price` decimal(10,2) NOT NULL,
|   `sale_date` date NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_product_date` (`product_id`,`sale_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_78faf820_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_78faf820_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_78faf820_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_78faf820_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_1797b203
< ok 0 0
> exec
| CREATE DATABASE pa_mm_1797b203
< ok 1 0
> exec
| USE pa_mm_1797b203
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     category VARCHAR(50),
|     price DECIMAL(10, 2),
|     sku VARCHAR(50) NOT NULL,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1797b203`.`products` ADD KEY `idx_category` (`category`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1797b203`.`products` ADD KEY `idx_name` (`name`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1797b203`.`products` ADD KEY `idx_price` (`price`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1797b203`.`products` ADD UNIQUE KEY `uk_sku` (`sku`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1797b203`.`products` ADD CONSTRAINT `chk_name_length` CHECK ((char_length(`name`) >= 3))
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1797b203`.`products` ADD CONSTRAINT `chk_price_positive` CHECK ((`price` > 0))
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_1797b203_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_1797b203_to
< ok 1 0
> exec
| USE pa_mm_1797b203_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     category VARCHAR(50),
|     price DECIMAL(10, 2),
|     sku VARCHAR(50) NOT NULL,
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_sku (sku),
|     INDEX idx_name (name),
|     INDEX idx_category (category),
|     INDEX idx_price (price),
|     CONSTRAINT chk_price_positive CHECK (price > 0),
|     CONSTRAINT chk_name_length CHECK (CHAR_LENGTH(name) >= 3)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_1797b203
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_1797b203
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_1797b203`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `category` varchar(50) DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `sku` varchar(50) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sku` (`sku`),
|   KEY `idx_category` (`category`),
|   KEY `idx_name` (`name`),
|   KEY `idx_price` (`price`),
|   CONSTRAINT `chk_name_length` CHECK ((char_length(`name`) >= 3)),
|   CONSTRAINT `chk_price_positive` CHECK ((`price` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_1797b203
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_1797b203
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_1797b203
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_1797b203
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_1797b203_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_1797b203_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_1797b203_to`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `category` varchar(50) DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `sku` varchar(50) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_sku` (`sku`),
|   KEY `idx_name` (`name`),
|   KEY `idx_category` (`category`),
|   KEY `idx_price` (`price`),
|   CONSTRAINT `chk_name_length` CHECK ((char_length(`name`) >= 3)),
|   CONSTRAINT `chk_price_positive` CHECK ((`price` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_1797b203_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_1797b203_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_1797b203_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_1797b203_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_11bb39c3
< ok 0 0
> exec
| CREATE DATABASE pa_mm_11bb39c3
< ok 1 0
> exec
| USE pa_mm_11bb39c3
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
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE PROCEDURE `pa_mm_11bb39c3`.`CalculateDiscount`(IN `amount` DECIMAL(10,2), OUT `discount` DECIMAL(10,2))
| BEGIN
|     SET discount = amount * 0.1;
| END
< ok 0 0
> exec
| CREATE FUNCTION `pa_mm_11bb39c3`.`GetMonthlySalesTotal`(`year` INT, `month` INT) RETURNS decimal(10,2)
|     DETERMINISTIC
|     READS SQL DATA
| BEGIN
|     DECLARE total DECIMAL(10, 2);
|     SELECT COALESCE(SUM(sale_amount), 0) INTO total
|     FROM sales
|     WHERE YEAR(sale_date) = year AND MONTH(sale_date) = month;
|     RETURN total;
| END
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `pa_mm_11bb39c3`.`monthly_sales` AS select year(`sales`.`sale_date`) AS `year`,month(`sales`.`sale_date`) AS `month`,sum(`sales`.`sale_amount`) AS `total_sales` from `sales` group by year(`sales`.`sale_date`),month(`sales`.`sale_date`)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `pa_mm_11bb39c3`.`top_products` AS select `sales`.`product_name` AS `product_name`,count(0) AS `sale_count`,sum(`sales`.`sale_amount`) AS `total_revenue` from `sales` group by `sales`.`product_name` order by sum(`sales`.`sale_amount`) desc
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_11bb39c3_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_11bb39c3_to
< ok 1 0
> exec
| USE pa_mm_11bb39c3_to
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
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
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
| CREATE FUNCTION GetMonthlySalesTotal(year INT, month INT)
| RETURNS DECIMAL(10, 2)
| DETERMINISTIC
| READS SQL DATA
| BEGIN
|     DECLARE total DECIMAL(10, 2);
|     SELECT COALESCE(SUM(sale_amount), 0) INTO total
|     FROM sales
|     WHERE YEAR(sale_date) = year AND MONTH(sale_date) = month;
|     RETURN total;
| END
< ok 0 0
> exec
| CREATE PROCEDURE CalculateDiscount(IN amount DECIMAL(10, 2), OUT discount DECIMAL(10, 2))
| BEGIN
|     SET discount = amount * 0.1;
| END
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_11bb39c3
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_11bb39c3
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| monthly_sales
: bytes
| VIEW
< row
: bytes
| sales
: bytes
| BASE TABLE
< row
: bytes
| top_products
: bytes
| VIEW
> query
| SHOW CREATE VIEW `pa_mm_11bb39c3`.`monthly_sales`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| monthly_sales
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `monthly_sales` AS select year(`sales`.`sale_date`) AS `year`,month(`sales`.`sale_date`) AS `month`,sum(`sales`.`sale_amount`) AS `total_sales` from `sales` group by year(`sales`.`sale_date`),month(`sales`.`sale_date`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `pa_mm_11bb39c3`.`sales`
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW `pa_mm_11bb39c3`.`top_products`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| top_products
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `top_products` AS select `sales`.`product_name` AS `product_name`,count(0) AS `sale_count`,sum(`sales`.`sale_amount`) AS `total_revenue` from `sales` group by `sales`.`product_name` order by sum(`sales`.`sale_amount`) desc
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_11bb39c3
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
< row
: bytes
| FUNCTION
: bytes
| GetMonthlySalesTotal
: bytes
| decimal(10,2)
: bytes
| BEGIN
|     DECLARE total DECIMAL(10, 2);
|     SELECT COALESCE(SUM(sale_amount), 0) INTO total
|     FROM sales
|     WHERE YEAR(sale_date) = year AND MONTH(sale_date) = month;
|     RETURN total;
| END
: bytes
| YES
: bytes
| READS SQL DATA
: bytes
| DEFINER
: bytes
|
< row
: bytes
| PROCEDURE
: bytes
| CalculateDiscount
: null
: bytes
| BEGIN
|     SET discount = amount * 0.1;
| END
: bytes
| NO
: bytes
| CONTAINS SQL
: bytes
| DEFINER
: bytes
|
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_11bb39c3
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< row
: bytes
| CalculateDiscount
: uint64 1
: bytes
| IN
: bytes
| amount
: bytes
| decimal(10,2)
< row
: bytes
| CalculateDiscount
: uint64 2
: bytes
| OUT
: bytes
| discount
: bytes
| decimal(10,2)
< row
: bytes
| GetMonthlySalesTotal
: uint64 0
: null
: null
: bytes
| decimal(10,2)
< row
: bytes
| GetMonthlySalesTotal
: uint64 1
: bytes
| IN
: bytes
| year
: bytes
| int
< row
: bytes
| GetMonthlySalesTotal
: uint64 2
: bytes
| IN
: bytes
| month
: bytes
| int
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_11bb39c3
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_11bb39c3
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_11bb39c3_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_11bb39c3_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| monthly_sales
: bytes
| VIEW
< row
: bytes
| sales
: bytes
| BASE TABLE
< row
: bytes
| top_products
: bytes
| VIEW
> query
| SHOW CREATE VIEW `pa_mm_11bb39c3_to`.`monthly_sales`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| monthly_sales
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `monthly_sales` AS select year(`sales`.`sale_date`) AS `year`,month(`sales`.`sale_date`) AS `month`,sum(`sales`.`sale_amount`) AS `total_sales` from `sales` group by year(`sales`.`sale_date`),month(`sales`.`sale_date`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `pa_mm_11bb39c3_to`.`sales`
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW `pa_mm_11bb39c3_to`.`top_products`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| top_products
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `top_products` AS select `sales`.`product_name` AS `product_name`,count(0) AS `sale_count`,sum(`sales`.`sale_amount`) AS `total_revenue` from `sales` group by `sales`.`product_name` order by sum(`sales`.`sale_amount`) desc
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_11bb39c3_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
< row
: bytes
| FUNCTION
: bytes
| GetMonthlySalesTotal
: bytes
| decimal(10,2)
: bytes
| BEGIN
|     DECLARE total DECIMAL(10, 2);
|     SELECT COALESCE(SUM(sale_amount), 0) INTO total
|     FROM sales
|     WHERE YEAR(sale_date) = year AND MONTH(sale_date) = month;
|     RETURN total;
| END
: bytes
| YES
: bytes
| READS SQL DATA
: bytes
| DEFINER
: bytes
|
< row
: bytes
| PROCEDURE
: bytes
| CalculateDiscount
: null
: bytes
| BEGIN
|     SET discount = amount * 0.1;
| END
: bytes
| NO
: bytes
| CONTAINS SQL
: bytes
| DEFINER
: bytes
|
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_11bb39c3_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< row
: bytes
| CalculateDiscount
: uint64 1
: bytes
| IN
: bytes
| amount
: bytes
| decimal(10,2)
< row
: bytes
| CalculateDiscount
: uint64 2
: bytes
| OUT
: bytes
| discount
: bytes
| decimal(10,2)
< row
: bytes
| GetMonthlySalesTotal
: uint64 0
: null
: null
: bytes
| decimal(10,2)
< row
: bytes
| GetMonthlySalesTotal
: uint64 1
: bytes
| IN
: bytes
| year
: bytes
| int
< row
: bytes
| GetMonthlySalesTotal
: uint64 2
: bytes
| IN
: bytes
| month
: bytes
| int
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_11bb39c3_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_11bb39c3_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_76d6ceec
< ok 0 0
> exec
| CREATE DATABASE pa_mm_76d6ceec
< ok 1 0
> exec
| USE pa_mm_76d6ceec
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     price DECIMAL(10, 2) NOT NULL,
|     description TEXT NOT NULL,
|     category VARCHAR(50),
|     is_active TINYINT(1) DEFAULT 1,
|     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|     stock_quantity INT DEFAULT 0,
|     weight DECIMAL(5, 2),
|     PRIMARY KEY (id),
|     INDEX idx_category (category),
|     INDEX idx_price (price),
|     INDEX idx_created_at (created_at),
|     UNIQUE INDEX uk_name_category (name, category),
|     CONSTRAINT chk_price_positive CHECK (price > 0),
|     CONSTRAINT chk_stock_non_negative CHECK (stock_quantity >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_76d6ceec`.`products` DROP CHECK `chk_price_positive`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_76d6ceec`.`products` DROP CHECK `chk_stock_non_negative`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_76d6ceec`.`products` DROP INDEX `idx_created_at`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_76d6ceec`.`products` DROP INDEX `uk_name_category`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_76d6ceec`.`products` DROP COLUMN `created_at`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_76d6ceec`.`products` DROP COLUMN `stock_quantity`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_76d6ceec`.`products` DROP COLUMN `updated_at`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_76d6ceec`.`products` DROP COLUMN `weight`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_76d6ceec`.`products` MODIFY COLUMN `category` varchar(30) DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_76d6ceec`.`products` MODIFY COLUMN `description` text
< ok 0 0
> exec
| ALTER TABLE `pa_mm_76d6ceec`.`products` MODIFY COLUMN `name` varchar(50) NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_76d6ceec`.`products` MODIFY COLUMN `price` decimal(8,2) NOT NULL
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_76d6ceec_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_76d6ceec_to
< ok 1 0
> exec
| USE pa_mm_76d6ceec_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(50) NOT NULL,
|     price DECIMAL(8, 2) NOT NULL,
|     description TEXT,
|     category VARCHAR(30),
|     is_active TINYINT(1) DEFAULT 1,
|     PRIMARY KEY (id),
|     INDEX idx_category (category),
|     INDEX idx_price (price)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_76d6ceec
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_76d6ceec
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_76d6ceec`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(50) NOT NULL,
|   `price` decimal(8,2) NOT NULL,
|   `description` text,
|   `category` varchar(30) DEFAULT NULL,
|   `is_active` tinyint(1) DEFAULT '1',
|   PRIMARY KEY (`id`),
|   KEY `idx_category` (`category`),
|   KEY `idx_price` (`price`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_76d6ceec
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_76d6ceec
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_76d6ceec
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_76d6ceec
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_76d6ceec_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_76d6ceec_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_76d6ceec_to`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(50) NOT NULL,
|   `price` decimal(8,2) NOT NULL,
|   `description` text,
|   `category` varchar(30) DEFAULT NULL,
|   `is_active` tinyint(1) DEFAULT '1',
|   PRIMARY KEY (`id`),
|   KEY `idx_category` (`category`),
|   KEY `idx_price` (`price`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_76d6ceec_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_76d6ceec_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_76d6ceec_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_76d6ceec_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_c711790b
< ok 0 0
> exec
| CREATE DATABASE pa_mm_c711790b
< ok 1 0
> exec
| USE pa_mm_c711790b
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| --
| -- Table structure for `authors`
| --
| CREATE TABLE `authors` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `email` varchar(100) DEFAULT (NULL),
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| --
| -- Table structure for `books`
| --
| CREATE TABLE `books` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `title` varchar(200) NOT NULL,
|   `author_id` int NOT NULL,
|   `isbn` varchar(20) DEFAULT (NULL),
|   `published_year` int DEFAULT (NULL),
|   `price` decimal(8,2) DEFAULT (NULL),
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_isbn` (`isbn`),
|   KEY `idx_author` (`author_id`),
|   KEY `idx_year` (`published_year`),
|   CONSTRAINT `fk_author_new` FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
|   CONSTRAINT `chk_price_positive` CHECK (`price` > 0),
|   CONSTRAINT `chk_year_extended` CHECK ((`published_year` >= 1000) and (`published_year` <= 2030)),
|   CONSTRAINT `chk_title_length` CHECK (char_length(`title`) >= 3)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_c711790b`.`books` DROP FOREIGN KEY `fk_author_new`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_c711790b`.`books` DROP CHECK `chk_title_length`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_c711790b`.`books` DROP CHECK `chk_year_extended`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_c711790b`.`books` ADD CONSTRAINT `chk_year_valid` CHECK (((`published_year` >= 1000) and (`published_year` <= 2100)))
< ok 0 0
> exec
| ALTER TABLE `pa_mm_c711790b`.`books` ADD CONSTRAINT `fk_author` FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_c711790b_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_c711790b_to
< ok 1 0
> exec
| USE pa_mm_c711790b_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE authors (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     email VARCHAR(100),
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_email (email)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE books (
|     id INT NOT NULL AUTO_INCREMENT,
|     title VARCHAR(200) NOT NULL,
|     author_id INT NOT NULL,
|     isbn VARCHAR(20),
|     published_year INT,
|     price DECIMAL(8, 2),
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_isbn (isbn),
|     INDEX idx_author (author_id),
|     INDEX idx_year (published_year),
|     CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES authors(id),
|     CONSTRAINT chk_year_valid CHECK (published_year >= 1000 AND published_year <= 2100),
|     CONSTRAINT chk_price_positive CHECK (price > 0)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_c711790b
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_c711790b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| authors
: bytes
| BASE TABLE
< row
: bytes
| books
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_c711790b`.`authors`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| authors
: bytes
| CREATE TABLE `authors` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `email` varchar(100) DEFAULT (NULL),
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_c711790b`.`books`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| books
: bytes
| CREATE TABLE `books` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `title` varchar(200) NOT NULL,
|   `author_id` int NOT NULL,
|   `isbn` varchar(20) DEFAULT (NULL),
|   `published_year` int DEFAULT (NULL),
|   `price` decimal(8,2) DEFAULT (NULL),
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_isbn` (`isbn`),
|   KEY `idx_author` (`author_id`),
|   KEY `idx_year` (`published_year`),
|   CONSTRAINT `fk_author` FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`),
|   CONSTRAINT `chk_price_positive` CHECK ((`price` > 0)),
|   CONSTRAINT `chk_year_valid` CHECK (((`published_year` >= 1000) and (`published_year` <= 2100)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_c711790b
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_c711790b
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_c711790b
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_c711790b
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_c711790b_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_c711790b_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| authors
: bytes
| BASE TABLE
< row
: bytes
| books
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_c711790b_to`.`authors`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| authors
: bytes
| CREATE TABLE `authors` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `email` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_c711790b_to`.`books`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| books
: bytes
| CREATE TABLE `books` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `title` varchar(200) NOT NULL,
|   `author_id` int NOT NULL,
|   `isbn` varchar(20) DEFAULT NULL,
|   `published_year` int DEFAULT NULL,
|   `price` decimal(8,2) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_isbn` (`isbn`),
|   KEY `idx_author` (`author_id`),
|   KEY `idx_year` (`published_year`),
|   CONSTRAINT `fk_author` FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`),
|   CONSTRAINT `chk_price_positive` CHECK ((`price` > 0)),
|   CONSTRAINT `chk_year_valid` CHECK (((`published_year` >= 1000) and (`published_year` <= 2100)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_c711790b_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_c711790b_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_c711790b_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_c711790b_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_6858194b
< ok 0 0
> exec
| CREATE DATABASE pa_mm_6858194b
< ok 1 0
> exec
| USE pa_mm_6858194b
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE articles (
|     id INT NOT NULL AUTO_INCREMENT,
|     title VARCHAR(200) NOT NULL,
|     content TEXT NOT NULL,
|     tags VARCHAR(500),
|     PRIMARY KEY (id),
|     FULLTEXT idx_fulltext_content (title, content),
|     FULLTEXT idx_fulltext_tags (tags),
|     INDEX idx_title_tags (title(50), tags(50))
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE locations (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     position POINT NOT NULL,
|     PRIMARY KEY (id),
|     SPATIAL INDEX idx_spatial_position (position)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_6858194b`.`articles` DROP INDEX `idx_fulltext_content`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_6858194b`.`articles` DROP INDEX `idx_fulltext_tags`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_6858194b`.`articles` DROP INDEX `idx_title_tags`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_6858194b`.`locations` DROP INDEX `idx_spatial_position`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_6858194b_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_6858194b_to
< ok 1 0
> exec
| USE pa_mm_6858194b_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE articles (
|     id INT NOT NULL AUTO_INCREMENT,
|     title VARCHAR(200) NOT NULL,
|     content TEXT NOT NULL,
|     tags VARCHAR(500),
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE locations (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     position POINT NOT NULL,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_6858194b
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_6858194b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| articles
: bytes
| BASE TABLE
< row
: bytes
| locations
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_6858194b`.`articles`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| articles
: bytes
| CREATE TABLE `articles` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `title` varchar(200) NOT NULL,
|   `content` text NOT NULL,
|   `tags` varchar(500) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_6858194b`.`locations`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| locations
: bytes
| CREATE TABLE `locations` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `position` point NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_6858194b
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_6858194b
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_6858194b
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_6858194b
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_6858194b_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_6858194b_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| articles
: bytes
| BASE TABLE
< row
: bytes
| locations
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_6858194b_to`.`articles`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| articles
: bytes
| CREATE TABLE `articles` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `title` varchar(200) NOT NULL,
|   `content` text NOT NULL,
|   `tags` varchar(500) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_6858194b_to`.`locations`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| locations
: bytes
| CREATE TABLE `locations` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `position` point NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_6858194b_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_6858194b_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_6858194b_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_6858194b_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_668e0e2b
< ok 0 0
> exec
| CREATE DATABASE pa_mm_668e0e2b
< ok 1 0
> exec
| USE pa_mm_668e0e2b
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| --
| -- Table structure for `customers`
| --
| CREATE TABLE `customers` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `preferred_order_id` int DEFAULT (NULL),
|   PRIMARY KEY (`id`),
|   KEY `fk_preferred_order` (`preferred_order_id`),
|   CONSTRAINT `fk_preferred_order` FOREIGN KEY (`preferred_order_id`) REFERENCES `orders` (`id`) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| --
| -- Table structure for `order_items`
| --
| CREATE TABLE `order_items` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `order_id` int NOT NULL,
|   `product_name` varchar(100) NOT NULL,
|   `quantity` int NOT NULL,
|   `unit_price` decimal(10,2) NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_order` (`order_id`),
|   CONSTRAINT `fk_order_item` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| --
| -- Table structure for `orders`
| --
| CREATE TABLE `orders` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `customer_id` int NOT NULL,
|   `order_date` datetime DEFAULT CURRENT_TIMESTAMP,
|   `total_amount` decimal(10,2) DEFAULT (NULL),
|   PRIMARY KEY (`id`),
|   KEY `fk_customer` (`customer_id`),
|   CONSTRAINT `fk_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE TRIGGER `trg_update_order_total`
| AFTER INSERT ON `order_items` FOR EACH ROW
| BEGIN
|     UPDATE orders
|     SET total_amount = (
|         SELECT SUM(quantity * unit_price)
|         FROM order_items
|         WHERE order_id = NEW.order_id
|     )
|     WHERE id = NEW.order_id;
| END
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_668e0e2b`.`customers` DROP FOREIGN KEY `fk_preferred_order`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_668e0e2b`.`order_items` DROP FOREIGN KEY `fk_order_item`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_668e0e2b`.`orders` DROP FOREIGN KEY `fk_customer`
< ok 0 0
> exec
| DROP TABLE `pa_mm_668e0e2b`.`order_items`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_668e0e2b`.`customers` DROP INDEX `fk_preferred_order`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_668e0e2b`.`orders` DROP INDEX `fk_customer`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_668e0e2b_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_668e0e2b_to
< ok 1 0
> exec
| USE pa_mm_668e0e2b_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE customers (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     preferred_order_id INT,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE orders (
|     id INT NOT NULL AUTO_INCREMENT,
|     customer_id INT NOT NULL,
|     order_date DATETIME DEFAULT CURRENT_TIMESTAMP,
|     total_amount DECIMAL(10, 2),
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_668e0e2b
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_668e0e2b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| customers
: bytes
| BASE TABLE
< row
: bytes
| orders
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_668e0e2b`.`customers`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| customers
: bytes
| CREATE TABLE `customers` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `preferred_order_id` int DEFAULT (NULL),
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_668e0e2b`.`orders`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| orders
: bytes
| CREATE TABLE `orders` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `customer_id` int NOT NULL,
|   `order_date` datetime DEFAULT CURRENT_TIMESTAMP,
|   `total_amount` decimal(10,2) DEFAULT (NULL),
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_668e0e2b
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_668e0e2b
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_668e0e2b
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_668e0e2b
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_668e0e2b_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_668e0e2b_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| customers
: bytes
| BASE TABLE
< row
: bytes
| orders
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_668e0e2b_to`.`customers`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| customers
: bytes
| CREATE TABLE `customers` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `preferred_order_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_668e0e2b_to`.`orders`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| orders
: bytes
| CREATE TABLE `orders` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `customer_id` int NOT NULL,
|   `order_date` datetime DEFAULT CURRENT_TIMESTAMP,
|   `total_amount` decimal(10,2) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_668e0e2b_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_668e0e2b_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_668e0e2b_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_668e0e2b_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_1a8e1c71
< ok 0 0
> exec
| CREATE DATABASE pa_mm_1a8e1c71
< ok 1 0
> exec
| USE pa_mm_1a8e1c71
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| --
| -- Temporary view structure for `regional_sales`
| --
| -- `regional_sales` references invalid table(s) or column(s) or function(s) or definer/invoker of view lack rights to use them
|
| --
| -- Table structure for `sales_data`
| --
| CREATE TABLE `sales_data` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `sale_date` date NOT NULL,
|   `customer_id` int NOT NULL,
|   `product_id` int NOT NULL,
|   `amount` decimal(10,2) NOT NULL,
|   `region` varchar(50) NOT NULL,
|   PRIMARY KEY (`id`, `sale_date`),
|   KEY `idx_customer` (`customer_id`),
|   KEY `idx_product` (`product_id`),
|   KEY `idx_region_date` (`region`, `sale_date`),
|   CONSTRAINT `chk_amount_positive` CHECK (`amount` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| --
| -- View structure for `regional_sales`
| --
| CREATE OR REPLACE VIEW `regional_sales` AS select `sales_data`.`region` AS `region`,year(`sales_data`.`sale_date`) AS `year`,month(`sales_data`.`sale_date`) AS `month`,sum(`sales_data`.`amount`) AS `total_sales`,count(distinct `sales_data`.`customer_id`) AS `unique_customers` from `sales_data` group by `sales_data`.`region`,year(`sales_data`.`sale_date`),month(`sales_data`.`sale_date`)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1a8e1c71`.`sales_data` DROP CHECK `chk_amount_positive`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1a8e1c71`.`sales_data` DROP INDEX `idx_customer`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1a8e1c71`.`sales_data` DROP INDEX `idx_product`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1a8e1c71`.`sales_data` DROP INDEX `idx_region_date`
< ok 0 0
> exec
| DROP VIEW IF EXISTS `pa_mm_1a8e1c71`.`regional_sales`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1a8e1c71`.`sales_data` /*!50100 PARTITION BY RANGE (year(`sale_date`))
| (PARTITION p2020 VALUES LESS THAN (2021) ENGINE = innodb,
|  PARTITION p2021 VALUES LESS THAN (2022) ENGINE = innodb,
|  PARTITION p2022 VALUES LESS THAN (2023) ENGINE = innodb,
|  PARTITION p2023 VALUES LESS THAN (2024) ENGINE = innodb,
|  PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = innodb) */
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_1a8e1c71_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_1a8e1c71_to
< ok 1 0
> exec
| USE pa_mm_1a8e1c71_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE sales_data (
|     id INT NOT NULL AUTO_INCREMENT,
|     sale_date DATE NOT NULL,
|     customer_id INT NOT NULL,
|     product_id INT NOT NULL,
|     amount DECIMAL(10, 2) NOT NULL,
|     region VARCHAR(50) NOT NULL,
|     PRIMARY KEY (id, sale_date)
| ) ENGINE=InnoDB
| PARTITION BY RANGE (YEAR(sale_date)) (
|     PARTITION p2020 VALUES LESS THAN (2021),
|     PARTITION p2021 VALUES LESS THAN (2022),
|     PARTITION p2022 VALUES LESS THAN (2023),
|     PARTITION p2023 VALUES LESS THAN (2024),
|     PARTITION p_future VALUES LESS THAN MAXVALUE
| )
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_1a8e1c71
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_1a8e1c71
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| sales_data
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_1a8e1c71`.`sales_data`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sales_data
: bytes
| CREATE TABLE `sales_data` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `sale_date` date NOT NULL,
|   `customer_id` int NOT NULL,
|   `product_id` int NOT NULL,
|   `amount` decimal(10,2) NOT NULL,
|   `region` varchar(50) NOT NULL,
|   PRIMARY KEY (`id`,`sale_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`sale_date`))
| (PARTITION p2020 VALUES LESS THAN (2021) ENGINE = InnoDB,
|  PARTITION p2021 VALUES LESS THAN (2022) ENGINE = InnoDB,
|  PARTITION p2022 VALUES LESS THAN (2023) ENGINE = InnoDB,
|  PARTITION p2023 VALUES LESS THAN (2024) ENGINE = InnoDB,
|  PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_1a8e1c71
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_1a8e1c71
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_1a8e1c71
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_1a8e1c71
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_1a8e1c71_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_1a8e1c71_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| sales_data
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_1a8e1c71_to`.`sales_data`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sales_data
: bytes
| CREATE TABLE `sales_data` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `sale_date` date NOT NULL,
|   `customer_id` int NOT NULL,
|   `product_id` int NOT NULL,
|   `amount` decimal(10,2) NOT NULL,
|   `region` varchar(50) NOT NULL,
|   PRIMARY KEY (`id`,`sale_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`sale_date`))
| (PARTITION p2020 VALUES LESS THAN (2021) ENGINE = InnoDB,
|  PARTITION p2021 VALUES LESS THAN (2022) ENGINE = InnoDB,
|  PARTITION p2022 VALUES LESS THAN (2023) ENGINE = InnoDB,
|  PARTITION p2023 VALUES LESS THAN (2024) ENGINE = InnoDB,
|  PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_1a8e1c71_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_1a8e1c71_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_1a8e1c71_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_1a8e1c71_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_a162ef44
< ok 0 0
> exec
| CREATE DATABASE pa_mm_a162ef44
< ok 1 0
> exec
| USE pa_mm_a162ef44
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| --
| -- Table structure for `ads_member_asset_by_account_channel_d`
| --
| CREATE TABLE `ads_member_asset_by_account_channel_d` (
|   `dt` date NOT NULL,
|   `member_id` bigint NOT NULL,
|   `account_type` varchar(50) NOT NULL,
|   `net_asset_amount` decimal(38,10) DEFAULT (NULL)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_a162ef44`.`ads_member_asset_by_account_channel_d` ADD PRIMARY KEY (`dt`,`member_id`,`account_type`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_a162ef44`.`ads_member_asset_by_account_channel_d` ADD KEY `sindex` (`dt`,`member_id`,`account_type`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_a162ef44`.`ads_member_asset_by_account_channel_d` /*!50100 PARTITION BY RANGE (to_days(`dt`))
| (PARTITION p0 VALUES LESS THAN (739000) ENGINE = innodb,
|  PARTITION p1 VALUES LESS THAN (739001) ENGINE = innodb,
|  PARTITION p2 VALUES LESS THAN (739002) ENGINE = innodb,
|  PARTITION p3 VALUES LESS THAN (739003) ENGINE = innodb,
|  PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = innodb) */
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_a162ef44_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_a162ef44_to
< ok 1 0
> exec
| USE pa_mm_a162ef44_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE ads_member_asset_by_account_channel_d (
|     dt DATE NOT NULL,
|     member_id BIGINT NOT NULL,
|     account_type VARCHAR(50) NOT NULL,
|     net_asset_amount DECIMAL(38,10) DEFAULT NULL,
|     PRIMARY KEY (dt, member_id, account_type),
|     KEY sindex (dt, member_id, account_type)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
| PARTITION BY RANGE (TO_DAYS(dt)) (
|     PARTITION p0 VALUES LESS THAN (739000),
|     PARTITION p1 VALUES LESS THAN (739001),
|     PARTITION p2 VALUES LESS THAN (739002),
|     PARTITION p3 VALUES LESS THAN (739003),
|     PARTITION p_future VALUES LESS THAN MAXVALUE
| )
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_a162ef44
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_a162ef44
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| ads_member_asset_by_account_channel_d
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_a162ef44`.`ads_member_asset_by_account_channel_d`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ads_member_asset_by_account_channel_d
: bytes
| CREATE TABLE `ads_member_asset_by_account_channel_d` (
|   `dt` date NOT NULL,
|   `member_id` bigint NOT NULL,
|   `account_type` varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
|   `net_asset_amount` decimal(38,10) DEFAULT (NULL),
|   PRIMARY KEY (`dt`,`member_id`,`account_type`),
|   KEY `sindex` (`dt`,`member_id`,`account_type`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
| /*!50100 PARTITION BY RANGE (to_days(`dt`))
| (PARTITION p0 VALUES LESS THAN (739000) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (739001) ENGINE = InnoDB,
|  PARTITION p2 VALUES LESS THAN (739002) ENGINE = InnoDB,
|  PARTITION p3 VALUES LESS THAN (739003) ENGINE = InnoDB,
|  PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_a162ef44
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_a162ef44
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_a162ef44
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_a162ef44
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_a162ef44_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_a162ef44_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| ads_member_asset_by_account_channel_d
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_a162ef44_to`.`ads_member_asset_by_account_channel_d`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ads_member_asset_by_account_channel_d
: bytes
| CREATE TABLE `ads_member_asset_by_account_channel_d` (
|   `dt` date NOT NULL,
|   `member_id` bigint NOT NULL,
|   `account_type` varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
|   `net_asset_amount` decimal(38,10) DEFAULT NULL,
|   PRIMARY KEY (`dt`,`member_id`,`account_type`),
|   KEY `sindex` (`dt`,`member_id`,`account_type`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
| /*!50100 PARTITION BY RANGE (to_days(`dt`))
| (PARTITION p0 VALUES LESS THAN (739000) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (739001) ENGINE = InnoDB,
|  PARTITION p2 VALUES LESS THAN (739002) ENGINE = InnoDB,
|  PARTITION p3 VALUES LESS THAN (739003) ENGINE = InnoDB,
|  PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_a162ef44_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_a162ef44_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_a162ef44_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_a162ef44_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_d40c3e05
< ok 0 0
> exec
| CREATE DATABASE pa_mm_d40c3e05
< ok 1 0
> exec
| USE pa_mm_d40c3e05
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE `pa_mm_d40c3e05`.`ads_member_asset_by_account_channel_d` (
|   `dt` date NOT NULL,
|   `member_id` bigint NOT NULL,
|   `account_type` varchar(50) NOT NULL,
|   `net_asset_amount` decimal(38,10) DEFAULT NULL,
|   PRIMARY KEY (`dt`,`member_id`,`account_type`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| ALTER TABLE `pa_mm_d40c3e05`.`ads_member_asset_by_account_channel_d` ADD KEY `sindex` (`dt`,`member_id`,`account_type`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_d40c3e05`.`ads_member_asset_by_account_channel_d` /*!50100 PARTITION BY RANGE (to_days(`dt`))
| (PARTITION p0 VALUES LESS THAN (739000) ENGINE = innodb,
|  PARTITION p1 VALUES LESS THAN (739001) ENGINE = innodb,
|  PARTITION p2 VALUES LESS THAN (739002) ENGINE = innodb,
|  PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = innodb) */
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_d40c3e05_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_d40c3e05_to
< ok 1 0
> exec
| USE pa_mm_d40c3e05_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE ads_member_asset_by_account_channel_d (
|     dt DATE NOT NULL,
|     member_id BIGINT NOT NULL,
|     account_type VARCHAR(50) NOT NULL,
|     net_asset_amount DECIMAL(38,10) DEFAULT NULL,
|     PRIMARY KEY (dt, member_id, account_type),
|     KEY sindex (dt, member_id, account_type)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
| PARTITION BY RANGE (TO_DAYS(dt)) (
|     PARTITION p0 VALUES LESS THAN (739000),
|     PARTITION p1 VALUES LESS THAN (739001),
|     PARTITION p2 VALUES LESS THAN (739002),
|     PARTITION p_future VALUES LESS THAN MAXVALUE
| )
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_d40c3e05
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_d40c3e05
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| ads_member_asset_by_account_channel_d
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_d40c3e05`.`ads_member_asset_by_account_channel_d`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ads_member_asset_by_account_channel_d
: bytes
| CREATE TABLE `ads_member_asset_by_account_channel_d` (
|   `dt` date NOT NULL,
|   `member_id` bigint NOT NULL,
|   `account_type` varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
|   `net_asset_amount` decimal(38,10) DEFAULT NULL,
|   PRIMARY KEY (`dt`,`member_id`,`account_type`),
|   KEY `sindex` (`dt`,`member_id`,`account_type`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
| /*!50100 PARTITION BY RANGE (to_days(`dt`))
| (PARTITION p0 VALUES LESS THAN (739000) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (739001) ENGINE = InnoDB,
|  PARTITION p2 VALUES LESS THAN (739002) ENGINE = InnoDB,
|  PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_d40c3e05
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_d40c3e05
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_d40c3e05
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_d40c3e05
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_d40c3e05_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_d40c3e05_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| ads_member_asset_by_account_channel_d
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_d40c3e05_to`.`ads_member_asset_by_account_channel_d`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ads_member_asset_by_account_channel_d
: bytes
| CREATE TABLE `ads_member_asset_by_account_channel_d` (
|   `dt` date NOT NULL,
|   `member_id` bigint NOT NULL,
|   `account_type` varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
|   `net_asset_amount` decimal(38,10) DEFAULT NULL,
|   PRIMARY KEY (`dt`,`member_id`,`account_type`),
|   KEY `sindex` (`dt`,`member_id`,`account_type`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
| /*!50100 PARTITION BY RANGE (to_days(`dt`))
| (PARTITION p0 VALUES LESS THAN (739000) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN (739001) ENGINE = InnoDB,
|  PARTITION p2 VALUES LESS THAN (739002) ENGINE = InnoDB,
|  PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_d40c3e05_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_d40c3e05_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_d40c3e05_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_d40c3e05_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_20cd7d09
< ok 0 0
> exec
| CREATE DATABASE pa_mm_20cd7d09
< ok 1 0
> exec
| USE pa_mm_20cd7d09
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     price DECIMAL(10, 2) NOT NULL,
|     tax_rate DECIMAL(5, 2) DEFAULT 8.0,
|     price_with_tax DECIMAL(10, 2) AS (price * (1 + tax_rate / 100)) STORED,
|     name_upper VARCHAR(100) AS (UPPER(name)) VIRTUAL,
|     PRIMARY KEY (id),
|     INDEX idx_price_with_tax (price_with_tax),
|     INDEX idx_name_upper (name_upper)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE order_summary (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL,
|     unit_price DECIMAL(10, 2) NOT NULL,
|     total_price DECIMAL(10, 2) AS (quantity * unit_price) STORED,
|     order_date DATE NOT NULL,
|     order_year INT AS (YEAR(order_date)) VIRTUAL,
|     order_month INT AS (MONTH(order_date)) VIRTUAL,
|     PRIMARY KEY (id),
|     INDEX idx_year_month (order_year, order_month)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE `pa_mm_20cd7d09`.`order_summary`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_20cd7d09`.`products` DROP INDEX `idx_name_upper`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_20cd7d09`.`products` DROP INDEX `idx_price_with_tax`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_20cd7d09`.`products` DROP COLUMN `name_upper`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_20cd7d09`.`products` DROP COLUMN `price_with_tax`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_20cd7d09_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_20cd7d09_to
< ok 1 0
> exec
| USE pa_mm_20cd7d09_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     price DECIMAL(10, 2) NOT NULL,
|     tax_rate DECIMAL(5, 2) DEFAULT 8.0,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_20cd7d09
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_20cd7d09
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_20cd7d09`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `tax_rate` decimal(5,2) DEFAULT '8.00',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_20cd7d09
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_20cd7d09
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_20cd7d09
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_20cd7d09
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_20cd7d09_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_20cd7d09_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_20cd7d09_to`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `tax_rate` decimal(5,2) DEFAULT '8.00',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_20cd7d09_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_20cd7d09_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_20cd7d09_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_20cd7d09_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_25188286
< ok 0 0
> exec
| CREATE DATABASE pa_mm_25188286
< ok 1 0
> exec
| USE pa_mm_25188286
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE users (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     email VARCHAR(100) NOT NULL,
|     preferences JSON,
|     metadata JSON,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     details JSON NOT NULL,
|     tags JSON,
|     brand VARCHAR(50) AS (details->>'$.brand') VIRTUAL,
|     PRIMARY KEY (id),
|     INDEX idx_product_brand ((CAST(details->>'$.brand' AS CHAR(50)))),
|     INDEX idx_product_price ((CAST(details->>'$.price' AS DECIMAL(10,2)))),
|     INDEX idx_brand (brand),
|     CONSTRAINT chk_details_not_empty CHECK (JSON_TYPE(details) IS NOT NULL)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE `pa_mm_25188286`.`products`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_25188286`.`users` DROP COLUMN `metadata`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_25188286`.`users` DROP COLUMN `preferences`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_25188286_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_25188286_to
< ok 1 0
> exec
| USE pa_mm_25188286_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE users (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     email VARCHAR(100) NOT NULL,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_25188286
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_25188286
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| users
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_25188286`.`users`
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
|   `email` varchar(100) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_25188286
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_25188286
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_25188286
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_25188286
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_25188286_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_25188286_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| users
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_25188286_to`.`users`
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
|   `email` varchar(100) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_25188286_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_25188286_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_25188286_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_25188286_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_8cb89de1
< ok 0 0
> exec
| CREATE DATABASE pa_mm_8cb89de1
< ok 1 0
> exec
| USE pa_mm_8cb89de1
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE inventory (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL DEFAULT 0,
|     last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_product (product_id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE inventory_log (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     old_quantity INT,
|     new_quantity INT,
|     change_type VARCHAR(20),
|     changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TRIGGER trg_inventory_insert
| AFTER INSERT ON inventory
| FOR EACH ROW
| BEGIN
|     INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|     VALUES (NEW.product_id, NULL, NEW.quantity, 'INSERT');
| END
< ok 0 0
> exec
| CREATE TRIGGER trg_inventory_update
| AFTER UPDATE ON inventory
| FOR EACH ROW
| BEGIN
|     IF OLD.quantity != NEW.quantity THEN
|         INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|         VALUES (NEW.product_id, OLD.quantity, NEW.quantity, 'UPDATE');
|     END IF;
| END
< ok 0 0
> exec
| CREATE TRIGGER trg_inventory_delete
| BEFORE DELETE ON inventory
| FOR EACH ROW
| BEGIN
|     INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|     VALUES (OLD.product_id, OLD.quantity, NULL, 'DELETE');
| END
< ok 0 0
> exec
| CREATE VIEW inventory_changes AS
| SELECT
|     product_id,
|     COUNT(*) AS change_count,
|     MAX(changed_at) AS last_changed
| FROM inventory_log
| GROUP BY product_id
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP VIEW IF EXISTS `pa_mm_8cb89de1`.`inventory_changes`
< ok 0 0
> exec
| DROP TRIGGER `pa_mm_8cb89de1`.`trg_inventory_delete`
< ok 0 0
> exec
| DROP TRIGGER `pa_mm_8cb89de1`.`trg_inventory_insert`
< ok 0 0
> exec
| DROP TRIGGER `pa_mm_8cb89de1`.`trg_inventory_update`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_8cb89de1_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_8cb89de1_to
< ok 1 0
> exec
| USE pa_mm_8cb89de1_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE inventory (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL DEFAULT 0,
|     last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_product (product_id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE inventory_log (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     old_quantity INT,
|     new_quantity INT,
|     change_type VARCHAR(20),
|     changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_8cb89de1
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_8cb89de1
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| inventory
: bytes
| BASE TABLE
< row
: bytes
| inventory_log
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_8cb89de1`.`inventory`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| inventory
: bytes
| CREATE TABLE `inventory` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL DEFAULT '0',
|   `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_product` (`product_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_8cb89de1`.`inventory_log`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| inventory_log
: bytes
| CREATE TABLE `inventory_log` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `old_quantity` int DEFAULT NULL,
|   `new_quantity` int DEFAULT NULL,
|   `change_type` varchar(20) DEFAULT NULL,
|   `changed_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_8cb89de1
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_8cb89de1
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_8cb89de1
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_8cb89de1
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_8cb89de1_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_8cb89de1_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| inventory
: bytes
| BASE TABLE
< row
: bytes
| inventory_log
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_8cb89de1_to`.`inventory`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| inventory
: bytes
| CREATE TABLE `inventory` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL DEFAULT '0',
|   `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_product` (`product_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_8cb89de1_to`.`inventory_log`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| inventory_log
: bytes
| CREATE TABLE `inventory_log` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `old_quantity` int DEFAULT NULL,
|   `new_quantity` int DEFAULT NULL,
|   `change_type` varchar(20) DEFAULT NULL,
|   `changed_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_8cb89de1_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_8cb89de1_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_8cb89de1_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_8cb89de1_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_b3bd91bb
< ok 0 0
> exec
| CREATE DATABASE pa_mm_b3bd91bb
< ok 1 0
> exec
| USE pa_mm_b3bd91bb
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| --
| -- Table structure for `international_content`
| --
| CREATE TABLE `international_content` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `english_text` varchar(1000) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
|   `chinese_text` varchar(1000) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
|   `emoji_text` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| --
| -- Table structure for `messages`
| --
| CREATE TABLE `messages` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `content` text NOT NULL,
|   `title` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
|   `summary` varchar(500) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_title` (`title`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE `pa_mm_b3bd91bb`.`international_content`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_b3bd91bb`.`messages` DROP INDEX `idx_title`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_b3bd91bb`.`messages` DROP COLUMN `summary`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_b3bd91bb`.`messages` DROP COLUMN `title`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_b3bd91bb_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_b3bd91bb_to
< ok 1 0
> exec
| USE pa_mm_b3bd91bb_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE messages (
|     id INT NOT NULL AUTO_INCREMENT,
|     content TEXT NOT NULL,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_b3bd91bb
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_b3bd91bb
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| messages
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_b3bd91bb`.`messages`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| messages
: bytes
| CREATE TABLE `messages` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `content` text COLLATE utf8mb4_general_ci NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_b3bd91bb
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_b3bd91bb
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_b3bd91bb
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_b3bd91bb
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_b3bd91bb_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_b3bd91bb_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| messages
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_b3bd91bb_to`.`messages`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| messages
: bytes
| CREATE TABLE `messages` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `content` text COLLATE utf8mb4_general_ci NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_b3bd91bb_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_b3bd91bb_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_b3bd91bb_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_b3bd91bb_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_90b1f429
< ok 0 0
> exec
| CREATE DATABASE pa_mm_90b1f429
< ok 1 0
> exec
| USE pa_mm_90b1f429
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE system_stats (
|     id INT NOT NULL AUTO_INCREMENT,
|     stat_name VARCHAR(100) NOT NULL,
|     stat_value DECIMAL(10, 2),
|     recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE EVENT IF NOT EXISTS evt_cleanup_old_stats
| ON SCHEDULE EVERY 1 DAY
| STARTS CURRENT_TIMESTAMP
| DO
| BEGIN
|     DELETE FROM system_stats
|     WHERE recorded_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
| END
< ok 0 0
> exec
| CREATE PROCEDURE RecordSystemStat(
|     IN p_stat_name VARCHAR(100),
|     IN p_stat_value DECIMAL(10, 2)
| )
| BEGIN
|     DECLARE EXIT HANDLER FOR SQLEXCEPTION
|     BEGIN
|         ROLLBACK;
|         RESIGNAL;
|     END;
|
|     START TRANSACTION;
|
|     INSERT INTO system_stats (stat_name, stat_value)
|     VALUES (p_stat_name, p_stat_value);
|
|     COMMIT;
| END
< ok 0 0
> exec
| CREATE FUNCTION GetAverageStat(p_stat_name VARCHAR(100))
| RETURNS DECIMAL(10, 2)
| DETERMINISTIC
| READS SQL DATA
| BEGIN
|     DECLARE avg_value DECIMAL(10, 2);
|
|     SELECT AVG(stat_value) INTO avg_value
|     FROM system_stats
|     WHERE stat_name = p_stat_name
|     AND recorded_at >= DATE_SUB(NOW(), INTERVAL 7 DAY);
|
|     RETURN COALESCE(avg_value, 0);
| END
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP FUNCTION `pa_mm_90b1f429`.`GetAverageStat`
< ok 0 0
> exec
| DROP PROCEDURE `pa_mm_90b1f429`.`RecordSystemStat`
< ok 0 0
> exec
| DROP EVENT IF EXISTS `evt_cleanup_old_stats`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_90b1f429_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_90b1f429_to
< ok 1 0
> exec
| USE pa_mm_90b1f429_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE system_stats (
|     id INT NOT NULL AUTO_INCREMENT,
|     stat_name VARCHAR(100) NOT NULL,
|     stat_value DECIMAL(10, 2),
|     recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_90b1f429
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_90b1f429
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| system_stats
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_90b1f429`.`system_stats`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| system_stats
: bytes
| CREATE TABLE `system_stats` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `stat_name` varchar(100) NOT NULL,
|   `stat_value` decimal(10,2) DEFAULT NULL,
|   `recorded_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_90b1f429
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_90b1f429
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_90b1f429
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_90b1f429
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_90b1f429_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_90b1f429_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| system_stats
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_90b1f429_to`.`system_stats`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| system_stats
: bytes
| CREATE TABLE `system_stats` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `stat_name` varchar(100) NOT NULL,
|   `stat_value` decimal(10,2) DEFAULT NULL,
|   `recorded_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_90b1f429_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_90b1f429_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_90b1f429_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_90b1f429_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_f38f765b
< ok 0 0
> exec
| CREATE DATABASE pa_mm_f38f765b
< ok 1 0
> exec
| USE pa_mm_f38f765b
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| --
| -- Table structure for `categories`
| --
| CREATE TABLE `categories` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(50) NOT NULL COMMENT 'Category display name',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Product categories for organization'
< ok 0 0
> exec
| --
| -- Table structure for `products`
| --
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL COMMENT 'Product display name',
|   `price` decimal(10,2) NOT NULL COMMENT 'Product price in USD',
|   `description` text DEFAULT NULL COMMENT 'Detailed product description',
|   `category_id` int DEFAULT NULL COMMENT 'Foreign key reference to categories table',
|   `supplier_id` int DEFAULT NULL COMMENT 'Reference to supplier providing this product',
|   PRIMARY KEY (`id`),
|   KEY `idx_category` (`category_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Product catalog table containing all product information'
< ok 0 0
> exec
| --
| -- Table structure for `suppliers`
| --
| CREATE TABLE `suppliers` (
|   `id` int NOT NULL AUTO_INCREMENT COMMENT 'Unique supplier identifier',
|   `company_name` varchar(100) NOT NULL COMMENT 'Legal company name',
|   `contact_email` varchar(100) DEFAULT NULL COMMENT 'Primary contact email address',
|   `phone` varchar(20) DEFAULT NULL COMMENT 'Business phone number',
|   `address` text DEFAULT NULL COMMENT 'Full business address',
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation timestamp',
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`contact_email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Supplier information and contact details'
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE `pa_mm_f38f765b`.`suppliers`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_f38f765b`.`products` DROP COLUMN `supplier_id`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_f38f765b`.`categories` COMMENT=''
< ok 0 0
> exec
| ALTER TABLE `pa_mm_f38f765b`.`products` COMMENT=''
< ok 0 0
> exec
| ALTER TABLE `pa_mm_f38f765b`.`categories` MODIFY COLUMN `name` varchar(50) NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_f38f765b`.`products` MODIFY COLUMN `category_id` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_f38f765b`.`products` MODIFY COLUMN `description` text
< ok 0 0
> exec
| ALTER TABLE `pa_mm_f38f765b`.`products` MODIFY COLUMN `name` varchar(100) NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_f38f765b`.`products` MODIFY COLUMN `price` decimal(10,2) NOT NULL
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_f38f765b_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_f38f765b_to
< ok 1 0
> exec
| USE pa_mm_f38f765b_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     price DECIMAL(10, 2) NOT NULL,
|     description TEXT,
|     category_id INT,
|     PRIMARY KEY (id),
|     INDEX idx_category (category_id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE categories (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(50) NOT NULL,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_f38f765b
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_f38f765b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| categories
: bytes
| BASE TABLE
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_f38f765b`.`categories`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| categories
: bytes
| CREATE TABLE `categories` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(50) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_f38f765b`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `description` text,
|   `category_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_category` (`category_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_f38f765b
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_f38f765b
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_f38f765b
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_f38f765b
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_f38f765b_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_f38f765b_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| categories
: bytes
| BASE TABLE
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_f38f765b_to`.`categories`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| categories
: bytes
| CREATE TABLE `categories` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(50) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_f38f765b_to`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `description` text,
|   `category_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_category` (`category_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_f38f765b_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_f38f765b_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_f38f765b_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_f38f765b_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_cd13c195
< ok 0 0
> exec
| CREATE DATABASE pa_mm_cd13c195
< ok 1 0
> exec
| USE pa_mm_cd13c195
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| --
| -- Table structure for `order_items`
| --
| CREATE TABLE `order_items` (
|   `id` int NOT NULL AUTO_INCREMENT COMMENT 'Line item identifier',
|   `order_id` int NOT NULL COMMENT 'Reference to parent order',
|   `product_name` varchar(150) NOT NULL COMMENT 'Full product name including variant details',
|   `quantity` int NOT NULL COMMENT 'Quantity ordered (must be positive)',
|   `unit_price` decimal(10,2) NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_order` (`order_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| --
| -- Table structure for `orders`
| --
| CREATE TABLE `orders` (
|   `id` int NOT NULL AUTO_INCREMENT COMMENT 'Order unique identifier',
|   `customer_name` varchar(100) NOT NULL COMMENT 'Full name of the purchasing customer',
|   `order_date` date NOT NULL,
|   `total_amount` decimal(10,2) DEFAULT NULL,
|   `status` varchar(20) DEFAULT 'pending' COMMENT 'Order processing status (pending, processing, shipped, delivered, cancelled)',
|   `notes` text DEFAULT NULL COMMENT 'Special delivery instructions and customer notes',
|   PRIMARY KEY (`id`),
|   KEY `idx_date` (`order_date`),
|   KEY `idx_status` (`status`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Customer purchase orders with tracking information'
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_cd13c195`.`order_items` COMMENT='Individual line items for each order'
< ok 0 0
> exec
| ALTER TABLE `pa_mm_cd13c195`.`orders` COMMENT='Customer orders and order details'
< ok 0 0
> exec
| ALTER TABLE `pa_mm_cd13c195`.`order_items` MODIFY COLUMN `product_name` varchar(100) NOT NULL COMMENT 'Name of ordered product'
< ok 0 0
> exec
| ALTER TABLE `pa_mm_cd13c195`.`order_items` MODIFY COLUMN `quantity` int NOT NULL COMMENT 'Number of items ordered'
< ok 0 0
> exec
| ALTER TABLE `pa_mm_cd13c195`.`order_items` MODIFY COLUMN `unit_price` decimal(10,2) NOT NULL COMMENT 'Price per individual item'
< ok 0 0
> exec
| ALTER TABLE `pa_mm_cd13c195`.`orders` MODIFY COLUMN `customer_name` varchar(100) NOT NULL COMMENT 'Customer full name'
< ok 0 0
> exec
| ALTER TABLE `pa_mm_cd13c195`.`orders` MODIFY COLUMN `notes` text COMMENT 'Additional order notes or special instructions'
< ok 0 0
> exec
| ALTER TABLE `pa_mm_cd13c195`.`orders` MODIFY COLUMN `order_date` date NOT NULL COMMENT 'Date when order was placed'
< ok 0 0
> exec
| ALTER TABLE `pa_mm_cd13c195`.`orders` MODIFY COLUMN `status` varchar(20) DEFAULT 'pending' COMMENT 'Current order status'
< ok 0 0
> exec
| ALTER TABLE `pa_mm_cd13c195`.`orders` MODIFY COLUMN `total_amount` decimal(10,2) DEFAULT NULL COMMENT 'Total order amount in USD'
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_cd13c195_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_cd13c195_to
< ok 1 0
> exec
| USE pa_mm_cd13c195_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE orders (
|     id INT NOT NULL AUTO_INCREMENT COMMENT 'Order unique identifier',
|     customer_name VARCHAR(100) NOT NULL COMMENT 'Customer full name',
|     order_date DATE NOT NULL COMMENT 'Date when order was placed',
|     total_amount DECIMAL(10, 2) COMMENT 'Total order amount in USD',
|     status VARCHAR(20) DEFAULT 'pending' COMMENT 'Current order status',
|     notes TEXT COMMENT 'Additional order notes or special instructions',
|     PRIMARY KEY (id),
|     INDEX idx_date (order_date),
|     INDEX idx_status (status)
| ) ENGINE=InnoDB COMMENT = 'Customer orders and order details'
< ok 0 0
> exec
| CREATE TABLE order_items (
|     id INT NOT NULL AUTO_INCREMENT COMMENT 'Line item identifier',
|     order_id INT NOT NULL COMMENT 'Reference to parent order',
|     product_name VARCHAR(100) NOT NULL COMMENT 'Name of ordered product',
|     quantity INT NOT NULL COMMENT 'Number of items ordered',
|     unit_price DECIMAL(10, 2) NOT NULL COMMENT 'Price per individual item',
|     PRIMARY KEY (id),
|     INDEX idx_order (order_id)
| ) ENGINE=InnoDB COMMENT = 'Individual line items for each order'
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_cd13c195
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_cd13c195
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| order_items
: bytes
| BASE TABLE
< row
: bytes
| orders
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_cd13c195`.`order_items`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| order_items
: bytes
| CREATE TABLE `order_items` (
|   `id` int NOT NULL AUTO_INCREMENT COMMENT 'Line item identifier',
|   `order_id` int NOT NULL COMMENT 'Reference to parent order',
|   `product_name` varchar(100) NOT NULL COMMENT 'Name of ordered product',
|   `quantity` int NOT NULL COMMENT 'Number of items ordered',
|   `unit_price` decimal(10,2) NOT NULL COMMENT 'Price per individual item',
|   PRIMARY KEY (`id`),
|   KEY `idx_order` (`order_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Individual line items for each order'
> query
| SHOW CREATE TABLE `pa_mm_cd13c195`.`orders`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| orders
: bytes
| CREATE TABLE `orders` (
|   `id` int NOT NULL AUTO_INCREMENT COMMENT 'Order unique identifier',
|   `customer_name` varchar(100) NOT NULL COMMENT 'Customer full name',
|   `order_date` date NOT NULL COMMENT 'Date when order was placed',
|   `total_amount` decimal(10,2) DEFAULT NULL COMMENT 'Total order amount in USD',
|   `status` varchar(20) DEFAULT 'pending' COMMENT 'Current order status',
|   `notes` text COMMENT 'Additional order notes or special instructions',
|   PRIMARY KEY (`id`),
|   KEY `idx_date` (`order_date`),
|   KEY `idx_status` (`status`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Customer orders and order details'
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_cd13c195
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_cd13c195
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_cd13c195
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_cd13c195
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_cd13c195_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_cd13c195_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| order_items
: bytes
| BASE TABLE
< row
: bytes
| orders
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_cd13c195_to`.`order_items`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| order_items
: bytes
| CREATE TABLE `order_items` (
|   `id` int NOT NULL AUTO_INCREMENT COMMENT 'Line item identifier',
|   `order_id` int NOT NULL COMMENT 'Reference to parent order',
|   `product_name` varchar(100) NOT NULL COMMENT 'Name of ordered product',
|   `quantity` int NOT NULL COMMENT 'Number of items ordered',
|   `unit_price` decimal(10,2) NOT NULL COMMENT 'Price per individual item',
|   PRIMARY KEY (`id`),
|   KEY `idx_order` (`order_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Individual line items for each order'
> query
| SHOW CREATE TABLE `pa_mm_cd13c195_to`.`orders`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| orders
: bytes
| CREATE TABLE `orders` (
|   `id` int NOT NULL AUTO_INCREMENT COMMENT 'Order unique identifier',
|   `customer_name` varchar(100) NOT NULL COMMENT 'Customer full name',
|   `order_date` date NOT NULL COMMENT 'Date when order was placed',
|   `total_amount` decimal(10,2) DEFAULT NULL COMMENT 'Total order amount in USD',
|   `status` varchar(20) DEFAULT 'pending' COMMENT 'Current order status',
|   `notes` text COMMENT 'Additional order notes or special instructions',
|   PRIMARY KEY (`id`),
|   KEY `idx_date` (`order_date`),
|   KEY `idx_status` (`status`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Customer orders and order details'
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_cd13c195_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_cd13c195_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_cd13c195_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_cd13c195_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_bfb1dfc0
< ok 0 0
> exec
| CREATE DATABASE pa_mm_bfb1dfc0
< ok 1 0
> exec
| USE pa_mm_bfb1dfc0
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| --
| -- Table structure for `analytics`
| --
| CREATE TABLE `analytics` (
|   `id` int NOT NULL AUTO_INCREMENT COMMENT 'Primary key (auto-increment)',
|   `event_name` varchar(100) NOT NULL COMMENT 'Event identifier - format: "page_view", "button_click", etc.',
|   `event_data` json COMMENT 'Event payload: {"user_id": 123, "timestamp": "2023-12-01T10:30:00Z", "properties": {...}}',
|   `ip_address` varchar(45) DEFAULT NULL COMMENT 'Client IP address (IPv4: xxx.xxx.xxx.xxx or IPv6: xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx)',
|   `user_agent` text DEFAULT NULL COMMENT 'Browser user agent string - may contain "Mozilla/5.0", various browser/OS info',
|   `referrer` varchar(500) DEFAULT NULL COMMENT 'HTTP referrer URL - where user came from (can be NULL)',
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT 'Event timestamp - UTC timezone',
|   `processed` tinyint(1) DEFAULT '0' COMMENT 'Processing status: TRUE = processed, FALSE = pending',
|   PRIMARY KEY (`id`),
|   KEY `idx_event_date` (`event_name`, `created_at`),
|   KEY `idx_processed` (`processed`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Analytics events tracking - stores user interactions & system events. Data retention: 2 years. Access level: "admin" & "analyst" roles only.'
< ok 0 0
> exec
| --
| -- Table structure for `posts`
| --
| CREATE TABLE `posts` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `user_id` int NOT NULL,
|   `title` varchar(200) NOT NULL COMMENT 'Post title - supports emojis 📝✨🔥💡🎉 and Unicode characters',
|   `content` longtext DEFAULT NULL COMMENT 'Post content in HTML format: <p>, <strong>, <em>, <a href="...">, <img src="..."/>',
|   `metadata` json COMMENT 'Post metadata: {"tags": ["tag1", "tag2"], "category": "tech", "featured": true, "views": 0}',
|   PRIMARY KEY (`id`),
|   KEY `idx_user` (`user_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Blog posts and articles - supports international content (中文, العربية, Русский, 日本語, 한국어, Français, Español, Deutsch)'
< ok 0 0
> exec
| --
| -- Table structure for `users`
| --
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) NOT NULL COMMENT 'User''s chosen "display name" for the platform',
|   `email` varchar(100) NOT NULL COMMENT 'Primary email address - must be "unique" across all users',
|   `bio` text DEFAULT NULL COMMENT 'User biography text
| Can contain multiple lines
| and various formatting',
|   `preferences` json COMMENT 'User settings: theme, notifications, privacy & security options (@, #, $, %, ^, &, *, +, =, |, \\, /, ?, <, >)',
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='User accounts - stores user''s personal information and preferences'
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE `pa_mm_bfb1dfc0`.`analytics`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_bfb1dfc0`.`posts` COMMENT=''
< ok 0 0
> exec
| ALTER TABLE `pa_mm_bfb1dfc0`.`users` COMMENT=''
< ok 0 0
> exec
| ALTER TABLE `pa_mm_bfb1dfc0`.`posts` MODIFY COLUMN `content` longtext
< ok 0 0
> exec
| ALTER TABLE `pa_mm_bfb1dfc0`.`posts` MODIFY COLUMN `metadata` json
< ok 0 0
> exec
| ALTER TABLE `pa_mm_bfb1dfc0`.`posts` MODIFY COLUMN `title` varchar(200) NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_bfb1dfc0`.`users` MODIFY COLUMN `bio` text
< ok 0 0
> exec
| ALTER TABLE `pa_mm_bfb1dfc0`.`users` MODIFY COLUMN `email` varchar(100) NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_bfb1dfc0`.`users` MODIFY COLUMN `preferences` json
< ok 0 0
> exec
| ALTER TABLE `pa_mm_bfb1dfc0`.`users` MODIFY COLUMN `username` varchar(50) NOT NULL
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_bfb1dfc0_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_bfb1dfc0_to
< ok 1 0
> exec
| USE pa_mm_bfb1dfc0_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE users (
|     id INT NOT NULL AUTO_INCREMENT,
|     username VARCHAR(50) NOT NULL,
|     email VARCHAR(100) NOT NULL,
|     bio TEXT,
|     preferences JSON,
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_email (email)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE posts (
|     id INT NOT NULL AUTO_INCREMENT,
|     user_id INT NOT NULL,
|     title VARCHAR(200) NOT NULL,
|     content LONGTEXT,
|     metadata JSON,
|     PRIMARY KEY (id),
|     INDEX idx_user (user_id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_bfb1dfc0
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_bfb1dfc0
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
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
| SHOW CREATE TABLE `pa_mm_bfb1dfc0`.`posts`
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
|   `content` longtext,
|   `metadata` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_user` (`user_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_bfb1dfc0`.`users`
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
|   `bio` text,
|   `preferences` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_bfb1dfc0
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_bfb1dfc0
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_bfb1dfc0
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_bfb1dfc0
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_bfb1dfc0_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_bfb1dfc0_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
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
| SHOW CREATE TABLE `pa_mm_bfb1dfc0_to`.`posts`
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
|   `content` longtext,
|   `metadata` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_user` (`user_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_bfb1dfc0_to`.`users`
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
|   `bio` text,
|   `preferences` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_bfb1dfc0_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_bfb1dfc0_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_bfb1dfc0_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_bfb1dfc0_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_ffa6ed7d
< ok 0 0
> exec
| CREATE DATABASE pa_mm_ffa6ed7d
< ok 1 0
> exec
| USE pa_mm_ffa6ed7d
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE `pa_mm_ffa6ed7d`.`authors` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `email` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE `pa_mm_ffa6ed7d`.`books` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `title` varchar(200) NOT NULL,
|   `author_id` int NOT NULL,
|   `isbn` varchar(20) DEFAULT NULL,
|   `published_year` int DEFAULT NULL,
|   `price` decimal(8,2) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE `pa_mm_ffa6ed7d`.`products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(50) NOT NULL,
|   `price` decimal(8,2) NOT NULL,
|   `description` text,
|   `category` varchar(30) DEFAULT NULL,
|   `is_active` tinyint(1) DEFAULT '1',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_mm_ffa6ed7d`.`authors` ADD UNIQUE KEY `uk_email` (`email`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_ffa6ed7d`.`books` ADD KEY `idx_author` (`author_id`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_ffa6ed7d`.`books` ADD KEY `idx_year` (`published_year`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_ffa6ed7d`.`books` ADD UNIQUE KEY `uk_isbn` (`isbn`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_ffa6ed7d`.`products` ADD KEY `idx_category` (`category`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_ffa6ed7d`.`products` ADD KEY `idx_price` (`price`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_ffa6ed7d`.`books` ADD CONSTRAINT `chk_price_positive` CHECK ((`price` > 0))
< ok 0 0
> exec
| ALTER TABLE `pa_mm_ffa6ed7d`.`books` ADD CONSTRAINT `chk_year_valid` CHECK (((`published_year` >= 1000) and (`published_year` <= 2100)))
< ok 0 0
> exec
| ALTER TABLE `pa_mm_ffa6ed7d`.`books` ADD CONSTRAINT `fk_author` FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_ffa6ed7d_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_ffa6ed7d_to
< ok 1 0
> exec
| USE pa_mm_ffa6ed7d_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE authors (
|   id int NOT NULL AUTO_INCREMENT,
|   name varchar(100) NOT NULL,
|   email varchar(100) DEFAULT NULL,
|   PRIMARY KEY (id),
|   UNIQUE KEY uk_email (email)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE books (
|   id int NOT NULL AUTO_INCREMENT,
|   title varchar(200) NOT NULL,
|   author_id int NOT NULL,
|   isbn varchar(20) DEFAULT NULL,
|   published_year int DEFAULT NULL,
|   price decimal(8,2) DEFAULT NULL,
|   PRIMARY KEY (id),
|   KEY idx_author (author_id),
|   KEY idx_year (published_year),
|   UNIQUE KEY uk_isbn (isbn),
|   CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES authors (id),
|   CONSTRAINT chk_price_positive CHECK (price > 0),
|   CONSTRAINT chk_year_valid CHECK ((published_year >= 1000) and (published_year <= 2100))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE products (
|   id int NOT NULL AUTO_INCREMENT,
|   name varchar(50) NOT NULL,
|   price decimal(8,2) NOT NULL,
|   description text DEFAULT NULL,
|   category varchar(30) DEFAULT NULL,
|   is_active tinyint(1) DEFAULT '1',
|   PRIMARY KEY (id),
|   KEY idx_category (category),
|   KEY idx_price (price)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_ffa6ed7d
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_ffa6ed7d
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| authors
: bytes
| BASE TABLE
< row
: bytes
| books
: bytes
| BASE TABLE
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_ffa6ed7d`.`authors`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| authors
: bytes
| CREATE TABLE `authors` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) COLLATE utf8mb4_general_ci NOT NULL,
|   `email` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
> query
| SHOW CREATE TABLE `pa_mm_ffa6ed7d`.`books`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| books
: bytes
| CREATE TABLE `books` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `title` varchar(200) COLLATE utf8mb4_general_ci NOT NULL,
|   `author_id` int NOT NULL,
|   `isbn` varchar(20) COLLATE utf8mb4_general_ci DEFAULT NULL,
|   `published_year` int DEFAULT NULL,
|   `price` decimal(8,2) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_isbn` (`isbn`),
|   KEY `idx_author` (`author_id`),
|   KEY `idx_year` (`published_year`),
|   CONSTRAINT `fk_author` FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`),
|   CONSTRAINT `chk_price_positive` CHECK ((`price` > 0)),
|   CONSTRAINT `chk_year_valid` CHECK (((`published_year` >= 1000) and (`published_year` <= 2100)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
> query
| SHOW CREATE TABLE `pa_mm_ffa6ed7d`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(50) NOT NULL,
|   `price` decimal(8,2) NOT NULL,
|   `description` text,
|   `category` varchar(30) DEFAULT NULL,
|   `is_active` tinyint(1) DEFAULT '1',
|   PRIMARY KEY (`id`),
|   KEY `idx_category` (`category`),
|   KEY `idx_price` (`price`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_ffa6ed7d
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_ffa6ed7d
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_ffa6ed7d
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_ffa6ed7d
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_ffa6ed7d_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_ffa6ed7d_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| authors
: bytes
| BASE TABLE
< row
: bytes
| books
: bytes
| BASE TABLE
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_ffa6ed7d_to`.`authors`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| authors
: bytes
| CREATE TABLE `authors` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) COLLATE utf8mb4_general_ci NOT NULL,
|   `email` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
> query
| SHOW CREATE TABLE `pa_mm_ffa6ed7d_to`.`books`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| books
: bytes
| CREATE TABLE `books` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `title` varchar(200) COLLATE utf8mb4_general_ci NOT NULL,
|   `author_id` int NOT NULL,
|   `isbn` varchar(20) COLLATE utf8mb4_general_ci DEFAULT NULL,
|   `published_year` int DEFAULT NULL,
|   `price` decimal(8,2) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_isbn` (`isbn`),
|   KEY `idx_author` (`author_id`),
|   KEY `idx_year` (`published_year`),
|   CONSTRAINT `fk_author` FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`),
|   CONSTRAINT `chk_price_positive` CHECK ((`price` > 0)),
|   CONSTRAINT `chk_year_valid` CHECK (((`published_year` >= 1000) and (`published_year` <= 2100)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
> query
| SHOW CREATE TABLE `pa_mm_ffa6ed7d_to`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(50) NOT NULL,
|   `price` decimal(8,2) NOT NULL,
|   `description` text,
|   `category` varchar(30) DEFAULT NULL,
|   `is_active` tinyint(1) DEFAULT '1',
|   PRIMARY KEY (`id`),
|   KEY `idx_category` (`category`),
|   KEY `idx_price` (`price`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_ffa6ed7d_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_ffa6ed7d_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_ffa6ed7d_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_ffa6ed7d_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_217cfbd9
< ok 0 0
> exec
| CREATE DATABASE pa_mm_217cfbd9
< ok 1 0
> exec
| USE pa_mm_217cfbd9
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_217cfbd9`.`posts` DROP FOREIGN KEY `fk_user`
< ok 0 0
> exec
| DROP TABLE `pa_mm_217cfbd9`.`posts`
< ok 0 0
> exec
| DROP TABLE `pa_mm_217cfbd9`.`users`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_217cfbd9_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_217cfbd9_to
< ok 1 0
> exec
| USE pa_mm_217cfbd9_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_217cfbd9
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_217cfbd9
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_217cfbd9
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_217cfbd9
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_217cfbd9
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_217cfbd9
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_217cfbd9_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_217cfbd9_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_217cfbd9_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_217cfbd9_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_217cfbd9_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_217cfbd9_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_8c191378
< ok 0 0
> exec
| CREATE DATABASE pa_mm_8c191378
< ok 1 0
> exec
| USE pa_mm_8c191378
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE `pa_mm_8c191378`.`comments` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `post_id` int NOT NULL,
|   `user_id` int NOT NULL,
|   `content` text NOT NULL,
|   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| ALTER TABLE `pa_mm_8c191378`.`users` ADD COLUMN `is_active` tinyint(1) DEFAULT '1'
< ok 0 0
> exec
| ALTER TABLE `pa_mm_8c191378`.`comments` ADD KEY `idx_post_user` (`post_id`,`user_id`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_8c191378`.`users` ADD KEY `idx_email_active` (`email`,`is_active`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_8c191378`.`posts` ADD CONSTRAINT `chk_title_length` CHECK ((char_length(`title`) > 0))
< ok 0 0
> exec
| ALTER TABLE `pa_mm_8c191378`.`comments` ADD CONSTRAINT `fk_comment_post` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE
< ok 0 0
> exec
| ALTER TABLE `pa_mm_8c191378`.`comments` ADD CONSTRAINT `fk_comment_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_8c191378_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_8c191378_to
< ok 1 0
> exec
| USE pa_mm_8c191378_to
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_8c191378
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_8c191378
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| comments
: bytes
| BASE TABLE
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
| SHOW CREATE TABLE `pa_mm_8c191378`.`comments`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| comments
: bytes
| CREATE TABLE `comments` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `post_id` int NOT NULL,
|   `user_id` int NOT NULL,
|   `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
|   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   KEY `idx_post_user` (`post_id`,`user_id`),
|   KEY `fk_comment_user` (`user_id`),
|   CONSTRAINT `fk_comment_post` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
|   CONSTRAINT `fk_comment_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SHOW CREATE TABLE `pa_mm_8c191378`.`posts`
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
|   `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `content` text COLLATE utf8mb4_unicode_ci,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_user_id` (`user_id`),
|   CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
|   CONSTRAINT `chk_title_length` CHECK ((char_length(`title`) > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SHOW CREATE TABLE `pa_mm_8c191378`.`users`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   `is_active` tinyint(1) DEFAULT '1',
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`),
|   KEY `idx_username` (`username`),
|   KEY `idx_email_active` (`email`,`is_active`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_8c191378
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_8c191378
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_8c191378
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_8c191378
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_8c191378_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_8c191378_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| comments
: bytes
| BASE TABLE
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
| SHOW CREATE TABLE `pa_mm_8c191378_to`.`comments`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| comments
: bytes
| CREATE TABLE `comments` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `post_id` int NOT NULL,
|   `user_id` int NOT NULL,
|   `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
|   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   KEY `idx_post_user` (`post_id`,`user_id`),
|   KEY `fk_comment_user` (`user_id`),
|   CONSTRAINT `fk_comment_post` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
|   CONSTRAINT `fk_comment_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SHOW CREATE TABLE `pa_mm_8c191378_to`.`posts`
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
|   `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `content` text COLLATE utf8mb4_unicode_ci,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_user_id` (`user_id`),
|   CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
|   CONSTRAINT `chk_title_length` CHECK ((char_length(`title`) > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SHOW CREATE TABLE `pa_mm_8c191378_to`.`users`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
|   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   `is_active` tinyint(1) DEFAULT '1',
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`),
|   KEY `idx_username` (`username`),
|   KEY `idx_email_active` (`email`,`is_active`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_8c191378_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_8c191378_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_8c191378_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_8c191378_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_857bdb94
< ok 0 0
> exec
| CREATE DATABASE pa_mm_857bdb94
< ok 1 0
> exec
| USE pa_mm_857bdb94
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     price DECIMAL(10, 2) NOT NULL,
|     stock INT DEFAULT 0,
|     PRIMARY KEY (id),
|     INDEX idx_name (name)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE orders (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL,
|     total DECIMAL(10, 2),
|     order_date DATETIME DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE PROCEDURE `pa_mm_857bdb94`.`GetProductInventory`(IN `product_name` VARCHAR(100))
| BEGIN
|     SELECT * FROM product_inventory
|     WHERE name LIKE CONCAT('%', product_name, '%');
| END
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `pa_mm_857bdb94`.`low_stock_products` AS select `products`.`id` AS `id`,`products`.`name` AS `name`,`products`.`price` AS `price`,`products`.`stock` AS `stock` from `products` where (`products`.`stock` < 10)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `pa_mm_857bdb94`.`product_inventory` AS select `p`.`id` AS `id`,`p`.`name` AS `name`,`p`.`price` AS `price`,`p`.`stock` AS `stock`,coalesce(sum(`o`.`quantity`),0) AS `total_ordered` from (`products` `p` left join `orders` `o` on((`p`.`id` = `o`.`product_id`))) group by `p`.`id`,`p`.`name`,`p`.`price`,`p`.`stock`
< ok 0 0
> exec
| CREATE TRIGGER `pa_mm_857bdb94`.`update_order_total` BEFORE INSERT ON `pa_mm_857bdb94`.`orders` FOR EACH ROW BEGIN
|     DECLARE product_price DECIMAL(10, 2);
|     SELECT price INTO product_price FROM products WHERE id = NEW.product_id;
|     SET NEW.total = NEW.quantity * product_price;
| END
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_857bdb94_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_857bdb94_to
< ok 1 0
> exec
| USE pa_mm_857bdb94_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     price DECIMAL(10, 2) NOT NULL,
|     stock INT DEFAULT 0,
|     PRIMARY KEY (id),
|     INDEX idx_name (name)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE orders (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL,
|     total DECIMAL(10, 2),
|     order_date DATETIME DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE VIEW product_inventory AS
| SELECT
|     p.id,
|     p.name,
|     p.price,
|     p.stock,
|     COALESCE(SUM(o.quantity), 0) AS total_ordered
| FROM products p
| LEFT JOIN orders o ON p.id = o.product_id
| GROUP BY p.id, p.name, p.price, p.stock
< ok 0 0
> exec
| CREATE VIEW low_stock_products AS
| SELECT * FROM products WHERE stock < 10
< ok 0 0
> exec
| CREATE TRIGGER update_order_total
| BEFORE INSERT ON orders
| FOR EACH ROW
| BEGIN
|     DECLARE product_price DECIMAL(10, 2);
|     SELECT price INTO product_price FROM products WHERE id = NEW.product_id;
|     SET NEW.total = NEW.quantity * product_price;
| END
< ok 0 0
> exec
| CREATE PROCEDURE GetProductInventory(IN product_name VARCHAR(100))
| BEGIN
|     SELECT * FROM product_inventory
|     WHERE name LIKE CONCAT('%', product_name, '%');
| END
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_857bdb94
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_857bdb94
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| low_stock_products
: bytes
| VIEW
< row
: bytes
| orders
: bytes
| BASE TABLE
< row
: bytes
| product_inventory
: bytes
| VIEW
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `pa_mm_857bdb94`.`low_stock_products`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| low_stock_products
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `low_stock_products` AS select `products`.`id` AS `id`,`products`.`name` AS `name`,`products`.`price` AS `price`,`products`.`stock` AS `stock` from `products` where (`products`.`stock` < 10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `pa_mm_857bdb94`.`orders`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| orders
: bytes
| CREATE TABLE `orders` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL,
|   `total` decimal(10,2) DEFAULT NULL,
|   `order_date` datetime DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   KEY `fk_product` (`product_id`),
|   CONSTRAINT `fk_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW `pa_mm_857bdb94`.`product_inventory`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| product_inventory
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `product_inventory` AS select `p`.`id` AS `id`,`p`.`name` AS `name`,`p`.`price` AS `price`,`p`.`stock` AS `stock`,coalesce(sum(`o`.`quantity`),0) AS `total_ordered` from (`products` `p` left join `orders` `o` on((`p`.`id` = `o`.`product_id`))) group by `p`.`id`,`p`.`name`,`p`.`price`,`p`.`stock`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `pa_mm_857bdb94`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `stock` int DEFAULT '0',
|   PRIMARY KEY (`id`),
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_857bdb94
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
< row
: bytes
| PROCEDURE
: bytes
| GetProductInventory
: null
: bytes
| BEGIN
|     SELECT * FROM product_inventory
|     WHERE name LIKE CONCAT('%', product_name, '%');
| END
: bytes
| NO
: bytes
| CONTAINS SQL
: bytes
| DEFINER
: bytes
|
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_857bdb94
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< row
: bytes
| GetProductInventory
: uint64 1
: bytes
| IN
: bytes
| product_name
: bytes
| varchar(100)
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_857bdb94
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| update_order_total
: bytes
| INSERT
: bytes
| orders
: bytes
| BEFORE
: int64 1
: bytes
| BEGIN
|     DECLARE product_price DECIMAL(10, 2);
|     SELECT price INTO product_price FROM products WHERE id = NEW.product_id;
|     SET NEW.total = NEW.quantity * product_price;
| END
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_857bdb94
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_857bdb94_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_857bdb94_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| low_stock_products
: bytes
| VIEW
< row
: bytes
| orders
: bytes
| BASE TABLE
< row
: bytes
| product_inventory
: bytes
| VIEW
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `pa_mm_857bdb94_to`.`low_stock_products`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| low_stock_products
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `low_stock_products` AS select `products`.`id` AS `id`,`products`.`name` AS `name`,`products`.`price` AS `price`,`products`.`stock` AS `stock` from `products` where (`products`.`stock` < 10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `pa_mm_857bdb94_to`.`orders`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| orders
: bytes
| CREATE TABLE `orders` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL,
|   `total` decimal(10,2) DEFAULT NULL,
|   `order_date` datetime DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   KEY `fk_product` (`product_id`),
|   CONSTRAINT `fk_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW `pa_mm_857bdb94_to`.`product_inventory`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| product_inventory
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `product_inventory` AS select `p`.`id` AS `id`,`p`.`name` AS `name`,`p`.`price` AS `price`,`p`.`stock` AS `stock`,coalesce(sum(`o`.`quantity`),0) AS `total_ordered` from (`products` `p` left join `orders` `o` on((`p`.`id` = `o`.`product_id`))) group by `p`.`id`,`p`.`name`,`p`.`price`,`p`.`stock`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `pa_mm_857bdb94_to`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `stock` int DEFAULT '0',
|   PRIMARY KEY (`id`),
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_857bdb94_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
< row
: bytes
| PROCEDURE
: bytes
| GetProductInventory
: null
: bytes
| BEGIN
|     SELECT * FROM product_inventory
|     WHERE name LIKE CONCAT('%', product_name, '%');
| END
: bytes
| NO
: bytes
| CONTAINS SQL
: bytes
| DEFINER
: bytes
|
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_857bdb94_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< row
: bytes
| GetProductInventory
: uint64 1
: bytes
| IN
: bytes
| product_name
: bytes
| varchar(100)
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_857bdb94_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| update_order_total
: bytes
| INSERT
: bytes
| orders
: bytes
| BEFORE
: int64 1
: bytes
| BEGIN
|     DECLARE product_price DECIMAL(10, 2);
|     SELECT price INTO product_price FROM products WHERE id = NEW.product_id;
|     SET NEW.total = NEW.quantity * product_price;
| END
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_857bdb94_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_a42f2c1
< ok 0 0
> exec
| CREATE DATABASE pa_mm_a42f2c1
< ok 1 0
> exec
| USE pa_mm_a42f2c1
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE sales (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL,
|     unit_price DECIMAL(10, 2) NOT NULL,
|     sale_date DATE NOT NULL,
|     PRIMARY KEY (id),
|     INDEX idx_product_date (product_id, sale_date)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE `pa_mm_a42f2c1`.`sales_summary` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `month_year` varchar(7) NOT NULL,
|   `total_sales` decimal(10,2) DEFAULT NULL,
|   `tax_amount` decimal(10,2) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_mm_a42f2c1`.`sales_summary` ADD UNIQUE KEY `uk_product_month` (`product_id`,`month_year`)
< ok 0 0
> exec
| CREATE FUNCTION `pa_mm_a42f2c1`.`CalculateTotalSales`(`p_product_id` INT, `p_start_date` DATE, `p_end_date` DATE) RETURNS decimal(10,2)
|     DETERMINISTIC
|     READS SQL DATA
| BEGIN
|     DECLARE total DECIMAL(10, 2);
|
|     SELECT COALESCE(SUM(quantity * unit_price), 0) INTO total
|     FROM sales
|     WHERE product_id = p_product_id
|     AND sale_date BETWEEN p_start_date AND p_end_date;
|
|     RETURN total;
| END
< ok 0 0
> exec
| CREATE FUNCTION `pa_mm_a42f2c1`.`GetTaxAmount`(`amount` DECIMAL(10,2), `tax_rate` DECIMAL(5,2)) RETURNS decimal(10,2)
|     DETERMINISTIC
|     NO SQL
| BEGIN
|     RETURN amount * (tax_rate / 100);
| END
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_a42f2c1_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_a42f2c1_to
< ok 1 0
> exec
| USE pa_mm_a42f2c1_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE sales (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL,
|     unit_price DECIMAL(10, 2) NOT NULL,
|     sale_date DATE NOT NULL,
|     PRIMARY KEY (id),
|     INDEX idx_product_date (product_id, sale_date)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE sales_summary (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     month_year VARCHAR(7) NOT NULL,
|     total_sales DECIMAL(10, 2),
|     tax_amount DECIMAL(10, 2),
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_product_month (product_id, month_year)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE FUNCTION CalculateTotalSales(p_product_id INT, p_start_date DATE, p_end_date DATE)
| RETURNS DECIMAL(10, 2)
| DETERMINISTIC
| READS SQL DATA
| BEGIN
|     DECLARE total DECIMAL(10, 2);
|
|     SELECT COALESCE(SUM(quantity * unit_price), 0) INTO total
|     FROM sales
|     WHERE product_id = p_product_id
|     AND sale_date BETWEEN p_start_date AND p_end_date;
|
|     RETURN total;
| END
< ok 0 0
> exec
| CREATE FUNCTION GetTaxAmount(amount DECIMAL(10, 2), tax_rate DECIMAL(5, 2))
| RETURNS DECIMAL(10, 2)
| DETERMINISTIC
| NO SQL
| BEGIN
|     RETURN amount * (tax_rate / 100);
| END
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_a42f2c1
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_a42f2c1
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| sales
: bytes
| BASE TABLE
< row
: bytes
| sales_summary
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_a42f2c1`.`sales`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sales
: bytes
| CREATE TABLE `sales` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL,
|   `unit_price` decimal(10,2) NOT NULL,
|   `sale_date` date NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_product_date` (`product_id`,`sale_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_a42f2c1`.`sales_summary`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sales_summary
: bytes
| CREATE TABLE `sales_summary` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `month_year` varchar(7) NOT NULL,
|   `total_sales` decimal(10,2) DEFAULT NULL,
|   `tax_amount` decimal(10,2) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_product_month` (`product_id`,`month_year`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_a42f2c1
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
< row
: bytes
| FUNCTION
: bytes
| CalculateTotalSales
: bytes
| decimal(10,2)
: bytes
| BEGIN
|     DECLARE total DECIMAL(10, 2);
|
|     SELECT COALESCE(SUM(quantity * unit_price), 0) INTO total
|     FROM sales
|     WHERE product_id = p_product_id
|     AND sale_date BETWEEN p_start_date AND p_end_date;
|
|     RETURN total;
| END
: bytes
| YES
: bytes
| READS SQL DATA
: bytes
| DEFINER
: bytes
|
< row
: bytes
| FUNCTION
: bytes
| GetTaxAmount
: bytes
| decimal(10,2)
: bytes
| BEGIN
|     RETURN amount * (tax_rate / 100);
| END
: bytes
| YES
: bytes
| NO SQL
: bytes
| DEFINER
: bytes
|
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_a42f2c1
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< row
: bytes
| CalculateTotalSales
: uint64 0
: null
: null
: bytes
| decimal(10,2)
< row
: bytes
| CalculateTotalSales
: uint64 1
: bytes
| IN
: bytes
| p_product_id
: bytes
| int
< row
: bytes
| CalculateTotalSales
: uint64 2
: bytes
| IN
: bytes
| p_start_date
: bytes
| date
< row
: bytes
| CalculateTotalSales
: uint64 3
: bytes
| IN
: bytes
| p_end_date
: bytes
| date
< row
: bytes
| GetTaxAmount
: uint64 0
: null
: null
: bytes
| decimal(10,2)
< row
: bytes
| GetTaxAmount
: uint64 1
: bytes
| IN
: bytes
| amount
: bytes
| decimal(10,2)
< row
: bytes
| GetTaxAmount
: uint64 2
: bytes
| IN
: bytes
| tax_rate
: bytes
| decimal(5,2)
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_a42f2c1
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_a42f2c1
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_a42f2c1_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_a42f2c1_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| sales
: bytes
| BASE TABLE
< row
: bytes
| sales_summary
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_a42f2c1_to`.`sales`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sales
: bytes
| CREATE TABLE `sales` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL,
|   `unit_price` decimal(10,2) NOT NULL,
|   `sale_date` date NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_product_date` (`product_id`,`sale_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_a42f2c1_to`.`sales_summary`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sales_summary
: bytes
| CREATE TABLE `sales_summary` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `month_year` varchar(7) NOT NULL,
|   `total_sales` decimal(10,2) DEFAULT NULL,
|   `tax_amount` decimal(10,2) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_product_month` (`product_id`,`month_year`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_a42f2c1_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
< row
: bytes
| FUNCTION
: bytes
| CalculateTotalSales
: bytes
| decimal(10,2)
: bytes
| BEGIN
|     DECLARE total DECIMAL(10, 2);
|
|     SELECT COALESCE(SUM(quantity * unit_price), 0) INTO total
|     FROM sales
|     WHERE product_id = p_product_id
|     AND sale_date BETWEEN p_start_date AND p_end_date;
|
|     RETURN total;
| END
: bytes
| YES
: bytes
| READS SQL DATA
: bytes
| DEFINER
: bytes
|
< row
: bytes
| FUNCTION
: bytes
| GetTaxAmount
: bytes
| decimal(10,2)
: bytes
| BEGIN
|     RETURN amount * (tax_rate / 100);
| END
: bytes
| YES
: bytes
| NO SQL
: bytes
| DEFINER
: bytes
|
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_a42f2c1_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< row
: bytes
| CalculateTotalSales
: uint64 0
: null
: null
: bytes
| decimal(10,2)
< row
: bytes
| CalculateTotalSales
: uint64 1
: bytes
| IN
: bytes
| p_product_id
: bytes
| int
< row
: bytes
| CalculateTotalSales
: uint64 2
: bytes
| IN
: bytes
| p_start_date
: bytes
| date
< row
: bytes
| CalculateTotalSales
: uint64 3
: bytes
| IN
: bytes
| p_end_date
: bytes
| date
< row
: bytes
| GetTaxAmount
: uint64 0
: null
: null
: bytes
| decimal(10,2)
< row
: bytes
| GetTaxAmount
: uint64 1
: bytes
| IN
: bytes
| amount
: bytes
| decimal(10,2)
< row
: bytes
| GetTaxAmount
: uint64 2
: bytes
| IN
: bytes
| tax_rate
: bytes
| decimal(5,2)
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_a42f2c1_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_a42f2c1_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_1474c605
< ok 0 0
> exec
| CREATE DATABASE pa_mm_1474c605
< ok 1 0
> exec
| USE pa_mm_1474c605
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     category VARCHAR(50),
|     price DECIMAL(10, 2),
|     sku VARCHAR(50) NOT NULL,
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_sku (sku),
|     INDEX idx_name (name),
|     INDEX idx_category (category),
|     INDEX idx_price (price),
|     CONSTRAINT chk_price_positive CHECK (price > 0),
|     CONSTRAINT chk_name_length CHECK (CHAR_LENGTH(name) >= 3)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1474c605`.`products` DROP CHECK `chk_name_length`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1474c605`.`products` DROP CHECK `chk_price_positive`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1474c605`.`products` DROP INDEX `idx_category`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1474c605`.`products` DROP INDEX `idx_name`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1474c605`.`products` DROP INDEX `idx_price`
< ok 0 0
> exec
| ALTER TABLE `pa_mm_1474c605`.`products` DROP INDEX `uk_sku`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_1474c605_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_1474c605_to
< ok 1 0
> exec
| USE pa_mm_1474c605_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     category VARCHAR(50),
|     price DECIMAL(10, 2),
|     sku VARCHAR(50) NOT NULL,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_1474c605
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_1474c605
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_1474c605`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `category` varchar(50) DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `sku` varchar(50) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_1474c605
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_1474c605
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_1474c605
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_1474c605
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_1474c605_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_1474c605_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_1474c605_to`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `category` varchar(50) DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `sku` varchar(50) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_1474c605_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_1474c605_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_1474c605_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_1474c605_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_4dfac7b
< ok 0 0
> exec
| CREATE DATABASE pa_mm_4dfac7b
< ok 1 0
> exec
| USE pa_mm_4dfac7b
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(50) NOT NULL,
|     price DECIMAL(8, 2) NOT NULL,
|     description TEXT,
|     category VARCHAR(30),
|     is_active TINYINT(1) DEFAULT 1,
|     PRIMARY KEY (id),
|     INDEX idx_category (category),
|     INDEX idx_price (price)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4dfac7b`.`products` MODIFY COLUMN `category` varchar(50) DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4dfac7b`.`products` ADD COLUMN `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4dfac7b`.`products` MODIFY COLUMN `description` text NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4dfac7b`.`products` MODIFY COLUMN `name` varchar(100) NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4dfac7b`.`products` MODIFY COLUMN `price` decimal(10,2) NOT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4dfac7b`.`products` ADD COLUMN `stock_quantity` int DEFAULT '0'
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4dfac7b`.`products` ADD COLUMN `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4dfac7b`.`products` ADD COLUMN `weight` decimal(5,2) DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4dfac7b`.`products` ADD KEY `idx_created_at` (`created_at`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4dfac7b`.`products` ADD UNIQUE KEY `uk_name_category` (`name`,`category`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4dfac7b`.`products` ADD CONSTRAINT `chk_price_positive` CHECK ((`price` > 0))
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4dfac7b`.`products` ADD CONSTRAINT `chk_stock_non_negative` CHECK ((`stock_quantity` >= 0))
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_4dfac7b_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_4dfac7b_to
< ok 1 0
> exec
| USE pa_mm_4dfac7b_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     price DECIMAL(10, 2) NOT NULL,
|     description TEXT NOT NULL,
|     category VARCHAR(50),
|     is_active TINYINT(1) DEFAULT 1,
|     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|     stock_quantity INT DEFAULT 0,
|     weight DECIMAL(5, 2),
|     PRIMARY KEY (id),
|     INDEX idx_category (category),
|     INDEX idx_price (price),
|     INDEX idx_created_at (created_at),
|     UNIQUE INDEX uk_name_category (name, category),
|     CONSTRAINT chk_price_positive CHECK (price > 0),
|     CONSTRAINT chk_stock_non_negative CHECK (stock_quantity >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_4dfac7b
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_4dfac7b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_4dfac7b`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `description` text NOT NULL,
|   `category` varchar(50) DEFAULT NULL,
|   `is_active` tinyint(1) DEFAULT '1',
|   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   `stock_quantity` int DEFAULT '0',
|   `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `weight` decimal(5,2) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_name_category` (`name`,`category`),
|   KEY `idx_category` (`category`),
|   KEY `idx_price` (`price`),
|   KEY `idx_created_at` (`created_at`),
|   CONSTRAINT `chk_price_positive` CHECK ((`price` > 0)),
|   CONSTRAINT `chk_stock_non_negative` CHECK ((`stock_quantity` >= 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_4dfac7b
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_4dfac7b
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_4dfac7b
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_4dfac7b
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_4dfac7b_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_4dfac7b_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_4dfac7b_to`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `description` text NOT NULL,
|   `category` varchar(50) DEFAULT NULL,
|   `is_active` tinyint(1) DEFAULT '1',
|   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `stock_quantity` int DEFAULT '0',
|   `weight` decimal(5,2) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_name_category` (`name`,`category`),
|   KEY `idx_category` (`category`),
|   KEY `idx_price` (`price`),
|   KEY `idx_created_at` (`created_at`),
|   CONSTRAINT `chk_price_positive` CHECK ((`price` > 0)),
|   CONSTRAINT `chk_stock_non_negative` CHECK ((`stock_quantity` >= 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_4dfac7b_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_4dfac7b_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_4dfac7b_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_4dfac7b_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_482448c6
< ok 0 0
> exec
| CREATE DATABASE pa_mm_482448c6
< ok 1 0
> exec
| USE pa_mm_482448c6
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE articles (
|     id INT NOT NULL AUTO_INCREMENT,
|     title VARCHAR(200) NOT NULL,
|     content TEXT NOT NULL,
|     tags VARCHAR(500),
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE locations (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     position POINT NOT NULL,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_482448c6`.`articles` ADD FULLTEXT KEY `idx_fulltext_content` (`title`,`content`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_482448c6`.`articles` ADD FULLTEXT KEY `idx_fulltext_tags` (`tags`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_482448c6`.`articles` ADD KEY `idx_title_tags` (`title`(50),`tags`(50))
< ok 0 0
> exec
| ALTER TABLE `pa_mm_482448c6`.`locations` ADD SPATIAL KEY `idx_spatial_position` (`position`)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_482448c6_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_482448c6_to
< ok 1 0
> exec
| USE pa_mm_482448c6_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE articles (
|     id INT NOT NULL AUTO_INCREMENT,
|     title VARCHAR(200) NOT NULL,
|     content TEXT NOT NULL,
|     tags VARCHAR(500),
|     PRIMARY KEY (id),
|     FULLTEXT idx_fulltext_content (title, content),
|     FULLTEXT idx_fulltext_tags (tags),
|     INDEX idx_title_tags (title(50), tags(50))
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE locations (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     position POINT NOT NULL,
|     PRIMARY KEY (id),
|     SPATIAL INDEX idx_spatial_position (position)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_482448c6
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_482448c6
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| articles
: bytes
| BASE TABLE
< row
: bytes
| locations
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_482448c6`.`articles`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| articles
: bytes
| CREATE TABLE `articles` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `title` varchar(200) NOT NULL,
|   `content` text NOT NULL,
|   `tags` varchar(500) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_title_tags` (`title`(50),`tags`(50)),
|   FULLTEXT KEY `idx_fulltext_content` (`title`,`content`),
|   FULLTEXT KEY `idx_fulltext_tags` (`tags`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_482448c6`.`locations`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| locations
: bytes
| CREATE TABLE `locations` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `position` point NOT NULL,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `idx_spatial_position` (`position`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_482448c6
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_482448c6
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_482448c6
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_482448c6
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_482448c6_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_482448c6_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| articles
: bytes
| BASE TABLE
< row
: bytes
| locations
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_482448c6_to`.`articles`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| articles
: bytes
| CREATE TABLE `articles` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `title` varchar(200) NOT NULL,
|   `content` text NOT NULL,
|   `tags` varchar(500) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_title_tags` (`title`(50),`tags`(50)),
|   FULLTEXT KEY `idx_fulltext_content` (`title`,`content`),
|   FULLTEXT KEY `idx_fulltext_tags` (`tags`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_482448c6_to`.`locations`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| locations
: bytes
| CREATE TABLE `locations` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `position` point NOT NULL,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `idx_spatial_position` (`position`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_482448c6_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_482448c6_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_482448c6_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_482448c6_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_970419a3
< ok 0 0
> exec
| CREATE DATABASE pa_mm_970419a3
< ok 1 0
> exec
| USE pa_mm_970419a3
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     price DECIMAL(10, 2) NOT NULL,
|     tax_rate DECIMAL(5, 2) DEFAULT 8.0,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE `pa_mm_970419a3`.`order_summary` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL,
|   `unit_price` decimal(10,2) NOT NULL,
|   `total_price` decimal(10,2) GENERATED ALWAYS AS (((`quantity` * `unit_price`))) STORED,
|   `order_date` date NOT NULL,
|   `order_year` int GENERATED ALWAYS AS (year(`order_date`)) VIRTUAL,
|   `order_month` int GENERATED ALWAYS AS (month(`order_date`)) VIRTUAL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_mm_970419a3`.`products` ADD COLUMN `name_upper` varchar(100) GENERATED ALWAYS AS (upper(`name`)) VIRTUAL
< ok 0 0
> exec
| ALTER TABLE `pa_mm_970419a3`.`products` ADD COLUMN `price_with_tax` decimal(10,2) GENERATED ALWAYS AS (((`price` * ((1 + ((`tax_rate` / 100))))))) STORED
< ok 0 0
> exec
| ALTER TABLE `pa_mm_970419a3`.`order_summary` ADD KEY `idx_year_month` (`order_year`,`order_month`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_970419a3`.`products` ADD KEY `idx_name_upper` (`name_upper`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_970419a3`.`products` ADD KEY `idx_price_with_tax` (`price_with_tax`)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_970419a3_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_970419a3_to
< ok 1 0
> exec
| USE pa_mm_970419a3_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     price DECIMAL(10, 2) NOT NULL,
|     tax_rate DECIMAL(5, 2) DEFAULT 8.0,
|     price_with_tax DECIMAL(10, 2) AS (price * (1 + tax_rate / 100)) STORED,
|     name_upper VARCHAR(100) AS (UPPER(name)) VIRTUAL,
|     PRIMARY KEY (id),
|     INDEX idx_price_with_tax (price_with_tax),
|     INDEX idx_name_upper (name_upper)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE order_summary (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL,
|     unit_price DECIMAL(10, 2) NOT NULL,
|     total_price DECIMAL(10, 2) AS (quantity * unit_price) STORED,
|     order_date DATE NOT NULL,
|     order_year INT AS (YEAR(order_date)) VIRTUAL,
|     order_month INT AS (MONTH(order_date)) VIRTUAL,
|     PRIMARY KEY (id),
|     INDEX idx_year_month (order_year, order_month)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_970419a3
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_970419a3
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| order_summary
: bytes
| BASE TABLE
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_970419a3`.`order_summary`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| order_summary
: bytes
| CREATE TABLE `order_summary` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL,
|   `unit_price` decimal(10,2) NOT NULL,
|   `total_price` decimal(10,2) GENERATED ALWAYS AS ((`quantity` * `unit_price`)) STORED,
|   `order_date` date NOT NULL,
|   `order_year` int GENERATED ALWAYS AS (year(`order_date`)) VIRTUAL,
|   `order_month` int GENERATED ALWAYS AS (month(`order_date`)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   KEY `idx_year_month` (`order_year`,`order_month`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_970419a3`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `tax_rate` decimal(5,2) DEFAULT '8.00',
|   `name_upper` varchar(100) GENERATED ALWAYS AS (upper(`name`)) VIRTUAL,
|   `price_with_tax` decimal(10,2) GENERATED ALWAYS AS ((`price` * (1 + (`tax_rate` / 100)))) STORED,
|   PRIMARY KEY (`id`),
|   KEY `idx_price_with_tax` (`price_with_tax`),
|   KEY `idx_name_upper` (`name_upper`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_970419a3
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_970419a3
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_970419a3
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_970419a3
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_970419a3_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_970419a3_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| order_summary
: bytes
| BASE TABLE
< row
: bytes
| products
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_970419a3_to`.`order_summary`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| order_summary
: bytes
| CREATE TABLE `order_summary` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL,
|   `unit_price` decimal(10,2) NOT NULL,
|   `total_price` decimal(10,2) GENERATED ALWAYS AS ((`quantity` * `unit_price`)) STORED,
|   `order_date` date NOT NULL,
|   `order_year` int GENERATED ALWAYS AS (year(`order_date`)) VIRTUAL,
|   `order_month` int GENERATED ALWAYS AS (month(`order_date`)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   KEY `idx_year_month` (`order_year`,`order_month`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_970419a3_to`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `tax_rate` decimal(5,2) DEFAULT '8.00',
|   `price_with_tax` decimal(10,2) GENERATED ALWAYS AS ((`price` * (1 + (`tax_rate` / 100)))) STORED,
|   `name_upper` varchar(100) GENERATED ALWAYS AS (upper(`name`)) VIRTUAL,
|   PRIMARY KEY (`id`),
|   KEY `idx_price_with_tax` (`price_with_tax`),
|   KEY `idx_name_upper` (`name_upper`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_970419a3_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_970419a3_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_970419a3_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_970419a3_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_de41a819
< ok 0 0
> exec
| CREATE DATABASE pa_mm_de41a819
< ok 1 0
> exec
| USE pa_mm_de41a819
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE users (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     email VARCHAR(100) NOT NULL,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE `pa_mm_de41a819`.`products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `details` json NOT NULL,
|   `tags` json,
|   `brand` varchar(50) GENERATED ALWAYS AS (json_unquote(json_extract(`details`,_utf8mb4'$.brand'))) VIRTUAL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_mm_de41a819`.`users` ADD COLUMN `metadata` json
< ok 0 0
> exec
| ALTER TABLE `pa_mm_de41a819`.`users` ADD COLUMN `preferences` json
< ok 0 0
> exec
| ALTER TABLE `pa_mm_de41a819`.`products` ADD KEY `idx_brand` (`brand`)
< ok 0 0
> exec
| ALTER TABLE `pa_mm_de41a819`.`products` ADD KEY `idx_product_brand` ((cast(json_unquote(json_extract(`details`,_utf8mb4'$.brand')) as char(50) charset utf8mb4)))
< ok 0 0
> exec
| ALTER TABLE `pa_mm_de41a819`.`products` ADD KEY `idx_product_price` ((cast(json_unquote(json_extract(`details`,_utf8mb4'$.price')) as decimal(10,2))))
< ok 0 0
> exec
| ALTER TABLE `pa_mm_de41a819`.`products` ADD CONSTRAINT `chk_details_not_empty` CHECK ((json_type(`details`) is not null))
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_de41a819_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_de41a819_to
< ok 1 0
> exec
| USE pa_mm_de41a819_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE users (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     email VARCHAR(100) NOT NULL,
|     preferences JSON,
|     metadata JSON,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE products (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) NOT NULL,
|     details JSON NOT NULL,
|     tags JSON,
|     brand VARCHAR(50) AS (details->>'$.brand') VIRTUAL,
|     PRIMARY KEY (id),
|     INDEX idx_product_brand ((CAST(details->>'$.brand' AS CHAR(50)))),
|     INDEX idx_product_price ((CAST(details->>'$.price' AS DECIMAL(10,2)))),
|     INDEX idx_brand (brand),
|     CONSTRAINT chk_details_not_empty CHECK (JSON_TYPE(details) IS NOT NULL)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_de41a819
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_de41a819
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| products
: bytes
| BASE TABLE
< row
: bytes
| users
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_de41a819`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `details` json NOT NULL,
|   `tags` json DEFAULT NULL,
|   `brand` varchar(50) GENERATED ALWAYS AS (json_unquote(json_extract(`details`,_utf8mb4'$.brand'))) VIRTUAL,
|   PRIMARY KEY (`id`),
|   KEY `idx_brand` (`brand`),
|   KEY `idx_product_brand` ((cast(json_unquote(json_extract(`details`,_utf8mb4'$.brand')) as char(50) charset utf8mb4))),
|   KEY `idx_product_price` ((cast(json_unquote(json_extract(`details`,_utf8mb4'$.price')) as decimal(10,2)))),
|   CONSTRAINT `chk_details_not_empty` CHECK ((json_type(`details`) is not null))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_de41a819`.`users`
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
|   `email` varchar(100) NOT NULL,
|   `metadata` json DEFAULT NULL,
|   `preferences` json DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_de41a819
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_de41a819
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_de41a819
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_de41a819
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_de41a819_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_de41a819_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| products
: bytes
| BASE TABLE
< row
: bytes
| users
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_de41a819_to`.`products`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| products
: bytes
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `details` json NOT NULL,
|   `tags` json DEFAULT NULL,
|   `brand` varchar(50) GENERATED ALWAYS AS (json_unquote(json_extract(`details`,_utf8mb4'$.brand'))) VIRTUAL,
|   PRIMARY KEY (`id`),
|   KEY `idx_product_brand` ((cast(json_unquote(json_extract(`details`,_utf8mb4'$.brand')) as char(50) charset utf8mb4))),
|   KEY `idx_product_price` ((cast(json_unquote(json_extract(`details`,_utf8mb4'$.price')) as decimal(10,2)))),
|   KEY `idx_brand` (`brand`),
|   CONSTRAINT `chk_details_not_empty` CHECK ((json_type(`details`) is not null))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_de41a819_to`.`users`
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
|   `email` varchar(100) NOT NULL,
|   `preferences` json DEFAULT NULL,
|   `metadata` json DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_de41a819_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_de41a819_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_de41a819_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_de41a819_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_ed63de90
< ok 0 0
> exec
| CREATE DATABASE pa_mm_ed63de90
< ok 1 0
> exec
| USE pa_mm_ed63de90
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE inventory (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL DEFAULT 0,
|     last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_product (product_id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE inventory_log (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     old_quantity INT,
|     new_quantity INT,
|     change_type VARCHAR(20),
|     changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `pa_mm_ed63de90`.`inventory_changes` AS select `inventory_log`.`product_id` AS `product_id`,count(0) AS `change_count`,max(`inventory_log`.`changed_at`) AS `last_changed` from `inventory_log` group by `inventory_log`.`product_id`
< ok 0 0
> exec
| CREATE TRIGGER `pa_mm_ed63de90`.`trg_inventory_delete` BEFORE DELETE ON `pa_mm_ed63de90`.`inventory` FOR EACH ROW BEGIN
|     INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|     VALUES (OLD.product_id, OLD.quantity, NULL, 'DELETE');
| END
< ok 0 0
> exec
| CREATE TRIGGER `pa_mm_ed63de90`.`trg_inventory_insert` AFTER INSERT ON `pa_mm_ed63de90`.`inventory` FOR EACH ROW BEGIN
|     INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|     VALUES (NEW.product_id, NULL, NEW.quantity, 'INSERT');
| END
< ok 0 0
> exec
| CREATE TRIGGER `pa_mm_ed63de90`.`trg_inventory_update` AFTER UPDATE ON `pa_mm_ed63de90`.`inventory` FOR EACH ROW BEGIN
|     IF OLD.quantity != NEW.quantity THEN
|         INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|         VALUES (NEW.product_id, OLD.quantity, NEW.quantity, 'UPDATE');
|     END IF;
| END
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_ed63de90_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_ed63de90_to
< ok 1 0
> exec
| USE pa_mm_ed63de90_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE inventory (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     quantity INT NOT NULL DEFAULT 0,
|     last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|     PRIMARY KEY (id),
|     UNIQUE KEY uk_product (product_id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE inventory_log (
|     id INT NOT NULL AUTO_INCREMENT,
|     product_id INT NOT NULL,
|     old_quantity INT,
|     new_quantity INT,
|     change_type VARCHAR(20),
|     changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TRIGGER trg_inventory_insert
| AFTER INSERT ON inventory
| FOR EACH ROW
| BEGIN
|     INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|     VALUES (NEW.product_id, NULL, NEW.quantity, 'INSERT');
| END
< ok 0 0
> exec
| CREATE TRIGGER trg_inventory_update
| AFTER UPDATE ON inventory
| FOR EACH ROW
| BEGIN
|     IF OLD.quantity != NEW.quantity THEN
|         INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|         VALUES (NEW.product_id, OLD.quantity, NEW.quantity, 'UPDATE');
|     END IF;
| END
< ok 0 0
> exec
| CREATE TRIGGER trg_inventory_delete
| BEFORE DELETE ON inventory
| FOR EACH ROW
| BEGIN
|     INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|     VALUES (OLD.product_id, OLD.quantity, NULL, 'DELETE');
| END
< ok 0 0
> exec
| CREATE VIEW inventory_changes AS
| SELECT
|     product_id,
|     COUNT(*) AS change_count,
|     MAX(changed_at) AS last_changed
| FROM inventory_log
| GROUP BY product_id
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_ed63de90
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_ed63de90
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| inventory
: bytes
| BASE TABLE
< row
: bytes
| inventory_changes
: bytes
| VIEW
< row
: bytes
| inventory_log
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_ed63de90`.`inventory`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| inventory
: bytes
| CREATE TABLE `inventory` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL DEFAULT '0',
|   `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_product` (`product_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW `pa_mm_ed63de90`.`inventory_changes`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| inventory_changes
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `inventory_changes` AS select `inventory_log`.`product_id` AS `product_id`,count(0) AS `change_count`,max(`inventory_log`.`changed_at`) AS `last_changed` from `inventory_log` group by `inventory_log`.`product_id`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `pa_mm_ed63de90`.`inventory_log`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| inventory_log
: bytes
| CREATE TABLE `inventory_log` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `old_quantity` int DEFAULT NULL,
|   `new_quantity` int DEFAULT NULL,
|   `change_type` varchar(20) DEFAULT NULL,
|   `changed_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_ed63de90
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_ed63de90
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_ed63de90
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| trg_inventory_delete
: bytes
| DELETE
: bytes
| inventory
: bytes
| BEFORE
: int64 1
: bytes
| BEGIN
|     INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|     VALUES (OLD.product_id, OLD.quantity, NULL, 'DELETE');
| END
< row
: bytes
| trg_inventory_insert
: bytes
| INSERT
: bytes
| inventory
: bytes
| AFTER
: int64 1
: bytes
| BEGIN
|     INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|     VALUES (NEW.product_id, NULL, NEW.quantity, 'INSERT');
| END
< row
: bytes
| trg_inventory_update
: bytes
| UPDATE
: bytes
| inventory
: bytes
| AFTER
: int64 1
: bytes
| BEGIN
|     IF OLD.quantity != NEW.quantity THEN
|         INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|         VALUES (NEW.product_id, OLD.quantity, NEW.quantity, 'UPDATE');
|     END IF;
| END
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_ed63de90
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_ed63de90_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_ed63de90_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| inventory
: bytes
| BASE TABLE
< row
: bytes
| inventory_changes
: bytes
| VIEW
< row
: bytes
| inventory_log
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_ed63de90_to`.`inventory`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| inventory
: bytes
| CREATE TABLE `inventory` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL DEFAULT '0',
|   `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_product` (`product_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW `pa_mm_ed63de90_to`.`inventory_changes`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| inventory_changes
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `inventory_changes` AS select `inventory_log`.`product_id` AS `product_id`,count(0) AS `change_count`,max(`inventory_log`.`changed_at`) AS `last_changed` from `inventory_log` group by `inventory_log`.`product_id`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `pa_mm_ed63de90_to`.`inventory_log`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| inventory_log
: bytes
| CREATE TABLE `inventory_log` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `old_quantity` int DEFAULT NULL,
|   `new_quantity` int DEFAULT NULL,
|   `change_type` varchar(20) DEFAULT NULL,
|   `changed_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_ed63de90_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_ed63de90_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_ed63de90_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| trg_inventory_delete
: bytes
| DELETE
: bytes
| inventory
: bytes
| BEFORE
: int64 1
: bytes
| BEGIN
|     INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|     VALUES (OLD.product_id, OLD.quantity, NULL, 'DELETE');
| END
< row
: bytes
| trg_inventory_insert
: bytes
| INSERT
: bytes
| inventory
: bytes
| AFTER
: int64 1
: bytes
| BEGIN
|     INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|     VALUES (NEW.product_id, NULL, NEW.quantity, 'INSERT');
| END
< row
: bytes
| trg_inventory_update
: bytes
| UPDATE
: bytes
| inventory
: bytes
| AFTER
: int64 1
: bytes
| BEGIN
|     IF OLD.quantity != NEW.quantity THEN
|         INSERT INTO inventory_log (product_id, old_quantity, new_quantity, change_type)
|         VALUES (NEW.product_id, OLD.quantity, NEW.quantity, 'UPDATE');
|     END IF;
| END
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_ed63de90_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS pa_mm_2f3a1b76
< ok 0 0
> exec
| CREATE DATABASE pa_mm_2f3a1b76
< ok 1 0
> exec
| USE pa_mm_2f3a1b76
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE system_stats (
|     id INT NOT NULL AUTO_INCREMENT,
|     stat_name VARCHAR(100) NOT NULL,
|     stat_value DECIMAL(10, 2),
|     recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE FUNCTION `pa_mm_2f3a1b76`.`GetAverageStat`(`p_stat_name` VARCHAR(100)) RETURNS decimal(10,2)
|     DETERMINISTIC
|     READS SQL DATA
| BEGIN
|     DECLARE avg_value DECIMAL(10, 2);
|
|     SELECT AVG(stat_value) INTO avg_value
|     FROM system_stats
|     WHERE stat_name = p_stat_name
|     AND recorded_at >= DATE_SUB(NOW(), INTERVAL 7 DAY);
|
|     RETURN COALESCE(avg_value, 0);
| END
< ok 0 0
> exec
| CREATE PROCEDURE `pa_mm_2f3a1b76`.`RecordSystemStat`(IN `p_stat_name` VARCHAR(100), IN `p_stat_value` DECIMAL(10,2))
| BEGIN
|     DECLARE EXIT HANDLER FOR SQLEXCEPTION
|     BEGIN
|         ROLLBACK;
|         RESIGNAL;
|     END;
|
|     START TRANSACTION;
|
|     INSERT INTO system_stats (stat_name, stat_value)
|     VALUES (p_stat_name, p_stat_value);
|
|     COMMIT;
| END
< ok 0 0
> exec
| CREATE EVENT `evt_cleanup_old_stats` ON SCHEDULE EVERY 1 DAY STARTS '2026-11-14 11:18:16' ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO BEGIN
|     DELETE FROM system_stats
|     WHERE recorded_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
| END
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_2f3a1b76_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_2f3a1b76_to
< ok 1 0
> exec
| USE pa_mm_2f3a1b76_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE system_stats (
|     id INT NOT NULL AUTO_INCREMENT,
|     stat_name VARCHAR(100) NOT NULL,
|     stat_value DECIMAL(10, 2),
|     recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE EVENT IF NOT EXISTS evt_cleanup_old_stats
| ON SCHEDULE EVERY 1 DAY
| STARTS CURRENT_TIMESTAMP
| DO
| BEGIN
|     DELETE FROM system_stats
|     WHERE recorded_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
| END
< ok 0 0
> exec
| CREATE PROCEDURE RecordSystemStat(
|     IN p_stat_name VARCHAR(100),
|     IN p_stat_value DECIMAL(10, 2)
| )
| BEGIN
|     DECLARE EXIT HANDLER FOR SQLEXCEPTION
|     BEGIN
|         ROLLBACK;
|         RESIGNAL;
|     END;
|
|     START TRANSACTION;
|
|     INSERT INTO system_stats (stat_name, stat_value)
|     VALUES (p_stat_name, p_stat_value);
|
|     COMMIT;
| END
< ok 0 0
> exec
| CREATE FUNCTION GetAverageStat(p_stat_name VARCHAR(100))
| RETURNS DECIMAL(10, 2)
| DETERMINISTIC
| READS SQL DATA
| BEGIN
|     DECLARE avg_value DECIMAL(10, 2);
|
|     SELECT AVG(stat_value) INTO avg_value
|     FROM system_stats
|     WHERE stat_name = p_stat_name
|     AND recorded_at >= DATE_SUB(NOW(), INTERVAL 7 DAY);
|
|     RETURN COALESCE(avg_value, 0);
| END
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_2f3a1b76
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_2f3a1b76
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| system_stats
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_2f3a1b76`.`system_stats`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| system_stats
: bytes
| CREATE TABLE `system_stats` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `stat_name` varchar(100) NOT NULL,
|   `stat_value` decimal(10,2) DEFAULT NULL,
|   `recorded_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_2f3a1b76
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
< row
: bytes
| FUNCTION
: bytes
| GetAverageStat
: bytes
| decimal(10,2)
: bytes
| BEGIN
|     DECLARE avg_value DECIMAL(10, 2);
|
|     SELECT AVG(stat_value) INTO avg_value
|     FROM system_stats
|     WHERE stat_name = p_stat_name
|     AND recorded_at >= DATE_SUB(NOW(), INTERVAL 7 DAY);
|
|     RETURN COALESCE(avg_value, 0);
| END
: bytes
| YES
: bytes
| READS SQL DATA
: bytes
| DEFINER
: bytes
|
< row
: bytes
| PROCEDURE
: bytes
| RecordSystemStat
: null
: bytes
| BEGIN
|     DECLARE EXIT HANDLER FOR SQLEXCEPTION
|     BEGIN
|         ROLLBACK;
|         RESIGNAL;
|     END;
|
|     START TRANSACTION;
|
|     INSERT INTO system_stats (stat_name, stat_value)
|     VALUES (p_stat_name, p_stat_value);
|
|     COMMIT;
| END
: bytes
| NO
: bytes
| CONTAINS SQL
: bytes
| DEFINER
: bytes
|
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_2f3a1b76
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< row
: bytes
| GetAverageStat
: uint64 0
: null
: null
: bytes
| decimal(10,2)
< row
: bytes
| GetAverageStat
: uint64 1
: bytes
| IN
: bytes
| p_stat_name
: bytes
| varchar(100)
< row
: bytes
| RecordSystemStat
: uint64 1
: bytes
| IN
: bytes
| p_stat_name
: bytes
| varchar(100)
< row
: bytes
| RecordSystemStat
: uint64 2
: bytes
| IN
: bytes
| p_stat_value
: bytes
| decimal(10,2)
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_2f3a1b76
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_2f3a1b76
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
< row
: bytes
| evt_cleanup_old_stats
: bytes
| BEGIN
|     DELETE FROM system_stats
|     WHERE recorded_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
| END
: bytes
| RECURRING
: bytes
| 1
: bytes
| DAY
: bytes
| ENABLED
: bytes
| NOT PRESERVE
> exec
| USE pa_mm_2f3a1b76_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_2f3a1b76_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| system_stats
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_2f3a1b76_to`.`system_stats`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| system_stats
: bytes
| CREATE TABLE `system_stats` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `stat_name` varchar(100) NOT NULL,
|   `stat_value` decimal(10,2) DEFAULT NULL,
|   `recorded_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_2f3a1b76_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
< row
: bytes
| FUNCTION
: bytes
| GetAverageStat
: bytes
| decimal(10,2)
: bytes
| BEGIN
|     DECLARE avg_value DECIMAL(10, 2);
|
|     SELECT AVG(stat_value) INTO avg_value
|     FROM system_stats
|     WHERE stat_name = p_stat_name
|     AND recorded_at >= DATE_SUB(NOW(), INTERVAL 7 DAY);
|
|     RETURN COALESCE(avg_value, 0);
| END
: bytes
| YES
: bytes
| READS SQL DATA
: bytes
| DEFINER
: bytes
|
< row
: bytes
| PROCEDURE
: bytes
| RecordSystemStat
: null
: bytes
| BEGIN
|     DECLARE EXIT HANDLER FOR SQLEXCEPTION
|     BEGIN
|         ROLLBACK;
|         RESIGNAL;
|     END;
|
|     START TRANSACTION;
|
|     INSERT INTO system_stats (stat_name, stat_value)
|     VALUES (p_stat_name, p_stat_value);
|
|     COMMIT;
| END
: bytes
| NO
: bytes
| CONTAINS SQL
: bytes
| DEFINER
: bytes
|
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_2f3a1b76_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< row
: bytes
| GetAverageStat
: uint64 0
: null
: null
: bytes
| decimal(10,2)
< row
: bytes
| GetAverageStat
: uint64 1
: bytes
| IN
: bytes
| p_stat_name
: bytes
| varchar(100)
< row
: bytes
| RecordSystemStat
: uint64 1
: bytes
| IN
: bytes
| p_stat_name
: bytes
| varchar(100)
< row
: bytes
| RecordSystemStat
: uint64 2
: bytes
| IN
: bytes
| p_stat_value
: bytes
| decimal(10,2)
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_2f3a1b76_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_2f3a1b76_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
< row
: bytes
| evt_cleanup_old_stats
: bytes
| BEGIN
|     DELETE FROM system_stats
|     WHERE recorded_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
| END
: bytes
| RECURRING
: bytes
| 1
: bytes
| DAY
: bytes
| ENABLED
: bytes
| NOT PRESERVE
> exec
| DROP DATABASE IF EXISTS pa_mm_4a125487
< ok 0 0
> exec
| CREATE DATABASE pa_mm_4a125487
< ok 1 0
> exec
| USE pa_mm_4a125487
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| --
| -- Table structure for `some_table`
| --
| CREATE TABLE `some_table` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `some_table_chk_1` CHECK (`a` in (1,2,3))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| --
| -- Table structure for `t`
| --
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` char(255) DEFAULT (NULL),
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| --
| -- Table structure for `t1`
| --
| CREATE TABLE `t1` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   `b` int DEFAULT (NULL),
|   `c` int DEFAULT (NULL),
|   PRIMARY KEY (`id`),
|   KEY `b` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| --
| -- Table structure for `test_table`
| --
| CREATE TABLE `test_table` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) DEFAULT (NULL),
|   `email` varchar(100) DEFAULT (NULL),
|   `status` varchar(20) DEFAULT ('active'),
|   `age` int DEFAULT (NULL),
|   `score` decimal(5,2) DEFAULT (NULL),
|   PRIMARY KEY (`id`),
|   KEY `idx_name` (`name`),
|   KEY `idx_status_hash` (`status`),
|   UNIQUE KEY `uk_email` (`email`),
|   CONSTRAINT `chk_age` CHECK (`age` >= 18),
|   CONSTRAINT `chk_status` CHECK (`status` in ('active','inactive','pending'))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE `pa_mm_4a125487`.`test_table` MODIFY COLUMN `status` varchar(20) DEFAULT 'active'
< ok 0 0
> exec
| DROP DATABASE IF EXISTS pa_mm_4a125487_to
< ok 0 0
> exec
| CREATE DATABASE pa_mm_4a125487_to
< ok 1 0
> exec
| USE pa_mm_4a125487_to
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| -- Test that identical schemas don't generate unnecessary migrations
| -- This tests the fixes for primary key, index type, check constraints, and DEFAULT NULL
|
| CREATE TABLE test_table (
|     id INT NOT NULL AUTO_INCREMENT,
|     name VARCHAR(100) DEFAULT NULL,
|     email VARCHAR(100),
|     status VARCHAR(20) DEFAULT 'active',
|     age INT DEFAULT NULL,
|     score DECIMAL(5,2),
|     PRIMARY KEY (id),
|     INDEX idx_name (name),
|     INDEX idx_status_hash (status) USING HASH,
|     UNIQUE KEY uk_email (email),
|     CONSTRAINT chk_age CHECK (age >= 18),
|     CONSTRAINT chk_status CHECK (status IN ('active', 'inactive', 'pending'))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| -- Table with various column types and defaults
| CREATE TABLE t (
|     id INT NOT NULL AUTO_INCREMENT,
|     name CHAR(255) DEFAULT NULL,
|     PRIMARY KEY (id)
| ) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t1 (
|     id INT NOT NULL,
|     a INT NOT NULL,
|     b INT DEFAULT NULL,
|     c INT,
|     PRIMARY KEY (id),
|     INDEX b (b) USING HASH
| ) ENGINE=InnoDB
< ok 0 0
> exec
| -- Table with check constraint that has spaces
| CREATE TABLE some_table (
|     id INT PRIMARY KEY,
|     a INT NOT NULL,
|     CONSTRAINT some_table_chk_1 CHECK (a IN (1, 2, 3))
| ) ENGINE=InnoDB
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| USE pa_mm_4a125487
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_4a125487
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| some_table
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| t1
: bytes
| BASE TABLE
< row
: bytes
| test_table
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_4a125487`.`some_table`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| some_table
: bytes
| CREATE TABLE `some_table` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `some_table_chk_1` CHECK ((`a` in (1,2,3)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_4a125487`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` char(255) DEFAULT (NULL),
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_4a125487`.`t1`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t1
: bytes
| CREATE TABLE `t1` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   `b` int DEFAULT (NULL),
|   `c` int DEFAULT (NULL),
|   PRIMARY KEY (`id`),
|   KEY `b` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_4a125487`.`test_table`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| test_table
: bytes
| CREATE TABLE `test_table` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) DEFAULT (NULL),
|   `email` varchar(100) DEFAULT (NULL),
|   `status` varchar(20) DEFAULT 'active',
|   `age` int DEFAULT (NULL),
|   `score` decimal(5,2) DEFAULT (NULL),
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`),
|   KEY `idx_name` (`name`),
|   KEY `idx_status_hash` (`status`),
|   CONSTRAINT `chk_age` CHECK ((`age` >= 18)),
|   CONSTRAINT `chk_status` CHECK ((`status` in (_utf8mb4'active',_utf8mb4'inactive',_utf8mb4'pending')))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_4a125487
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_4a125487
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_4a125487
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_4a125487
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
> exec
| USE pa_mm_4a125487_to
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| pa_mm_4a125487_to
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| some_table
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| t1
: bytes
| BASE TABLE
< row
: bytes
| test_table
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `pa_mm_4a125487_to`.`some_table`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| some_table
: bytes
| CREATE TABLE `some_table` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `some_table_chk_1` CHECK ((`a` in (1,2,3)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_4a125487_to`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` char(255) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_4a125487_to`.`t1`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t1
: bytes
| CREATE TABLE `t1` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   `b` int DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `b` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `pa_mm_4a125487_to`.`test_table`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| test_table
: bytes
| CREATE TABLE `test_table` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) DEFAULT NULL,
|   `email` varchar(100) DEFAULT NULL,
|   `status` varchar(20) DEFAULT 'active',
|   `age` int DEFAULT NULL,
|   `score` decimal(5,2) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`),
|   KEY `idx_name` (`name`),
|   KEY `idx_status_hash` (`status`),
|   CONSTRAINT `chk_age` CHECK ((`age` >= 18)),
|   CONSTRAINT `chk_status` CHECK ((`status` in (_utf8mb4'active',_utf8mb4'inactive',_utf8mb4'pending')))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME
: string
| pa_mm_4a125487_to
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
< column "ROUTINE_DEFINITION" "TEXT"
< column "IS_DETERMINISTIC" "VARCHAR"
< column "SQL_DATA_ACCESS" "ENUM"
< column "SECURITY_TYPE" "ENUM"
< column "ROUTINE_COMMENT" "TEXT"
> query
| SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION
: string
| pa_mm_4a125487_to
< result
< column "SPECIFIC_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "PARAMETER_MODE" "VARCHAR"
< column "PARAMETER_NAME" "VARCHAR"
< column "DTD_IDENTIFIER" "TEXT"
> query
| SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| pa_mm_4a125487_to
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_ORDER" "UNSIGNED INT"
< column "ACTION_STATEMENT" "TEXT"
> query
| SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME
: string
| pa_mm_4a125487_to
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_DEFINITION" "TEXT"
< column "EVENT_TYPE" "VARCHAR"
< column "INTERVAL_VALUE" "VARCHAR"
< column "INTERVAL_FIELD" "ENUM"
< column "STATUS" "ENUM"
< column "ON_COMPLETION" "VARCHAR"
