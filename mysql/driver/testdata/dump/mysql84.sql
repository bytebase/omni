SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
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
-- Temporary view structure for `v_order_totals`
--
CREATE VIEW `v_order_totals` AS SELECT
  1 AS `acct_id`,
  1 AS `order_count`;

--
-- Table structure for `accounts`
--
CREATE TABLE `accounts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `email` varchar(190) NOT NULL,
  `bio` text DEFAULT NULL,
  `status` enum('open','closed') NOT NULL DEFAULT 'open',
  `flags` set('alpha','beta') DEFAULT 'alpha',
  `balance` decimal(12,4) NOT NULL DEFAULT '0.0000',
  `raw_token` binary(8) NOT NULL DEFAULT 'ab',
  `raw_bytes` varbinary(16) NOT NULL DEFAULT 0x01,
  `secret_code` varchar(32) DEFAULT 'hidden' /*!80023 INVISIBLE */,
  `meta` json,
  `coords` point NOT NULL /*!80003 SRID 4326 */,
  `geom` geometry /*!80003 SRID 0 */,
  `email_lower` varchar(190) GENERATED ALWAYS AS (lower(`email`)) VIRTUAL,
  `display_tag` varchar(220) GENERATED ALWAYS AS (concat(`email`,_utf8mb4':',_utf8mb4'v1')) STORED,
  `external_id` varchar(36) DEFAULT (uuid()),
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `note` varchar(100) DEFAULT NULL COMMENT 'owner\'s note \\ path',
  PRIMARY KEY (`id`),
  FULLTEXT KEY `ft_bio` (`bio`),
  KEY `idx_bio_prefix` (`bio`(32)),
  KEY `idx_created_desc` (`created_at` DESC) COMMENT 'newest first',
  KEY `idx_expr` ((lower(`status`))),
  KEY `idx_invisible` (`status`) /*!80000 INVISIBLE */,
  SPATIAL KEY `sp_coords` (`coords`),
  UNIQUE KEY `uk_email` (`email`),
  CONSTRAINT `chk_balance` CHECK ((`balance` >= 0) and (`status` in (_utf8mb4'open',_utf8mb4'closed')))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='customer\'s accounts';

--
-- Table structure for `metrics_hash`
--
CREATE TABLE `metrics_hash` (
  `id` int NOT NULL,
  `val` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
/*!50100 PARTITION BY HASH (`id`)
PARTITIONS 4 */;

--
-- Table structure for `metrics_range`
--
CREATE TABLE `metrics_range` (
  `ts` int NOT NULL,
  `region` int NOT NULL,
  `val` bigint NOT NULL,
  PRIMARY KEY (`ts`, `region`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
/*!50100 PARTITION BY RANGE (`ts`)
SUBPARTITION BY HASH (`region`)
SUBPARTITIONS 2
(PARTITION `p0` VALUES LESS THAN (100) ENGINE=InnoDB COMMENT = 'first',
 PARTITION `p1` VALUES LESS THAN MAXVALUE ENGINE=InnoDB COMMENT = 'last') */;

--
-- Table structure for `orders`
--
CREATE TABLE `orders` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `account_id` bigint NOT NULL,
  `parent_order_id` bigint DEFAULT NULL,
  `total` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_account` (`account_id`),
  KEY `idx_parent` (`parent_order_id`),
  CONSTRAINT `fk_orders_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `fk_orders_parent` FOREIGN KEY (`parent_order_id`) REFERENCES `orders` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Function structure for `fn_format_tag`
--
SET character_set_client = utf8mb4;
SET character_set_results = utf8mb4;
SET collation_connection = utf8mb4_general_ci;
SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';
DELIMITER ;;
CREATE FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64)
    NO SQL
    DETERMINISTIC
    COMMENT 'formats a tag'
BEGIN
	RETURN CONCAT(p_prefix, '-', p_id);
END;;
DELIMITER ;

--
-- View structure for `v_named_cols`
--
CREATE OR REPLACE VIEW `v_named_cols` (`renamed_id`) AS select `accounts`.`id` AS `id` from `accounts`;

--
-- View structure for `v_named_cols_ref`
--
CREATE OR REPLACE VIEW `v_named_cols_ref` AS select `v_named_cols`.`renamed_id` AS `renamed_id` from `v_named_cols`;

--
-- View structure for `v_open_accounts`
--
CREATE OR REPLACE VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where (`accounts`.`status` = 'open');

--
-- View structure for `v_order_totals`
--
CREATE OR REPLACE VIEW `v_order_totals` AS select `a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`v_open_accounts` `a` left join `orders` `o` on((`o`.`account_id` = `a`.`acct_id`))) group by `a`.`acct_id`;

--
-- Procedure structure for `sp_touch_order`
--
SET character_set_client = utf8mb4;
SET character_set_results = utf8mb4;
SET collation_connection = utf8mb4_general_ci;
SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';
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
SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';
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
SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';
DELIMITER ;;
CREATE TRIGGER `trg_orders_bi`
BEFORE INSERT ON `orders` FOR EACH ROW
BEGIN
	IF NEW.total < 0 THEN
		SET NEW.total = 0;
	END IF;
END;;
DELIMITER ;

SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
