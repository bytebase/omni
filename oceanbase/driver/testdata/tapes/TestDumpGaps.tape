-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDumpGaps/datetime(6)_default_current_timestamp_loses_precision default#1
> exec
| CREATE TABLE dt6_t (id INT PRIMARY KEY, ts DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6));
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
| omni_gap_dt6
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
| dt6_t
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
| omni_gap_dt6
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
| dt6_t
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
: bytes
|
: bytes
|
: null
< row
: bytes
| dt6_t
: bytes
| ts
: uint64 2
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
| omni_gap_dt6
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
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
| omni_gap_dt6
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
| omni_gap_dt6
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
| omni_gap_dt6
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_gap_dt6
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| omni_gap_dt6
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "TEXT"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT"
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
| omni_gap_dt6
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
| dt6_t
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
> query
|
| 		SELECT MVIEW_NAME, CONTAINER_NAME, QUERY
| 		FROM oceanbase.DBA_MVIEWS
| 		WHERE OWNER = ?
| 		ORDER BY MVIEW_NAME
: string
| omni_gap_dt6
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
| omni_gap_dt6
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| USE `omni_gap_dt6_restored`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `dt6_t` (
|   `id` int NOT NULL,
|   `ts` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< error
| 1067 42000 Invalid default value for 'ts'

== TestDumpGaps/event default#1
> exec
| CREATE TABLE ev_t (id INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE EVENT ev_1 ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DISABLE DO DELETE FROM ev_t;
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
| omni_gap_ev
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
| ev_t
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
| omni_gap_ev
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
| ev_t
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
| omni_gap_ev
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
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
| omni_gap_ev
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
| omni_gap_ev
< result
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| ev_1
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
| omni_gap_ev
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_gap_ev
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| omni_gap_ev
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "TEXT"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT"
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
| omni_gap_ev
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
| ev_t
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
> query
|
| 		SELECT MVIEW_NAME, CONTAINER_NAME, QUERY
| 		FROM oceanbase.DBA_MVIEWS
| 		WHERE OWNER = ?
| 		ORDER BY MVIEW_NAME
: string
| omni_gap_ev
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
| omni_gap_ev
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| USE `omni_gap_ev_restored`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `ev_t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
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
| omni_gap_ev_restored
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
| ev_t
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
| omni_gap_ev_restored
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
| ev_t
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
| omni_gap_ev_restored
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
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
| omni_gap_ev_restored
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
| omni_gap_ev_restored
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
| omni_gap_ev_restored
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_gap_ev_restored
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| omni_gap_ev_restored
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "TEXT"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT"
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
| omni_gap_ev_restored
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
| ev_t
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
> query
|
| 		SELECT MVIEW_NAME, CONTAINER_NAME, QUERY
| 		FROM oceanbase.DBA_MVIEWS
| 		WHERE OWNER = ?
| 		ORDER BY MVIEW_NAME
: string
| omni_gap_ev_restored
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
| omni_gap_ev_restored
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestDumpGaps/generated_column_carries_COLUMN_DEFAULT default#1
> exec
| CREATE TABLE gen_t (id INT PRIMARY KEY, a VARCHAR(32), b VARCHAR(32) GENERATED ALWAYS AS (LOWER(a)) VIRTUAL);
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
| omni_gap_gen
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
| gen_t
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
| omni_gap_gen
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
| gen_t
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
: bytes
|
: bytes
|
: null
< row
: bytes
| gen_t
: bytes
| a
: uint64 2
: null
: bytes
| YES
: bytes
| varchar(32)
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
| gen_t
: bytes
| b
: uint64 3
: bytes
| 'LOWER(`a`)'
: bytes
| YES
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| ''
: bytes
| LOWER(`a`)
: bytes
| VIRTUAL GENERATED
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
| omni_gap_gen
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
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
| omni_gap_gen
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
| omni_gap_gen
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
| omni_gap_gen
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_gap_gen
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| omni_gap_gen
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "TEXT"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT"
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
| omni_gap_gen
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
| gen_t
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
> query
|
| 		SELECT MVIEW_NAME, CONTAINER_NAME, QUERY
| 		FROM oceanbase.DBA_MVIEWS
| 		WHERE OWNER = ?
| 		ORDER BY MVIEW_NAME
: string
| omni_gap_gen
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
| omni_gap_gen
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| USE `omni_gap_gen_restored`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `gen_t` (
|   `id` int NOT NULL,
|   `a` varchar(32) DEFAULT NULL,
|   `b` varchar(32) GENERATED ALWAYS AS (LOWER(`a`)) VIRTUAL DEFAULT 'LOWER(`a`)',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< error
| 1064 42000 …near 'DEFAULT 'LOWER(`a`)',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8m' at line 4

== TestDumpGaps/materialized_view default#1
> exec
| CREATE TABLE base_t (id INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE MATERIALIZED VIEW mv_t AS SELECT id FROM base_t;
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
| omni_gap_mv
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
| base_t
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
| omni_gap_mv
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
| base_t
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
: bytes
|
: bytes
|
: null
< row
: bytes
| mv_t
: bytes
| id
: uint64 1
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
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_gap_mv
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
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
| omni_gap_mv
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
| omni_gap_mv
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
| omni_gap_mv
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_gap_mv
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| omni_gap_mv
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "TEXT"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT"
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
| omni_gap_mv
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
| base_t
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
| ''
: bytes
| utf8mb4
> query
|
| 		SELECT MVIEW_NAME, CONTAINER_NAME, QUERY
| 		FROM oceanbase.DBA_MVIEWS
| 		WHERE OWNER = ?
| 		ORDER BY MVIEW_NAME
: string
| omni_gap_mv
< result
< column "MVIEW_NAME" "VARCHAR" "string"
< column "CONTAINER_NAME" "VARCHAR" "string"
< column "QUERY" "TEXT" "string"
< row
: bytes
| mv_t
: bytes
| __mv_container_0000000000000001
: bytes
| select `omni_gap_mv`.`base_t`.`id` AS `id` from `omni_gap_mv`.`base_t`
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_gap_mv
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| USE `omni_gap_mv_restored`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `base_t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE `__mv_container_0000000000000001` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
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
| omni_gap_mv_restored
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
| base_t
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
| omni_gap_mv_restored
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
| base_t
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
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_gap_mv_restored
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
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
| omni_gap_mv_restored
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
| omni_gap_mv_restored
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
| omni_gap_mv_restored
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_gap_mv_restored
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| omni_gap_mv_restored
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "TEXT"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT"
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
| omni_gap_mv_restored
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
| base_t
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
| ''
: bytes
| utf8mb4
> query
|
| 		SELECT MVIEW_NAME, CONTAINER_NAME, QUERY
| 		FROM oceanbase.DBA_MVIEWS
| 		WHERE OWNER = ?
| 		ORDER BY MVIEW_NAME
: string
| omni_gap_mv_restored
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
| omni_gap_mv_restored
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestDumpGaps/trigger_ACTION_STATEMENT_repeats_header default#1
> exec
| CREATE TABLE trg_t (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TRIGGER trg_1 BEFORE INSERT ON trg_t FOR EACH ROW SET NEW.v = 1;
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
| omni_gap_trg
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
| trg_t
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
| omni_gap_trg
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
| trg_t
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
: bytes
|
: bytes
|
: null
< row
: bytes
| trg_t
: bytes
| v
: uint64 2
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
| omni_gap_trg
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
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
| omni_gap_trg
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "TEXT"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| trg_1
: bytes
| trg_t
: bytes
| INSERT
: bytes
| BEFORE
: bytes
| TRIGGER trg_1 BEFORE INSERT ON trg_t FOR EACH ROW SET NEW.v = 1
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
| omni_gap_trg
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
| omni_gap_trg
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| omni_gap_trg
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| omni_gap_trg
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "VARCHAR"
< column "SQL_MODE" "TEXT"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT"
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
| omni_gap_trg
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
| trg_t
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
> query
|
| 		SELECT MVIEW_NAME, CONTAINER_NAME, QUERY
| 		FROM oceanbase.DBA_MVIEWS
| 		WHERE OWNER = ?
| 		ORDER BY MVIEW_NAME
: string
| omni_gap_trg
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
| omni_gap_trg
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| USE `omni_gap_trg_restored`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `trg_t` (
|   `id` int NOT NULL,
|   `v` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
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
| CREATE TRIGGER `trg_1`
| BEFORE INSERT ON `trg_t` FOR EACH ROW
| TRIGGER trg_1 BEFORE INSERT ON trg_t FOR EACH ROW SET NEW.v = 1
< error
| 1064 42000 …near 'TRIGGER trg_1 BEFORE INSERT ON trg_t FOR EACH ROW SET NEW.v = 1' at line 3
