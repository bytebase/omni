-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestWalkThroughSweepTiDB/change0 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE A AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| CREATE SEQUENCE IF NOT EXISTS b;
< ok 0 0
> exec
| INSERT INTO c (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP TABLE q.v, a;
< error
| 1051 42S02 Unknown table 'q.v,shop.a'
> exec
| ALTER TABLE b ADD COLUMN x int, DROP COLUMN X;
< error
| 1347 HY000 'shop.b' is not BASE TABLE
> exec
| DROP SEQUENCE q.v;
< error
| 4139 HY000 Unknown SEQUENCE: 'q.v'
> exec
| ALTER TABLE s CHANGE x y int, MODIFY X bigint;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE b RENAME COLUMN z TO z;
< error
| 1347 HY000 'shop.b' is not BASE TABLE
> exec
| DROP SEQUENCE v, q.A;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.v,q.A'
> exec
| DROP SEQUENCE IF EXISTS c, A;
< ok 0 0
> exec
| ALTER TABLE A CHANGE x x int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE INDEX j ON v (X);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change0 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| b
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| SEQUENCE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change0 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP TEMPORARY TABLE s, A;
< error
| 1051 42S02 Unknown table 'shop.s,shop.A'
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP INDEX i ON b;
< error
| 1091 42000 index i doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| INSERT INTO q.b (x) VALUES (1);
< error
| 1146 42S02 Table 'q.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE INDEX IF NOT EXISTS x ON q.A (y);
< error
| 1146 42S02 Table 'q.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| TRUNCATE TABLE q.c;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE s AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS v;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP INDEX i ON q.s;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS q.v (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TEMPORARY TABLE v (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE s AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE q.a DROP INDEX x;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP SEQUENCE IF EXISTS q.b;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE SEQUENCE v;
< error
| 1050 42S01 Table 'shop.v' already exists
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1105 HY000 insert into sequence b is not supported now
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE v (id int PRIMARY KEY, x int, y int);
< error
| 1050 42S01 Table 'shop.v' already exists
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE VIEW q.v AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE v ADD CONSTRAINT I UNIQUE INDEX j (X);
< error
| 8200 HY000 TiDB doesn't support ALTER TABLE for local temporary table
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist

== TestWalkThroughSweepTiDB/change1 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE UNIQUE INDEX IF NOT EXISTS i ON q.c (y);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DROP SEQUENCE q.A, v;
< error
| 4139 HY000 Unknown SEQUENCE: 'q.A,shop.v'
> exec
| TRUNCATE TABLE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE A ADD COLUMN x int, DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| DROP TABLE IF EXISTS A;
< ok 0 0
> exec
| ALTER TABLE b ADD COLUMN X int, DROP COLUMN y;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE b DROP FOREIGN KEY fk_1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE a LIKE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP VIEW a;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| ALTER TABLE c DROP COLUMN z, ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE q.v DROP INDEX j;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| DROP SEQUENCE IF EXISTS v;
< ok 0 0
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change1 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			

== TestWalkThroughSweepTiDB/change1 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`v`
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE s AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> query
| SHOW CREATE TABLE `q`.`v`
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE c DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `q`.`v`
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE q.s LIKE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`v`
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP INDEX I ON a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE b DROP INDEX I;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE a ADD UNIQUE INDEX i (x);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE a LIKE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TEMPORARY TABLE v (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE a RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP INDEX i ON q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE q.A LIKE q.v;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE v ADD COLUMN X int, DROP COLUMN y;
< error
| 8200 HY000 TiDB doesn't support ALTER TABLE for local temporary table
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| TRUNCATE TABLE v;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE q.b CHANGE x x int;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| TRUNCATE TABLE A;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE c DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP INDEX i ON c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE INDEX i ON s (X);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP SEQUENCE b, q.v;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.b,q.v'
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist

== TestWalkThroughSweepTiDB/change10 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE v DROP INDEX x, ADD INDEX j (id);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE A CHANGE x y int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| RENAME TABLE A TO q.b, v TO q.s;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE q.A CHANGE x y int, MODIFY y bigint;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| CREATE INDEX i ON s (y);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE q.s ADD CONSTRAINT f FOREIGN KEY (x) REFERENCES s (id);
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| CREATE VIEW q.b AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TEMPORARY TABLE q.b (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TEMPORARY TABLE v (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change10 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	a	I	
| index	shop	s	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	a	I	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	a	I	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	a	I	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	a	I	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	a	I	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	a	I	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE

== TestWalkThroughSweepTiDB/change10 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE IF EXISTS b;
< ok 0 0
> exec
| CREATE OR REPLACE VIEW q.c AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| RENAME TABLE q.b TO c, a TO b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE v DROP INDEX j;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP INDEX I ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| TRUNCATE TABLE s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE VIEW q.c AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE c RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| INSERT INTO s (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| DROP VIEW v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| ALTER TABLE c DROP COLUMN IF EXISTS y, ADD COLUMN y int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE s (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP INDEX i ON q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE v DROP FOREIGN KEY g;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE a (x int, y int, KEY I (x));
< ok 0 0
> exec
| ALTER TABLE s ADD UNIQUE INDEX x (X);
< ok 0 0
> exec
| ALTER TABLE q.v DROP FOREIGN KEY f, ADD COLUMN X int;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| INSERT INTO q.s (x) VALUES (1);
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| ALTER TABLE q.s ADD CONSTRAINT g FOREIGN KEY g (x) REFERENCES b (id);
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| CREATE VIEW a AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'shop.a' already exists
> exec
| CREATE SEQUENCE q.A;
< error
| 1049 42000 Unknown database 'q'

== TestWalkThroughSweepTiDB/change11 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| DROP TABLE a;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| CREATE SEQUENCE s;
< ok 0 0
> exec
| CREATE TABLE b LIKE q.v;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| ALTER TABLE s ADD CONSTRAINT g FOREIGN KEY f (x) REFERENCES A (id);
< error
| 1347 HY000 'shop.s' is not BASE TABLE
> exec
| TRUNCATE TABLE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS a, b;
< ok 0 0
> exec
| CREATE TABLE s (id int PRIMARY KEY, x int, y int);
< error
| 1050 42S01 Table 'shop.s' already exists
> exec
| ALTER TABLE s DROP PRIMARY KEY;
< error
| 1347 HY000 'shop.s' is not BASE TABLE
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change11 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| s
: bytes
| SEQUENCE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		s	SEQUENCE

== TestWalkThroughSweepTiDB/change11 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE a LIKE A;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| RENAME TABLE b TO c;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE A ADD UNIQUE INDEX x (x);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE VIEW q.A AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE q.v ADD PRIMARY KEY (x);
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| ALTER TABLE v ADD CONSTRAINT f FOREIGN KEY (x) REFERENCES b (id);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP SEQUENCE q.A;
< error
| 4139 HY000 Unknown SEQUENCE: 'q.A'
> exec
| ALTER TABLE v ADD COLUMN z int, DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE s ADD PRIMARY KEY (x);
< error
| 1347 HY000 'shop.s' is not BASE TABLE
> exec
| CREATE INDEX IF NOT EXISTS j ON b (x);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS a LIKE c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE SEQUENCE s;
< error
| 1050 42S01 Table 'shop.s' already exists
> exec
| ALTER TABLE c DROP INDEX i;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE UNIQUE INDEX IF NOT EXISTS I ON v (z);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS s (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE A AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| CREATE TABLE IF NOT EXISTS s (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE q.b AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| CREATE SEQUENCE q.a;
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP INDEX x ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist

== TestWalkThroughSweepTiDB/change12 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE b ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TABLE IF EXISTS b, s;
< ok 0 0
> exec
| ALTER TABLE q.a ADD FOREIGN KEY F (x) REFERENCES q.v (id);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE v DROP INDEX x;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS c LIKE q.A;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE q.v RENAME COLUMN x TO y;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| TRUNCATE TABLE A;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS b (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE VIEW a AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE TABLE A (x int, y int, KEY i (x), CONSTRAINT f FOREIGN KEY (y) REFERENCES a (id));
< error
| 3734 HY000 Failed to add the foreign key constraint. Missing column 'id' for constraint 'f' in the referenced table 'a'
> exec
| CREATE TABLE c AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| CREATE TABLE s AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP TEMPORARY TABLE q.a;
< error
| 1051 42S02 Unknown table 'q.a'
> exec
| CREATE TABLE q.b AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change12 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| b
: bytes
| PRIMARY
: int64 0
< row
: bytes
| b
: bytes
| x
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| x
< row
: bytes
| b
: bytes
| id
< row
: bytes
| b
: bytes
| x
< row
: bytes
| b
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| a
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| VIEW
: bytes
|
< row
: bytes
| b
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE

== TestWalkThroughSweepTiDB/change12 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE s ADD COLUMN y int, DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE A DROP FOREIGN KEY g;
< error
| 1347 HY000 'shop.A' is not BASE TABLE
> exec
| ALTER TABLE s ADD COLUMN X int, DROP COLUMN y;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE UNIQUE INDEX j ON A (x);
< error
| 1105 HY000 Unknown charset 
> exec
| ALTER TABLE b ADD COLUMN x int AFTER x;
< error
| 1060 42S21 Duplicate column name 'x'
> exec
| DROP INDEX j ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE v ADD FOREIGN KEY F (x) REFERENCES b (id);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE b DROP INDEX j, ADD INDEX I (id);
< error
| 1091 42000 index j doesn't exist
> exec
| DROP VIEW IF EXISTS a;
< ok 0 0
> exec
| ALTER TABLE v DROP INDEX x, ADD INDEX x (id);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE s ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE A DROP INDEX I, ADD INDEX I (id);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE a DROP INDEX i;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP VIEW IF EXISTS q.a;
< ok 0 0
> exec
| ALTER TABLE a CHANGE X X int, MODIFY X bigint;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE c DROP FOREIGN KEY g;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE INDEX x ON a (y);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE s AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| TRUNCATE TABLE s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| DROP INDEX I ON A;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE q.a DROP INDEX j;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE q.v AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP VIEW IF EXISTS c;
< ok 0 0
> exec
| ALTER TABLE s ADD CONSTRAINT F FOREIGN KEY g (x) REFERENCES s (id);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE q.v DROP COLUMN IF EXISTS z;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| ALTER TABLE c ADD COLUMN x int, DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.c' doesn't exist

== TestWalkThroughSweepTiDB/change13 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE s RENAME COLUMN z TO X;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE SEQUENCE q.s;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE c DROP FOREIGN KEY f;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP INDEX j ON q.c;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| INSERT INTO q.A (x) VALUES (1);
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| INSERT INTO A (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE VIEW A AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE INDEX i ON q.A (y);
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| UPDATE s SET x = 1;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE q.v RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| CREATE DATABASE q;
< ok 0 0
> exec
| CREATE TABLE q.c (x int, y int, KEY j (x));
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change13 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| A
: bytes
| x
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| A
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| A
: bytes
| VIEW
: bytes
|
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	c	x	
| column	q	c	y	
| column	shop	A	x	
| database	q			
| database	shop			
| database	test			
| index	q	c	j	
| relation	q		c	BASE TABLE
| relation	shop		A	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	c	x	
| column	q	c	y	
| column	shop	A	x	
| database	q			
| database	shop			
| database	test			
| index	q	c	j	
| relation	q		c	BASE TABLE
| relation	shop		A	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	c	x	
| column	q	c	y	
| column	shop	A	x	
| database	q			
| database	shop			
| database	test			
| index	q	c	j	
| relation	q		c	BASE TABLE
| relation	shop		A	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	c	x	
| column	q	c	y	
| column	shop	A	x	
| database	q			
| database	shop			
| database	test			
| index	q	c	j	
| relation	q		c	BASE TABLE
| relation	shop		A	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	c	x	
| column	q	c	y	
| column	shop	A	x	
| database	q			
| database	shop			
| database	test			
| index	q	c	j	
| relation	q		c	BASE TABLE
| relation	shop		A	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	c	x	
| column	q	c	y	
| column	shop	A	x	
| database	q			
| database	shop			
| database	test			
| index	q	c	j	
| relation	q		c	BASE TABLE
| relation	shop		A	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	q	c	x	
| column	q	c	y	
| database	q			
| database	shop			
| database	test			
| index	q	c	j	
| relation	q		c	BASE TABLE

== TestWalkThroughSweepTiDB/change13 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`v`
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE b RENAME COLUMN X TO x;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `q`.`v`
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE q.c ADD COLUMN y int;
< error
| 1060 42S21 Duplicate column name 'y'
> query
| SHOW CREATE TABLE `q`.`v`
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE INDEX j ON s (X);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SHOW CREATE TABLE `q`.`v`
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE q.c CHANGE X y int;
< error
| 1060 42S21 Duplicate column name 'y'
> query
| SHOW CREATE TABLE `q`.`v`
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE q.A DROP COLUMN IF EXISTS x;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> query
| SHOW CREATE TABLE `q`.`v`
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP VIEW q.a, a;
< error
| 1051 42S02 Unknown table 'q.a'
> query
| SHOW CREATE TABLE `q`.`v`
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist

== TestWalkThroughSweepTiDB/change14 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a, v;
< ok 0 0
> exec
| ALTER TABLE q.c DROP COLUMN IF EXISTS X;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TABLE q.c (x int, y int, KEY x (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| INSERT INTO A (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| DROP VIEW IF EXISTS b;
< ok 0 0
> exec
| DROP SEQUENCE v;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.v'
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec
| UPDATE s SET x = 1;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| INSERT INTO a (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE q.s DROP INDEX i;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change14 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| relation	shop		A	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| relation	shop		A	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| relation	shop		A	BASE TABLE
| table foreign key	shop	A	g	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| relation	shop		A	BASE TABLE
| table foreign key	shop	A	g	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| relation	shop		A	BASE TABLE
| table foreign key	shop	A	g	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| relation	shop		A	BASE TABLE
| table foreign key	shop	A	g	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| relation	shop		A	BASE TABLE
| table foreign key	shop	A	g	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| relation	shop		A	BASE TABLE
| table foreign key	shop	A	g	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| relation	shop		A	BASE TABLE
| table foreign key	shop	A	g	

== TestWalkThroughSweepTiDB/change14 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE INDEX I (X);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| USE shop;
< ok 0 0
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE b ADD COLUMN x int, DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE v AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| ALTER TABLE c DROP COLUMN IF EXISTS X, ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| ALTER TABLE A RENAME COLUMN x TO y;
< error
| 1060 42S21 Duplicate column name 'y'
> exec
| ALTER TABLE A ADD CONSTRAINT g FOREIGN KEY (x) REFERENCES b (id);
< ok 0 0
> exec
| ALTER TABLE c DROP COLUMN IF EXISTS X, ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE c DROP FOREIGN KEY g;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE A DROP FOREIGN KEY fk_1;
< error
| 1091 42000 Can't DROP 'fk_1'; check that column/key exists
> exec
| CREATE TABLE A AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP INDEX x ON a;
< error
| 1553 HY000 Cannot drop index 'x': needed in a foreign key constraint
> exec
| RENAME TABLE q.s TO b;
< error
| 1146 42S02 Table 'q.s' doesn't exist

== TestWalkThroughSweepTiDB/change15 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE A ADD INDEX i (y);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE TABLE q.b LIKE q.v;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS v (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| ALTER TABLE s RENAME COLUMN X TO y;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| UPDATE c SET x = 1;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE q.A CHANGE X x int;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE s ADD COLUMN y int, DROP COLUMN y;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE TABLE q.a AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| ALTER TABLE s DROP INDEX i;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE b LIKE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE q.a ADD FOREIGN KEY f (x) REFERENCES A (id);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE INDEX j ON q.c (x);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE q.A ADD COLUMN z int, DROP COLUMN X;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE c ADD COLUMN X int, DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change15 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			

== TestWalkThroughSweepTiDB/change15 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE INDEX j ON s (z);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE q.A RENAME COLUMN z TO y;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE a CHANGE X z int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN X int, DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE q.b CHANGE y X int;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| INSERT INTO q.A (x) VALUES (1);
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN z int AFTER y;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE a DROP INDEX j, ADD INDEX i (id);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE a ADD COLUMN z int, DROP COLUMN y;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE A LIKE A;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE q.a ADD INDEX i (x);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE q.c DROP INDEX j;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE q.v DROP COLUMN IF EXISTS y;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE a ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| TRUNCATE TABLE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TABLE q.b;
< error
| 1051 42S02 Unknown table 'q.b'
> exec
| ALTER TABLE a DROP INDEX i;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP TABLE s;
< error
| 1051 42S02 Unknown table 'shop.s'

== TestWalkThroughSweepTiDB/change16 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE b DROP COLUMN y, ADD COLUMN X int;
< error
| 1060 42S21 Duplicate column name 'X'
> exec
| ALTER TABLE v ADD COLUMN x int AFTER z;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE v RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE c CHANGE x x int, CHANGE X z int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE q.c ADD UNIQUE INDEX x (y);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE INDEX x ON s (X);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE SEQUENCE q.a;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE SEQUENCE IF NOT EXISTS c;
< ok 0 0
> exec
| ALTER TABLE A RENAME COLUMN x TO z;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE TABLE v (x int, y int, KEY I (x));
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change16 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| b
: bytes
| PRIMARY
: int64 0
< row
: bytes
| v
: bytes
| I
: int64 1
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| id
< row
: bytes
| b
: bytes
| x
< row
: bytes
| b
: bytes
| y
< row
: bytes
| v
: bytes
| x
< row
: bytes
| v
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| c
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| c
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| v
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	v	I	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	v	I	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	v	I	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	v	I	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	v	I	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	v	I	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	v	I	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| index	shop	v	I	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		s	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		s	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		s	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		s	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		s	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		s	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		s	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		s	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		s	SEQUENCE
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change16 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE s (x int, y int, KEY x (x), FOREIGN KEY g (y) REFERENCES a (id));
< error
| 1824 HY000 Failed to open the referenced table 'a'
> exec
| CREATE TABLE c (x int, y int, KEY x (x));
< error
| 1050 42S01 Table 'shop.c' already exists
> exec
| TRUNCATE TABLE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| RENAME TABLE c TO c;
< error
| 1050 42S01 Table 'shop.c' already exists
> exec
| ALTER TABLE A DROP FOREIGN KEY f, DROP FOREIGN KEY fk_1;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE s RENAME COLUMN x TO y;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE q.a DROP COLUMN IF EXISTS X, ADD COLUMN y int;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP INDEX i ON v;
< ok 0 0
> exec
| ALTER TABLE A RENAME COLUMN z TO X, ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE a RENAME COLUMN X TO x;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE SEQUENCE s;
< ok 0 0
> exec
| ALTER TABLE c DROP INDEX i;
< error
| 1347 HY000 'shop.c' is not BASE TABLE
> exec
| ALTER TABLE q.a RENAME TO s;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP INDEX I ON q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| CREATE SEQUENCE IF NOT EXISTS s;
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN X int AFTER X;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE b RENAME TO c;
< error
| 1050 42S01 Table 'shop.c' already exists
> exec
| CREATE TABLE IF NOT EXISTS v LIKE b;
< ok 0 0
> exec
| TRUNCATE TABLE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist

== TestWalkThroughSweepTiDB/change17 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| TRUNCATE TABLE q.s;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| CREATE UNIQUE INDEX IF NOT EXISTS I ON s (y);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE q.b ADD COLUMN y int;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN X int, DROP COLUMN y;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE v ADD FOREIGN KEY g (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| INSERT INTO c (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE b CHANGE X x int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| RENAME TABLE A TO a, v TO v;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change17 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| database	test			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| database	test			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| database	test			
| relation	shop		a	SEQUENCE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| database	test			
| relation	shop		a	SEQUENCE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| database	test			
| relation	shop		a	SEQUENCE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| database	test			
| relation	shop		a	SEQUENCE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| database	test			
| relation	shop		a	SEQUENCE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| database	test			
| relation	shop		a	SEQUENCE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| database	test			
| relation	shop		a	SEQUENCE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| database	test			
| relation	shop		a	SEQUENCE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| database	test			
| relation	shop		a	SEQUENCE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| database	test			
| relation	shop		a	SEQUENCE
| relation	shop		c	VIEW

== TestWalkThroughSweepTiDB/change17 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP VIEW v, A;
< error
| 1051 42S02 Unknown table 'shop.v,shop.A'
> exec
| DROP VIEW q.b, c;
< error
| 1051 42S02 Unknown table 'q.b,shop.c'
> exec
| RENAME TABLE q.c TO a;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE s RENAME COLUMN y TO y;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE VIEW c AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE a DROP INDEX i, ADD INDEX j (id);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE SEQUENCE a;
< ok 0 0
> exec
| DROP VIEW IF EXISTS b, A;
< error
| 1347 HY000 'shop.A' is not VIEW
> exec
| ALTER TABLE c RENAME COLUMN z TO z;
< error
| 1347 HY000 'shop.c' is not BASE TABLE
> exec
| INSERT INTO A (x) VALUES (1);
< error
| 1105 HY000 insert into sequence a is not supported now
> exec
| DROP INDEX j ON s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE c DROP FOREIGN KEY fk_1, ADD COLUMN X int;
< error
| 1347 HY000 'shop.c' is not BASE TABLE
> exec
| CREATE TABLE IF NOT EXISTS q.a (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE A DROP COLUMN X;
< error
| 1347 HY000 'shop.A' is not BASE TABLE
> exec
| CREATE INDEX IF NOT EXISTS I ON s (X);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| DROP SEQUENCE c, A;
< error
| 1347 HY000 'shop.c' is not SEQUENCE

== TestWalkThroughSweepTiDB/change18 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| UPDATE v SET x = 1;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE VIEW b AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE TABLE c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE b CHANGE y x int, MODIFY y bigint;
< error
| 1347 HY000 'shop.b' is not BASE TABLE
> exec
| CREATE UNIQUE INDEX j ON c (x);
< ok 0 0
> exec
| DROP VIEW IF EXISTS A, q.s;
< ok 0 0
> exec
| ALTER TABLE s RENAME COLUMN z TO y;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| DROP INDEX I ON s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| DROP TABLE q.b;
< error
| 1051 42S02 Unknown table 'q.b'
> exec
| ALTER TABLE c DROP INDEX i, ADD INDEX j (id);
< error
| 1091 42000 index i doesn't exist
> exec
| ALTER TABLE a RENAME COLUMN x TO z, ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change18 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| c
: bytes
| PRIMARY
: int64 0
< row
: bytes
| c
: bytes
| j
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| x
< row
: bytes
| c
: bytes
| id
< row
: bytes
| c
: bytes
| x
< row
: bytes
| c
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| b
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| VIEW
: bytes
|
< row
: bytes
| c
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	j	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	j	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	j	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	j	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	j	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	i	
| index	shop	a	j	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	i	
| index	shop	a	j	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	i	
| index	shop	a	j	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	z	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	i	
| index	shop	a	j	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	z	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	i	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	z	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	i	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	z	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	i	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	a	y	
| column	shop	a	z	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	z	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	x	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	a	y	
| column	shop	a	z	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	z	
| column	shop	s	x	
| column	shop	v	y	
| column	shop	v	z	
| database	shop			
| database	test			
| index	shop	a	x	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	a	y	
| column	shop	a	z	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	z	
| column	shop	s	x	
| column	shop	v	y	
| column	shop	v	z	
| database	shop			
| database	test			
| index	shop	a	x	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	a	y	
| column	shop	a	z	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	z	
| column	shop	s	x	
| column	shop	v	y	
| column	shop	v	z	
| database	shop			
| database	test			
| index	shop	a	x	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	a	y	
| column	shop	a	z	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	z	
| column	shop	s	x	
| column	shop	v	y	
| column	shop	v	z	
| database	shop			
| database	test			
| index	shop	a	x	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	a	y	
| column	shop	a	z	
| column	shop	b	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	z	
| column	shop	s	x	
| column	shop	v	y	
| column	shop	v	z	
| database	shop			
| database	test			
| index	shop	a	x	
| index	shop	c	PRIMARY	
| index	shop	c	j	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change18 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE q.A DROP COLUMN X;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1105 HY000 insert into view b is not supported now
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| CREATE TABLE q.A (x int, y int, KEY j (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP INDEX x ON q.c;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE a ADD COLUMN X int, DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE c CHANGE z z int;
< error
| 1054 42S22 Unknown column 'z' in 'c'
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE a (x int, y int, KEY j (x));
< ok 0 0
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| TRUNCATE TABLE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE OR REPLACE VIEW s AS SELECT 1 AS x;
< ok 0 0
> exec
| TRUNCATE TABLE a;
< ok 0 0
> exec
| ALTER TABLE A ADD CONSTRAINT i UNIQUE INDEX I (X);
< ok 0 0
> exec
| CREATE TABLE c LIKE a;
< error
| 1050 42S01 Table 'shop.c' already exists
> exec
| DROP VIEW IF EXISTS c, s;
< error
| 1347 HY000 'shop.c' is not VIEW
> exec
| ALTER TABLE c CHANGE y z int;
< ok 0 0
> exec
| ALTER TABLE A ADD CONSTRAINT x UNIQUE INDEX x (y);
< ok 0 0
> exec
| ALTER TABLE v ADD COLUMN y int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE q.c RENAME COLUMN z TO y;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE a ADD COLUMN z int, DROP COLUMN x;
< ok 0 0
> exec
| CREATE TABLE v LIKE a;
< ok 0 0
> exec
| CREATE TABLE q.b (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE VIEW A AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'shop.A' already exists
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| ALTER TABLE q.s CHANGE x x int;
< error
| 1146 42S02 Table 'q.s' doesn't exist

== TestWalkThroughSweepTiDB/change19 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE b CHANGE y x int, CHANGE y X int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE a RENAME COLUMN X TO X;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE v (x int, y int, KEY i (x));
< ok 0 0
> exec
| DROP VIEW IF EXISTS c;
< ok 0 0
> exec
| ALTER TABLE b DROP COLUMN IF EXISTS x;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE INDEX IF NOT EXISTS i ON v (X);
< ok 0 0
> exec
| CREATE TABLE b (x int, y int, KEY I (x));
< ok 0 0
> exec
| CREATE TABLE v AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| ALTER TABLE a ADD INDEX I (X);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change19 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| v
: bytes
| i
: int64 1
< row
: bytes
| b
: bytes
| I
: int64 1
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| x
< row
: bytes
| b
: bytes
| y
< row
: bytes
| v
: bytes
| x
< row
: bytes
| v
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	I	
| index	shop	v	i	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change19 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE q.c ADD COLUMN X int, DROP COLUMN y;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE a DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE s DROP FOREIGN KEY f, DROP FOREIGN KEY fk_1;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE v RENAME TO v;
< ok 0 0
> exec
| ALTER TABLE b DROP FOREIGN KEY F;
< error
| 1091 42000 Can't DROP 'F'; check that column/key exists
> exec
| CREATE TABLE IF NOT EXISTS a LIKE s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE TABLE q.s LIKE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP TEMPORARY TABLE v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| ALTER TABLE b ENGINE = InnoDB;
< ok 0 0
> exec
| CREATE SEQUENCE q.a;
< error
| 1049 42000 Unknown database 'q'
> exec
| UPDATE s SET x = 1;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE A ENGINE = InnoDB;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE v DROP FOREIGN KEY f, DROP FOREIGN KEY g;
< error
| 1091 42000 Can't DROP 'f'; check that column/key exists
> exec
| INSERT INTO A (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE q.c DROP INDEX I;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DROP SEQUENCE IF EXISTS q.A, b;
< ok 0 0
> exec
| CREATE VIEW a AS SELECT 1 AS x;
< ok 0 0
> exec
| DROP TABLE IF EXISTS q.A;
< ok 0 0
> exec
| RENAME TABLE s TO A, q.c TO c;
< error
| 1050 42S01 Table 'shop.A' already exists
> exec
| DROP TEMPORARY TABLE IF EXISTS a;
< ok 0 0
> exec
| USE shop;
< ok 0 0
> exec
| CREATE OR REPLACE VIEW a AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE s ADD UNIQUE INDEX i (x);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE c RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS s LIKE q.v;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| DROP TEMPORARY TABLE c;
< error
| 1051 42S02 Unknown table 'shop.c'
> exec
| DROP INDEX I ON v;
< ok 0 0

== TestWalkThroughSweepTiDB/change2 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| UPDATE a SET x = 1;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE q.s (x int, y int, KEY I (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE q.a LIKE A;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE A ADD INDEX j (X);
< ok 0 0
> exec
| CREATE TEMPORARY TABLE b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE q.c (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE a DROP COLUMN IF EXISTS x;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change2 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| A
: bytes
| PRIMARY
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| A
: bytes
| id
< row
: bytes
| A
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| A
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	A	id	
| column	shop	A	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	A	id	
| column	shop	A	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	b	fk_1	
| index	shop	b	x	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		s	VIEW
| table foreign key	shop	b	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	j	
| index	shop	b	fk_1	
| index	shop	b	x	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		s	VIEW
| table foreign key	shop	b	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	j	
| index	shop	b	fk_1	
| index	shop	b	x	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| table foreign key	shop	b	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	j	
| index	shop	b	fk_1	
| index	shop	b	x	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| table foreign key	shop	b	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	j	
| index	shop	b	fk_1	
| index	shop	b	x	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| table foreign key	shop	b	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	j	
| index	shop	b	fk_1	
| index	shop	b	x	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| table foreign key	shop	b	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	j	
| index	shop	b	fk_1	
| index	shop	b	x	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| table foreign key	shop	b	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	j	
| index	shop	b	fk_1	
| index	shop	b	x	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| table foreign key	shop	b	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	j	
| index	shop	b	fk_1	
| index	shop	b	x	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| table foreign key	shop	b	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	id	
| column	shop	A	y	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	j	
| index	shop	b	fk_1	
| index	shop	b	x	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	SEQUENCE
| table foreign key	shop	b	fk_1	

== TestWalkThroughSweepTiDB/change2 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE VIEW s AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE b ADD CONSTRAINT fk_1 FOREIGN KEY (x) REFERENCES q.s (id);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS q.c LIKE s;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE b DROP INDEX I;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE b (x int, y int, KEY x (x), CONSTRAINT fk_1 FOREIGN KEY g (y) REFERENCES a (id));
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT j UNIQUE INDEX I (y);
< ok 0 0
> exec
| CREATE SEQUENCE c;
< ok 0 0
> exec
| ALTER TABLE q.a DROP FOREIGN KEY fk_1;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP TABLE IF EXISTS c;
< ok 0 0
> exec
| DROP VIEW IF EXISTS s;
< ok 0 0
> exec
| CREATE SEQUENCE a;
< error
| 1050 42S01 Table 'shop.a' already exists
> exec
| RENAME TABLE v TO q.c, A TO b;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE a (x int, y int, KEY I (x), CONSTRAINT fk_1 FOREIGN KEY (y) REFERENCES a (id));
< error
| 3734 HY000 Failed to add the foreign key constraint. Missing column 'id' for constraint 'fk_1' in the referenced table 'a'
> exec
| TRUNCATE TABLE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| USE shop;
< ok 0 0

== TestWalkThroughSweepTiDB/change20 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE q.a ADD COLUMN z int AFTER X;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| TRUNCATE TABLE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE s DROP INDEX j, ADD INDEX x (id);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| DROP SEQUENCE c, s;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.c,shop.s'
> exec
| CREATE INDEX I ON s (z);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| DROP TABLE a, a;
< error
| 1051 42S02 Unknown table 'shop.a,shop.a'
> exec
| ALTER TABLE v CHANGE y x int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE v DROP COLUMN IF EXISTS x;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE c ADD UNIQUE INDEX x (X);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change20 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	i	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	i	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	i	
| relation	shop		b	BASE TABLE

== TestWalkThroughSweepTiDB/change20 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE q.A LIKE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| RENAME TABLE v TO c;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE c DROP INDEX x;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE s DROP FOREIGN KEY fk_1, ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| DROP INDEX x ON q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE q.b ADD COLUMN z int;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| DROP VIEW v, s;
< error
| 1051 42S02 Unknown table 'shop.v,shop.s'
> exec
| CREATE UNIQUE INDEX j ON s (x);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE a RENAME COLUMN x TO x;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE A DROP COLUMN z, ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE b ADD PRIMARY KEY (x);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN y int, DROP COLUMN z;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE a DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE b (x int, y int, KEY i (x));
< ok 0 0
> exec
| CREATE TABLE s LIKE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE v ADD COLUMN x int, DROP COLUMN z;
< error
| 1146 42S02 Table 'shop.v' doesn't exist

== TestWalkThroughSweepTiDB/change21 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE q.c ADD COLUMN z int;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE INDEX I ON c (X);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE a (x int, y int, KEY j (x));
< ok 0 0
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE q.a (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE b ADD COLUMN y int, DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| INSERT INTO q.A (x) VALUES (1);
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| CREATE SEQUENCE IF NOT EXISTS a;
< ok 0 0
> exec
| ALTER TABLE A CHANGE X x int;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change21 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| a
: bytes
| j
: int64 1
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| x
< row
: bytes
| a
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| database	test			
| index	shop	a	j	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| database	test			
| index	shop	a	j	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			

== TestWalkThroughSweepTiDB/change21 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE q.c AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP TABLE q.c, a;
< error
| 1051 42S02 Unknown table 'q.c'

== TestWalkThroughSweepTiDB/change22 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| DROP SEQUENCE v;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.v'
> exec
| DROP TABLE q.s, b;
< error
| 1051 42S02 Unknown table 'q.s,shop.b'
> exec
| INSERT INTO s (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE q.a DROP FOREIGN KEY fk_1, DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE q.b (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE A RENAME COLUMN X TO x;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE q.A ADD FOREIGN KEY fk_1 (x) REFERENCES q.A (id);
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE s ENGINE = InnoDB;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE s ADD COLUMN z int, DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE INDEX x ON c (X);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE s RENAME COLUMN z TO y;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change22 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| v
: bytes
| PRIMARY
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| v
: bytes
| id
< row
: bytes
| v
: bytes
| x
< row
: bytes
| v
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| v
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| relation	shop		c	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change22 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE TABLE b (x int, y int, KEY j (x), CONSTRAINT g FOREIGN KEY g (y) REFERENCES a (id));
< error
| 1824 HY000 Failed to open the referenced table 'a'
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE c RENAME COLUMN y TO z;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE q.b ADD COLUMN z int;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS q.A (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| UPDATE q.A SET x = 1;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE A ENGINE = InnoDB;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE s RENAME COLUMN y TO z, ADD COLUMN y int;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE SEQUENCE IF NOT EXISTS c;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| DROP TABLE s;
< error
| 1051 42S02 Unknown table 'shop.s'
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| DROP SEQUENCE b;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.b'
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS v LIKE s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE A ENGINE = InnoDB;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE OR REPLACE VIEW s AS SELECT 1 AS x;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE b RENAME TO a;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE TABLE v AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE A DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE DATABASE q;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| DROP TABLE s;
< error
| 1051 42S02 Unknown table 'shop.s'
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE TABLE s (x int, y int, KEY i (x));
< error
| 1050 42S01 Table 'shop.s' already exists
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE TABLE q.b LIKE c;
< error
| 1347 HY000 'q.c' is not BASE TABLE
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE INDEX IF NOT EXISTS I ON q.v (z);
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE A ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| RENAME TABLE v TO A, q.a TO q.s;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist

== TestWalkThroughSweepTiDB/change23 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| DROP SEQUENCE a;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.a'
> exec
| RENAME TABLE b TO a;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| UPDATE b SET x = 1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE VIEW b AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE OR REPLACE VIEW s AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE s ADD COLUMN y int;
< error
| 1347 HY000 'shop.s' is not BASE TABLE
> exec
| CREATE TABLE v (x int, y int, KEY x (x));
< ok 0 0
> exec
| DROP VIEW IF EXISTS v;
< error
| 1347 HY000 'shop.v' is not VIEW
> exec
| CREATE DATABASE q;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change23 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| v
: bytes
| x
: int64 1
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| x
< row
: bytes
| s
: bytes
| x
< row
: bytes
| v
: bytes
| x
< row
: bytes
| v
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| b
< row
: bytes
| s
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| VIEW
: bytes
|
< row
: bytes
| s
: bytes
| VIEW
: bytes
|
< row
: bytes
| v
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| column	shop	v	z	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| column	shop	v	z	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| column	shop	v	z	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| column	shop	v	z	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| column	shop	v	z	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| column	shop	v	z	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	b	x	
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| column	shop	v	z	
| database	q			
| database	shop			
| database	test			
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change23 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `s` (`x`) AS SELECT 1 AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| ALTER TABLE s ADD INDEX j (z);
< error
| 1347 HY000 'shop.s' is not BASE TABLE
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `s` (`x`) AS SELECT 1 AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS s (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| TRUNCATE TABLE s;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| ALTER TABLE A ADD INDEX i (x);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| TRUNCATE TABLE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| DROP TABLE IF EXISTS a;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| ALTER TABLE c DROP INDEX i;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| ALTER TABLE v ADD COLUMN X int;
< error
| 1060 42S21 Duplicate column name 'X'
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| ALTER TABLE c ADD UNIQUE INDEX i (z);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| CREATE SEQUENCE A;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| ALTER TABLE q.s DROP COLUMN IF EXISTS y;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| DROP TABLE b, q.b;
< error
| 1051 42S02 Unknown table 'shop.b,q.b'
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| ALTER TABLE v ADD COLUMN z int AFTER x;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| ALTER TABLE q.v RENAME TO c;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| ALTER TABLE a ADD COLUMN x int, DROP COLUMN z;
< error
| 1347 HY000 'shop.a' is not BASE TABLE
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| ALTER TABLE A ADD COLUMN y int, DROP COLUMN y;
< error
| 1347 HY000 'shop.A' is not BASE TABLE
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| TRUNCATE TABLE q.c;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec
| CREATE TABLE v (x int, y int, KEY j (x));
< error
| 1050 42S01 Table 'shop.v' already exists
> query
| SHOW CREATE TABLE `q`.`s`
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE TEMPORARY TABLE `s` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`s`
< error
| 1146 42S02 Table 'test.s' doesn't exist

== TestWalkThroughSweepTiDB/change24 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| USE shop;
< ok 0 0
> exec
| ALTER TABLE A CHANGE y z int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE TABLE q.A LIKE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP INDEX I ON c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE v ADD CONSTRAINT F FOREIGN KEY (x) REFERENCES s (id);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE A DROP COLUMN IF EXISTS y, ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| DROP TEMPORARY TABLE b, a;
< error
| 1051 42S02 Unknown table 'shop.b,shop.a'
> exec
| INSERT INTO q.v (x) VALUES (1);
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| TRUNCATE TABLE s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE s DROP INDEX i, ADD INDEX I (id);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN z TO x;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP SEQUENCE a;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.a'
> exec
| CREATE UNIQUE INDEX IF NOT EXISTS j ON s (x);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change24 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE

== TestWalkThroughSweepTiDB/change24 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| INSERT INTO a (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| CREATE TABLE q.a AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| INSERT INTO q.c (x) VALUES (1);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE q.c ADD COLUMN x int, DROP COLUMN z;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE SEQUENCE b;
< ok 0 0
> exec
| ALTER TABLE v DROP COLUMN IF EXISTS x;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE q.A ADD COLUMN z int;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| RENAME TABLE q.A TO q.s;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 0 0
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| TRUNCATE TABLE c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE v RENAME COLUMN y TO X;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| TRUNCATE TABLE q.s;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| UPDATE q.A SET x = 1;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN X int;
< error
| 1347 HY000 'shop.b' is not BASE TABLE
> exec
| ALTER TABLE A RENAME COLUMN z TO X, ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE q.A DROP INDEX j;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE c RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.c' doesn't exist

== TestWalkThroughSweepTiDB/change25 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE b ADD COLUMN x int, DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE SEQUENCE v;
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS v;
< ok 0 0
> exec
| CREATE INDEX I ON q.A (x);
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE q.a DROP COLUMN y;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE VIEW c AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE s ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE v ENGINE = InnoDB;
< error
| 1347 HY000 'shop.v' is not BASE TABLE
> exec
| DROP SEQUENCE q.b;
< error
| 4139 HY000 Unknown SEQUENCE: 'q.b'
> exec
| ALTER TABLE A RENAME COLUMN X TO z, ADD COLUMN y int;
< error
| 1060 42S21 Duplicate column name 'y'
> exec
| DROP VIEW IF EXISTS q.b;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change25 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| A
: bytes
| PRIMARY
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| A
: bytes
| id
< row
: bytes
| A
: bytes
| x
< row
: bytes
| A
: bytes
| y
< row
: bytes
| c
: bytes
| x
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| v
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| A
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| c
: bytes
| VIEW
: bytes
|
< row
: bytes
| v
: bytes
| SEQUENCE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	i	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	i	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	i	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	i	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	i	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	i	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	i	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	i	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	i	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	i	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	i	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	i	
| index	shop	b	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	VIEW
| relation	shop		v	SEQUENCE

== TestWalkThroughSweepTiDB/change25 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE s DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| RENAME TABLE q.A TO q.b;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| DROP TABLE b;
< error
| 1051 42S02 Unknown table 'shop.b'
> exec
| TRUNCATE TABLE q.A;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE v DROP INDEX j;
< error
| 1347 HY000 'shop.v' is not BASE TABLE
> exec
| CREATE INDEX I ON c (x);
< error
| 1105 HY000 Unknown charset 
> exec
| ALTER TABLE s ADD COLUMN x int, DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE c CHANGE z X int, CHANGE x x int;
< error
| 1347 HY000 'shop.c' is not BASE TABLE
> exec
| CREATE TABLE q.c (x int, y int, KEY x (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE q.b (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE q.A DROP COLUMN y;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE UNIQUE INDEX i ON A (y);
< ok 0 0
> exec
| DROP INDEX i ON q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE q.v RENAME COLUMN x TO y;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| DROP INDEX i ON c;
< error
| 1091 42000 index i doesn't exist
> exec
| CREATE TABLE q.s LIKE A;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE v DROP COLUMN IF EXISTS x, ADD COLUMN y int;
< error
| 1347 HY000 'shop.v' is not BASE TABLE
> exec
| ALTER TABLE q.b ADD COLUMN z int, DROP COLUMN z;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| CREATE TABLE A AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP SEQUENCE b;
< error
| 1347 HY000 'shop.b' is not SEQUENCE
> exec
| UPDATE A SET x = 1;
< ok 0 0
> exec
| CREATE TABLE q.c LIKE b;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE c (x int, y int, KEY i (x));
< error
| 1050 42S01 Table 'shop.c' already exists

== TestWalkThroughSweepTiDB/change26 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| RENAME TABLE q.s TO s, s TO s;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| ALTER TABLE c ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE VIEW s AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE s RENAME COLUMN y TO z;
< error
| 1347 HY000 'shop.s' is not BASE TABLE
> exec
| CREATE TABLE b LIKE s;
< error
| 1347 HY000 'shop.s' is not BASE TABLE
> exec
| CREATE TABLE a (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE a (x int, y int, KEY j (x));
< error
| 1050 42S01 Table 'shop.a' already exists
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS s (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE v RENAME COLUMN z TO y;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE v CHANGE y x int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE A (x int, y int, KEY j (x));
< error
| 1050 42S01 Table 'shop.A' already exists
> exec
| ALTER TABLE c RENAME COLUMN z TO X;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change26 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| a
: bytes
| PRIMARY
: int64 0
< row
: bytes
| a
: bytes
| x
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| id
< row
: bytes
| a
: bytes
| x
< row
: bytes
| a
: bytes
| y
< row
: bytes
| s
: bytes
| x
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| s
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| s
: bytes
| VIEW
: bytes
|
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW

== TestWalkThroughSweepTiDB/change26 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| TRUNCATE TABLE q.s;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| CREATE TABLE A (x int, y int, KEY j (x));
< error
| 1050 42S01 Table 'shop.A' already exists
> exec
| CREATE TABLE IF NOT EXISTS q.v (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP INDEX j ON q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE v LIKE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE b RENAME COLUMN X TO x;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE s DROP INDEX I;
< error
| 1347 HY000 'shop.s' is not BASE TABLE
> exec
| DROP INDEX j ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP VIEW IF EXISTS a, q.b;
< error
| 1347 HY000 'shop.a' is not VIEW
> exec
| CREATE INDEX i ON q.v (x);
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| UPDATE q.s SET x = 1;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE q.b DROP INDEX j;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| DROP SEQUENCE a;
< error
| 1347 HY000 'shop.a' is not SEQUENCE

== TestWalkThroughSweepTiDB/change27 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| DROP INDEX j ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE c DROP INDEX x, ADD INDEX j (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE q.s DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| CREATE TABLE a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE A RENAME COLUMN y TO y;
< ok 0 0
> exec
| CREATE UNIQUE INDEX I ON q.b (y);
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE q.a CHANGE x x int;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE q.s ADD COLUMN y int;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| ALTER TABLE v CHANGE X X int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| INSERT INTO a (x) VALUES (1);
< error
| 1364 HY000 Field 'id' doesn't have a default value
> exec
| DROP INDEX I ON q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change27 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| a
: bytes
| PRIMARY
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| id
< row
: bytes
| a
: bytes
| x
< row
: bytes
| a
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		v	VIEW

== TestWalkThroughSweepTiDB/change27 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| RENAME TABLE a TO q.s;
< error
| 1025 HY000 Error on rename of 'shop.a' to 'q.s' (errno: 168 - Database `q` doesn't exist)
> exec
| TRUNCATE TABLE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE VIEW v AS SELECT 1 AS x;
< ok 0 0
> exec
| DROP VIEW IF EXISTS q.s;
< ok 0 0
> exec
| ALTER TABLE A CHANGE X x int;
< ok 0 0
> exec
| TRUNCATE TABLE s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE c ADD CONSTRAINT F FOREIGN KEY (x) REFERENCES b (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP INDEX x ON q.s;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| UPDATE v SET x = 1;
< error
| 1288 HY000 The target table v of the UPDATE is not updatable
> exec
| CREATE TABLE b LIKE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE OR REPLACE VIEW b AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE TABLE v (id int PRIMARY KEY, x int, y int);
< error
| 1050 42S01 Table 'shop.v' already exists
> exec
| CREATE TABLE A (x int, y int, KEY x (x));
< error
| 1050 42S01 Table 'shop.A' already exists
> exec
| ALTER TABLE q.b ADD COLUMN x int, DROP COLUMN z;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1105 HY000 insert into view v is not supported now
> exec
| ALTER TABLE s DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| RENAME TABLE c TO a;
< error
| 1050 42S01 Table 'shop.a' already exists
> exec
| DROP INDEX x ON s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| INSERT INTO c (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.c' doesn't exist

== TestWalkThroughSweepTiDB/change28 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE A ADD FOREIGN KEY g (x) REFERENCES q.s (id);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| DROP TEMPORARY TABLE v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| DROP SEQUENCE s, c;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.s,shop.c'
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP VIEW a;
< error
| 1347 HY000 'shop.a' is not VIEW
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change28 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| a
: bytes
| PRIMARY
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| id
< row
: bytes
| a
: bytes
| x
< row
: bytes
| a
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change28 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE c ADD FOREIGN KEY F (x) REFERENCES q.a (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE INDEX j ON A (z);
< error
| 1072 42000 column does not exist: z
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| UPDATE q.a SET x = 1;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE s (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE a DROP INDEX i;
< error
| 1091 42000 index i doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE a RENAME COLUMN z TO x;
< error
| 1054 42S22 Unknown column 'z' in 'a'
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE s AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE c DROP COLUMN IF EXISTS z;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE q.A ADD INDEX x (y);
< error
| 1146 42S02 Table 'q.A' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| DROP TABLE q.s;
< error
| 1051 42S02 Unknown table 'q.s'
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE q.v DROP FOREIGN KEY f, ADD COLUMN X int;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TEMPORARY TABLE q.c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| DROP TEMPORARY TABLE q.a;
< error
| 1051 42S02 Unknown table 'q.a'
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| INSERT INTO A (x) VALUES (1);
< error
| 1364 HY000 Field 'id' doesn't have a default value
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE q.s ADD COLUMN z int, DROP COLUMN X;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE b DROP INDEX x, ADD INDEX j (id);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE q.A ADD PRIMARY KEY (x);
< error
| 1146 42S02 Table 'q.A' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| SET foreign_key_checks = 1;
< ok 0 0

== TestWalkThroughSweepTiDB/change29 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE s (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS v (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| ALTER TABLE q.b ADD COLUMN y int, DROP COLUMN x;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| RENAME TABLE q.A TO q.b;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| DROP TABLE IF EXISTS v;
< ok 0 0
> exec
| CREATE UNIQUE INDEX I ON c (y);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE UNIQUE INDEX i ON b (x);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS c LIKE q.s;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change29 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| s
: bytes
| PRIMARY
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| s
: bytes
| id
< row
: bytes
| s
: bytes
| x
< row
: bytes
| s
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| s
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	PRIMARY	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	s	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	s	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	b	j	
| index	shop	c	x	
| index	shop	s	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change29 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| RENAME TABLE q.A TO a;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE s COMMENT 'x', RENAME TO c;
< error
| 8200 HY000 Unsupported multi schema change for rename table
> exec
| ALTER TABLE q.a DROP FOREIGN KEY F, ADD COLUMN X int;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| UPDATE a SET x = 1;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE b RENAME COLUMN x TO y;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE b CHANGE z x int, RENAME COLUMN x TO X;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE c DROP FOREIGN KEY fk_1, DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE b (x int, y int, KEY j (x));
< ok 0 0
> exec
| UPDATE q.a SET x = 1;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE c (x int, y int, KEY x (x));
< ok 0 0
> exec
| ALTER TABLE q.v RENAME COLUMN X TO X;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| CREATE TABLE A LIKE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE a ADD PRIMARY KEY (x);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE v RENAME COLUMN X TO z, ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| UPDATE c SET x = 1;
< ok 0 0
> exec
| ALTER TABLE A CHANGE y z int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN y, ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE q.A AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| CREATE INDEX j ON q.c (y);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE VIEW v AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'shop.v' already exists
> exec
| ALTER TABLE a ADD COLUMN X int, DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int, y int);
< error
| 1050 42S01 Table 'shop.b' already exists
> exec
| ALTER TABLE v RENAME COLUMN X TO X, ADD COLUMN y int;
< error
| 1060 42S21 Duplicate column name 'y'

== TestWalkThroughSweepTiDB/change3 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| DROP SEQUENCE A;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.A'
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS c (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| ALTER TABLE q.A DROP FOREIGN KEY F, ADD COLUMN y int;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| RENAME TABLE v TO a, b TO A;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE A ADD CONSTRAINT i UNIQUE INDEX I (y);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP INDEX I ON a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP INDEX x ON q.A;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS q.b (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE INDEX j ON a (X);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE q.a (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE SEQUENCE A;
< ok 0 0
> exec
| ALTER TABLE v COMMENT 'x', RENAME TO a;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change3 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| A
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| A
: bytes
| SEQUENCE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		A	SEQUENCE
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		b	SEQUENCE
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change3 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| RENAME TABLE q.c TO a;
< error
| 1050 42S01 Table 'shop.a' already exists
> exec
| ALTER TABLE q.c RENAME COLUMN X TO X;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE q.v RENAME COLUMN x TO x;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| CREATE TABLE v (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP VIEW IF EXISTS q.v;
< ok 0 0
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE SEQUENCE b;
< ok 0 0
> exec ^
< error
| 1050 42S01 Table 'shop.b' already exists
> exec
| DROP SEQUENCE c;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.c'
> exec
| USE shop;
< ok 0 0
> exec
| DROP TABLE IF EXISTS a;
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS q.s, A;
< ok 0 0
> exec
| CREATE INDEX I ON b (y);
< error
| 1072 42000 column does not exist: y
> exec
| CREATE TABLE q.A (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'

== TestWalkThroughSweepTiDB/change30 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| DROP SEQUENCE v;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.v'
> exec
| USE shop;
< ok 0 0
> exec
| DROP VIEW IF EXISTS q.c;
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a;
< ok 0 0
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE q.c (x int, y int, KEY i (x), CONSTRAINT fk_1 FOREIGN KEY g (y) REFERENCES a (id));
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TEMPORARY TABLE q.a (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE a ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP INDEX x ON q.v;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| DROP SEQUENCE a;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.a'
> exec
| ALTER TABLE v ADD COLUMN y int, DROP COLUMN y;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP VIEW IF EXISTS q.a;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change30 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| b
: bytes
| PRIMARY
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| id
< row
: bytes
| b
: bytes
| x
< row
: bytes
| b
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE

== TestWalkThroughSweepTiDB/change30 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE a DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| INSERT INTO q.a (x) VALUES (1);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE s DROP INDEX x;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| DROP VIEW q.b, c;
< error
| 1051 42S02 Unknown table 'q.b,shop.c'
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| RENAME TABLE a TO s;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE TEMPORARY TABLE q.a (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE TEMPORARY TABLE A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE a ADD COLUMN z int, DROP COLUMN x;
< error
| 8200 HY000 TiDB doesn't support ALTER TABLE for local temporary table
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| USE shop;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE v DROP COLUMN IF EXISTS X;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE q.a RENAME COLUMN z TO y, ADD COLUMN X int;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS c LIKE q.A;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE a ADD CONSTRAINT g FOREIGN KEY (x) REFERENCES s (id);
< error
| 8200 HY000 TiDB doesn't support ALTER TABLE for local temporary table
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| DROP TABLE q.s, v;
< error
| 1051 42S02 Unknown table 'q.s,shop.v'
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN y int, DROP COLUMN X;
< error
| 1060 42S21 Duplicate column name 'y'
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| RENAME TABLE c TO a, b TO v;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| INSERT INTO c (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist

== TestWalkThroughSweepTiDB/change31 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| UPDATE s SET x = 1;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| TRUNCATE TABLE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE b DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE q.s CHANGE X x int;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| ALTER TABLE s DROP FOREIGN KEY g;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS a (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN X TO x, ADD COLUMN z int;
< error
| 8200 HY000 TiDB doesn't support ALTER TABLE for local temporary table
> exec
| ALTER TABLE q.s DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change31 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		A	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		A	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		A	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		A	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		A	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		A	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	I	
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change31 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE SEQUENCE q.a;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE v (x int, y int, KEY I (x));
< ok 0 0
> exec
| ALTER TABLE b ADD CONSTRAINT fk_1 FOREIGN KEY F (x) REFERENCES s (id);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| UPDATE s SET x = 1;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE q.b DROP FOREIGN KEY fk_1;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| TRUNCATE TABLE v;
< ok 0 0
> exec
| DROP VIEW q.A, A;
< error
| 1051 42S02 Unknown table 'q.A,shop.A'
> exec
| ALTER TABLE q.b ADD COLUMN X int AFTER X;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE A CHANGE X z int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE s CHANGE z z int, MODIFY X bigint;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE c DROP INDEX i;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| RENAME TABLE a TO v, A TO a;
< error
| 1050 42S01 Table 'shop.v' already exists
> exec
| USE shop;
< ok 0 0
> exec
| DROP SEQUENCE q.s;
< error
| 4139 HY000 Unknown SEQUENCE: 'q.s'
> exec
| CREATE SEQUENCE IF NOT EXISTS A;
< ok 0 0
> exec
| ALTER TABLE A DROP INDEX x, ADD INDEX I (id);
< error
| 1347 HY000 'shop.A' is not BASE TABLE
> exec
| TRUNCATE TABLE s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE INDEX I ON b (z);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| INSERT INTO q.b (x) VALUES (1);
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE VIEW b AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE s ADD CONSTRAINT fk_1 FOREIGN KEY F (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'shop.s' doesn't exist

== TestWalkThroughSweepTiDB/change32 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE b ADD FOREIGN KEY f (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE s DROP FOREIGN KEY g;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE q.c DROP INDEX I;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE c DROP INDEX j, ADD INDEX j (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE v RENAME COLUMN X TO X;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE s RENAME COLUMN x TO z;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE SEQUENCE IF NOT EXISTS v;
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 0 0
> exec
| RENAME TABLE a TO q.b, b TO s;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN x int, DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| UPDATE A SET x = 1;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE A ADD COLUMN X int, DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE TABLE q.A (x int, y int, KEY x (x), FOREIGN KEY f (y) REFERENCES a (id));
< error
| 3734 HY000 Failed to add the foreign key constraint. Missing column 'id' for constraint 'f' in the referenced table 'a'
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1105 HY000 insert into sequence v is not supported now
> exec
| TRUNCATE TABLE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change32 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| v
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| v
: bytes
| SEQUENCE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		b	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		b	VIEW
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	x	
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		A	VIEW
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	x	
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		A	VIEW
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	x	
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		A	VIEW
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	x	
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		A	VIEW
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	x	
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		A	VIEW
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	x	
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		A	VIEW
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	x	
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		A	VIEW
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	x	
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		A	VIEW
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	x	
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		A	VIEW
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	x	
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		A	VIEW
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	A	x	
| column	shop	b	x	
| database	q			
| database	shop			
| database	test			
| relation	shop		A	VIEW
| relation	shop		b	VIEW
| relation	shop		c	SEQUENCE
| relation	shop		v	SEQUENCE

== TestWalkThroughSweepTiDB/change32 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE b DROP INDEX x, ADD INDEX i (id);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| INSERT INTO c (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP SEQUENCE b, a;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.b,shop.a'
> exec
| ALTER TABLE q.b DROP COLUMN IF EXISTS x;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| CREATE VIEW b AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE TABLE b AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| CREATE SEQUENCE IF NOT EXISTS c;
< ok 0 0
> exec
| ALTER TABLE q.s DROP INDEX x, ADD INDEX i (id);
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| CREATE VIEW A AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE c DROP FOREIGN KEY fk_1, DROP FOREIGN KEY F;
< error
| 1347 HY000 'shop.c' is not BASE TABLE
> exec
| CREATE TABLE s AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| ALTER TABLE q.v ADD INDEX j (x);
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| USE q;
< ok 0 0
> exec
| CREATE DATABASE q;
< error
| 1007 HY000 Can't create database 'q'; database exists
> exec
| RENAME TABLE A TO b;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| TRUNCATE TABLE q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE q.c DROP FOREIGN KEY fk_1;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE UNIQUE INDEX j ON A (z);
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| DROP INDEX i ON v;
< error
| 1146 42S02 Table 'q.v' doesn't exist

== TestWalkThroughSweepTiDB/change33 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE s (x int, y int, KEY i (x));
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS s;
< ok 0 0
> exec
| ALTER TABLE v CHANGE y y int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE v CHANGE z y int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE q.v CHANGE y X int;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| CREATE TEMPORARY TABLE a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| UPDATE a SET x = 1;
< ok 0 0
> exec
| DROP INDEX I ON a;
< error
| 8200 HY000 TiDB doesn't support DROP INDEX for local temporary table
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP SEQUENCE A;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.A'
> exec
| ALTER TABLE a ADD CONSTRAINT j UNIQUE INDEX x (y);
< error
| 8200 HY000 TiDB doesn't support ALTER TABLE for local temporary table
> exec
| DROP VIEW IF EXISTS b, c;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change33 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| s
: bytes
| i
: int64 1
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| s
: bytes
| x
< row
: bytes
| s
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| s
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	s	i	
| relation	shop		s	BASE TABLE

== TestWalkThroughSweepTiDB/change33 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE s RENAME COLUMN z TO X, ADD COLUMN X int;
< error
| 1054 42S22 Unknown column 'z' in 's'
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| DROP TABLE IF EXISTS A;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE q.a ADD COLUMN X int, DROP COLUMN x;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE TEMPORARY TABLE q.a (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE q.s ADD CONSTRAINT fk_1 FOREIGN KEY (x) REFERENCES q.a (id);
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| CREATE TABLE q.a AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE q.a RENAME COLUMN y TO x;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE A DROP COLUMN IF EXISTS X, ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE v DROP INDEX i, ADD INDEX x (id);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| RENAME TABLE A TO a, s TO v;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE c DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| TRUNCATE TABLE c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE b DROP COLUMN X, ADD COLUMN y int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| INSERT INTO q.v (x) VALUES (1);
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| ALTER TABLE q.v ADD COLUMN z int AFTER z;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist

== TestWalkThroughSweepTiDB/change34 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE s (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE TABLE A AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| ALTER TABLE b ADD COLUMN z int AFTER X;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE A (x int, y int, KEY j (x));
< ok 0 0
> exec
| ALTER TABLE v RENAME COLUMN X TO z;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE c ADD CONSTRAINT f FOREIGN KEY g (x) REFERENCES s (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE OR REPLACE VIEW s AS SELECT 1 AS x;
< error
| 1347 HY000 'shop.s' is not VIEW
> exec
| CREATE VIEW A AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'shop.A' already exists
> exec
| DROP VIEW A, A;
< error
| 1347 HY000 'shop.A' is not VIEW
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change34 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| s
: bytes
| PRIMARY
: int64 0
< row
: bytes
| s
: bytes
| x
: int64 0
< row
: bytes
| A
: bytes
| j
: int64 1
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| A
: bytes
| x
< row
: bytes
| A
: bytes
| y
< row
: bytes
| s
: bytes
| id
< row
: bytes
| s
: bytes
| x
< row
: bytes
| s
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| A
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| s
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	s	id	
| column	shop	s	x	
| column	shop	s	y	
| database	shop			
| database	test			
| index	shop	A	j	
| index	shop	b	PRIMARY	
| index	shop	s	PRIMARY	
| index	shop	s	x	
| relation	shop		A	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		s	BASE TABLE
| relation	shop		v	SEQUENCE
| table foreign key	shop	s	f	

== TestWalkThroughSweepTiDB/change34 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE q.a AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| CREATE SEQUENCE s;
< error
| 1050 42S01 Table 'shop.s' already exists
> exec
| ALTER TABLE b DROP INDEX i;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE c DROP INDEX x;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE q.c ADD COLUMN X int, DROP COLUMN z;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE SEQUENCE v;
< ok 0 0
> exec
| ALTER TABLE s ADD CONSTRAINT f FOREIGN KEY F (x) REFERENCES s (id);
< ok 0 0
> exec
| TRUNCATE TABLE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| ALTER TABLE q.c RENAME COLUMN z TO x;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| TRUNCATE TABLE a;
< ok 0 0
> exec
| USE shop;
< ok 0 0
> exec
| CREATE TABLE s (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1050 42S01 Table 'shop.s' already exists
> exec
| TRUNCATE TABLE A;
< ok 0 0
> exec
| ALTER TABLE A ADD CONSTRAINT i UNIQUE INDEX j (z);
< error
| 1072 42000 column does not exist: z
> exec
| CREATE UNIQUE INDEX IF NOT EXISTS j ON q.a (X);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP VIEW s;
< error
| 1347 HY000 'shop.s' is not VIEW
> exec
| ALTER TABLE q.v DROP INDEX I;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS c LIKE q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE a RENAME INDEX i TO j;
< error
| 1176 42000 Key 'i' doesn't exist in table 'A'
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| UPDATE b SET x = 1;
< ok 0 0

== TestWalkThroughSweepTiDB/change35 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| UPDATE b SET x = 1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE q.s ADD CONSTRAINT F FOREIGN KEY fk_1 (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| CREATE INDEX i ON a (X);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE s DROP COLUMN y;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE OR REPLACE VIEW s AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN y, ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE a DROP FOREIGN KEY g, ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE SEQUENCE b;
< ok 0 0
> exec
| CREATE TABLE v (x int, y int, KEY i (x));
< ok 0 0
> exec
| ALTER TABLE v RENAME COLUMN X TO y;
< error
| 1060 42S21 Duplicate column name 'y'
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change35 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| v
: bytes
| i
: int64 1
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| s
: bytes
| x
< row
: bytes
| v
: bytes
| x
< row
: bytes
| v
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| s
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| b
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| s
: bytes
| VIEW
: bytes
|
< row
: bytes
| v
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	s	x	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	s	x	
| database	shop			
| database	test			
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW

== TestWalkThroughSweepTiDB/change35 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| TRUNCATE TABLE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE c ADD UNIQUE INDEX j (x);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE VIEW q.c AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE s (x int, y int, KEY I (x));
< error
| 1050 42S01 Table 'shop.s' already exists
> exec
| CREATE UNIQUE INDEX IF NOT EXISTS x ON b (z);
< error
| 1072 42000 column does not exist: z
> exec
| ALTER TABLE a DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE b DROP INDEX I;
< error
| 1347 HY000 'shop.b' is not BASE TABLE
> exec
| ALTER TABLE s ADD CONSTRAINT j UNIQUE INDEX j (x);
< error
| 1347 HY000 'shop.s' is not BASE TABLE
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1050 42S01 Table 'shop.b' already exists
> exec
| DROP TABLE IF EXISTS c;
< ok 0 0
> exec
| ALTER TABLE s DROP FOREIGN KEY f;
< error
| 1347 HY000 'shop.s' is not BASE TABLE
> exec
| ALTER TABLE s ADD CONSTRAINT I UNIQUE INDEX x (X);
< error
| 1347 HY000 'shop.s' is not BASE TABLE
> exec
| ALTER TABLE A DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE c ADD COLUMN X int, DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE v ADD COLUMN x int;
< error
| 1060 42S21 Duplicate column name 'x'
> exec
| TRUNCATE TABLE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP INDEX j ON A;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| DROP TABLE q.c, v;
< error
| 1051 42S02 Unknown table 'q.c'

== TestWalkThroughSweepTiDB/change36 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS c;
< ok 0 0
> exec
| ALTER TABLE b DROP INDEX i, ADD INDEX x (id);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP VIEW IF EXISTS s;
< ok 0 0
> exec
| ALTER TABLE A DROP FOREIGN KEY fk_1;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE q.c (x int, y int, KEY i (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE s ADD COLUMN x int, DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change36 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| c
: bytes
| PRIMARY
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| c
: bytes
| id
< row
: bytes
| c
: bytes
| x
< row
: bytes
| c
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| c
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	c	PRIMARY	
| relation	shop		a	SEQUENCE
| relation	shop		c	BASE TABLE

== TestWalkThroughSweepTiDB/change36 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE SEQUENCE a;
< ok 0 0
> exec
| DROP VIEW q.b, c;
< error
| 1347 HY000 'shop.c' is not VIEW
> exec
| CREATE TABLE q.A (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE q.v LIKE a;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE v RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec ^
< ok 0 0
> exec
| CREATE TABLE v (x int, y int, KEY j (x), CONSTRAINT F FOREIGN KEY (y) REFERENCES a (id));
< error
| 3734 HY000 Failed to add the foreign key constraint. Missing column 'id' for constraint 'F' in the referenced table 'a'
> exec
| CREATE OR REPLACE VIEW q.s AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE b ADD COLUMN z int AFTER z;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE q.A ADD FOREIGN KEY F (x) REFERENCES q.c (id);
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| CREATE SEQUENCE q.c;
< error
| 1049 42000 Unknown database 'q'
> exec
| TRUNCATE TABLE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE a ADD COLUMN x int;
< error
| 1347 HY000 'shop.a' is not BASE TABLE
> exec
| ALTER TABLE v ADD INDEX I (z);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE q.s ADD COLUMN X int, DROP COLUMN y;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| DROP INDEX i ON c;
< error
| 1091 42000 index i doesn't exist

== TestWalkThroughSweepTiDB/change37 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE v DROP INDEX I;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE A DROP FOREIGN KEY fk_1;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS q.s;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE UNIQUE INDEX x ON a (x);
< error
| 8200 HY000 TiDB doesn't support CREATE INDEX for local temporary table
> exec
| ALTER TABLE c ADD CONSTRAINT I UNIQUE INDEX i (x);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE b AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| CREATE TABLE q.s AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP VIEW IF EXISTS v, A;
< ok 0 0
> exec
| TRUNCATE TABLE a;
< ok 0 0
> exec
| UPDATE b SET x = 1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE b (x int, y int, KEY j (x));
< ok 0 0
> exec
| CREATE SEQUENCE q.v;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE v ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change37 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| b
: bytes
| j
: int64 1
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| x
< row
: bytes
| b
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| database	test			
| index	shop	b	j	
| relation	shop		b	BASE TABLE

== TestWalkThroughSweepTiDB/change37 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE q.b AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP INDEX j ON s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE A DROP FOREIGN KEY fk_1;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| DROP TABLE q.a;
< error
| 1051 42S02 Unknown table 'q.a'
> exec
| CREATE TABLE q.a LIKE c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP TABLE a, A;
< error
| 1051 42S02 Unknown table 'shop.a,shop.A'
> exec
| CREATE TABLE q.b AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP VIEW IF EXISTS s;
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS b, s;
< ok 0 0
> exec
| ALTER TABLE s DROP INDEX i;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE TABLE v LIKE q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| TRUNCATE TABLE A;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| DROP VIEW q.c, v;
< error
| 1051 42S02 Unknown table 'q.c,shop.v'
> exec
| DROP TEMPORARY TABLE s, q.v;
< error
| 1051 42S02 Unknown table 'shop.s,q.v'
> exec
| DROP VIEW q.A;
< error
| 1051 42S02 Unknown table 'q.A'
> exec
| CREATE TABLE s (x int, y int, KEY j (x), CONSTRAINT fk_1 FOREIGN KEY f (y) REFERENCES a (id));
< error
| 1824 HY000 Failed to open the referenced table 'a'

== TestWalkThroughSweepTiDB/change38 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE a RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE q.s LIKE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE A CHANGE z z int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| DROP VIEW q.c, a;
< error
| 1051 42S02 Unknown table 'q.c,shop.a'
> exec
| ALTER TABLE b DROP INDEX j;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP INDEX I ON s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE TABLE A LIKE c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE A (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE TABLE c AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP SEQUENCE A;
< error
| 1347 HY000 'shop.A' is not SEQUENCE
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change38 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| c
: bytes
| PRIMARY
: int64 0
< row
: bytes
| A
: bytes
| PRIMARY
: int64 0
< row
: bytes
| A
: bytes
| x
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| A
: bytes
| id
< row
: bytes
| A
: bytes
| x
< row
: bytes
| A
: bytes
| y
< row
: bytes
| c
: bytes
| id
< row
: bytes
| c
: bytes
| x
< row
: bytes
| c
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| A
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| c
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	A	id	
| column	shop	A	x	
| column	shop	A	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| database	test			
| index	shop	A	PRIMARY	
| index	shop	A	x	
| index	shop	c	PRIMARY	
| relation	shop		A	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW

== TestWalkThroughSweepTiDB/change38 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP TEMPORARY TABLE IF EXISTS a;
< ok 0 0
> exec
| RENAME TABLE A TO c;
< error
| 1050 42S01 Table 'shop.c' already exists
> exec
| CREATE SEQUENCE A;
< error
| 1050 42S01 Table 'shop.A' already exists
> exec
| CREATE SEQUENCE q.A;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE q.c ADD CONSTRAINT i UNIQUE INDEX x (X);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE s RENAME COLUMN y TO X;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE UNIQUE INDEX IF NOT EXISTS I ON c (z);
< error
| 1072 42000 column does not exist: z
> exec
| DROP SEQUENCE v;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.v'
> exec
| ALTER TABLE q.c CHANGE y z int, CHANGE x y int;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DROP SEQUENCE IF EXISTS c, q.A;
< ok 0 0
> exec
| DROP VIEW A, s;
< error
| 1347 HY000 'shop.A' is not VIEW
> exec
| CREATE VIEW v AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE TABLE q.s (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE s ADD UNIQUE INDEX i (x);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| UPDATE c SET x = 1;
< ok 0 0
> exec
| ALTER TABLE b RENAME COLUMN X TO x;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE q.s (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE s DROP COLUMN IF EXISTS z;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| DROP INDEX I ON v;
< error
| 1091 42000 index I doesn't exist
> exec
| CREATE TABLE b LIKE a;
< ok 0 0

== TestWalkThroughSweepTiDB/change39 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| DROP INDEX i ON a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 0 0
> exec
| ALTER TABLE q.s ADD UNIQUE INDEX I (z);
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| CREATE TABLE q.v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE q.s ADD COLUMN y int;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| DROP INDEX j ON q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE v CHANGE y z int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| USE shop;
< ok 0 0
> exec
| ALTER TABLE s ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE A RENAME COLUMN z TO y;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE OR REPLACE VIEW s AS SELECT 1 AS x;
< ok 0 0
> exec
| DROP INDEX j ON b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE v CHANGE z x int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change39 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| s
: bytes
| x
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| s
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| s
: bytes
| VIEW
: bytes
|
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	q	v	id	
| column	q	v	x	
| column	q	v	y	
| column	shop	s	x	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	q	v	id	
| column	q	v	x	
| column	q	v	y	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	q	v	id	
| column	q	v	x	
| column	q	v	y	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	q	v	id	
| column	q	v	x	
| column	q	v	y	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	q	v	id	
| column	q	v	x	
| column	q	v	y	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	q	v	id	
| column	q	v	x	
| column	q	v	y	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	q	v	id	
| column	q	v	x	
| column	q	v	y	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	q	v	id	
| column	q	v	x	
| column	q	v	y	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	q	v	id	
| column	q	v	x	
| column	q	v	y	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	q	v	id	
| column	q	v	x	
| column	q	v	y	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	q	s	x	
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		s	VIEW
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	q	s	x	
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		s	VIEW
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	q	s	x	
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		s	VIEW
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	q	s	x	
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		s	VIEW
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	q	s	x	
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| relation	q		s	VIEW
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	q	s	x	
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	j	
| relation	q		s	VIEW
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	q	s	x	
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	j	
| relation	q		s	VIEW
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	q	s	x	
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	j	
| relation	q		s	VIEW
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	q	s	x	
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	j	
| relation	q		s	VIEW
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	q	s	x	
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	j	
| relation	q		s	VIEW
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	q	s	x	
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	j	
| relation	q		s	VIEW
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	q	s	x	
| column	q	v	id	
| column	q	v	x	
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	v	PRIMARY	
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	j	
| relation	q		s	VIEW
| relation	q		v	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	s	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	SEQUENCE
| relation	shop		s	VIEW
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change39 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE a (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE A CHANGE z z int;
< error
| 1054 42S22 Unknown column 'z' in 'A'
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP VIEW IF EXISTS q.b, v;
< error
| 1347 HY000 'shop.v' is not VIEW
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE q.A LIKE s;
< error
| 1347 HY000 'q.s' is not BASE TABLE
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS c;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| TRUNCATE TABLE q.s;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| RENAME TABLE c TO a, v TO A;
< error
| 1050 42S01 Table 'shop.a' already exists
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE q.c CHANGE y z int;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE q.v DROP COLUMN y;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| UPDATE s SET x = 1;
< error
| 1288 HY000 The target table s of the UPDATE is not updatable
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE q.s CHANGE y X int;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TABLE s (x int, y int, KEY i (x));
< error
| 1050 42S01 Table 'shop.s' already exists
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE VIEW q.s AS SELECT 1 AS x;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE DATABASE q;
< error
| 1007 HY000 Can't create database 'q'; database exists
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| UPDATE b SET x = 1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE b ADD UNIQUE INDEX x (z);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP SEQUENCE s;
< error
| 1347 HY000 'shop.s' is not SEQUENCE
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE v ADD UNIQUE INDEX j (X);
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `j` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE TEMPORARY TABLE v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE q.b RENAME TO v;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE b DROP COLUMN IF EXISTS z;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE A DROP FOREIGN KEY fk_1;
< error
| 1091 42000 Can't DROP 'fk_1'; check that column/key exists
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE q.s ADD UNIQUE INDEX i (y);
< error
| 1347 HY000 'q.s' is not BASE TABLE
> query
| SHOW CREATE TABLE `q`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| CREATE SEQUENCE IF NOT EXISTS b;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec
| ALTER TABLE s DROP FOREIGN KEY F;
< error
| 1347 HY000 'shop.s' is not BASE TABLE
> query
| SHOW CREATE TABLE `shop`.`v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TEMPORARY TABLE `v` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`v`
< error
| 1146 42S02 Table 'test.v' doesn't exist

== TestWalkThroughSweepTiDB/change4 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE c (x int, y int, KEY x (x));
< ok 0 0
> exec
| CREATE UNIQUE INDEX x ON q.a (z);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP TABLE v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| ALTER TABLE c DROP COLUMN X, ADD COLUMN x int;
< error
| 1060 42S21 Duplicate column name 'x'
> exec
| USE shop;
< ok 0 0
> exec
| ALTER TABLE c DROP INDEX j;
< error
| 1091 42000 index j doesn't exist
> exec
| INSERT INTO s (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE q.A LIKE s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE TABLE q.c (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change4 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| c
: bytes
| x
: int64 1
< row
: bytes
| a
: bytes
| PRIMARY
: int64 0
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| id
< row
: bytes
| a
: bytes
| x
< row
: bytes
| a
: bytes
| y
< row
: bytes
| c
: bytes
| x
< row
: bytes
| c
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| c
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	SEQUENCE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| database	test			
| index	shop	a	PRIMARY	
| index	shop	c	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	SEQUENCE
| relation	shop		c	BASE TABLE

== TestWalkThroughSweepTiDB/change4 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP SEQUENCE b;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.b'
> exec
| RENAME TABLE q.a TO b, b TO a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| UPDATE A SET x = 1;
< ok 0 0
> exec
| CREATE TABLE q.a (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE b ADD COLUMN y int, DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE SEQUENCE A;
< error
| 1050 42S01 Table 'shop.A' already exists
> exec
| TRUNCATE TABLE A;
< ok 0 0
> exec
| CREATE OR REPLACE VIEW a AS SELECT 1 AS x;
< error
| 1347 HY000 'shop.a' is not VIEW
> exec
| CREATE TABLE c (x int, y int, KEY i (x));
< error
| 1050 42S01 Table 'shop.c' already exists
> exec
| CREATE TABLE q.A (x int, y int, KEY I (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE INDEX IF NOT EXISTS j ON a (z);
< error
| 1072 42000 column does not exist: z
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS v;
< ok 0 0
> exec
| CREATE SEQUENCE b;
< ok 0 0
> exec
| RENAME TABLE c TO q.A;
< error
| 1025 HY000 Error on rename of 'shop.c' to 'q.A' (errno: 168 - Database `q` doesn't exist)
> exec
| CREATE TABLE v AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE a RENAME TO a;
< ok 0 0
> exec
| DROP INDEX I ON b;
< error
| 1091 42000 index I doesn't exist

== TestWalkThroughSweepTiDB/change5 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| TRUNCATE TABLE q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE A DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE v DROP COLUMN y, ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP INDEX i ON b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP INDEX I ON b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE q.v DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| DROP SEQUENCE v;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.v'
> exec
| ALTER TABLE a CHANGE z x int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE b LIKE q.s;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change5 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	q	b	x	
| column	q	b	y	
| database	q			
| database	shop			
| database	test			
| index	q	b	I	
| relation	q		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	q	b	x	
| column	q	b	y	
| database	q			
| database	shop			
| database	test			
| index	q	b	I	
| relation	q		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	q	b	x	
| column	q	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	b	I	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	q	b	x	
| column	q	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	b	I	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	q	b	x	
| column	q	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	b	I	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	q	b	x	
| column	q	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	b	I	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	q	b	x	
| column	q	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	b	I	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	q	b	x	
| column	q	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	b	I	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	q	b	x	
| column	q	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	b	I	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	q	b	x	
| column	q	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	b	I	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	q	b	x	
| column	q	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	b	I	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	q	b	x	
| column	q	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| database	test			
| index	q	b	I	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		c	SEQUENCE
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change5 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| UPDATE q.A SET x = 1;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| INSERT INTO a (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE q.A AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| INSERT INTO A (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| TRUNCATE TABLE q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TEMPORARY TABLE a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS q.c (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE s DROP COLUMN z;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| DROP VIEW c, s;
< error
| 1051 42S02 Unknown table 'shop.c,shop.s'
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE q.b (x int, y int, KEY I (x));
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< error
| 1146 42S02 Table 'q.c' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TEMPORARY TABLE q.c (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE v (x int, y int, KEY j (x));
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE a ADD CONSTRAINT F FOREIGN KEY (x) REFERENCES b (id);
< error
| 8200 HY000 TiDB doesn't support ALTER TABLE for local temporary table
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE SEQUENCE c;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE SEQUENCE `c` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| DROP SEQUENCE a;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.a'
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE SEQUENCE `c` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE q.A ADD PRIMARY KEY (x);
< error
| 1146 42S02 Table 'q.A' doesn't exist
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE SEQUENCE `c` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| UPDATE b SET x = 1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE SEQUENCE `c` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE SEQUENCE `c` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE SEQUENCE `c` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE VIEW c AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'shop.c' already exists
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE SEQUENCE `c` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE SEQUENCE q.b;
< error
| 1050 42S01 Table 'q.b' already exists
> query
| SHOW CREATE TABLE `q`.`a`
< error
| 1146 42S02 Table 'q.a' doesn't exist
> query
| SHOW CREATE TABLE `q`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `x` (`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TEMPORARY TABLE `a` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE SEQUENCE `c` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist

== TestWalkThroughSweepTiDB/change6 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TEMPORARY TABLE a (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| RENAME TABLE b TO q.a;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE a DROP INDEX x;
< error
| 8200 HY000 TiDB doesn't support ALTER TABLE for local temporary table
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE A ADD COLUMN z int AFTER y;
< error
| 8200 HY000 TiDB doesn't support ALTER TABLE for local temporary table
> exec
| ALTER TABLE a ADD INDEX j (z);
< error
| 8200 HY000 TiDB doesn't support ALTER TABLE for local temporary table
> exec
| ALTER TABLE q.s DROP COLUMN IF EXISTS z;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| CREATE TABLE v (x int, y int, KEY i (x));
< ok 0 0
> exec
| UPDATE c SET x = 1;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP TABLE IF EXISTS c, s;
< ok 0 0
> exec
| ALTER TABLE A DROP COLUMN x;
< error
| 8200 HY000 TiDB doesn't support ALTER TABLE for local temporary table
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change6 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< row
: bytes
| v
: bytes
| i
: int64 1
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
< row
: bytes
| v
: bytes
| x
< row
: bytes
| v
: bytes
| y
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
< row
: bytes
| v
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| database	test			
| index	shop	v	i	
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepTiDB/change6 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE q.a AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE SEQUENCE v;
< error
| 1050 42S01 Table 'shop.v' already exists
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE q.v ADD COLUMN z int, DROP COLUMN x;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE a ADD COLUMN X int, DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE v (x int, y int, KEY I (x));
< error
| 1050 42S01 Table 'shop.v' already exists
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE A RENAME COLUMN y TO y, ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE v LIKE v;
< error
| 1050 42S01 Table 'shop.v' already exists
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| DROP TABLE IF EXISTS a;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE A ADD FOREIGN KEY F (x) REFERENCES v (id);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE c LIKE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE v ADD COLUMN x int, DROP COLUMN X;
< error
| 1060 42S21 Duplicate column name 'x'
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE b AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE q.b RENAME COLUMN x TO X;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TABLE c AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE b COMMENT 'x', RENAME TO A;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| UPDATE q.b SET x = 1;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE A ADD CONSTRAINT g FOREIGN KEY g (x) REFERENCES A (id);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| UPDATE b SET x = 1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE s ADD INDEX j (x);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN z int, DROP COLUMN z;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| TRUNCATE TABLE q.s;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec
| ALTER TABLE s ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| A
: bytes
| CREATE TEMPORARY TABLE `A` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`c`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TEMPORARY TABLE `c` (
|   `id` int NOT NULL,
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> query
| SHOW CREATE TABLE `test`.`c`
< error
| 1146 42S02 Table 'test.c' doesn't exist

== TestWalkThroughSweepTiDB/change7 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| DROP SEQUENCE q.c;
< error
| 4139 HY000 Unknown SEQUENCE: 'q.c'
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP VIEW b, A;
< error
| 1051 42S02 Unknown table 'shop.b,shop.A'
> exec
| RENAME TABLE a TO a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| DROP TABLE a, A;
< error
| 1051 42S02 Unknown table 'shop.a,shop.A'
> exec
| ALTER TABLE q.s CHANGE y x int;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| ALTER TABLE c DROP COLUMN X, ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE c ADD CONSTRAINT fk_1 FOREIGN KEY f (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE q.A DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change7 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| database	test			
| relation	shop		A	SEQUENCE
| relation	shop		b	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW

== TestWalkThroughSweepTiDB/change7 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE VIEW v AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE v DROP FOREIGN KEY F, DROP FOREIGN KEY f;
< error
| 1347 HY000 'shop.v' is not BASE TABLE
> exec
| CREATE TABLE IF NOT EXISTS q.v LIKE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| TRUNCATE TABLE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE INDEX i ON v (x);
< error
| 1105 HY000 Unknown charset 
> exec
| ALTER TABLE b RENAME COLUMN y TO z;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS v LIKE s;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| CREATE SEQUENCE A;
< ok 0 0
> exec
| CREATE VIEW q.a AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE q.s (x int, y int, KEY I (x), CONSTRAINT f FOREIGN KEY f (y) REFERENCES a (id));
< error
| 1049 42000 Unknown database 'q'
> exec
| RENAME TABLE b TO a;
< error
| 1050 42S01 Table 'shop.a' already exists
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec
| ALTER TABLE v ADD COLUMN z int;
< error
| 1347 HY000 'shop.v' is not BASE TABLE
> exec
| ALTER TABLE q.s RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| ALTER TABLE q.v ADD FOREIGN KEY f (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| CREATE TABLE s AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| CREATE OR REPLACE VIEW b AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE TABLE v (id int PRIMARY KEY, x int, y int);
< error
| 1050 42S01 Table 'shop.v' already exists
> exec
| DROP VIEW q.b, c;
< error
| 1051 42S02 Unknown table 'q.b,shop.c'
> exec
| DROP SEQUENCE s, s;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.s,shop.s'
> exec
| CREATE OR REPLACE VIEW c AS SELECT 1 AS x;
< ok 0 0
> exec
| DROP TABLE b;
< error
| 1051 42S02 Unknown table 'shop.b'
> exec
| CREATE SEQUENCE c;
< error
| 1050 42S01 Table 'shop.c' already exists
> exec
| UPDATE q.c SET x = 1;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE q.v ADD COLUMN x int;
< error
| 1146 42S02 Table 'q.v' doesn't exist

== TestWalkThroughSweepTiDB/change8 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| DROP TEMPORARY TABLE c;
< error
| 1051 42S02 Unknown table 'shop.c'
> exec
| ALTER TABLE b ADD UNIQUE INDEX I (X);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP INDEX I ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE q.a RENAME COLUMN z TO x, ADD COLUMN x int;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE UNIQUE INDEX i ON a (y);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP INDEX j ON A;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE TABLE s AS SELECT 1 AS x;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| TRUNCATE TABLE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP SEQUENCE A, q.b;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.A,q.b'
> exec
| ALTER TABLE a DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE q.b (x int, y int, KEY j (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE b RENAME COLUMN z TO y;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE v RENAME COLUMN z TO z;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change8 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	q		c	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	q		c	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	q		c	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	q		c	SEQUENCE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| database	q			
| database	shop			
| database	test			
| relation	q		c	SEQUENCE

== TestWalkThroughSweepTiDB/change8 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE a CHANGE z z int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE q.s LIKE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP TABLE IF EXISTS b;
< ok 0 0
> exec
| ALTER TABLE s DROP INDEX j;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| DROP SEQUENCE c, c;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.c,shop.c'
> exec
| ALTER TABLE q.a DROP FOREIGN KEY g, DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE VIEW q.a AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP INDEX j ON q.c;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE q.s DROP INDEX x;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| CREATE DATABASE q;
< ok 0 0
> exec
| ALTER TABLE s CHANGE X y int, RENAME COLUMN z TO z;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE s RENAME COLUMN X TO X, ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE q.b ADD COLUMN y int;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE SEQUENCE q.c;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS A LIKE q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| DROP TABLE IF EXISTS c;
< ok 0 0
> exec
| DROP VIEW IF EXISTS c, s;
< ok 0 0
> exec
| ALTER TABLE c ADD COLUMN y int, DROP COLUMN y;
< error
| 1146 42S02 Table 'shop.c' doesn't exist

== TestWalkThroughSweepTiDB/change9 default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE s DROP INDEX x;
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| RENAME TABLE q.v TO c;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| ALTER TABLE c ADD COLUMN x int, DROP COLUMN y;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP TABLE q.s, s;
< error
| 1051 42S02 Unknown table 'q.s,shop.s'
> exec
| CREATE VIEW q.a AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE A DROP INDEX j, ADD INDEX x (id);
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE q.b ADD COLUMN x int AFTER X;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec
| DROP TABLE b;
< error
| 1051 42S02 Unknown table 'shop.b'
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughSweepTiDB/change9 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, KEY_NAME, NON_UNIQUE FROM information_schema.TIDB_INDEXES WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "KEY_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
> query
| SELECT TABLE_NAME, IFNULL(COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR"
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT fks.TABLE_NAME, fks.CONSTRAINT_NAME, fks.REFERENCED_TABLE_NAME
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		FROM information_schema.CHECK_CONSTRAINTS cc
| 			JOIN information_schema.TABLE_CONSTRAINTS tc
| 				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = tc.TABLE_SCHEMA
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| shop
< result
< column "SEQUENCE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE, IFNULL(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "IFNULL(ENGINE, '')" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'table foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| ) o WHERE LOWER(o.db) NOT IN ('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	q	a	id	
| column	q	a	x	
| column	q	a	y	
| database	q			
| database	shop			
| database	test			
| index	q	a	PRIMARY	
| relation	q		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	q	a	id	
| column	q	a	x	
| column	q	a	y	
| database	q			
| database	shop			
| database	test			
| index	q	a	PRIMARY	
| relation	q		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	q	a	id	
| column	q	a	x	
| column	q	a	y	
| database	q			
| database	shop			
| database	test			
| index	q	a	PRIMARY	
| relation	q		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	q	a	id	
| column	q	a	x	
| column	q	a	y	
| database	q			
| database	shop			
| database	test			
| index	q	a	PRIMARY	
| relation	q		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| database	test			

== TestWalkThroughSweepTiDB/change9 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| RENAME TABLE v TO A;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| INSERT INTO c (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS q.a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE c ADD COLUMN z int AFTER z;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS s;
< ok 0 0
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| DROP VIEW IF EXISTS q.c;
< ok 0 0
> exec
| ALTER TABLE q.s ADD COLUMN X int AFTER x;
< error
| 1146 42S02 Table 'q.s' doesn't exist
> exec
| ALTER TABLE q.A ADD COLUMN X int, DROP COLUMN y;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE s ADD CONSTRAINT i UNIQUE INDEX i (X);
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> exec
| ALTER TABLE a RENAME COLUMN y TO y, ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP SEQUENCE b;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.b'
> exec
| ALTER TABLE A DROP INDEX I;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE OR REPLACE VIEW q.s AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE OR REPLACE VIEW q.c AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP VIEW A;
< error
| 1051 42S02 Unknown table 'shop.A'
> exec
| CREATE TABLE s (x int, y int, KEY x (x), FOREIGN KEY f (y) REFERENCES a (id));
< ok 0 0
