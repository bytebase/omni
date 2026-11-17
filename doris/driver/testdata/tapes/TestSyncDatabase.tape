-- pin: default dyrnq/doris:3.1.4 sha256:88f99819c0877567b77d78cc8367f0e2dd4b6c1602fe17c6f4b1e332d8e315dd 3.1.4-rc02

== TestSyncDatabase default#1
> exec
| CREATE TABLE accounts (
| 	id BIGINT NOT NULL,
| 	email VARCHAR(190) NOT NULL COMMENT 'login email',
| 	bio TEXT NULL,
| 	status VARCHAR(16) NOT NULL DEFAULT 'open',
| 	balance DECIMAL(12, 4) NOT NULL DEFAULT '0.0000',
| 	flag TINYINT NOT NULL DEFAULT '1',
| 	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at DATETIME(6) NULL DEFAULT CURRENT_TIMESTAMP(6),
| 	note VARCHAR(100) NULL COMMENT 'owner''s note',
| 	dt DATE NULL,
| 	payload JSON NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(id, email)
| COMMENT 'customer accounts'
| DISTRIBUTED BY HASH(id) BUCKETS 1
| PROPERTIES("replication_num" = "1")
< ok 0 0
> exec
| CREATE TABLE orders (
| 	id BIGINT NOT NULL,
| 	account_id BIGINT NOT NULL,
| 	total INT NULL DEFAULT '0',
| 	created_on DATE NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(id)
| DISTRIBUTED BY HASH(id) BUCKETS 1
| PROPERTIES("replication_num" = "1")
< ok 0 0
> exec
| CREATE TABLE metrics (
| 	ts DATETIME NOT NULL,
| 	val BIGINT NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(ts)
| PARTITION BY RANGE(ts) (
| 	PARTITION p2025 VALUES [('2025-01-01 00:00:00'), ('2026-01-01 00:00:00')),
| 	PARTITION p2026 VALUES [('2026-01-01 00:00:00'), ('2027-01-01 00:00:00'))
| )
| DISTRIBUTED BY HASH(ts) BUCKETS 1
| PROPERTIES("replication_num" = "1")
< ok 0 0
> exec
| CREATE VIEW v_open_accounts AS
| 	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open'
< ok 0 0
> exec
| CREATE MATERIALIZED VIEW mv_account_status
| DISTRIBUTED BY HASH(id) BUCKETS 1
| PROPERTIES("replication_num" = "1")
| AS SELECT id, status FROM accounts
< ok 0 0
> exec
| CREATE TABLE events (
| 	id BIGINT NOT NULL,
| 	name VARCHAR(64) NOT NULL COMMENT 'event name',
| 	payload TEXT NULL,
| 	slug VARCHAR(96) AS (concat(name, '-evt')),
| 	INDEX idx_name (name) USING INVERTED,
| 	INDEX idx_payload (payload) USING INVERTED PROPERTIES('parser' = 'english')
| ) ENGINE=OLAP
| DUPLICATE KEY(id)
| DISTRIBUTED BY RANDOM BUCKETS 3
| PROPERTIES("replication_num" = "1", "disable_auto_compaction" = "true")
< ok 0 0
> exec
| CREATE TABLE daily_rollup (
| 	day DATE NOT NULL,
| 	region VARCHAR(16) NOT NULL,
| 	amount DECIMAL(12, 2) SUM NULL,
| 	peak INT MAX NULL,
| 	label VARCHAR(32) REPLACE NULL,
| 	visitors HLL HLL_UNION NOT NULL,
| 	seen BITMAP BITMAP_UNION NOT NULL
| ) ENGINE=OLAP
| AGGREGATE KEY(day, region)
| DISTRIBUTED BY HASH(day) BUCKETS 2
| PROPERTIES("replication_num" = "1")
< ok 0 0
> exec
| CREATE TABLE tickets (
| 	id BIGINT NOT NULL AUTO_INCREMENT(1000),
| 	title VARCHAR(100) NULL,
| 	slug VARCHAR(120) AS (concat(title, '-t')),
| 	status VARCHAR(16) NOT NULL DEFAULT 'open'
| ) ENGINE=OLAP
| UNIQUE KEY(id)
| DISTRIBUTED BY HASH(id) BUCKETS 2
| PROPERTIES("replication_num" = "1")
< ok 0 0
> exec
| CREATE TABLE samples (
| 	id BIGINT NOT NULL,
| 	value INT NULL
| ) ENGINE=OLAP
| UNIQUE KEY(id)
| DISTRIBUTED BY HASH(id) BUCKETS AUTO
| PROPERTIES("replication_num" = "1")
< ok 0 0
> exec
| CREATE PROCEDURE bump_total(IN p_id BIGINT, IN delta INT) COMMENT 'add to an order total' BEGIN UPDATE samples SET value = value + delta WHERE id = p_id; END
< ok 0 0
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			IFNULL(COLUMN_NAME, ''),
| 			ORDINAL_POSITION,
| 			IFNULL(COLUMN_KEY, ''),
| 			COLUMN_DEFAULT,
| 			IS_NULLABLE,
| 			COLUMN_TYPE,
| 			IFNULL(CHARACTER_SET_NAME, ''),
| 			IFNULL(COLLATION_NAME, ''),
| 			COLUMN_COMMENT,
| 			EXTRA
| 		FROM information_schema.columns
| 		WHERE TABLE_SCHEMA = 'omni_sync_corpus'
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "CHAR"
< column "IFNULL(COLUMN_NAME, '')" "CHAR"
< column "ORDINAL_POSITION" "BIGINT"
< column "IFNULL(COLUMN_KEY, '')" "CHAR"
< column "COLUMN_DEFAULT" "CHAR"
< column "IS_NULLABLE" "CHAR"
< column "COLUMN_TYPE" "CHAR"
< column "IFNULL(CHARACTER_SET_NAME, '')" "CHAR"
< column "IFNULL(COLLATION_NAME, '')" "CHAR"
< column "COLUMN_COMMENT" "CHAR"
< column "EXTRA" "CHAR"
< row
: bytes
| accounts
: bytes
| id
: int64 1
: bytes
| DUP
: null
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| accounts
: bytes
| email
: int64 2
: bytes
| DUP
: null
: bytes
| NO
: bytes
| varchar(190)
: bytes
|
: bytes
|
: bytes
| login email
: bytes
|
< row
: bytes
| accounts
: bytes
| bio
: int64 3
: bytes
|
: null
: bytes
| YES
: bytes
| string
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| accounts
: bytes
| status
: int64 4
: bytes
|
: bytes
| open
: bytes
| NO
: bytes
| varchar(16)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| accounts
: bytes
| balance
: int64 5
: bytes
|
: bytes
| 0.0000
: bytes
| NO
: bytes
| decimalv3(12, 4)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| accounts
: bytes
| flag
: int64 6
: bytes
|
: bytes
| 1
: bytes
| NO
: bytes
| tinyint(4)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| accounts
: bytes
| created_at
: int64 7
: bytes
|
: bytes
| CURRENT_TIMESTAMP
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| accounts
: bytes
| updated_at
: int64 8
: bytes
|
: bytes
| CURRENT_TIMESTAMP(6)
: bytes
| YES
: bytes
| datetime(6)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| accounts
: bytes
| note
: int64 9
: bytes
|
: null
: bytes
| YES
: bytes
| varchar(100)
: bytes
|
: bytes
|
: bytes
| owner''s note
: bytes
|
< row
: bytes
| accounts
: bytes
| dt
: int64 10
: bytes
|
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| accounts
: bytes
| payload
: int64 11
: bytes
|
: null
: bytes
| YES
: bytes
| json
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| daily_rollup
: bytes
| day
: int64 1
: bytes
| AGG
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| daily_rollup
: bytes
| region
: int64 2
: bytes
| AGG
: null
: bytes
| NO
: bytes
| varchar(16)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| daily_rollup
: bytes
| amount
: int64 3
: bytes
|
: null
: bytes
| YES
: bytes
| decimalv3(12, 2)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| daily_rollup
: bytes
| peak
: int64 4
: bytes
|
: null
: bytes
| YES
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| daily_rollup
: bytes
| label
: int64 5
: bytes
|
: null
: bytes
| YES
: bytes
| varchar(32)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| daily_rollup
: bytes
| visitors
: int64 6
: bytes
|
: bytes
|q "\x00"
: bytes
| NO
: bytes
| hll
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| daily_rollup
: bytes
| seen
: int64 7
: bytes
|
: bytes
|q "\x00"
: bytes
| NO
: bytes
| bitmap
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| events
: bytes
| id
: int64 1
: bytes
| DUP
: null
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| events
: bytes
| name
: int64 2
: bytes
|
: null
: bytes
| NO
: bytes
| varchar(64)
: bytes
|
: bytes
|
: bytes
| event name
: bytes
|
< row
: bytes
| events
: bytes
| payload
: int64 3
: bytes
|
: null
: bytes
| YES
: bytes
| string
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| events
: bytes
| slug
: int64 4
: bytes
|
: null
: bytes
| YES
: bytes
| varchar(96)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| metrics
: bytes
| ts
: int64 1
: bytes
| DUP
: null
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| metrics
: bytes
| val
: int64 2
: bytes
|
: null
: bytes
| YES
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| mv_account_status
: bytes
| id
: int64 1
: bytes
| DUP
: null
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| mv_account_status
: bytes
| status
: int64 2
: bytes
| DUP
: null
: bytes
| NO
: bytes
| varchar(16)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| orders
: bytes
| id
: int64 1
: bytes
| DUP
: null
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| orders
: bytes
| account_id
: int64 2
: bytes
|
: null
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| orders
: bytes
| total
: int64 3
: bytes
|
: bytes
| 0
: bytes
| YES
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| orders
: bytes
| created_on
: int64 4
: bytes
|
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| samples
: bytes
| id
: int64 1
: bytes
| UNI
: null
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| samples
: bytes
| value
: int64 2
: bytes
|
: null
: bytes
| YES
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| tickets
: bytes
| id
: int64 1
: bytes
| UNI
: null
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| tickets
: bytes
| title
: int64 2
: bytes
|
: null
: bytes
| YES
: bytes
| varchar(100)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| tickets
: bytes
| slug
: int64 3
: bytes
|
: null
: bytes
| YES
: bytes
| varchar(120)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| tickets
: bytes
| status
: int64 4
: bytes
|
: bytes
| open
: bytes
| NO
: bytes
| varchar(16)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| v_open_accounts
: bytes
| acct_id
: int64 1
: bytes
|
: null
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| v_open_accounts
: bytes
| acct_email
: int64 2
: bytes
|
: null
: bytes
| NO
: bytes
| varchar(190)
: bytes
|
: bytes
|
: bytes
|
: bytes
|
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = 'omni_sync_corpus'
< result
< column "TABLE_NAME" "CHAR"
< column "VIEW_DEFINITION" "CHAR"
< row
: bytes
| v_open_accounts
: bytes
| SELECT `internal`.`omni_sync_corpus`.`accounts`.`id` AS `acct_id`, `internal`.`omni_sync_corpus`.`accounts`.`email` AS `acct_email` FROM `internal`.`omni_sync_corpus`.`accounts` WHERE `internal`.`omni_sync_corpus`.`accounts`.`status` = 'open'
> query
| SELECT Name, RefreshInfo, QuerySql FROM mv_infos("database"="omni_sync_corpus") ORDER BY Name
< result
< column "Name" "TEXT"
< column "RefreshInfo" "TEXT"
< column "QuerySql" "TEXT"
< row
: bytes
| mv_account_status
: bytes
| BUILD IMMEDIATE REFRESH AUTO ON MANUAL
: bytes
| SELECT `internal`.`omni_sync_corpus`.`accounts`.`id`, `internal`.`omni_sync_corpus`.`accounts`.`status` FROM `internal`.`omni_sync_corpus`.`accounts`
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			TABLE_TYPE,
| 			IFNULL(TABLE_ROWS, 0),
| 			IFNULL(DATA_LENGTH, 0),
| 			IFNULL(INDEX_LENGTH, 0),
| 			IFNULL(DATA_FREE, 0),
| 			IFNULL(CREATE_OPTIONS, ''),
| 			IFNULL(TABLE_COMMENT, '')
| 		FROM information_schema.TABLES
| 		WHERE TABLE_SCHEMA = 'omni_sync_corpus'
| 		ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "CHAR"
< column "IFNULL(TABLE_ROWS, 0)" "BIGINT"
< column "IFNULL(DATA_LENGTH, 0)" "BIGINT"
< column "IFNULL(INDEX_LENGTH, 0)" "BIGINT"
< column "IFNULL(DATA_FREE, 0)" "BIGINT"
< column "IFNULL(CREATE_OPTIONS, '')" "CHAR"
< column "IFNULL(TABLE_COMMENT, '')" "CHAR"
< row
: bytes
| accounts
: bytes
| BASE TABLE
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| customer accounts
< row
: bytes
| daily_rollup
: bytes
| BASE TABLE
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
|
< row
: bytes
| events
: bytes
| BASE TABLE
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
|
< row
: bytes
| metrics
: bytes
| BASE TABLE
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
|
< row
: bytes
| mv_account_status
: bytes
| BASE TABLE
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
|
< row
: bytes
| orders
: bytes
| BASE TABLE
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
|
< row
: bytes
| samples
: bytes
| BASE TABLE
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
|
< row
: bytes
| tickets
: bytes
| BASE TABLE
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
|
< row
: bytes
| v_open_accounts
: bytes
| VIEW
: int64 -1
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
|
> query
| SHOW CREATE TABLE `omni_sync_corpus`.`accounts`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| accounts
: bytes
| CREATE TABLE `accounts` (
|   `id` bigint NOT NULL,
|   `email` varchar(190) NOT NULL COMMENT "login email",
|   `bio` text NULL,
|   `status` varchar(16) NOT NULL DEFAULT "open",
|   `balance` decimal(12,4) NOT NULL DEFAULT "0.0000",
|   `flag` tinyint NOT NULL DEFAULT "1",
|   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   `updated_at` datetime(6) NULL DEFAULT CURRENT_TIMESTAMP(6),
|   `note` varchar(100) NULL COMMENT "owner''s note",
|   `dt` date NULL,
|   `payload` json NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`id`, `email`)
| COMMENT 'customer accounts'
| DISTRIBUTED BY HASH(`id`) BUCKETS 1
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728"
| );
> query
| SHOW INDEX FROM `omni_sync_corpus`.`accounts`
< result
< column "Table" "CHAR"
< column "Non_unique" "CHAR"
< column "Key_name" "CHAR"
< column "Seq_in_index" "CHAR"
< column "Column_name" "CHAR"
< column "Collation" "CHAR"
< column "Cardinality" "CHAR"
< column "Sub_part" "CHAR"
< column "Packed" "CHAR"
< column "Null" "CHAR"
< column "Index_type" "CHAR"
< column "Comment" "CHAR"
< column "Properties" "CHAR"
> query
| SHOW CREATE TABLE `omni_sync_corpus`.`daily_rollup`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| daily_rollup
: bytes
| CREATE TABLE `daily_rollup` (
|   `day` date NOT NULL,
|   `region` varchar(16) NOT NULL,
|   `amount` decimal(12,2) SUM NULL,
|   `peak` int MAX NULL,
|   `label` varchar(32) REPLACE NULL,
|   `visitors` hll HLL_UNION NOT NULL,
|   `seen` bitmap BITMAP_UNION NOT NULL DEFAULT BITMAP_EMPTY
| ) ENGINE=OLAP
| AGGREGATE KEY(`day`, `region`)
| DISTRIBUTED BY HASH(`day`) BUCKETS 2
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728"
| );
> query
| SHOW INDEX FROM `omni_sync_corpus`.`daily_rollup`
< result
< column "Table" "CHAR"
< column "Non_unique" "CHAR"
< column "Key_name" "CHAR"
< column "Seq_in_index" "CHAR"
< column "Column_name" "CHAR"
< column "Collation" "CHAR"
< column "Cardinality" "CHAR"
< column "Sub_part" "CHAR"
< column "Packed" "CHAR"
< column "Null" "CHAR"
< column "Index_type" "CHAR"
< column "Comment" "CHAR"
< column "Properties" "CHAR"
> query
| SHOW CREATE TABLE `omni_sync_corpus`.`events`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| events
: bytes
| CREATE TABLE `events` (
|   `id` bigint NOT NULL,
|   `name` varchar(64) NOT NULL COMMENT "event name",
|   `payload` text NULL,
|   `slug` varchar(96) AS (concat(`name`, '-evt')) NULL,
|   INDEX idx_name (`name`) USING INVERTED,
|   INDEX idx_payload (`payload`) USING INVERTED PROPERTIES("lower_case" = "true", "parser" = "english", "support_phrase" = "true")
| ) ENGINE=OLAP
| DUPLICATE KEY(`id`)
| DISTRIBUTED BY RANDOM BUCKETS 3
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "true",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728"
| );
> query
| SHOW INDEX FROM `omni_sync_corpus`.`events`
< result
< column "Table" "CHAR"
< column "Non_unique" "CHAR"
< column "Key_name" "CHAR"
< column "Seq_in_index" "CHAR"
< column "Column_name" "CHAR"
< column "Collation" "CHAR"
< column "Cardinality" "CHAR"
< column "Sub_part" "CHAR"
< column "Packed" "CHAR"
< column "Null" "CHAR"
< column "Index_type" "CHAR"
< column "Comment" "CHAR"
< column "Properties" "CHAR"
< row
: bytes
| omni_sync_corpus.events
: bytes
|
: bytes
| idx_name
: bytes
|
: bytes
| name
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
| INVERTED
: bytes
|
: bytes
|
< row
: bytes
| omni_sync_corpus.events
: bytes
|
: bytes
| idx_payload
: bytes
|
: bytes
| payload
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
| INVERTED
: bytes
|
: bytes
| ("lower_case" = "true", "parser" = "english", "support_phrase" = "true")
> query
| SHOW CREATE TABLE `omni_sync_corpus`.`metrics`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| metrics
: bytes
| CREATE TABLE `metrics` (
|   `ts` datetime NOT NULL,
|   `val` bigint NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`ts`)
| PARTITION BY RANGE(`ts`)
| (PARTITION p2025 VALUES [('2025-01-01 00:00:00'), ('2026-01-01 00:00:00')),
| PARTITION p2026 VALUES [('2026-01-01 00:00:00'), ('2027-01-01 00:00:00')))
| DISTRIBUTED BY HASH(`ts`) BUCKETS 1
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728"
| );
> query
| SHOW INDEX FROM `omni_sync_corpus`.`metrics`
< result
< column "Table" "CHAR"
< column "Non_unique" "CHAR"
< column "Key_name" "CHAR"
< column "Seq_in_index" "CHAR"
< column "Column_name" "CHAR"
< column "Collation" "CHAR"
< column "Cardinality" "CHAR"
< column "Sub_part" "CHAR"
< column "Packed" "CHAR"
< column "Null" "CHAR"
< column "Index_type" "CHAR"
< column "Comment" "CHAR"
< column "Properties" "CHAR"
> query
| SHOW CREATE TABLE `omni_sync_corpus`.`orders`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| orders
: bytes
| CREATE TABLE `orders` (
|   `id` bigint NOT NULL,
|   `account_id` bigint NOT NULL,
|   `total` int NULL DEFAULT "0",
|   `created_on` date NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`id`)
| DISTRIBUTED BY HASH(`id`) BUCKETS 1
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728"
| );
> query
| SHOW INDEX FROM `omni_sync_corpus`.`orders`
< result
< column "Table" "CHAR"
< column "Non_unique" "CHAR"
< column "Key_name" "CHAR"
< column "Seq_in_index" "CHAR"
< column "Column_name" "CHAR"
< column "Collation" "CHAR"
< column "Cardinality" "CHAR"
< column "Sub_part" "CHAR"
< column "Packed" "CHAR"
< column "Null" "CHAR"
< column "Index_type" "CHAR"
< column "Comment" "CHAR"
< column "Properties" "CHAR"
> query
| SHOW CREATE TABLE `omni_sync_corpus`.`samples`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| samples
: bytes
| CREATE TABLE `samples` (
|   `id` bigint NOT NULL,
|   `value` int NULL
| ) ENGINE=OLAP
| UNIQUE KEY(`id`)
| DISTRIBUTED BY HASH(`id`) BUCKETS AUTO
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "enable_unique_key_merge_on_write" = "true",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728",
| "enable_mow_light_delete" = "false"
| );
> query
| SHOW INDEX FROM `omni_sync_corpus`.`samples`
< result
< column "Table" "CHAR"
< column "Non_unique" "CHAR"
< column "Key_name" "CHAR"
< column "Seq_in_index" "CHAR"
< column "Column_name" "CHAR"
< column "Collation" "CHAR"
< column "Cardinality" "CHAR"
< column "Sub_part" "CHAR"
< column "Packed" "CHAR"
< column "Null" "CHAR"
< column "Index_type" "CHAR"
< column "Comment" "CHAR"
< column "Properties" "CHAR"
> query
| SHOW CREATE TABLE `omni_sync_corpus`.`tickets`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| tickets
: bytes
| CREATE TABLE `tickets` (
|   `id` bigint NOT NULL AUTO_INCREMENT(1000),
|   `title` varchar(100) NULL,
|   `slug` varchar(120) AS (concat(`title`, '-t')) NULL,
|   `status` varchar(16) NOT NULL DEFAULT "open"
| ) ENGINE=OLAP
| UNIQUE KEY(`id`)
| DISTRIBUTED BY HASH(`id`) BUCKETS 2
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "enable_unique_key_merge_on_write" = "true",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728",
| "enable_mow_light_delete" = "false"
| );
> query
| SHOW INDEX FROM `omni_sync_corpus`.`tickets`
< result
< column "Table" "CHAR"
< column "Non_unique" "CHAR"
< column "Key_name" "CHAR"
< column "Seq_in_index" "CHAR"
< column "Column_name" "CHAR"
< column "Collation" "CHAR"
< column "Cardinality" "CHAR"
< column "Sub_part" "CHAR"
< column "Packed" "CHAR"
< column "Null" "CHAR"
< column "Index_type" "CHAR"
< column "Comment" "CHAR"
< column "Properties" "CHAR"
> query
|
| 		SELECT ROUTINE_NAME
| 		FROM information_schema.ROUTINES
| 		WHERE ROUTINE_SCHEMA = 'omni_sync_corpus' AND ROUTINE_TYPE = 'PROCEDURE'
| 		ORDER BY ROUTINE_NAME
< result
< column "ROUTINE_NAME" "CHAR"
< row
: bytes
| BUMP_TOTAL
> query
| SHOW CREATE PROCEDURE `omni_sync_corpus`.`BUMP_TOTAL`
< result
< column "Procedure" "TEXT"
< column "Create Procedure" "TEXT"
< row
: bytes
| BUMP_TOTAL
: bytes
| CREATE PROCEDURE bump_total(IN p_id BIGINT, IN delta INT) COMMENT 'add to an order total' BEGIN UPDATE samples SET value = value + delta WHERE id = p_id; END
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = 'omni_sync_corpus'
< result
< column "DEFAULT_CHARACTER_SET_NAME" "CHAR"
< column "DEFAULT_COLLATION_NAME" "CHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
