SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
--
-- Temporary view structure for `v_mv_emails`
--
CREATE VIEW `v_mv_emails` AS SELECT
  1 AS `email`;

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
-- Table structure for `accounts`
--
CREATE TABLE `accounts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `email` varchar(190) NOT NULL,
  `bio` text DEFAULT NULL,
  `status` enum('open','closed') NOT NULL DEFAULT 'open',
  `literal_ts` varchar(40) NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
  `balance` decimal(12,4) NOT NULL DEFAULT '0.0000',
  `email_lower` varchar(190) GENERATED ALWAYS AS (LOWER(`email`)) VIRTUAL,
  `display_tag` varchar(220) GENERATED ALWAYS AS (CONCAT(`email`,':','v1')) STORED,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `reviewed_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `note` varchar(100) DEFAULT NULL COMMENT 'owner\'s note',
  PRIMARY KEY (`id`),
  KEY `idx_bio_prefix` (`bio`(32)),
  UNIQUE KEY `uk_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='customer accounts';

--
-- Table structure for `orders`
--
CREATE TABLE `orders` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `account_id` bigint NOT NULL,
  `total` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_account` (`account_id`),
  CONSTRAINT `fk_orders_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Table structure for `p_auto`
--
CREATE TABLE `p_auto` (
  `a` int NOT NULL,
  PRIMARY KEY (`a`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Table structure for `p_hash`
--
CREATE TABLE `p_hash` (
  `a` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
/*!50100 PARTITION BY HASH (`a`)
(PARTITION `p0` ENGINE=InnoDB,
 PARTITION `p1` ENGINE=InnoDB,
 PARTITION `p2` ENGINE=InnoDB) */;

--
-- Table structure for `p_key`
--
CREATE TABLE `p_key` (
  `a` int NOT NULL,
  `b` int DEFAULT NULL,
  PRIMARY KEY (`a`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
/*!50100 PARTITION BY KEY ()
(PARTITION `p0` ENGINE=InnoDB,
 PARTITION `p1` ENGINE=InnoDB) */;

--
-- Table structure for `p_list`
--
CREATE TABLE `p_list` (
  `a` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
/*!50100 PARTITION BY LIST (`a`)
(PARTITION `p_low` VALUES IN (1,2) ENGINE=InnoDB,
 PARTITION `p_rest` VALUES IN (DEFAULT) ENGINE=InnoDB) */;

--
-- Table structure for `p_list_cols`
--
CREATE TABLE `p_list_cols` (
  `a` int DEFAULT NULL,
  `b` char(2) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
/*!50500 PARTITION BY LIST COLUMNS (`a`, `b`)
(PARTITION `p0` VALUES IN ((1,'a'),(2,'b')) ENGINE=InnoDB,
 PARTITION `p1` VALUES IN ((3,'c')) ENGINE=InnoDB) */;

--
-- Table structure for `p_range`
--
CREATE TABLE `p_range` (
  `id` int DEFAULT NULL,
  `sold` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
/*!50100 PARTITION BY RANGE (year(`sold`))
(PARTITION `p2020` VALUES LESS THAN (2021) ENGINE=InnoDB,
 PARTITION `p_max` VALUES LESS THAN MAXVALUE ENGINE=InnoDB) */;

--
-- Table structure for `p_range_cols`
--
CREATE TABLE `p_range_cols` (
  `a` int DEFAULT NULL,
  `b` varchar(10) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
/*!50500 PARTITION BY RANGE COLUMNS (`a`, `b`)
(PARTITION `p0` VALUES LESS THAN (10,'m') ENGINE=InnoDB,
 PARTITION `p1` VALUES LESS THAN (MAXVALUE,MAXVALUE) ENGINE=InnoDB) */;

--
-- Table structure for `p_sub`
--
CREATE TABLE `p_sub` (
  `a` int DEFAULT NULL,
  `b` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
/*!50100 PARTITION BY RANGE (`a`)
SUBPARTITION BY HASH (`b`)
(PARTITION `p0` VALUES LESS THAN (10) ENGINE=InnoDB
 (SUBPARTITION `p0sp0` ENGINE=InnoDB,
  SUBPARTITION `p0sp1` ENGINE=InnoDB),
 PARTITION `p1` VALUES LESS THAN (20) ENGINE=InnoDB
 (SUBPARTITION `p1sp0` ENGINE=InnoDB,
  SUBPARTITION `p1sp1` ENGINE=InnoDB)) */;

--
-- Table structure for `p_sub_named`
--
CREATE TABLE `p_sub_named` (
  `a` int DEFAULT NULL,
  `b` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
/*!50100 PARTITION BY LIST (`a`)
SUBPARTITION BY RANGE (`b`)
(PARTITION `p0` VALUES IN (1) ENGINE=InnoDB
 (SUBPARTITION `s0` VALUES LESS THAN (5) ENGINE=InnoDB,
  SUBPARTITION `s1` VALUES LESS THAN MAXVALUE ENGINE=InnoDB),
 PARTITION `p1` VALUES IN (2) ENGINE=InnoDB
 (SUBPARTITION `s2` VALUES LESS THAN MAXVALUE ENGINE=InnoDB)) */;

--
-- Table structure for `ext_feed`
--
CREATE EXTERNAL TABLE `ext_feed` (
  `c1` int GENERATED ALWAYS AS (metadata$filecol1),
  `c2` varchar(10) GENERATED ALWAYS AS (metadata$filecol2) COMMENT 'owner\'s note'
)
LOCATION = 'file:///tmp'
PATTERN = '.*[.]csv'
AUTO_REFRESH = IMMEDIATE
FORMAT (TYPE = 'CSV', ENCODING = 'utf8mb4');

--
-- Function structure for `fn_escaped`
--
SET character_set_client = utf8mb4;
SET character_set_results = utf8mb4;
SET collation_connection = utf8mb4_general_ci;
SET sql_mode = 'NO_BACKSLASH_ESCAPES,STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER';
DELIMITER ;;
CREATE FUNCTION `fn_escaped`(p VARCHAR(16)) RETURNS varchar(32)
    NO SQL
    DETERMINISTIC
    COMMENT 'it''s C:\\path'
RETURN CONCAT(p, '-escaped');;
DELIMITER ;

--
-- Function structure for `fn_format_tag`
--
SET character_set_client = utf8mb4;
SET character_set_results = utf8mb4;
SET collation_connection = utf8mb4_general_ci;
SET sql_mode = 'STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER';
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
-- Materialized view structure for `mv_accounts`
--
SET sql_mode = DEFAULT;
CREATE MATERIALIZED VIEW `mv_accounts` COMMENT 'owners\' materialized accounts' AS select `omni_sync_corpus`.`accounts`.`id` AS `id`,`omni_sync_corpus`.`accounts`.`email` AS `email` from `omni_sync_corpus`.`accounts`;

CREATE INDEX `idx_mv_email` ON `mv_accounts` (`email`);

--
-- View structure for `v_open_accounts`
--
CREATE OR REPLACE VIEW `v_open_accounts` AS select `omni_sync_corpus`.`accounts`.`id` AS `acct_id`,`omni_sync_corpus`.`accounts`.`email` AS `acct_email` from `omni_sync_corpus`.`accounts` where (`omni_sync_corpus`.`accounts`.`status` = 'open');

--
-- Materialized view structure for `mv_open_accounts`
--
SET sql_mode = DEFAULT;
CREATE MATERIALIZED VIEW `mv_open_accounts` AS select `omni_sync_corpus`.`v_open_accounts`.`acct_id` AS `acct_id` from `omni_sync_corpus`.`v_open_accounts`;

--
-- View structure for `v_mv_emails`
--
CREATE OR REPLACE VIEW `v_mv_emails` AS select `mv_accounts`.`email` AS `email` from `omni_sync_corpus`.`mv_accounts` `mv_accounts`;

--
-- View structure for `v_named_cols`
--
CREATE OR REPLACE VIEW `v_named_cols` AS select `omni_sync_corpus`.`accounts`.`id` AS `renamed_id` from `omni_sync_corpus`.`accounts`;

--
-- View structure for `v_named_cols_ref`
--
CREATE OR REPLACE VIEW `v_named_cols_ref` AS select `omni_sync_corpus`.`v_named_cols`.`renamed_id` AS `renamed_id` from `omni_sync_corpus`.`v_named_cols`;

--
-- Procedure structure for `sp_touch_order`
--
SET character_set_client = utf8mb4;
SET character_set_results = utf8mb4;
SET collation_connection = utf8mb4_general_ci;
SET sql_mode = 'STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER';
DELIMITER ;;
CREATE PROCEDURE `sp_touch_order`(IN p_id BIGINT, OUT p_total INT)
    READS SQL DATA
    COMMENT 'reads an order''s total from C:\\orders'
BEGIN
	SELECT total INTO p_total FROM orders WHERE id = p_id;
END;;
DELIMITER ;

--
-- Event structure for `ev_daily_check`
--
SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';
SET time_zone = DEFAULT;
DELIMITER ;;
CREATE EVENT `ev_daily_check` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO DELETE FROM orders WHERE total < 0;;
DELIMITER ;

--
-- Event structure for `ev_once`
--
SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';
SET time_zone = DEFAULT;
DELIMITER ;;
CREATE EVENT `ev_once` ON SCHEDULE AT '2030-06-01 12:00:00' ON COMPLETION PRESERVE COMMENT 'owner''s \\new once' DO DELETE FROM orders;;
DELIMITER ;

--
-- Trigger structure for `trg_orders_bi`
--
SET character_set_client = utf8mb4;
SET character_set_results = utf8mb4;
SET collation_connection = utf8mb4_general_ci;
SET sql_mode = 'STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER';
DELIMITER ;;
CREATE TRIGGER `trg_orders_bi`
BEFORE INSERT ON `orders` FOR EACH ROW
BEGIN IF NEW.total < 0 THEN SET NEW.total = 0; END IF; END;;
DELIMITER ;

SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
