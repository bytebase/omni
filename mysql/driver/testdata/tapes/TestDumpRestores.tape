-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46
-- pin: lctn1 mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46
-- pin: mysql84 mysql:8.4.13 sha256:a81d043930e51745522bd5c430961cf8286f0a23f0a9a2a6e0108f7ea71b0902 mysql-8.4.13

== TestDumpRestores/default default#1
> exec
| CREATE TABLE accounts (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	email VARCHAR(190) NOT NULL,
| 	bio TEXT,
| 	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
| 	flags SET('alpha', 'beta') DEFAULT 'alpha',
| 	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
| 	raw_token BINARY(8) NOT NULL DEFAULT 'ab',
| 	raw_bytes VARBINARY(16) NOT NULL DEFAULT 0x0100fe,
| 	secret_code VARCHAR(32) INVISIBLE DEFAULT 'hidden',
| 	meta JSON,
| 	coords POINT NOT NULL SRID 4326,
| 	geom GEOMETRY SRID 0,
| 	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL,
| 	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED,
| 	external_id VARCHAR(36) DEFAULT (UUID()),
| 	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
| 	note VARCHAR(100) COMMENT 'owner''s note \\ path',
| 	UNIQUE KEY uk_email (email),
| 	KEY idx_bio_prefix (bio(32)),
| 	KEY idx_created_desc (created_at DESC) COMMENT 'newest first',
| 	KEY idx_invisible (status) INVISIBLE,
| 	KEY idx_expr ((LOWER(status))),
| 	FULLTEXT KEY ft_bio (bio),
| 	SPATIAL INDEX sp_coords (coords),
| 	CONSTRAINT chk_balance CHECK (balance >= 0 AND status IN ('open', 'closed'))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='customer''s accounts';
< ok 0 0
> exec
| CREATE TABLE orders (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	account_id BIGINT NOT NULL,
| 	parent_order_id BIGINT DEFAULT NULL,
| 	total INT NOT NULL,
| 	KEY idx_account (account_id),
| 	KEY idx_parent (parent_order_id),
| 	CONSTRAINT fk_orders_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT,
| 	CONSTRAINT fk_orders_parent FOREIGN KEY (parent_order_id) REFERENCES orders (id) ON DELETE SET NULL ON UPDATE NO ACTION
| ) ENGINE=InnoDB;
< ok 0 0
> exec
| CREATE TABLE metrics_range (
| 	ts INT NOT NULL,
| 	region INT NOT NULL,
| 	val BIGINT NOT NULL,
| 	PRIMARY KEY (ts, region)
| ) PARTITION BY RANGE (ts) SUBPARTITION BY HASH (region) SUBPARTITIONS 2 (
| 	PARTITION p0 VALUES LESS THAN (100) COMMENT 'first',
| 	PARTITION p1 VALUES LESS THAN (MAXVALUE) COMMENT 'last'
| );
< ok 0 0
> exec
| CREATE TABLE metrics_hash (
| 	id INT NOT NULL PRIMARY KEY,
| 	val INT NOT NULL
| ) PARTITION BY HASH (id) PARTITIONS 4;
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
| CREATE VIEW v_order_totals AS
| 	SELECT a.acct_id, COUNT(o.id) AS order_count
| 	FROM v_open_accounts a
| 	LEFT JOIN orders o ON o.account_id = a.acct_id
| 	GROUP BY a.acct_id;
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
| 8.0.46
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				INDEX_NAME,
| 				COLUMN_NAME,
| 				COLLATION,
| 				IFNULL(SUB_PART, -1),
| 				EXPRESSION,
| 				SEQ_IN_INDEX,
| 				INDEX_TYPE,
| 				NON_UNIQUE,
| 				CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END,
| 				INDEX_COMMENT
| 			FROM information_schema.STATISTICS
| 			WHERE TABLE_SCHEMA = ?
| 			ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "INT" "int32"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT" "int32"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| ft_bio
: bytes
| bio
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
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
: null
: int64 1
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
| idx_created_desc
: bytes
| created_at
: bytes
| D
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
| newest first
< row
: bytes
| accounts
: bytes
| idx_expr
: null
: bytes
| A
: int64 -1
: bytes
| lower(`status`)
: int64 1
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
| idx_invisible
: bytes
| status
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 0
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| accounts
: bytes
| sp_coords
: bytes
| coords
: bytes
| A
: int64 32
: null
: int64 1
: bytes
| SPATIAL
: int64 1
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| metrics_range
: bytes
| PRIMARY
: bytes
| region
: bytes
| A
: int64 -1
: null
: int64 2
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
: null
: int64 1
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
| idx_parent
: bytes
| parent_order_id
: bytes
| A
: int64 -1
: null
: int64 1
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
: null
: int64 1
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
| 			SRS_ID
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED INT" "uint32"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB" "[]uint8"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
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
: null
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
| utf8mb4_0900_ai_ci
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
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
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
| utf8mb4_0900_ai_ci
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
| flags
: int64 5
: bytes
| 'alpha'
: bytes
| YES
: bytes
| set('alpha','beta')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
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
: null
< row
: bytes
| accounts
: bytes
| raw_token
: int64 7
: bytes
| '0x6162'
: bytes
| NO
: bytes
| binary(8)
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
| raw_bytes
: int64 8
: bytes
| '0x01'
: bytes
| NO
: bytes
| varbinary(16)
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
| secret_code
: int64 9
: bytes
| 'hidden'
: bytes
| YES
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
| INVISIBLE
: null
< row
: bytes
| accounts
: bytes
| meta
: int64 10
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
| coords
: int64 11
: null
: bytes
| NO
: bytes
| point
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
: int64 4326
< row
: bytes
| accounts
: bytes
| geom
: int64 12
: null
: bytes
| YES
: bytes
| geometry
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
: int64 0
< row
: bytes
| accounts
: bytes
| email_lower
: int64 13
: null
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| lower(`email`)
: bytes
| VIRTUAL GENERATED
: null
< row
: bytes
| accounts
: bytes
| display_tag
: int64 14
: null
: bytes
| YES
: bytes
| varchar(220)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| concat(`email`,_utf8mb4\':\',_utf8mb4\'v1\')
: bytes
| STORED GENERATED
: null
< row
: bytes
| accounts
: bytes
| external_id
: int64 15
: bytes
| 'uuid()'
: bytes
| YES
: bytes
| varchar(36)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| accounts
: bytes
| created_at
: int64 16
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
| DEFAULT_GENERATED
: null
< row
: bytes
| accounts
: bytes
| updated_at
: int64 17
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
: null
< row
: bytes
| accounts
: bytes
| note
: int64 18
: null
: bytes
| YES
: bytes
| varchar(100)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 'owner\'s note \\ path'
: bytes
|
: bytes
|
: null
< row
: bytes
| metrics_hash
: bytes
| id
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
: null
< row
: bytes
| metrics_hash
: bytes
| val
: int64 2
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
: null
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
: null
< row
: bytes
| metrics_range
: bytes
| region
: int64 2
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
: null
< row
: bytes
| metrics_range
: bytes
| val
: int64 3
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
: null
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
: null
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
: null
< row
: bytes
| orders
: bytes
| parent_order_id
: int64 3
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
: null
< row
: bytes
| orders
: bytes
| total
: int64 4
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
: null
< row
: bytes
| v_named_cols
: bytes
| renamed_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_named_cols_ref
: bytes
| renamed_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_open_accounts
: bytes
| acct_id
: int64 1
: bytes
| '0'
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
: null
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
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| v_order_totals
: bytes
| acct_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_order_totals
: bytes
| order_count
: int64 2
: bytes
| '0'
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
: null
> query
|
| 			SELECT
| 				tc.TABLE_NAME,
| 				cc.CONSTRAINT_NAME,
| 				cc.CHECK_CLAUSE
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| omni_dump_roundtrip
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
| ((`balance` >= 0) and (`status` in (_utf8mb4\'open\',_utf8mb4\'closed\')))
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
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_named_cols
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `id` from `omni_dump_roundtrip`.`accounts`
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
< row
: bytes
| v_order_totals
: bytes
| select `omni_dump_roundtrip`.`a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`omni_dump_roundtrip`.`v_open_accounts` `a` left join `omni_dump_roundtrip`.`orders` `o` on((`o`.`account_id` = `omni_dump_roundtrip`.`a`.`acct_id`))) group by `omni_dump_roundtrip`.`a`.`acct_id`
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols` (`renamed_id`) AS select `accounts`.`id` AS `id` from `accounts`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols_ref`
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
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_open_accounts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where (`accounts`.`status` = 'open')
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_order_totals`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_order_totals
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_order_totals` AS select `a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`v_open_accounts` `a` left join `orders` `o` on((`o`.`account_id` = `a`.`acct_id`))) group by `a`.`acct_id`
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
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| SYSTEM
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
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
< row
: bytes
| orders
: bytes
| fk_orders_parent
: bytes
| orders
: bytes
| SET NULL
: bytes
| NO ACTION
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
< column "COLUMN_NAME" "VARCHAR"
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
< row
: bytes
| omni_dump_roundtrip
: bytes
| orders
: bytes
| fk_orders_parent
: bytes
| parent_order_id
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
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< column "PARTITION_COMMENT" "TEXT" "string"
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
: bytes
|
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
: bytes
|
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
: bytes
|
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
: bytes
|
< row
: bytes
| metrics_range
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| 100
: bytes
| first
< row
: bytes
| metrics_range
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| 100
: bytes
| first
< row
: bytes
| metrics_range
: bytes
| p1
: bytes
| p1sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| MAXVALUE
: bytes
| last
< row
: bytes
| metrics_range
: bytes
| p1
: bytes
| p1sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| MAXVALUE
: bytes
| last
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_hash
: bytes
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_range
: bytes
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `region` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`,`region`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`ts`)
| SUBPARTITION BY HASH (`region`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (100) COMMENT = 'first' ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE COMMENT = 'last' ENGINE = InnoDB) */
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
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
< row
: bytes
| fn_format_tag
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| formats a tag
< row
: bytes
| sp_touch_order
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| reads order total
> query
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "TABLE_TYPE" "ENUM" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR" "string"
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
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 131072
: bytes
| 0
: bytes
|
: bytes
| 'customer\'s accounts'
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
| utf8mb4_0900_ai_ci
: bytes
| 0
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
| utf8mb4_0900_ai_ci
: bytes
| 0
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
| orders
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
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
< row
: bytes
| v_order_totals
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
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE `omni_dump_roundtrip`
< ok 8 0
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
| CREATE VIEW `v_order_totals` AS SELECT
|   1 AS `acct_id`,
|   1 AS `order_count`
< ok 0 0
> exec
| CREATE TABLE `accounts` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `email` varchar(190) NOT NULL,
|   `bio` text DEFAULT NULL,
|   `status` enum('open','closed') NOT NULL DEFAULT 'open',
|   `flags` set('alpha','beta') DEFAULT 'alpha',
|   `balance` decimal(12,4) NOT NULL DEFAULT '0.0000',
|   `raw_token` binary(8) NOT NULL DEFAULT 'ab',
|   `raw_bytes` varbinary(16) NOT NULL DEFAULT 0x01,
|   `secret_code` varchar(32) DEFAULT 'hidden' /*!80023 INVISIBLE */,
|   `meta` json,
|   `coords` point NOT NULL /*!80003 SRID 4326 */,
|   `geom` geometry /*!80003 SRID 0 */,
|   `email_lower` varchar(190) GENERATED ALWAYS AS (lower(`email`)) VIRTUAL,
|   `display_tag` varchar(220) GENERATED ALWAYS AS (concat(`email`,_utf8mb4':',_utf8mb4'v1')) STORED,
|   `external_id` varchar(36) DEFAULT (uuid()),
|   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
|   `note` varchar(100) DEFAULT NULL COMMENT 'owner\'s note \\ path',
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ft_bio` (`bio`),
|   KEY `idx_bio_prefix` (`bio`(32)),
|   KEY `idx_created_desc` (`created_at` DESC) COMMENT 'newest first',
|   KEY `idx_expr` ((lower(`status`))),
|   KEY `idx_invisible` (`status`) /*!80000 INVISIBLE */,
|   SPATIAL KEY `sp_coords` (`coords`),
|   UNIQUE KEY `uk_email` (`email`),
|   CONSTRAINT `chk_balance` CHECK ((`balance` >= 0) and (`status` in (_utf8mb4'open',_utf8mb4'closed')))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='customer\'s accounts'
< ok 0 0
> exec
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
< ok 0 0
> exec
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `region` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`, `region`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`ts`)
| SUBPARTITION BY HASH (`region`)
| SUBPARTITIONS 2
| (PARTITION `p0` VALUES LESS THAN (100) ENGINE=InnoDB COMMENT = 'first',
|  PARTITION `p1` VALUES LESS THAN MAXVALUE ENGINE=InnoDB COMMENT = 'last') */
< ok 0 0
> exec
| CREATE TABLE `orders` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `account_id` bigint NOT NULL,
|   `parent_order_id` bigint DEFAULT NULL,
|   `total` int NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_account` (`account_id`),
|   KEY `idx_parent` (`parent_order_id`),
|   CONSTRAINT `fk_orders_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
|   CONSTRAINT `fk_orders_parent` FOREIGN KEY (`parent_order_id`) REFERENCES `orders` (`id`) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
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
| CREATE OR REPLACE VIEW `v_named_cols` (`renamed_id`) AS select `accounts`.`id` AS `id` from `accounts`
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_named_cols_ref` AS select `v_named_cols`.`renamed_id` AS `renamed_id` from `v_named_cols`
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where (`accounts`.`status` = 'open')
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_order_totals` AS select `a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`v_open_accounts` `a` left join `orders` `o` on((`o`.`account_id` = `a`.`acct_id`))) group by `a`.`acct_id`
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
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
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
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
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
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
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
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR" "string"
< row
: bytes
| 8.0.46
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				INDEX_NAME,
| 				COLUMN_NAME,
| 				COLLATION,
| 				IFNULL(SUB_PART, -1),
| 				EXPRESSION,
| 				SEQ_IN_INDEX,
| 				INDEX_TYPE,
| 				NON_UNIQUE,
| 				CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END,
| 				INDEX_COMMENT
| 			FROM information_schema.STATISTICS
| 			WHERE TABLE_SCHEMA = ?
| 			ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "INT" "int32"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT" "int32"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| ft_bio
: bytes
| bio
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
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
: null
: int64 1
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
| idx_created_desc
: bytes
| created_at
: bytes
| D
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
| newest first
< row
: bytes
| accounts
: bytes
| idx_expr
: null
: bytes
| A
: int64 -1
: bytes
| lower(`status`)
: int64 1
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
| idx_invisible
: bytes
| status
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 0
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| accounts
: bytes
| sp_coords
: bytes
| coords
: bytes
| A
: int64 32
: null
: int64 1
: bytes
| SPATIAL
: int64 1
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| metrics_range
: bytes
| PRIMARY
: bytes
| region
: bytes
| A
: int64 -1
: null
: int64 2
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
: null
: int64 1
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
| idx_parent
: bytes
| parent_order_id
: bytes
| A
: int64 -1
: null
: int64 1
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
: null
: int64 1
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
| 			SRS_ID
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED INT" "uint32"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB" "[]uint8"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
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
: null
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
| utf8mb4_0900_ai_ci
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
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
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
| utf8mb4_0900_ai_ci
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
| flags
: int64 5
: bytes
| 'alpha'
: bytes
| YES
: bytes
| set('alpha','beta')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
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
: null
< row
: bytes
| accounts
: bytes
| raw_token
: int64 7
: bytes
| '0x6162'
: bytes
| NO
: bytes
| binary(8)
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
| raw_bytes
: int64 8
: bytes
| '0x01'
: bytes
| NO
: bytes
| varbinary(16)
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
| secret_code
: int64 9
: bytes
| 'hidden'
: bytes
| YES
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
| INVISIBLE
: null
< row
: bytes
| accounts
: bytes
| meta
: int64 10
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
| coords
: int64 11
: null
: bytes
| NO
: bytes
| point
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
: int64 4326
< row
: bytes
| accounts
: bytes
| geom
: int64 12
: null
: bytes
| YES
: bytes
| geometry
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
: int64 0
< row
: bytes
| accounts
: bytes
| email_lower
: int64 13
: null
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| lower(`email`)
: bytes
| VIRTUAL GENERATED
: null
< row
: bytes
| accounts
: bytes
| display_tag
: int64 14
: null
: bytes
| YES
: bytes
| varchar(220)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| concat(`email`,_utf8mb4\':\',_utf8mb4\'v1\')
: bytes
| STORED GENERATED
: null
< row
: bytes
| accounts
: bytes
| external_id
: int64 15
: bytes
| 'uuid()'
: bytes
| YES
: bytes
| varchar(36)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| accounts
: bytes
| created_at
: int64 16
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
| DEFAULT_GENERATED
: null
< row
: bytes
| accounts
: bytes
| updated_at
: int64 17
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
: null
< row
: bytes
| accounts
: bytes
| note
: int64 18
: null
: bytes
| YES
: bytes
| varchar(100)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 'owner\'s note \\ path'
: bytes
|
: bytes
|
: null
< row
: bytes
| metrics_hash
: bytes
| id
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
: null
< row
: bytes
| metrics_hash
: bytes
| val
: int64 2
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
: null
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
: null
< row
: bytes
| metrics_range
: bytes
| region
: int64 2
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
: null
< row
: bytes
| metrics_range
: bytes
| val
: int64 3
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
: null
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
: null
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
: null
< row
: bytes
| orders
: bytes
| parent_order_id
: int64 3
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
: null
< row
: bytes
| orders
: bytes
| total
: int64 4
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
: null
< row
: bytes
| v_named_cols
: bytes
| renamed_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_named_cols_ref
: bytes
| renamed_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_open_accounts
: bytes
| acct_id
: int64 1
: bytes
| '0'
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
: null
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
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| v_order_totals
: bytes
| acct_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_order_totals
: bytes
| order_count
: int64 2
: bytes
| '0'
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
: null
> query
|
| 			SELECT
| 				tc.TABLE_NAME,
| 				cc.CONSTRAINT_NAME,
| 				cc.CHECK_CLAUSE
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| omni_dump_roundtrip
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
| ((`balance` >= 0) and (`status` in (_utf8mb4\'open\',_utf8mb4\'closed\')))
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
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_named_cols
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `id` from `omni_dump_roundtrip`.`accounts`
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
< row
: bytes
| v_order_totals
: bytes
| select `omni_dump_roundtrip`.`a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`omni_dump_roundtrip`.`v_open_accounts` `a` left join `omni_dump_roundtrip`.`orders` `o` on((`o`.`account_id` = `omni_dump_roundtrip`.`a`.`acct_id`))) group by `omni_dump_roundtrip`.`a`.`acct_id`
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols` (`renamed_id`) AS select `accounts`.`id` AS `id` from `accounts`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols_ref`
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
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_open_accounts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where (`accounts`.`status` = 'open')
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_order_totals`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_order_totals
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_order_totals` AS select `a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`v_open_accounts` `a` left join `orders` `o` on((`o`.`account_id` = `a`.`acct_id`))) group by `a`.`acct_id`
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
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| SYSTEM
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
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
< row
: bytes
| orders
: bytes
| fk_orders_parent
: bytes
| orders
: bytes
| SET NULL
: bytes
| NO ACTION
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
< column "COLUMN_NAME" "VARCHAR"
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
< row
: bytes
| omni_dump_roundtrip
: bytes
| orders
: bytes
| fk_orders_parent
: bytes
| parent_order_id
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
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< column "PARTITION_COMMENT" "TEXT" "string"
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
: bytes
|
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
: bytes
|
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
: bytes
|
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
: bytes
|
< row
: bytes
| metrics_range
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| 100
: bytes
| first
< row
: bytes
| metrics_range
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| 100
: bytes
| first
< row
: bytes
| metrics_range
: bytes
| p1
: bytes
| p1sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| MAXVALUE
: bytes
| last
< row
: bytes
| metrics_range
: bytes
| p1
: bytes
| p1sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| MAXVALUE
: bytes
| last
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_hash
: bytes
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_range
: bytes
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `region` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`,`region`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`ts`)
| SUBPARTITION BY HASH (`region`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (100) COMMENT = 'first' ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE COMMENT = 'last' ENGINE = InnoDB) */
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
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
< row
: bytes
| fn_format_tag
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| formats a tag
< row
: bytes
| sp_touch_order
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| reads order total
> query
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64) CHARSET utf8mb4
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
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "TABLE_TYPE" "ENUM" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR" "string"
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
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 131072
: bytes
| 0
: bytes
|
: bytes
| 'customer\'s accounts'
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
| utf8mb4_0900_ai_ci
: bytes
| 0
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
| utf8mb4_0900_ai_ci
: bytes
| 0
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
| orders
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
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
< row
: bytes
| v_order_totals
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
| utf8mb4_0900_ai_ci

== TestDumpRestores/lctn1 lctn1#1
> exec
| CREATE TABLE accounts (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	email VARCHAR(190) NOT NULL,
| 	bio TEXT,
| 	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
| 	flags SET('alpha', 'beta') DEFAULT 'alpha',
| 	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
| 	raw_token BINARY(8) NOT NULL DEFAULT 'ab',
| 	raw_bytes VARBINARY(16) NOT NULL DEFAULT 0x0100fe,
| 	secret_code VARCHAR(32) INVISIBLE DEFAULT 'hidden',
| 	meta JSON,
| 	coords POINT NOT NULL SRID 4326,
| 	geom GEOMETRY SRID 0,
| 	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL,
| 	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED,
| 	external_id VARCHAR(36) DEFAULT (UUID()),
| 	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
| 	note VARCHAR(100) COMMENT 'owner''s note \\ path',
| 	UNIQUE KEY uk_email (email),
| 	KEY idx_bio_prefix (bio(32)),
| 	KEY idx_created_desc (created_at DESC) COMMENT 'newest first',
| 	KEY idx_invisible (status) INVISIBLE,
| 	KEY idx_expr ((LOWER(status))),
| 	FULLTEXT KEY ft_bio (bio),
| 	SPATIAL INDEX sp_coords (coords),
| 	CONSTRAINT chk_balance CHECK (balance >= 0 AND status IN ('open', 'closed'))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='customer''s accounts';
< ok 0 0
> exec
| CREATE TABLE orders (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	account_id BIGINT NOT NULL,
| 	parent_order_id BIGINT DEFAULT NULL,
| 	total INT NOT NULL,
| 	KEY idx_account (account_id),
| 	KEY idx_parent (parent_order_id),
| 	CONSTRAINT fk_orders_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT,
| 	CONSTRAINT fk_orders_parent FOREIGN KEY (parent_order_id) REFERENCES orders (id) ON DELETE SET NULL ON UPDATE NO ACTION
| ) ENGINE=InnoDB;
< ok 0 0
> exec
| CREATE TABLE metrics_range (
| 	ts INT NOT NULL,
| 	region INT NOT NULL,
| 	val BIGINT NOT NULL,
| 	PRIMARY KEY (ts, region)
| ) PARTITION BY RANGE (ts) SUBPARTITION BY HASH (region) SUBPARTITIONS 2 (
| 	PARTITION p0 VALUES LESS THAN (100) COMMENT 'first',
| 	PARTITION p1 VALUES LESS THAN (MAXVALUE) COMMENT 'last'
| );
< ok 0 0
> exec
| CREATE TABLE metrics_hash (
| 	id INT NOT NULL PRIMARY KEY,
| 	val INT NOT NULL
| ) PARTITION BY HASH (id) PARTITIONS 4;
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
| CREATE VIEW v_order_totals AS
| 	SELECT a.acct_id, COUNT(o.id) AS order_count
| 	FROM v_open_accounts a
| 	LEFT JOIN orders o ON o.account_id = a.acct_id
| 	GROUP BY a.acct_id;
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
| 8.0.46
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				INDEX_NAME,
| 				COLUMN_NAME,
| 				COLLATION,
| 				IFNULL(SUB_PART, -1),
| 				EXPRESSION,
| 				SEQ_IN_INDEX,
| 				INDEX_TYPE,
| 				NON_UNIQUE,
| 				CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END,
| 				INDEX_COMMENT
| 			FROM information_schema.STATISTICS
| 			WHERE TABLE_SCHEMA = ?
| 			ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "INT" "int32"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT" "int32"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| ft_bio
: bytes
| bio
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
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
: null
: int64 1
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
| idx_created_desc
: bytes
| created_at
: bytes
| D
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
| newest first
< row
: bytes
| accounts
: bytes
| idx_expr
: null
: bytes
| A
: int64 -1
: bytes
| lower(`status`)
: int64 1
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
| idx_invisible
: bytes
| status
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 0
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| accounts
: bytes
| sp_coords
: bytes
| coords
: bytes
| A
: int64 32
: null
: int64 1
: bytes
| SPATIAL
: int64 1
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| metrics_range
: bytes
| PRIMARY
: bytes
| region
: bytes
| A
: int64 -1
: null
: int64 2
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
: null
: int64 1
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
| idx_parent
: bytes
| parent_order_id
: bytes
| A
: int64 -1
: null
: int64 1
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
: null
: int64 1
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
| 			SRS_ID
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED INT" "uint32"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB" "[]uint8"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
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
: null
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
| utf8mb4_0900_ai_ci
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
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
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
| utf8mb4_0900_ai_ci
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
| flags
: int64 5
: bytes
| 'alpha'
: bytes
| YES
: bytes
| set('alpha','beta')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
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
: null
< row
: bytes
| accounts
: bytes
| raw_token
: int64 7
: bytes
| '0x6162'
: bytes
| NO
: bytes
| binary(8)
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
| raw_bytes
: int64 8
: bytes
| '0x01'
: bytes
| NO
: bytes
| varbinary(16)
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
| secret_code
: int64 9
: bytes
| 'hidden'
: bytes
| YES
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
| INVISIBLE
: null
< row
: bytes
| accounts
: bytes
| meta
: int64 10
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
| coords
: int64 11
: null
: bytes
| NO
: bytes
| point
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
: int64 4326
< row
: bytes
| accounts
: bytes
| geom
: int64 12
: null
: bytes
| YES
: bytes
| geometry
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
: int64 0
< row
: bytes
| accounts
: bytes
| email_lower
: int64 13
: null
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| lower(`email`)
: bytes
| VIRTUAL GENERATED
: null
< row
: bytes
| accounts
: bytes
| display_tag
: int64 14
: null
: bytes
| YES
: bytes
| varchar(220)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| concat(`email`,_utf8mb4\':\',_utf8mb4\'v1\')
: bytes
| STORED GENERATED
: null
< row
: bytes
| accounts
: bytes
| external_id
: int64 15
: bytes
| 'uuid()'
: bytes
| YES
: bytes
| varchar(36)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| accounts
: bytes
| created_at
: int64 16
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
| DEFAULT_GENERATED
: null
< row
: bytes
| accounts
: bytes
| updated_at
: int64 17
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
: null
< row
: bytes
| accounts
: bytes
| note
: int64 18
: null
: bytes
| YES
: bytes
| varchar(100)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 'owner\'s note \\ path'
: bytes
|
: bytes
|
: null
< row
: bytes
| metrics_hash
: bytes
| id
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
: null
< row
: bytes
| metrics_hash
: bytes
| val
: int64 2
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
: null
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
: null
< row
: bytes
| metrics_range
: bytes
| region
: int64 2
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
: null
< row
: bytes
| metrics_range
: bytes
| val
: int64 3
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
: null
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
: null
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
: null
< row
: bytes
| orders
: bytes
| parent_order_id
: int64 3
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
: null
< row
: bytes
| orders
: bytes
| total
: int64 4
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
: null
< row
: bytes
| v_named_cols
: bytes
| renamed_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_named_cols_ref
: bytes
| renamed_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_open_accounts
: bytes
| acct_id
: int64 1
: bytes
| '0'
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
: null
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
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| v_order_totals
: bytes
| acct_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_order_totals
: bytes
| order_count
: int64 2
: bytes
| '0'
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
: null
> query
|
| 			SELECT
| 				tc.TABLE_NAME,
| 				cc.CONSTRAINT_NAME,
| 				cc.CHECK_CLAUSE
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| accounts
: bytes
| chk_balance
: bytes
| ((`balance` >= 0) and (`status` in (_utf8mb4\'open\',_utf8mb4\'closed\')))
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
| v_named_cols
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `id` from `omni_dump_roundtrip`.`accounts`
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
< row
: bytes
| v_order_totals
: bytes
| select `omni_dump_roundtrip`.`a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`omni_dump_roundtrip`.`v_open_accounts` `a` left join `omni_dump_roundtrip`.`orders` `o` on((`o`.`account_id` = `omni_dump_roundtrip`.`a`.`acct_id`))) group by `omni_dump_roundtrip`.`a`.`acct_id`
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols` (`renamed_id`) AS select `accounts`.`id` AS `id` from `accounts`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols_ref`
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
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_open_accounts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where (`accounts`.`status` = 'open')
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_order_totals`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_order_totals
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_order_totals` AS select `a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`v_open_accounts` `a` left join `orders` `o` on((`o`.`account_id` = `a`.`acct_id`))) group by `a`.`acct_id`
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
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| SYSTEM
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
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
< row
: bytes
| orders
: bytes
| fk_orders_parent
: bytes
| orders
: bytes
| SET NULL
: bytes
| NO ACTION
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
| orders
: bytes
| fk_orders_account
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
| fk_orders_parent
: bytes
| parent_order_id
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
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< column "PARTITION_COMMENT" "TEXT" "string"
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
: bytes
|
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
: bytes
|
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
: bytes
|
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
: bytes
|
< row
: bytes
| metrics_range
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| 100
: bytes
| first
< row
: bytes
| metrics_range
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| 100
: bytes
| first
< row
: bytes
| metrics_range
: bytes
| p1
: bytes
| p1sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| MAXVALUE
: bytes
| last
< row
: bytes
| metrics_range
: bytes
| p1
: bytes
| p1sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| MAXVALUE
: bytes
| last
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_hash
: bytes
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_range
: bytes
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `region` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`,`region`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`ts`)
| SUBPARTITION BY HASH (`region`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (100) COMMENT = 'first' ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE COMMENT = 'last' ENGINE = InnoDB) */
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
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
< row
: bytes
| fn_format_tag
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| formats a tag
< row
: bytes
| sp_touch_order
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| reads order total
> query
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
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
< column "TABLE_TYPE" "ENUM" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR" "string"
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
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 131072
: bytes
| 0
: bytes
|
: bytes
| 'customer\'s accounts'
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
| utf8mb4_0900_ai_ci
: bytes
| 0
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
| utf8mb4_0900_ai_ci
: bytes
| 0
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
| orders
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
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
< row
: bytes
| v_order_totals
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
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE `omni_dump_roundtrip`
< ok 8 0
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
| CREATE VIEW `v_order_totals` AS SELECT
|   1 AS `acct_id`,
|   1 AS `order_count`
< ok 0 0
> exec
| CREATE TABLE `accounts` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `email` varchar(190) NOT NULL,
|   `bio` text DEFAULT NULL,
|   `status` enum('open','closed') NOT NULL DEFAULT 'open',
|   `flags` set('alpha','beta') DEFAULT 'alpha',
|   `balance` decimal(12,4) NOT NULL DEFAULT '0.0000',
|   `raw_token` binary(8) NOT NULL DEFAULT 'ab',
|   `raw_bytes` varbinary(16) NOT NULL DEFAULT 0x01,
|   `secret_code` varchar(32) DEFAULT 'hidden' /*!80023 INVISIBLE */,
|   `meta` json,
|   `coords` point NOT NULL /*!80003 SRID 4326 */,
|   `geom` geometry /*!80003 SRID 0 */,
|   `email_lower` varchar(190) GENERATED ALWAYS AS (lower(`email`)) VIRTUAL,
|   `display_tag` varchar(220) GENERATED ALWAYS AS (concat(`email`,_utf8mb4':',_utf8mb4'v1')) STORED,
|   `external_id` varchar(36) DEFAULT (uuid()),
|   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
|   `note` varchar(100) DEFAULT NULL COMMENT 'owner\'s note \\ path',
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ft_bio` (`bio`),
|   KEY `idx_bio_prefix` (`bio`(32)),
|   KEY `idx_created_desc` (`created_at` DESC) COMMENT 'newest first',
|   KEY `idx_expr` ((lower(`status`))),
|   KEY `idx_invisible` (`status`) /*!80000 INVISIBLE */,
|   SPATIAL KEY `sp_coords` (`coords`),
|   UNIQUE KEY `uk_email` (`email`),
|   CONSTRAINT `chk_balance` CHECK ((`balance` >= 0) and (`status` in (_utf8mb4'open',_utf8mb4'closed')))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='customer\'s accounts'
< ok 0 0
> exec
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
< ok 0 0
> exec
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `region` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`, `region`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`ts`)
| SUBPARTITION BY HASH (`region`)
| SUBPARTITIONS 2
| (PARTITION `p0` VALUES LESS THAN (100) ENGINE=InnoDB COMMENT = 'first',
|  PARTITION `p1` VALUES LESS THAN MAXVALUE ENGINE=InnoDB COMMENT = 'last') */
< ok 0 0
> exec
| CREATE TABLE `orders` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `account_id` bigint NOT NULL,
|   `parent_order_id` bigint DEFAULT NULL,
|   `total` int NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_account` (`account_id`),
|   KEY `idx_parent` (`parent_order_id`),
|   CONSTRAINT `fk_orders_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
|   CONSTRAINT `fk_orders_parent` FOREIGN KEY (`parent_order_id`) REFERENCES `orders` (`id`) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
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
| CREATE OR REPLACE VIEW `v_named_cols` (`renamed_id`) AS select `accounts`.`id` AS `id` from `accounts`
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_named_cols_ref` AS select `v_named_cols`.`renamed_id` AS `renamed_id` from `v_named_cols`
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where (`accounts`.`status` = 'open')
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_order_totals` AS select `a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`v_open_accounts` `a` left join `orders` `o` on((`o`.`account_id` = `a`.`acct_id`))) group by `a`.`acct_id`
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
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
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
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
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
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
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
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR" "string"
< row
: bytes
| 8.0.46
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				INDEX_NAME,
| 				COLUMN_NAME,
| 				COLLATION,
| 				IFNULL(SUB_PART, -1),
| 				EXPRESSION,
| 				SEQ_IN_INDEX,
| 				INDEX_TYPE,
| 				NON_UNIQUE,
| 				CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END,
| 				INDEX_COMMENT
| 			FROM information_schema.STATISTICS
| 			WHERE TABLE_SCHEMA = ?
| 			ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "INT" "int32"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT" "int32"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| ft_bio
: bytes
| bio
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
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
: null
: int64 1
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
| idx_created_desc
: bytes
| created_at
: bytes
| D
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
| newest first
< row
: bytes
| accounts
: bytes
| idx_expr
: null
: bytes
| A
: int64 -1
: bytes
| lower(`status`)
: int64 1
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
| idx_invisible
: bytes
| status
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 0
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| accounts
: bytes
| sp_coords
: bytes
| coords
: bytes
| A
: int64 32
: null
: int64 1
: bytes
| SPATIAL
: int64 1
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| metrics_range
: bytes
| PRIMARY
: bytes
| region
: bytes
| A
: int64 -1
: null
: int64 2
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
: null
: int64 1
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
| idx_parent
: bytes
| parent_order_id
: bytes
| A
: int64 -1
: null
: int64 1
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
: null
: int64 1
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
| 			SRS_ID
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED INT" "uint32"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB" "[]uint8"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
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
: null
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
| utf8mb4_0900_ai_ci
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
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
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
| utf8mb4_0900_ai_ci
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
| flags
: int64 5
: bytes
| 'alpha'
: bytes
| YES
: bytes
| set('alpha','beta')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
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
: null
< row
: bytes
| accounts
: bytes
| raw_token
: int64 7
: bytes
| '0x6162'
: bytes
| NO
: bytes
| binary(8)
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
| raw_bytes
: int64 8
: bytes
| '0x01'
: bytes
| NO
: bytes
| varbinary(16)
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
| secret_code
: int64 9
: bytes
| 'hidden'
: bytes
| YES
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
| INVISIBLE
: null
< row
: bytes
| accounts
: bytes
| meta
: int64 10
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
| coords
: int64 11
: null
: bytes
| NO
: bytes
| point
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
: int64 4326
< row
: bytes
| accounts
: bytes
| geom
: int64 12
: null
: bytes
| YES
: bytes
| geometry
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
: int64 0
< row
: bytes
| accounts
: bytes
| email_lower
: int64 13
: null
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| lower(`email`)
: bytes
| VIRTUAL GENERATED
: null
< row
: bytes
| accounts
: bytes
| display_tag
: int64 14
: null
: bytes
| YES
: bytes
| varchar(220)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| concat(`email`,_utf8mb4\':\',_utf8mb4\'v1\')
: bytes
| STORED GENERATED
: null
< row
: bytes
| accounts
: bytes
| external_id
: int64 15
: bytes
| 'uuid()'
: bytes
| YES
: bytes
| varchar(36)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| accounts
: bytes
| created_at
: int64 16
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
| DEFAULT_GENERATED
: null
< row
: bytes
| accounts
: bytes
| updated_at
: int64 17
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
: null
< row
: bytes
| accounts
: bytes
| note
: int64 18
: null
: bytes
| YES
: bytes
| varchar(100)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 'owner\'s note \\ path'
: bytes
|
: bytes
|
: null
< row
: bytes
| metrics_hash
: bytes
| id
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
: null
< row
: bytes
| metrics_hash
: bytes
| val
: int64 2
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
: null
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
: null
< row
: bytes
| metrics_range
: bytes
| region
: int64 2
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
: null
< row
: bytes
| metrics_range
: bytes
| val
: int64 3
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
: null
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
: null
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
: null
< row
: bytes
| orders
: bytes
| parent_order_id
: int64 3
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
: null
< row
: bytes
| orders
: bytes
| total
: int64 4
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
: null
< row
: bytes
| v_named_cols
: bytes
| renamed_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_named_cols_ref
: bytes
| renamed_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_open_accounts
: bytes
| acct_id
: int64 1
: bytes
| '0'
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
: null
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
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| v_order_totals
: bytes
| acct_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_order_totals
: bytes
| order_count
: int64 2
: bytes
| '0'
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
: null
> query
|
| 			SELECT
| 				tc.TABLE_NAME,
| 				cc.CONSTRAINT_NAME,
| 				cc.CHECK_CLAUSE
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| accounts
: bytes
| chk_balance
: bytes
| ((`balance` >= 0) and (`status` in (_utf8mb4\'open\',_utf8mb4\'closed\')))
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
| v_named_cols
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `id` from `omni_dump_roundtrip`.`accounts`
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
< row
: bytes
| v_order_totals
: bytes
| select `omni_dump_roundtrip`.`a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`omni_dump_roundtrip`.`v_open_accounts` `a` left join `omni_dump_roundtrip`.`orders` `o` on((`o`.`account_id` = `omni_dump_roundtrip`.`a`.`acct_id`))) group by `omni_dump_roundtrip`.`a`.`acct_id`
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols` (`renamed_id`) AS select `accounts`.`id` AS `id` from `accounts`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols_ref`
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
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_open_accounts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where (`accounts`.`status` = 'open')
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_order_totals`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_order_totals
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_order_totals` AS select `a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`v_open_accounts` `a` left join `orders` `o` on((`o`.`account_id` = `a`.`acct_id`))) group by `a`.`acct_id`
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
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| SYSTEM
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
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
< row
: bytes
| orders
: bytes
| fk_orders_parent
: bytes
| orders
: bytes
| SET NULL
: bytes
| NO ACTION
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
| orders
: bytes
| fk_orders_account
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
| fk_orders_parent
: bytes
| parent_order_id
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
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< column "PARTITION_COMMENT" "TEXT" "string"
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
: bytes
|
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
: bytes
|
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
: bytes
|
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
: bytes
|
< row
: bytes
| metrics_range
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| 100
: bytes
| first
< row
: bytes
| metrics_range
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| 100
: bytes
| first
< row
: bytes
| metrics_range
: bytes
| p1
: bytes
| p1sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| MAXVALUE
: bytes
| last
< row
: bytes
| metrics_range
: bytes
| p1
: bytes
| p1sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| MAXVALUE
: bytes
| last
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_hash
: bytes
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_range
: bytes
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `region` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`,`region`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`ts`)
| SUBPARTITION BY HASH (`region`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (100) COMMENT = 'first' ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE COMMENT = 'last' ENGINE = InnoDB) */
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
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
< row
: bytes
| fn_format_tag
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| formats a tag
< row
: bytes
| sp_touch_order
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| reads order total
> query
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64) CHARSET utf8mb4
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
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
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
< column "TABLE_TYPE" "ENUM" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR" "string"
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
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 131072
: bytes
| 0
: bytes
|
: bytes
| 'customer\'s accounts'
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
| utf8mb4_0900_ai_ci
: bytes
| 0
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
| utf8mb4_0900_ai_ci
: bytes
| 0
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
| orders
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
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
< row
: bytes
| v_order_totals
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
| utf8mb4_0900_ai_ci

== TestDumpRestores/mysql84 mysql84#1
> exec
| CREATE TABLE accounts (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	email VARCHAR(190) NOT NULL,
| 	bio TEXT,
| 	status ENUM('open', 'closed') NOT NULL DEFAULT 'open',
| 	flags SET('alpha', 'beta') DEFAULT 'alpha',
| 	balance DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
| 	raw_token BINARY(8) NOT NULL DEFAULT 'ab',
| 	raw_bytes VARBINARY(16) NOT NULL DEFAULT 0x0100fe,
| 	secret_code VARCHAR(32) INVISIBLE DEFAULT 'hidden',
| 	meta JSON,
| 	coords POINT NOT NULL SRID 4326,
| 	geom GEOMETRY SRID 0,
| 	email_lower VARCHAR(190) GENERATED ALWAYS AS (LOWER(email)) VIRTUAL,
| 	display_tag VARCHAR(220) GENERATED ALWAYS AS (CONCAT(email, ':', 'v1')) STORED,
| 	external_id VARCHAR(36) DEFAULT (UUID()),
| 	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
| 	note VARCHAR(100) COMMENT 'owner''s note \\ path',
| 	UNIQUE KEY uk_email (email),
| 	KEY idx_bio_prefix (bio(32)),
| 	KEY idx_created_desc (created_at DESC) COMMENT 'newest first',
| 	KEY idx_invisible (status) INVISIBLE,
| 	KEY idx_expr ((LOWER(status))),
| 	FULLTEXT KEY ft_bio (bio),
| 	SPATIAL INDEX sp_coords (coords),
| 	CONSTRAINT chk_balance CHECK (balance >= 0 AND status IN ('open', 'closed'))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='customer''s accounts';
< ok 0 0
> exec
| CREATE TABLE orders (
| 	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
| 	account_id BIGINT NOT NULL,
| 	parent_order_id BIGINT DEFAULT NULL,
| 	total INT NOT NULL,
| 	KEY idx_account (account_id),
| 	KEY idx_parent (parent_order_id),
| 	CONSTRAINT fk_orders_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE RESTRICT,
| 	CONSTRAINT fk_orders_parent FOREIGN KEY (parent_order_id) REFERENCES orders (id) ON DELETE SET NULL ON UPDATE NO ACTION
| ) ENGINE=InnoDB;
< ok 0 0
> exec
| CREATE TABLE metrics_range (
| 	ts INT NOT NULL,
| 	region INT NOT NULL,
| 	val BIGINT NOT NULL,
| 	PRIMARY KEY (ts, region)
| ) PARTITION BY RANGE (ts) SUBPARTITION BY HASH (region) SUBPARTITIONS 2 (
| 	PARTITION p0 VALUES LESS THAN (100) COMMENT 'first',
| 	PARTITION p1 VALUES LESS THAN (MAXVALUE) COMMENT 'last'
| );
< ok 0 0
> exec
| CREATE TABLE metrics_hash (
| 	id INT NOT NULL PRIMARY KEY,
| 	val INT NOT NULL
| ) PARTITION BY HASH (id) PARTITIONS 4;
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
| CREATE VIEW v_order_totals AS
| 	SELECT a.acct_id, COUNT(o.id) AS order_count
| 	FROM v_open_accounts a
| 	LEFT JOIN orders o ON o.account_id = a.acct_id
| 	GROUP BY a.acct_id;
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
| 8.4.13
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				INDEX_NAME,
| 				COLUMN_NAME,
| 				COLLATION,
| 				IFNULL(SUB_PART, -1),
| 				EXPRESSION,
| 				SEQ_IN_INDEX,
| 				INDEX_TYPE,
| 				NON_UNIQUE,
| 				CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END,
| 				INDEX_COMMENT
| 			FROM information_schema.STATISTICS
| 			WHERE TABLE_SCHEMA = ?
| 			ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "INT" "int32"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT" "int32"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| ft_bio
: bytes
| bio
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
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
: null
: int64 1
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
| idx_created_desc
: bytes
| created_at
: bytes
| D
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
| newest first
< row
: bytes
| accounts
: bytes
| idx_expr
: null
: bytes
| A
: int64 -1
: bytes
| lower(`status`)
: int64 1
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
| idx_invisible
: bytes
| status
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 0
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| accounts
: bytes
| sp_coords
: bytes
| coords
: bytes
| A
: int64 32
: null
: int64 1
: bytes
| SPATIAL
: int64 1
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| metrics_range
: bytes
| PRIMARY
: bytes
| region
: bytes
| A
: int64 -1
: null
: int64 2
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
: null
: int64 1
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
| idx_parent
: bytes
| parent_order_id
: bytes
| A
: int64 -1
: null
: int64 1
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
: null
: int64 1
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
| 			SRS_ID
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED INT" "uint32"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB" "[]uint8"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
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
: null
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
| utf8mb4_0900_ai_ci
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
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
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
| utf8mb4_0900_ai_ci
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
| flags
: int64 5
: bytes
| 'alpha'
: bytes
| YES
: bytes
| set('alpha','beta')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
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
: null
< row
: bytes
| accounts
: bytes
| raw_token
: int64 7
: bytes
| '0x6162'
: bytes
| NO
: bytes
| binary(8)
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
| raw_bytes
: int64 8
: bytes
| '0x01'
: bytes
| NO
: bytes
| varbinary(16)
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
| secret_code
: int64 9
: bytes
| 'hidden'
: bytes
| YES
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
| INVISIBLE
: null
< row
: bytes
| accounts
: bytes
| meta
: int64 10
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
| coords
: int64 11
: null
: bytes
| NO
: bytes
| point
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
: int64 4326
< row
: bytes
| accounts
: bytes
| geom
: int64 12
: null
: bytes
| YES
: bytes
| geometry
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
: int64 0
< row
: bytes
| accounts
: bytes
| email_lower
: int64 13
: null
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| lower(`email`)
: bytes
| VIRTUAL GENERATED
: null
< row
: bytes
| accounts
: bytes
| display_tag
: int64 14
: null
: bytes
| YES
: bytes
| varchar(220)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| concat(`email`,_utf8mb4\':\',_utf8mb4\'v1\')
: bytes
| STORED GENERATED
: null
< row
: bytes
| accounts
: bytes
| external_id
: int64 15
: bytes
| 'uuid()'
: bytes
| YES
: bytes
| varchar(36)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| accounts
: bytes
| created_at
: int64 16
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
| DEFAULT_GENERATED
: null
< row
: bytes
| accounts
: bytes
| updated_at
: int64 17
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
: null
< row
: bytes
| accounts
: bytes
| note
: int64 18
: null
: bytes
| YES
: bytes
| varchar(100)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 'owner\'s note \\ path'
: bytes
|
: bytes
|
: null
< row
: bytes
| metrics_hash
: bytes
| id
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
: null
< row
: bytes
| metrics_hash
: bytes
| val
: int64 2
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
: null
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
: null
< row
: bytes
| metrics_range
: bytes
| region
: int64 2
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
: null
< row
: bytes
| metrics_range
: bytes
| val
: int64 3
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
: null
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
: null
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
: null
< row
: bytes
| orders
: bytes
| parent_order_id
: int64 3
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
: null
< row
: bytes
| orders
: bytes
| total
: int64 4
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
: null
< row
: bytes
| v_named_cols
: bytes
| renamed_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_named_cols_ref
: bytes
| renamed_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_open_accounts
: bytes
| acct_id
: int64 1
: bytes
| '0'
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
: null
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
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| v_order_totals
: bytes
| acct_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_order_totals
: bytes
| order_count
: int64 2
: bytes
| '0'
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
: null
> query
|
| 			SELECT
| 				tc.TABLE_NAME,
| 				cc.CONSTRAINT_NAME,
| 				cc.CHECK_CLAUSE
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| omni_dump_roundtrip
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
| ((`balance` >= 0) and (`status` in (_utf8mb4\'open\',_utf8mb4\'closed\')))
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
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_named_cols
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `id` from `omni_dump_roundtrip`.`accounts`
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
< row
: bytes
| v_order_totals
: bytes
| select `omni_dump_roundtrip`.`a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`omni_dump_roundtrip`.`v_open_accounts` `a` left join `omni_dump_roundtrip`.`orders` `o` on((`o`.`account_id` = `omni_dump_roundtrip`.`a`.`acct_id`))) group by `omni_dump_roundtrip`.`a`.`acct_id`
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols` (`renamed_id`) AS select `accounts`.`id` AS `id` from `accounts`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols_ref`
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
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_open_accounts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where (`accounts`.`status` = 'open')
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_order_totals`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_order_totals
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_order_totals` AS select `a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`v_open_accounts` `a` left join `orders` `o` on((`o`.`account_id` = `a`.`acct_id`))) group by `a`.`acct_id`
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
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| SYSTEM
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
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
< row
: bytes
| orders
: bytes
| fk_orders_parent
: bytes
| orders
: bytes
| SET NULL
: bytes
| NO ACTION
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
< column "COLUMN_NAME" "VARCHAR"
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
< row
: bytes
| omni_dump_roundtrip
: bytes
| orders
: bytes
| fk_orders_parent
: bytes
| parent_order_id
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
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< column "PARTITION_COMMENT" "TEXT" "string"
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
: bytes
|
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
: bytes
|
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
: bytes
|
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
: bytes
|
< row
: bytes
| metrics_range
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| 100
: bytes
| first
< row
: bytes
| metrics_range
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| 100
: bytes
| first
< row
: bytes
| metrics_range
: bytes
| p1
: bytes
| p1sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| MAXVALUE
: bytes
| last
< row
: bytes
| metrics_range
: bytes
| p1
: bytes
| p1sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| MAXVALUE
: bytes
| last
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_hash
: bytes
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_range
: bytes
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `region` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`,`region`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`ts`)
| SUBPARTITION BY HASH (`region`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (100) COMMENT = 'first' ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE COMMENT = 'last' ENGINE = InnoDB) */
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
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
< row
: bytes
| fn_format_tag
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| formats a tag
< row
: bytes
| sp_touch_order
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| reads order total
> query
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "TABLE_TYPE" "ENUM" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR" "string"
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
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 131072
: bytes
| 0
: bytes
|
: bytes
| 'customer\'s accounts'
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
| utf8mb4_0900_ai_ci
: bytes
| 0
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
| utf8mb4_0900_ai_ci
: bytes
| 0
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
| orders
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
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
< row
: bytes
| v_order_totals
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
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE `omni_dump_roundtrip`
< ok 8 0
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
| CREATE VIEW `v_order_totals` AS SELECT
|   1 AS `acct_id`,
|   1 AS `order_count`
< ok 0 0
> exec
| CREATE TABLE `accounts` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `email` varchar(190) NOT NULL,
|   `bio` text DEFAULT NULL,
|   `status` enum('open','closed') NOT NULL DEFAULT 'open',
|   `flags` set('alpha','beta') DEFAULT 'alpha',
|   `balance` decimal(12,4) NOT NULL DEFAULT '0.0000',
|   `raw_token` binary(8) NOT NULL DEFAULT 'ab',
|   `raw_bytes` varbinary(16) NOT NULL DEFAULT 0x01,
|   `secret_code` varchar(32) DEFAULT 'hidden' /*!80023 INVISIBLE */,
|   `meta` json,
|   `coords` point NOT NULL /*!80003 SRID 4326 */,
|   `geom` geometry /*!80003 SRID 0 */,
|   `email_lower` varchar(190) GENERATED ALWAYS AS (lower(`email`)) VIRTUAL,
|   `display_tag` varchar(220) GENERATED ALWAYS AS (concat(`email`,_utf8mb4':',_utf8mb4'v1')) STORED,
|   `external_id` varchar(36) DEFAULT (uuid()),
|   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
|   `note` varchar(100) DEFAULT NULL COMMENT 'owner\'s note \\ path',
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ft_bio` (`bio`),
|   KEY `idx_bio_prefix` (`bio`(32)),
|   KEY `idx_created_desc` (`created_at` DESC) COMMENT 'newest first',
|   KEY `idx_expr` ((lower(`status`))),
|   KEY `idx_invisible` (`status`) /*!80000 INVISIBLE */,
|   SPATIAL KEY `sp_coords` (`coords`),
|   UNIQUE KEY `uk_email` (`email`),
|   CONSTRAINT `chk_balance` CHECK ((`balance` >= 0) and (`status` in (_utf8mb4'open',_utf8mb4'closed')))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='customer\'s accounts'
< ok 0 0
> exec
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
< ok 0 0
> exec
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `region` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`, `region`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`ts`)
| SUBPARTITION BY HASH (`region`)
| SUBPARTITIONS 2
| (PARTITION `p0` VALUES LESS THAN (100) ENGINE=InnoDB COMMENT = 'first',
|  PARTITION `p1` VALUES LESS THAN MAXVALUE ENGINE=InnoDB COMMENT = 'last') */
< ok 0 0
> exec
| CREATE TABLE `orders` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `account_id` bigint NOT NULL,
|   `parent_order_id` bigint DEFAULT NULL,
|   `total` int NOT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_account` (`account_id`),
|   KEY `idx_parent` (`parent_order_id`),
|   CONSTRAINT `fk_orders_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
|   CONSTRAINT `fk_orders_parent` FOREIGN KEY (`parent_order_id`) REFERENCES `orders` (`id`) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
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
| CREATE OR REPLACE VIEW `v_named_cols` (`renamed_id`) AS select `accounts`.`id` AS `id` from `accounts`
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_named_cols_ref` AS select `v_named_cols`.`renamed_id` AS `renamed_id` from `v_named_cols`
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where (`accounts`.`status` = 'open')
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_order_totals` AS select `a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`v_open_accounts` `a` left join `orders` `o` on((`o`.`account_id` = `a`.`acct_id`))) group by `a`.`acct_id`
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
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
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
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
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
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
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
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR" "string"
< row
: bytes
| 8.4.13
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				INDEX_NAME,
| 				COLUMN_NAME,
| 				COLLATION,
| 				IFNULL(SUB_PART, -1),
| 				EXPRESSION,
| 				SEQ_IN_INDEX,
| 				INDEX_TYPE,
| 				NON_UNIQUE,
| 				CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END,
| 				INDEX_COMMENT
| 			FROM information_schema.STATISTICS
| 			WHERE TABLE_SCHEMA = ?
| 			ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "INT" "int32"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT" "int32"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| accounts
: bytes
| ft_bio
: bytes
| bio
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
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
: null
: int64 1
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
| idx_created_desc
: bytes
| created_at
: bytes
| D
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
| newest first
< row
: bytes
| accounts
: bytes
| idx_expr
: null
: bytes
| A
: int64 -1
: bytes
| lower(`status`)
: int64 1
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
| idx_invisible
: bytes
| status
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 0
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| accounts
: bytes
| sp_coords
: bytes
| coords
: bytes
| A
: int64 32
: null
: int64 1
: bytes
| SPATIAL
: int64 1
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
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
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| metrics_range
: bytes
| PRIMARY
: bytes
| region
: bytes
| A
: int64 -1
: null
: int64 2
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
: null
: int64 1
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
| idx_parent
: bytes
| parent_order_id
: bytes
| A
: int64 -1
: null
: int64 1
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
: null
: int64 1
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
| 			SRS_ID
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED INT" "uint32"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB" "[]uint8"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
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
: null
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
| utf8mb4_0900_ai_ci
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
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
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
| utf8mb4_0900_ai_ci
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
| flags
: int64 5
: bytes
| 'alpha'
: bytes
| YES
: bytes
| set('alpha','beta')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
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
: null
< row
: bytes
| accounts
: bytes
| raw_token
: int64 7
: bytes
| '0x6162'
: bytes
| NO
: bytes
| binary(8)
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
| raw_bytes
: int64 8
: bytes
| '0x01'
: bytes
| NO
: bytes
| varbinary(16)
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
| secret_code
: int64 9
: bytes
| 'hidden'
: bytes
| YES
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
| INVISIBLE
: null
< row
: bytes
| accounts
: bytes
| meta
: int64 10
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
| coords
: int64 11
: null
: bytes
| NO
: bytes
| point
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
: int64 4326
< row
: bytes
| accounts
: bytes
| geom
: int64 12
: null
: bytes
| YES
: bytes
| geometry
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
: int64 0
< row
: bytes
| accounts
: bytes
| email_lower
: int64 13
: null
: bytes
| YES
: bytes
| varchar(190)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| lower(`email`)
: bytes
| VIRTUAL GENERATED
: null
< row
: bytes
| accounts
: bytes
| display_tag
: int64 14
: null
: bytes
| YES
: bytes
| varchar(220)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| concat(`email`,_utf8mb4\':\',_utf8mb4\'v1\')
: bytes
| STORED GENERATED
: null
< row
: bytes
| accounts
: bytes
| external_id
: int64 15
: bytes
| 'uuid()'
: bytes
| YES
: bytes
| varchar(36)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| accounts
: bytes
| created_at
: int64 16
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
| DEFAULT_GENERATED
: null
< row
: bytes
| accounts
: bytes
| updated_at
: int64 17
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
: null
< row
: bytes
| accounts
: bytes
| note
: int64 18
: null
: bytes
| YES
: bytes
| varchar(100)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 'owner\'s note \\ path'
: bytes
|
: bytes
|
: null
< row
: bytes
| metrics_hash
: bytes
| id
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
: null
< row
: bytes
| metrics_hash
: bytes
| val
: int64 2
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
: null
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
: null
< row
: bytes
| metrics_range
: bytes
| region
: int64 2
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
: null
< row
: bytes
| metrics_range
: bytes
| val
: int64 3
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
: null
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
: null
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
: null
< row
: bytes
| orders
: bytes
| parent_order_id
: int64 3
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
: null
< row
: bytes
| orders
: bytes
| total
: int64 4
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
: null
< row
: bytes
| v_named_cols
: bytes
| renamed_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_named_cols_ref
: bytes
| renamed_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_open_accounts
: bytes
| acct_id
: int64 1
: bytes
| '0'
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
: null
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
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| v_order_totals
: bytes
| acct_id
: int64 1
: bytes
| '0'
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
: null
< row
: bytes
| v_order_totals
: bytes
| order_count
: int64 2
: bytes
| '0'
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
: null
> query
|
| 			SELECT
| 				tc.TABLE_NAME,
| 				cc.CONSTRAINT_NAME,
| 				cc.CHECK_CLAUSE
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| omni_dump_roundtrip
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
| ((`balance` >= 0) and (`status` in (_utf8mb4\'open\',_utf8mb4\'closed\')))
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
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_named_cols
: bytes
| select `omni_dump_roundtrip`.`accounts`.`id` AS `id` from `omni_dump_roundtrip`.`accounts`
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
< row
: bytes
| v_order_totals
: bytes
| select `omni_dump_roundtrip`.`a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`omni_dump_roundtrip`.`v_open_accounts` `a` left join `omni_dump_roundtrip`.`orders` `o` on((`o`.`account_id` = `omni_dump_roundtrip`.`a`.`acct_id`))) group by `omni_dump_roundtrip`.`a`.`acct_id`
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_named_cols
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_named_cols` (`renamed_id`) AS select `accounts`.`id` AS `id` from `accounts`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_named_cols_ref`
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
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_open_accounts`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_open_accounts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_open_accounts` AS select `accounts`.`id` AS `acct_id`,`accounts`.`email` AS `acct_email` from `accounts` where (`accounts`.`status` = 'open')
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`v_order_totals`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_order_totals
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_order_totals` AS select `a`.`acct_id` AS `acct_id`,count(`o`.`id`) AS `order_count` from (`v_open_accounts` `a` left join `orders` `o` on((`o`.`account_id` = `a`.`acct_id`))) group by `a`.`acct_id`
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
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ev_daily_check
: bytes
| SYSTEM
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE EVENT `omni_dump_roundtrip`.`ev_daily_check`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
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
< row
: bytes
| orders
: bytes
| fk_orders_parent
: bytes
| orders
: bytes
| SET NULL
: bytes
| NO ACTION
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
< column "COLUMN_NAME" "VARCHAR"
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
< row
: bytes
| omni_dump_roundtrip
: bytes
| orders
: bytes
| fk_orders_parent
: bytes
| parent_order_id
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
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< column "PARTITION_COMMENT" "TEXT" "string"
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
: bytes
|
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
: bytes
|
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
: bytes
|
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
: bytes
|
< row
: bytes
| metrics_range
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| 100
: bytes
| first
< row
: bytes
| metrics_range
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| 100
: bytes
| first
< row
: bytes
| metrics_range
: bytes
| p1
: bytes
| p1sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| MAXVALUE
: bytes
| last
< row
: bytes
| metrics_range
: bytes
| p1
: bytes
| p1sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `ts`
: bytes
| `region`
: bytes
| MAXVALUE
: bytes
| last
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_hash`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_hash
: bytes
| CREATE TABLE `metrics_hash` (
|   `id` int NOT NULL,
|   `val` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 4 */
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`metrics_range`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metrics_range
: bytes
| CREATE TABLE `metrics_range` (
|   `ts` int NOT NULL,
|   `region` int NOT NULL,
|   `val` bigint NOT NULL,
|   PRIMARY KEY (`ts`,`region`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`ts`)
| SUBPARTITION BY HASH (`region`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (100) COMMENT = 'first' ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE COMMENT = 'last' ENGINE = InnoDB) */
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
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
< row
: bytes
| fn_format_tag
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| formats a tag
< row
: bytes
| sp_touch_order
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| reads order total
> query
| SHOW CREATE FUNCTION `omni_dump_roundtrip`.`fn_format_tag`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_format_tag`(p_prefix VARCHAR(16) CHARSET utf8mb4, p_id BIGINT) RETURNS varchar(64) CHARSET utf8mb4
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
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`sp_touch_order`
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
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
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
| utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR" "string"
< column "TABLE_TYPE" "ENUM" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR" "string"
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
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 131072
: bytes
| 0
: bytes
|
: bytes
| 'customer\'s accounts'
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
| utf8mb4_0900_ai_ci
: bytes
| 0
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
| utf8mb4_0900_ai_ci
: bytes
| 0
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
| orders
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
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
< row
: bytes
| v_order_totals
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
| utf8mb4_0900_ai_ci
