-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestOracle_DiffMatchesServer default#1
> exec
| SET GLOBAL tidb_enable_check_constraint = ON
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int(11), b tinyint(4))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` tinyint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, b tinyint)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` tinyint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
< row
: bytes
| t
: bytes
| b
: int64 2
: null
: bytes
| YES
: bytes
| tinyint
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
< row
: bytes
| t
: bytes
| b
: int64 2
: null
: bytes
| YES
: bytes
| tinyint
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a boolean)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` tinyint(1) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a tinyint(1))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` tinyint(1) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
: bytes
| a
: int64 1
: null
: bytes
| YES
: bytes
| tinyint(1)
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
: bytes
| a
: int64 1
: null
: bytes
| YES
: bytes
| tinyint(1)
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a varchar(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
: bytes
| a
: int64 1
: null
: bytes
| YES
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
: bytes
| a
: int64 1
: null
: bytes
| YES
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int DEFAULT 5, b decimal(5,2) DEFAULT 1.5)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT '5',
|   `b` decimal(5,2) DEFAULT '1.5'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int DEFAULT '5', b decimal(5,2) DEFAULT '1.5')
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT '5',
|   `b` decimal(5,2) DEFAULT '1.5'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
: bytes
| a
: int64 1
: bytes
| '5'
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
< row
: bytes
| t
: bytes
| b
: int64 2
: bytes
| '1.5'
: bytes
| YES
: bytes
| decimal(5,2)
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
: bytes
| a
: int64 1
: bytes
| '5'
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
< row
: bytes
| t
: bytes
| b
: int64 2
: bytes
| '1.5'
: bytes
| YES
: bytes
| decimal(5,2)
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int DEFAULT NULL, b varchar(3) NULL)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` varchar(3) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, b varchar(3))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` varchar(3) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
< row
: bytes
| t
: bytes
| b
: int64 2
: null
: bytes
| YES
: bytes
| varchar(3)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
< row
: bytes
| t
: bytes
| b
: int64 2
: null
: bytes
| YES
: bytes
| varchar(3)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, KEY (a))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, KEY a (a))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| a
: bytes
| a
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| a
: bytes
| a
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, b int, KEY ka (a), KEY kb (b))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `ka` (`a`),
|   KEY `kb` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, b int, KEY kb (b), KEY ka (a))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `kb` (`b`),
|   KEY `ka` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| ka
: bytes
| a
: int64 1
: bytes
|
: null
: null
< row
: bytes
| t
: bytes
| kb
: bytes
| b
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| ka
: bytes
| a
: int64 1
: bytes
|
: null
: null
< row
: bytes
| t
: bytes
| kb
: bytes
| b
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (id int NOT NULL, PRIMARY KEY (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int) ENGINE=MyISAM ROW_FORMAT=COMPACT
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, KEY k (a) USING HASH)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, KEY k (a))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| a
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| a
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a text, FULLTEXT KEY ft (a))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` text DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a text)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` text DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
: bytes
| a
: int64 1
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
: bytes
| a
: int64 1
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int) PARTITION BY LINEAR HASH (a) PARTITIONS 4
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY HASH (`a`) PARTITIONS 4
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int) PARTITION BY HASH (a) PARTITIONS 4
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY HASH (`a`) PARTITIONS 4
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| t
: bytes
| p0
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p1
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p2
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p3
: bytes
| HASH
: bytes
| `a`
: bytes
|
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
| partitioned
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| t
: bytes
| p0
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p1
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p2
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p3
: bytes
| HASH
: bytes
| `a`
: bytes
|
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
| partitioned
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (id int AUTO_INCREMENT PRIMARY KEY) AUTO_INCREMENT=100
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=100
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (id int AUTO_INCREMENT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
| auto_increment
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
| auto_increment
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, CONSTRAINT c CHECK (a>0))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   CONSTRAINT `c` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, CONSTRAINT c CHECK (`a` > 0))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   CONSTRAINT `c` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| t
: bytes
| c
: bytes
| (`a` > 0)
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| t
: bytes
| c
: bytes
| (`a` > 0)
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| omni_diff_a
: string
| v
< result
< column "CHECK_OPTION" "VARCHAR" "string"
< row
: bytes
| CASCADED
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`a`) AS SELECT `a` AS `a` FROM `omni_diff_a`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a FROM omni_diff_b.t
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| omni_diff_b
: string
| v
< result
< column "CHECK_OPTION" "VARCHAR" "string"
< row
: bytes
| CASCADED
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`a`) AS SELECT `t`.`a` AS `a` FROM `omni_diff_b`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
< row
: bytes
| v
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v
: bytes
| SELECT `a` AS `a` FROM `omni_diff_a`.`t`
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
< row
: bytes
| v
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 'VIEW'
: bytes
|
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
< row
: bytes
| v
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v
: bytes
| SELECT `t`.`a` AS `a` FROM `omni_diff_b`.`t`
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
< row
: bytes
| v
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 'VIEW'
: bytes
|
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE SEQUENCE s START 1 MINVALUE 1 INCREMENT 1 CACHE 1000 NOCYCLE
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| s
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
< row
: bytes
| s
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: int64 0
: int64 1000
: bytes
|
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| s
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
< row
: bytes
| s
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: int64 0
: int64 1000
: bytes
|
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE p (id int PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id int PRIMARY KEY, pid int, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| c
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `pid` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `fk` (`pid`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE p (id int PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id int PRIMARY KEY, pid int, KEY fk (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| c
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `pid` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `fk` (`pid`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| c
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
< row
: bytes
| c
: bytes
| fk
: bytes
| pid
: int64 1
: bytes
|
: null
: null
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| c
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
< row
: bytes
| c
: bytes
| pid
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
< row
: bytes
| p
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR" "string"
< column "UPDATE_RULE" "VARCHAR" "string"
< column "MATCH_OPTION" "VARCHAR" "string"
< row
: bytes
| c
: bytes
| fk
: bytes
| pid
: bytes
| omni_diff_a
: bytes
| p
: bytes
| id
: bytes
| NO ACTION
: bytes
| NO ACTION
: bytes
| NONE
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| c
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
|
< row
: bytes
| p
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| c
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
< row
: bytes
| c
: bytes
| fk
: bytes
| pid
: int64 1
: bytes
|
: null
: null
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| c
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
< row
: bytes
| c
: bytes
| pid
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
< row
: bytes
| p
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR" "string"
< column "UPDATE_RULE" "VARCHAR" "string"
< column "MATCH_OPTION" "VARCHAR" "string"
< row
: bytes
| c
: bytes
| fk
: bytes
| pid
: bytes
| omni_diff_b
: bytes
| p
: bytes
| id
: bytes
| NO ACTION
: bytes
| NO ACTION
: bytes
| NONE
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| c
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
|
< row
: bytes
| p
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a decimal(5,2) DEFAULT 1.5)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` decimal(5,2) DEFAULT '1.5'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a decimal(5,2) DEFAULT 1.50)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` decimal(5,2) DEFAULT '1.50'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
: bytes
| a
: int64 1
: bytes
| '1.5'
: bytes
| YES
: bytes
| decimal(5,2)
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
: bytes
| a
: int64 1
: bytes
| '1.50'
: bytes
| YES
: bytes
| decimal(5,2)
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a bigint)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` bigint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
: bytes
| a
: int64 1
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int NOT NULL)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
: bytes
| a
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int DEFAULT 1)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT '1'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int DEFAULT 2)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT '2'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
: bytes
| a
: int64 1
: bytes
| '1'
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
: bytes
| a
: int64 1
: bytes
| '2'
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int COMMENT 'x')
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL COMMENT 'x'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int COMMENT 'it''s')
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL COMMENT 'it''s'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
| 'x'
: bytes
|
: bytes
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
| 'it\'s'
: bytes
|
: bytes
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a varchar(5))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(5) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a varchar(5) COLLATE utf8mb4_general_ci)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(5) COLLATE utf8mb4_general_ci DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
: bytes
| a
: int64 1
: null
: bytes
| YES
: bytes
| varchar(5)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
: bytes
| a
: int64 1
: null
: bytes
| YES
: bytes
| varchar(5)
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a timestamp DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` timestamp DEFAULT CURRENT_TIMESTAMP
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
: bytes
| a
: int64 1
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| YES
: bytes
| timestamp
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
: bytes
| a
: int64 1
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| YES
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, b int AS (a + 1))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int GENERATED ALWAYS AS (`a` + 1) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, b int AS (a + 2))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int GENERATED ALWAYS AS (`a` + 2) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
< row
: bytes
| t
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
| `a` + 1
: bytes
| VIRTUAL GENERATED
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
< row
: bytes
| t
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
| `a` + 2
: bytes
| VIRTUAL GENERATED
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, b int)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, b int, KEY k (a))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, b int, KEY k (a, b))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `k` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| a
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| a
: int64 1
: bytes
|
: null
: null
< row
: bytes
| t
: bytes
| k
: bytes
| b
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, KEY k (a))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, UNIQUE KEY k (a))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   UNIQUE KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| a
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| a
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, KEY k (a))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, KEY k (a) INVISIBLE)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   KEY `k` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| a
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| a
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, KEY k (a))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, KEY k (a) COMMENT 'c')
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   KEY `k` (`a`) COMMENT 'c'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| a
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| a
: int64 1
: bytes
| c
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a varchar(20), KEY k (a))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(20) DEFAULT NULL,
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a varchar(20), KEY k (a(5)))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(20) DEFAULT NULL,
|   KEY `k` (`a`(5))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| a
: int64 1
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
: bytes
| a
: int64 1
: null
: bytes
| YES
: bytes
| varchar(20)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| a
: int64 1
: bytes
|
: int64 5
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
: bytes
| a
: int64 1
: null
: bytes
| YES
: bytes
| varchar(20)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, KEY k ((a + 1)))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   KEY `k` ((`a` + 1))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, KEY k ((a + 2)))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   KEY `k` ((`a` + 2))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| NULL
: int64 1
: bytes
|
: null
: bytes
| `a` + 1
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| k
: bytes
| NULL
: int64 1
: bytes
|
: null
: bytes
| `a` + 2
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (id int, PRIMARY KEY (id) CLUSTERED)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (id int, PRIMARY KEY (id) NONCLUSTERED)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (id bigint AUTO_RANDOM(5) PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` bigint NOT NULL /*T![auto_rand] AUTO_RANDOM(5) */,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (id bigint AUTO_RANDOM(6) PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` bigint NOT NULL /*T![auto_rand] AUTO_RANDOM(6) */,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| PK_AUTO_RANDOM_BITS=5
: bytes
| CLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
|
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| PK_AUTO_RANDOM_BITS=6
: bytes
| CLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int) SHARD_ROW_ID_BITS=4
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T! SHARD_ROW_ID_BITS=4 */
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| SHARD_BITS=4
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int) AUTO_ID_CACHE=100
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![auto_id_cache] AUTO_ID_CACHE=100 */
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a datetime) TTL = a + INTERVAL 1 DAY
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![ttl] TTL=`a` + INTERVAL 1 DAY */ /*T![ttl] TTL_ENABLE='ON' */ /*T![ttl] TTL_JOB_INTERVAL='1h' */
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a datetime) TTL = a + INTERVAL 2 DAY
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![ttl] TTL=`a` + INTERVAL 2 DAY */ /*T![ttl] TTL_ENABLE='ON' */ /*T![ttl] TTL_JOB_INTERVAL='1h' */
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
: bytes
| a
: int64 1
: null
: bytes
| YES
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
: bytes
| a
: int64 1
: null
: bytes
| YES
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int) COMMENT 'x'
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='x'
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int) COMMENT 'y'
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='y'
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| 'x'
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| 'y'
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int) COLLATE utf8mb4_general_ci
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_general_ci
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int) PARTITION BY HASH (a) PARTITIONS 4
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY HASH (`a`) PARTITIONS 4
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int) PARTITION BY HASH (a) PARTITIONS 8
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY HASH (`a`) PARTITIONS 8
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| t
: bytes
| p0
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p1
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p2
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p3
: bytes
| HASH
: bytes
| `a`
: bytes
|
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
| partitioned
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| t
: bytes
| p0
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p1
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p2
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p3
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p4
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p5
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p6
: bytes
| HASH
: bytes
| `a`
: bytes
|
< row
: bytes
| t
: bytes
| p7
: bytes
| HASH
: bytes
| `a`
: bytes
|
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
| partitioned
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY RANGE (`a`)
| (PARTITION `p0` VALUES LESS THAN (10))
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (20))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY RANGE (`a`)
| (PARTITION `p0` VALUES LESS THAN (20))
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| t
: bytes
| p0
: bytes
| RANGE
: bytes
| `a`
: bytes
| 10
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
| partitioned
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| t
: bytes
| p0
: bytes
| RANGE
: bytes
| `a`
: bytes
| 20
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
| partitioned
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE p (id int PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (pid int, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| c
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `pid` int DEFAULT NULL,
|   KEY `fk` (`pid`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE p (id int PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (pid int, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id) ON DELETE CASCADE)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| c
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `pid` int DEFAULT NULL,
|   KEY `fk` (`pid`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| c
: bytes
| fk
: bytes
| pid
: int64 1
: bytes
|
: null
: null
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| c
: bytes
| pid
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
< row
: bytes
| p
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR" "string"
< column "UPDATE_RULE" "VARCHAR" "string"
< column "MATCH_OPTION" "VARCHAR" "string"
< row
: bytes
| c
: bytes
| fk
: bytes
| pid
: bytes
| omni_diff_a
: bytes
| p
: bytes
| id
: bytes
| NO ACTION
: bytes
| NO ACTION
: bytes
| NONE
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| c
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
< row
: bytes
| p
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| c
: bytes
| fk
: bytes
| pid
: int64 1
: bytes
|
: null
: null
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| id
: int64 0
: bytes
|
: null
: null
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| c
: bytes
| pid
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
< row
: bytes
| p
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR" "string"
< column "UPDATE_RULE" "VARCHAR" "string"
< column "MATCH_OPTION" "VARCHAR" "string"
< row
: bytes
| c
: bytes
| fk
: bytes
| pid
: bytes
| omni_diff_b
: bytes
| p
: bytes
| id
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| c
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
< row
: bytes
| p
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, CONSTRAINT c CHECK (a > 0))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   CONSTRAINT `c` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, CONSTRAINT c CHECK (a > 1))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   CONSTRAINT `c` CHECK ((`a` > 1))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| t
: bytes
| c
: bytes
| (`a` > 0)
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| t
: bytes
| c
: bytes
| (`a` > 1)
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, CONSTRAINT c CHECK (a > 0))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   CONSTRAINT `c` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, CONSTRAINT c CHECK (a > 0) NOT ENFORCED)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   CONSTRAINT `c` CHECK ((`a` > 0)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| t
: bytes
| c
: bytes
| (`a` > 0)
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| t
: bytes
| c
: bytes
| (`a` > 0)
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int, b int)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| omni_diff_a
: string
| v
< result
< column "CHECK_OPTION" "VARCHAR" "string"
< row
: bytes
| CASCADED
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`a`) AS SELECT `a` AS `a` FROM `omni_diff_a`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int, b int)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT b AS a FROM t
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| omni_diff_b
: string
| v
< result
< column "CHECK_OPTION" "VARCHAR" "string"
< row
: bytes
| CASCADED
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`a`) AS SELECT `b` AS `a` FROM `omni_diff_b`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
< row
: bytes
| t
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
< row
: bytes
| v
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v
: bytes
| SELECT `a` AS `a` FROM `omni_diff_a`.`t`
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
< row
: bytes
| v
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 'VIEW'
: bytes
|
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
< row
: bytes
| t
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
< row
: bytes
| v
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v
: bytes
| SELECT `b` AS `a` FROM `omni_diff_b`.`t`
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
< row
: bytes
| v
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 'VIEW'
: bytes
|
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WITH CASCADED CHECK OPTION
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| omni_diff_a
: string
| v
< result
< column "CHECK_OPTION" "VARCHAR" "string"
< row
: bytes
| CASCADED
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`a`) AS SELECT `a` AS `a` FROM `omni_diff_a`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WITH LOCAL CHECK OPTION
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| omni_diff_b
: string
| v
< result
< column "CHECK_OPTION" "VARCHAR" "string"
< row
: bytes
| LOCAL
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`a`) AS SELECT `a` AS `a` FROM `omni_diff_b`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
< row
: bytes
| v
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v
: bytes
| SELECT `a` AS `a` FROM `omni_diff_a`.`t`
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
< row
: bytes
| v
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 'VIEW'
: bytes
|
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
< row
: bytes
| v
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v
: bytes
| SELECT `a` AS `a` FROM `omni_diff_b`.`t`
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
< row
: bytes
| v
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 'VIEW'
: bytes
|
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| omni_diff_b
: string
| v
< result
< column "CHECK_OPTION" "VARCHAR" "string"
< row
: bytes
| CASCADED
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`a`) AS SELECT `a` AS `a` FROM `omni_diff_b`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
< row
: bytes
| v
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v
: bytes
| SELECT `a` AS `a` FROM `omni_diff_b`.`t`
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
< row
: bytes
| v
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 'VIEW'
: bytes
|
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE SEQUENCE s INCREMENT 2
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 2 cache 1000 nocycle ENGINE=InnoDB
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| s
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
< row
: bytes
| s
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: int64 0
: int64 1000
: bytes
|
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| s
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
< row
: bytes
| s
: int64 1
: int64 1
: int64 9223372036854775806
: int64 2
: int64 0
: int64 1000
: bytes
|
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE SEQUENCE s NOCACHE
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| s
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
< row
: bytes
| s
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: int64 0
: int64 1000
: bytes
|
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| s
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
< row
: bytes
| s
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: int64 0
: int64 0
: bytes
|
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE SEQUENCE s COMMENT 'c'
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB COMMENT='c'
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| s
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
< row
: bytes
| s
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: int64 0
: int64 1000
: bytes
|
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| s
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
< row
: bytes
| s
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: int64 0
: int64 1000
: bytes
| c
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| s
: bytes
| SEQUENCE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
< row
: bytes
| s
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: int64 0
: int64 1000
: bytes
|
> exec
| DROP DATABASE IF EXISTS omni_diff_a
< ok 0 0
> exec
| CREATE DATABASE omni_diff_a
< ok 0 0
> exec
| USE omni_diff_a
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_diff_b
< ok 0 0
> exec
| CREATE DATABASE omni_diff_b
< ok 0 0
> exec
| USE omni_diff_b
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> exec
| CREATE TABLE u (a int)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| u
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `u`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_a
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
< row
: bytes
| t
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_a
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_a
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_a
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_a
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
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_diff_b
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
< row
: bytes
| t
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
< row
: bytes
| u
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
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_diff_b
< result
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_diff_b
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLE_COLLATION, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLE_COMMENT, ''))" "VARCHAR"
< column "IFNULL(TIDB_ROW_ID_SHARDING_INFO, '')" "VARCHAR" "string"
< column "IFNULL(TIDB_PK_TYPE, '')" "VARCHAR" "string"
< column "IFNULL(CREATE_OPTIONS, '')" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
< row
: bytes
| u
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_diff_b
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_diff_b
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT" "int64"
< column "CYCLE" "TINYINT" "int8"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR" "string"
> exec
| USE test
< ok 0 0
> exec
| SET GLOBAL tidb_enable_check_constraint = OFF
< ok 0 0
