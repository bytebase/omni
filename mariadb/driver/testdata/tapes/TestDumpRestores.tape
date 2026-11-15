-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8
-- pin: mariadb1011 mariadb:10.11.19 sha256:7d829180b5b2a9ba74f6242876ab401a0d7dd338a807e522ef92bc9698760f20 mariadb-10.11.19
-- pin: mariadb107 mariadb:10.7.8 sha256:9a48ac9f196f3d4fd6fea2cab59a49df9e7ca459bf14b2f7b85a0e38a5454571 mariadb-10.7.8
-- pin: mariadb114 mariadb:11.4.13 sha256:5e73baf8179b29f115101bb102763a47a017361f3e9792daccba6646c3ab1532 mariadb-11.4.13

== TestDumpRestores/default default#1
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
< column "VERSION()" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "SEQ_IN_INDEX" "UNSIGNED INT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "MEDIUMTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "LONGBLOB"
< column "EXTRA" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_open_accounts
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip`.`accounts` where `omni_dump_roundtrip`.`accounts`.`status` = 'open'
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "TIME_ZONE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_dump_roundtrip
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
: bytes
| id
< row
: bytes
| omni_dump_roundtrip
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
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
| omni_dump_roundtrip
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "LONGTEXT"
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
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "TEXT"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
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
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_custom`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
< row
: int64 100
: int64 10
: int64 1000
: int64 5
: uint64 20
: int64 1
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_default`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
< row
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: uint64 1000
: int64 0
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_typed`
< result
< column "start_value" "UNSIGNED INT"
< column "minimum_value" "UNSIGNED INT"
< column "maximum_value" "UNSIGNED INT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
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
| omni_dump_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE `omni_dump_roundtrip`
< ok 10 0
> exec
| CREATE DATABASE `omni_dump_roundtrip`
< ok 1 0
> exec
| USE `omni_dump_roundtrip`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| SET @OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@@EXPLICIT_DEFAULTS_FOR_TIMESTAMP, EXPLICIT_DEFAULTS_FOR_TIMESTAMP=1
< ok 0 0
> exec
| CREATE VIEW `v_open_accounts` AS SELECT
|   1 AS `acct_id`,
|   1 AS `acct_email`
< ok 0 0
> exec
| CREATE SEQUENCE `seq_custom` start with 100 minvalue 10 maxvalue 1000 increment by 5 cache 20 cycle COMMENT='sequence comment'
< ok 0 0
> exec
| CREATE SEQUENCE `seq_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle
< ok 0 0
> exec
| CREATE SEQUENCE `seq_typed` as int unsigned start with 7 minvalue 1 maxvalue 4000000000 increment by 1 cache 1000 nocycle
< ok 0 0
> exec
| CREATE TABLE `accounts` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `email` varchar(190) NOT NULL,
|   `bio` text DEFAULT NULL,
|   `status` enum('open','closed') NOT NULL DEFAULT 'open',
|   `balance` decimal(12,4) NOT NULL DEFAULT 0.0000,
|   `secret_code` varchar(32) DEFAULT 'hidden' INVISIBLE,
|   `email_lower` varchar(190) GENERATED ALWAYS AS (lcase(`email`)) VIRTUAL,
|   `display_tag` varchar(220) GENERATED ALWAYS AS (concat(`email`,':','v1')) STORED,
|   `created_at` datetime NOT NULL DEFAULT current_timestamp(),
|   `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6) ON UPDATE CURRENT_TIMESTAMP(6),
|   `note` varchar(100) DEFAULT NULL COMMENT 'owner\'s note',
|   PRIMARY KEY (`id`),
|   KEY `idx_bio_prefix` (`bio`(32)),
|   UNIQUE KEY `uk_email` (`email`),
|   CONSTRAINT `chk_balance` CHECK (`balance` >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='customer accounts'
< ok 0 0
> exec
| CREATE TABLE `explicit_period` (
|   `id` int NOT NULL,
|   `valid_from` timestamp(6) GENERATED ALWAYS AS ROW START NOT NULL,
|   `valid_to` timestamp(6) GENERATED ALWAYS AS ROW END NOT NULL,
|   PERIOD FOR SYSTEM_TIME (`valid_from`, `valid_to`),
|   PRIMARY KEY (`id`, `valid_to`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING
< ok 0 0
> exec
| CREATE TABLE `history` (
|   `id` int NOT NULL,
|   `amount` int DEFAULT NULL,
|   `account_id` bigint DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_history_account` (`account_id`),
|   KEY `idx_amount` (`amount`),
|   CONSTRAINT `fk_history_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`),
|   CONSTRAINT `chk_amount` CHECK (`amount` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='versioned' WITH SYSTEM VERSIONING
< ok 0 0
> exec
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY HASH (`id`)
| PARTITIONS 4
< ok 0 0
> exec
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY RANGE (`ts`)
| (PARTITION `p_low` VALUES LESS THAN (100) ENGINE=InnoDB,
|  PARTITION `p_high` VALUES LESS THAN MAXVALUE ENGINE=InnoDB)
< ok 0 0
> exec
| CREATE TABLE `orders` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `account_id` bigint NOT NULL,
|   `total` int NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_account` (`account_id`),
|   CONSTRAINT `fk_orders_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| SET @saved_cs_client = @@character_set_client, @saved_cs_results = @@character_set_results, @saved_col_connection = @@collation_connection, @saved_sql_mode = @@sql_mode
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64)
|     NO SQL
|     DETERMINISTIC
|     COMMENT 'formats a tag'
| BEGIN
| 	RETURN CONCAT(p_prefix, '-', p_id);
| END
< ok 0 0
> exec
| SET character_set_client = @saved_cs_client, character_set_results = @saved_cs_results, collation_connection = @saved_col_connection, sql_mode = @saved_sql_mode
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where `accounts`.`status` = 'open'
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE PROCEDURE `sp_touch_order`(IN p_id BIGINT, OUT p_total INT)
|     COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| SET time_zone = 'SYSTEM'
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` EVENT `ev_daily_check` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO BEGIN
| 	DELETE FROM orders WHERE total < 0;
| END
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE TRIGGER `trg_orders_bi`
| BEFORE INSERT ON `orders` FOR EACH ROW
| BEGIN
| 	IF NEW.total < 0 THEN
| 		SET NEW.total = 0;
| 	END IF;
| END
< ok 0 0
> exec
| SET EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "SEQ_IN_INDEX" "UNSIGNED INT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "MEDIUMTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "LONGBLOB"
< column "EXTRA" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_open_accounts
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip`.`accounts` where `omni_dump_roundtrip`.`accounts`.`status` = 'open'
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "TIME_ZONE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_dump_roundtrip
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
: bytes
| id
< row
: bytes
| omni_dump_roundtrip
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
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
| omni_dump_roundtrip
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "LONGTEXT"
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
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_format_tag
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64) CHARSET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
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
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "TEXT"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
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
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_custom`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
< row
: int64 100
: int64 10
: int64 1000
: int64 5
: uint64 20
: int64 1
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_default`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
< row
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: uint64 1000
: int64 0
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_typed`
< result
< column "start_value" "UNSIGNED INT"
< column "minimum_value" "UNSIGNED INT"
< column "maximum_value" "UNSIGNED INT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
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
| omni_dump_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci

== TestDumpRestores/mariadb1011 mariadb1011#1
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
< column "VERSION()" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "SEQ_IN_INDEX" "UNSIGNED INT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "MEDIUMTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "LONGBLOB"
< column "EXTRA" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_open_accounts
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip`.`accounts` where `omni_dump_roundtrip`.`accounts`.`status` = 'open'
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "TIME_ZONE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_dump_roundtrip
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
: bytes
| id
< row
: bytes
| omni_dump_roundtrip
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
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
| omni_dump_roundtrip
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "LONGTEXT"
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
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "TEXT"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
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
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_custom`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
< row
: int64 100
: int64 10
: int64 1000
: int64 5
: uint64 20
: int64 1
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_default`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
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
| omni_dump_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE `omni_dump_roundtrip`
< ok 9 0
> exec
| CREATE DATABASE `omni_dump_roundtrip`
< ok 1 0
> exec
| USE `omni_dump_roundtrip`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| SET @OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@@EXPLICIT_DEFAULTS_FOR_TIMESTAMP, EXPLICIT_DEFAULTS_FOR_TIMESTAMP=1
< ok 0 0
> exec
| CREATE VIEW `v_open_accounts` AS SELECT
|   1 AS `acct_id`,
|   1 AS `acct_email`
< ok 0 0
> exec
| CREATE SEQUENCE `seq_custom` start with 100 minvalue 10 maxvalue 1000 increment by 5 cache 20 cycle COMMENT='sequence comment'
< ok 0 0
> exec
| CREATE SEQUENCE `seq_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle
< ok 0 0
> exec
| CREATE TABLE `accounts` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `email` varchar(190) NOT NULL,
|   `bio` text DEFAULT NULL,
|   `status` enum('open','closed') NOT NULL DEFAULT 'open',
|   `balance` decimal(12,4) NOT NULL DEFAULT 0.0000,
|   `secret_code` varchar(32) DEFAULT 'hidden' INVISIBLE,
|   `email_lower` varchar(190) GENERATED ALWAYS AS (lcase(`email`)) VIRTUAL,
|   `display_tag` varchar(220) GENERATED ALWAYS AS (concat(`email`,':','v1')) STORED,
|   `created_at` datetime NOT NULL DEFAULT current_timestamp(),
|   `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6) ON UPDATE CURRENT_TIMESTAMP(6),
|   `note` varchar(100) DEFAULT NULL COMMENT 'owner\'s note',
|   PRIMARY KEY (`id`),
|   KEY `idx_bio_prefix` (`bio`(32)),
|   UNIQUE KEY `uk_email` (`email`),
|   CONSTRAINT `chk_balance` CHECK (`balance` >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='customer accounts'
< ok 0 0
> exec
| CREATE TABLE `explicit_period` (
|   `id` int NOT NULL,
|   `valid_from` timestamp(6) GENERATED ALWAYS AS ROW START NOT NULL,
|   `valid_to` timestamp(6) GENERATED ALWAYS AS ROW END NOT NULL,
|   PERIOD FOR SYSTEM_TIME (`valid_from`, `valid_to`),
|   PRIMARY KEY (`id`, `valid_to`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING
< ok 0 0
> exec
| CREATE TABLE `history` (
|   `id` int NOT NULL,
|   `amount` int DEFAULT NULL,
|   `account_id` bigint DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_history_account` (`account_id`),
|   KEY `idx_amount` (`amount`),
|   CONSTRAINT `fk_history_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`),
|   CONSTRAINT `chk_amount` CHECK (`amount` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='versioned' WITH SYSTEM VERSIONING
< ok 0 0
> exec
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY HASH (`id`)
| PARTITIONS 4
< ok 0 0
> exec
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY RANGE (`ts`)
| (PARTITION `p_low` VALUES LESS THAN (100) ENGINE=InnoDB,
|  PARTITION `p_high` VALUES LESS THAN MAXVALUE ENGINE=InnoDB)
< ok 0 0
> exec
| CREATE TABLE `orders` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `account_id` bigint NOT NULL,
|   `total` int NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_account` (`account_id`),
|   CONSTRAINT `fk_orders_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| SET @saved_cs_client = @@character_set_client, @saved_cs_results = @@character_set_results, @saved_col_connection = @@collation_connection, @saved_sql_mode = @@sql_mode
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_general_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64)
|     NO SQL
|     DETERMINISTIC
|     COMMENT 'formats a tag'
| BEGIN
| 	RETURN CONCAT(p_prefix, '-', p_id);
| END
< ok 0 0
> exec
| SET character_set_client = @saved_cs_client, character_set_results = @saved_cs_results, collation_connection = @saved_col_connection, sql_mode = @saved_sql_mode
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where `accounts`.`status` = 'open'
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_general_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE PROCEDURE `sp_touch_order`(IN p_id BIGINT, OUT p_total INT)
|     COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_general_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| SET time_zone = 'SYSTEM'
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` EVENT `ev_daily_check` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO BEGIN
| 	DELETE FROM orders WHERE total < 0;
| END
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_general_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE TRIGGER `trg_orders_bi`
| BEFORE INSERT ON `orders` FOR EACH ROW
| BEGIN
| 	IF NEW.total < 0 THEN
| 		SET NEW.total = 0;
| 	END IF;
| END
< ok 0 0
> exec
| SET EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "SEQ_IN_INDEX" "UNSIGNED INT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "MEDIUMTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "LONGBLOB"
< column "EXTRA" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_open_accounts
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip`.`accounts` where `omni_dump_roundtrip`.`accounts`.`status` = 'open'
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "TIME_ZONE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_dump_roundtrip
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
: bytes
| id
< row
: bytes
| omni_dump_roundtrip
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
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
| omni_dump_roundtrip
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "LONGTEXT"
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
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_format_tag
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64) CHARSET utf8mb4 COLLATE utf8mb4_general_ci
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
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "TEXT"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
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
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_custom`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
< row
: int64 100
: int64 10
: int64 1000
: int64 5
: uint64 20
: int64 1
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_default`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
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
| omni_dump_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestDumpRestores/mariadb107 mariadb107#1
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
< column "VERSION()" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "SEQ_IN_INDEX" "BIGINT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "MEDIUMTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "LONGBLOB"
< column "EXTRA" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_open_accounts
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip`.`accounts` where `omni_dump_roundtrip`.`accounts`.`status` = 'open'
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "TIME_ZONE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_dump_roundtrip
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
: bytes
| id
< row
: bytes
| omni_dump_roundtrip
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
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
| omni_dump_roundtrip
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "LONGTEXT"
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
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "TEXT"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
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
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_custom`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
< row
: int64 100
: int64 10
: int64 1000
: int64 5
: uint64 20
: int64 1
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_default`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
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
| omni_dump_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE `omni_dump_roundtrip`
< ok 9 0
> exec
| CREATE DATABASE `omni_dump_roundtrip`
< ok 1 0
> exec
| USE `omni_dump_roundtrip`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| SET @OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@@EXPLICIT_DEFAULTS_FOR_TIMESTAMP, EXPLICIT_DEFAULTS_FOR_TIMESTAMP=1
< ok 0 0
> exec
| CREATE VIEW `v_open_accounts` AS SELECT
|   1 AS `acct_id`,
|   1 AS `acct_email`
< ok 0 0
> exec
| CREATE SEQUENCE `seq_custom` start with 100 minvalue 10 maxvalue 1000 increment by 5 cache 20 cycle COMMENT='sequence comment'
< ok 0 0
> exec
| CREATE SEQUENCE `seq_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle
< ok 0 0
> exec
| CREATE TABLE `accounts` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `email` varchar(190) NOT NULL,
|   `bio` text DEFAULT NULL,
|   `status` enum('open','closed') NOT NULL DEFAULT 'open',
|   `balance` decimal(12,4) NOT NULL DEFAULT 0.0000,
|   `secret_code` varchar(32) DEFAULT 'hidden' INVISIBLE,
|   `email_lower` varchar(190) GENERATED ALWAYS AS (lcase(`email`)) VIRTUAL,
|   `display_tag` varchar(220) GENERATED ALWAYS AS (concat(`email`,':','v1')) STORED,
|   `created_at` datetime NOT NULL DEFAULT current_timestamp(),
|   `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6) ON UPDATE CURRENT_TIMESTAMP(6),
|   `note` varchar(100) DEFAULT NULL COMMENT 'owner\'s note',
|   PRIMARY KEY (`id`),
|   KEY `idx_bio_prefix` (`bio`(32)),
|   UNIQUE KEY `uk_email` (`email`),
|   CONSTRAINT `chk_balance` CHECK (`balance` >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='customer accounts'
< ok 0 0
> exec
| CREATE TABLE `explicit_period` (
|   `id` int NOT NULL,
|   `valid_from` timestamp(6) GENERATED ALWAYS AS ROW START NOT NULL,
|   `valid_to` timestamp(6) GENERATED ALWAYS AS ROW END NOT NULL,
|   PERIOD FOR SYSTEM_TIME (`valid_from`, `valid_to`),
|   PRIMARY KEY (`id`, `valid_to`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING
< ok 0 0
> exec
| CREATE TABLE `history` (
|   `id` int NOT NULL,
|   `amount` int DEFAULT NULL,
|   `account_id` bigint DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_history_account` (`account_id`),
|   KEY `idx_amount` (`amount`),
|   CONSTRAINT `fk_history_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`),
|   CONSTRAINT `chk_amount` CHECK (`amount` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='versioned' WITH SYSTEM VERSIONING
< ok 0 0
> exec
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY HASH (`id`)
| PARTITIONS 4
< ok 0 0
> exec
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY RANGE (`ts`)
| (PARTITION `p_low` VALUES LESS THAN (100) ENGINE=InnoDB,
|  PARTITION `p_high` VALUES LESS THAN MAXVALUE ENGINE=InnoDB)
< ok 0 0
> exec
| CREATE TABLE `orders` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `account_id` bigint NOT NULL,
|   `total` int NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_account` (`account_id`),
|   CONSTRAINT `fk_orders_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| SET @saved_cs_client = @@character_set_client, @saved_cs_results = @@character_set_results, @saved_col_connection = @@collation_connection, @saved_sql_mode = @@sql_mode
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_general_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64)
|     NO SQL
|     DETERMINISTIC
|     COMMENT 'formats a tag'
| BEGIN
| 	RETURN CONCAT(p_prefix, '-', p_id);
| END
< ok 0 0
> exec
| SET character_set_client = @saved_cs_client, character_set_results = @saved_cs_results, collation_connection = @saved_col_connection, sql_mode = @saved_sql_mode
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where `accounts`.`status` = 'open'
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_general_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE PROCEDURE `sp_touch_order`(IN p_id BIGINT, OUT p_total INT)
|     COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_general_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| SET time_zone = 'SYSTEM'
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` EVENT `ev_daily_check` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO BEGIN
| 	DELETE FROM orders WHERE total < 0;
| END
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_general_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE TRIGGER `trg_orders_bi`
| BEFORE INSERT ON `orders` FOR EACH ROW
| BEGIN
| 	IF NEW.total < 0 THEN
| 		SET NEW.total = 0;
| 	END IF;
| END
< ok 0 0
> exec
| SET EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "SEQ_IN_INDEX" "BIGINT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "MEDIUMTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "LONGBLOB"
< column "EXTRA" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_open_accounts
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip`.`accounts` where `omni_dump_roundtrip`.`accounts`.`status` = 'open'
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "TIME_ZONE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_dump_roundtrip
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
: bytes
| id
< row
: bytes
| omni_dump_roundtrip
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
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
| omni_dump_roundtrip
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "LONGTEXT"
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
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_format_tag
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64) CHARSET utf8mb4 COLLATE utf8mb4_general_ci
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
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "TEXT"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
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
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_custom`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
< row
: int64 100
: int64 10
: int64 1000
: int64 5
: uint64 20
: int64 1
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_default`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
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
| omni_dump_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestDumpRestores/mariadb114 mariadb114#1
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
< column "VERSION()" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "SEQ_IN_INDEX" "UNSIGNED INT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "MEDIUMTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "LONGBLOB"
< column "EXTRA" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_open_accounts
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip`.`accounts` where `omni_dump_roundtrip`.`accounts`.`status` = 'open'
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "TIME_ZONE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_dump_roundtrip
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
: bytes
| id
< row
: bytes
| omni_dump_roundtrip
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
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
| omni_dump_roundtrip
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "LONGTEXT"
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
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "TEXT"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
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
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_custom`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
< row
: int64 100
: int64 10
: int64 1000
: int64 5
: uint64 20
: int64 1
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_default`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
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
| omni_dump_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE `omni_dump_roundtrip`
< ok 9 0
> exec
| CREATE DATABASE `omni_dump_roundtrip`
< ok 1 0
> exec
| USE `omni_dump_roundtrip`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| SET @OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@@EXPLICIT_DEFAULTS_FOR_TIMESTAMP, EXPLICIT_DEFAULTS_FOR_TIMESTAMP=1
< ok 0 0
> exec
| CREATE VIEW `v_open_accounts` AS SELECT
|   1 AS `acct_id`,
|   1 AS `acct_email`
< ok 0 0
> exec
| CREATE SEQUENCE `seq_custom` start with 100 minvalue 10 maxvalue 1000 increment by 5 cache 20 cycle COMMENT='sequence comment'
< ok 0 0
> exec
| CREATE SEQUENCE `seq_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle
< ok 0 0
> exec
| CREATE TABLE `accounts` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `email` varchar(190) NOT NULL,
|   `bio` text DEFAULT NULL,
|   `status` enum('open','closed') NOT NULL DEFAULT 'open',
|   `balance` decimal(12,4) NOT NULL DEFAULT 0.0000,
|   `secret_code` varchar(32) DEFAULT 'hidden' INVISIBLE,
|   `email_lower` varchar(190) GENERATED ALWAYS AS (lcase(`email`)) VIRTUAL,
|   `display_tag` varchar(220) GENERATED ALWAYS AS (concat(`email`,':','v1')) STORED,
|   `created_at` datetime NOT NULL DEFAULT current_timestamp(),
|   `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6) ON UPDATE CURRENT_TIMESTAMP(6),
|   `note` varchar(100) DEFAULT NULL COMMENT 'owner\'s note',
|   PRIMARY KEY (`id`),
|   KEY `idx_bio_prefix` (`bio`(32)),
|   UNIQUE KEY `uk_email` (`email`),
|   CONSTRAINT `chk_balance` CHECK (`balance` >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='customer accounts'
< ok 0 0
> exec
| CREATE TABLE `explicit_period` (
|   `id` int NOT NULL,
|   `valid_from` timestamp(6) GENERATED ALWAYS AS ROW START NOT NULL,
|   `valid_to` timestamp(6) GENERATED ALWAYS AS ROW END NOT NULL,
|   PERIOD FOR SYSTEM_TIME (`valid_from`, `valid_to`),
|   PRIMARY KEY (`id`, `valid_to`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci WITH SYSTEM VERSIONING
< ok 0 0
> exec
| CREATE TABLE `history` (
|   `id` int NOT NULL,
|   `amount` int DEFAULT NULL,
|   `account_id` bigint DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_history_account` (`account_id`),
|   KEY `idx_amount` (`amount`),
|   CONSTRAINT `fk_history_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`),
|   CONSTRAINT `chk_amount` CHECK (`amount` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='versioned' WITH SYSTEM VERSIONING
< ok 0 0
> exec
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY HASH (`id`)
| PARTITIONS 4
< ok 0 0
> exec
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
| PARTITION BY RANGE (`ts`)
| (PARTITION `p_low` VALUES LESS THAN (100) ENGINE=InnoDB,
|  PARTITION `p_high` VALUES LESS THAN MAXVALUE ENGINE=InnoDB)
< ok 0 0
> exec
| CREATE TABLE `orders` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `account_id` bigint NOT NULL,
|   `total` int NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_account` (`account_id`),
|   CONSTRAINT `fk_orders_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> exec
| SET @saved_cs_client = @@character_set_client, @saved_cs_results = @@character_set_results, @saved_col_connection = @@collation_connection, @saved_sql_mode = @@sql_mode
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64)
|     NO SQL
|     DETERMINISTIC
|     COMMENT 'formats a tag'
| BEGIN
| 	RETURN CONCAT(p_prefix, '-', p_id);
| END
< ok 0 0
> exec
| SET character_set_client = @saved_cs_client, character_set_results = @saved_cs_results, collation_connection = @saved_col_connection, sql_mode = @saved_sql_mode
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where `accounts`.`status` = 'open'
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE PROCEDURE `sp_touch_order`(IN p_id BIGINT, OUT p_total INT)
|     COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| SET time_zone = 'SYSTEM'
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` EVENT `ev_daily_check` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO BEGIN
| 	DELETE FROM orders WHERE total < 0;
| END
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE TRIGGER `trg_orders_bi`
| BEFORE INSERT ON `orders` FOR EACH ROW
| BEGIN
| 	IF NEW.total < 0 THEN
| 		SET NEW.total = 0;
| 	END IF;
| END
< ok 0 0
> exec
| SET EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "SEQ_IN_INDEX" "UNSIGNED INT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "MEDIUMTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "LONGBLOB"
< column "EXTRA" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_open_accounts
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip`.`accounts` where `omni_dump_roundtrip`.`accounts`.`status` = 'open'
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "TIME_ZONE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
< result
< column "Event" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "time_zone" "VARCHAR"
< column "Create Event" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_dump_roundtrip
: bytes
| history
: bytes
| fk_history_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
: bytes
| id
< row
: bytes
| omni_dump_roundtrip
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip
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
| omni_dump_roundtrip
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "VARCHAR"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "LONGTEXT"
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
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_format_tag
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64) CHARSET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
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
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "TEXT"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
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
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_custom`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
< row
: int64 100
: int64 10
: int64 1000
: int64 5
: uint64 20
: int64 1
> query
| SELECT start_value, minimum_value, maximum_value, increment, cache_size, cycle_option FROM `omni_dump_roundtrip`.`seq_default`
< result
< column "start_value" "BIGINT"
< column "minimum_value" "BIGINT"
< column "maximum_value" "BIGINT"
< column "increment" "BIGINT"
< column "cache_size" "UNSIGNED BIGINT"
< column "cycle_option" "UNSIGNED TINYINT"
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
| omni_dump_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
