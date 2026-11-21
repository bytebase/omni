-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46
-- pin: mysql84 mysql:8.4.13 sha256:a81d043930e51745522bd5c430961cf8286f0a23f0a9a2a6e0108f7ea71b0902 mysql-8.4.13

== TestDumpRestoresTablespacesKeyBlockSizesAndKeyAlgorithms/default default#1
> exec
| CREATE TABLESPACE ts_sync ADD DATAFILE 'ts_sync.ibd';
< ok 0 0
> exec
| CREATE TABLE t_ts (id INT PRIMARY KEY) TABLESPACE ts_sync;
< ok 0 0
> exec
| CREATE TABLE `t-ts` (id INT PRIMARY KEY) TABLESPACE ts_sync;
< ok 0 0
> exec
| CREATE TABLE t_sys (id INT PRIMARY KEY) TABLESPACE innodb_system;
< ok 0 0
> exec
| CREATE TABLE t_plain (id INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t_kbs (a INT, b INT, KEY ka (a) KEY_BLOCK_SIZE=4, KEY kb (b)) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8;
< ok 0 0
> exec
| CREATE TABLE t_myisam (a INT, KEY ka (a) KEY_BLOCK_SIZE=1024) ENGINE=MyISAM;
< ok 0 0
> exec
| CREATE TABLE p_key (a INT) PARTITION BY KEY ALGORITHM=1 (a) PARTITIONS 2;
< ok 0 0
> exec
| CREATE TABLE p_sub (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY KEY ALGORITHM=1 (b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10));
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
| omni_sync_show_create_roundtrip
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
| t-ts
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
| t_kbs
: bytes
| ka
: bytes
| a
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
| t_kbs
: bytes
| kb
: bytes
| b
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
| t_myisam
: bytes
| ka
: bytes
| a
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
| t_plain
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
| t_sys
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
| t_ts
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
| omni_sync_show_create_roundtrip
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
| p_key
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| p_sub
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| p_sub
: bytes
| b
: int64 2
: null
: bytes
| YES
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
| t-ts
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
| t_kbs
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| t_kbs
: bytes
| b
: int64 2
: null
: bytes
| YES
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
| t_myisam
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| t_plain
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
| t_sys
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
| t_ts
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
| 				cc.CHECK_CLAUSE,
| 				tc.ENFORCED
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| omni_sync_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< column "ENFORCED" "VARCHAR" "string"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
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
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| omni_sync_show_create_roundtrip
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
| p_key
: bytes
| p0
: null
: bytes
| KEY
: null
: bytes
| `a`
: null
: null
: bytes
|
< row
: bytes
| p_key
: bytes
| p1
: null
: bytes
| KEY
: null
: bytes
| `a`
: null
: null
: bytes
|
< row
: bytes
| p_sub
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| KEY
: bytes
| `a`
: bytes
| `b`
: bytes
| 10
: bytes
|
< row
: bytes
| p_sub
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| KEY
: bytes
| `a`
: bytes
| `b`
: bytes
| 10
: bytes
|
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`p_key`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| p_key
: bytes
| CREATE TABLE `p_key` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY */ /*!50611 ALGORITHM = 1 */ /*!50100 (a)
| PARTITIONS 2 */
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`p_sub`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| p_sub
: bytes
| CREATE TABLE `p_sub` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`a`)
| SUBPARTITION BY KEY */ /*!50611 ALGORITHM = 1 */ /*!50100 (b)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB) */
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
| omni_sync_show_create_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
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
| omni_sync_show_create_roundtrip
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
| p_key
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
| ''
: bytes
| utf8mb4
< row
: bytes
| p_sub
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
| ''
: bytes
| utf8mb4
< row
: bytes
| t-ts
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
| ''
: bytes
| utf8mb4
< row
: bytes
| t_kbs
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 8192
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=COMPRESSED KEY_BLOCK_SIZE=8
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t_myisam
: bytes
| BASE TABLE
: bytes
| MyISAM
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 0
: bytes
| 1024
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
| t_plain
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
| ''
: bytes
| utf8mb4
< row
: bytes
| t_sys
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
| 6291456
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t_ts
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
| ''
: bytes
| utf8mb4
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| omni_sync_show_create_roundtrip
< result
< column "NAME" "VARCHAR" "string"
< row
: bytes
| omni_sync_show_create_roundtrip/t_ts
< row
: bytes
| omni_sync_show_create_roundtrip/t@002dts
< row
: bytes
| omni_sync_show_create_roundtrip/t_sys
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t-ts`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t-ts
: bytes
| CREATE TABLE `t-ts` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `ts_sync` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_kbs`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_kbs
: bytes
| CREATE TABLE `t_kbs` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `ka` (`a`) KEY_BLOCK_SIZE=4,
|   KEY `kb` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_myisam`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_myisam
: bytes
| CREATE TABLE `t_myisam` (
|   `a` int DEFAULT NULL,
|   KEY `ka` (`a`) KEY_BLOCK_SIZE=1024
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_sys`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_sys
: bytes
| CREATE TABLE `t_sys` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `innodb_system` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_ts`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_ts
: bytes
| CREATE TABLE `t_ts` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `ts_sync` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_sync_show_create_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE `omni_sync_show_create_roundtrip`
< ok 8 0
> exec
| CREATE DATABASE `omni_sync_show_create_roundtrip`
< ok 1 0
> exec
| USE `omni_sync_show_create_roundtrip`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| /*!80030 SET @OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@@SQL_GENERATE_INVISIBLE_PRIMARY_KEY, SQL_GENERATE_INVISIBLE_PRIMARY_KEY=OFF */
< ok 0 0
> exec
| CREATE TABLE `p_key` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY ALGORITHM = 1 (`a`)
| PARTITIONS 2 */
< ok 0 0
> exec
| /*!80030 SET SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY */
< ok 0 0
> exec
| /*!80030 SET @OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@@SQL_GENERATE_INVISIBLE_PRIMARY_KEY, SQL_GENERATE_INVISIBLE_PRIMARY_KEY=OFF */
< ok 0 0
> exec
| CREATE TABLE `p_sub` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`a`)
| SUBPARTITION BY KEY ALGORITHM = 1 (`b`)
| SUBPARTITIONS 2
| (PARTITION `p0` VALUES LESS THAN (10) ENGINE=InnoDB) */
< ok 0 0
> exec
| /*!80030 SET SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY */
< ok 0 0
> exec
| CREATE TABLE `t-ts` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci /*!50100 TABLESPACE `ts_sync` */
< ok 0 0
> exec
| /*!80030 SET @OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@@SQL_GENERATE_INVISIBLE_PRIMARY_KEY, SQL_GENERATE_INVISIBLE_PRIMARY_KEY=OFF */
< ok 0 0
> exec
| CREATE TABLE `t_kbs` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `ka` (`a`) KEY_BLOCK_SIZE=4,
|   KEY `kb` (`b`)
| ) ENGINE=InnoDB row_format=COMPRESSED KEY_BLOCK_SIZE=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| /*!80030 SET SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY */
< ok 0 0
> exec
| /*!80030 SET @OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@@SQL_GENERATE_INVISIBLE_PRIMARY_KEY, SQL_GENERATE_INVISIBLE_PRIMARY_KEY=OFF */
< ok 0 0
> exec
| CREATE TABLE `t_myisam` (
|   `a` int DEFAULT NULL,
|   KEY `ka` (`a`) KEY_BLOCK_SIZE=1024
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| /*!80030 SET SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY */
< ok 0 0
> exec
| CREATE TABLE `t_plain` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE TABLE `t_sys` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci /*!50100 TABLESPACE `innodb_system` */
< ok 0 0
> exec
| CREATE TABLE `t_ts` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci /*!50100 TABLESPACE `ts_sync` */
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
| omni_sync_show_create_roundtrip
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
| t-ts
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
| t_kbs
: bytes
| ka
: bytes
| a
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
| t_kbs
: bytes
| kb
: bytes
| b
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
| t_myisam
: bytes
| ka
: bytes
| a
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
| t_plain
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
| t_sys
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
| t_ts
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
| omni_sync_show_create_roundtrip
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
| p_key
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| p_sub
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| p_sub
: bytes
| b
: int64 2
: null
: bytes
| YES
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
| t-ts
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
| t_kbs
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| t_kbs
: bytes
| b
: int64 2
: null
: bytes
| YES
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
| t_myisam
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| t_plain
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
| t_sys
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
| t_ts
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
| 				cc.CHECK_CLAUSE,
| 				tc.ENFORCED
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| omni_sync_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< column "ENFORCED" "VARCHAR" "string"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
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
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| omni_sync_show_create_roundtrip
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
| p_key
: bytes
| p0
: null
: bytes
| KEY
: null
: bytes
| `a`
: null
: null
: bytes
|
< row
: bytes
| p_key
: bytes
| p1
: null
: bytes
| KEY
: null
: bytes
| `a`
: null
: null
: bytes
|
< row
: bytes
| p_sub
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| KEY
: bytes
| `a`
: bytes
| `b`
: bytes
| 10
: bytes
|
< row
: bytes
| p_sub
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| KEY
: bytes
| `a`
: bytes
| `b`
: bytes
| 10
: bytes
|
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`p_key`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| p_key
: bytes
| CREATE TABLE `p_key` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY */ /*!50611 ALGORITHM = 1 */ /*!50100 (a)
| PARTITIONS 2 */
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`p_sub`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| p_sub
: bytes
| CREATE TABLE `p_sub` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`a`)
| SUBPARTITION BY KEY */ /*!50611 ALGORITHM = 1 */ /*!50100 (b)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB) */
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
| omni_sync_show_create_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
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
| omni_sync_show_create_roundtrip
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
| p_key
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
| ''
: bytes
| utf8mb4
< row
: bytes
| p_sub
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
| ''
: bytes
| utf8mb4
< row
: bytes
| t-ts
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
| ''
: bytes
| utf8mb4
< row
: bytes
| t_kbs
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 8192
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=COMPRESSED KEY_BLOCK_SIZE=8
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t_myisam
: bytes
| BASE TABLE
: bytes
| MyISAM
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 0
: bytes
| 1024
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
| t_plain
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
| ''
: bytes
| utf8mb4
< row
: bytes
| t_sys
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
| 6291456
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t_ts
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
| ''
: bytes
| utf8mb4
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| omni_sync_show_create_roundtrip
< result
< column "NAME" "VARCHAR" "string"
< row
: bytes
| omni_sync_show_create_roundtrip/t@002dts
< row
: bytes
| omni_sync_show_create_roundtrip/t_sys
< row
: bytes
| omni_sync_show_create_roundtrip/t_ts
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t-ts`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t-ts
: bytes
| CREATE TABLE `t-ts` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `ts_sync` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_kbs`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_kbs
: bytes
| CREATE TABLE `t_kbs` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `ka` (`a`) KEY_BLOCK_SIZE=4,
|   KEY `kb` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_myisam`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_myisam
: bytes
| CREATE TABLE `t_myisam` (
|   `a` int DEFAULT NULL,
|   KEY `ka` (`a`) KEY_BLOCK_SIZE=1024
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_sys`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_sys
: bytes
| CREATE TABLE `t_sys` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `innodb_system` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_ts`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_ts
: bytes
| CREATE TABLE `t_ts` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `ts_sync` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_sync_show_create_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestDumpRestoresTablespacesKeyBlockSizesAndKeyAlgorithms/default default#2
> exec
| DROP TABLESPACE ts_sync
< ok 0 0

== TestDumpRestoresTablespacesKeyBlockSizesAndKeyAlgorithms/mysql84 mysql84#1
> exec
| CREATE TABLESPACE ts_sync ADD DATAFILE 'ts_sync.ibd';
< ok 0 0
> exec
| CREATE TABLE t_ts (id INT PRIMARY KEY) TABLESPACE ts_sync;
< ok 0 0
> exec
| CREATE TABLE `t-ts` (id INT PRIMARY KEY) TABLESPACE ts_sync;
< ok 0 0
> exec
| CREATE TABLE t_sys (id INT PRIMARY KEY) TABLESPACE innodb_system;
< ok 0 0
> exec
| CREATE TABLE t_plain (id INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t_kbs (a INT, b INT, KEY ka (a) KEY_BLOCK_SIZE=4, KEY kb (b)) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8;
< ok 0 0
> exec
| CREATE TABLE t_myisam (a INT, KEY ka (a) KEY_BLOCK_SIZE=1024) ENGINE=MyISAM;
< ok 0 0
> exec
| CREATE TABLE p_key (a INT) PARTITION BY KEY ALGORITHM=1 (a) PARTITIONS 2;
< ok 0 0
> exec
| CREATE TABLE p_sub (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY KEY ALGORITHM=1 (b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10));
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
| omni_sync_show_create_roundtrip
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
| t-ts
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
| t_kbs
: bytes
| ka
: bytes
| a
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
| t_kbs
: bytes
| kb
: bytes
| b
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
| t_myisam
: bytes
| ka
: bytes
| a
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
| t_plain
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
| t_sys
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
| t_ts
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
| omni_sync_show_create_roundtrip
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
| p_key
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| p_sub
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| p_sub
: bytes
| b
: int64 2
: null
: bytes
| YES
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
| t-ts
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
| t_kbs
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| t_kbs
: bytes
| b
: int64 2
: null
: bytes
| YES
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
| t_myisam
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| t_plain
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
| t_sys
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
| t_ts
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
| 				cc.CHECK_CLAUSE,
| 				tc.ENFORCED
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| omni_sync_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< column "ENFORCED" "VARCHAR" "string"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
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
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| omni_sync_show_create_roundtrip
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
| p_key
: bytes
| p0
: null
: bytes
| KEY
: null
: bytes
| `a`
: null
: null
: bytes
|
< row
: bytes
| p_key
: bytes
| p1
: null
: bytes
| KEY
: null
: bytes
| `a`
: null
: null
: bytes
|
< row
: bytes
| p_sub
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| KEY
: bytes
| `a`
: bytes
| `b`
: bytes
| 10
: bytes
|
< row
: bytes
| p_sub
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| KEY
: bytes
| `a`
: bytes
| `b`
: bytes
| 10
: bytes
|
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`p_key`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| p_key
: bytes
| CREATE TABLE `p_key` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY */ /*!50611 ALGORITHM = 1 */ /*!50100 (a)
| PARTITIONS 2 */
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`p_sub`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| p_sub
: bytes
| CREATE TABLE `p_sub` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`a`)
| SUBPARTITION BY KEY */ /*!50611 ALGORITHM = 1 */ /*!50100 (b)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB) */
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
| omni_sync_show_create_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
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
| omni_sync_show_create_roundtrip
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
| p_key
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
| ''
: bytes
| utf8mb4
< row
: bytes
| p_sub
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
| ''
: bytes
| utf8mb4
< row
: bytes
| t-ts
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
| ''
: bytes
| utf8mb4
< row
: bytes
| t_kbs
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 8192
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=COMPRESSED KEY_BLOCK_SIZE=8
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t_myisam
: bytes
| BASE TABLE
: bytes
| MyISAM
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 0
: bytes
| 1024
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
| t_plain
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
| ''
: bytes
| utf8mb4
< row
: bytes
| t_sys
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
| 6291456
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t_ts
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
| ''
: bytes
| utf8mb4
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| omni_sync_show_create_roundtrip
< result
< column "NAME" "VARCHAR" "string"
< row
: bytes
| omni_sync_show_create_roundtrip/t_ts
< row
: bytes
| omni_sync_show_create_roundtrip/t@002dts
< row
: bytes
| omni_sync_show_create_roundtrip/t_sys
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t-ts`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t-ts
: bytes
| CREATE TABLE `t-ts` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `ts_sync` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_kbs`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_kbs
: bytes
| CREATE TABLE `t_kbs` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `ka` (`a`) KEY_BLOCK_SIZE=4,
|   KEY `kb` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_myisam`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_myisam
: bytes
| CREATE TABLE `t_myisam` (
|   `a` int DEFAULT NULL,
|   KEY `ka` (`a`) KEY_BLOCK_SIZE=1024
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_sys`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_sys
: bytes
| CREATE TABLE `t_sys` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `innodb_system` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_ts`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_ts
: bytes
| CREATE TABLE `t_ts` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `ts_sync` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_sync_show_create_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE `omni_sync_show_create_roundtrip`
< ok 8 0
> exec
| CREATE DATABASE `omni_sync_show_create_roundtrip`
< ok 1 0
> exec
| USE `omni_sync_show_create_roundtrip`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| /*!80030 SET @OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@@SQL_GENERATE_INVISIBLE_PRIMARY_KEY, SQL_GENERATE_INVISIBLE_PRIMARY_KEY=OFF */
< ok 0 0
> exec
| CREATE TABLE `p_key` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY ALGORITHM = 1 (`a`)
| PARTITIONS 2 */
< ok 0 0
> exec
| /*!80030 SET SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY */
< ok 0 0
> exec
| /*!80030 SET @OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@@SQL_GENERATE_INVISIBLE_PRIMARY_KEY, SQL_GENERATE_INVISIBLE_PRIMARY_KEY=OFF */
< ok 0 0
> exec
| CREATE TABLE `p_sub` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`a`)
| SUBPARTITION BY KEY ALGORITHM = 1 (`b`)
| SUBPARTITIONS 2
| (PARTITION `p0` VALUES LESS THAN (10) ENGINE=InnoDB) */
< ok 0 0
> exec
| /*!80030 SET SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY */
< ok 0 0
> exec
| CREATE TABLE `t-ts` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci /*!50100 TABLESPACE `ts_sync` */
< ok 0 0
> exec
| /*!80030 SET @OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@@SQL_GENERATE_INVISIBLE_PRIMARY_KEY, SQL_GENERATE_INVISIBLE_PRIMARY_KEY=OFF */
< ok 0 0
> exec
| CREATE TABLE `t_kbs` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `ka` (`a`) KEY_BLOCK_SIZE=4,
|   KEY `kb` (`b`)
| ) ENGINE=InnoDB row_format=COMPRESSED KEY_BLOCK_SIZE=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| /*!80030 SET SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY */
< ok 0 0
> exec
| /*!80030 SET @OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@@SQL_GENERATE_INVISIBLE_PRIMARY_KEY, SQL_GENERATE_INVISIBLE_PRIMARY_KEY=OFF */
< ok 0 0
> exec
| CREATE TABLE `t_myisam` (
|   `a` int DEFAULT NULL,
|   KEY `ka` (`a`) KEY_BLOCK_SIZE=1024
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| /*!80030 SET SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY */
< ok 0 0
> exec
| CREATE TABLE `t_plain` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE TABLE `t_sys` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci /*!50100 TABLESPACE `innodb_system` */
< ok 0 0
> exec
| CREATE TABLE `t_ts` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci /*!50100 TABLESPACE `ts_sync` */
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
| omni_sync_show_create_roundtrip
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
| t-ts
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
| t_kbs
: bytes
| ka
: bytes
| a
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
| t_kbs
: bytes
| kb
: bytes
| b
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
| t_myisam
: bytes
| ka
: bytes
| a
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
| t_plain
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
| t_sys
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
| t_ts
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
| omni_sync_show_create_roundtrip
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
| p_key
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| p_sub
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| p_sub
: bytes
| b
: int64 2
: null
: bytes
| YES
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
| t-ts
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
| t_kbs
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| t_kbs
: bytes
| b
: int64 2
: null
: bytes
| YES
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
| t_myisam
: bytes
| a
: int64 1
: null
: bytes
| YES
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
| t_plain
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
| t_sys
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
| t_ts
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
| 				cc.CHECK_CLAUSE,
| 				tc.ENFORCED
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| omni_sync_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< column "ENFORCED" "VARCHAR" "string"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_sync_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
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
| omni_sync_show_create_roundtrip
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
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
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| omni_sync_show_create_roundtrip
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
| p_key
: bytes
| p0
: null
: bytes
| KEY
: null
: bytes
| `a`
: null
: null
: bytes
|
< row
: bytes
| p_key
: bytes
| p1
: null
: bytes
| KEY
: null
: bytes
| `a`
: null
: null
: bytes
|
< row
: bytes
| p_sub
: bytes
| p0
: bytes
| p0sp0
: bytes
| RANGE
: bytes
| KEY
: bytes
| `a`
: bytes
| `b`
: bytes
| 10
: bytes
|
< row
: bytes
| p_sub
: bytes
| p0
: bytes
| p0sp1
: bytes
| RANGE
: bytes
| KEY
: bytes
| `a`
: bytes
| `b`
: bytes
| 10
: bytes
|
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`p_key`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| p_key
: bytes
| CREATE TABLE `p_key` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY */ /*!50611 ALGORITHM = 1 */ /*!50100 (a)
| PARTITIONS 2 */
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`p_sub`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| p_sub
: bytes
| CREATE TABLE `p_sub` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`a`)
| SUBPARTITION BY KEY */ /*!50611 ALGORITHM = 1 */ /*!50100 (b)
| SUBPARTITIONS 2
| (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB) */
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
| omni_sync_show_create_roundtrip
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
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
| omni_sync_show_create_roundtrip
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
| p_key
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
| ''
: bytes
| utf8mb4
< row
: bytes
| p_sub
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
| ''
: bytes
| utf8mb4
< row
: bytes
| t-ts
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
| ''
: bytes
| utf8mb4
< row
: bytes
| t_kbs
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 8192
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=COMPRESSED KEY_BLOCK_SIZE=8
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t_myisam
: bytes
| BASE TABLE
: bytes
| MyISAM
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 0
: bytes
| 1024
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
| t_plain
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
| ''
: bytes
| utf8mb4
< row
: bytes
| t_sys
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
| 6291456
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t_ts
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
| ''
: bytes
| utf8mb4
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| omni_sync_show_create_roundtrip
< result
< column "NAME" "VARCHAR" "string"
< row
: bytes
| omni_sync_show_create_roundtrip/t@002dts
< row
: bytes
| omni_sync_show_create_roundtrip/t_sys
< row
: bytes
| omni_sync_show_create_roundtrip/t_ts
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t-ts`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t-ts
: bytes
| CREATE TABLE `t-ts` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `ts_sync` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_kbs`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_kbs
: bytes
| CREATE TABLE `t_kbs` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `ka` (`a`) KEY_BLOCK_SIZE=4,
|   KEY `kb` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_myisam`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_myisam
: bytes
| CREATE TABLE `t_myisam` (
|   `a` int DEFAULT NULL,
|   KEY `ka` (`a`) KEY_BLOCK_SIZE=1024
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_sys`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_sys
: bytes
| CREATE TABLE `t_sys` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `innodb_system` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `omni_sync_show_create_roundtrip`.`t_ts`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| t_ts
: bytes
| CREATE TABLE `t_ts` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `ts_sync` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_sync_show_create_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestDumpRestoresTablespacesKeyBlockSizesAndKeyAlgorithms/mysql84 mysql84#2
> exec
| DROP TABLESPACE ts_sync
< ok 0 0
