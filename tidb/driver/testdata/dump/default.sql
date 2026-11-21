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
-- Sequence structure for `seq_orders`
--
CREATE SEQUENCE `seq_orders` start with 10 minvalue 1 maxvalue 10000 increment by 2 cache 50 cycle ENGINE=InnoDB COMMENT='order seq';

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
  `email_lower` varchar(190) GENERATED ALWAYS AS (lower(`email`)) VIRTUAL,
  `display_tag` varchar(220) GENERATED ALWAYS AS (concat(`email`, _utf8mb4':', _utf8mb4'v1')) STORED,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `note` varchar(100) DEFAULT NULL COMMENT 'owner\'s note',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `idx_bio_prefix` (`bio`(32)),
  KEY `idx_expr` ((lower(`status`))),
  UNIQUE KEY `uk_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='customer accounts';

--
-- Table structure for `metrics_range`
--
CREATE TABLE `metrics_range` (
  `ts` int NOT NULL,
  `val` bigint NOT NULL,
  PRIMARY KEY (`ts`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
/*!50100 PARTITION BY RANGE (`ts`)
(PARTITION p_low VALUES LESS THAN (100) ENGINE=InnoDB COMMENT 'first',
 PARTITION p_high VALUES LESS THAN MAXVALUE ENGINE=InnoDB COMMENT 'last') */;

--
-- Table structure for `metrics_shard`
--
CREATE TABLE `metrics_shard` (
  `ts` int NOT NULL,
  `val` bigint NOT NULL,
  PRIMARY KEY (`ts`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
/*!50100 PARTITION BY RANGE (`ts`)
(PARTITION p10 VALUES LESS THAN (10) ENGINE=InnoDB COMMENT 'first tenth',
 PARTITION p2 VALUES LESS THAN (20) ENGINE=InnoDB,
 PARTITION p1 VALUES LESS THAN (30) ENGINE=InnoDB) */;

--
-- Table structure for `sharded_events`
--
CREATE TABLE `sharded_events` (
  `id` bigint NOT NULL /*T![auto_rand] AUTO_RANDOM(5) */,
  `account_id` bigint NOT NULL,
  `payload` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `idx_account` (`account_id`),
  CONSTRAINT `fk_events_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

--
-- View structure for `v_named_cols`
--
DROP VIEW IF EXISTS `v_named_cols`;
CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `v_named_cols` (`renamed_id`) AS SELECT `id` AS `id` FROM `omni_sync_corpus`.`accounts`;

--
-- View structure for `v_named_cols_ref`
--
DROP VIEW IF EXISTS `v_named_cols_ref`;
CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `v_named_cols_ref` AS SELECT `renamed_id` AS `renamed_id` FROM `omni_sync_corpus`.`v_named_cols`;

--
-- View structure for `v_open_accounts`
--
DROP VIEW IF EXISTS `v_open_accounts`;
CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `v_open_accounts` AS SELECT `id` AS `acct_id`,`email` AS `acct_email` FROM `omni_sync_corpus`.`accounts` WHERE `status`=_UTF8MB4'open';

SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
