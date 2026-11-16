-- pin: default starrocks/allin1-ubuntu:3.4.10 sha256:8bc803ea7a5bb490c77117a456b32b792ccae698b1f0e47dfcafdac41542566d 3.4.10

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
| 	updated_at DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
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
| REFRESH ASYNC
| PROPERTIES("replication_num" = "1")
| AS SELECT id, status FROM accounts
< ok 0 0
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			IFNULL(COLUMN_NAME, ''),
| 			ORDINAL_POSITION,
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
< column "TABLE_NAME" "VARCHAR"
< column "ifnull(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "BIGINT"
< column "COLUMN_DEFAULT" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "VARCHAR"
< column "ifnull(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "ifnull(COLLATION_NAME, '')" "VARCHAR"
< column "COLUMN_COMMENT" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| accounts
: bytes
| id
: int64 1
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
: null
: bytes
| YES
: bytes
| varchar(65533)
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
| 0.0000
: bytes
| NO
: bytes
| decimal(12, 4)
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
| CURRENT_TIMESTAMP
: bytes
| YES
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
| note
: int64 9
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
| owner's note
: bytes
|
< row
: bytes
| accounts
: bytes
| dt
: int64 10
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
| metrics
: bytes
| ts
: int64 1
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
| v_open_accounts
: bytes
| acct_id
: int64 1
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
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "VARCHAR"
< row
: bytes
| v_open_accounts
: bytes
| SELECT `omni_sync_corpus`.`accounts`.`id` AS `acct_id`, `omni_sync_corpus`.`accounts`.`email` AS `acct_email`
| FROM `omni_sync_corpus`.`accounts`
| WHERE `omni_sync_corpus`.`accounts`.`status` = 'open'
> query
| SELECT TABLE_NAME, IFNULL(REFRESH_TYPE, ''), IFNULL(MATERIALIZED_VIEW_DEFINITION, '') FROM information_schema.materialized_views WHERE TABLE_SCHEMA = 'omni_sync_corpus'
< error
| 1064 HY000 FE RPC failure, address=TNetworkAddress(hostname=127.0.0.1, port=9020), reason=Internal error processing listMaterializedViewStatus: BE:10002, host: unknown
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			TABLE_TYPE,
| 			IFNULL(ENGINE, ''),
| 			IFNULL(TABLE_COLLATION, ''),
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
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "ifnull(ENGINE, '')" "VARCHAR"
< column "ifnull(TABLE_COLLATION, '')" "VARCHAR"
< column "ifnull(TABLE_ROWS, 0)" "BIGINT"
< column "ifnull(DATA_LENGTH, 0)" "BIGINT"
< column "ifnull(INDEX_LENGTH, 0)" "BIGINT"
< column "ifnull(DATA_FREE, 0)" "BIGINT"
< column "ifnull(CREATE_OPTIONS, '')" "VARCHAR"
< column "ifnull(TABLE_COMMENT, '')" "VARCHAR"
< row
: bytes
| accounts
: bytes
| BASE TABLE
: bytes
| StarRocks
: bytes
| utf8_general_ci
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
| metrics
: bytes
| BASE TABLE
: bytes
| StarRocks
: bytes
| utf8_general_ci
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
| VIEW
: bytes
| StarRocks
: bytes
| utf8_general_ci
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
: bytes
| StarRocks
: bytes
| utf8_general_ci
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
: bytes
|
: bytes
| utf8_general_ci
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
|
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = 'omni_sync_corpus'
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8
: bytes
| utf8_general_ci
