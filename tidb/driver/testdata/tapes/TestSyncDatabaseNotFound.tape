-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSyncDatabaseNotFound default#1
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR" "string"
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
| omni_no_such_database
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "SEQ_IN_INDEX" "BIGINT"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
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
| omni_no_such_database
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "BIGINT"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< column "GENERATION_EXPRESSION" "TEXT" "string"
< column "EXTRA" "VARCHAR"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION,
| 			IFNULL(SECURITY_TYPE, ''),
| 			IFNULL(CHECK_OPTION, '')
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_no_such_database
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< column "IFNULL(SECURITY_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CHECK_OPTION, '')" "VARCHAR" "string"
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
| omni_no_such_database
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR" "string"
< column "UPDATE_RULE" "VARCHAR" "string"
< column "MATCH_OPTION" "VARCHAR" "string"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			CHECK_CLAUSE
| 		FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ?
| 		ORDER BY TABLE_NAME, CONSTRAINT_NAME
| 	
: string
| omni_no_such_database
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			REPLICA_COUNT,
| 			IFNULL(LOCATION_LABELS, '')
| 		FROM information_schema.TIFLASH_REPLICA
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
| 	
: string
| omni_no_such_database
< result
< column "TABLE_NAME" "VARCHAR"
< column "REPLICA_COUNT" "BIGINT"
< column "IFNULL(LOCATION_LABELS, '')" "VARCHAR" "string"
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
| omni_no_such_database
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
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
| omni_no_such_database
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "SUBPARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< column "PARTITION_COMMENT" "VARCHAR"
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
| 			IFNULL(TIDB_PK_TYPE, ''),
| 			IFNULL(TIDB_PLACEMENT_POLICY_NAME, '')
| 		FROM information_schema.TABLES
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_no_such_database
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_ROWS, 0)" "BIGINT" "int64"
< column "IFNULL(DATA_LENGTH, 0)" "BIGINT" "int64"
< column "IFNULL(INDEX_LENGTH, 0)" "BIGINT" "int64"
< column "IFNULL(DATA_FREE, 0)" "BIGINT" "int64"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PLACEMENT_POLICY_NAME, '')" "VARCHAR" "string"
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME,
| 			IFNULL(TIDB_PLACEMENT_POLICY_NAME, '')
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_no_such_database
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< column "IFNULL(TIDB_PLACEMENT_POLICY_NAME, '')" "VARCHAR" "string"
