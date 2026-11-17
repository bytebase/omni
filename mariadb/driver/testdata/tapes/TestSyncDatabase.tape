-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8
-- pin: mariadb1011 mariadb:10.11.19 sha256:7d829180b5b2a9ba74f6242876ab401a0d7dd338a807e522ef92bc9698760f20 mariadb-10.11.19
-- pin: mariadb107 mariadb:10.7.8 sha256:9a48ac9f196f3d4fd6fea2cab59a49df9e7ca459bf14b2f7b85a0e38a5454571 mariadb-10.7.8
-- pin: mariadb114 mariadb:11.4.13 sha256:5e73baf8179b29f115101bb102763a47a017361f3e9792daccba6646c3ab1532 mariadb-11.4.13

== TestSyncDatabase/default default#1
> exec
| CREATE TABLE accounts (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	email VARCHAR(190) NOT NULL,
| 	bio TEXT,
| 	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
| 	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
| 	secret_code VARCHAR(32) INVISIBLE DEFAULT 'hidden',
| 	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL,
| 	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED,
| 	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
| 	note VARCHAR(100) COMMENT 'owner''s note',
| 	UNIQUE KEY uk_email (email),
| 	KEY idx_bio_prefix (bio(32)),
| 	CONSTRAINT chk_balance CHECK (balance >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='customer accounts';
< ok 0 0
> exec
| CREATE TABLE orders (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	account_id BIGINT NOT NULL,
| 	total INT NOT NULL,
| 	KEY idx_account (account_id),
| 	CONSTRAINT fk_orders_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
< ok 0 0
> exec
| CREATE TABLE history (
| 	id INT PRIMARY KEY,
| 	amount INT,
| 	account_id BIGINT,
| 	KEY idx_amount (amount),
| 	CONSTRAINT fk_history_account FOREIGN KEY (account_id) REFERENCES accounts (id),
| 	CONSTRAINT chk_amount CHECK (amount > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING COMMENT='versioned';
< ok 0 0
> exec
| CREATE TABLE explicit_period (
| 	id INT PRIMARY KEY,
| 	valid_from TIMESTAMP(6) GENERATED ALWAYS AS ROW START,
| 	valid_to TIMESTAMP(6) GENERATED ALWAYS AS ROW END,
| 	PERIOD FOR SYSTEM_TIME (valid_from, valid_to)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING;
< ok 0 0
> exec
| CREATE TABLE metrics_range (
| 	ts INT NOT NULL,
| 	val BIGINT NOT NULL,
| 	PRIMARY KEY (ts)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY RANGE (ts) (
| 	PARTITION p_low VALUES LESS THAN (100),
| 	PARTITION p_high VALUES LESS THAN MAXVALUE
| );
< ok 0 0
> exec
| CREATE TABLE metrics_hash (
| 	id INT NOT NULL PRIMARY KEY,
| 	val INT NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY HASH (id) PARTITIONS 4;
< ok 0 0
> exec
| CREATE SEQUENCE seq_default;
< ok 0 0
> exec
| CREATE SEQUENCE seq_custom START WITH 100 INCREMENT BY 5 MINVALUE 10 MAXVALUE 1000 CACHE 20 CYCLE COMMENT='sequence comment';
< ok 0 0
> exec
| CREATE VIEW v_open_accounts AS
| 	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open';
< ok 0 0
> exec
| CREATE VIEW v_named_cols (`renamed_id`) AS SELECT id FROM accounts;
< ok 0 0
> exec
| CREATE VIEW v_named_cols_ref AS
| 	SELECT renamed_id FROM v_named_cols;
< ok 0 0
> exec
| CREATE FUNCTION fn_format_tag(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT)
| RETURNS VARCHAR(64) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci
| DETERMINISTIC
| NO SQL
| COMMENT 'formats a tag'
| BEGIN
| 	RETURN CONCAT(p_prefix, '-', p_id);
| END;
< ok 0 0
> exec
| CREATE PROCEDURE sp_touch_order(IN p_id BIGINT, OUT p_total INT)
| COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END;
< ok 0 0
> exec
| CREATE TRIGGER trg_orders_bi
| BEFORE INSERT ON orders
| FOR EACH ROW
| BEGIN
| 	IF NEW.total < 0 THEN
| 		SET NEW.total = 0;
| 	END IF;
| END;
< ok 0 0
> exec
| CREATE EVENT ev_daily_check
| ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00'
| DISABLE
| DO
| BEGIN
| 	DELETE FROM orders WHERE total < 0;
| END;
< ok 0 0
> exec
| CREATE SEQUENCE seq_typed AS INT UNSIGNED START WITH 7 MAXVALUE 4000000000;
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR" "string"
< row
: bytes
| 11.8.8-MariaDB-ubu2404
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			INDEX_NAME,
| 			COLUMN_NAME,
| 			COLLATION,
| 			IFNULL(SUB_PART, -1),
| 			SEQ_IN_INDEX,
| 			INDEX_TYPE,
| 			NON_UNIQUE,
| 			INDEX_COMMENT
| 		FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "BIGINT" "int64"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| idx_bio_prefix
: bytes
| bio
: bytes
| A
: int64 32
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| accounts
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| accounts
: bytes
| uk_email
: bytes
| email
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| explicit_period
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| explicit_period
: bytes
| PRIMARY
: bytes
| valid_to
: bytes
| A
: int64 -1
: int64 2
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| history
: bytes
| idx_amount
: bytes
| amount
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| history
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| metrics_hash
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| metrics_range
: bytes
| PRIMARY
: bytes
| ts
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| orders
: bytes
| idx_account
: bytes
| account_id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| orders
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
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
| 			convert(GENERATION_EXPRESSION using BINARY),
| 			EXTRA
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT" "uint64"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "MEDIUMTEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "LONGBLOB" "[]uint8"
< column "EXTRA" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| id
: uint64 1
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
| ''
: null
: bytes
| auto_increment
< row
: bytes
| accounts
: bytes
| email
: uint64 2
: null
: bytes
| NO
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| bio
: uint64 3
: bytes
| 'NULL'
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| status
: uint64 4
: bytes
| '\'open\''
: bytes
| NO
: bytes
| enum('open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| balance
: uint64 5
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
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| secret_code
: uint64 6
: bytes
| '\'hidden\''
: bytes
| YES
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
| INVISIBLE
< row
: bytes
| accounts
: bytes
| email_lower
: uint64 7
: bytes
| 'NULL'
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
| lcase(`email`)
: bytes
| VIRTUAL GENERATED
< row
: bytes
| accounts
: bytes
| display_tag
: uint64 8
: bytes
| 'NULL'
: bytes
| YES
: bytes
| varchar(220)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
| concat(`email`,':','v1')
: bytes
| STORED GENERATED
< row
: bytes
| accounts
: bytes
| created_at
: uint64 9
: bytes
| 'current_timestamp()'
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
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| updated_at
: uint64 10
: bytes
| 'current_timestamp(6)'
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
: null
: bytes
| on update current_timestamp(6)
< row
: bytes
| accounts
: bytes
| note
: uint64 11
: bytes
| 'NULL'
: bytes
| YES
: bytes
| varchar(100)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| 'owner\'s note'
: null
: bytes
|
< row
: bytes
| explicit_period
: bytes
| id
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| explicit_period
: bytes
| valid_from
: uint64 2
: null
: bytes
| NO
: bytes
| timestamp(6)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
| ROW START
: bytes
| STORED GENERATED
< row
: bytes
| explicit_period
: bytes
| valid_to
: uint64 3
: null
: bytes
| NO
: bytes
| timestamp(6)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
| ROW END
: bytes
| STORED GENERATED
< row
: bytes
| history
: bytes
| id
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| history
: bytes
| amount
: uint64 2
: bytes
| 'NULL'
: bytes
| YES
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| history
: bytes
| account_id
: uint64 3
: bytes
| 'NULL'
: bytes
| YES
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_hash
: bytes
| id
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_hash
: bytes
| val
: uint64 2
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_range
: bytes
| ts
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_range
: bytes
| val
: uint64 2
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
| ''
: null
: bytes
|
< row
: bytes
| orders
: bytes
| id
: uint64 1
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
| ''
: null
: bytes
| auto_increment
< row
: bytes
| orders
: bytes
| account_id
: uint64 2
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
| ''
: null
: bytes
|
< row
: bytes
| orders
: bytes
| total
: uint64 3
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| next_not_cached_value
: uint64 1
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| minimum_value
: uint64 2
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| maximum_value
: uint64 3
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| start_value
: uint64 4
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'start value when sequences is created or value if RESTART is used'
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| increment
: uint64 5
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'increment value'
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| cache_size
: uint64 6
: null
: bytes
| NO
: bytes
| bigint(21) unsigned
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| cycle_option
: uint64 7
: null
: bytes
| NO
: bytes
| tinyint(1) unsigned
: bytes
|
: bytes
|
: bytes
| '0 if no cycles are allowed, 1 if the sequence should begin a new cycle when maximum_value is passed'
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| cycle_count
: uint64 8
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'How many cycles have been done'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| next_not_cached_value
: uint64 1
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| minimum_value
: uint64 2
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| maximum_value
: uint64 3
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| start_value
: uint64 4
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'start value when sequences is created or value if RESTART is used'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| increment
: uint64 5
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'increment value'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| cache_size
: uint64 6
: null
: bytes
| NO
: bytes
| bigint(21) unsigned
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| cycle_option
: uint64 7
: null
: bytes
| NO
: bytes
| tinyint(1) unsigned
: bytes
|
: bytes
|
: bytes
| '0 if no cycles are allowed, 1 if the sequence should begin a new cycle when maximum_value is passed'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| cycle_count
: uint64 8
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'How many cycles have been done'
: null
: bytes
|
< row
: bytes
| seq_typed
: bytes
| next_not_cached_value
: uint64 1
: null
: bytes
| NO
: bytes
| int(12) unsigned
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_typed
: bytes
| minimum_value
: uint64 2
: null
: bytes
| NO
: bytes
| int(12) unsigned
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_typed
: bytes
| maximum_value
: uint64 3
: null
: bytes
| NO
: bytes
| int(12) unsigned
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_typed
: bytes
| start_value
: uint64 4
: null
: bytes
| NO
: bytes
| int(12) unsigned
: bytes
|
: bytes
|
: bytes
| 'start value when sequences is created or value if RESTART is used'
: null
: bytes
|
< row
: bytes
| seq_typed
: bytes
| increment
: uint64 5
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'increment value'
: null
: bytes
|
< row
: bytes
| seq_typed
: bytes
| cache_size
: uint64 6
: null
: bytes
| NO
: bytes
| bigint(21) unsigned
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_typed
: bytes
| cycle_option
: uint64 7
: null
: bytes
| NO
: bytes
| tinyint(1) unsigned
: bytes
|
: bytes
|
: bytes
| '0 if no cycles are allowed, 1 if the sequence should begin a new cycle when maximum_value is passed'
: null
: bytes
|
< row
: bytes
| seq_typed
: bytes
| cycle_count
: uint64 8
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'How many cycles have been done'
: null
: bytes
|
< row
: bytes
| v_named_cols
: bytes
| renamed_id
: uint64 1
: bytes
| '0'
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| v_named_cols_ref
: bytes
| renamed_id
: uint64 1
: bytes
| '0'
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| v_open_accounts
: bytes
| acct_id
: uint64 1
: bytes
| '0'
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| v_open_accounts
: bytes
| acct_email
: uint64 2
: null
: bytes
| NO
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				CONSTRAINT_NAME,
| 				CHECK_CLAUSE
| 			FROM information_schema.CHECK_CONSTRAINTS
| 			WHERE CONSTRAINT_SCHEMA = ?
| 			ORDER BY TABLE_NAME, CONSTRAINT_NAME
| 		
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| accounts
: bytes
| chk_balance
: bytes
| `balance` >= 0
< row
: bytes
| history
: bytes
| chk_amount
: bytes
| `amount` > 0
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v_named_cols
: bytes
| select `omni_sync_corpus`.`accounts`.`id` AS `renamed_id` from `omni_sync_corpus`.`accounts`
< row
: bytes
| v_named_cols_ref
: bytes
| select `v_named_cols`.`renamed_id` AS `renamed_id` from `omni_sync_corpus`.`v_named_cols`
< row
: bytes
| v_open_accounts
: bytes
| select `omni_sync_corpus`.`accounts`.`id` AS `acct_id`,`omni_sync_corpus`.`accounts`.`email` AS `acct_email` from `omni_sync_corpus`.`accounts` where `omni_sync_corpus`.`accounts`.`status` = 'open'
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_named_cols`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols` AS select `accounts`.`id` AS `renamed_id` from `accounts`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_named_cols_ref`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols_ref
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols_ref` AS select `v_named_cols`.`renamed_id` AS `renamed_id` from `v_named_cols`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_open_accounts`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_open_accounts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where `accounts`.`status` = 'open'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
|
| 	SELECT
| 		TRIGGER_NAME,
| 		EVENT_OBJECT_TABLE,
| 		EVENT_MANIPULATION,
| 		ACTION_TIMING,
| 		ACTION_STATEMENT,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.TRIGGERS
| 	WHERE EVENT_OBJECT_SCHEMA = ?
| 	ORDER BY EVENT_OBJECT_TABLE ASC, EVENT_MANIPULATION ASC, ACTION_TIMING ASC, ACTION_ORDER ASC;
| 	
: string
| omni_sync_corpus
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "VARCHAR" "string"
< column "ACTION_TIMING" "VARCHAR" "string"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "VARCHAR" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| trg_orders_bi
: bytes
| orders
: bytes
| INSERT
: bytes
| BEFORE
: bytes
| BEGIN
| 	IF NEW.total < 0 THEN
| 		SET NEW.total = 0;
| 	END IF;
| END
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
|
| 	SELECT
| 		EVENT_NAME,
| 		TIME_ZONE,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.EVENTS
| 	WHERE EVENT_SCHEMA = ?
| 	ORDER BY EVENT_NAME ASC;
| 	
: string
| omni_sync_corpus
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "VARCHAR" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| SYSTEM
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE EVENT `omni_sync_corpus`.`ev_daily_check`
< result
< column "Event" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "time_zone" "VARCHAR" "string"
< column "Create Event" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_daily_check` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO BEGIN
| 	DELETE FROM orders WHERE total < 0;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			REFERENCED_TABLE_NAME,
| 			DELETE_RULE,
| 			UPDATE_RULE,
| 			MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ?;
| 	
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR" "string"
< column "UPDATE_RULE" "VARCHAR" "string"
< column "MATCH_OPTION" "VARCHAR" "string"
< row
: bytes
| history
: bytes
| fk_history_account
: bytes
| accounts
: bytes
| RESTRICT
: bytes
| RESTRICT
: bytes
| NONE
< row
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| accounts
: bytes
| CASCADE
: bytes
| RESTRICT
: bytes
| NONE
> query
|
| 		SELECT
| 			TABLE_SCHEMA,
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			COLUMN_NAME,
| 			REFERENCED_TABLE_SCHEMA,
| 			REFERENCED_COLUMN_NAME
| 		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
| 		WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION;
| 	
: string
| omni_sync_corpus
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_sync_corpus
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| omni_sync_corpus
: bytes
| id
< row
: bytes
| omni_sync_corpus
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_sync_corpus
: bytes
| id
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "SUBPARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| metrics_hash
: bytes
| p0
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
< row
: bytes
| metrics_hash
: bytes
| p1
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
< row
: bytes
| metrics_hash
: bytes
| p2
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
< row
: bytes
| metrics_hash
: bytes
| p3
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
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
| SHOW CREATE TABLE `omni_sync_corpus`.`metrics_hash`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_hash
: bytes
| CREATE TABLE `metrics_hash` (
|   `id` int(11) NOT NULL,
|   `val` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
|  PARTITION BY HASH (`id`)
| PARTITIONS 4
> query
| SHOW CREATE TABLE `omni_sync_corpus`.`metrics_range`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_range
: bytes
| CREATE TABLE `metrics_range` (
|   `ts` int(11) NOT NULL,
|   `val` bigint(20) NOT NULL,
|   PRIMARY KEY (`ts`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
|  PARTITION BY RANGE (`ts`)
| (PARTITION `p_low` VALUES LESS THAN (100) ENGINE = InnoDB,
|  PARTITION `p_high` VALUES LESS THAN MAXVALUE ENGINE = InnoDB)
> query
|
| 		SELECT
| 			ROUTINE_NAME,
| 			ROUTINE_TYPE,
| 			SQL_MODE,
| 			CHARACTER_SET_CLIENT,
| 			COLLATION_CONNECTION,
| 			DATABASE_COLLATION,
| 			IFNULL(ROUTINE_COMMENT, '')
| 		FROM
| 			INFORMATION_SCHEMA.ROUTINES
| 		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE IN ('FUNCTION', 'PROCEDURE')
| 		ORDER BY ROUTINE_TYPE, ROUTINE_NAME;
| 	
: string
| omni_sync_corpus
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "VARCHAR" "string"
< column "SQL_MODE" "VARCHAR" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "LONGTEXT" "string"
< row
: bytes
| fn_format_tag
: bytes
| FUNCTION
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| formats a tag
< row
: bytes
| sp_touch_order
: bytes
| PROCEDURE
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| reads order total
> query
| SHOW CREATE FUNCTION `omni_sync_corpus`.`fn_format_tag`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| fn_format_tag
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci
|     NO SQL
|     DETERMINISTIC
|     COMMENT 'formats a tag'
| BEGIN
| 	RETURN CONCAT(p_prefix, '-', p_id);
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE PROCEDURE `omni_sync_corpus`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| sp_touch_order
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `sp_touch_order`(IN p_id BIGINT, OUT p_total INT)
|     COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
|
| 		SELECT
| 			TABLES.TABLE_NAME,
| 			TABLES.TABLE_TYPE,
| 			IFNULL(TABLES.ENGINE, ''),
| 			IFNULL(TABLES.TABLE_COLLATION, ''),
| 			IFNULL(TABLES.TABLE_ROWS, 0),
| 			IFNULL(TABLES.DATA_LENGTH, 0),
| 			IFNULL(TABLES.INDEX_LENGTH, 0),
| 			IFNULL(TABLES.DATA_FREE, 0),
| 			IFNULL(TABLES.CREATE_OPTIONS, ''),
| 			QUOTE(IFNULL(TABLES.TABLE_COMMENT, '')),
| 			IFNULL(CCSA.CHARACTER_SET_NAME, '')
| 		FROM information_schema.TABLES TABLES
| 		LEFT JOIN information_schema.COLLATION_CHARACTER_SET_APPLICABILITY CCSA
| 		ON TABLES.TABLE_COLLATION = CCSA.COLLATION_NAME
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "TABLE_TYPE" "VARCHAR" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "TEXT" "string"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| 'customer accounts'
: bytes
| utf8mb4
< row
: bytes
| explicit_period
: bytes
| SYSTEM VERSIONED
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| history
: bytes
| SYSTEM VERSIONED
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| 'versioned'
: bytes
| utf8mb4
< row
: bytes
| metrics_hash
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 4
: bytes
| 65536
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| metrics_range
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 2
: bytes
| 32768
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| orders
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| seq_custom
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| 1
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'sequence comment'
: bytes
|
< row
: bytes
| seq_default
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| 1
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
|
< row
: bytes
| seq_typed
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| 1
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
|
< row
: bytes
| v_named_cols
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| v_named_cols_ref
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
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
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_sync_corpus`.`seq_custom`
< result
< column "start_value" "BIGINT" "int64"
< column "minimum_value" "BIGINT" "int64"
< column "maximum_value" "BIGINT" "int64"
< column "increment" "BIGINT" "int64"
< column "cache_size" "UNSIGNED BIGINT" "uint64"
< column "cycle_option" "UNSIGNED TINYINT" "uint8"
< row
: int64 100
: int64 10
: int64 1000
: int64 5
: uint64 20
: int64 1
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_sync_corpus`.`seq_default`
< result
< column "start_value" "BIGINT" "int64"
< column "minimum_value" "BIGINT" "int64"
< column "maximum_value" "BIGINT" "int64"
< column "increment" "BIGINT" "int64"
< column "cache_size" "UNSIGNED BIGINT" "uint64"
< column "cycle_option" "UNSIGNED TINYINT" "uint8"
< row
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: uint64 1000
: int64 0
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_sync_corpus`.`seq_typed`
< result
< column "start_value" "UNSIGNED INT" "uint32"
< column "minimum_value" "UNSIGNED INT" "uint32"
< column "maximum_value" "UNSIGNED INT" "uint32"
< column "increment" "BIGINT" "int64"
< column "cache_size" "UNSIGNED BIGINT" "uint64"
< column "cycle_option" "UNSIGNED TINYINT" "uint8"
< row
: int64 7
: int64 1
: int64 4000000000
: int64 1
: uint64 1000
: int64 0
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
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci

== TestSyncDatabase/mariadb1011 mariadb1011#1
> exec
| CREATE TABLE accounts (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	email VARCHAR(190) NOT NULL,
| 	bio TEXT,
| 	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
| 	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
| 	secret_code VARCHAR(32) INVISIBLE DEFAULT 'hidden',
| 	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL,
| 	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED,
| 	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
| 	note VARCHAR(100) COMMENT 'owner''s note',
| 	UNIQUE KEY uk_email (email),
| 	KEY idx_bio_prefix (bio(32)),
| 	CONSTRAINT chk_balance CHECK (balance >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='customer accounts';
< ok 0 0
> exec
| CREATE TABLE orders (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	account_id BIGINT NOT NULL,
| 	total INT NOT NULL,
| 	KEY idx_account (account_id),
| 	CONSTRAINT fk_orders_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
< ok 0 0
> exec
| CREATE TABLE history (
| 	id INT PRIMARY KEY,
| 	amount INT,
| 	account_id BIGINT,
| 	KEY idx_amount (amount),
| 	CONSTRAINT fk_history_account FOREIGN KEY (account_id) REFERENCES accounts (id),
| 	CONSTRAINT chk_amount CHECK (amount > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING COMMENT='versioned';
< ok 0 0
> exec
| CREATE TABLE explicit_period (
| 	id INT PRIMARY KEY,
| 	valid_from TIMESTAMP(6) GENERATED ALWAYS AS ROW START,
| 	valid_to TIMESTAMP(6) GENERATED ALWAYS AS ROW END,
| 	PERIOD FOR SYSTEM_TIME (valid_from, valid_to)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING;
< ok 0 0
> exec
| CREATE TABLE metrics_range (
| 	ts INT NOT NULL,
| 	val BIGINT NOT NULL,
| 	PRIMARY KEY (ts)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY RANGE (ts) (
| 	PARTITION p_low VALUES LESS THAN (100),
| 	PARTITION p_high VALUES LESS THAN MAXVALUE
| );
< ok 0 0
> exec
| CREATE TABLE metrics_hash (
| 	id INT NOT NULL PRIMARY KEY,
| 	val INT NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY HASH (id) PARTITIONS 4;
< ok 0 0
> exec
| CREATE SEQUENCE seq_default;
< ok 0 0
> exec
| CREATE SEQUENCE seq_custom START WITH 100 INCREMENT BY 5 MINVALUE 10 MAXVALUE 1000 CACHE 20 CYCLE COMMENT='sequence comment';
< ok 0 0
> exec
| CREATE VIEW v_open_accounts AS
| 	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open';
< ok 0 0
> exec
| CREATE VIEW v_named_cols (`renamed_id`) AS SELECT id FROM accounts;
< ok 0 0
> exec
| CREATE VIEW v_named_cols_ref AS
| 	SELECT renamed_id FROM v_named_cols;
< ok 0 0
> exec
| CREATE FUNCTION fn_format_tag(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT)
| RETURNS VARCHAR(64) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci
| DETERMINISTIC
| NO SQL
| COMMENT 'formats a tag'
| BEGIN
| 	RETURN CONCAT(p_prefix, '-', p_id);
| END;
< ok 0 0
> exec
| CREATE PROCEDURE sp_touch_order(IN p_id BIGINT, OUT p_total INT)
| COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END;
< ok 0 0
> exec
| CREATE TRIGGER trg_orders_bi
| BEFORE INSERT ON orders
| FOR EACH ROW
| BEGIN
| 	IF NEW.total < 0 THEN
| 		SET NEW.total = 0;
| 	END IF;
| END;
< ok 0 0
> exec
| CREATE EVENT ev_daily_check
| ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00'
| DISABLE
| DO
| BEGIN
| 	DELETE FROM orders WHERE total < 0;
| END;
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR" "string"
< row
: bytes
| 10.11.19-MariaDB-ubu2204
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			INDEX_NAME,
| 			COLUMN_NAME,
| 			COLLATION,
| 			IFNULL(SUB_PART, -1),
| 			SEQ_IN_INDEX,
| 			INDEX_TYPE,
| 			NON_UNIQUE,
| 			INDEX_COMMENT
| 		FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "BIGINT" "int64"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| idx_bio_prefix
: bytes
| bio
: bytes
| A
: int64 32
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| accounts
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| accounts
: bytes
| uk_email
: bytes
| email
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| explicit_period
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| explicit_period
: bytes
| PRIMARY
: bytes
| valid_to
: bytes
| A
: int64 -1
: int64 2
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| history
: bytes
| idx_amount
: bytes
| amount
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| history
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| metrics_hash
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| metrics_range
: bytes
| PRIMARY
: bytes
| ts
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| orders
: bytes
| idx_account
: bytes
| account_id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| orders
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
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
| 			convert(GENERATION_EXPRESSION using BINARY),
| 			EXTRA
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT" "uint64"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "MEDIUMTEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "LONGBLOB" "[]uint8"
< column "EXTRA" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| id
: uint64 1
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
| ''
: null
: bytes
| auto_increment
< row
: bytes
| accounts
: bytes
| email
: uint64 2
: null
: bytes
| NO
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| bio
: uint64 3
: bytes
| 'NULL'
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| status
: uint64 4
: bytes
| '\'open\''
: bytes
| NO
: bytes
| enum('open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| balance
: uint64 5
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
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| secret_code
: uint64 6
: bytes
| '\'hidden\''
: bytes
| YES
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
| INVISIBLE
< row
: bytes
| accounts
: bytes
| email_lower
: uint64 7
: bytes
| 'NULL'
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
| lcase(`email`)
: bytes
| VIRTUAL GENERATED
< row
: bytes
| accounts
: bytes
| display_tag
: uint64 8
: bytes
| 'NULL'
: bytes
| YES
: bytes
| varchar(220)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
| concat(`email`,':','v1')
: bytes
| STORED GENERATED
< row
: bytes
| accounts
: bytes
| created_at
: uint64 9
: bytes
| 'current_timestamp()'
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
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| updated_at
: uint64 10
: bytes
| 'current_timestamp(6)'
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
: null
: bytes
| on update current_timestamp(6)
< row
: bytes
| accounts
: bytes
| note
: uint64 11
: bytes
| 'NULL'
: bytes
| YES
: bytes
| varchar(100)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| 'owner\'s note'
: null
: bytes
|
< row
: bytes
| explicit_period
: bytes
| id
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| explicit_period
: bytes
| valid_from
: uint64 2
: null
: bytes
| NO
: bytes
| timestamp(6)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
| ROW START
: bytes
| STORED GENERATED
< row
: bytes
| explicit_period
: bytes
| valid_to
: uint64 3
: null
: bytes
| NO
: bytes
| timestamp(6)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
| ROW END
: bytes
| STORED GENERATED
< row
: bytes
| history
: bytes
| id
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| history
: bytes
| amount
: uint64 2
: bytes
| 'NULL'
: bytes
| YES
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| history
: bytes
| account_id
: uint64 3
: bytes
| 'NULL'
: bytes
| YES
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_hash
: bytes
| id
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_hash
: bytes
| val
: uint64 2
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_range
: bytes
| ts
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_range
: bytes
| val
: uint64 2
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
| ''
: null
: bytes
|
< row
: bytes
| orders
: bytes
| id
: uint64 1
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
| ''
: null
: bytes
| auto_increment
< row
: bytes
| orders
: bytes
| account_id
: uint64 2
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
| ''
: null
: bytes
|
< row
: bytes
| orders
: bytes
| total
: uint64 3
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| next_not_cached_value
: uint64 1
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| minimum_value
: uint64 2
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| maximum_value
: uint64 3
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| start_value
: uint64 4
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'start value when sequences is created or value if RESTART is used'
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| increment
: uint64 5
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'increment value'
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| cache_size
: uint64 6
: null
: bytes
| NO
: bytes
| bigint(21) unsigned
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| cycle_option
: uint64 7
: null
: bytes
| NO
: bytes
| tinyint(1) unsigned
: bytes
|
: bytes
|
: bytes
| '0 if no cycles are allowed, 1 if the sequence should begin a new cycle when maximum_value is passed'
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| cycle_count
: uint64 8
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'How many cycles have been done'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| next_not_cached_value
: uint64 1
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| minimum_value
: uint64 2
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| maximum_value
: uint64 3
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| start_value
: uint64 4
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'start value when sequences is created or value if RESTART is used'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| increment
: uint64 5
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'increment value'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| cache_size
: uint64 6
: null
: bytes
| NO
: bytes
| bigint(21) unsigned
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| cycle_option
: uint64 7
: null
: bytes
| NO
: bytes
| tinyint(1) unsigned
: bytes
|
: bytes
|
: bytes
| '0 if no cycles are allowed, 1 if the sequence should begin a new cycle when maximum_value is passed'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| cycle_count
: uint64 8
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'How many cycles have been done'
: null
: bytes
|
< row
: bytes
| v_named_cols
: bytes
| renamed_id
: uint64 1
: bytes
| '0'
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| v_named_cols_ref
: bytes
| renamed_id
: uint64 1
: bytes
| '0'
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| v_open_accounts
: bytes
| acct_id
: uint64 1
: bytes
| '0'
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| v_open_accounts
: bytes
| acct_email
: uint64 2
: null
: bytes
| NO
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				CONSTRAINT_NAME,
| 				CHECK_CLAUSE
| 			FROM information_schema.CHECK_CONSTRAINTS
| 			WHERE CONSTRAINT_SCHEMA = ?
| 			ORDER BY TABLE_NAME, CONSTRAINT_NAME
| 		
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| accounts
: bytes
| chk_balance
: bytes
| `balance` >= 0
< row
: bytes
| history
: bytes
| chk_amount
: bytes
| `amount` > 0
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v_named_cols
: bytes
| select `omni_sync_corpus`.`accounts`.`id` AS `renamed_id` from `omni_sync_corpus`.`accounts`
< row
: bytes
| v_named_cols_ref
: bytes
| select `v_named_cols`.`renamed_id` AS `renamed_id` from `omni_sync_corpus`.`v_named_cols`
< row
: bytes
| v_open_accounts
: bytes
| select `omni_sync_corpus`.`accounts`.`id` AS `acct_id`,`omni_sync_corpus`.`accounts`.`email` AS `acct_email` from `omni_sync_corpus`.`accounts` where `omni_sync_corpus`.`accounts`.`status` = 'open'
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_named_cols`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols` AS select `accounts`.`id` AS `renamed_id` from `accounts`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_named_cols_ref`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols_ref
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols_ref` AS select `v_named_cols`.`renamed_id` AS `renamed_id` from `v_named_cols`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_open_accounts`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_open_accounts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where `accounts`.`status` = 'open'
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
|
| 	SELECT
| 		TRIGGER_NAME,
| 		EVENT_OBJECT_TABLE,
| 		EVENT_MANIPULATION,
| 		ACTION_TIMING,
| 		ACTION_STATEMENT,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.TRIGGERS
| 	WHERE EVENT_OBJECT_SCHEMA = ?
| 	ORDER BY EVENT_OBJECT_TABLE ASC, EVENT_MANIPULATION ASC, ACTION_TIMING ASC, ACTION_ORDER ASC;
| 	
: string
| omni_sync_corpus
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "VARCHAR" "string"
< column "ACTION_TIMING" "VARCHAR" "string"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "VARCHAR" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| trg_orders_bi
: bytes
| orders
: bytes
| INSERT
: bytes
| BEFORE
: bytes
| BEGIN
| 	IF NEW.total < 0 THEN
| 		SET NEW.total = 0;
| 	END IF;
| END
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
|
| 	SELECT
| 		EVENT_NAME,
| 		TIME_ZONE,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.EVENTS
| 	WHERE EVENT_SCHEMA = ?
| 	ORDER BY EVENT_NAME ASC;
| 	
: string
| omni_sync_corpus
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "VARCHAR" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| SYSTEM
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE EVENT `omni_sync_corpus`.`ev_daily_check`
< result
< column "Event" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "time_zone" "VARCHAR" "string"
< column "Create Event" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_daily_check` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO BEGIN
| 	DELETE FROM orders WHERE total < 0;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			REFERENCED_TABLE_NAME,
| 			DELETE_RULE,
| 			UPDATE_RULE,
| 			MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ?;
| 	
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR" "string"
< column "UPDATE_RULE" "VARCHAR" "string"
< column "MATCH_OPTION" "VARCHAR" "string"
< row
: bytes
| history
: bytes
| fk_history_account
: bytes
| accounts
: bytes
| RESTRICT
: bytes
| RESTRICT
: bytes
| NONE
< row
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| accounts
: bytes
| CASCADE
: bytes
| RESTRICT
: bytes
| NONE
> query
|
| 		SELECT
| 			TABLE_SCHEMA,
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			COLUMN_NAME,
| 			REFERENCED_TABLE_SCHEMA,
| 			REFERENCED_COLUMN_NAME
| 		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
| 		WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION;
| 	
: string
| omni_sync_corpus
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_sync_corpus
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| omni_sync_corpus
: bytes
| id
< row
: bytes
| omni_sync_corpus
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_sync_corpus
: bytes
| id
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "SUBPARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| metrics_hash
: bytes
| p0
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
< row
: bytes
| metrics_hash
: bytes
| p1
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
< row
: bytes
| metrics_hash
: bytes
| p2
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
< row
: bytes
| metrics_hash
: bytes
| p3
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
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
| SHOW CREATE TABLE `omni_sync_corpus`.`metrics_hash`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_hash
: bytes
| CREATE TABLE `metrics_hash` (
|   `id` int(11) NOT NULL,
|   `val` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
|  PARTITION BY HASH (`id`)
| PARTITIONS 4
> query
| SHOW CREATE TABLE `omni_sync_corpus`.`metrics_range`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_range
: bytes
| CREATE TABLE `metrics_range` (
|   `ts` int(11) NOT NULL,
|   `val` bigint(20) NOT NULL,
|   PRIMARY KEY (`ts`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
|  PARTITION BY RANGE (`ts`)
| (PARTITION `p_low` VALUES LESS THAN (100) ENGINE = InnoDB,
|  PARTITION `p_high` VALUES LESS THAN MAXVALUE ENGINE = InnoDB)
> query
|
| 		SELECT
| 			ROUTINE_NAME,
| 			ROUTINE_TYPE,
| 			SQL_MODE,
| 			CHARACTER_SET_CLIENT,
| 			COLLATION_CONNECTION,
| 			DATABASE_COLLATION,
| 			IFNULL(ROUTINE_COMMENT, '')
| 		FROM
| 			INFORMATION_SCHEMA.ROUTINES
| 		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE IN ('FUNCTION', 'PROCEDURE')
| 		ORDER BY ROUTINE_TYPE, ROUTINE_NAME;
| 	
: string
| omni_sync_corpus
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "VARCHAR" "string"
< column "SQL_MODE" "VARCHAR" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "LONGTEXT" "string"
< row
: bytes
| fn_format_tag
: bytes
| FUNCTION
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
: bytes
| formats a tag
< row
: bytes
| sp_touch_order
: bytes
| PROCEDURE
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
: bytes
| reads order total
> query
| SHOW CREATE FUNCTION `omni_sync_corpus`.`fn_format_tag`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| fn_format_tag
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci
|     NO SQL
|     DETERMINISTIC
|     COMMENT 'formats a tag'
| BEGIN
| 	RETURN CONCAT(p_prefix, '-', p_id);
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE PROCEDURE `omni_sync_corpus`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| sp_touch_order
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `sp_touch_order`(IN p_id BIGINT, OUT p_total INT)
|     COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> query
|
| 		SELECT
| 			TABLES.TABLE_NAME,
| 			TABLES.TABLE_TYPE,
| 			IFNULL(TABLES.ENGINE, ''),
| 			IFNULL(TABLES.TABLE_COLLATION, ''),
| 			IFNULL(TABLES.TABLE_ROWS, 0),
| 			IFNULL(TABLES.DATA_LENGTH, 0),
| 			IFNULL(TABLES.INDEX_LENGTH, 0),
| 			IFNULL(TABLES.DATA_FREE, 0),
| 			IFNULL(TABLES.CREATE_OPTIONS, ''),
| 			QUOTE(IFNULL(TABLES.TABLE_COMMENT, '')),
| 			IFNULL(CCSA.CHARACTER_SET_NAME, '')
| 		FROM information_schema.TABLES TABLES
| 		LEFT JOIN information_schema.COLLATION_CHARACTER_SET_APPLICABILITY CCSA
| 		ON TABLES.TABLE_COLLATION = CCSA.COLLATION_NAME
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "TABLE_TYPE" "VARCHAR" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "TEXT" "string"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| 'customer accounts'
: bytes
| utf8mb4
< row
: bytes
| explicit_period
: bytes
| SYSTEM VERSIONED
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| history
: bytes
| SYSTEM VERSIONED
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| 'versioned'
: bytes
| utf8mb4
< row
: bytes
| metrics_hash
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 4
: bytes
| 65536
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| metrics_range
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 2
: bytes
| 32768
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| orders
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| seq_custom
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_general_ci
: bytes
| 1
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'sequence comment'
: bytes
| utf8mb4
< row
: bytes
| seq_default
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_general_ci
: bytes
| 1
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| v_named_cols
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| v_named_cols_ref
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
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
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_sync_corpus`.`seq_custom`
< result
< column "start_value" "BIGINT" "int64"
< column "minimum_value" "BIGINT" "int64"
< column "maximum_value" "BIGINT" "int64"
< column "increment" "BIGINT" "int64"
< column "cache_size" "UNSIGNED BIGINT" "uint64"
< column "cycle_option" "UNSIGNED TINYINT" "uint8"
< row
: int64 100
: int64 10
: int64 1000
: int64 5
: uint64 20
: int64 1
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_sync_corpus`.`seq_default`
< result
< column "start_value" "BIGINT" "int64"
< column "minimum_value" "BIGINT" "int64"
< column "maximum_value" "BIGINT" "int64"
< column "increment" "BIGINT" "int64"
< column "cache_size" "UNSIGNED BIGINT" "uint64"
< column "cycle_option" "UNSIGNED TINYINT" "uint8"
< row
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: uint64 1000
: int64 0
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
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestSyncDatabase/mariadb107 mariadb107#1
> exec
| CREATE TABLE accounts (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	email VARCHAR(190) NOT NULL,
| 	bio TEXT,
| 	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
| 	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
| 	secret_code VARCHAR(32) INVISIBLE DEFAULT 'hidden',
| 	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL,
| 	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED,
| 	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
| 	note VARCHAR(100) COMMENT 'owner''s note',
| 	UNIQUE KEY uk_email (email),
| 	KEY idx_bio_prefix (bio(32)),
| 	CONSTRAINT chk_balance CHECK (balance >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='customer accounts';
< ok 0 0
> exec
| CREATE TABLE orders (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	account_id BIGINT NOT NULL,
| 	total INT NOT NULL,
| 	KEY idx_account (account_id),
| 	CONSTRAINT fk_orders_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
< ok 0 0
> exec
| CREATE TABLE history (
| 	id INT PRIMARY KEY,
| 	amount INT,
| 	account_id BIGINT,
| 	KEY idx_amount (amount),
| 	CONSTRAINT fk_history_account FOREIGN KEY (account_id) REFERENCES accounts (id),
| 	CONSTRAINT chk_amount CHECK (amount > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING COMMENT='versioned';
< ok 0 0
> exec
| CREATE TABLE explicit_period (
| 	id INT PRIMARY KEY,
| 	valid_from TIMESTAMP(6) GENERATED ALWAYS AS ROW START,
| 	valid_to TIMESTAMP(6) GENERATED ALWAYS AS ROW END,
| 	PERIOD FOR SYSTEM_TIME (valid_from, valid_to)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING;
< ok 0 0
> exec
| CREATE TABLE metrics_range (
| 	ts INT NOT NULL,
| 	val BIGINT NOT NULL,
| 	PRIMARY KEY (ts)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY RANGE (ts) (
| 	PARTITION p_low VALUES LESS THAN (100),
| 	PARTITION p_high VALUES LESS THAN MAXVALUE
| );
< ok 0 0
> exec
| CREATE TABLE metrics_hash (
| 	id INT NOT NULL PRIMARY KEY,
| 	val INT NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY HASH (id) PARTITIONS 4;
< ok 0 0
> exec
| CREATE SEQUENCE seq_default;
< ok 0 0
> exec
| CREATE SEQUENCE seq_custom START WITH 100 INCREMENT BY 5 MINVALUE 10 MAXVALUE 1000 CACHE 20 CYCLE COMMENT='sequence comment';
< ok 0 0
> exec
| CREATE VIEW v_open_accounts AS
| 	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open';
< ok 0 0
> exec
| CREATE VIEW v_named_cols (`renamed_id`) AS SELECT id FROM accounts;
< ok 0 0
> exec
| CREATE VIEW v_named_cols_ref AS
| 	SELECT renamed_id FROM v_named_cols;
< ok 0 0
> exec
| CREATE FUNCTION fn_format_tag(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT)
| RETURNS VARCHAR(64) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci
| DETERMINISTIC
| NO SQL
| COMMENT 'formats a tag'
| BEGIN
| 	RETURN CONCAT(p_prefix, '-', p_id);
| END;
< ok 0 0
> exec
| CREATE PROCEDURE sp_touch_order(IN p_id BIGINT, OUT p_total INT)
| COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END;
< ok 0 0
> exec
| CREATE TRIGGER trg_orders_bi
| BEFORE INSERT ON orders
| FOR EACH ROW
| BEGIN
| 	IF NEW.total < 0 THEN
| 		SET NEW.total = 0;
| 	END IF;
| END;
< ok 0 0
> exec
| CREATE EVENT ev_daily_check
| ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00'
| DISABLE
| DO
| BEGIN
| 	DELETE FROM orders WHERE total < 0;
| END;
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR" "string"
< row
: bytes
| 10.7.8-MariaDB-1:10.7.8+maria~ubu2004
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			INDEX_NAME,
| 			COLUMN_NAME,
| 			COLLATION,
| 			IFNULL(SUB_PART, -1),
| 			SEQ_IN_INDEX,
| 			INDEX_TYPE,
| 			NON_UNIQUE,
| 			INDEX_COMMENT
| 		FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "SEQ_IN_INDEX" "BIGINT" "int64"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "BIGINT" "int64"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| idx_bio_prefix
: bytes
| bio
: bytes
| A
: int64 32
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| accounts
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| accounts
: bytes
| uk_email
: bytes
| email
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| explicit_period
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| explicit_period
: bytes
| PRIMARY
: bytes
| valid_to
: bytes
| A
: int64 -1
: int64 2
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| history
: bytes
| idx_amount
: bytes
| amount
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| history
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| metrics_hash
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| metrics_range
: bytes
| PRIMARY
: bytes
| ts
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| orders
: bytes
| idx_account
: bytes
| account_id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| orders
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
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
| 			convert(GENERATION_EXPRESSION using BINARY),
| 			EXTRA
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT" "uint64"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "MEDIUMTEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "LONGBLOB" "[]uint8"
< column "EXTRA" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| id
: uint64 1
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
| ''
: null
: bytes
| auto_increment
< row
: bytes
| accounts
: bytes
| email
: uint64 2
: null
: bytes
| NO
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| bio
: uint64 3
: bytes
| 'NULL'
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| status
: uint64 4
: bytes
| '\'open\''
: bytes
| NO
: bytes
| enum('open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| balance
: uint64 5
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
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| secret_code
: uint64 6
: bytes
| '\'hidden\''
: bytes
| YES
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
| INVISIBLE
< row
: bytes
| accounts
: bytes
| email_lower
: uint64 7
: bytes
| 'NULL'
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
| lcase(`email`)
: bytes
| VIRTUAL GENERATED
< row
: bytes
| accounts
: bytes
| display_tag
: uint64 8
: bytes
| 'NULL'
: bytes
| YES
: bytes
| varchar(220)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
| concat(`email`,':','v1')
: bytes
| STORED GENERATED
< row
: bytes
| accounts
: bytes
| created_at
: uint64 9
: bytes
| 'current_timestamp()'
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
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| updated_at
: uint64 10
: bytes
| 'current_timestamp(6)'
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
: null
: bytes
| on update current_timestamp(6)
< row
: bytes
| accounts
: bytes
| note
: uint64 11
: bytes
| 'NULL'
: bytes
| YES
: bytes
| varchar(100)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| 'owner\'s note'
: null
: bytes
|
< row
: bytes
| explicit_period
: bytes
| id
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| explicit_period
: bytes
| valid_from
: uint64 2
: null
: bytes
| NO
: bytes
| timestamp(6)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
| ROW START
: bytes
| STORED GENERATED
< row
: bytes
| explicit_period
: bytes
| valid_to
: uint64 3
: null
: bytes
| NO
: bytes
| timestamp(6)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
| ROW END
: bytes
| STORED GENERATED
< row
: bytes
| history
: bytes
| id
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| history
: bytes
| amount
: uint64 2
: bytes
| 'NULL'
: bytes
| YES
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| history
: bytes
| account_id
: uint64 3
: bytes
| 'NULL'
: bytes
| YES
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_hash
: bytes
| id
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_hash
: bytes
| val
: uint64 2
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_range
: bytes
| ts
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_range
: bytes
| val
: uint64 2
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
| ''
: null
: bytes
|
< row
: bytes
| orders
: bytes
| id
: uint64 1
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
| ''
: null
: bytes
| auto_increment
< row
: bytes
| orders
: bytes
| account_id
: uint64 2
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
| ''
: null
: bytes
|
< row
: bytes
| orders
: bytes
| total
: uint64 3
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| next_not_cached_value
: uint64 1
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| minimum_value
: uint64 2
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| maximum_value
: uint64 3
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| start_value
: uint64 4
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'start value when sequences is created or value if RESTART is used'
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| increment
: uint64 5
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'increment value'
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| cache_size
: uint64 6
: null
: bytes
| NO
: bytes
| bigint(21) unsigned
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| cycle_option
: uint64 7
: null
: bytes
| NO
: bytes
| tinyint(1) unsigned
: bytes
|
: bytes
|
: bytes
| '0 if no cycles are allowed, 1 if the sequence should begin a new cycle when maximum_value is passed'
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| cycle_count
: uint64 8
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'How many cycles have been done'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| next_not_cached_value
: uint64 1
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| minimum_value
: uint64 2
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| maximum_value
: uint64 3
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| start_value
: uint64 4
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'start value when sequences is created or value if RESTART is used'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| increment
: uint64 5
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'increment value'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| cache_size
: uint64 6
: null
: bytes
| NO
: bytes
| bigint(21) unsigned
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| cycle_option
: uint64 7
: null
: bytes
| NO
: bytes
| tinyint(1) unsigned
: bytes
|
: bytes
|
: bytes
| '0 if no cycles are allowed, 1 if the sequence should begin a new cycle when maximum_value is passed'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| cycle_count
: uint64 8
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'How many cycles have been done'
: null
: bytes
|
< row
: bytes
| v_named_cols
: bytes
| renamed_id
: uint64 1
: bytes
| '0'
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| v_named_cols_ref
: bytes
| renamed_id
: uint64 1
: bytes
| '0'
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| v_open_accounts
: bytes
| acct_id
: uint64 1
: bytes
| '0'
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| v_open_accounts
: bytes
| acct_email
: uint64 2
: null
: bytes
| NO
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				CONSTRAINT_NAME,
| 				CHECK_CLAUSE
| 			FROM information_schema.CHECK_CONSTRAINTS
| 			WHERE CONSTRAINT_SCHEMA = ?
| 			ORDER BY TABLE_NAME, CONSTRAINT_NAME
| 		
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| accounts
: bytes
| chk_balance
: bytes
| `balance` >= 0
< row
: bytes
| history
: bytes
| chk_amount
: bytes
| `amount` > 0
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v_named_cols
: bytes
| select `omni_sync_corpus`.`accounts`.`id` AS `renamed_id` from `omni_sync_corpus`.`accounts`
< row
: bytes
| v_named_cols_ref
: bytes
| select `v_named_cols`.`renamed_id` AS `renamed_id` from `omni_sync_corpus`.`v_named_cols`
< row
: bytes
| v_open_accounts
: bytes
| select `omni_sync_corpus`.`accounts`.`id` AS `acct_id`,`omni_sync_corpus`.`accounts`.`email` AS `acct_email` from `omni_sync_corpus`.`accounts` where `omni_sync_corpus`.`accounts`.`status` = 'open'
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_named_cols`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols` AS select `accounts`.`id` AS `renamed_id` from `accounts`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_named_cols_ref`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols_ref
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols_ref` AS select `v_named_cols`.`renamed_id` AS `renamed_id` from `v_named_cols`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_open_accounts`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_open_accounts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where `accounts`.`status` = 'open'
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
|
| 	SELECT
| 		TRIGGER_NAME,
| 		EVENT_OBJECT_TABLE,
| 		EVENT_MANIPULATION,
| 		ACTION_TIMING,
| 		ACTION_STATEMENT,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.TRIGGERS
| 	WHERE EVENT_OBJECT_SCHEMA = ?
| 	ORDER BY EVENT_OBJECT_TABLE ASC, EVENT_MANIPULATION ASC, ACTION_TIMING ASC, ACTION_ORDER ASC;
| 	
: string
| omni_sync_corpus
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "VARCHAR" "string"
< column "ACTION_TIMING" "VARCHAR" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "VARCHAR" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| trg_orders_bi
: bytes
| orders
: bytes
| INSERT
: bytes
| BEFORE
: bytes
| BEGIN
| 	IF NEW.total < 0 THEN
| 		SET NEW.total = 0;
| 	END IF;
| END
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
|
| 	SELECT
| 		EVENT_NAME,
| 		TIME_ZONE,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.EVENTS
| 	WHERE EVENT_SCHEMA = ?
| 	ORDER BY EVENT_NAME ASC;
| 	
: string
| omni_sync_corpus
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "VARCHAR" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| SYSTEM
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE EVENT `omni_sync_corpus`.`ev_daily_check`
< result
< column "Event" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "time_zone" "VARCHAR" "string"
< column "Create Event" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_daily_check` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO BEGIN
| 	DELETE FROM orders WHERE total < 0;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			REFERENCED_TABLE_NAME,
| 			DELETE_RULE,
| 			UPDATE_RULE,
| 			MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ?;
| 	
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "VARCHAR" "string"
< column "UPDATE_RULE" "VARCHAR" "string"
< column "MATCH_OPTION" "VARCHAR" "string"
< row
: bytes
| history
: bytes
| fk_history_account
: bytes
| accounts
: bytes
| RESTRICT
: bytes
| RESTRICT
: bytes
| NONE
< row
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| accounts
: bytes
| CASCADE
: bytes
| RESTRICT
: bytes
| NONE
> query
|
| 		SELECT
| 			TABLE_SCHEMA,
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			COLUMN_NAME,
| 			REFERENCED_TABLE_SCHEMA,
| 			REFERENCED_COLUMN_NAME
| 		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
| 		WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION;
| 	
: string
| omni_sync_corpus
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_sync_corpus
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| omni_sync_corpus
: bytes
| id
< row
: bytes
| omni_sync_corpus
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_sync_corpus
: bytes
| id
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "SUBPARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| metrics_hash
: bytes
| p0
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
< row
: bytes
| metrics_hash
: bytes
| p1
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
< row
: bytes
| metrics_hash
: bytes
| p2
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
< row
: bytes
| metrics_hash
: bytes
| p3
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
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
| SHOW CREATE TABLE `omni_sync_corpus`.`metrics_hash`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_hash
: bytes
| CREATE TABLE `metrics_hash` (
|   `id` int(11) NOT NULL,
|   `val` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
|  PARTITION BY HASH (`id`)
| PARTITIONS 4
> query
| SHOW CREATE TABLE `omni_sync_corpus`.`metrics_range`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_range
: bytes
| CREATE TABLE `metrics_range` (
|   `ts` int(11) NOT NULL,
|   `val` bigint(20) NOT NULL,
|   PRIMARY KEY (`ts`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
|  PARTITION BY RANGE (`ts`)
| (PARTITION `p_low` VALUES LESS THAN (100) ENGINE = InnoDB,
|  PARTITION `p_high` VALUES LESS THAN MAXVALUE ENGINE = InnoDB)
> query
|
| 		SELECT
| 			ROUTINE_NAME,
| 			ROUTINE_TYPE,
| 			SQL_MODE,
| 			CHARACTER_SET_CLIENT,
| 			COLLATION_CONNECTION,
| 			DATABASE_COLLATION,
| 			IFNULL(ROUTINE_COMMENT, '')
| 		FROM
| 			INFORMATION_SCHEMA.ROUTINES
| 		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE IN ('FUNCTION', 'PROCEDURE')
| 		ORDER BY ROUTINE_TYPE, ROUTINE_NAME;
| 	
: string
| omni_sync_corpus
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "VARCHAR" "string"
< column "SQL_MODE" "VARCHAR" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "LONGTEXT" "string"
< row
: bytes
| fn_format_tag
: bytes
| FUNCTION
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
: bytes
| formats a tag
< row
: bytes
| sp_touch_order
: bytes
| PROCEDURE
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
: bytes
| reads order total
> query
| SHOW CREATE FUNCTION `omni_sync_corpus`.`fn_format_tag`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| fn_format_tag
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci
|     NO SQL
|     DETERMINISTIC
|     COMMENT 'formats a tag'
| BEGIN
| 	RETURN CONCAT(p_prefix, '-', p_id);
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE PROCEDURE `omni_sync_corpus`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| sp_touch_order
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `sp_touch_order`(IN p_id BIGINT, OUT p_total INT)
|     COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> query
|
| 		SELECT
| 			TABLES.TABLE_NAME,
| 			TABLES.TABLE_TYPE,
| 			IFNULL(TABLES.ENGINE, ''),
| 			IFNULL(TABLES.TABLE_COLLATION, ''),
| 			IFNULL(TABLES.TABLE_ROWS, 0),
| 			IFNULL(TABLES.DATA_LENGTH, 0),
| 			IFNULL(TABLES.INDEX_LENGTH, 0),
| 			IFNULL(TABLES.DATA_FREE, 0),
| 			IFNULL(TABLES.CREATE_OPTIONS, ''),
| 			QUOTE(IFNULL(TABLES.TABLE_COMMENT, '')),
| 			IFNULL(CCSA.CHARACTER_SET_NAME, '')
| 		FROM information_schema.TABLES TABLES
| 		LEFT JOIN information_schema.COLLATION_CHARACTER_SET_APPLICABILITY CCSA
| 		ON TABLES.TABLE_COLLATION = CCSA.COLLATION_NAME
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "TABLE_TYPE" "VARCHAR" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "TEXT" "string"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| 'customer accounts'
: bytes
| utf8mb4
< row
: bytes
| explicit_period
: bytes
| SYSTEM VERSIONED
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| history
: bytes
| SYSTEM VERSIONED
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| 'versioned'
: bytes
| utf8mb4
< row
: bytes
| metrics_hash
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 4
: bytes
| 65536
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| metrics_range
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 2
: bytes
| 32768
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| orders
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| seq_custom
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_general_ci
: bytes
| 1
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'sequence comment'
: bytes
| utf8mb4
< row
: bytes
| seq_default
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_general_ci
: bytes
| 1
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| v_named_cols
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| v_named_cols_ref
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
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
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_sync_corpus`.`seq_custom`
< result
< column "start_value" "BIGINT" "int64"
< column "minimum_value" "BIGINT" "int64"
< column "maximum_value" "BIGINT" "int64"
< column "increment" "BIGINT" "int64"
< column "cache_size" "UNSIGNED BIGINT" "uint64"
< column "cycle_option" "UNSIGNED TINYINT" "uint8"
< row
: int64 100
: int64 10
: int64 1000
: int64 5
: uint64 20
: int64 1
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_sync_corpus`.`seq_default`
< result
< column "start_value" "BIGINT" "int64"
< column "minimum_value" "BIGINT" "int64"
< column "maximum_value" "BIGINT" "int64"
< column "increment" "BIGINT" "int64"
< column "cache_size" "UNSIGNED BIGINT" "uint64"
< column "cycle_option" "UNSIGNED TINYINT" "uint8"
< row
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: uint64 1000
: int64 0
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
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestSyncDatabase/mariadb114 mariadb114#1
> exec
| CREATE TABLE accounts (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	email VARCHAR(190) NOT NULL,
| 	bio TEXT,
| 	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
| 	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
| 	secret_code VARCHAR(32) INVISIBLE DEFAULT 'hidden',
| 	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL,
| 	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED,
| 	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
| 	note VARCHAR(100) COMMENT 'owner''s note',
| 	UNIQUE KEY uk_email (email),
| 	KEY idx_bio_prefix (bio(32)),
| 	CONSTRAINT chk_balance CHECK (balance >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='customer accounts';
< ok 0 0
> exec
| CREATE TABLE orders (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	account_id BIGINT NOT NULL,
| 	total INT NOT NULL,
| 	KEY idx_account (account_id),
| 	CONSTRAINT fk_orders_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
< ok 0 0
> exec
| CREATE TABLE history (
| 	id INT PRIMARY KEY,
| 	amount INT,
| 	account_id BIGINT,
| 	KEY idx_amount (amount),
| 	CONSTRAINT fk_history_account FOREIGN KEY (account_id) REFERENCES accounts (id),
| 	CONSTRAINT chk_amount CHECK (amount > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING COMMENT='versioned';
< ok 0 0
> exec
| CREATE TABLE explicit_period (
| 	id INT PRIMARY KEY,
| 	valid_from TIMESTAMP(6) GENERATED ALWAYS AS ROW START,
| 	valid_to TIMESTAMP(6) GENERATED ALWAYS AS ROW END,
| 	PERIOD FOR SYSTEM_TIME (valid_from, valid_to)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING;
< ok 0 0
> exec
| CREATE TABLE metrics_range (
| 	ts INT NOT NULL,
| 	val BIGINT NOT NULL,
| 	PRIMARY KEY (ts)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY RANGE (ts) (
| 	PARTITION p_low VALUES LESS THAN (100),
| 	PARTITION p_high VALUES LESS THAN MAXVALUE
| );
< ok 0 0
> exec
| CREATE TABLE metrics_hash (
| 	id INT NOT NULL PRIMARY KEY,
| 	val INT NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY HASH (id) PARTITIONS 4;
< ok 0 0
> exec
| CREATE SEQUENCE seq_default;
< ok 0 0
> exec
| CREATE SEQUENCE seq_custom START WITH 100 INCREMENT BY 5 MINVALUE 10 MAXVALUE 1000 CACHE 20 CYCLE COMMENT='sequence comment';
< ok 0 0
> exec
| CREATE VIEW v_open_accounts AS
| 	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open';
< ok 0 0
> exec
| CREATE VIEW v_named_cols (`renamed_id`) AS SELECT id FROM accounts;
< ok 0 0
> exec
| CREATE VIEW v_named_cols_ref AS
| 	SELECT renamed_id FROM v_named_cols;
< ok 0 0
> exec
| CREATE FUNCTION fn_format_tag(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT)
| RETURNS VARCHAR(64) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci
| DETERMINISTIC
| NO SQL
| COMMENT 'formats a tag'
| BEGIN
| 	RETURN CONCAT(p_prefix, '-', p_id);
| END;
< ok 0 0
> exec
| CREATE PROCEDURE sp_touch_order(IN p_id BIGINT, OUT p_total INT)
| COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END;
< ok 0 0
> exec
| CREATE TRIGGER trg_orders_bi
| BEFORE INSERT ON orders
| FOR EACH ROW
| BEGIN
| 	IF NEW.total < 0 THEN
| 		SET NEW.total = 0;
| 	END IF;
| END;
< ok 0 0
> exec
| CREATE EVENT ev_daily_check
| ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00'
| DISABLE
| DO
| BEGIN
| 	DELETE FROM orders WHERE total < 0;
| END;
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR" "string"
< row
: bytes
| 11.4.13-MariaDB-ubu2404
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			INDEX_NAME,
| 			COLUMN_NAME,
| 			COLLATION,
| 			IFNULL(SUB_PART, -1),
| 			SEQ_IN_INDEX,
| 			INDEX_TYPE,
| 			NON_UNIQUE,
| 			INDEX_COMMENT
| 		FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "BIGINT" "int64"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| idx_bio_prefix
: bytes
| bio
: bytes
| A
: int64 32
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| accounts
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| accounts
: bytes
| uk_email
: bytes
| email
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| explicit_period
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| explicit_period
: bytes
| PRIMARY
: bytes
| valid_to
: bytes
| A
: int64 -1
: int64 2
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| history
: bytes
| idx_amount
: bytes
| amount
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| history
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| metrics_hash
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| metrics_range
: bytes
| PRIMARY
: bytes
| ts
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
< row
: bytes
| orders
: bytes
| idx_account
: bytes
| account_id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 1
: bytes
|
< row
: bytes
| orders
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: int64 1
: bytes
| BTREE
: int64 0
: bytes
|
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
| 			convert(GENERATION_EXPRESSION using BINARY),
| 			EXTRA
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT" "uint64"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "MEDIUMTEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "LONGBLOB" "[]uint8"
< column "EXTRA" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| id
: uint64 1
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
| ''
: null
: bytes
| auto_increment
< row
: bytes
| accounts
: bytes
| email
: uint64 2
: null
: bytes
| NO
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| bio
: uint64 3
: bytes
| 'NULL'
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| status
: uint64 4
: bytes
| '\'open\''
: bytes
| NO
: bytes
| enum('open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| balance
: uint64 5
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
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| secret_code
: uint64 6
: bytes
| '\'hidden\''
: bytes
| YES
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
| INVISIBLE
< row
: bytes
| accounts
: bytes
| email_lower
: uint64 7
: bytes
| 'NULL'
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
| lcase(`email`)
: bytes
| VIRTUAL GENERATED
< row
: bytes
| accounts
: bytes
| display_tag
: uint64 8
: bytes
| 'NULL'
: bytes
| YES
: bytes
| varchar(220)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
| concat(`email`,':','v1')
: bytes
| STORED GENERATED
< row
: bytes
| accounts
: bytes
| created_at
: uint64 9
: bytes
| 'current_timestamp()'
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
: null
: bytes
|
< row
: bytes
| accounts
: bytes
| updated_at
: uint64 10
: bytes
| 'current_timestamp(6)'
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
: null
: bytes
| on update current_timestamp(6)
< row
: bytes
| accounts
: bytes
| note
: uint64 11
: bytes
| 'NULL'
: bytes
| YES
: bytes
| varchar(100)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| 'owner\'s note'
: null
: bytes
|
< row
: bytes
| explicit_period
: bytes
| id
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| explicit_period
: bytes
| valid_from
: uint64 2
: null
: bytes
| NO
: bytes
| timestamp(6)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
| ROW START
: bytes
| STORED GENERATED
< row
: bytes
| explicit_period
: bytes
| valid_to
: uint64 3
: null
: bytes
| NO
: bytes
| timestamp(6)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
| ROW END
: bytes
| STORED GENERATED
< row
: bytes
| history
: bytes
| id
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| history
: bytes
| amount
: uint64 2
: bytes
| 'NULL'
: bytes
| YES
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| history
: bytes
| account_id
: uint64 3
: bytes
| 'NULL'
: bytes
| YES
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_hash
: bytes
| id
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_hash
: bytes
| val
: uint64 2
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_range
: bytes
| ts
: uint64 1
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| metrics_range
: bytes
| val
: uint64 2
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
| ''
: null
: bytes
|
< row
: bytes
| orders
: bytes
| id
: uint64 1
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
| ''
: null
: bytes
| auto_increment
< row
: bytes
| orders
: bytes
| account_id
: uint64 2
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
| ''
: null
: bytes
|
< row
: bytes
| orders
: bytes
| total
: uint64 3
: null
: bytes
| NO
: bytes
| int(11)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| next_not_cached_value
: uint64 1
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| minimum_value
: uint64 2
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| maximum_value
: uint64 3
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| start_value
: uint64 4
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'start value when sequences is created or value if RESTART is used'
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| increment
: uint64 5
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'increment value'
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| cache_size
: uint64 6
: null
: bytes
| NO
: bytes
| bigint(21) unsigned
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| cycle_option
: uint64 7
: null
: bytes
| NO
: bytes
| tinyint(1) unsigned
: bytes
|
: bytes
|
: bytes
| '0 if no cycles are allowed, 1 if the sequence should begin a new cycle when maximum_value is passed'
: null
: bytes
|
< row
: bytes
| seq_custom
: bytes
| cycle_count
: uint64 8
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'How many cycles have been done'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| next_not_cached_value
: uint64 1
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| minimum_value
: uint64 2
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| maximum_value
: uint64 3
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| start_value
: uint64 4
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'start value when sequences is created or value if RESTART is used'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| increment
: uint64 5
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'increment value'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| cache_size
: uint64 6
: null
: bytes
| NO
: bytes
| bigint(21) unsigned
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| cycle_option
: uint64 7
: null
: bytes
| NO
: bytes
| tinyint(1) unsigned
: bytes
|
: bytes
|
: bytes
| '0 if no cycles are allowed, 1 if the sequence should begin a new cycle when maximum_value is passed'
: null
: bytes
|
< row
: bytes
| seq_default
: bytes
| cycle_count
: uint64 8
: null
: bytes
| NO
: bytes
| bigint(21)
: bytes
|
: bytes
|
: bytes
| 'How many cycles have been done'
: null
: bytes
|
< row
: bytes
| v_named_cols
: bytes
| renamed_id
: uint64 1
: bytes
| '0'
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| v_named_cols_ref
: bytes
| renamed_id
: uint64 1
: bytes
| '0'
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| v_open_accounts
: bytes
| acct_id
: uint64 1
: bytes
| '0'
: bytes
| NO
: bytes
| bigint(20)
: bytes
|
: bytes
|
: bytes
| ''
: null
: bytes
|
< row
: bytes
| v_open_accounts
: bytes
| acct_email
: uint64 2
: null
: bytes
| NO
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: null
: bytes
|
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				CONSTRAINT_NAME,
| 				CHECK_CLAUSE
| 			FROM information_schema.CHECK_CONSTRAINTS
| 			WHERE CONSTRAINT_SCHEMA = ?
| 			ORDER BY TABLE_NAME, CONSTRAINT_NAME
| 		
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| accounts
: bytes
| chk_balance
: bytes
| `balance` >= 0
< row
: bytes
| history
: bytes
| chk_amount
: bytes
| `amount` > 0
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v_named_cols
: bytes
| select `omni_sync_corpus`.`accounts`.`id` AS `renamed_id` from `omni_sync_corpus`.`accounts`
< row
: bytes
| v_named_cols_ref
: bytes
| select `v_named_cols`.`renamed_id` AS `renamed_id` from `omni_sync_corpus`.`v_named_cols`
< row
: bytes
| v_open_accounts
: bytes
| select `omni_sync_corpus`.`accounts`.`id` AS `acct_id`,`omni_sync_corpus`.`accounts`.`email` AS `acct_email` from `omni_sync_corpus`.`accounts` where `omni_sync_corpus`.`accounts`.`status` = 'open'
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_named_cols`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols` AS select `accounts`.`id` AS `renamed_id` from `accounts`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_named_cols_ref`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols_ref
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols_ref` AS select `v_named_cols`.`renamed_id` AS `renamed_id` from `v_named_cols`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_open_accounts`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_open_accounts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where `accounts`.`status` = 'open'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
|
| 	SELECT
| 		TRIGGER_NAME,
| 		EVENT_OBJECT_TABLE,
| 		EVENT_MANIPULATION,
| 		ACTION_TIMING,
| 		ACTION_STATEMENT,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.TRIGGERS
| 	WHERE EVENT_OBJECT_SCHEMA = ?
| 	ORDER BY EVENT_OBJECT_TABLE ASC, EVENT_MANIPULATION ASC, ACTION_TIMING ASC, ACTION_ORDER ASC;
| 	
: string
| omni_sync_corpus
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "VARCHAR" "string"
< column "ACTION_TIMING" "VARCHAR" "string"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "VARCHAR" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| trg_orders_bi
: bytes
| orders
: bytes
| INSERT
: bytes
| BEFORE
: bytes
| BEGIN
| 	IF NEW.total < 0 THEN
| 		SET NEW.total = 0;
| 	END IF;
| END
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
|
| 	SELECT
| 		EVENT_NAME,
| 		TIME_ZONE,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.EVENTS
| 	WHERE EVENT_SCHEMA = ?
| 	ORDER BY EVENT_NAME ASC;
| 	
: string
| omni_sync_corpus
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "VARCHAR" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| SYSTEM
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE EVENT `omni_sync_corpus`.`ev_daily_check`
< result
< column "Event" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "time_zone" "VARCHAR" "string"
< column "Create Event" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_daily_check` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO BEGIN
| 	DELETE FROM orders WHERE total < 0;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			REFERENCED_TABLE_NAME,
| 			DELETE_RULE,
| 			UPDATE_RULE,
| 			MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ?;
| 	
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR" "string"
< column "UPDATE_RULE" "VARCHAR" "string"
< column "MATCH_OPTION" "VARCHAR" "string"
< row
: bytes
| history
: bytes
| fk_history_account
: bytes
| accounts
: bytes
| RESTRICT
: bytes
| RESTRICT
: bytes
| NONE
< row
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| accounts
: bytes
| CASCADE
: bytes
| RESTRICT
: bytes
| NONE
> query
|
| 		SELECT
| 			TABLE_SCHEMA,
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			COLUMN_NAME,
| 			REFERENCED_TABLE_SCHEMA,
| 			REFERENCED_COLUMN_NAME
| 		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
| 		WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION;
| 	
: string
| omni_sync_corpus
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_sync_corpus
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| omni_sync_corpus
: bytes
| id
< row
: bytes
| omni_sync_corpus
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_sync_corpus
: bytes
| id
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "SUBPARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| metrics_hash
: bytes
| p0
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
< row
: bytes
| metrics_hash
: bytes
| p1
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
< row
: bytes
| metrics_hash
: bytes
| p2
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
< row
: bytes
| metrics_hash
: bytes
| p3
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: null
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
| SHOW CREATE TABLE `omni_sync_corpus`.`metrics_hash`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_hash
: bytes
| CREATE TABLE `metrics_hash` (
|   `id` int(11) NOT NULL,
|   `val` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
|  PARTITION BY HASH (`id`)
| PARTITIONS 4
> query
| SHOW CREATE TABLE `omni_sync_corpus`.`metrics_range`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_range
: bytes
| CREATE TABLE `metrics_range` (
|   `ts` int(11) NOT NULL,
|   `val` bigint(20) NOT NULL,
|   PRIMARY KEY (`ts`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
|  PARTITION BY RANGE (`ts`)
| (PARTITION `p_low` VALUES LESS THAN (100) ENGINE = InnoDB,
|  PARTITION `p_high` VALUES LESS THAN MAXVALUE ENGINE = InnoDB)
> query
|
| 		SELECT
| 			ROUTINE_NAME,
| 			ROUTINE_TYPE,
| 			SQL_MODE,
| 			CHARACTER_SET_CLIENT,
| 			COLLATION_CONNECTION,
| 			DATABASE_COLLATION,
| 			IFNULL(ROUTINE_COMMENT, '')
| 		FROM
| 			INFORMATION_SCHEMA.ROUTINES
| 		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE IN ('FUNCTION', 'PROCEDURE')
| 		ORDER BY ROUTINE_TYPE, ROUTINE_NAME;
| 	
: string
| omni_sync_corpus
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "VARCHAR" "string"
< column "SQL_MODE" "VARCHAR" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "LONGTEXT" "string"
< row
: bytes
| fn_format_tag
: bytes
| FUNCTION
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| formats a tag
< row
: bytes
| sp_touch_order
: bytes
| PROCEDURE
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| reads order total
> query
| SHOW CREATE FUNCTION `omni_sync_corpus`.`fn_format_tag`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| fn_format_tag
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci
|     NO SQL
|     DETERMINISTIC
|     COMMENT 'formats a tag'
| BEGIN
| 	RETURN CONCAT(p_prefix, '-', p_id);
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE PROCEDURE `omni_sync_corpus`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| sp_touch_order
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `sp_touch_order`(IN p_id BIGINT, OUT p_total INT)
|     COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> query
|
| 		SELECT
| 			TABLES.TABLE_NAME,
| 			TABLES.TABLE_TYPE,
| 			IFNULL(TABLES.ENGINE, ''),
| 			IFNULL(TABLES.TABLE_COLLATION, ''),
| 			IFNULL(TABLES.TABLE_ROWS, 0),
| 			IFNULL(TABLES.DATA_LENGTH, 0),
| 			IFNULL(TABLES.INDEX_LENGTH, 0),
| 			IFNULL(TABLES.DATA_FREE, 0),
| 			IFNULL(TABLES.CREATE_OPTIONS, ''),
| 			QUOTE(IFNULL(TABLES.TABLE_COMMENT, '')),
| 			IFNULL(CCSA.CHARACTER_SET_NAME, '')
| 		FROM information_schema.TABLES TABLES
| 		LEFT JOIN information_schema.COLLATION_CHARACTER_SET_APPLICABILITY CCSA
| 		ON TABLES.TABLE_COLLATION = CCSA.COLLATION_NAME
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "TABLE_TYPE" "VARCHAR" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "TEXT" "string"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| 'customer accounts'
: bytes
| utf8mb4
< row
: bytes
| explicit_period
: bytes
| SYSTEM VERSIONED
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| history
: bytes
| SYSTEM VERSIONED
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| 'versioned'
: bytes
| utf8mb4
< row
: bytes
| metrics_hash
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 4
: bytes
| 65536
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| metrics_range
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 2
: bytes
| 32768
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| orders
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| seq_custom
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| 1
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'sequence comment'
: bytes
|
< row
: bytes
| seq_default
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| 1
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
|
< row
: bytes
| v_named_cols
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| v_named_cols_ref
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
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
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_sync_corpus`.`seq_custom`
< result
< column "start_value" "BIGINT" "int64"
< column "minimum_value" "BIGINT" "int64"
< column "maximum_value" "BIGINT" "int64"
< column "increment" "BIGINT" "int64"
< column "cache_size" "UNSIGNED BIGINT" "uint64"
< column "cycle_option" "UNSIGNED TINYINT" "uint8"
< row
: int64 100
: int64 10
: int64 1000
: int64 5
: uint64 20
: int64 1
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_sync_corpus`.`seq_default`
< result
< column "start_value" "BIGINT" "int64"
< column "minimum_value" "BIGINT" "int64"
< column "maximum_value" "BIGINT" "int64"
< column "increment" "BIGINT" "int64"
< column "cache_size" "UNSIGNED BIGINT" "uint64"
< column "cycle_option" "UNSIGNED TINYINT" "uint8"
< row
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: uint64 1000
: int64 0
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
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
