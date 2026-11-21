-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestDatabaseDefinitionCorpusRunsOnMySQL default#1
> exec
| DROP DATABASE IF EXISTS def_corpus_0
< ok 0 0
> exec
| CREATE DATABASE def_corpus_0
< ok 1 0
> exec
| USE def_corpus_0
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `posts` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `user_id` int NOT NULL,
|   `title` varchar(200) NOT NULL,
|   `content` text DEFAULT NULL,
|   `published_at` datetime DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_user_published` (`user_id`, `published_at`),
|   CONSTRAINT `posts_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE TABLE `users` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `username` varchar(50) NOT NULL,
|   `email` varchar(100) NOT NULL,
|   `age` int DEFAULT NULL,
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
|   `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `profile` json,
|   `is_active` tinyint(1) DEFAULT '1',
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `username` (`username`),
|   KEY `idx_email` (`email`),
|   KEY `idx_created_at` (`created_at`),
|   CONSTRAINT `users_chk_1` CHECK (`age` >= 18)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> exec
| DROP DATABASE IF EXISTS def_corpus_1
< ok 0 0
> exec
| CREATE DATABASE def_corpus_1
< ok 1 0
> exec
| USE def_corpus_1
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `inventory` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `product_id` int NOT NULL,
|   `warehouse` varchar(50) NOT NULL,
|   `quantity` int NOT NULL DEFAULT '0',
|   `last_updated` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_product_warehouse` (`product_id`, `warehouse`),
|   CONSTRAINT `inventory_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE TABLE `products` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `price` decimal(10,2) NOT NULL,
|   `tax_rate` decimal(5,2) NOT NULL DEFAULT '0.08',
|   `price_with_tax` decimal(10,2) GENERATED ALWAYS AS ((`price` * ((1 + `tax_rate`)))) STORED,
|   `description` text DEFAULT NULL,
|   `tags` json,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `idx_fulltext` (`name`, `description`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> exec
| DROP DATABASE IF EXISTS def_corpus_2
< ok 0 0
> exec
| CREATE DATABASE def_corpus_2
< ok 1 0
> exec
| USE def_corpus_2
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `order_history` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `order_id` int NOT NULL,
|   `old_status` varchar(20) DEFAULT NULL,
|   `new_status` varchar(20) DEFAULT NULL,
|   `changed_at` timestamp DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   KEY `order_id` (`order_id`),
|   CONSTRAINT `order_history_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE TABLE `orders` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `order_number` varchar(20) NOT NULL,
|   `customer_name` varchar(100) NOT NULL,
|   `total_amount` decimal(10,2) NOT NULL,
|   `status` varchar(20) DEFAULT 'pending',
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `order_number` (`order_number`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `pending_orders` AS select `orders`.`id` AS `id`,`orders`.`order_number` AS `order_number`,`orders`.`customer_name` AS `customer_name`,`orders`.`total_amount` AS `total_amount`,`orders`.`created_at` AS `created_at` from `orders` where (`orders`.`status` = 'pending') order by `orders`.`created_at` desc
< ok 0 0
> exec
| CREATE TRIGGER `order_status_change`
| AFTER UPDATE ON `orders` FOR EACH ROW
| BEGIN
|     IF OLD.status != NEW.status THEN
|         INSERT INTO order_history (order_id, old_status, new_status)
|         VALUES (NEW.id, OLD.status, NEW.status);
|     END IF;
| END
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> exec
| DROP DATABASE IF EXISTS def_corpus_3
< ok 0 0
> exec
| CREATE DATABASE def_corpus_3
< ok 1 0
> exec
| USE def_corpus_3
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `accounts` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `account_number` varchar(20) NOT NULL,
|   `balance` decimal(15,2) NOT NULL DEFAULT '0.00',
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `account_number` (`account_number`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE FUNCTION `calculate_interest`(principal DECIMAL(15,2), rate DECIMAL(5,4), years INT) RETURNS decimal(15,2) READS SQL DATA DETERMINISTIC
| BEGIN
|     RETURN principal * POW(1 + rate, years);
| END
< ok 0 0
> exec
| CREATE PROCEDURE `transfer_funds`(IN from_account VARCHAR(20), IN to_account VARCHAR(20), IN amount DECIMAL(15,2))
| BEGIN
|     DECLARE from_balance DECIMAL(15, 2);
|
|     START TRANSACTION;
|
|     SELECT balance INTO from_balance
|     FROM accounts
|     WHERE account_number = from_account
|     FOR UPDATE;
|
|     IF from_balance >= amount THEN
|         UPDATE accounts
|         SET balance = balance - amount
|         WHERE account_number = from_account;
|
|         UPDATE accounts
|         SET balance = balance + amount
|         WHERE account_number = to_account;
|
|         COMMIT;
|     ELSE
|         ROLLBACK;
|         SIGNAL SQLSTATE '45000'
|         SET MESSAGE_TEXT = 'Insufficient funds';
|     END IF;
| END
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> exec
| DROP DATABASE IF EXISTS def_corpus_4
< ok 0 0
> exec
| CREATE DATABASE def_corpus_4
< ok 1 0
> exec
| USE def_corpus_4
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `customer_regions` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `customer_name` varchar(100) NOT NULL,
|   `region` varchar(20) NOT NULL,
|   `sales_amount` decimal(10,2) DEFAULT NULL,
|   PRIMARY KEY (`id`, `region`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY LIST COLUMNS (region)
| (PARTITION `p_north` VALUES IN ('north','northeast','northwest') ENGINE=InnoDB,
|  PARTITION `p_south` VALUES IN ('south','southeast','southwest') ENGINE=InnoDB,
|  PARTITION `p_east` VALUES IN ('east') ENGINE=InnoDB,
|  PARTITION `p_west` VALUES IN ('west') ENGINE=InnoDB,
|  PARTITION `p_central` VALUES IN ('central') ENGINE=InnoDB) */
< ok 0 0
> exec
| CREATE TABLE `employees` (
|   `id` int NOT NULL,
|   `name` varchar(100) NOT NULL,
|   `department_id` int NOT NULL,
|   `hired_date` date DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| (PARTITION `p0` ENGINE=InnoDB,
|  PARTITION `p1` ENGINE=InnoDB,
|  PARTITION `p2` ENGINE=InnoDB,
|  PARTITION `p3` ENGINE=InnoDB) */
< ok 0 0
> exec
| CREATE TABLE `order_archive` (
|   `order_id` int NOT NULL,
|   `order_date` date NOT NULL,
|   `customer_id` int NOT NULL,
|   `status` varchar(20) NOT NULL,
|   `total_amount` decimal(10,2) DEFAULT NULL,
|   PRIMARY KEY (`order_id`, `order_date`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY RANGE COLUMNS (order_date)
| (PARTITION `p_2022_q1` VALUES LESS THAN ('2022-04-01') ENGINE=InnoDB,
|  PARTITION `p_2022_q2` VALUES LESS THAN ('2022-07-01') ENGINE=InnoDB,
|  PARTITION `p_2022_q3` VALUES LESS THAN ('2022-10-01') ENGINE=InnoDB,
|  PARTITION `p_2022_q4` VALUES LESS THAN ('2023-01-01') ENGINE=InnoDB,
|  PARTITION `p_2023_and_later` VALUES LESS THAN MAXVALUE ENGINE=InnoDB) */
< ok 0 0
> exec
| CREATE TABLE `sales` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `sale_date` date NOT NULL,
|   `product_id` int NOT NULL,
|   `quantity` int NOT NULL,
|   `amount` decimal(10,2) NOT NULL,
|   PRIMARY KEY (`id`, `sale_date`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`sale_date`))
| (PARTITION `p2022` VALUES LESS THAN (2023) ENGINE=InnoDB,
|  PARTITION `p2023` VALUES LESS THAN (2024) ENGINE=InnoDB,
|  PARTITION `p2024` VALUES LESS THAN (2025) ENGINE=InnoDB,
|  PARTITION `p_future` VALUES LESS THAN MAXVALUE ENGINE=InnoDB) */
< ok 0 0
> exec
| CREATE TABLE `user_sessions` (
|   `session_id` varchar(64) NOT NULL,
|   `user_id` int NOT NULL,
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
|   `last_activity` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (`session_id`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY (session_id)
| (PARTITION `p0` ENGINE=InnoDB,
|  PARTITION `p1` ENGINE=InnoDB,
|  PARTITION `p2` ENGINE=InnoDB,
|  PARTITION `p3` ENGINE=InnoDB,
|  PARTITION `p4` ENGINE=InnoDB,
|  PARTITION `p5` ENGINE=InnoDB,
|  PARTITION `p6` ENGINE=InnoDB,
|  PARTITION `p7` ENGINE=InnoDB) */
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> exec
| DROP DATABASE IF EXISTS def_corpus_5
< ok 0 0
> exec
| CREATE DATABASE def_corpus_5
< ok 1 0
> exec
| USE def_corpus_5
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `daily_stats` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `stat_date` date NOT NULL,
|   `total_orders` int DEFAULT '0',
|   `total_revenue` decimal(15,2) DEFAULT '0.00',
|   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `stat_date` (`stat_date`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` EVENT `update_daily_stats` ON SCHEDULE EVERY 1 DAY
| STARTS CURRENT_TIMESTAMP ON COMPLETION NOT PRESERVE ENABLE DO INSERT INTO daily_stats (stat_date, total_orders, total_revenue)
|     VALUES (CURDATE() - INTERVAL 1 DAY, 0, 0.00)
|     ON DUPLICATE KEY UPDATE
|         total_orders = VALUES(total_orders),
|         total_revenue = VALUES(total_revenue)
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> exec
| DROP DATABASE IF EXISTS def_corpus_6
< ok 0 0
> exec
| CREATE DATABASE def_corpus_6
< ok 1 0
> exec
| USE def_corpus_6
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `translations` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `language_code` varchar(5) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
|   `content_key` varchar(100) COLLATE utf8mb4_bin NOT NULL,
|   `translation` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_lang_key` (`language_code`, `content_key`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> exec
| DROP DATABASE IF EXISTS def_corpus_7
< ok 0 0
> exec
| CREATE DATABASE def_corpus_7
< ok 1 0
> exec
| USE def_corpus_7
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `transactions` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `type` varchar(10) DEFAULT NULL,
|   `amount` decimal(10,2) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `c1` CHECK ((`type` = 'A') and (`amount` is not null))
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE FUNCTION `lookup_token`(p_a BIGINT, p_b VARCHAR(36)) RETURNS char(36) CHARSET utf8mb4 NO SQL DETERMINISTIC
| BEGIN
|     RETURN p_b;
| END
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> exec
| DROP DATABASE IF EXISTS def_corpus_8
< ok 0 0
> exec
| CREATE DATABASE def_corpus_8
< ok 1 0
> exec
| USE def_corpus_8
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `department` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `parent_id` (`parent_id`),
|   CONSTRAINT `department_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `department` (`id`) ON DELETE SET NULL
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE TABLE `employee` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `department_id` int DEFAULT NULL,
|   `manager_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `department_id` (`department_id`),
|   KEY `manager_id` (`manager_id`),
|   CONSTRAINT `employee_ibfk_1` FOREIGN KEY (`department_id`) REFERENCES `department` (`id`) ON DELETE SET NULL,
|   CONSTRAINT `employee_ibfk_2` FOREIGN KEY (`manager_id`) REFERENCES `employee` (`id`) ON DELETE SET NULL
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE TABLE `project` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) NOT NULL,
|   `lead_id` int NOT NULL,
|   `department_id` int NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `lead_id` (`lead_id`),
|   KEY `department_id` (`department_id`),
|   CONSTRAINT `project_ibfk_1` FOREIGN KEY (`lead_id`) REFERENCES `employee` (`id`),
|   CONSTRAINT `project_ibfk_2` FOREIGN KEY (`department_id`) REFERENCES `department` (`id`)
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE TABLE `project_member` (
|   `project_id` int NOT NULL,
|   `employee_id` int NOT NULL,
|   `role` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`project_id`, `employee_id`),
|   KEY `employee_id` (`employee_id`),
|   CONSTRAINT `project_member_ibfk_1` FOREIGN KEY (`project_id`) REFERENCES `project` (`id`) ON DELETE CASCADE,
|   CONSTRAINT `project_member_ibfk_2` FOREIGN KEY (`employee_id`) REFERENCES `employee` (`id`) ON DELETE CASCADE
| ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
