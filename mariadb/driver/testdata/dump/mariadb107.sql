SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@@EXPLICIT_DEFAULTS_FOR_TIMESTAMP, EXPLICIT_DEFAULTS_FOR_TIMESTAMP=1;
--
-- Temporary view structure for `v_named_cols`
--
CREATE VIEW `v_named_cols` AS SELECT
  1 AS `renamed_id`;

--
-- Temporary view structure for `v_named_cols_ref`
--
CREATE VIEW `v_named_cols_ref` AS SELECT
  1 AS `renamed_id`;

--
-- Temporary view structure for `v_open_accounts`
--
CREATE VIEW `v_open_accounts` AS SELECT
  1 AS `acct_id`,
  1 AS `acct_email`;

--
-- Sequence structure for `seq_custom`
--
CREATE SEQUENCE `seq_custom` start with 100 minvalue 10 maxvalue 1000 increment by 5 cache 20 cycle COMMENT='sequence comment';

--
-- Sequence structure for `seq_default`
--
CREATE SEQUENCE `seq_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle;

--
-- Table structure for `accounts`
--
CREATE TABLE `accounts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `email` varchar(190) NOT NULL,
  `bio` text DEFAULT NULL,
  `status` enum('open','closed') NOT NULL DEFAULT 'open',
  `balance` decimal(12,4) NOT NULL DEFAULT 0.0000,
  `secret_code` varchar(32) DEFAULT 'hidden' INVISIBLE,
  `email_lower` varchar(190) GENERATED ALWAYS AS (lcase(`email`)) VIRTUAL,
  `display_tag` varchar(220) GENERATED ALWAYS AS (concat(`email`,':','v1')) STORED,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `note` varchar(100) DEFAULT NULL COMMENT 'owner\'s note',
  PRIMARY KEY (`id`),
  KEY `idx_bio_prefix` (`bio`(32)),
  UNIQUE KEY `uk_email` (`email`),
  CONSTRAINT `chk_balance` CHECK (`balance` >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='customer accounts';

--
-- Table structure for `explicit_period`
--
CREATE TABLE `explicit_period` (
  `id` int NOT NULL,
  `valid_from` timestamp(6) GENERATED ALWAYS AS ROW START NOT NULL,
  `valid_to` timestamp(6) GENERATED ALWAYS AS ROW END NOT NULL,
  PERIOD FOR SYSTEM_TIME (`valid_from`, `valid_to`),
  PRIMARY KEY (`id`, `valid_to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING;

--
-- Table structure for `history`
--
CREATE TABLE `history` (
  `id` int NOT NULL,
  `amount` int DEFAULT NULL,
  `account_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_history_account` (`account_id`),
  KEY `idx_amount` (`amount`),
  CONSTRAINT `fk_history_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`),
  CONSTRAINT `chk_amount` CHECK (`amount` > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='versioned' WITH SYSTEM VERSIONING;

--
-- Table structure for `metrics_hash`
--
CREATE TABLE `metrics_hash` (
  `id` int NOT NULL,
  `val` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY HASH (`id`)
PARTITIONS 4;

--
-- Table structure for `metrics_range`
--
CREATE TABLE `metrics_range` (
  `ts` int NOT NULL,
  `val` bigint NOT NULL,
  PRIMARY KEY (`ts`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY RANGE (`ts`)
(PARTITION `p_low` VALUES LESS THAN (100) ENGINE=InnoDB COMMENT = 'first',
 PARTITION `p_high` VALUES LESS THAN MAXVALUE ENGINE=InnoDB COMMENT = 'last');

--
-- Table structure for `orders`
--
CREATE TABLE `orders` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `account_id` bigint NOT NULL,
  `total` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_account` (`account_id`),
  CONSTRAINT `fk_orders_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET @saved_cs_client = @@character_set_client, @saved_cs_results = @@character_set_results, @saved_col_connection = @@collation_connection, @saved_sql_mode = @@sql_mode;
--
-- Function structure for `fn_format_tag`
--
SET character_set_client = utf8mb4;
SET character_set_results = utf8mb4;
SET collation_connection = utf8mb4_general_ci;
SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';
DELIMITER ;;
CREATE FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64)
    NO SQL
    DETERMINISTIC
    COMMENT 'formats a tag'
BEGIN
	RETURN CONCAT(p_prefix, '-', p_id);
END;;
DELIMITER ;

SET character_set_client = @saved_cs_client, character_set_results = @saved_cs_results, collation_connection = @saved_col_connection, sql_mode = @saved_sql_mode;
--
-- View structure for `v_named_cols`
--
CREATE OR REPLACE VIEW `v_named_cols` AS select `accounts`.`id` AS `renamed_id` from `accounts`;

--
-- View structure for `v_named_cols_ref`
--
CREATE OR REPLACE VIEW `v_named_cols_ref` AS select `v_named_cols`.`renamed_id` AS `renamed_id` from `v_named_cols`;

--
-- View structure for `v_open_accounts`
--
CREATE OR REPLACE VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where `accounts`.`status` = 'open';

--
-- Procedure structure for `sp_touch_order`
--
SET character_set_client = utf8mb4;
SET character_set_results = utf8mb4;
SET collation_connection = utf8mb4_general_ci;
SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';
DELIMITER ;;
CREATE PROCEDURE `sp_touch_order`(IN p_id BIGINT, OUT p_total INT)
    COMMENT 'reads order total'
BEGIN
	SELECT total INTO p_total FROM orders WHERE id = p_id;
END;;
DELIMITER ;

--
-- Event structure for `ev_daily_check`
--
SET character_set_client = utf8mb4;
SET character_set_results = utf8mb4;
SET collation_connection = utf8mb4_general_ci;
SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';
SET time_zone = 'SYSTEM';
DELIMITER ;;
CREATE DEFINER=`root`@`%` EVENT `ev_daily_check` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO BEGIN
	DELETE FROM orders WHERE total < 0;
END;;
DELIMITER ;

--
-- Trigger structure for `trg_orders_bi`
--
SET character_set_client = utf8mb4;
SET character_set_results = utf8mb4;
SET collation_connection = utf8mb4_general_ci;
SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';
DELIMITER ;;
CREATE TRIGGER `trg_orders_bi`
BEFORE INSERT ON `orders` FOR EACH ROW
BEGIN
	IF NEW.total < 0 THEN
		SET NEW.total = 0;
	END IF;
END;;
DELIMITER ;

SET EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
