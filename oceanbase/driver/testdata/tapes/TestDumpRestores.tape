-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDumpRestores default#1
> exec
| CREATE TABLE accounts (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	email VARCHAR(190) NOT NULL,
| 	bio TEXT,
| 	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
| 	literal_ts VARCHAR(40) NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
| 	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
| 	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
| CREATE EXTERNAL TABLE ext_feed (c1 INT, c2 VARCHAR(10)) LOCATION = '/tmp' PATTERN = '.*[.]csv' FORMAT = (TYPE = 'CSV');
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR" "string"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "" "VARCHAR" "string"
< column "SEQ_IN_INDEX" "UNSIGNED BIGINT" "uint64"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "1" "BIGINT" "int64"
< column "INDEX_COMMENT" "TEXT" "string"
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
| omni_dump_roundtrip
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
< column "convert(GENERATION_EXPRESSION using BINARY)" "VARBINARY" "[]uint8"
< column "EXTRA" "VARCHAR"
< column "NULL" "NULL" "*interface {}"
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
| literal_ts
: uint64 5
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| varchar(40)
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
: uint64 6
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
| created_at
: uint64 7
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
| note
: uint64 8
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
| v_named_cols
: bytes
| renamed_id
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
| v_named_cols_ref
: bytes
| renamed_id
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v_named_cols
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `renamed_id` from `omni_dump_roundtrip`.`accounts`
< row
: bytes
| v_named_cols_ref
: bytes
| select `omni_dump_roundtrip`.`v_named_cols`.`renamed_id` AS `renamed_id` from `omni_dump_roundtrip`.`v_named_cols`
< row
: bytes
| v_open_accounts
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip`.`accounts` where (`omni_dump_roundtrip`.`accounts`.`status` = 'open')
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols
: bytes
| CREATE VIEW `v_named_cols` AS select `omni_dump_roundtrip`.`accounts`.`id` AS `renamed_id` from `omni_dump_roundtrip`.`accounts`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols_ref`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols_ref
: bytes
| CREATE VIEW `v_named_cols_ref` AS select `omni_dump_roundtrip`.`v_named_cols`.`renamed_id` AS `renamed_id` from `omni_dump_roundtrip`.`v_named_cols`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_open_accounts
: bytes
| CREATE VIEW `v_open_accounts` AS select `omni_dump_roundtrip`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip`.`accounts` where (`omni_dump_roundtrip`.`accounts`.`status` = 'open')
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
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "TEXT"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
> query
|
| 		SELECT DISTINCT EVENT_NAME
| 		FROM information_schema.EVENTS
| 		WHERE EVENT_SCHEMA = ?
| 		ORDER BY EVENT_NAME
: string
| omni_dump_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
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
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
< column "ROUTINE_NAME" "VARCHAR" "string"
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
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
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
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
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
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" precision 20 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" precision 20 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" precision 20 scale 0
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
| v_named_cols
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
| v_named_cols_ref
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
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`ext_feed`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| ext_feed
: bytes
| CREATE EXTERNAL TABLE `ext_feed` (
|   `c1` int(11) GENERATED ALWAYS AS (metadata$filecol1),
|   `c2` varchar(10) GENERATED ALWAYS AS (metadata$filecol2)
| ) 
| LOCATION='file:///tmp'
| PATTERN='.*[.]csv'
| AUTO_REFRESH = OFF
| FORMAT (
|   TYPE = 'CSV',
|   ENCODING = 'utf8mb4'
| ) DEFAULT CHARSET = utf8mb4
> query
|
| 		SELECT MVIEW_NAME, CONTAINER_NAME, QUERY
| 		FROM oceanbase.DBA_MVIEWS
| 		WHERE OWNER = ?
| 		ORDER BY MVIEW_NAME
: string
| omni_dump_roundtrip
< result
< column "MVIEW_NAME" "VARCHAR" "string"
< column "CONTAINER_NAME" "VARCHAR" "string"
< column "QUERY" "TEXT" "string"
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
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| USE `omni_dump_roundtrip_restored`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE VIEW `v_named_cols` AS SELECT
|   1 AS `renamed_id`
< ok 0 0
> exec
| CREATE VIEW `v_named_cols_ref` AS SELECT
|   1 AS `renamed_id`
< ok 0 0
> exec
| CREATE VIEW `v_open_accounts` AS SELECT
|   1 AS `acct_id`,
|   1 AS `acct_email`
< ok 0 0
> exec
| CREATE TABLE `accounts` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `email` varchar(190) NOT NULL,
|   `bio` text DEFAULT NULL,
|   `status` enum('open','closed') NOT NULL DEFAULT 'open',
|   `literal_ts` varchar(40) NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
|   `balance` decimal(12,4) NOT NULL DEFAULT '0.0000',
|   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   `note` varchar(100) DEFAULT NULL COMMENT 'owner\'s note',
|   PRIMARY KEY (`id`),
|   KEY `idx_bio_prefix` (`bio`(32)),
|   UNIQUE KEY `uk_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='customer accounts'
< ok 0 0
> exec
| CREATE TABLE `orders` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `account_id` bigint NOT NULL,
|   `total` int NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_account` (`account_id`),
|   CONSTRAINT `fk_orders_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE EXTERNAL TABLE `ext_feed` (
|   `c1` int GENERATED ALWAYS AS (metadata$filecol1),
|   `c2` varchar(10) GENERATED ALWAYS AS (metadata$filecol2)
| )
| LOCATION = 'file:///tmp'
| PATTERN = '.*[.]csv'
| FORMAT (TYPE = 'CSV', ENCODING = 'utf8mb4')
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
| SET sql_mode = 'STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER'
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
| CREATE OR REPLACE VIEW `v_named_cols` AS select `omni_dump_roundtrip_restored`.`accounts`.`id` AS `renamed_id` from `omni_dump_roundtrip_restored`.`accounts`
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_named_cols_ref` AS select `omni_dump_roundtrip_restored`.`v_named_cols`.`renamed_id` AS `renamed_id` from `omni_dump_roundtrip_restored`.`v_named_cols`
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_open_accounts` AS select `omni_dump_roundtrip_restored`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip_restored`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip_restored`.`accounts` where (`omni_dump_roundtrip_restored`.`accounts`.`status` = 'open')
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
| SET sql_mode = 'STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER'
< ok 0 0
> exec
| CREATE PROCEDURE `sp_touch_order`(IN p_id BIGINT, OUT p_total INT)
|     READS SQL DATA
|     COMMENT 'reads order total'
| BEGIN
| 	SELECT total INTO p_total FROM orders WHERE id = p_id;
| END
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
| omni_dump_roundtrip_restored
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR" "string"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "" "VARCHAR" "string"
< column "SEQ_IN_INDEX" "UNSIGNED BIGINT" "uint64"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "1" "BIGINT" "int64"
< column "INDEX_COMMENT" "TEXT" "string"
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
| omni_dump_roundtrip_restored
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
< column "convert(GENERATION_EXPRESSION using BINARY)" "VARBINARY" "[]uint8"
< column "EXTRA" "VARCHAR"
< column "NULL" "NULL" "*interface {}"
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
| literal_ts
: uint64 5
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| varchar(40)
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
: uint64 6
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
| created_at
: uint64 7
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
| note
: uint64 8
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
| v_named_cols
: bytes
| renamed_id
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
| v_named_cols_ref
: bytes
| renamed_id
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
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_dump_roundtrip_restored
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v_named_cols
: bytes
| select `omni_dump_roundtrip_restored`.`accounts`.`id` AS `renamed_id` from `omni_dump_roundtrip_restored`.`accounts`
< row
: bytes
| v_named_cols_ref
: bytes
| select `omni_dump_roundtrip_restored`.`v_named_cols`.`renamed_id` AS `renamed_id` from `omni_dump_roundtrip_restored`.`v_named_cols`
< row
: bytes
| v_open_accounts
: bytes
| select `omni_dump_roundtrip_restored`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip_restored`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip_restored`.`accounts` where (`omni_dump_roundtrip_restored`.`accounts`.`status` = 'open')
> query
| SHOW CREATE VIEW `omni_dump_roundtrip_restored`.`v_named_cols`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols
: bytes
| CREATE VIEW `v_named_cols` AS select `omni_dump_roundtrip_restored`.`accounts`.`id` AS `renamed_id` from `omni_dump_roundtrip_restored`.`accounts`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip_restored`.`v_named_cols_ref`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols_ref
: bytes
| CREATE VIEW `v_named_cols_ref` AS select `omni_dump_roundtrip_restored`.`v_named_cols`.`renamed_id` AS `renamed_id` from `omni_dump_roundtrip_restored`.`v_named_cols`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip_restored`.`v_open_accounts`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_open_accounts
: bytes
| CREATE VIEW `v_open_accounts` AS select `omni_dump_roundtrip_restored`.`accounts`.`id` AS `acct_id`,`omni_dump_roundtrip_restored`.`accounts`.`email` AS `acct_email` from `omni_dump_roundtrip_restored`.`accounts` where (`omni_dump_roundtrip_restored`.`accounts`.`status` = 'open')
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
| omni_dump_roundtrip_restored
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "TEXT"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
> query
|
| 		SELECT DISTINCT EVENT_NAME
| 		FROM information_schema.EVENTS
| 		WHERE EVENT_SCHEMA = ?
| 		ORDER BY EVENT_NAME
: string
| omni_dump_roundtrip_restored
< result
< column "EVENT_NAME" "VARCHAR"
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
| omni_dump_roundtrip_restored
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
| omni_dump_roundtrip_restored
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_dump_roundtrip_restored
: bytes
| orders
: bytes
| fk_orders_account
: bytes
| account_id
: bytes
| omni_dump_roundtrip_restored
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
| omni_dump_roundtrip_restored
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
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
| SHOW CREATE FUNCTION `omni_dump_roundtrip_restored`.`fn_format_tag`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
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
| SHOW CREATE PROCEDURE `omni_dump_roundtrip_restored`.`sp_touch_order`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
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
| omni_dump_roundtrip_restored
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" precision 20 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" precision 20 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" precision 20 scale 0
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
| v_named_cols
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
| v_named_cols_ref
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
> query
| SHOW CREATE TABLE `omni_dump_roundtrip_restored`.`ext_feed`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| ext_feed
: bytes
| CREATE EXTERNAL TABLE `ext_feed` (
|   `c1` int(11) GENERATED ALWAYS AS (metadata$filecol1),
|   `c2` varchar(10) GENERATED ALWAYS AS (metadata$filecol2)
| ) 
| LOCATION='file:///tmp'
| PATTERN='.*[.]csv'
| AUTO_REFRESH = OFF
| FORMAT (
|   TYPE = 'CSV',
|   ENCODING = 'utf8mb4'
| ) DEFAULT CHARSET = utf8mb4
> query
|
| 		SELECT MVIEW_NAME, CONTAINER_NAME, QUERY
| 		FROM oceanbase.DBA_MVIEWS
| 		WHERE OWNER = ?
| 		ORDER BY MVIEW_NAME
: string
| omni_dump_roundtrip_restored
< result
< column "MVIEW_NAME" "VARCHAR" "string"
< column "CONTAINER_NAME" "VARCHAR" "string"
< column "QUERY" "TEXT" "string"
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_dump_roundtrip_restored
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
