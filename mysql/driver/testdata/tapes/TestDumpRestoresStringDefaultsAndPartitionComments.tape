-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46
-- pin: mysql84 mysql:8.4.13 sha256:a81d043930e51745522bd5c430961cf8286f0a23f0a9a2a6e0108f7ea71b0902 mysql-8.4.13

== TestDumpRestoresStringDefaultsAndPartitionComments/default default#1
> exec
| CREATE TABLE literal_defaults (
| 	id INT PRIMARY KEY,
| 	a VARCHAR(40) NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
| 	b VARCHAR(40) NOT NULL DEFAULT 'current_date_label',
| 	ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
| ) COMMENT='PARTITIONS 9 SUBPARTITIONS 7';
< ok 0 0
> exec
| CREATE TABLE hash_commented (
| 	id INT NOT NULL PRIMARY KEY COMMENT 'PARTITIONS 9'
| ) COMMENT='PARTITIONS 11' PARTITION BY HASH (id) PARTITIONS 2;
< ok 0 0
> exec
| CREATE TABLE range_commented (
| 	id INT NOT NULL PRIMARY KEY
| ) COMMENT='PARTITIONS 9 SUBPARTITIONS 7'
| PARTITION BY RANGE (id) SUBPARTITION BY HASH (id) SUBPARTITIONS 2 (
| 	PARTITION p0 VALUES LESS THAN (100),
| 	PARTITION p1 VALUES LESS THAN (MAXVALUE)
| );
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT"
< column "INDEX_COMMENT" "VARCHAR"
< row
: bytes
| hash_commented
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
| literal_defaults
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
| range_commented
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED INT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
< row
: bytes
| hash_commented
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
| 'PARTITIONS 9'
: bytes
|
: bytes
|
: null
< row
: bytes
| literal_defaults
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
| literal_defaults
: bytes
| a
: int64 2
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| varchar(40)
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
| literal_defaults
: bytes
| b
: int64 3
: bytes
| 'current_date_label'
: bytes
| NO
: bytes
| varchar(40)
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
| literal_defaults
: bytes
| ts
: int64 4
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
| range_commented
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
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
| omni_sync_defect_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "SET"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "TIME_ZONE" "VARCHAR"
< column "SQL_MODE" "SET"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "ENUM"
< column "UPDATE_RULE" "ENUM"
< column "MATCH_OPTION" "ENUM"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| hash_commented
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
| hash_commented
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
| range_commented
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| 100
< row
: bytes
| range_commented
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| 100
< row
: bytes
| range_commented
: bytes
| p1
: bytes
| p1sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| MAXVALUE
< row
: bytes
| range_commented
: bytes
| p1
: bytes
| p1sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| MAXVALUE
> query
| SHOW CREATE TABLE `omni_sync_defect_roundtrip`.`hash_commented`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| hash_commented
: bytes
| CREATE TABLE `hash_commented` (
|   `id` int NOT NULL COMMENT 'PARTITIONS 9',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 11'
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 2 */
> query
| SHOW CREATE TABLE `omni_sync_defect_roundtrip`.`range_commented`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| range_commented
: bytes
| CREATE TABLE `range_commented` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 9 SUBPARTITIONS 7'
| /*!50100 PARTITION BY RANGE (`id`)
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (100) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
| omni_sync_defect_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "ENUM"
< column "SQL_MODE" "SET"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
< row
: bytes
| hash_commented
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 32768
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| 'PARTITIONS 11'
: bytes
| utf8mb4
< row
: bytes
| literal_defaults
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
| 0
: bytes
| 0
: bytes
|
: bytes
| 'PARTITIONS 9 SUBPARTITIONS 7'
: bytes
| utf8mb4
< row
: bytes
| range_commented
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
| 'PARTITIONS 9 SUBPARTITIONS 7'
: bytes
| utf8mb4
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_sync_defect_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE `omni_sync_defect_roundtrip`
< ok 3 0
> exec
| CREATE DATABASE `omni_sync_defect_roundtrip`
< ok 1 0
> exec
| USE `omni_sync_defect_roundtrip`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `hash_commented` (
|   `id` int NOT NULL COMMENT 'PARTITIONS 9',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 11'
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 2 */
< ok 0 0
> exec
| CREATE TABLE `literal_defaults` (
|   `id` int NOT NULL,
|   `a` varchar(40) NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
|   `b` varchar(40) NOT NULL DEFAULT 'current_date_label',
|   `ts` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 9 SUBPARTITIONS 7'
< ok 0 0
> exec
| CREATE TABLE `range_commented` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 9 SUBPARTITIONS 7'
| /*!50100 PARTITION BY RANGE (`id`)
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION `p0` VALUES LESS THAN (100) ENGINE=InnoDB,
|  PARTITION `p1` VALUES LESS THAN MAXVALUE ENGINE=InnoDB) */
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT"
< column "INDEX_COMMENT" "VARCHAR"
< row
: bytes
| hash_commented
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
| literal_defaults
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
| range_commented
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED INT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
< row
: bytes
| hash_commented
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
| 'PARTITIONS 9'
: bytes
|
: bytes
|
: null
< row
: bytes
| literal_defaults
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
| literal_defaults
: bytes
| a
: int64 2
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| varchar(40)
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
| literal_defaults
: bytes
| b
: int64 3
: bytes
| 'current_date_label'
: bytes
| NO
: bytes
| varchar(40)
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
| literal_defaults
: bytes
| ts
: int64 4
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
| range_commented
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
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
| omni_sync_defect_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "SET"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "TIME_ZONE" "VARCHAR"
< column "SQL_MODE" "SET"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "ENUM"
< column "UPDATE_RULE" "ENUM"
< column "MATCH_OPTION" "ENUM"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| hash_commented
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
| hash_commented
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
| range_commented
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| 100
< row
: bytes
| range_commented
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| 100
< row
: bytes
| range_commented
: bytes
| p1
: bytes
| p1sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| MAXVALUE
< row
: bytes
| range_commented
: bytes
| p1
: bytes
| p1sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| MAXVALUE
> query
| SHOW CREATE TABLE `omni_sync_defect_roundtrip`.`hash_commented`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| hash_commented
: bytes
| CREATE TABLE `hash_commented` (
|   `id` int NOT NULL COMMENT 'PARTITIONS 9',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 11'
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 2 */
> query
| SHOW CREATE TABLE `omni_sync_defect_roundtrip`.`range_commented`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| range_commented
: bytes
| CREATE TABLE `range_commented` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 9 SUBPARTITIONS 7'
| /*!50100 PARTITION BY RANGE (`id`)
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (100) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
| omni_sync_defect_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "ENUM"
< column "SQL_MODE" "SET"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
< row
: bytes
| hash_commented
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 32768
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| 'PARTITIONS 11'
: bytes
| utf8mb4
< row
: bytes
| literal_defaults
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
| 0
: bytes
| 0
: bytes
|
: bytes
| 'PARTITIONS 9 SUBPARTITIONS 7'
: bytes
| utf8mb4
< row
: bytes
| range_commented
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
| 'PARTITIONS 9 SUBPARTITIONS 7'
: bytes
| utf8mb4
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_sync_defect_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestDumpRestoresStringDefaultsAndPartitionComments/mysql84 mysql84#1
> exec
| CREATE TABLE literal_defaults (
| 	id INT PRIMARY KEY,
| 	a VARCHAR(40) NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
| 	b VARCHAR(40) NOT NULL DEFAULT 'current_date_label',
| 	ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
| ) COMMENT='PARTITIONS 9 SUBPARTITIONS 7';
< ok 0 0
> exec
| CREATE TABLE hash_commented (
| 	id INT NOT NULL PRIMARY KEY COMMENT 'PARTITIONS 9'
| ) COMMENT='PARTITIONS 11' PARTITION BY HASH (id) PARTITIONS 2;
< ok 0 0
> exec
| CREATE TABLE range_commented (
| 	id INT NOT NULL PRIMARY KEY
| ) COMMENT='PARTITIONS 9 SUBPARTITIONS 7'
| PARTITION BY RANGE (id) SUBPARTITION BY HASH (id) SUBPARTITIONS 2 (
| 	PARTITION p0 VALUES LESS THAN (100),
| 	PARTITION p1 VALUES LESS THAN (MAXVALUE)
| );
< ok 0 0
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT"
< column "INDEX_COMMENT" "VARCHAR"
< row
: bytes
| hash_commented
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
| literal_defaults
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
| range_commented
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED INT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
< row
: bytes
| hash_commented
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
| 'PARTITIONS 9'
: bytes
|
: bytes
|
: null
< row
: bytes
| literal_defaults
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
| literal_defaults
: bytes
| a
: int64 2
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| varchar(40)
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
| literal_defaults
: bytes
| b
: int64 3
: bytes
| 'current_date_label'
: bytes
| NO
: bytes
| varchar(40)
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
| literal_defaults
: bytes
| ts
: int64 4
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
| range_commented
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
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
| omni_sync_defect_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "SET"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "TIME_ZONE" "VARCHAR"
< column "SQL_MODE" "SET"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "ENUM"
< column "UPDATE_RULE" "ENUM"
< column "MATCH_OPTION" "ENUM"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| hash_commented
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
| hash_commented
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
| range_commented
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| 100
< row
: bytes
| range_commented
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| 100
< row
: bytes
| range_commented
: bytes
| p1
: bytes
| p1sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| MAXVALUE
< row
: bytes
| range_commented
: bytes
| p1
: bytes
| p1sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| MAXVALUE
> query
| SHOW CREATE TABLE `omni_sync_defect_roundtrip`.`hash_commented`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| hash_commented
: bytes
| CREATE TABLE `hash_commented` (
|   `id` int NOT NULL COMMENT 'PARTITIONS 9',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 11'
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 2 */
> query
| SHOW CREATE TABLE `omni_sync_defect_roundtrip`.`range_commented`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| range_commented
: bytes
| CREATE TABLE `range_commented` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 9 SUBPARTITIONS 7'
| /*!50100 PARTITION BY RANGE (`id`)
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (100) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
| omni_sync_defect_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "ENUM"
< column "SQL_MODE" "SET"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
< row
: bytes
| hash_commented
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 32768
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| 'PARTITIONS 11'
: bytes
| utf8mb4
< row
: bytes
| literal_defaults
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
| 0
: bytes
| 0
: bytes
|
: bytes
| 'PARTITIONS 9 SUBPARTITIONS 7'
: bytes
| utf8mb4
< row
: bytes
| range_commented
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
| 'PARTITIONS 9 SUBPARTITIONS 7'
: bytes
| utf8mb4
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_sync_defect_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE `omni_sync_defect_roundtrip`
< ok 3 0
> exec
| CREATE DATABASE `omni_sync_defect_roundtrip`
< ok 1 0
> exec
| USE `omni_sync_defect_roundtrip`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `hash_commented` (
|   `id` int NOT NULL COMMENT 'PARTITIONS 9',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 11'
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 2 */
< ok 0 0
> exec
| CREATE TABLE `literal_defaults` (
|   `id` int NOT NULL,
|   `a` varchar(40) NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
|   `b` varchar(40) NOT NULL DEFAULT 'current_date_label',
|   `ts` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 9 SUBPARTITIONS 7'
< ok 0 0
> exec
| CREATE TABLE `range_commented` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 9 SUBPARTITIONS 7'
| /*!50100 PARTITION BY RANGE (`id`)
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION `p0` VALUES LESS THAN (100) ENGINE=InnoDB,
|  PARTITION `p1` VALUES LESS THAN MAXVALUE ENGINE=InnoDB) */
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT"
< column "INDEX_COMMENT" "VARCHAR"
< row
: bytes
| hash_commented
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
| literal_defaults
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
| range_commented
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED INT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
< row
: bytes
| hash_commented
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
| 'PARTITIONS 9'
: bytes
|
: bytes
|
: null
< row
: bytes
| literal_defaults
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
| literal_defaults
: bytes
| a
: int64 2
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| varchar(40)
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
| literal_defaults
: bytes
| b
: int64 3
: bytes
| 'current_date_label'
: bytes
| NO
: bytes
| varchar(40)
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
| literal_defaults
: bytes
| ts
: int64 4
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
| range_commented
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "TEXT"
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
| omni_sync_defect_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "EVENT_MANIPULATION" "ENUM"
< column "ACTION_TIMING" "ENUM"
< column "ACTION_STATEMENT" "TEXT"
< column "SQL_MODE" "SET"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "TIME_ZONE" "VARCHAR"
< column "SQL_MODE" "SET"
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "ENUM"
< column "UPDATE_RULE" "ENUM"
< column "MATCH_OPTION" "ENUM"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| hash_commented
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
| hash_commented
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
| range_commented
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| 100
< row
: bytes
| range_commented
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| 100
< row
: bytes
| range_commented
: bytes
| p1
: bytes
| p1sp0
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| MAXVALUE
< row
: bytes
| range_commented
: bytes
| p1
: bytes
| p1sp1
: bytes
| RANGE
: bytes
| HASH
: bytes
| `id`
: bytes
| `id`
: bytes
| MAXVALUE
> query
| SHOW CREATE TABLE `omni_sync_defect_roundtrip`.`hash_commented`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| hash_commented
: bytes
| CREATE TABLE `hash_commented` (
|   `id` int NOT NULL COMMENT 'PARTITIONS 9',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 11'
| /*!50100 PARTITION BY HASH (`id`)
| PARTITIONS 2 */
> query
| SHOW CREATE TABLE `omni_sync_defect_roundtrip`.`range_commented`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| range_commented
: bytes
| CREATE TABLE `range_commented` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PARTITIONS 9 SUBPARTITIONS 7'
| /*!50100 PARTITION BY RANGE (`id`)
| SUBPARTITION BY HASH (`id`)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (100) ENGINE = InnoDB,
|  PARTITION p1 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
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
| omni_sync_defect_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR"
< column "ROUTINE_TYPE" "ENUM"
< column "SQL_MODE" "SET"
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
| omni_sync_defect_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL"
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL"
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR"
< row
: bytes
| hash_commented
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 32768
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| 'PARTITIONS 11'
: bytes
| utf8mb4
< row
: bytes
| literal_defaults
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
| 0
: bytes
| 0
: bytes
|
: bytes
| 'PARTITIONS 9 SUBPARTITIONS 7'
: bytes
| utf8mb4
< row
: bytes
| range_commented
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
| 'PARTITIONS 9 SUBPARTITIONS 7'
: bytes
| utf8mb4
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_sync_defect_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
