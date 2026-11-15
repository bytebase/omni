CREATE TABLE accounts (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	email VARCHAR(190) NOT NULL,
	bio TEXT,
	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL,
	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
	note VARCHAR(100) COMMENT 'owner''s note',
	UNIQUE KEY uk_email (email),
	KEY idx_bio_prefix (bio(32)),
	KEY idx_expr ((LOWER(status)))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='customer accounts';

CREATE TABLE sharded_events (
	id BIGINT PRIMARY KEY AUTO_RANDOM(5, 64),
	account_id BIGINT NOT NULL,
	payload VARCHAR(128),
	KEY idx_account (account_id),
	CONSTRAINT fk_events_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE metrics_range (
	ts INT NOT NULL,
	val BIGINT NOT NULL,
	PRIMARY KEY (ts)
) PARTITION BY RANGE (ts) (
	PARTITION p0 VALUES LESS THAN (100),
	PARTITION p1 VALUES LESS THAN MAXVALUE
);

CREATE SEQUENCE seq_orders START WITH 10 INCREMENT BY 2 MINVALUE 1 MAXVALUE 10000 CACHE 50 CYCLE COMMENT='order seq'; -- not restorable: TestDumpGaps

CREATE VIEW v_open_accounts AS
	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open';
