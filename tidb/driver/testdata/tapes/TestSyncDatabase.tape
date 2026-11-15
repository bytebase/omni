-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSyncDatabase default#1
> exec
| USE `omni_sync_corpus`
< ok 0 0
> exec
| CREATE TABLE accounts (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	email VARCHAR(190) NOT NULL,
| 	bio TEXT,
| 	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
| 	literal_ts VARCHAR(40) NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
| 	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
| 	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL,
| 	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED,
| 	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
| 	note VARCHAR(100) COMMENT 'owner''s note',
| 	UNIQUE KEY uk_email (email),
| 	KEY idx_bio_prefix (bio(32)),
| 	KEY idx_expr ((LOWER(status)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='customer accounts';
< ok 0 0
> exec
| CREATE TABLE sharded_events (
| 	id BIGINT PRIMARY KEY AUTO_RANDOM(5, 64),
| 	account_id BIGINT NOT NULL,
| 	payload VARCHAR(128),
| 	KEY idx_account (account_id),
| 	CONSTRAINT fk_events_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
< ok 0 0
> exec
| CREATE TABLE metrics_range (
| 	ts INT NOT NULL,
| 	val BIGINT NOT NULL,
| 	PRIMARY KEY (ts)
| ) PARTITION BY RANGE (ts) (
| 	PARTITION p_low VALUES LESS THAN (100),
| 	PARTITION p_high VALUES LESS THAN MAXVALUE
| );
< ok 0 0
> exec
| CREATE SEQUENCE seq_orders START WITH 10 INCREMENT BY 2 MINVALUE 1 MAXVALUE 10000 CACHE 50 CYCLE COMMENT='order seq'; -- not restorable: TestDumpGaps
< ok 0 0
> exec
| CREATE VIEW v_open_accounts AS
| 	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open';
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
< row
: bytes
| 8.0.11-TiDB-v8.5.5
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			KEY_NAME,
| 			COLUMN_NAME,
| 			SEQ_IN_INDEX,
| 			NON_UNIQUE,
| 			INDEX_COMMENT,
| 			SUB_PART,
| 			EXPRESSION
| 		FROM information_schema.TIDB_INDEXES
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "SEQ_IN_INDEX" "BIGINT"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| accounts
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| accounts
: bytes
| idx_bio_prefix
: bytes
| bio
: int64 1
: int64 1
: bytes
|
: int64 32
: null
< row
: bytes
| accounts
: bytes
| idx_expr
: bytes
| NULL
: int64 1
: int64 1
: bytes
|
: null
: bytes
| lower(`status`)
< row
: bytes
| accounts
: bytes
| uk_email
: bytes
| email
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| metrics_range
: bytes
| PRIMARY
: bytes
| ts
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| sharded_events
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| sharded_events
: bytes
| idx_account
: bytes
| account_id
: int64 1
: int64 1
: bytes
|
: null
: null
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			IFNULL(COLUMN_NAME, ''),
| 			ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE,
| 			COLUMN_TYPE,
| 			IFNULL(CHARACTER_SET_NAME, ''),
| 			IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT),
| 			GENERATION_EXPRESSION,
| 			EXTRA
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "BIGINT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "GENERATION_EXPRESSION" "TEXT"
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
| bigint
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
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
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
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
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
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
| 'open'
: bytes
| NO
: bytes
| enum('open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
< row
: bytes
| accounts
: bytes
| literal_ts
: int64 5
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
< row
: bytes
| accounts
: bytes
| balance
: int64 6
: bytes
| '0.0000'
: bytes
| NO
: bytes
| decimal(12,4)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
< row
: bytes
| accounts
: bytes
| email_lower
: int64 7
: null
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| lower(`email`)
: bytes
| VIRTUAL GENERATED
< row
: bytes
| accounts
: bytes
| display_tag
: int64 8
: null
: bytes
| YES
: bytes
| varchar(220)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| concat(`email`, _utf8mb4':', _utf8mb4'v1')
: bytes
| STORED GENERATED
< row
: bytes
| accounts
: bytes
| created_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
< row
: bytes
| accounts
: bytes
| updated_at
: int64 10
: bytes
| 'CURRENT_TIMESTAMP(6)'
: bytes
| NO
: bytes
| datetime(6)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP(6)
< row
: bytes
| accounts
: bytes
| note
: int64 11
: null
: bytes
| YES
: bytes
| varchar(100)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| 'owner\'s note'
: bytes
|
: bytes
|
< row
: bytes
| metrics_range
: bytes
| ts
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
< row
: bytes
| metrics_range
: bytes
| val
: int64 2
: null
: bytes
| NO
: bytes
| bigint
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
< row
: bytes
| sharded_events
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| bigint
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
< row
: bytes
| sharded_events
: bytes
| account_id
: int64 2
: null
: bytes
| NO
: bytes
| bigint
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
< row
: bytes
| sharded_events
: bytes
| payload
: int64 3
: null
: bytes
| YES
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
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
| YES
: bytes
| bigint
: bytes
|
: bytes
|
: bytes
| ''
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
| YES
: bytes
| varchar(190)
: bytes
|
: bytes
|
: bytes
| ''
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
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_open_accounts
: bytes
| SELECT `id` AS `acct_id`,`email` AS `acct_email` FROM `omni_sync_corpus`.`accounts` WHERE `status`=_UTF8MB4'open'
> query
|
| 		SELECT
| 			kcu.TABLE_SCHEMA,
| 			fks.TABLE_NAME,
| 			fks.CONSTRAINT_NAME,
| 			kcu.COLUMN_NAME,
| 			kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME,
| 			kcu.REFERENCED_COLUMN_NAME,
| 			fks.DELETE_RULE,
| 			fks.UPDATE_RULE,
| 			fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION;
| 	
: string
| omni_sync_corpus
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
< row
: bytes
| omni_sync_corpus
: bytes
| sharded_events
: bytes
| fk_events_account
: bytes
| account_id
: bytes
| omni_sync_corpus
: bytes
| accounts
: bytes
| id
: bytes
| CASCADE
: bytes
| RESTRICT
: bytes
| NONE
> query
|
| 		SELECT
| 			tc.TABLE_NAME,
| 			cc.CONSTRAINT_NAME,
| 			cc.CHECK_CLAUSE
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = ?
| 		ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 	
: string
| omni_sync_corpus
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
|
| 		SELECT
| 			SEQUENCE_NAME,
| 			START,
| 			MIN_VALUE,
| 			MAX_VALUE,
| 			INCREMENT,
| 			CYCLE,
| 			CACHE_VALUE,
| 			IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES
| 		WHERE SEQUENCE_SCHEMA = ?
| 		ORDER BY SEQUENCE_NAME
: string
| omni_sync_corpus
< result
< column "SEQUENCE_NAME" "VARCHAR"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT"
< column "CYCLE" "TINYINT"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR"
< row
: bytes
| seq_orders
: int64 10
: int64 1
: int64 10000
: int64 2
: int64 1
: int64 50
: bytes
| order seq
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			PARTITION_NAME,
| 			SUBPARTITION_NAME,
| 			PARTITION_METHOD,
| 			SUBPARTITION_METHOD,
| 			PARTITION_EXPRESSION,
| 			SUBPARTITION_EXPRESSION,
| 			PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "SUBPARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| metrics_range
: bytes
| p_low
: null
: bytes
| RANGE
: null
: bytes
| `ts`
: null
: bytes
| 100
< row
: bytes
| metrics_range
: bytes
| p_high
: null
: bytes
| RANGE
: null
: bytes
| `ts`
: null
: bytes
| MAXVALUE
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
| 			QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''),
| 			IFNULL(TIDB_PK_TYPE, '')
| 		FROM information_schema.TABLES
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLE_ROWS, 0)" "BIGINT"
< column "IFNULL(DATA_LENGTH, 0)" "BIGINT"
< column "IFNULL(INDEX_LENGTH, 0)" "BIGINT"
< column "IFNULL(DATA_FREE, 0)" "BIGINT"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR"
< row
: bytes
| accounts
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| 'customer accounts'
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
< row
: bytes
| metrics_range
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
| partitioned
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
< row
: bytes
| seq_orders
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 1
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
< row
: bytes
| sharded_events
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| ''
: bytes
| PK_AUTO_RANDOM_BITS=5
: bytes
| CLUSTERED
< row
: bytes
| v_open_accounts
: bytes
| VIEW
: bytes
|
: bytes
|
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
: bytes
| NONCLUSTERED
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_sync_corpus
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
