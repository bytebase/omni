CREATE TABLE accounts (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	email VARCHAR(190) NOT NULL,
	bio TEXT,
	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
	literal_ts VARCHAR(40) NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL,
	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
	reviewed_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
	note VARCHAR(100) COMMENT 'owner''s note',
	UNIQUE KEY uk_email (email),
	KEY idx_bio_prefix (bio(32))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='customer accounts';

CREATE TABLE orders (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	account_id BIGINT NOT NULL,
	total INT NOT NULL,
	KEY idx_account (account_id),
	CONSTRAINT fk_orders_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE EXTERNAL TABLE ext_feed (c1 INT, c2 VARCHAR(10)) LOCATION = '/tmp' PATTERN = '.*[.]csv' FORMAT = (TYPE = 'CSV');

CREATE MATERIALIZED VIEW mv_accounts COMMENT 'owners'' materialized accounts' AS SELECT id, email FROM accounts;

CREATE INDEX idx_mv_email ON mv_accounts (email);

CREATE VIEW v_open_accounts AS
	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open';

CREATE VIEW v_named_cols (`renamed_id`) AS SELECT id FROM accounts;

CREATE VIEW v_named_cols_ref AS
	SELECT renamed_id FROM v_named_cols;

CREATE FUNCTION fn_format_tag(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT)
RETURNS VARCHAR(64) CHARSET utf8mb4 COLLATE utf8mb4_general_ci
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

SET SESSION sql_mode = 'NO_BACKSLASH_ESCAPES,STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER'; CREATE FUNCTION fn_escaped(p VARCHAR(16)) RETURNS VARCHAR(32) CHARSET utf8mb4 DETERMINISTIC NO SQL RETURN CONCAT(p, '-escaped'); SET SESSION sql_mode = DEFAULT;

CREATE TRIGGER trg_orders_bi BEFORE INSERT ON orders FOR EACH ROW BEGIN IF NEW.total < 0 THEN SET NEW.total = 0; END IF; END;

CREATE EVENT ev_daily_check ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DISABLE DO DELETE FROM orders WHERE total < 0;

CREATE EVENT ev_once ON SCHEDULE AT '2030-06-01 12:00:00' ON COMPLETION PRESERVE ENABLE COMMENT 'owner''s \\new once' DO DELETE FROM orders;

CREATE VIEW v_mv_emails AS SELECT email FROM mv_accounts;

CREATE MATERIALIZED VIEW mv_open_accounts AS SELECT acct_id FROM v_open_accounts;
