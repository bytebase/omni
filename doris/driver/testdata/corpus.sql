CREATE TABLE accounts (
	id BIGINT NOT NULL,
	email VARCHAR(190) NOT NULL COMMENT 'login email',
	bio TEXT NULL,
	status VARCHAR(16) NOT NULL DEFAULT 'open',
	balance DECIMAL(12, 4) NOT NULL DEFAULT '0.0000',
	flag TINYINT NOT NULL DEFAULT '1',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME(6) NULL DEFAULT CURRENT_TIMESTAMP(6),
	note VARCHAR(100) NULL COMMENT 'owner''s note',
	dt DATE NULL,
	payload JSON NULL
) ENGINE=OLAP
DUPLICATE KEY(id, email)
COMMENT 'customer accounts'
DISTRIBUTED BY HASH(id) BUCKETS 1
PROPERTIES("replication_num" = "1")

CREATE TABLE orders (
	id BIGINT NOT NULL,
	account_id BIGINT NOT NULL,
	total INT NULL DEFAULT '0',
	created_on DATE NULL
) ENGINE=OLAP
DUPLICATE KEY(id)
DISTRIBUTED BY HASH(id) BUCKETS 1
PROPERTIES("replication_num" = "1")

CREATE TABLE metrics (
	ts DATETIME NOT NULL,
	val BIGINT NULL
) ENGINE=OLAP
DUPLICATE KEY(ts)
PARTITION BY RANGE(ts) (
	PARTITION p2025 VALUES [('2025-01-01 00:00:00'), ('2026-01-01 00:00:00')),
	PARTITION p2026 VALUES [('2026-01-01 00:00:00'), ('2027-01-01 00:00:00'))
)
DISTRIBUTED BY HASH(ts) BUCKETS 1
PROPERTIES("replication_num" = "1")

CREATE VIEW v_open_accounts AS
	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open'

CREATE MATERIALIZED VIEW mv_account_status
DISTRIBUTED BY HASH(id) BUCKETS 1
PROPERTIES("replication_num" = "1")
AS SELECT id, status FROM accounts

CREATE TABLE events (
	id BIGINT NOT NULL,
	name VARCHAR(64) NOT NULL COMMENT 'event name',
	payload TEXT NULL,
	slug VARCHAR(96) AS (concat(name, '-evt')),
	INDEX idx_name (name) USING INVERTED,
	INDEX idx_payload (payload) USING INVERTED PROPERTIES('parser' = 'english')
) ENGINE=OLAP
DUPLICATE KEY(id)
DISTRIBUTED BY RANDOM BUCKETS 3
PROPERTIES("replication_num" = "1", "disable_auto_compaction" = "true")

CREATE TABLE daily_rollup (
	day DATE NOT NULL,
	region VARCHAR(16) NOT NULL,
	amount DECIMAL(12, 2) SUM NULL,
	peak INT MAX NULL,
	label VARCHAR(32) REPLACE NULL,
	visitors HLL HLL_UNION NOT NULL,
	seen BITMAP BITMAP_UNION NOT NULL
) ENGINE=OLAP
AGGREGATE KEY(day, region)
DISTRIBUTED BY HASH(day) BUCKETS 2
PROPERTIES("replication_num" = "1")

CREATE TABLE tickets (
	id BIGINT NOT NULL AUTO_INCREMENT(1000),
	title VARCHAR(100) NULL,
	slug VARCHAR(120) AS (concat(title, '-t')),
	status VARCHAR(16) NOT NULL DEFAULT 'open'
) ENGINE=OLAP
UNIQUE KEY(id)
DISTRIBUTED BY HASH(id) BUCKETS 2
PROPERTIES("replication_num" = "1")

CREATE TABLE samples (
	id BIGINT NOT NULL,
	value INT NULL
) ENGINE=OLAP
UNIQUE KEY(id)
DISTRIBUTED BY HASH(id) BUCKETS AUTO
PROPERTIES("replication_num" = "1")

CREATE PROCEDURE bump_total(IN p_id BIGINT, IN delta INT) COMMENT 'add to an order total' BEGIN UPDATE samples SET value = value + delta WHERE id = p_id; END
