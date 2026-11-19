-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainerLoadMetadataMatchesTiDB default#1
> exec
| SET GLOBAL tidb_enable_check_constraint = ON
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_mdload
< ok 0 0
> exec
| CREATE DATABASE omni_mdload
< ok 0 0
> exec
| USE omni_mdload
< ok 0 0
> exec
| CREATE TABLE p_range (id int, d date) PARTITION BY RANGE (YEAR(d)) (PARTITION z1 VALUES LESS THAN (2000), PARTITION a2 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| CREATE TABLE p_range_columns (a int, b varchar(10)) PARTITION BY RANGE COLUMNS (a, b) (PARTITION p0 VALUES LESS THAN (10, 'm'), PARTITION p1 VALUES LESS THAN (MAXVALUE, MAXVALUE))
< ok 0 0
> exec
| CREATE TABLE p_list (id int, r int) PARTITION BY LIST (r) (PARTITION pa VALUES IN (1, 2), PARTITION pb VALUES IN (3, NULL), PARTITION pd VALUES IN (DEFAULT))
< ok 0 0
> exec
| CREATE TABLE p_list_columns (id int, c varchar(10)) PARTITION BY LIST COLUMNS (c) (PARTITION pa VALUES IN ('a', 'b'), PARTITION pb VALUES IN ('c'))
< ok 0 0
> exec
| CREATE TABLE p_hash (id int) PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> exec
| CREATE TABLE p_hash_named (id int) PARTITION BY HASH (id) (PARTITION x0, PARTITION x1)
< ok 0 0
> exec
| CREATE TABLE p_key (id int PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 3
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_mdload
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
| p_key
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
| omni_mdload
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
| p_hash
: bytes
| id
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
| p_hash_named
: bytes
| id
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
| p_key
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
| p_list
: bytes
| id
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
| p_list
: bytes
| r
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
| p_list_columns
: bytes
| id
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
| p_list_columns
: bytes
| c
: int64 2
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
< row
: bytes
| p_range
: bytes
| id
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
| p_range
: bytes
| d
: int64 2
: null
: bytes
| YES
: bytes
| date
: bytes
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
| p_range_columns
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
| p_range_columns
: bytes
| b
: int64 2
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
| omni_mdload
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
| omni_mdload
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
| omni_mdload
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_mdload
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "TEXT"
< column "PARTITION_DESCRIPTION" "TEXT"
< row
: bytes
| p_hash
: bytes
| p0
: bytes
| HASH
: bytes
| `id`
: bytes
|
< row
: bytes
| p_hash
: bytes
| p1
: bytes
| HASH
: bytes
| `id`
: bytes
|
< row
: bytes
| p_hash
: bytes
| p2
: bytes
| HASH
: bytes
| `id`
: bytes
|
< row
: bytes
| p_hash
: bytes
| p3
: bytes
| HASH
: bytes
| `id`
: bytes
|
< row
: bytes
| p_hash_named
: bytes
| x0
: bytes
| HASH
: bytes
| `id`
: bytes
|
< row
: bytes
| p_hash_named
: bytes
| x1
: bytes
| HASH
: bytes
| `id`
: bytes
|
< row
: bytes
| p_key
: bytes
| p0
: bytes
| KEY
: bytes
| `id`
: bytes
|
< row
: bytes
| p_key
: bytes
| p1
: bytes
| KEY
: bytes
| `id`
: bytes
|
< row
: bytes
| p_key
: bytes
| p2
: bytes
| KEY
: bytes
| `id`
: bytes
|
< row
: bytes
| p_list
: bytes
| pa
: bytes
| LIST
: bytes
| `r`
: bytes
| 1,2
< row
: bytes
| p_list
: bytes
| pb
: bytes
| LIST
: bytes
| `r`
: bytes
| 3,NULL
< row
: bytes
| p_list
: bytes
| pd
: bytes
| LIST
: bytes
| `r`
: bytes
| DEFAULT
< row
: bytes
| p_list_columns
: bytes
| pa
: bytes
| LIST COLUMNS
: bytes
| `c`
: bytes
| 'a','b'
< row
: bytes
| p_list_columns
: bytes
| pb
: bytes
| LIST COLUMNS
: bytes
| `c`
: bytes
| 'c'
< row
: bytes
| p_range
: bytes
| z1
: bytes
| RANGE
: bytes
| YEAR(`d`)
: bytes
| 2000
< row
: bytes
| p_range
: bytes
| a2
: bytes
| RANGE
: bytes
| YEAR(`d`)
: bytes
| 2010
< row
: bytes
| p_range
: bytes
| pmax
: bytes
| RANGE
: bytes
| YEAR(`d`)
: bytes
| MAXVALUE
< row
: bytes
| p_range_columns
: bytes
| p0
: bytes
| RANGE COLUMNS
: bytes
| `a`,`b`
: bytes
| 10,'m'
< row
: bytes
| p_range_columns
: bytes
| p1
: bytes
| RANGE COLUMNS
: bytes
| `a`,`b`
: bytes
| MAXVALUE,MAXVALUE
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, ''), IFNULL(TABLE_COLLATION, ''), QUOTE(IFNULL(TABLE_COMMENT, '')),
| 			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''), IFNULL(TIDB_PK_TYPE, ''), IFNULL(CREATE_OPTIONS, '')
| 		FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_mdload
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
| p_hash
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
< row
: bytes
| p_hash_named
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
< row
: bytes
| p_key
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
| partitioned
< row
: bytes
| p_list
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
< row
: bytes
| p_list_columns
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
< row
: bytes
| p_range
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
< row
: bytes
| p_range_columns
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
| omni_mdload
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SHOW CREATE TABLE `p_hash`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p_hash
: bytes
| CREATE TABLE `p_hash` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY HASH (`id`) PARTITIONS 4
> query
| SHOW CREATE TABLE `p_hash_named`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p_hash_named
: bytes
| CREATE TABLE `p_hash_named` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY HASH (`id`)
| (PARTITION `x0`,
|  PARTITION `x1`)
> query
| SHOW CREATE TABLE `p_key`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p_key
: bytes
| CREATE TABLE `p_key` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`id`) PARTITIONS 3
> query
| SHOW CREATE TABLE `p_list`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p_list
: bytes
| CREATE TABLE `p_list` (
|   `id` int DEFAULT NULL,
|   `r` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY LIST (`r`)
| (PARTITION `pa` VALUES IN (1,2),
|  PARTITION `pb` VALUES IN (3,NULL),
|  PARTITION `pd` DEFAULT)
> query
| SHOW CREATE TABLE `p_list_columns`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p_list_columns
: bytes
| CREATE TABLE `p_list_columns` (
|   `id` int DEFAULT NULL,
|   `c` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY LIST COLUMNS(`c`)
| (PARTITION `pa` VALUES IN ('a','b'),
|  PARTITION `pb` VALUES IN ('c'))
> query
| SHOW CREATE TABLE `p_range`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p_range
: bytes
| CREATE TABLE `p_range` (
|   `id` int DEFAULT NULL,
|   `d` date DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY RANGE (YEAR(`d`))
| (PARTITION `z1` VALUES LESS THAN (2000),
|  PARTITION `a2` VALUES LESS THAN (2010),
|  PARTITION `pmax` VALUES LESS THAN (MAXVALUE))
> query
| SHOW CREATE TABLE `p_range_columns`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p_range_columns
: bytes
| CREATE TABLE `p_range_columns` (
|   `a` int DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY RANGE COLUMNS(`a`,`b`)
| (PARTITION `p0` VALUES LESS THAN (10,'m'),
|  PARTITION `p1` VALUES LESS THAN (MAXVALUE,MAXVALUE))
> exec
| USE test
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_mdload
< ok 0 0
> exec
| CREATE DATABASE omni_mdload
< ok 0 0
> exec
| USE omni_mdload
< ok 0 0
> exec
| CREATE TABLE ar_bits (id bigint PRIMARY KEY AUTO_RANDOM(5), v int)
< ok 0 0
> exec
| CREATE TABLE ar_range (id bigint AUTO_RANDOM(4, 54), v int, PRIMARY KEY (id) CLUSTERED)
< ok 0 0
> exec
| CREATE TABLE ar_default (id bigint AUTO_RANDOM, PRIMARY KEY (id))
< ok 0 0
> exec
| CREATE TABLE sharded (id int, v int) SHARD_ROW_ID_BITS=4
< ok 0 0
> exec
| CREATE TABLE sharded_pk (id int, v int, PRIMARY KEY (id) NONCLUSTERED) SHARD_ROW_ID_BITS=3
< ok 0 0
> exec
| CREATE TABLE nonclustered (a int, b varchar(10), PRIMARY KEY (a, b) NONCLUSTERED)
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_mdload
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
| ar_bits
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
| ar_default
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
| ar_range
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
| nonclustered
: bytes
| PRIMARY
: bytes
| a
: int64 0
: bytes
|
: null
: null
< row
: bytes
| nonclustered
: bytes
| PRIMARY
: bytes
| b
: int64 0
: bytes
|
: null
: null
< row
: bytes
| sharded_pk
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
| omni_mdload
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
| ar_bits
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
< row
: bytes
| ar_bits
: bytes
| v
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
| ar_default
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
< row
: bytes
| ar_range
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
< row
: bytes
| ar_range
: bytes
| v
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
| nonclustered
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
< row
: bytes
| nonclustered
: bytes
| b
: int64 2
: null
: bytes
| NO
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
< row
: bytes
| sharded
: bytes
| id
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
| sharded
: bytes
| v
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
| sharded_pk
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
| sharded_pk
: bytes
| v
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
| omni_mdload
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
| omni_mdload
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
| omni_mdload
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_mdload
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
| omni_mdload
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
| ar_bits
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
< row
: bytes
| ar_default
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
< row
: bytes
| ar_range
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| PK_AUTO_RANDOM_BITS=4, RANGE BITS=54
: bytes
| CLUSTERED
: bytes
|
< row
: bytes
| nonclustered
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
| sharded
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
< row
: bytes
| sharded_pk
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| SHARD_BITS=3
: bytes
| NONCLUSTERED
: bytes
|
> query
| SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?
: string
| omni_mdload
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SHOW CREATE TABLE `ar_bits`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ar_bits
: bytes
| CREATE TABLE `ar_bits` (
|   `id` bigint NOT NULL /*T![auto_rand] AUTO_RANDOM(5) */,
|   `v` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `ar_default`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ar_default
: bytes
| CREATE TABLE `ar_default` (
|   `id` bigint NOT NULL /*T![auto_rand] AUTO_RANDOM(5) */,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `ar_range`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ar_range
: bytes
| CREATE TABLE `ar_range` (
|   `id` bigint NOT NULL /*T![auto_rand] AUTO_RANDOM(4, 54) */,
|   `v` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `nonclustered`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| nonclustered
: bytes
| CREATE TABLE `nonclustered` (
|   `a` int NOT NULL,
|   `b` varchar(10) NOT NULL,
|   PRIMARY KEY (`a`,`b`) /*T![clustered_index] NONCLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `sharded`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sharded
: bytes
| CREATE TABLE `sharded` (
|   `id` int DEFAULT NULL,
|   `v` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T! SHARD_ROW_ID_BITS=4 */
> query
| SHOW CREATE TABLE `sharded_pk`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sharded_pk
: bytes
| CREATE TABLE `sharded_pk` (
|   `id` int NOT NULL,
|   `v` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T! SHARD_ROW_ID_BITS=3 */
> exec
| USE test
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_mdload
< ok 0 0
> exec
| CREATE DATABASE omni_mdload
< ok 0 0
> exec
| USE omni_mdload
< ok 0 0
> exec
| CREATE TABLE ex (id int PRIMARY KEY, j json, name varchar(20),
|   KEY ix_quote ((cast(json_unquote(json_extract(j, '$."it''s"')) as char(10)))),
|   KEY ix_backslash ((cast(json_extract(j, '$."a\\b"') as char(10)))),
|   KEY ix_mixed ((lower(name)), id),
|   KEY ix_prefix (name(4), id),
|   KEY ix_array ((cast(j->'$.tags' as unsigned array))))
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_mdload
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
| ex
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
| ex
: bytes
| ix_array
: bytes
| NULL
: int64 1
: bytes
|
: null
: bytes
| cast(json_extract(`j`, _utf8mb4'$.tags') as unsigned array)
< row
: bytes
| ex
: bytes
| ix_backslash
: bytes
| NULL
: int64 1
: bytes
|
: null
: bytes
| cast(json_extract(`j`, _utf8mb4'$."a\\b"') as char(10))
< row
: bytes
| ex
: bytes
| ix_mixed
: bytes
| NULL
: int64 1
: bytes
|
: null
: bytes
| lower(`name`)
< row
: bytes
| ex
: bytes
| ix_mixed
: bytes
| id
: int64 1
: bytes
|
: null
: null
< row
: bytes
| ex
: bytes
| ix_prefix
: bytes
| name
: int64 1
: bytes
|
: int64 4
: null
< row
: bytes
| ex
: bytes
| ix_prefix
: bytes
| id
: int64 1
: bytes
|
: null
: null
< row
: bytes
| ex
: bytes
| ix_quote
: bytes
| NULL
: int64 1
: bytes
|
: null
: bytes
| cast(json_unquote(json_extract(`j`, _utf8mb4'$."it''s"')) as char(10))
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, ''), ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE, COLUMN_TYPE, IFNULL(CHARACTER_SET_NAME, ''), IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT), GENERATION_EXPRESSION, EXTRA
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mdload
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
| ex
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
| ex
: bytes
| j
: int64 2
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
< row
: bytes
| ex
: bytes
| name
: int64 3
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
| omni_mdload
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
| omni_mdload
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
| omni_mdload
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_mdload
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
| omni_mdload
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
| ex
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
| omni_mdload
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SHOW CREATE TABLE `ex`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ex
: bytes
| CREATE TABLE `ex` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `ix_quote` ((cast(json_unquote(json_extract(`j`, _utf8mb4'$."it''s"')) as char(10)))),
|   KEY `ix_backslash` ((cast(json_extract(`j`, _utf8mb4'$."a\\b"') as char(10)))),
|   KEY `ix_mixed` ((lower(`name`)),`id`),
|   KEY `ix_prefix` (`name`(4),`id`),
|   KEY `ix_array` ((cast(json_extract(`j`, _utf8mb4'$.tags') as unsigned array)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| USE test
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_mdload
< ok 0 0
> exec
| CREATE DATABASE omni_mdload
< ok 0 0
> exec
| USE omni_mdload
< ok 0 0
> exec
| CREATE TABLE parent (id int PRIMARY KEY, k int COMMENT 'it''s a \\ path', a int, b int, UNIQUE KEY uk (k), UNIQUE KEY uab (a, b)) COMMENT 'table''s \\ comment'
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, pid int, k int, a int, b int,
|   CONSTRAINT fk_cascade FOREIGN KEY (pid) REFERENCES parent (id) ON DELETE CASCADE,
|   CONSTRAINT fk_qualified FOREIGN KEY (k) REFERENCES omni_mdload.parent (k) ON UPDATE SET NULL,
|   CONSTRAINT fk_pair FOREIGN KEY (a, b) REFERENCES parent (a, b))
< ok 0 0
> exec
| CREATE VIEW v_join AS SELECT c.id, p.k AS pk, omni_mdload.c.pid FROM omni_mdload.child c JOIN parent p ON c.pid = p.id
< ok 0 0
> exec
| CREATE VIEW v_alias AS SELECT omni_mdload.id FROM parent omni_mdload WHERE omni_mdload.k > 1
< ok 0 0
> exec
| CREATE VIEW v_nested AS SELECT id, count(*) AS n FROM v_join GROUP BY id
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_mdload
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
| child
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
| child
: bytes
| fk_cascade
: bytes
| pid
: int64 1
: bytes
|
: null
: null
< row
: bytes
| child
: bytes
| fk_pair
: bytes
| a
: int64 1
: bytes
|
: null
: null
< row
: bytes
| child
: bytes
| fk_pair
: bytes
| b
: int64 1
: bytes
|
: null
: null
< row
: bytes
| child
: bytes
| fk_qualified
: bytes
| k
: int64 1
: bytes
|
: null
: null
< row
: bytes
| parent
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
| parent
: bytes
| uab
: bytes
| a
: int64 0
: bytes
|
: null
: null
< row
: bytes
| parent
: bytes
| uab
: bytes
| b
: int64 0
: bytes
|
: null
: null
< row
: bytes
| parent
: bytes
| uk
: bytes
| k
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
| omni_mdload
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
| child
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
| child
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
| child
: bytes
| k
: int64 3
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
| child
: bytes
| a
: int64 4
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
| child
: bytes
| b
: int64 5
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
| parent
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
| parent
: bytes
| k
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
| 'it\'s a \\ path'
: bytes
|
: bytes
|
< row
: bytes
| parent
: bytes
| a
: int64 3
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
| parent
: bytes
| b
: int64 4
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
| v_alias
: bytes
| id
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
| v_join
: bytes
| id
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
| v_join
: bytes
| pk
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
| v_join
: bytes
| pid
: int64 3
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
| v_nested
: bytes
| id
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
| v_nested
: bytes
| n
: int64 2
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
| omni_mdload
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT" "string"
< row
: bytes
| v_alias
: bytes
| SELECT `omni_mdload`.`id` AS `id` FROM `omni_mdload`.`parent` AS `omni_mdload` WHERE `omni_mdload`.`k`>1
< row
: bytes
| v_join
: bytes
| SELECT `c`.`id` AS `id`,`p`.`k` AS `pk`,`omni_mdload`.`c`.`pid` AS `pid` FROM `omni_mdload`.`child` AS `c` JOIN `omni_mdload`.`parent` AS `p` ON `c`.`pid`=`p`.`id`
< row
: bytes
| v_nested
: bytes
| SELECT `id` AS `id`,COUNT(1) AS `n` FROM `omni_mdload`.`v_join` GROUP BY `id`
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME, kcu.REFERENCED_COLUMN_NAME, fks.DELETE_RULE, fks.UPDATE_RULE, fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA AND fks.TABLE_NAME = kcu.TABLE_NAME AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| omni_mdload
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
| child
: bytes
| fk_cascade
: bytes
| pid
: bytes
| omni_mdload
: bytes
| parent
: bytes
| id
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| child
: bytes
| fk_pair
: bytes
| a
: bytes
| omni_mdload
: bytes
| parent
: bytes
| a
: bytes
| NO ACTION
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| child
: bytes
| fk_pair
: bytes
| b
: bytes
| omni_mdload
: bytes
| parent
: bytes
| b
: bytes
| NO ACTION
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| child
: bytes
| fk_qualified
: bytes
| k
: bytes
| omni_mdload
: bytes
| parent
: bytes
| k
: bytes
| NO ACTION
: bytes
| SET NULL
: bytes
| NONE
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ? ORDER BY TABLE_NAME, CONSTRAINT_NAME
: string
| omni_mdload
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_mdload
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
| omni_mdload
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
| child
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
| parent
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: bytes
| 'table\'s \\ comment'
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
|
< row
: bytes
| v_alias
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
< row
: bytes
| v_join
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
< row
: bytes
| v_nested
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
| omni_mdload
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SHOW CREATE TABLE `child`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| child
: bytes
| CREATE TABLE `child` (
|   `id` int NOT NULL,
|   `pid` int DEFAULT NULL,
|   `k` int DEFAULT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `fk_cascade` (`pid`),
|   KEY `fk_qualified` (`k`),
|   KEY `fk_pair` (`a`,`b`),
|   CONSTRAINT `fk_cascade` FOREIGN KEY (`pid`) REFERENCES `parent` (`id`) ON DELETE CASCADE,
|   CONSTRAINT `fk_qualified` FOREIGN KEY (`k`) REFERENCES `parent` (`k`) ON UPDATE SET NULL,
|   CONSTRAINT `fk_pair` FOREIGN KEY (`a`,`b`) REFERENCES `parent` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `parent`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| parent
: bytes
| CREATE TABLE `parent` (
|   `id` int NOT NULL,
|   `k` int DEFAULT NULL COMMENT 'it''s a \\ path',
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `uk` (`k`),
|   UNIQUE KEY `uab` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='table''s \\ comment'
> query
| SHOW CREATE VIEW `v_alias`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_alias
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_alias` (`id`) AS SELECT `omni_mdload`.`id` AS `id` FROM `omni_mdload`.`parent` AS `omni_mdload` WHERE `omni_mdload`.`k`>1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `v_join`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_join
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_join` (`id`, `pk`, `pid`) AS SELECT `c`.`id` AS `id`,`p`.`k` AS `pk`,`omni_mdload`.`c`.`pid` AS `pid` FROM `omni_mdload`.`child` AS `c` JOIN `omni_mdload`.`parent` AS `p` ON `c`.`pid`=`p`.`id`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `v_nested`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_nested
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_nested` (`id`, `n`) AS SELECT `id` AS `id`,COUNT(1) AS `n` FROM `omni_mdload`.`v_join` GROUP BY `id`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| USE test
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_mdload
< ok 0 0
> exec
| CREATE DATABASE omni_mdload
< ok 0 0
> exec
| USE omni_mdload
< ok 0 0
> exec
| CREATE TABLE ck (id int PRIMARY KEY, c varchar(10), g int, CONSTRAINT ck_named CHECK (g > 0 AND c <> 'it''s'), CHECK (g < 100))
< ok 0 0
> exec
| CREATE TABLE defaults (id int PRIMARY KEY,
|   a timestamp(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
|   b datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   c varchar(10) DEFAULT 'it''s',
|   e varchar(36) DEFAULT (uuid()),
|   f date DEFAULT (curdate()),
|   g int DEFAULT '5',
|   h enum('x','y') NOT NULL DEFAULT 'y',
|   i decimal(10,2) DEFAULT '1.50' COMMENT 'price')
< ok 0 0
> exec
| CREATE TABLE gen_cols (a int, b int GENERATED ALWAYS AS (a * 2) VIRTUAL, c varchar(10) GENERATED ALWAYS AS (concat('it''s', a)) STORED)
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_mdload
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
| ck
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
| defaults
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
| omni_mdload
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
| ck
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
| ck
: bytes
| c
: int64 2
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
< row
: bytes
| ck
: bytes
| g
: int64 3
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
| defaults
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
| defaults
: bytes
| a
: int64 2
: bytes
| 'CURRENT_TIMESTAMP(3)'
: bytes
| YES
: bytes
| timestamp(3)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP(3)
< row
: bytes
| defaults
: bytes
| b
: int64 3
: bytes
| 'CURRENT_TIMESTAMP'
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
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| defaults
: bytes
| c
: int64 4
: bytes
| 'it\'s'
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
< row
: bytes
| defaults
: bytes
| e
: int64 5
: bytes
| 'uuid()'
: bytes
| YES
: bytes
| varchar(36)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
< row
: bytes
| defaults
: bytes
| f
: int64 6
: bytes
| 'CURRENT_DATE'
: bytes
| YES
: bytes
| date
: bytes
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
| defaults
: bytes
| g
: int64 7
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
| defaults
: bytes
| h
: int64 8
: bytes
| 'y'
: bytes
| NO
: bytes
| enum('x','y')
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
< row
: bytes
| defaults
: bytes
| i
: int64 9
: bytes
| '1.50'
: bytes
| YES
: bytes
| decimal(10,2)
: bytes
|
: bytes
|
: bytes
| 'price'
: bytes
|
: bytes
|
< row
: bytes
| gen_cols
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
| gen_cols
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
| `a` * 2
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_cols
: bytes
| c
: int64 3
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
| concat(_utf8mb4'it''s', `a`)
: bytes
| STORED GENERATED
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_mdload
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
| omni_mdload
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
| omni_mdload
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| ck
: bytes
| ck_chk_1
: bytes
| (`g` < 100)
< row
: bytes
| ck
: bytes
| ck_named
: bytes
| (`g` > 0 and `c` != _utf8mb4'it''s')
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_mdload
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
| omni_mdload
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
| ck
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
| defaults
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
| gen_cols
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
| omni_mdload
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SHOW CREATE TABLE `ck`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ck
: bytes
| CREATE TABLE `ck` (
|   `id` int NOT NULL,
|   `c` varchar(10) DEFAULT NULL,
|   `g` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   CONSTRAINT `ck_named` CHECK ((`g` > 0 and `c` != _utf8mb4'it''s')),
|   CONSTRAINT `ck_chk_1` CHECK ((`g` < 100))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `defaults`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| defaults
: bytes
| CREATE TABLE `defaults` (
|   `id` int NOT NULL,
|   `a` timestamp(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
|   `b` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `c` varchar(10) DEFAULT 'it''s',
|   `e` varchar(36) DEFAULT (uuid()),
|   `f` date DEFAULT (CURRENT_DATE),
|   `g` int DEFAULT '5',
|   `h` enum('x','y') NOT NULL DEFAULT 'y',
|   `i` decimal(10,2) DEFAULT '1.50' COMMENT 'price',
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `gen_cols`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| gen_cols
: bytes
| CREATE TABLE `gen_cols` (
|   `a` int DEFAULT NULL,
|   `b` int GENERATED ALWAYS AS (`a` * 2) VIRTUAL,
|   `c` varchar(10) GENERATED ALWAYS AS (concat(_utf8mb4'it''s', `a`)) STORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| USE test
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_mdload
< ok 0 0
> exec
| CREATE DATABASE omni_mdload
< ok 0 0
> exec
| USE omni_mdload
< ok 0 0
> exec
| CREATE TABLE gen_exprs (a int, b int, s varchar(20), j json,
|  g1 int AS (a * (b + 1)),
|  g2 int AS ((a + b) * 2),
|  g3 int AS (a + b * 2),
|  g4 varchar(30) AS (concat('it''s', s)),
|  g5 varchar(30) AS (CASE WHEN a < 5 THEN 'x' WHEN a IS NULL THEN 'n' ELSE 'y' END) STORED,
|  g6 int AS (a IN (1, 2, 3)),
|  g7 int AS (s LIKE 'a%'),
|  g8 int AS (a BETWEEN 1 AND 5),
|  g9 char(10) AS (CAST(a AS char)),
|  g10 int AS (IF(a > 0, a, -a)),
|  g11 varchar(20) AS (concat('a\\b', s)),
|  g12 int AS (NOT (a > 1)),
|  g13 int AS (a IS NOT NULL),
|  g14 varchar(40) AS (j->>'$.x'),
|  g15 int AS (-a + 1),
|  g16 int AS (a DIV 2 MOD 3),
|  g17 int AS (coalesce(a, b, 0)))
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_mdload
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
| omni_mdload
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
| gen_exprs
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
| gen_exprs
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
| gen_exprs
: bytes
| s
: int64 3
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
< row
: bytes
| gen_exprs
: bytes
| j
: int64 4
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
< row
: bytes
| gen_exprs
: bytes
| g1
: int64 5
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
| `a` * (`b` + 1)
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g2
: int64 6
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
| (`a` + `b`) * 2
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g3
: int64 7
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
| `a` + `b` * 2
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g4
: int64 8
: null
: bytes
| YES
: bytes
| varchar(30)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| concat(_utf8mb4'it''s', `s`)
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g5
: int64 9
: null
: bytes
| YES
: bytes
| varchar(30)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| case when `a` < 5 then _utf8mb4'x' when `a` is null then _utf8mb4'n' else _utf8mb4'y' end
: bytes
| STORED GENERATED
< row
: bytes
| gen_exprs
: bytes
| g6
: int64 10
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
| `a` in (1,2,3)
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g7
: int64 11
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
| `s` like _utf8mb4'a%'
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g8
: int64 12
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
| `a` between 1 and 5
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g9
: int64 13
: null
: bytes
| YES
: bytes
| char(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| cast(`a` as char)
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g10
: int64 14
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
| if(`a` > 0, `a`, -`a`)
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g11
: int64 15
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
| concat(_utf8mb4'a\\b', `s`)
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g12
: int64 16
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
| not (`a` > 1)
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g13
: int64 17
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
| `a` is not null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g14
: int64 18
: null
: bytes
| YES
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
| json_unquote(json_extract(`j`, _utf8mb4'$.x'))
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g15
: int64 19
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
| -`a` + 1
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g16
: int64 20
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
| `a` div 2 % 3
: bytes
| VIRTUAL GENERATED
< row
: bytes
| gen_exprs
: bytes
| g17
: int64 21
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
| coalesce(`a`, `b`, 0)
: bytes
| VIRTUAL GENERATED
> query
| SELECT TABLE_NAME, VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_mdload
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
| omni_mdload
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
| omni_mdload
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_mdload
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
| omni_mdload
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
| gen_exprs
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
| omni_mdload
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SHOW CREATE TABLE `gen_exprs`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| gen_exprs
: bytes
| CREATE TABLE `gen_exprs` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   `j` json DEFAULT NULL,
|   `g1` int GENERATED ALWAYS AS (`a` * (`b` + 1)) VIRTUAL,
|   `g2` int GENERATED ALWAYS AS ((`a` + `b`) * 2) VIRTUAL,
|   `g3` int GENERATED ALWAYS AS (`a` + `b` * 2) VIRTUAL,
|   `g4` varchar(30) GENERATED ALWAYS AS (concat(_utf8mb4'it''s', `s`)) VIRTUAL,
|   `g5` varchar(30) GENERATED ALWAYS AS (case when `a` < 5 then _utf8mb4'x' when `a` is null then _utf8mb4'n' else _utf8mb4'y' end) STORED,
|   `g6` int GENERATED ALWAYS AS (`a` in (1,2,3)) VIRTUAL,
|   `g7` int GENERATED ALWAYS AS (`s` like _utf8mb4'a%') VIRTUAL,
|   `g8` int GENERATED ALWAYS AS (`a` between 1 and 5) VIRTUAL,
|   `g9` char(10) GENERATED ALWAYS AS (cast(`a` as char)) VIRTUAL,
|   `g10` int GENERATED ALWAYS AS (if(`a` > 0, `a`, -`a`)) VIRTUAL,
|   `g11` varchar(20) GENERATED ALWAYS AS (concat(_utf8mb4'a\\b', `s`)) VIRTUAL,
|   `g12` int GENERATED ALWAYS AS (not (`a` > 1)) VIRTUAL,
|   `g13` int GENERATED ALWAYS AS (`a` is not null) VIRTUAL,
|   `g14` varchar(40) GENERATED ALWAYS AS (json_unquote(json_extract(`j`, _utf8mb4'$.x'))) VIRTUAL,
|   `g15` int GENERATED ALWAYS AS (-`a` + 1) VIRTUAL,
|   `g16` int GENERATED ALWAYS AS (`a` div 2 % 3) VIRTUAL,
|   `g17` int GENERATED ALWAYS AS (coalesce(`a`, `b`, 0)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| USE test
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_mdload
< ok 0 0
> exec
| CREATE DATABASE omni_mdload
< ok 0 0
> exec
| USE omni_mdload
< ok 0 0
> exec
| CREATE TABLE bits (id int PRIMARY KEY, b1 bit(1) NOT NULL DEFAULT b'0', b8 bit(8) DEFAULT b'101', b16 bit(16) DEFAULT 5, b64 bit(64))
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, COLUMN_NAME, NON_UNIQUE, INDEX_COMMENT, SUB_PART, EXPRESSION
| 		FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX
: string
| omni_mdload
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
| bits
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
| omni_mdload
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
| bits
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
| bits
: bytes
| b1
: int64 2
: bytes
| 'b\'0\''
: bytes
| NO
: bytes
| bit(1)
: bytes
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
| bits
: bytes
| b8
: int64 3
: bytes
| 'b\'101\''
: bytes
| YES
: bytes
| bit(8)
: bytes
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
| bits
: bytes
| b16
: int64 4
: bytes
| 'b\'101\''
: bytes
| YES
: bytes
| bit(16)
: bytes
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
| bits
: bytes
| b64
: int64 5
: null
: bytes
| YES
: bytes
| bit(64)
: bytes
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
| omni_mdload
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
| omni_mdload
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
| omni_mdload
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
> query
| SELECT TABLE_NAME, PARTITION_NAME, PARTITION_METHOD, PARTITION_EXPRESSION, PARTITION_DESCRIPTION
| 		FROM INFORMATION_SCHEMA.PARTITIONS WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
: string
| omni_mdload
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
| omni_mdload
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
| bits
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
| omni_mdload
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SHOW CREATE TABLE `bits`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| bits
: bytes
| CREATE TABLE `bits` (
|   `id` int NOT NULL,
|   `b1` bit(1) NOT NULL DEFAULT b'0',
|   `b8` bit(8) DEFAULT b'101',
|   `b16` bit(16) DEFAULT b'101',
|   `b64` bit(64) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| USE test
< ok 0 0
> exec
| SET GLOBAL tidb_enable_check_constraint = OFF
< ok 0 0
