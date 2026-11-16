CREATE TABLE accounts (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	email VARCHAR(190) NOT NULL,
	bio TEXT,
	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
	flags SET('alpha', 'beta') DEFAULT 'alpha',
	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
	raw_token BINARY(8) NOT NULL DEFAULT 'ab',
	raw_bytes VARBINARY(16) NOT NULL DEFAULT 0x0100fe,
	secret_code VARCHAR(32) INVISIBLE DEFAULT 'hidden',
	meta JSON,
	coords POINT NOT NULL SRID 4326,
	geom GEOMETRY SRID 0,
	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL,
	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED,
	external_id VARCHAR(36) DEFAULT (UUID()),
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
	note VARCHAR(100) COMMENT 'owner''s note \\ path',
	UNIQUE KEY uk_email (email),
	KEY idx_bio_prefix (bio(32)),
	KEY idx_created_desc (created_at DESC) COMMENT 'newest first',
	KEY idx_invisible (status) INVISIBLE,
	KEY idx_expr ((LOWER(status))),
	FULLTEXT KEY ft_bio (bio),
	SPATIAL INDEX sp_coords (coords),
	CONSTRAINT chk_balance CHECK (balance >= 0 AND status IN ('open', 'closed'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='customer''s accounts';

CREATE TABLE orders (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	account_id BIGINT NOT NULL,
	parent_order_id BIGINT DEFAULT NULL,
	total INT NOT NULL,
	KEY idx_account (account_id),
	KEY idx_parent (parent_order_id),
	CONSTRAINT fk_orders_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT,
	CONSTRAINT fk_orders_parent FOREIGN KEY (parent_order_id) REFERENCES orders (id) ON DELETE SET NULL ON UPDATE NO ACTION
) ENGINE=InnoDB;

CREATE TABLE metrics_range (
	ts INT NOT NULL,
	region INT NOT NULL,
	val BIGINT NOT NULL,
	PRIMARY KEY (ts, region)
) PARTITION BY RANGE (ts) SUBPARTITION BY HASH (region) SUBPARTITIONS 2 (
	PARTITION p0 VALUES LESS THAN (100),
	PARTITION p1 VALUES LESS THAN (MAXVALUE)
);

CREATE TABLE metrics_hash (
	id INT NOT NULL PRIMARY KEY,
	val INT NOT NULL
) PARTITION BY HASH (id) PARTITIONS 4;

CREATE VIEW v_open_accounts AS
	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open';

CREATE VIEW v_named_cols (`renamed_id`) AS SELECT id FROM accounts; -- not restorable: TestDumpGaps

CREATE VIEW v_order_totals AS
	SELECT a.acct_id, COUNT(o.id) AS order_count
	FROM v_open_accounts a
	LEFT JOIN orders o ON o.account_id = a.acct_id
	GROUP BY a.acct_id;

CREATE FUNCTION fn_format_tag(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT)
RETURNS VARCHAR(64) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci
DETERMINISTIC
NO SQL
COMMENT 'formats a tag'
BEGIN
	RETURN CONCAT(p_prefix, '-', p_id);
END;

CREATE PROCEDURE sp_touch_order(IN p_id BIGINT, OUT p_total INT)
COMMENT 'reads order total'
BEGIN
	SELECT total INTO p_total FROM orders WHERE id = p_id;
END;

CREATE TRIGGER trg_orders_bi
BEFORE INSERT ON orders
FOR EACH ROW
BEGIN
	IF NEW.total < 0 THEN
		SET NEW.total = 0;
	END IF;
END;

CREATE EVENT ev_daily_check
ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00'
DISABLE
DO
BEGIN
	DELETE FROM orders WHERE total < 0;
END;
