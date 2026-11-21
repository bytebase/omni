-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDumpRestoresExtendedTiDBObjects default#1
> exec
| USE `omni_ext_roundtrip`
< ok 0 0
> exec
| SET GLOBAL tidb_enable_check_constraint = 1
< ok 0 0
> exec
| CREATE TABLE t_chk_a (
| 	id INT PRIMARY KEY,
| 	a INT NOT NULL,
| 	CONSTRAINT chk_a CHECK (a > 0)
| );
< ok 0 0
> exec
| CREATE TABLE t_chk_b (
| 	id INT PRIMARY KEY,
| 	b INT NOT NULL,
| 	CONSTRAINT chk_b CHECK (b < 100) NOT ENFORCED
| );
< ok 0 0
> exec
| CREATE GLOBAL TEMPORARY TABLE gtt_del (
| 	id INT PRIMARY KEY
| ) ON COMMIT DELETE ROWS;
< ok 0 0
> exec
| CREATE TABLE t_part_global (
| 	id INT NOT NULL,
| 	k INT NOT NULL,
| 	PRIMARY KEY (id) NONCLUSTERED,
| 	UNIQUE KEY uk_k (k) /*T![global_index] GLOBAL */
| ) PARTITION BY HASH (id) PARTITIONS 2;
< ok 0 0
> exec
| CREATE TABLE t_cmp_ws (
| 	id INT PRIMARY KEY
| ) COMPRESSION='custom value';
< ok 0 0
> exec
| CREATE TABLE t_cmp_bt (
| 	id INT PRIMARY KEY
| ) COMPRESSION='`custom`';
< ok 0 0
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
| omni_ext_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "SEQ_IN_INDEX" "BIGINT"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| gtt_del
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| t_chk_a
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| t_chk_b
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| t_cmp_bt
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| t_cmp_ws
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| t_part_global
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| t_part_global
: bytes
| uk_k
: bytes
| k
: int64 1
: int64 0
: bytes
|
: null
: null
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
| omni_ext_roundtrip
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
| gtt_del
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
| t_chk_a
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
| t_chk_a
: bytes
| a
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
< row
: bytes
| t_chk_b
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
| t_chk_b
: bytes
| b
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
< row
: bytes
| t_cmp_bt
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
| t_cmp_ws
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
| t_part_global
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
| t_part_global
: bytes
| k
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
| omni_ext_roundtrip
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
| omni_ext_roundtrip
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
| omni_ext_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| t_chk_a
: bytes
| chk_a
: bytes
| (`a` > 0)
< row
: bytes
| t_chk_b
: bytes
| chk_b
: bytes
| (`b` < 100)
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
| omni_ext_roundtrip
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
| omni_ext_roundtrip
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
| omni_ext_roundtrip
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
< row
: bytes
| t_part_global
: bytes
| p0
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: bytes
|
: bytes
|
< row
: bytes
| t_part_global
: bytes
| p1
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: bytes
|
: bytes
|
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
| omni_ext_roundtrip
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
< row
: bytes
| gtt_del
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
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
| t_chk_a
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
| p_sync_ext
< row
: bytes
| t_chk_b
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
| p_sync_ext
< row
: bytes
| t_cmp_bt
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
| p_sync_ext
< row
: bytes
| t_cmp_ws
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
| p_sync_ext
< row
: bytes
| t_part_global
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
| partitioned
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
| p_sync_ext
> query
| SHOW CREATE TABLE `omni_ext_roundtrip`.`gtt_del`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| gtt_del
: bytes
| CREATE GLOBAL TEMPORARY TABLE `gtt_del` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ON COMMIT DELETE ROWS
> query
| SHOW CREATE TABLE `omni_ext_roundtrip`.`t_chk_a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chk_a
: bytes
| CREATE TABLE `t_chk_a` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   CONSTRAINT `chk_a` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
> query
| SHOW CREATE TABLE `omni_ext_roundtrip`.`t_chk_b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chk_b
: bytes
| CREATE TABLE `t_chk_b` (
|   `id` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   CONSTRAINT `chk_b` CHECK ((`b` < 100)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
> query
| SHOW CREATE TABLE `omni_ext_roundtrip`.`t_cmp_bt`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_cmp_bt
: bytes
| CREATE TABLE `t_cmp_bt` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMPRESSION='`custom`' /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
> query
| SHOW CREATE TABLE `omni_ext_roundtrip`.`t_cmp_ws`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_cmp_ws
: bytes
| CREATE TABLE `t_cmp_ws` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMPRESSION='custom value' /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
> query
| SHOW CREATE TABLE `omni_ext_roundtrip`.`t_part_global`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_global
: bytes
| CREATE TABLE `t_part_global` (
|   `id` int NOT NULL,
|   `k` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */,
|   UNIQUE KEY `uk_k` (`k`) /*T![global_index] GLOBAL */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
| PARTITION BY HASH (`id`) PARTITIONS 2
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME,
| 			IFNULL(TIDB_PLACEMENT_POLICY_NAME, '')
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_ext_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< column "IFNULL(TIDB_PLACEMENT_POLICY_NAME, '')" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| p_sync_ext
> exec
| DROP DATABASE `omni_ext_roundtrip`
< ok 0 0
> exec
| CREATE DATABASE `omni_ext_roundtrip`
< ok 0 0
> exec
| USE `omni_ext_roundtrip`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| ALTER DATABASE PLACEMENT POLICY=`p_sync_ext`
< ok 0 0
> exec
| CREATE GLOBAL TEMPORARY TABLE `gtt_del` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ON COMMIT DELETE ROWS
< ok 0 0
> exec
| CREATE TABLE `t_chk_a` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   CONSTRAINT `chk_a` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
< ok 0 0
> exec
| CREATE TABLE `t_chk_b` (
|   `id` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   CONSTRAINT `chk_b` CHECK (`b` < 100) NOT ENFORCED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
< ok 0 0
> exec
| CREATE TABLE `t_cmp_bt` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMPRESSION='`custom`' /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
< ok 0 0
> exec
| CREATE TABLE `t_cmp_ws` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMPRESSION='custom value' /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
< ok 0 0
> exec
| CREATE TABLE `t_part_global` (
|   `id` int NOT NULL,
|   `k` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */,
|   UNIQUE KEY `uk_k` (`k`) /*T![global_index] GLOBAL */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
| /*!50100 PARTITION BY HASH (`id`)
| (PARTITION p0 ENGINE=InnoDB,
|  PARTITION p1 ENGINE=InnoDB) */
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
| omni_ext_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "SEQ_IN_INDEX" "BIGINT"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_COMMENT" "VARCHAR"
< column "SUB_PART" "BIGINT"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| gtt_del
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| t_chk_a
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| t_chk_b
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| t_cmp_bt
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| t_cmp_ws
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| t_part_global
: bytes
| PRIMARY
: bytes
| id
: int64 1
: int64 0
: bytes
|
: null
: null
< row
: bytes
| t_part_global
: bytes
| uk_k
: bytes
| k
: int64 1
: int64 0
: bytes
|
: null
: null
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
| omni_ext_roundtrip
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
| gtt_del
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
| t_chk_a
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
| t_chk_a
: bytes
| a
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
< row
: bytes
| t_chk_b
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
| t_chk_b
: bytes
| b
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
< row
: bytes
| t_cmp_bt
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
| t_cmp_ws
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
| t_part_global
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
| t_part_global
: bytes
| k
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
| omni_ext_roundtrip
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
| omni_ext_roundtrip
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
| omni_ext_roundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< row
: bytes
| t_chk_a
: bytes
| chk_a
: bytes
| (`a` > 0)
< row
: bytes
| t_chk_b
: bytes
| chk_b
: bytes
| (`b` < 100)
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
| omni_ext_roundtrip
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
| omni_ext_roundtrip
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
| omni_ext_roundtrip
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
< row
: bytes
| t_part_global
: bytes
| p0
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: bytes
|
: bytes
|
< row
: bytes
| t_part_global
: bytes
| p1
: null
: bytes
| HASH
: null
: bytes
| `id`
: null
: bytes
|
: bytes
|
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
| omni_ext_roundtrip
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
< row
: bytes
| gtt_del
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
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
| t_chk_a
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
| p_sync_ext
< row
: bytes
| t_chk_b
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
| p_sync_ext
< row
: bytes
| t_cmp_bt
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
| p_sync_ext
< row
: bytes
| t_cmp_ws
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
|
: bytes
| ''
: bytes
| NOT_SHARDED(PK_IS_HANDLE)
: bytes
| CLUSTERED
: bytes
| p_sync_ext
< row
: bytes
| t_part_global
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_bin
: int64 0
: int64 0
: int64 0
: int64 0
: bytes
| partitioned
: bytes
| ''
: bytes
| NOT_SHARDED
: bytes
| NONCLUSTERED
: bytes
| p_sync_ext
> query
| SHOW CREATE TABLE `omni_ext_roundtrip`.`gtt_del`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| gtt_del
: bytes
| CREATE GLOBAL TEMPORARY TABLE `gtt_del` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ON COMMIT DELETE ROWS
> query
| SHOW CREATE TABLE `omni_ext_roundtrip`.`t_chk_a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chk_a
: bytes
| CREATE TABLE `t_chk_a` (
|   `id` int NOT NULL,
|   `a` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   CONSTRAINT `chk_a` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
> query
| SHOW CREATE TABLE `omni_ext_roundtrip`.`t_chk_b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chk_b
: bytes
| CREATE TABLE `t_chk_b` (
|   `id` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   CONSTRAINT `chk_b` CHECK ((`b` < 100)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
> query
| SHOW CREATE TABLE `omni_ext_roundtrip`.`t_cmp_bt`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_cmp_bt
: bytes
| CREATE TABLE `t_cmp_bt` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMPRESSION='`custom`' /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
> query
| SHOW CREATE TABLE `omni_ext_roundtrip`.`t_cmp_ws`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_cmp_ws
: bytes
| CREATE TABLE `t_cmp_ws` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMPRESSION='custom value' /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
> query
| SHOW CREATE TABLE `omni_ext_roundtrip`.`t_part_global`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_global
: bytes
| CREATE TABLE `t_part_global` (
|   `id` int NOT NULL,
|   `k` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */,
|   UNIQUE KEY `uk_k` (`k`) /*T![global_index] GLOBAL */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![placement] PLACEMENT POLICY=`p_sync_ext` */
| PARTITION BY HASH (`id`) PARTITIONS 2
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME,
| 			IFNULL(TIDB_PLACEMENT_POLICY_NAME, '')
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| omni_ext_roundtrip
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< column "IFNULL(TIDB_PLACEMENT_POLICY_NAME, '')" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| p_sync_ext
> exec
| SET GLOBAL tidb_enable_check_constraint = 0
< ok 0 0
