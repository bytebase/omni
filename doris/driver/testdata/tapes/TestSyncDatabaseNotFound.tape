-- pin: default dyrnq/doris:3.1.4 sha256:88f99819c0877567b77d78cc8367f0e2dd4b6c1602fe17c6f4b1e332d8e315dd 3.1.4-rc02

== TestSyncDatabaseNotFound default#1
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			IFNULL(COLUMN_NAME, ''),
| 			ORDINAL_POSITION,
| 			IFNULL(COLUMN_KEY, ''),
| 			COLUMN_DEFAULT,
| 			IS_NULLABLE,
| 			COLUMN_TYPE,
| 			IFNULL(CHARACTER_SET_NAME, ''),
| 			IFNULL(COLLATION_NAME, ''),
| 			COLUMN_COMMENT,
| 			EXTRA
| 		FROM information_schema.columns
| 		WHERE TABLE_SCHEMA = 'omni_no_such_database'
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "CHAR"
< column "IFNULL(COLUMN_NAME, '')" "CHAR"
< column "ORDINAL_POSITION" "BIGINT"
< column "IFNULL(COLUMN_KEY, '')" "CHAR"
< column "COLUMN_DEFAULT" "CHAR"
< column "IS_NULLABLE" "CHAR"
< column "COLUMN_TYPE" "CHAR"
< column "IFNULL(CHARACTER_SET_NAME, '')" "CHAR"
< column "IFNULL(COLLATION_NAME, '')" "CHAR"
< column "COLUMN_COMMENT" "CHAR"
< column "EXTRA" "CHAR"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = 'omni_no_such_database'
< result
< column "TABLE_NAME" "CHAR"
< column "VIEW_DEFINITION" "CHAR"
> query
| SELECT Name, RefreshInfo, QuerySql FROM mv_infos("database"="omni_no_such_database") ORDER BY Name
< error
| 1105 HY000 errCode = 2, detailMessage = omni_trace_1[INTERNAL_ERROR]TStatus: errCode = 2, detailMessage = Unknown database 'omni_no_such_database'
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			TABLE_TYPE,
| 			IFNULL(TABLE_ROWS, 0),
| 			IFNULL(DATA_LENGTH, 0),
| 			IFNULL(INDEX_LENGTH, 0),
| 			IFNULL(DATA_FREE, 0),
| 			IFNULL(CREATE_OPTIONS, ''),
| 			IFNULL(TABLE_COMMENT, '')
| 		FROM information_schema.TABLES
| 		WHERE TABLE_SCHEMA = 'omni_no_such_database'
| 		ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "CHAR"
< column "IFNULL(TABLE_ROWS, 0)" "BIGINT"
< column "IFNULL(DATA_LENGTH, 0)" "BIGINT"
< column "IFNULL(INDEX_LENGTH, 0)" "BIGINT"
< column "IFNULL(DATA_FREE, 0)" "BIGINT"
< column "IFNULL(CREATE_OPTIONS, '')" "CHAR"
< column "IFNULL(TABLE_COMMENT, '')" "CHAR"
> query
|
| 		SELECT ROUTINE_NAME
| 		FROM information_schema.ROUTINES
| 		WHERE ROUTINE_SCHEMA = 'omni_no_such_database' AND ROUTINE_TYPE = 'PROCEDURE'
| 		ORDER BY ROUTINE_NAME
< result
< column "ROUTINE_NAME" "CHAR"
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = 'omni_no_such_database'
< result
< column "DEFAULT_CHARACTER_SET_NAME" "CHAR"
< column "DEFAULT_COLLATION_NAME" "CHAR"
