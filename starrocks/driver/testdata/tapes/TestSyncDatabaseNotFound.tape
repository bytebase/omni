-- pin: default starrocks/allin1-ubuntu:3.4.10 sha256:8bc803ea7a5bb490c77117a456b32b792ccae698b1f0e47dfcafdac41542566d 3.4.10

== TestSyncDatabaseNotFound default#1
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			IFNULL(COLUMN_NAME, ''),
| 			ORDINAL_POSITION,
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
< column "TABLE_NAME" "VARCHAR"
< column "ifnull(COLUMN_NAME, '')" "VARCHAR"
< column "ORDINAL_POSITION" "BIGINT"
< column "COLUMN_DEFAULT" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_TYPE" "VARCHAR"
< column "ifnull(CHARACTER_SET_NAME, '')" "VARCHAR"
< column "ifnull(COLLATION_NAME, '')" "VARCHAR"
< column "COLUMN_COMMENT" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = 'omni_no_such_database'
< result
< column "TABLE_NAME" "VARCHAR"
< column "VIEW_DEFINITION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(REFRESH_TYPE, ''), IFNULL(MATERIALIZED_VIEW_DEFINITION, '') FROM information_schema.materialized_views WHERE TABLE_SCHEMA = 'omni_no_such_database'
< result
< column "TABLE_NAME" "VARCHAR"
< column "ifnull(REFRESH_TYPE, '')" "VARCHAR"
< column "ifnull(MATERIALIZED_VIEW_DEFINITION, '')" "VARCHAR"
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
| 			IFNULL(TABLE_COMMENT, '')
| 		FROM information_schema.TABLES
| 		WHERE TABLE_SCHEMA = 'omni_no_such_database'
| 		ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "ifnull(ENGINE, '')" "VARCHAR"
< column "ifnull(TABLE_COLLATION, '')" "VARCHAR"
< column "ifnull(TABLE_ROWS, 0)" "BIGINT"
< column "ifnull(DATA_LENGTH, 0)" "BIGINT"
< column "ifnull(INDEX_LENGTH, 0)" "BIGINT"
< column "ifnull(DATA_FREE, 0)" "BIGINT"
< column "ifnull(CREATE_OPTIONS, '')" "VARCHAR"
< column "ifnull(TABLE_COMMENT, '')" "VARCHAR"
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = 'omni_no_such_database'
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
