-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestWalkThroughAgainstTiDB/ADMIN default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/ADMIN default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/ADMIN default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ADMIN SHOW DDL;
< ok 0 0
> exec
| ADMIN SHOW DDL JOBS 5;
< ok 0 0
> exec
| ADMIN CHECK TABLE t, v, sq;
< ok 0 0
> exec
| ADMIN CHECK INDEX t i;
< ok 0 0
> exec
| ADMIN SHOW t NEXT_ROW_ID;
< ok 0 0
> exec
| ADMIN CLEANUP TABLE LOCK t;
< ok 0 0
> exec
| ADMIN CHECK TABLE nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| ADMIN CHECKSUM TABLE nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| ADMIN SHOW nosuch NEXT_ROW_ID;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| ADMIN CLEANUP TABLE LOCK nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| ADMIN CHECK INDEX nosuch i;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| ADMIN RECOVER INDEX nosuch i;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| ADMIN CLEANUP INDEX nosuch i;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| ADMIN FLUSH BINDINGS;
< ok 0 0
> exec
| ADMIN SHOW SLOW TOP 3;
< ok 0 0
> exec
| DROP TABLE nokey;
< ok 0 0
> exec
| ADMIN CHECK TABLE nokey;
< error
| 1146 42S02 Table 'shop.nokey' doesn't exist

== TestWalkThroughAgainstTiDB/ALTER_TABLE_judges_the_table_as_it_was default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/ALTER_TABLE_judges_the_table_as_it_was default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	j	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/ALTER_TABLE_judges_the_table_as_it_was default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (a);
< error
| 1068 42000 Multiple primary key defined
> exec
| ALTER TABLE t ADD COLUMN a int, DROP COLUMN a;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| ALTER TABLE t ADD COLUMN z int, DROP COLUMN z;
< error
| 1091 42000 Can't DROP 'z'; check that column/key exists
> exec
| ALTER TABLE t DROP COLUMN a, ADD COLUMN a int;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| ALTER TABLE t DROP INDEX i, ADD INDEX i (b);
< error
| 1061 42000 Duplicate key name 'i'
> exec
| ALTER TABLE t RENAME COLUMN a TO b2, RENAME COLUMN b TO a;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| ALTER TABLE t ADD COLUMN y1 int, ADD COLUMN y2 int AFTER y1;
< error
| 1054 42S22 Unknown column 'y1' in 't'
> exec
| ALTER TABLE t ADD COLUMN x int, MODIFY x bigint;
< error
| 1054 42S22 Unknown column 'x' in 't'
> exec
| ALTER TABLE t RENAME COLUMN a TO c, ADD COLUMN a int;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| ALTER TABLE t ADD INDEX m (a), RENAME INDEX m TO m2;
< error
| 1176 42000 Key 'm' doesn't exist in table 't'
> exec
| ALTER TABLE t DROP INDEX i, ADD INDEX j (b);
< ok 0 0
> exec
| ALTER TABLE t DROP INDEX j, DROP COLUMN a;
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN a int;
< ok 0 0

== TestWalkThroughAgainstTiDB/ALTER_TABLE_subcommands default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/ALTER_TABLE_subcommands default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/ALTER_TABLE_subcommands default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE nosuch ADD COLUMN a int;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| ALTER TABLE v ADD COLUMN a int;
< error
| 1347 HY000 'shop.v' is not BASE TABLE
> exec
| ALTER TABLE sq ADD COLUMN a int;
< error
| 1347 HY000 'shop.sq' is not BASE TABLE
> exec
| ALTER TABLE sq COMMENT 'x';
< error
| 1347 HY000 'shop.sq' is not BASE TABLE
> exec
| ALTER TABLE t ADD COLUMN A int;
< error
| 1060 42S21 Duplicate column name 'A'
> exec
| ALTER TABLE t ADD COLUMN c int AFTER nosuch;
< error
| 1054 42S22 Unknown column 'nosuch' in 't'
> exec
| ALTER TABLE t DROP COLUMN nosuch;
< error
| 1091 42000 Can't DROP 'nosuch'; check that column/key exists
> exec
| ALTER TABLE t MODIFY COLUMN nosuch int;
< error
| 1054 42S22 Unknown column 'nosuch' in 't'
> exec
| ALTER TABLE t CHANGE nosuch c int;
< error
| 1054 42S22 Unknown column 'nosuch' in 't'
> exec
| ALTER TABLE t CHANGE a b int;
< error
| 1060 42S21 Duplicate column name 'b'
> exec
| ALTER TABLE t ALTER COLUMN nosuch SET DEFAULT 1;
< error
| 1054 42S22 Unknown column 'nosuch' in 't'
> exec
| ALTER TABLE t ADD INDEX I (b);
< error
| 1061 42000 Duplicate key name 'I'
> exec
| ALTER TABLE t ADD INDEX `PRIMARY` (b);
< error
| 1280 42000 Incorrect index name 'PRIMARY'
> exec
| ALTER TABLE t DROP INDEX nosuch;
< error
| 1091 42000 index nosuch doesn't exist
> exec
| ALTER TABLE t ADD PRIMARY KEY (a);
< error
| 1068 42000 Multiple primary key defined
> exec
| ALTER TABLE t RENAME INDEX nosuch TO j;
< error
| 1176 42000 Key 'nosuch' doesn't exist in table 't'
> exec
| ALTER TABLE t ALTER INDEX nosuch INVISIBLE;
< error
| 1176 42000 Key 'nosuch' doesn't exist in table 't'

== TestWalkThroughAgainstTiDB/CREATE_and_DROP_INDEX default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/CREATE_and_DROP_INDEX default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/CREATE_and_DROP_INDEX default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE INDEX k ON nosuch (a);
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| CREATE INDEX i ON t (b);
< error
| 1061 42000 Duplicate key name 'i'
> exec
| CREATE INDEX I ON t (b);
< error
| 1061 42000 Duplicate key name 'I'
> exec
| CREATE INDEX `PRIMARY` ON t (b);
< error
| 1280 42000 Incorrect index name 'PRIMARY'
> exec
| CREATE INDEX IF NOT EXISTS i ON t (b);
< ok 0 0
> exec
| DROP INDEX nosuch ON t;
< error
| 1091 42000 index nosuch doesn't exist
> exec
| DROP INDEX i ON nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| DROP INDEX i ON t;
< ok 0 0
> exec ^
< error
| 1091 42000 index i doesn't exist

== TestWalkThroughAgainstTiDB/DML_targets default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DML_targets default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DML_targets default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| INSERT INTO nosuch VALUES (1);
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| UPDATE nosuch SET a = 1;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| DELETE FROM nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| REPLACE INTO nosuch VALUES (1);
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| INSERT INTO nokey VALUES (1);
< ok 1 0
> exec
| UPDATE nokey SET a = 2;
< ok 1 0
> exec
| DELETE nokey FROM nokey JOIN r ON r.id = nokey.a;
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_COLUMN_IF_EXISTS default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_COLUMN_IF_EXISTS default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	t	q	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DROP_COLUMN_IF_EXISTS default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN IF EXISTS zz;
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN q int, DROP COLUMN IF EXISTS q;
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN q;
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN IF EXISTS a;
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN a;
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| ALTER TABLE t DROP COLUMN q;
< error
| 1091 42000 Can't DROP 'q'; check that column/key exists

== TestWalkThroughAgainstTiDB/DROP_FOREIGN_KEY_twice_drops_it default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_FOREIGN_KEY_twice_drops_it default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	

== TestWalkThroughAgainstTiDB/DROP_FOREIGN_KEY_twice_drops_it default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;
< error
| 1091 42000 Can't DROP 'fk1'; check that column/key exists

== TestWalkThroughAgainstTiDB/DROP_INDEX_IF_EXISTS default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_INDEX_IF_EXISTS default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DROP_INDEX_IF_EXISTS default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP INDEX IF EXISTS nosuch ON t;
< ok 0 0
> exec
| DROP INDEX IF EXISTS i ON nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| DROP INDEX IF EXISTS i ON t;
< ok 0 0
> exec ^
< ok 0 0
> exec
| CREATE INDEX i ON t (a);
< ok 0 0
> exec
| DROP INDEX IF EXISTS i ON t;
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_SEQUENCE_drops_the_sequences_before_a_missing_name default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_SEQUENCE_drops_the_sequences_before_a_missing_name default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DROP_SEQUENCE_drops_the_sequences_before_a_missing_name default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP SEQUENCE sq, nosuch;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.nosuch'

== TestWalkThroughAgainstTiDB/DROP_SEQUENCE_drops_the_sequences_before_a_table default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_SEQUENCE_drops_the_sequences_before_a_table default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DROP_SEQUENCE_drops_the_sequences_before_a_table default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP SEQUENCE sq, t;
< error
| 1347 HY000 'shop.t' is not SEQUENCE

== TestWalkThroughAgainstTiDB/DROP_TABLE default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_TABLE default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DROP_TABLE default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE nosuch;
< error
| 1051 42S02 Unknown table 'shop.nosuch'
> exec
| DROP TABLE v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| DROP TABLE sq;
< error
| 1051 42S02 Unknown table 'shop.sq'
> exec
| DROP TABLE IF EXISTS nosuch, nokey;
< ok 0 0
> exec
| DROP TABLE nokey;
< error
| 1051 42S02 Unknown table 'shop.nokey'
> exec
| DROP VIEW t;
< error
| 1347 HY000 'shop.t' is not VIEW
> exec
| DROP VIEW IF EXISTS t;
< error
| 1347 HY000 'shop.t' is not VIEW
> exec
| DROP VIEW IF EXISTS sq;
< error
| 1347 HY000 'shop.sq' is not VIEW
> exec
| DROP VIEW nosuch;
< error
| 1051 42S02 Unknown table 'shop.nosuch'
> exec
| DROP VIEW IF EXISTS nosuch, v;
< ok 0 0
> exec
| DROP VIEW v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| TRUNCATE TABLE nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| TRUNCATE TABLE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| TRUNCATE TABLE sq;
< error
| 1146 42S02 Table 'shop.sq' doesn't exist

== TestWalkThroughAgainstTiDB/DROP_TABLE_IF_EXISTS_leaves_a_sequence default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_TABLE_IF_EXISTS_leaves_a_sequence default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DROP_TABLE_IF_EXISTS_leaves_a_sequence default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE IF EXISTS sq;
< ok 0 0
> exec
| DROP SEQUENCE sq;
< ok 0 0
> exec ^
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.sq'

== TestWalkThroughAgainstTiDB/DROP_TABLE_IF_EXISTS_leaves_a_view default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 0 0
> exec
| CREATE VIEW q.v AS SELECT 1 AS a;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_TABLE_IF_EXISTS_leaves_a_view default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 27
: bytes
| column	q	v	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		v	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	q	v	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		v	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DROP_TABLE_IF_EXISTS_leaves_a_view default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE IF EXISTS q.v;
< ok 0 0
> exec
| DROP VIEW q.v;
< ok 0 0
> exec ^
< error
| 1051 42S02 Unknown table 'q.v'

== TestWalkThroughAgainstTiDB/DROP_TABLE_drops_the_tables_before_a_missing_name default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_TABLE_drops_the_tables_before_a_missing_name default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	

== TestWalkThroughAgainstTiDB/DROP_TABLE_drops_the_tables_before_a_missing_name default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE t, nosuch;
< error
| 1051 42S02 Unknown table 'shop.nosuch'

== TestWalkThroughAgainstTiDB/DROP_TABLE_drops_the_tables_before_a_missing_name,_listed_last default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_TABLE_drops_the_tables_before_a_missing_name,_listed_last default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DROP_TABLE_drops_the_tables_before_a_missing_name,_listed_last default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE nosuch, nokey;
< error
| 1051 42S02 Unknown table 'shop.nosuch'

== TestWalkThroughAgainstTiDB/DROP_TABLE_drops_the_tables_before_a_sequence default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_TABLE_drops_the_tables_before_a_sequence default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DROP_TABLE_drops_the_tables_before_a_sequence default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE nokey, sq;
< error
| 1051 42S02 Unknown table 'shop.sq'

== TestWalkThroughAgainstTiDB/DROP_TABLE_drops_the_tables_before_a_view default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_TABLE_drops_the_tables_before_a_view default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DROP_TABLE_drops_the_tables_before_a_view default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE nokey, v;
< error
| 1051 42S02 Unknown table 'shop.v'

== TestWalkThroughAgainstTiDB/DROP_VIEW_IF_EXISTS_drops_the_views_before_a_table default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_VIEW_IF_EXISTS_drops_the_views_before_a_table default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DROP_VIEW_IF_EXISTS_drops_the_views_before_a_table default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP VIEW IF EXISTS v, t;
< error
| 1347 HY000 'shop.t' is not VIEW

== TestWalkThroughAgainstTiDB/DROP_VIEW_drops_the_views_before_a_missing_name default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/DROP_VIEW_drops_the_views_before_a_missing_name default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/DROP_VIEW_drops_the_views_before_a_missing_name default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP VIEW v, nosuch;
< error
| 1051 42S02 Unknown table 'shop.nosuch'

== TestWalkThroughAgainstTiDB/ENGINE_is_ignored default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/ENGINE_is_ignored default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	my	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	my	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	my	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	my	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	my	PRIMARY	
| index	shop	my	f	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	my	f	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	my	PRIMARY	
| index	shop	my	f	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	my	f	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/ENGINE_is_ignored default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE my (id int PRIMARY KEY, a int, FOREIGN KEY (a) REFERENCES nosuch (id)) ENGINE=MyISAM;
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'
> exec
| CREATE TABLE my (id int PRIMARY KEY, a int) ENGINE=MyISAM;
< ok 0 0
> exec
| ALTER TABLE my DROP FOREIGN KEY nosuch;
< error
| 1091 42000 Can't DROP 'nosuch'; check that column/key exists
> exec
| ALTER TABLE my ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES nosuch (id);
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'
> exec
| ALTER TABLE my ENGINE=MyISAM;
< ok 0 0
> exec
| ALTER TABLE my ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id);
< ok 0 0
> exec
| ALTER TABLE my ADD CONSTRAINT F FOREIGN KEY (a) REFERENCES r (id);
< error
| 1826 HY000 Duplicate foreign key constraint name 'F'

== TestWalkThroughAgainstTiDB/EXPLAIN default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/EXPLAIN default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	c	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	c	
| column	shop	t	d	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	c	
| column	shop	t	d	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	c	
| column	shop	t	d	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	c	
| column	shop	t	d	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	c	
| column	shop	t	d	
| column	shop	t	f	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	c	
| column	shop	t	d	
| column	shop	t	f	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	c	
| column	shop	t	d	
| column	shop	t	f	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	c	
| column	shop	t	d	
| column	shop	t	f	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/EXPLAIN default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| EXPLAIN ALTER TABLE t ADD COLUMN c int;
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN c int;
< ok 0 0
> exec
| EXPLAIN ANALYZE ALTER TABLE t ADD COLUMN d int;
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN d int;
< error
| 1060 42S21 Duplicate column name 'd'
> exec
| EXPLAIN ANALYZE ALTER TABLE nosuch ADD COLUMN e int;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| EXPLAIN FORMAT = 'brief' ALTER TABLE t ADD COLUMN f int;
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN f int;
< ok 0 0
> exec
| EXPLAIN SELECT 1;
< ok 0 0
> exec
| EXPLAIN ANALYZE SELECT 1;
< ok 0 0
> exec
| EXPLAIN t;
< ok 0 0
> exec
| EXPLAIN ANALYZE UPDATE t SET a = 1;
< ok 0 0

== TestWalkThroughAgainstTiDB/LOCK_STATS_and_UNLOCK_STATS default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/LOCK_STATS_and_UNLOCK_STATS default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/LOCK_STATS_and_UNLOCK_STATS default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| LOCK STATS t;
< ok 0 0
> exec
| UNLOCK STATS t;
< ok 0 0
> exec
| LOCK STATS t, r;
< ok 0 0
> exec
| UNLOCK STATS t, r;
< ok 0 0
> exec
| LOCK STATS test.t;
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| UNLOCK STATS test.t;
< error
| 1146 42S02 Table 'test.t' doesn't exist

== TestWalkThroughAgainstTiDB/RENAME default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/RENAME default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t	PRIMARY	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t	PRIMARY	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t	PRIMARY	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t2	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t2	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t2	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	q	nokey	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| relation	q		nokey	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t2	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	q	nokey	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| relation	q		nokey	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t2	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	q	nokey	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| relation	q		nokey	BASE TABLE
| relation	q		sq	SEQUENCE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t2	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	q	nokey	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| relation	q		nokey	BASE TABLE
| relation	q		sq	SEQUENCE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t2	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	q	nokey	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| relation	q		nokey	BASE TABLE
| relation	q		sq	SEQUENCE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t2	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	q	nokey	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	q.t2	id	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	q.t2	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| relation	q		nokey	BASE TABLE
| relation	q		sq	SEQUENCE
| relation	shop		child	BASE TABLE
| relation	shop		q.t2	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	q	nokey	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	q.t2	id	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	q.t2	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| relation	q		nokey	BASE TABLE
| relation	q		sq	SEQUENCE
| relation	shop		child	BASE TABLE
| relation	shop		q.t2	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	q	nokey	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	q.t2	id	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	q.t2	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| relation	q		nokey	BASE TABLE
| relation	q		sq	SEQUENCE
| relation	shop		child	BASE TABLE
| relation	shop		q.t2	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	r	fk1	

== TestWalkThroughAgainstTiDB/RENAME default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| RENAME TABLE nosuch TO x;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| RENAME TABLE t TO r;
< error
| 1050 42S01 Table 'shop.r' already exists
> exec
| RENAME TABLE t TO tmp, r TO t, tmp TO r;
< ok 0 0
> exec
| ALTER TABLE r ADD COLUMN a int;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| ALTER TABLE t RENAME TO r;
< error
| 1050 42S01 Table 'shop.r' already exists
> exec
| ALTER TABLE t RENAME TO t2;
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN zz int;
< error
| 1146 42S02 Table 'shop.t' doesn't exist
> exec
| CREATE DATABASE q;
< ok 0 0
> exec
| RENAME TABLE nokey TO q.nokey;
< ok 0 0
> exec
| CREATE TABLE q.nokey (a int);
< error
| 1050 42S01 Table 'q.nokey' already exists
> exec
| RENAME TABLE sq TO q.sq;
< ok 0 0
> exec
| CREATE SEQUENCE q.sq;
< error
| 1050 42S01 Table 'q.sq' already exists
> exec
| RENAME TABLE t2 TO t3, nosuch TO t4;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| ALTER TABLE t2 RENAME TO `q.t2`;
< ok 0 0
> exec
| CREATE TABLE `q.t2` (a int);
< error
| 3734 HY000 Failed to add the foreign key constraint. Missing column 'id' for constraint 'fk_1' in the referenced table 'q.t2'
> exec
| RENAME TABLE r TO R;
< error
| 1050 42S01 Table 'shop.R' already exists
> exec
| ALTER TABLE r ADD COLUMN zz int, RENAME TO R;
< ok 0 0

== TestWalkThroughAgainstTiDB/RENAME_TO_beside_other_subcommands default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/RENAME_TO_beside_other_subcommands default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t2	a	
| column	shop	t2	b	
| column	shop	t2	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t2	PRIMARY	
| index	shop	t2	fk1	
| index	shop	t2	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t2	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t2	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t3	a	
| column	shop	t3	b	
| column	shop	t3	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t3	PRIMARY	
| index	shop	t3	fk1	
| index	shop	t3	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t3	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t3	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/RENAME_TO_beside_other_subcommands default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE t COMMENT 'x', RENAME TO t2;
< error
| 8200 HY000 Unsupported multi schema change for rename table
> exec
| ALTER TABLE t RENAME TO t2, ADD COLUMN c int;
< error
| 8200 HY000 Unsupported multi schema change for rename table
> exec
| ALTER TABLE t ENGINE = InnoDB, RENAME TO t2;
< error
| 8200 HY000 Unsupported multi schema change for rename table
> exec
| ALTER TABLE t RENAME INDEX i TO i2, RENAME TO t2;
< error
| 8200 HY000 Unsupported multi schema change for rename table
> exec
| ALTER TABLE t COMMENT 'x', RENAME TO t;
< ok 0 0
> exec
| ALTER TABLE t ALGORITHM = INSTANT, RENAME TO t2;
< ok 0 0
> exec
| ALTER TABLE t2 LOCK = NONE, RENAME TO t3;
< ok 0 0
> exec
| ALTER TABLE t3 RENAME TO t;
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN a int;
< error
| 1060 42S21 Duplicate column name 'a'

== TestWalkThroughAgainstTiDB/SET_foreign_key_checks_in_any_scope default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/SET_foreign_key_checks_in_any_scope default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	z	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	z	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	z	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	z	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	z	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	z	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	z	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	z	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 36
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z	a	
| column	shop	z3	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	z	fk_1	
| index	shop	z3	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| relation	shop		z3	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	z	fk_1	
| table foreign key	shop	z3	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 36
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z	a	
| column	shop	z3	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	z	fk_1	
| index	shop	z3	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| relation	shop		z3	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	z	fk_1	
| table foreign key	shop	z3	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 36
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z	a	
| column	shop	z3	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	z	fk_1	
| index	shop	z3	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| relation	shop		z3	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	z	fk_1	
| table foreign key	shop	z3	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 36
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z	a	
| column	shop	z3	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	z	fk_1	
| index	shop	z3	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| relation	shop		z3	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	z	fk_1	
| table foreign key	shop	z3	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 40
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z	a	
| column	shop	z3	a	
| column	shop	z5	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	z	fk_1	
| index	shop	z3	fk_1	
| index	shop	z5	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| relation	shop		z3	BASE TABLE
| relation	shop		z5	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	z	fk_1	
| table foreign key	shop	z3	fk_1	
| table foreign key	shop	z5	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 40
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z	a	
| column	shop	z3	a	
| column	shop	z5	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	z	fk_1	
| index	shop	z3	fk_1	
| index	shop	z5	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| relation	shop		z3	BASE TABLE
| relation	shop		z5	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	z	fk_1	
| table foreign key	shop	z3	fk_1	
| table foreign key	shop	z5	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 40
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z	a	
| column	shop	z3	a	
| column	shop	z5	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	z	fk_1	
| index	shop	z3	fk_1	
| index	shop	z5	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| relation	shop		z3	BASE TABLE
| relation	shop		z5	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	z	fk_1	
| table foreign key	shop	z3	fk_1	
| table foreign key	shop	z5	fk_1	

== TestWalkThroughAgainstTiDB/SET_foreign_key_checks_in_any_scope default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1, @@session.foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1, foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1;
< ok 0 0
> exec
| CREATE TABLE z2 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'
> exec
| SET @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE z3 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 0, @@session.foreign_key_checks = 1;
< ok 0 0
> exec
| CREATE TABLE z4 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'
> exec
| SET @@session.foreign_key_checks = 1, @@global.foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE z5 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1;
< ok 0 0
> exec
| CREATE TABLE w (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'

== TestWalkThroughAgainstTiDB/SET_foreign_key_checks_later_in_the_list default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/SET_foreign_key_checks_later_in_the_list default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	

== TestWalkThroughAgainstTiDB/SET_foreign_key_checks_later_in_the_list default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| SET @x = 1, GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< ok 0 0
> exec
| SET @x = 1, SESSION foreign_key_checks = 1;
< ok 0 0
> exec
| CREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'

== TestWalkThroughAgainstTiDB/SET_foreign_key_checks_to_a_word default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/SET_foreign_key_checks_to_a_word default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	w	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	w	fk_1	
| index	shop	y	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		w	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	w	fk_1	
| table foreign key	shop	y	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	w	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	w	fk_1	
| index	shop	y	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		w	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	w	fk_1	
| table foreign key	shop	y	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	w	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	w	fk_1	
| index	shop	y	fk_1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		w	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	w	fk_1	
| table foreign key	shop	y	fk_1	

== TestWalkThroughAgainstTiDB/SET_foreign_key_checks_to_a_word default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| SET foreign_key_checks = OFF;
< ok 0 0
> exec
| CREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< ok 0 0
> exec
| SET foreign_key_checks = ON;
< ok 0 0
> exec
| CREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'
> exec
| SET foreign_key_checks = 'off';
< ok 0 0
> exec
| CREATE TABLE w (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< ok 0 0
> exec
| SET foreign_key_checks = 'on';
< ok 0 0
> exec
| CREATE TABLE v (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'

== TestWalkThroughAgainstTiDB/USE_moves_unqualified_names default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/USE_moves_unqualified_names default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	q	t	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	q	t	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	q	t	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	q	t	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/USE_moves_unqualified_names default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 0 0
> exec
| USE q;
< ok 0 0
> exec
| CREATE TABLE t (a int);
< ok 0 0
> exec ^
< error
| 1050 42S01 Table 'q.t' already exists
> exec
| USE shop;
< ok 0 0
> exec
| CREATE TABLE t (a int);
< error
| 1050 42S01 Table 'shop.t' already exists

== TestWalkThroughAgainstTiDB/a_CHECK_constraint's_long_name default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0
> exec
| SET GLOBAL tidb_enable_check_constraint = OFF
< ok 0 0

== TestWalkThroughAgainstTiDB/a_CHECK_constraint's_long_name default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	k1	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		k1	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	k1	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		k1	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	k1	a	
| column	shop	k2	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		k1	BASE TABLE
| relation	shop		k2	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	k1	a	
| column	shop	k2	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		k1	BASE TABLE
| relation	shop		k2	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	k1	a	
| column	shop	k2	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		k1	BASE TABLE
| relation	shop		k2	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/a_CHECK_constraint's_long_name default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE k1 (a int, CONSTRAINT aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa CHECK (a > 0));
< ok 0 0
> exec
| ALTER TABLE k1 ADD CONSTRAINT aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaax CHECK (a > 1);
< ok 0 0
> exec
| CREATE TABLE k2 (a int, CHECK (a > 0), CONSTRAINT aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa CHECK (a > 1));
< ok 0 0
> exec
| ALTER TABLE k2 DROP CHECK aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa;
< ok 0 0
> exec
| CREATE TABLE aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa (a int);
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long

== TestWalkThroughAgainstTiDB/a_constraint's_name_and_its_index's default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/a_constraint's_name_and_its_index's default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	cx	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	x	cx	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	cx	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	cx	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fy	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y	fy	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fy	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	t	nk	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 33
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	u	a	
| column	shop	u	b	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	u	cu	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		u	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	t	nk	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	u	a	
| column	shop	u	b	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		u	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	t	nk	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	u	a	
| column	shop	u	b	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		u	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	t	nk	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 33
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	u	a	
| column	shop	u	b	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	u	cv	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		u	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	t	nk	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	u	a	
| column	shop	u	b	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		u	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	t	nk	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 36
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	u	a	
| column	shop	u	b	
| column	shop	v	id	
| column	shop	w	a	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	w	cw	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		u	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		w	BASE TABLE
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	t	nk	
| table foreign key	shop	w	cw	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 36
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	u	a	
| column	shop	u	b	
| column	shop	v	id	
| column	shop	w	a	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	w	cw	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		u	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		w	BASE TABLE
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	t	nk	
| table foreign key	shop	w	cw	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 39
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	u	a	
| column	shop	u	b	
| column	shop	v	id	
| column	shop	w	a	
| column	shop	w3	a	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	w	cw	
| index	shop	w3	cw3	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		u	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		w	BASE TABLE
| relation	shop		w3	BASE TABLE
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	t	nk	
| table foreign key	shop	w	cw	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 39
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	u	a	
| column	shop	u	b	
| column	shop	v	id	
| column	shop	w	a	
| column	shop	w3	a	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	w	cw	
| index	shop	w3	cw3	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		u	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		w	BASE TABLE
| relation	shop		w3	BASE TABLE
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	t	nk	
| table foreign key	shop	w	cw	

== TestWalkThroughAgainstTiDB/a_constraint's_name_and_its_index's default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a int, CONSTRAINT cx FOREIGN KEY fx (a) REFERENCES r (id));
< ok 0 0
> exec
| ALTER TABLE x DROP FOREIGN KEY cx;
< ok 0 0
> exec ^
< error
| 1091 42000 Can't DROP 'cx'; check that column/key exists
> exec
| ALTER TABLE x DROP INDEX cx;
< ok 0 0
> exec
| CREATE TABLE y (a int, FOREIGN KEY fy (a) REFERENCES r (id));
< ok 0 0
> exec
| ALTER TABLE y DROP FOREIGN KEY fy;
< ok 0 0
> exec
| ALTER TABLE y DROP INDEX fy;
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT nk FOREIGN KEY fk1 (a) REFERENCES r (id);
< ok 0 0
> exec
| CREATE TABLE u (a int, b int, CONSTRAINT cu UNIQUE KEY uu (a));
< ok 0 0
> exec
| ALTER TABLE u DROP INDEX cu;
< ok 0 0
> exec ^
< error
| 1091 42000 index cu doesn't exist
> exec
| ALTER TABLE u ADD CONSTRAINT cv UNIQUE INDEX uv (b);
< ok 0 0
> exec
| ALTER TABLE u DROP INDEX cv;
< ok 0 0
> exec
| CREATE TABLE w (a int, CONSTRAINT cw FOREIGN KEY aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa (a) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE w2 (a int, FOREIGN KEY aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa (a) REFERENCES r (id));
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec
| CREATE TABLE w3 (a int, CONSTRAINT cw3 UNIQUE KEY aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa (a));
< ok 0 0
> exec
| CREATE TABLE w4 (a int, UNIQUE KEY aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa (a));
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long

== TestWalkThroughAgainstTiDB/a_local_temporary_table_ends_the_walk default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/a_local_temporary_table_ends_the_walk default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/a_local_temporary_table_ends_the_walk default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE TEMPORARY TABLE tt (a int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`tt`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tt
: bytes
| CREATE TEMPORARY TABLE `tt` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| ALTER TABLE tt ADD COLUMN b int;
< error
| 8200 HY000 TiDB doesn't support ALTER TABLE for local temporary table
> query
| SHOW CREATE TABLE `shop`.`tt`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tt
: bytes
| CREATE TEMPORARY TABLE `tt` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE INDEX k ON tt (a);
< error
| 8200 HY000 TiDB doesn't support CREATE INDEX for local temporary table
> query
| SHOW CREATE TABLE `shop`.`tt`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tt
: bytes
| CREATE TEMPORARY TABLE `tt` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| RENAME TABLE tt TO tt2;
< error
| 8200 HY000 TiDB doesn't support RENAME TABLE for local temporary table
> query
| SHOW CREATE TABLE `shop`.`tt`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tt
: bytes
| CREATE TEMPORARY TABLE `tt` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE TABLE t (a int);
< error
| 1050 42S01 Table 'shop.t' already exists
> query
| SHOW CREATE TABLE `shop`.`tt`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tt
: bytes
| CREATE TEMPORARY TABLE `tt` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist

== TestWalkThroughAgainstTiDB/a_new_table's_text default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/a_new_table's_text default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/a_new_table's_text default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a int, a int);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| CREATE TABLE x (a int, b int, KEY k (a), KEY k (b));
< error
| 1061 42000 Duplicate key name 'k'
> exec
| CREATE TABLE x (a int PRIMARY KEY, b int, PRIMARY KEY (b));
< error
| 1068 42000 Multiple primary key defined
> exec
| CREATE TABLE x (a int, KEY `PRIMARY` (a));
< error
| 1280 42000 Incorrect index name 'PRIMARY'

== TestWalkThroughAgainstTiDB/a_new_table_holds_what_it_declares default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/a_new_table_holds_what_it_declares default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	x	fx	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	x	fx	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	x	fx	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	x	fx	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	x	fx	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	x	fx	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	x	fx	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	fx	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	fx	
| table foreign key	shop	t	fk1	
| table foreign key	shop	x	fx	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 31
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	fx	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	fx	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 31
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	fx	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	fx	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 36
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| column	shop	y	a	
| column	shop	y	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	fx	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| index	shop	y	b	
| index	shop	y	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	fx	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/a_new_table_holds_what_it_declares default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a int, b int UNIQUE, KEY k (a), CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id));
< ok 0 0
> exec
| ALTER TABLE x ADD COLUMN a int;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| ALTER TABLE x ADD INDEX k (b);
< error
| 1061 42000 Duplicate key name 'k'
> exec
| ALTER TABLE x DROP COLUMN zz;
< error
| 1091 42000 Can't DROP 'zz'; check that column/key exists
> exec
| ALTER TABLE x DROP INDEX k, DROP INDEX k;
< error
| 1553 HY000 Cannot drop index 'k': needed in a foreign key constraint
> exec
| ALTER TABLE x DROP PRIMARY KEY;
< error
| 1091 42000 Can't DROP 'PRIMARY'; check that column/key exists
> exec
| ALTER TABLE x ADD CONSTRAINT FX FOREIGN KEY (a) REFERENCES r (id);
< error
| 1826 HY000 Duplicate foreign key constraint name 'FX'
> exec
| ALTER TABLE nokey ADD CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id);
< ok 0 0
> exec
| ALTER TABLE x DROP FOREIGN KEY fx;
< ok 0 0
> exec ^
< error
| 1091 42000 Can't DROP 'fx'; check that column/key exists
> exec
| CREATE TABLE y LIKE x;
< ok 0 0
> exec
| ALTER TABLE y DROP INDEX k;
< ok 0 0

== TestWalkThroughAgainstTiDB/a_string_outside_ASCII_is_no_name default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/a_string_outside_ASCII_is_no_name default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/a_string_outside_ASCII_is_no_name default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE y (a varchar(10)) COMMENT '中文';
< ok 0 0
> exec
| CREATE TABLE y (a int);
< error
| 1050 42S01 Table 'shop.y' already exists
> exec
| INSERT INTO y VALUES ('é');
< ok 1 0
> exec
| INSERT INTO nosuch VALUES ('é');
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist

== TestWalkThroughAgainstTiDB/a_table_without_a_primary_key_has_none default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/a_table_without_a_primary_key_has_none default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	g	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		g	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	g	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		g	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	g	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		g	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	g	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	g	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		g	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	g	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	g	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		g	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/a_table_without_a_primary_key_has_none default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE nokey DROP PRIMARY KEY;
< error
| 1091 42000 Can't DROP 'PRIMARY'; check that column/key exists
> exec
| CREATE TABLE g (a int);
< ok 0 0
> exec
| ALTER TABLE g DROP PRIMARY KEY;
< error
| 1091 42000 Can't DROP 'PRIMARY'; check that column/key exists
> exec
| ALTER TABLE g DROP INDEX `PRIMARY`;
< error
| 1091 42000 Can't DROP 'PRIMARY'; check that column/key exists
> exec
| ALTER TABLE g ADD PRIMARY KEY (a);
< ok 0 0
> exec ^
< error
| 1068 42000 Multiple primary key defined

== TestWalkThroughAgainstTiDB/a_temporary_table's_definition default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/a_temporary_table's_definition default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/a_temporary_table's_definition default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE TEMPORARY TABLE tmp (a int, a int);
< error
| 1060 42S21 Duplicate column name 'a'
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE TEMPORARY TABLE tmp (a int, KEY k (a), KEY k (a));
< error
| 1061 42000 Duplicate key name 'k'
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE TEMPORARY TABLE tmp (a int PRIMARY KEY, PRIMARY KEY (a));
< error
| 1068 42000 Multiple primary key defined
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE TEMPORARY TABLE tmp (a int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`tmp`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tmp
: bytes
| CREATE TEMPORARY TABLE `tmp` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS tmp (a int, a int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`tmp`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tmp
: bytes
| CREATE TEMPORARY TABLE `tmp` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE TEMPORARY TABLE tmp (a int);
< error
| 1050 42S01 Table 'shop.tmp' already exists
> query
| SHOW CREATE TABLE `shop`.`tmp`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tmp
: bytes
| CREATE TEMPORARY TABLE `tmp` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE TEMPORARY TABLE TMP (a int);
< error
| 1050 42S01 Table 'shop.TMP' already exists
> query
| SHOW CREATE TABLE `shop`.`tmp`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tmp
: bytes
| CREATE TEMPORARY TABLE `tmp` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist

== TestWalkThroughAgainstTiDB/a_temporary_table_is_no_parent default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/a_temporary_table_is_no_parent default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/a_temporary_table_is_no_parent default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`r`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| r
: bytes
| CREATE TABLE `r` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tp`
< error
| 1146 42S02 Table 'shop.tp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`r`
< error
| 1146 42S02 Table 'test.r' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tp`
< error
| 1146 42S02 Table 'test.tp' doesn't exist
> exec
| CREATE TEMPORARY TABLE tp (id int PRIMARY KEY);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`r`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| r
: bytes
| CREATE TABLE `r` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tp`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tp
: bytes
| CREATE TEMPORARY TABLE `tp` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`r`
< error
| 1146 42S02 Table 'test.r' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tp`
< error
| 1146 42S02 Table 'test.tp' doesn't exist
> exec
| CREATE TABLE y1 (a int, FOREIGN KEY (a) REFERENCES tp (id));
< error
| 1824 HY000 Failed to open the referenced table 'tp'
> query
| SHOW CREATE TABLE `shop`.`r`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| r
: bytes
| CREATE TABLE `r` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tp`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tp
: bytes
| CREATE TEMPORARY TABLE `tp` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`r`
< error
| 1146 42S02 Table 'test.r' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tp`
< error
| 1146 42S02 Table 'test.tp' doesn't exist
> exec
| CREATE TEMPORARY TABLE r (id int PRIMARY KEY);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`r`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| r
: bytes
| CREATE TEMPORARY TABLE `r` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tp`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tp
: bytes
| CREATE TEMPORARY TABLE `tp` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`r`
< error
| 1146 42S02 Table 'test.r' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tp`
< error
| 1146 42S02 Table 'test.tp' doesn't exist
> exec
| CREATE TABLE y2 (a int, FOREIGN KEY (a) REFERENCES r (id));
< ok 0 0

== TestWalkThroughAgainstTiDB/a_view_stays_in_its_database default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/a_view_stays_in_its_database default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/a_view_stays_in_its_database default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 0 0
> exec
| RENAME TABLE v TO q.v;
< error
| 1450 HY000 Changing schema from 'shop' to 'q' is not allowed.
> exec
| DROP VIEW v;
< ok 0 0

== TestWalkThroughAgainstTiDB/an_index_on_a_view_or_a_sequence default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/an_index_on_a_view_or_a_sequence default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	v	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	sq	k	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	v	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	v	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	v	k	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/an_index_on_a_view_or_a_sequence default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE INDEX k ON v ((1 + 1));
< ok 0 0
> exec
| CREATE INDEX k ON sq ((1 + 1));
< ok 0 0
> exec
| DROP INDEX k ON sq;
< ok 0 0
> exec
| CREATE TABLE t (a int);
< error
| 1050 42S01 Table 'shop.t' already exists

== TestWalkThroughAgainstTiDB/an_unqualified_parent_is_in_the_child's_database default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/an_unqualified_parent_is_in_the_child's_database default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	q	parent	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	q	parent	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		parent	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 32
: bytes
| column	q	child	id	
| column	q	parent	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	q	child	fk_1	
| index	q	parent	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		child	BASE TABLE
| relation	q		parent	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	q	child	fk_1	
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 33
: bytes
| column	q	child	id	
| column	q	parent	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	q	child	fk_1	
| index	q	parent	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		child	BASE TABLE
| relation	q		parent	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	q	child	c2	
| table foreign key	q	child	fk_1	
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 33
: bytes
| column	q	child	id	
| column	q	parent	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	q	child	fk_1	
| index	q	parent	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		child	BASE TABLE
| relation	q		parent	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	q	child	c2	
| table foreign key	q	child	fk_1	
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/an_unqualified_parent_is_in_the_child's_database default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 0 0
> exec
| CREATE TABLE q.parent (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE q.child (id int, FOREIGN KEY (id) REFERENCES parent (id));
< ok 0 0
> exec
| ALTER TABLE q.child ADD CONSTRAINT c2 FOREIGN KEY (id) REFERENCES parent (id);
< ok 0 0
> exec
| CREATE TABLE q.orphan (id int, FOREIGN KEY (id) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'

== TestWalkThroughAgainstTiDB/create_a_table_that_exists default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/create_a_table_that_exists default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/create_a_table_that_exists default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE t (a int);
< error
| 1050 42S01 Table 'shop.t' already exists
> exec
| CREATE TABLE v (a int);
< error
| 1050 42S01 Table 'shop.v' already exists
> exec
| CREATE TABLE sq (a int);
< error
| 1050 42S01 Table 'shop.sq' already exists
> exec
| CREATE TABLE IF NOT EXISTS t (a int);
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS v (a int);
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS sq (a int, a int);
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS v LIKE nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| CREATE TABLE x LIKE nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| CREATE TABLE x LIKE v;
< error
| 1347 HY000 'shop.v' is not BASE TABLE
> exec
| CREATE TABLE x LIKE sq;
< error
| 1347 HY000 'shop.sq' is not BASE TABLE

== TestWalkThroughAgainstTiDB/databases default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/databases default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/databases default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE DATABASE shop;
< error
| 1007 HY000 Can't create database 'shop'; database exists
> exec
| CREATE DATABASE q;
< ok 0 0
> exec ^
< error
| 1007 HY000 Can't create database 'q'; database exists
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 0 0
> exec
| DROP DATABASE IF EXISTS nosuch;
< ok 0 0
> exec
| USE nosuch;
< error
| 1049 42000 Unknown database 'nosuch'
> exec
| CREATE TABLE q.x (a int);
< ok 0 0
> exec ^
< error
| 1050 42S01 Table 'q.x' already exists
> exec
| DROP DATABASE q;
< ok 0 0
> exec
| CREATE TABLE q.y (a int);
< error
| 1049 42000 Unknown database 'q'

== TestWalkThroughAgainstTiDB/foreign_keys default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/foreign_keys default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	t	fk1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 33
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 33
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	y00	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 33
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	y00	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	nokey	FK1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fk_1	
| index	shop	y00	fk1	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		y00	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	FK1	
| table foreign key	shop	y	fk_1	
| table foreign key	shop	y00	fk1	

== TestWalkThroughAgainstTiDB/foreign_keys default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE y0 (a int, FOREIGN KEY (a) REFERENCES v (id));
< error
| 3780 HY000 Referencing column 'a' and referenced column 'id' in foreign key constraint 'fk_1' are incompatible.
> exec
| CREATE TABLE y00 (a int, CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id));
< ok 0 0
> exec
| ALTER TABLE nokey DROP FOREIGN KEY fk1;
< error
| 1091 42000 Can't DROP 'fk1'; check that column/key exists
> exec
| ALTER TABLE nokey ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);
< ok 0 0
> exec
| ALTER TABLE nokey ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id);
< error
| 1826 HY000 Duplicate foreign key constraint name 'fk1'
> exec
| ALTER TABLE t ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);
< error
| 1826 HY000 Duplicate foreign key constraint name 'FK1'
> exec
| ALTER TABLE t DROP FOREIGN KEY nosuch;
< error
| 1091 42000 Can't DROP 'nosuch'; check that column/key exists
> exec
| ALTER TABLE child DROP FOREIGN KEY fk_1, ADD CONSTRAINT fk_1 FOREIGN KEY (b) REFERENCES r (id);
< error
| 1826 HY000 Duplicate foreign key constraint name 'fk_1'
> exec
| ALTER TABLE t DROP FOREIGN KEY fk1;
< ok 0 0
> exec ^
< error
| 1091 42000 Can't DROP 'fk1'; check that column/key exists
> exec
| CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'
> exec
| CREATE TABLE ys (a int, FOREIGN KEY (a) REFERENCES sq (id));
< error
| 3734 HY000 Failed to add the foreign key constraint. Missing column 'id' for constraint 'fk_1' in the referenced table 'sq'
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< ok 0 0
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec
| CREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'
> exec
| CREATE TABLE w (id int PRIMARY KEY, p int, FOREIGN KEY (p) REFERENCES w (id));
< ok 0 0

== TestWalkThroughAgainstTiDB/foreign_keys_named_twice default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/foreign_keys_named_twice default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/foreign_keys_named_twice default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a int, b int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT f FOREIGN KEY (b) REFERENCES r (id));
< error
| 1826 HY000 Duplicate foreign key constraint name 'f'
> exec
| CREATE TABLE x (a int, b int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT F FOREIGN KEY (a) REFERENCES r (id));
< error
| 1826 HY000 Duplicate foreign key constraint name 'F'
> exec
| CREATE TABLE x (a int, b int);
< ok 0 0

== TestWalkThroughAgainstTiDB/generated_foreign_key_names default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/generated_foreign_key_names default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	g	b	
| column	shop	g	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	g	PRIMARY	
| index	shop	g	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		g	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	g	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	h	b	
| column	shop	h	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	h	PRIMARY	
| index	shop	h	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		h	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	h	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	h	b	
| column	shop	h	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	h	PRIMARY	
| index	shop	h	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		h	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	h	b	
| column	shop	h	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	h	PRIMARY	
| index	shop	h	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		h	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 31
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	h	b	
| column	shop	h	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	h	PRIMARY	
| index	shop	h	fk_1	
| index	shop	nokey	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		h	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	nokey	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 31
: bytes
| column	shop	h	b	
| column	shop	h	id	
| column	shop	kid	b	
| column	shop	kid	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	h	PRIMARY	
| index	shop	h	fk_1	
| index	shop	kid	PRIMARY	
| index	shop	kid	fk_1	
| index	shop	nokey	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		h	BASE TABLE
| relation	shop		kid	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	kid	fk_1	
| table foreign key	shop	nokey	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	h	b	
| column	shop	h	id	
| column	shop	kid	b	
| column	shop	kid	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	h	PRIMARY	
| index	shop	h	fk_1	
| index	shop	kid	PRIMARY	
| index	shop	kid	fk_1	
| index	shop	nokey	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		h	BASE TABLE
| relation	shop		kid	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	nokey	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	h	b	
| column	shop	h	id	
| column	shop	kid	b	
| column	shop	kid	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	h	PRIMARY	
| index	shop	h	fk_1	
| index	shop	kid	PRIMARY	
| index	shop	kid	fk_1	
| index	shop	nokey	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		h	BASE TABLE
| relation	shop		kid	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	nokey	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	c2	b	
| column	shop	c2	id	
| column	shop	h	b	
| column	shop	h	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	c2	PRIMARY	
| index	shop	c2	fk_1	
| index	shop	h	PRIMARY	
| index	shop	h	fk_1	
| index	shop	nokey	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		c2	BASE TABLE
| relation	shop		h	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	nokey	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 31
: bytes
| column	shop	c2	b	
| column	shop	c2	id	
| column	shop	h	b	
| column	shop	h	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	c2	PRIMARY	
| index	shop	c2	fk_1	
| index	shop	h	PRIMARY	
| index	shop	h	fk_1	
| index	shop	nokey	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		c2	BASE TABLE
| relation	shop		h	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	c2	fk_1	
| table foreign key	shop	nokey	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 31
: bytes
| column	shop	c2	b	
| column	shop	c2	id	
| column	shop	h	b	
| column	shop	h	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	c2	PRIMARY	
| index	shop	c2	fk_1	
| index	shop	h	PRIMARY	
| index	shop	h	fk_1	
| index	shop	nokey	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		c2	BASE TABLE
| relation	shop		h	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	c2	fk_1	
| table foreign key	shop	nokey	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/generated_foreign_key_names default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE g (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| RENAME TABLE g TO h;
< ok 0 0
> exec
| ALTER TABLE h DROP FOREIGN KEY fk_1;
< ok 0 0
> exec ^
< error
| 1091 42000 Can't DROP 'fk_1'; check that column/key exists
> exec
| ALTER TABLE nokey ADD CONSTRAINT fk_1 FOREIGN KEY (a) REFERENCES r (id);
< ok 0 0
> exec
| RENAME TABLE child TO kid;
< ok 0 0
> exec
| ALTER TABLE kid DROP FOREIGN KEY fk_1;
< ok 0 0
> exec ^
< error
| 1091 42000 Can't DROP 'fk_1'; check that column/key exists
> exec
| ALTER TABLE kid RENAME TO c2;
< ok 0 0
> exec
| ALTER TABLE c2 ADD CONSTRAINT fk_1 FOREIGN KEY (b) REFERENCES r (id);
< ok 0 0
> exec ^
< error
| 1826 HY000 Duplicate foreign key constraint name 'fk_1'

== TestWalkThroughAgainstTiDB/name_case default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/name_case default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	Q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	Q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	Q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	Q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	Q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	Q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	Q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	Q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	Q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	Q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	Q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	Q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	Q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	Q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	Q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	Q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	Q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	Q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	Q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	Q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	Q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	Q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	Q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	Q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	Q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	Q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	Q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	Q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	Q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	

== TestWalkThroughAgainstTiDB/name_case default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE T (a int);
< error
| 1050 42S01 Table 'shop.T' already exists
> exec
| DROP TABLE NOKEY;
< ok 0 0
> exec
| CREATE DATABASE Q;
< ok 0 0
> exec
| CREATE DATABASE q;
< error
| 1007 HY000 Can't create database 'q'; database exists
> exec
| CREATE TABLE Q.x (a int);
< ok 0 0
> exec
| CREATE TABLE q.X (a int);
< error
| 1050 42S01 Table 'Q.X' already exists
> exec
| ALTER TABLE R ADD COLUMN ID int;
< error
| 1060 42S21 Duplicate column name 'ID'
> exec
| CREATE VIEW V AS SELECT 1 AS a;
< error
| 1050 42S01 Table 'shop.V' already exists
> exec
| CREATE SEQUENCE SQ;
< error
| 1050 42S01 Table 'shop.SQ' already exists
> exec
| DROP SEQUENCE Sq;
< ok 0 0
> exec
| DROP SEQUENCE sQ;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.sQ'
> exec
| ALTER TABLE T DROP FOREIGN KEY FK1;
< ok 0 0
> exec
| ALTER TABLE t DROP FOREIGN KEY fk1;
< error
| 1091 42000 Can't DROP 'fk1'; check that column/key exists

== TestWalkThroughAgainstTiDB/names_outside_ASCII default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| CREATE TABLE kel (`K` int, a int);
< ok 0 0
> exec
| CREATE TABLE `tK` (a int);
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/names_outside_ASCII default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| kel
: bytes
| K
< row
: bytes
| kel
: bytes
| a
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| tK
: bytes
| a
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| kel
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| tK
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	kel	K	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	tK	a	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		tK	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	tK	a	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		tK	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	tK	a	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		tK	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	café	
| column	shop	x	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	ké	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	ké	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/names_outside_ASCII default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE kel DROP COLUMN k;
< ok 0 0
> exec
| ALTER TABLE kel DROP COLUMN zz;
< error
| 1091 42000 Can't DROP 'zz'; check that column/key exists
> exec
| DROP TABLE tk;
< ok 0 0
> exec
| CREATE TABLE x (id int, café int, KEY ké (id));
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN CAFé;
< ok 0 0
> exec
| ALTER TABLE x DROP INDEX KÉ;
< ok 0 0
> exec
| ALTER TABLE x DROP INDEX Ké;
< error
| 1091 42000 index Ké doesn't exist

== TestWalkThroughAgainstTiDB/names_too_long default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/names_too_long default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/names_too_long default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa (a int);
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec
| ALTER TABLE t ADD COLUMN aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa int;
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec
| CREATE INDEX aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa ON t (a);
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec
| RENAME TABLE nokey TO aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa;
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec
| ALTER TABLE nokey RENAME TO aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa;
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec
| CREATE SEQUENCE aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa;
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec
| CREATE VIEW aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa AS SELECT 1 AS a;
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec
| DROP TABLE IF EXISTS aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa;
< ok 0 0
> exec
| DROP TABLE aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa;
< error
| 1051 42S02 Unknown table 'shop.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
> exec
| ALTER TABLE t CHANGE a aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa int;
< ok 0 0
> exec
| ALTER TABLE t RENAME COLUMN b TO baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa;
< ok 0 0

== TestWalkThroughAgainstTiDB/no_database_selected default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/no_database_selected default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	test			

== TestWalkThroughAgainstTiDB/no_database_selected default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`tx`
< error
| 1146 42S02 Table 'shop.tx' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tx`
< error
| 1146 42S02 Table 'test.tx' doesn't exist
> exec
| DROP DATABASE shop;
< ok 0 0
> query
| SHOW CREATE TABLE `test`.`tx`
< error
| 1146 42S02 Table 'test.tx' doesn't exist
> exec
| ALTER DATABASE CHARACTER SET utf8mb4;
< error
| 1046 3D000 No database selected
> query
| SHOW CREATE TABLE `test`.`tx`
< error
| 1146 42S02 Table 'test.tx' doesn't exist
> exec
| CREATE TABLE x (a int);
< error
| 1046 3D000 No database selected
> query
| SHOW CREATE TABLE `test`.`tx`
< error
| 1146 42S02 Table 'test.tx' doesn't exist
> exec
| CREATE TEMPORARY TABLE tx (a int);
< error
| 1046 3D000 No database selected
> query
| SHOW CREATE TABLE `test`.`tx`
< error
| 1146 42S02 Table 'test.tx' doesn't exist
> exec
| DROP TABLE IF EXISTS tx;
< error
| 1046 3D000 No database selected
> query
| SHOW CREATE TABLE `test`.`tx`
< error
| 1146 42S02 Table 'test.tx' doesn't exist
> exec
| CREATE VIEW vx AS SELECT 1 AS a;
< error
| 1046 3D000 No database selected
> query
| SHOW CREATE TABLE `test`.`tx`
< error
| 1146 42S02 Table 'test.tx' doesn't exist
> exec
| CREATE SEQUENCE sx;
< error
| 1046 3D000 No database selected
> query
| SHOW CREATE TABLE `test`.`tx`
< error
| 1146 42S02 Table 'test.tx' doesn't exist
> exec
| DROP SEQUENCE IF EXISTS sx;
< error
| 1046 3D000 No database selected
> query
| SHOW CREATE TABLE `test`.`tx`
< error
| 1146 42S02 Table 'test.tx' doesn't exist

== TestWalkThroughAgainstTiDB/objects_in_a_dropped_database default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/objects_in_a_dropped_database default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| database	z			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/objects_in_a_dropped_database default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE DATABASE z;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> query
| SHOW CREATE TABLE `z`.`tmp`
< error
| 1146 42S02 Table 'z.tmp' doesn't exist
> exec
| DROP DATABASE z;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE VIEW z.w AS SELECT 1 AS a;
< error
| 1049 42000 Unknown database 'z'
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE TEMPORARY TABLE z.tmp (a int);
< error
| 1049 42000 Unknown database 'z'
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE SEQUENCE z.s;
< error
| 1049 42000 Unknown database 'z'
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE TABLE z.t (a int);
< error
| 1049 42000 Unknown database 'z'
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist

== TestWalkThroughAgainstTiDB/repeated_subcommands default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/repeated_subcommands default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/repeated_subcommands default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN c int, ADD COLUMN c int;
< error
| 8200 HY000 Unsupported operate same column 'c'
> exec
| ALTER TABLE t DROP COLUMN a, DROP COLUMN a;
< error
| 8200 HY000 Unsupported operate same column 'a'
> exec
| ALTER TABLE t ADD INDEX k (a), ADD INDEX k (b);
< error
| 8200 HY000 Unsupported operate same index 'k'
> exec
| ALTER TABLE nokey ADD PRIMARY KEY (a), ADD PRIMARY KEY (a);
< error
| 8200 HY000 Unsupported operate same index 'primary'
> exec
| ALTER TABLE t DROP INDEX i, DROP INDEX I;
< error
| 8200 HY000 Unsupported operate same index 'i'
> exec
| ALTER TABLE t DROP PRIMARY KEY, DROP PRIMARY KEY;
< error
| 8200 HY000 Unsupported operate same index 'primary'
> exec
| ALTER TABLE t DROP PRIMARY KEY, DROP INDEX `PRIMARY`;
< error
| 8200 HY000 Unsupported operate same index 'primary'
> exec
| ALTER TABLE t CHANGE a z int, CHANGE b z int;
< error
| 8200 HY000 Unsupported operate same column 'z'
> exec
| ALTER TABLE t ADD COLUMN z int, RENAME COLUMN a TO z;
< error
| 8200 HY000 Unsupported operate same column 'z'
> exec
| ALTER TABLE t DROP COLUMN IF EXISTS zz, DROP COLUMN IF EXISTS zz;
< ok 0 0

== TestWalkThroughAgainstTiDB/sequences default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/sequences default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		s2	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		s3	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		s3	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 23
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/sequences default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE SEQUENCE `sq` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE SEQUENCE sq;
< error
| 1050 42S01 Table 'shop.sq' already exists
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE SEQUENCE `sq` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE SEQUENCE SQ;
< error
| 1050 42S01 Table 'shop.SQ' already exists
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE SEQUENCE `sq` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE SEQUENCE t;
< error
| 1050 42S01 Table 'shop.t' already exists
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE SEQUENCE `sq` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE SEQUENCE v;
< error
| 1050 42S01 Table 'shop.v' already exists
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE SEQUENCE `sq` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE SEQUENCE IF NOT EXISTS sq;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE SEQUENCE `sq` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE VIEW sq AS SELECT 1 AS a;
< error
| 1050 42S01 Table 'shop.sq' already exists
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE SEQUENCE `sq` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE OR REPLACE VIEW sq AS SELECT 1 AS a;
< error
| 1347 HY000 'shop.sq' is not VIEW
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE SEQUENCE `sq` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| DROP SEQUENCE t;
< error
| 1347 HY000 'shop.t' is not SEQUENCE
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE SEQUENCE `sq` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| DROP SEQUENCE v;
< error
| 1347 HY000 'shop.v' is not SEQUENCE
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE SEQUENCE `sq` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| DROP SEQUENCE nosuch;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.nosuch'
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE SEQUENCE `sq` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| DROP SEQUENCE IF EXISTS nosuch, t;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE SEQUENCE `sq` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE TEMPORARY TABLE sq (a int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE TEMPORARY TABLE `sq` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| DROP SEQUENCE sq;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE TEMPORARY TABLE `sq` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| DROP SEQUENCE sq;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.sq'
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE TEMPORARY TABLE `sq` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE SEQUENCE s2;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE TEMPORARY TABLE `sq` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| RENAME TABLE s2 TO s3;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE TEMPORARY TABLE `sq` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| DROP SEQUENCE s2;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.s2'
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE TEMPORARY TABLE `sq` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| DROP SEQUENCE s3;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE TEMPORARY TABLE `sq` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE TEMPORARY TABLE tt (a int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE TEMPORARY TABLE `sq` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tt
: bytes
| CREATE TEMPORARY TABLE `tt` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| DROP SEQUENCE tt;
< error
| 4139 HY000 Unknown SEQUENCE: 'shop.tt'
> query
| SHOW CREATE TABLE `shop`.`sq`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| sq
: bytes
| CREATE TEMPORARY TABLE `sq` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tt
: bytes
| CREATE TEMPORARY TABLE `tt` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`sq`
< error
| 1146 42S02 Table 'test.sq' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist

== TestWalkThroughAgainstTiDB/standalone_admin_statements default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/standalone_admin_statements default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/standalone_admin_statements default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| RECOMMEND INDEX SHOW;
< ok 0 0
> exec
| RECOMMEND INDEX RUN FOR 'select * from t where a = 1';
< ok 0 0
> exec
| SET CONFIG tikv nosuch.item = 1;
< ok 0 0
> exec
| ALTER RANGE global PLACEMENT POLICY DEFAULT;
< ok 0 0
> exec
| ALTER RANGE meta PLACEMENT POLICY SET DEFAULT;
< ok 0 0
> exec
| KILL 999999;
< ok 0 0
> exec
| KILL QUERY 999999;
< ok 0 0
> exec
| KILL TIDB 999999;
< ok 0 0
> exec
| KILL TIDB CONNECTION 999999;
< ok 0 0
> exec
| CANCEL BR JOB 1;
< ok 0 0
> exec
| SHOW BR JOB QUERY 1;
< ok 0 0
> exec
| SHOW BACKUP METADATA FROM 'noop://';
< ok 0 0
> exec
| FLUSH PRIVILEGES;
< ok 0 0
> exec
| DROP TABLE nosuch;
< error
| 1051 42S02 Unknown table 'shop.nosuch'

== TestWalkThroughAgainstTiDB/temporary_tables default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/temporary_tables default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/temporary_tables default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */,
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE TEMPORARY TABLE t (z int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TEMPORARY TABLE `t` (
|   `z` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| ALTER TABLE nokey ADD COLUMN z int;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TEMPORARY TABLE `t` (
|   `z` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| INSERT INTO t (z) VALUES (1);
< ok 1 0
> query
| SHOW CREATE TABLE `shop`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TEMPORARY TABLE `t` (
|   `z` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| DROP TABLE t;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */,
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| INSERT INTO t (id) VALUES (1);
< ok 1 0
> query
| SHOW CREATE TABLE `shop`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */,
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| DROP TEMPORARY TABLE t;
< error
| 1051 42S02 Unknown table 'shop.t'
> query
| SHOW CREATE TABLE `shop`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */,
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS t;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */,
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< error
| 1146 42S02 Table 'shop.tt' doesn't exist
> query
| SHOW CREATE TABLE `test`.`t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE TEMPORARY TABLE tt (a int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */,
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tt
: bytes
| CREATE TEMPORARY TABLE `tt` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| CREATE TEMPORARY TABLE tt (a int);
< error
| 1050 42S01 Table 'shop.tt' already exists
> query
| SHOW CREATE TABLE `shop`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */,
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tt
: bytes
| CREATE TEMPORARY TABLE `tt` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| TRUNCATE TABLE tt;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */,
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tt`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tt
: bytes
| CREATE TEMPORARY TABLE `tt` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `test`.`t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tt`
< error
| 1146 42S02 Table 'test.tt' doesn't exist
> exec
| INSERT INTO tt VALUES (1);
< ok 1 0

== TestWalkThroughAgainstTiDB/temporary_tables_outlive_their_database default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/temporary_tables_outlive_their_database default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/temporary_tables_outlive_their_database default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE DATABASE q;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`tmp`
< error
| 1146 42S02 Table 'q.tmp' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE TEMPORARY TABLE q.tmp (a int);
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`tmp`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tmp
: bytes
| CREATE TEMPORARY TABLE `tmp` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| DROP DATABASE q;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE DATABASE q;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`tmp`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tmp
: bytes
| CREATE TEMPORARY TABLE `tmp` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| DROP TEMPORARY TABLE q.tmp;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`tmp`
< error
| 1146 42S02 Table 'q.tmp' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| DROP TEMPORARY TABLE q.tmp;
< error
| 1051 42S02 Unknown table 'q.tmp'
> query
| SHOW CREATE TABLE `q`.`tmp`
< error
| 1146 42S02 Table 'q.tmp' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`tmp`
< error
| 1146 42S02 Table 'shop.tmp' doesn't exist
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist

== TestWalkThroughAgainstTiDB/unknown_by_design default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/unknown_by_design default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/unknown_by_design default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP DATABASE nosuch;
< error
| 1008 HY000 Can't drop database 'nosuch'; database doesn't exist
> exec
| RENAME TABLE t TO nodb.t;
< error
| 1025 HY000 Error on rename of 'shop.t' to 'nodb.t' (errno: 168 - Database `nodb` doesn't exist)
> exec
| CREATE VIEW nodb.w AS SELECT 1 AS a;
< error
| 1049 42000 Unknown database 'nodb'
> exec
| CREATE TABLE x AS SELECT 1 AS a;
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| CREATE TABLE x (a int, b int UNIQUE);
< ok 0 0
> exec
| ALTER TABLE x DROP INDEX zz;
< error
| 1091 42000 index zz doesn't exist
> exec
| ALTER TABLE nokey ADD COLUMN _tidb_rowid int;
< error
| 1166 42000 Incorrect column name '_tidb_rowid'
> exec
| ALTER TABLE nokey ADD COLUMN y int AFTER _tidb_rowid;
< error
| 1054 42S22 Unknown column '_tidb_rowid' in 'nokey'
> exec
| ALTER TABLE nokey DROP COLUMN _tidb_rowid;
< error
| 1091 42000 Can't DROP '_tidb_rowid'; check that column/key exists
> exec
| ALTER TABLE t ADD COLUMN q int, ADD INDEX kq (q);
< error
| 1072 42000 column does not exist: q
> exec
| ALTER TABLE t DROP COLUMN IF EXISTS a, DROP COLUMN a;
< error
| 8200 HY000 Unsupported operate same column 'a'
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1, MODIFY a int;
< error
| 8200 HY000 Unsupported operate same column 'a'
> exec
| ALTER TABLE r DROP PRIMARY KEY;
< error
| 8200 HY000 Unsupported drop primary key when the table is using clustered index
> exec
| ALTER TABLE t DROP CONSTRAINT nosuch;
< error
| 3940 HY000 Constraint 'nosuch' does not exist.
> exec
| INSERT INTO v VALUES (1);
< error
| 1105 HY000 insert into view v is not supported now
> exec
| CREATE TABLE c3 (id int, a int, FOREIGN KEY (id) REFERENCES r (id), CONSTRAINT fk_1 FOREIGN KEY (a) REFERENCES r (id));
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE PROCEDURE p() SELECT 1;
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec
| CREATE TABLE t (a int);
< error
| 1050 42S01 Table 'shop.t' already exists

== TestWalkThroughAgainstTiDB/vector_columns default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/vector_columns default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 26
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	c	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	c	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	c	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	c	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/vector_columns default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a VECTOR(3));
< ok 0 0
> exec
| CREATE TABLE x (a VECTOR);
< error
| 1050 42S01 Table 'shop.x' already exists
> exec
| ALTER TABLE x ADD COLUMN b VECTOR<FLOAT>(3);
< ok 0 0
> exec
| ALTER TABLE x ADD COLUMN b VECTOR;
< error
| 1060 42S21 Duplicate column name 'b'
> exec
| ALTER TABLE x MODIFY nosuch VECTOR(3);
< error
| 1054 42S22 Unknown column 'nosuch' in 'x'
> exec
| ALTER TABLE x CHANGE b c VECTOR(4);
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN b;
< error
| 1091 42000 Can't DROP 'b'; check that column/key exists
> exec
| CREATE VECTOR INDEX i ON nosuch ((VEC_COSINE_DISTANCE(a)));
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| ALTER TABLE nosuch ADD VECTOR INDEX i ((VEC_COSINE_DISTANCE(a)));
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist

== TestWalkThroughAgainstTiDB/views default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE SEQUENCE sq;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1
< ok 0 0

== TestWalkThroughAgainstTiDB/views default#2
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
| r
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
< row
: bytes
| t
: bytes
| i
: int64 1
< row
: bytes
| t
: bytes
| fk1
: int64 1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
< row
: bytes
| child
: bytes
| fk_1
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
| child
: bytes
| id
< row
: bytes
| child
: bytes
| b
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| r
: bytes
| id
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| a
< row
: bytes
| t
: bytes
| b
< row
: bytes
| v
: bytes
| id
> query
| SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ?
: string
| shop
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| v
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
< row
: bytes
| child
: bytes
| fk_1
: bytes
| r
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
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
| sq
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
| child
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| sq
: bytes
| SEQUENCE
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| v
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
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	id	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 24
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	v	a	
| database	shop			
| database	test			
| index	shop	child	PRIMARY	
| index	shop	child	fk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		sq	SEQUENCE
| relation	shop		t	BASE TABLE
| relation	shop		v	VIEW
| table foreign key	shop	child	fk_1	
| table foreign key	shop	t	fk1	

== TestWalkThroughAgainstTiDB/views default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE VIEW t AS SELECT 1 AS a;
< error
| 1050 42S01 Table 'shop.t' already exists
> exec
| CREATE VIEW v AS SELECT 1 AS a;
< error
| 1050 42S01 Table 'shop.v' already exists
> exec
| CREATE OR REPLACE VIEW t AS SELECT 1 AS a;
< error
| 1347 HY000 'shop.t' is not VIEW
> exec
| CREATE OR REPLACE VIEW v AS SELECT 1 AS a;
< ok 0 0
> exec
| DROP VIEW IF EXISTS t;
< error
| 1347 HY000 'shop.t' is not VIEW
> exec
| DROP VIEW t;
< error
| 1347 HY000 'shop.t' is not VIEW
