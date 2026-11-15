-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestSyncDatabase default#1
> exec
| CREATE TABLE accounts (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	email VARCHAR(190) NOT NULL,
| 	bio TEXT,
| 	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
| 	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
| 	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL, -- not restorable: TestDumpGaps
| 	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED, -- not restorable: TestDumpGaps
| 	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6), -- not restorable: TestDumpGaps
| 	note VARCHAR(100) COMMENT 'owner''s note',
| 	UNIQUE KEY uk_email (email),
| 	KEY idx_bio_prefix (bio(32))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='customer accounts';
< ok 0 0
> exec
| CREATE TABLE orders (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	account_id BIGINT NOT NULL,
| 	total INT NOT NULL,
| 	KEY idx_account (account_id),
| 	CONSTRAINT fk_orders_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
< ok 0 0
> exec
| CREATE EXTERNAL TABLE ext_feed (c1 INT, c2 VARCHAR(10)) LOCATION = '/tmp' FORMAT = (TYPE = 'CSV'); -- not restorable: TestDumpGaps
< ok 0 0
> exec
| CREATE MATERIALIZED VIEW mv_accounts COMMENT 'materialized accounts' AS SELECT id, email FROM accounts; -- not restorable: TestDumpGaps
< ok 0 0
> exec
| CREATE INDEX idx_mv_email ON mv_accounts (email); -- not restorable: TestDumpGaps
< ok 0 0
> exec
| CREATE VIEW v_open_accounts AS
| 	SELECT id AS acct_id, email AS acct_email FROM accounts WHERE status = 'open';
< ok 0 0
> exec
| CREATE FUNCTION fn_format_tag(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT)
| RETURNS VARCHAR(64) CHARSET utf8mb4 COLLATE utf8mb4_general_ci
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
| CREATE TRIGGER trg_orders_bi BEFORE INSERT ON orders FOR EACH ROW BEGIN IF NEW.total < 0 THEN SET NEW.total = 0; END IF; END; -- not restorable: TestDumpGaps
< ok 0 0
> exec
| CREATE EVENT ev_daily_check ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DISABLE DO DELETE FROM orders WHERE total < 0; -- not restorable: TestDumpGaps
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
< row
: bytes
| 5.7.25-OceanBase_CE-v4.3.5.6
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			INDEX_NAME,
| 			COLUMN_NAME,
| 			COLLATION,
| 			IFNULL(SUB_PART, -1),
| 			'',
| 			SEQ_IN_INDEX,
| 			INDEX_TYPE,
| 			NON_UNIQUE,
| 			1,
| 			INDEX_COMMENT
| 		FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "" "VARCHAR"
< column "SEQ_IN_INDEX" "UNSIGNED BIGINT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "1" "BIGINT"
< column "INDEX_COMMENT" "TEXT"
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
: bytes
|
: uint64 1
: bytes
| BTREE
: int64 1
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
: bytes
|
: uint64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
| customer accounts
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
: bytes
|
: uint64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: bytes
|
: uint64 1
: bytes
| BTREE
: int64 1
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
: bytes
|
: uint64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| __mv_container_0000000000000001
: bytes
| idx_mv_email
: bytes
| email
: bytes
| A
: int64 -1
: bytes
|
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
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
| 			EXTRA,
| 			NULL
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_sync_corpus
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "convert(GENERATION_EXPRESSION using BINARY)" "VARBINARY"
< column "EXTRA" "VARCHAR"
< column "NULL" "NULL"
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
: bytes
|
: bytes
| auto_increment
: null
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
| utf8mb4_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| accounts
: bytes
| bio
: uint64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| accounts
: bytes
| status
: uint64 4
: bytes
| 'open'
: bytes
| NO
: bytes
| enum('open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
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
: bytes
|
: bytes
|
: null
< row
: bytes
| accounts
: bytes
| email_lower
: uint64 6
: bytes
| 'LOWER(`email`)'
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| ''
: bytes
| LOWER(`email`)
: bytes
| VIRTUAL GENERATED
: null
< row
: bytes
| accounts
: bytes
| display_tag
: uint64 7
: bytes
| 'CONCAT(`email`,\':\',\'v1\')'
: bytes
| YES
: bytes
| varchar(220)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| ''
: bytes
| CONCAT(`email`,':','v1')
: bytes
| STORED GENERATED
: null
< row
: bytes
| accounts
: bytes
| created_at
: uint64 8
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
: null
< row
: bytes
| accounts
: bytes
| updated_at
: uint64 9
: bytes
| 'CURRENT_TIMESTAMP'
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
| on update current_timestamp(6)
: null
< row
: bytes
| accounts
: bytes
| note
: uint64 10
: null
: bytes
| YES
: bytes
| varchar(100)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| 'owner\'s note'
: bytes
|
: bytes
|
: null
< row
: bytes
| ext_feed
: bytes
| c1
: uint64 1
: bytes
| 'metadata$filecol1'
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
: bytes
| metadata$filecol1
: bytes
| STORED GENERATED
: null
< row
: bytes
| ext_feed
: bytes
| c2
: uint64 2
: bytes
| 'metadata$filecol2'
: bytes
| YES
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| ''
: bytes
| metadata$filecol2
: bytes
| STORED GENERATED
: null
< row
: bytes
| mv_accounts
: bytes
| id
: uint64 1
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
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| mv_accounts
: bytes
| email
: uint64 2
: null
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
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
: bytes
|
: bytes
| auto_increment
: null
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
: bytes
|
: bytes
|
: null
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
: bytes
|
: bytes
|
: null
< row
: bytes
| v_open_accounts
: bytes
| acct_id
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
: bytes
|
: bytes
|
: null
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
| utf8mb4_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| __mv_container_0000000000000001
: bytes
| id
: uint64 1
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
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| __mv_container_0000000000000001
: bytes
| email
: uint64 2
: null
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
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
| select `omni_sync_corpus`.`accounts`.`id` AS `acct_id`,`omni_sync_corpus`.`accounts`.`email` AS `acct_email` from `omni_sync_corpus`.`accounts` where (`omni_sync_corpus`.`accounts`.`status` = 'open')
> query
| SHOW CREATE VIEW `omni_sync_corpus`.`v_open_accounts`
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_open_accounts
: bytes
| CREATE VIEW `v_open_accounts` AS select `omni_sync_corpus`.`accounts`.`id` AS `acct_id`,`omni_sync_corpus`.`accounts`.`email` AS `acct_email` from `omni_sync_corpus`.`accounts` where (`omni_sync_corpus`.`accounts`.`status` = 'open')
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
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "TEXT"
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
| TRIGGER trg_orders_bi BEFORE INSERT ON orders FOR EACH ROW BEGIN IF NEW.total < 0 THEN SET NEW.total = 0; END IF; END
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
|
| 		SELECT DISTINCT EVENT_NAME
| 		FROM information_schema.EVENTS
| 		WHERE EVENT_SCHEMA = ?
| 		ORDER BY EVENT_NAME
: string
| omni_sync_corpus
< result
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| ev_daily_check
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
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "TEXT"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT"
< row
: bytes
| fn_format_tag
: bytes
| FUNCTION
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
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
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
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
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_format_tag
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64)
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
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| sp_touch_order
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` PROCEDURE `sp_touch_order`(IN p_id BIGINT, OUT p_total INT)
|     READS SQL DATA
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
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
< row
: bytes
| accounts
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_general_ci
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
| 'customer accounts'
: bytes
| utf8mb4
< row
: bytes
| ext_feed
: bytes
| EXTERNAL TABLE
: bytes
| MEMORY
: bytes
| utf8mb4_general_ci
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
| utf8mb4_general_ci
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
| ''
: bytes
| utf8mb4
< row
: bytes
| v_open_accounts
: bytes
| VIEW
: bytes
| MEMORY
: bytes
| utf8mb4_general_ci
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
| utf8mb4
< row
: bytes
| __mv_container_0000000000000001
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_general_ci
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
| 'materialized accounts'
: bytes
| utf8mb4
> query
|
| 		SELECT MVIEW_NAME, CONTAINER_NAME, QUERY
| 		FROM oceanbase.DBA_MVIEWS
| 		WHERE OWNER = ?
| 		ORDER BY MVIEW_NAME
: string
| omni_sync_corpus
< result
< column "MVIEW_NAME" "VARCHAR"
< column "CONTAINER_NAME" "VARCHAR"
< column "QUERY" "TEXT"
< row
: bytes
| mv_accounts
: bytes
| __mv_container_0000000000000001
: bytes
| select `omni_sync_corpus`.`accounts`.`id` AS `id`,`omni_sync_corpus`.`accounts`.`email` AS `email` from `omni_sync_corpus`.`accounts`
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
| utf8mb4_general_ci
