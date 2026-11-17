CREATE TABLE accounts (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	email VARCHAR(190) NOT NULL,
	bio TEXT,
	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
	secret_code VARCHAR(32) INVISIBLE DEFAULT 'hidden',
	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL,
	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
	note VARCHAR(100) COMMENT 'owner''s note',
	UNIQUE KEY uk_email (email),
	KEY idx_bio_prefix (bio(32)),
	CONSTRAINT chk_balance CHECK (balance >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='customer accounts';

CREATE TABLE orders (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	account_id BIGINT NOT NULL,
	total INT NOT NULL,
	KEY idx_account (account_id),
	CONSTRAINT fk_orders_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE history (
	id INT PRIMARY KEY,
	amount INT,
	account_id BIGINT,
	KEY idx_amount (amount),
	CONSTRAINT fk_history_account FOREIGN KEY (account_id) REFERENCES accounts (id),
	CONSTRAINT chk_amount CHECK (amount > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING COMMENT='versioned';

CREATE TABLE explicit_period (
	id INT PRIMARY KEY,
	valid_from TIMESTAMP(6) GENERATED ALWAYS AS ROW START,
	valid_to TIMESTAMP(6) GENERATED ALWAYS AS ROW END,
	PERIOD FOR SYSTEM_TIME (valid_from, valid_to)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING;

CREATE TABLE metrics_range (
	ts INT NOT NULL,
	val BIGINT NOT NULL,
	PRIMARY KEY (ts)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY RANGE (ts) (
	PARTITION p_low VALUES LESS THAN (100),
	PARTITION p_high VALUES LESS THAN MAXVALUE
);

CREATE TABLE metrics_hash (
	id INT NOT NULL PRIMARY KEY,
	val INT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY HASH (id) PARTITIONS 4;

CREATE SEQUENCE seq_default;

CREATE SEQUENCE seq_custom START WITH 100 INCREMENT BY 5 MINVALUE 10 MAXVALUE 1000 CACHE 20 CYCLE COMMENT='sequence comment';

CREATE VIEW v_open_accounts AS
	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open';

CREATE VIEW v_named_cols (`renamed_id`) AS SELECT id FROM accounts;

CREATE VIEW v_named_cols_ref AS
	SELECT renamed_id FROM v_named_cols;

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
