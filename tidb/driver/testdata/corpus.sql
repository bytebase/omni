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
	PARTITION p_low VALUES LESS THAN (100) COMMENT 'first',
	PARTITION p_high VALUES LESS THAN MAXVALUE COMMENT 'last'
);

-- A partition's name order is not its ordinal order: a renderer that sorted the
-- partitions by name would write p1 (30) before p10 (10), which the server
-- refuses with "VALUES LESS THAN value must be strictly increasing".
CREATE TABLE metrics_shard (
	ts INT NOT NULL,
	val BIGINT NOT NULL,
	PRIMARY KEY (ts)
) PARTITION BY RANGE (ts) (
	PARTITION p10 VALUES LESS THAN (10) COMMENT 'first tenth',
	PARTITION p2 VALUES LESS THAN (20),
	PARTITION p1 VALUES LESS THAN (30)
);

CREATE SEQUENCE seq_orders START WITH 10 INCREMENT BY 2 MINVALUE 1 MAXVALUE 10000 CACHE 50 CYCLE COMMENT='order seq';

CREATE VIEW v_open_accounts AS
	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open';

CREATE VIEW v_named_cols (`renamed_id`) AS SELECT id FROM accounts;

CREATE VIEW v_named_cols_ref AS
	SELECT renamed_id FROM v_named_cols;
