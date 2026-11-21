-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDumpRestoresColumnGroupsAndVectorIndexes default#1
> exec
| CREATE TABLE t_cg (a INT, b INT) WITH COLUMN GROUP(all columns, each column);
< ok 0 0
> exec
| CREATE TABLE t_each (a INT) WITH COLUMN GROUP(each column);
< ok 0 0
> exec
| CREATE TABLE t_vec (id INT PRIMARY KEY, v VECTOR(3), VECTOR INDEX vi(v) WITH (distance=l2, type=hnsw));
< ok 0 0
> exec
| CREATE TABLE p_quoted (`a,b` INT, c INT) PARTITION BY KEY (`a,b`, c) PARTITIONS 2;
< ok 0 0
> exec
| CREATE TABLE t_fn (code VARCHAR(10), KEY k1 ((CONCAT(code, 'a\'b'))), KEY k2 ((CONCAT(code, 'c\\d'))));
< ok 0 0
> exec
| CREATE TABLE t_ft (id INT PRIMARY KEY, body VARCHAR(200), FULLTEXT KEY fk (body) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3));
< ok 0 0
> exec
| CREATE TABLE t_ft_default (id INT PRIMARY KEY, body VARCHAR(200), FULLTEXT KEY fk2 (body));
< ok 0 0
> exec
| CREATE TABLE t_global (a INT, b INT, c INT, KEY k_plain (c), UNIQUE KEY k_global (b) GLOBAL) PARTITION BY HASH(a) PARTITIONS 2;
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
| omni_ob_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR" "string"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "EXPRESSION" "VARBINARY" "[]uint8"
< column "SEQ_IN_INDEX" "UNSIGNED BIGINT" "uint64"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "BIGINT" "int64"
< column "INDEX_COMMENT" "TEXT" "string"
< row
: bytes
| t_fn
: bytes
| k1
: bytes
| SYS_NC17$
: bytes
| A
: int64 -1
: bytes
| CONCAT(`code`,'a\'b')
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_fn
: bytes
| k2
: bytes
| SYS_NC18$
: bytes
| A
: int64 -1
: bytes
| CONCAT(`code`,'c\\d')
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_ft
: bytes
| fk
: bytes
| __word_segment_19_1795257251884077
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_ft
: bytes
| fk
: bytes
| __doc_id_1795257251883674
: bytes
| A
: int64 -1
: null
: uint64 2
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_ft
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t_ft_default
: bytes
| fk2
: bytes
| __word_segment_19_1795257252118926
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_ft_default
: bytes
| fk2
: bytes
| __doc_id_1795257252118560
: bytes
| A
: int64 -1
: null
: uint64 2
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_ft_default
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t_global
: bytes
| k_global
: bytes
| b
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t_global
: bytes
| k_plain
: bytes
| c
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| rowkey_vid_table
: bytes
| id
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi
: bytes
| __vid_1795257251225000
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi
: bytes
| __type_17_1795257251225474
: bytes
| A
: int64 -1
: null
: uint64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vid_rowkey_table
: bytes
| __vid_1795257251225000
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi_index_id_table
: bytes
| __scn_17_1795257251229216
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi_index_id_table
: bytes
| __vid_1795257251225000
: bytes
| A
: int64 -1
: null
: uint64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi_index_id_table
: bytes
| __type_17_1795257251225474
: bytes
| A
: int64 -1
: null
: uint64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi_index_snapshot_data_table
: bytes
| __key_17_1795257251230542
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi_index_snapshot_data_table
: bytes
| id
: bytes
| A
: int64 -1
: null
: uint64 2
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
| omni_ob_show_create_roundtrip
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
| p_quoted
: bytes
| a,b
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
| p_quoted
: bytes
| c
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
< row
: bytes
| t_cg
: bytes
| a
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
| t_cg
: bytes
| b
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
< row
: bytes
| t_each
: bytes
| a
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
| t_fn
: bytes
| code
: uint64 1
: null
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
|
: bytes
|
: null
< row
: bytes
| t_ft
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
| t_ft
: bytes
| body
: uint64 2
: null
: bytes
| YES
: bytes
| varchar(200)
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
| t_ft_default
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
| t_ft_default
: bytes
| body
: uint64 2
: null
: bytes
| YES
: bytes
| varchar(200)
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
| t_global
: bytes
| a
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
| t_global
: bytes
| b
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
< row
: bytes
| t_global
: bytes
| c
: uint64 3
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
| t_vec
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
| t_vec
: bytes
| v
: uint64 2
: null
: bytes
| YES
: bytes
| VECTOR(3)
: bytes
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
| omni_ob_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< column "ENFORCED" "VARCHAR"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION,
| 			IFNULL(CHECK_OPTION, '')
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_ob_show_create_roundtrip
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< column "IFNULL(CHECK_OPTION, '')" "VARCHAR"
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
| omni_ob_show_create_roundtrip
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
| 		SELECT
| 			EVENT_NAME,
| 			EVENT_TYPE,
| 			IFNULL(CAST(EXECUTE_AT AS CHAR), ''),
| 			IFNULL(INTERVAL_VALUE, ''),
| 			IFNULL(INTERVAL_FIELD, ''),
| 			IFNULL(CAST(STARTS AS CHAR), ''),
| 			IFNULL(CAST(ENDS AS CHAR), ''),
| 			IFNULL(STATUS, ''),
| 			IFNULL(ON_COMPLETION, ''),
| 			IFNULL(EVENT_COMMENT, ''),
| 			IFNULL(TIME_ZONE, ''),
| 			IFNULL(SQL_MODE, ''),
| 			IFNULL(EVENT_DEFINITION, '')
| 		FROM information_schema.EVENTS
| 		WHERE EVENT_SCHEMA = ?
| 		ORDER BY EVENT_NAME
: string
| omni_ob_show_create_roundtrip
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_TYPE" "VARCHAR" "string"
< column "IFNULL(CAST(EXECUTE_AT AS CHAR), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(CAST(STARTS AS CHAR), '')" "VARCHAR"
< column "IFNULL(CAST(ENDS AS CHAR), '')" "VARCHAR"
< column "IFNULL(STATUS, '')" "VARCHAR"
< column "IFNULL(ON_COMPLETION, '')" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "IFNULL(TIME_ZONE, '')" "VARCHAR"
< column "IFNULL(SQL_MODE, '')" "TEXT"
< column "IFNULL(EVENT_DEFINITION, '')" "TEXT"
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
| omni_ob_show_create_roundtrip
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
| omni_ob_show_create_roundtrip
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
| omni_ob_show_create_roundtrip
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
| omni_ob_show_create_roundtrip
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
| p_quoted
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
| t_cg
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
| t_each
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
| t_fn
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
| t_ft
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
| t_ft_default
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
| t_global
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
| t_vec
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
| omni_ob_show_create_roundtrip
< result
< column "MVIEW_NAME" "VARCHAR" "string"
< column "CONTAINER_NAME" "VARCHAR" "string"
< column "QUERY" "TEXT" "string"
> query
|
| 		SELECT MASTER, LOG_TABLE
| 		FROM oceanbase.DBA_MVIEW_LOGS
| 		WHERE LOG_OWNER = ?
| 		ORDER BY MASTER
: string
| omni_ob_show_create_roundtrip
< result
< column "MASTER" "VARCHAR" "string"
< column "LOG_TABLE" "VARCHAR" "string"
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip`.`p_quoted`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| p_quoted
: bytes
| CREATE TABLE `p_quoted` (
|   `a,b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key(`a,b`, `c`)
| (partition `p0`,
| partition `p1`)
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip`.`t_cg`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_cg
: bytes
| CREATE TABLE `t_cg` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns, each column)
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip`.`t_each`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_each
: bytes
| CREATE TABLE `t_each` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(each column)
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip`.`t_fn`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_fn
: bytes
| CREATE TABLE `t_fn` (
|   `code` varchar(10) DEFAULT NULL,
|   KEY `k1` ((CONCAT(`code`,'a\'b'))) BLOCK_SIZE 16384 LOCAL,
|   KEY `k2` ((CONCAT(`code`,'c\\d'))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip`.`t_ft`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_ft
: bytes
| CREATE TABLE `t_ft` (
|   `id` int(11) NOT NULL,
|   `body` varchar(200) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`body`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip`.`t_ft_default`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_ft_default
: bytes
| CREATE TABLE `t_ft_default` (
|   `id` int(11) NOT NULL,
|   `body` varchar(200) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk2` (`body`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip`.`t_global`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_global
: bytes
| CREATE TABLE `t_global` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   UNIQUE KEY `k_global` (`b`) BLOCK_SIZE 16384 GLOBAL
|  ,
|   KEY `k_plain` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by hash(`a`)
| (partition `p0`,
| partition `p1`)
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip`.`t_vec`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_vec
: bytes
| CREATE TABLE `t_vec` (
|   `id` int(11) NOT NULL,
|   `v` VECTOR(3) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   VECTOR KEY `vi` (`v`) WITH (DISTANCE=L2, TYPE=HNSW, LIB=VSAG, M=16, EF_CONSTRUCTION=200, EF_SEARCH=64) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_ob_show_create_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
|
| 		SELECT m.MVIEW_NAME, c.COLUMN_NAME, c.COLUMN_TYPE
| 		FROM oceanbase.DBA_MVIEWS m
| 		JOIN information_schema.COLUMNS c
| 		ON c.TABLE_SCHEMA = m.OWNER AND c.TABLE_NAME = m.CONTAINER_NAME
| 		WHERE m.OWNER = ?
| 		ORDER BY m.MVIEW_NAME, c.ORDINAL_POSITION
: string
| omni_ob_show_create_roundtrip
< result
< column "MVIEW_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT"
> exec
| USE `omni_ob_show_create_roundtrip_restored`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE TABLE `p_quoted` (
|   `a,b` int DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE `t_cg` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci WITH COLUMN GROUP(all columns, each column)
< ok 0 0
> exec
| CREATE TABLE `t_each` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci WITH COLUMN GROUP(each column)
< ok 0 0
> exec
| CREATE TABLE `t_fn` (
|   `code` varchar(10) DEFAULT NULL,
|   KEY `k1` ((CONCAT(`code`,'a\'b'))),
|   KEY `k2` ((CONCAT(`code`,'c\\d')))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE `t_ft` (
|   `id` int NOT NULL,
|   `body` varchar(200) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`body`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE `t_ft_default` (
|   `id` int NOT NULL,
|   `body` varchar(200) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk2` (`body`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE `t_global` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   UNIQUE KEY `k_global` (`b`) GLOBAL,
|   KEY `k_plain` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE `t_vec` (
|   `id` int NOT NULL,
|   `v` VECTOR(3) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   VECTOR KEY `vi` (`v`) WITH (DISTANCE=L2, TYPE=HNSW, LIB=VSAG, M=16, EF_CONSTRUCTION=200, EF_SEARCH=64)
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
| omni_ob_show_create_roundtrip_restored
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR" "string"
< column "IFNULL(SUB_PART, -1)" "BIGINT"
< column "EXPRESSION" "VARBINARY" "[]uint8"
< column "SEQ_IN_INDEX" "UNSIGNED BIGINT" "uint64"
< column "INDEX_TYPE" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "BIGINT" "int64"
< column "INDEX_COMMENT" "TEXT" "string"
< row
: bytes
| t_fn
: bytes
| k1
: bytes
| SYS_NC17$
: bytes
| A
: int64 -1
: bytes
| CONCAT(`code`,'a\'b')
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_fn
: bytes
| k2
: bytes
| SYS_NC18$
: bytes
| A
: int64 -1
: bytes
| CONCAT(`code`,'c\\d')
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_ft
: bytes
| fk
: bytes
| __word_segment_19_1795257255112943
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_ft
: bytes
| fk
: bytes
| __doc_id_1795257255112686
: bytes
| A
: int64 -1
: null
: uint64 2
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_ft
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t_ft_default
: bytes
| fk2
: bytes
| __word_segment_19_1795257255303298
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_ft_default
: bytes
| fk2
: bytes
| __doc_id_1795257255303034
: bytes
| A
: int64 -1
: null
: uint64 2
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_ft_default
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t_global
: bytes
| k_global
: bytes
| b
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t_global
: bytes
| k_plain
: bytes
| c
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| rowkey_vid_table
: bytes
| id
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi
: bytes
| __vid_1795257255686485
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi
: bytes
| __type_17_1795257255686813
: bytes
| A
: int64 -1
: null
: uint64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vid_rowkey_table
: bytes
| __vid_1795257255686485
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi_index_id_table
: bytes
| __scn_17_1795257255688734
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi_index_id_table
: bytes
| __vid_1795257255686485
: bytes
| A
: int64 -1
: null
: uint64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi_index_id_table
: bytes
| __type_17_1795257255686813
: bytes
| A
: int64 -1
: null
: uint64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi_index_snapshot_data_table
: bytes
| __key_17_1795257255689986
: bytes
| A
: int64 -1
: null
: uint64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t_vec
: bytes
| vi_index_snapshot_data_table
: bytes
| id
: bytes
| A
: int64 -1
: null
: uint64 2
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
| omni_ob_show_create_roundtrip_restored
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
| p_quoted
: bytes
| a,b
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
| p_quoted
: bytes
| c
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
< row
: bytes
| t_cg
: bytes
| a
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
| t_cg
: bytes
| b
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
< row
: bytes
| t_each
: bytes
| a
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
| t_fn
: bytes
| code
: uint64 1
: null
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
|
: bytes
|
: null
< row
: bytes
| t_ft
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
| t_ft
: bytes
| body
: uint64 2
: null
: bytes
| YES
: bytes
| varchar(200)
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
| t_ft_default
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
| t_ft_default
: bytes
| body
: uint64 2
: null
: bytes
| YES
: bytes
| varchar(200)
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
| t_global
: bytes
| a
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
| t_global
: bytes
| b
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
< row
: bytes
| t_global
: bytes
| c
: uint64 3
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
| t_vec
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
| t_vec
: bytes
| v
: uint64 2
: null
: bytes
| YES
: bytes
| VECTOR(3)
: bytes
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
| omni_ob_show_create_roundtrip_restored
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< column "ENFORCED" "VARCHAR"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION,
| 			IFNULL(CHECK_OPTION, '')
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| omni_ob_show_create_roundtrip_restored
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< column "IFNULL(CHECK_OPTION, '')" "VARCHAR"
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
| omni_ob_show_create_roundtrip_restored
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
| 		SELECT
| 			EVENT_NAME,
| 			EVENT_TYPE,
| 			IFNULL(CAST(EXECUTE_AT AS CHAR), ''),
| 			IFNULL(INTERVAL_VALUE, ''),
| 			IFNULL(INTERVAL_FIELD, ''),
| 			IFNULL(CAST(STARTS AS CHAR), ''),
| 			IFNULL(CAST(ENDS AS CHAR), ''),
| 			IFNULL(STATUS, ''),
| 			IFNULL(ON_COMPLETION, ''),
| 			IFNULL(EVENT_COMMENT, ''),
| 			IFNULL(TIME_ZONE, ''),
| 			IFNULL(SQL_MODE, ''),
| 			IFNULL(EVENT_DEFINITION, '')
| 		FROM information_schema.EVENTS
| 		WHERE EVENT_SCHEMA = ?
| 		ORDER BY EVENT_NAME
: string
| omni_ob_show_create_roundtrip_restored
< result
< column "EVENT_NAME" "VARCHAR"
< column "EVENT_TYPE" "VARCHAR" "string"
< column "IFNULL(CAST(EXECUTE_AT AS CHAR), '')" "VARCHAR"
< column "IFNULL(INTERVAL_VALUE, '')" "VARCHAR"
< column "IFNULL(INTERVAL_FIELD, '')" "VARCHAR"
< column "IFNULL(CAST(STARTS AS CHAR), '')" "VARCHAR"
< column "IFNULL(CAST(ENDS AS CHAR), '')" "VARCHAR"
< column "IFNULL(STATUS, '')" "VARCHAR"
< column "IFNULL(ON_COMPLETION, '')" "VARCHAR"
< column "IFNULL(EVENT_COMMENT, '')" "TEXT"
< column "IFNULL(TIME_ZONE, '')" "VARCHAR"
< column "IFNULL(SQL_MODE, '')" "TEXT"
< column "IFNULL(EVENT_DEFINITION, '')" "TEXT"
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
| omni_ob_show_create_roundtrip_restored
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
| omni_ob_show_create_roundtrip_restored
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
| omni_ob_show_create_roundtrip_restored
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
| omni_ob_show_create_roundtrip_restored
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
| p_quoted
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
| t_cg
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
| t_each
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
| t_fn
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
| t_ft
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
| t_ft_default
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
| t_global
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
| t_vec
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
| omni_ob_show_create_roundtrip_restored
< result
< column "MVIEW_NAME" "VARCHAR" "string"
< column "CONTAINER_NAME" "VARCHAR" "string"
< column "QUERY" "TEXT" "string"
> query
|
| 		SELECT MASTER, LOG_TABLE
| 		FROM oceanbase.DBA_MVIEW_LOGS
| 		WHERE LOG_OWNER = ?
| 		ORDER BY MASTER
: string
| omni_ob_show_create_roundtrip_restored
< result
< column "MASTER" "VARCHAR" "string"
< column "LOG_TABLE" "VARCHAR" "string"
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip_restored`.`p_quoted`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| p_quoted
: bytes
| CREATE TABLE `p_quoted` (
|   `a,b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip_restored`.`t_cg`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_cg
: bytes
| CREATE TABLE `t_cg` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns, each column)
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip_restored`.`t_each`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_each
: bytes
| CREATE TABLE `t_each` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(each column)
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip_restored`.`t_fn`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_fn
: bytes
| CREATE TABLE `t_fn` (
|   `code` varchar(10) DEFAULT NULL,
|   KEY `k1` ((CONCAT(`code`,'a\'b'))) BLOCK_SIZE 16384 LOCAL,
|   KEY `k2` ((CONCAT(`code`,'c\\d'))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip_restored`.`t_ft`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_ft
: bytes
| CREATE TABLE `t_ft` (
|   `id` int(11) NOT NULL,
|   `body` varchar(200) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`body`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip_restored`.`t_ft_default`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_ft_default
: bytes
| CREATE TABLE `t_ft_default` (
|   `id` int(11) NOT NULL,
|   `body` varchar(200) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk2` (`body`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip_restored`.`t_global`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_global
: bytes
| CREATE TABLE `t_global` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   UNIQUE KEY `k_global` (`b`) BLOCK_SIZE 16384 GLOBAL,
|   KEY `k_plain` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `omni_ob_show_create_roundtrip_restored`.`t_vec`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t_vec
: bytes
| CREATE TABLE `t_vec` (
|   `id` int(11) NOT NULL,
|   `v` VECTOR(3) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   VECTOR KEY `vi` (`v`) WITH (DISTANCE=L2, TYPE=HNSW, LIB=VSAG, M=16, EF_CONSTRUCTION=200, EF_SEARCH=64) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_ob_show_create_roundtrip_restored
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
|
| 		SELECT m.MVIEW_NAME, c.COLUMN_NAME, c.COLUMN_TYPE
| 		FROM oceanbase.DBA_MVIEWS m
| 		JOIN information_schema.COLUMNS c
| 		ON c.TABLE_SCHEMA = m.OWNER AND c.TABLE_NAME = m.CONTAINER_NAME
| 		WHERE m.OWNER = ?
| 		ORDER BY m.MVIEW_NAME, c.ORDINAL_POSITION
: string
| omni_ob_show_create_roundtrip_restored
< result
< column "MVIEW_NAME" "VARCHAR" "string"
< column "COLUMN_NAME" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT"
