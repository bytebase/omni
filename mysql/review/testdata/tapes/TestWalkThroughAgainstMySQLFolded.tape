-- pin: lctn1 mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestWalkThroughAgainstMySQLFolded/ALTER_TABLE_drops_first lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/ALTER_TABLE_drops_first lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	b	
| column	shop	t	c	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	b	
| column	shop	t	c	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/ALTER_TABLE_drops_first lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN a int, DROP COLUMN a;
< ok 0 0
> exec
| ALTER TABLE t DROP INDEX i, ADD INDEX i (b);
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN a;
< ok 0 0
> exec ^
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| ALTER TABLE t RENAME COLUMN b TO c, ADD COLUMN b int;
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN zz, ADD COLUMN y int;
< error
| 1091 42000 Can't DROP 'zz'; check that column/key exists

== TestWalkThroughAgainstMySQLFolded/ALTER_TABLE_subcommands lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/ALTER_TABLE_subcommands lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	PRIMARY2	
| index	shop	t	fk1	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	PRIMARY2	
| index	shop	t	fk1	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/ALTER_TABLE_subcommands lctn1#3
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
| 1091 42000 Can't DROP 'nosuch'; check that column/key exists
> exec
| ALTER TABLE t ADD PRIMARY KEY (a);
< error
| 1068 42000 Multiple primary key defined
> exec
| ALTER TABLE t RENAME INDEX nosuch TO j;
< error
| 1176 42000 Key 'nosuch' doesn't exist in table 't'
> exec
| ALTER TABLE t RENAME INDEX i TO PRIMARY2;
< ok 0 0
> exec
| ALTER TABLE t ALTER INDEX nosuch INVISIBLE;
< error
| 1176 42000 Key 'nosuch' doesn't exist in table 't'

== TestWalkThroughAgainstMySQLFolded/CREATE_and_DROP_INDEX lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/CREATE_and_DROP_INDEX lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/CREATE_and_DROP_INDEX lctn1#3
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
| CREATE INDEX k ON v (id);
< error
| 1347 HY000 'shop.v' is not BASE TABLE
> exec
| CREATE INDEX i ON t (b);
< error
| 1061 42000 Duplicate key name 'i'
> exec
| CREATE INDEX `PRIMARY` ON t (b);
< error
| 1280 42000 Incorrect index name 'PRIMARY'
> exec
| DROP INDEX nosuch ON t;
< error
| 1091 42000 Can't DROP 'nosuch'; check that column/key exists
> exec
| DROP INDEX i ON nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| DROP INDEX i ON t;
< ok 0 0
> exec ^
< error
| 1091 42000 Can't DROP 'i'; check that column/key exists

== TestWalkThroughAgainstMySQLFolded/DML_targets lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/DML_targets lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/DML_targets lctn1#3
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
| INSERT INTO v VALUES (1);
< ok 1 0
> exec
| INSERT INTO nokey VALUES (1);
< ok 1 0

== TestWalkThroughAgainstMySQLFolded/DROP_FUNCTION_with_no_database_selected lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/DROP_FUNCTION_with_no_database_selected lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			

== TestWalkThroughAgainstMySQLFolded/DROP_FUNCTION_with_no_database_selected lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP DATABASE shop;
< ok 6 0
> exec
| DROP FUNCTION IF EXISTS nosuch;
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS f;
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p;
< error
| 1046 3D000 No database selected
> exec
| DROP FUNCTION IF EXISTS shop.f;
< ok 0 0
> exec
| DROP VIEW IF EXISTS v;
< error
| 1046 3D000 No database selected
> exec
| DROP EVENT IF EXISTS e;
< error
| 1046 3D000 No database selected
> exec
| DROP TABLE IF EXISTS nokey;
< error
| 1046 3D000 No database selected

== TestWalkThroughAgainstMySQLFolded/DROP_FUNCTION_with_no_database_selected_may_find_a_loadable_one lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/DROP_FUNCTION_with_no_database_selected_may_find_a_loadable_one lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			

== TestWalkThroughAgainstMySQLFolded/DROP_FUNCTION_with_no_database_selected_may_find_a_loadable_one lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP DATABASE shop;
< ok 6 0
> exec
| DROP FUNCTION nosuch;
< error
| 1305 42000 FUNCTION (UDF) nosuch does not exist
> exec
| CREATE FUNCTION service_get_read_locks RETURNS INT SONAME 'locking_service.so';
< ok 0 0
> exec
| DROP FUNCTION service_get_read_locks;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/DROP_TABLE lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/DROP_TABLE lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 25
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/DROP_TABLE lctn1#3
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
> exec
| DROP TABLE v;
< error
| 1051 42S02 Unknown table 'shop.v'
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

== TestWalkThroughAgainstMySQLFolded/DROP_TABLE_IF_EXISTS_leaves_a_view lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| CREATE VIEW q.v AS SELECT 1 AS a;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/DROP_TABLE_IF_EXISTS_leaves_a_view lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	q	v	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		v	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	q	v	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		v	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/DROP_TABLE_IF_EXISTS_leaves_a_view lctn1#3
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

== TestWalkThroughAgainstMySQLFolded/RENAME lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/RENAME lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	r	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t	PRIMARY	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	r	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t	PRIMARY	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	r	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t	PRIMARY	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	r	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t2	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 29
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	r	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t2	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	r	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t2	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	q	nokey	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	r	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| procedure	shop		p	
| relation	q		nokey	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t2	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	q	nokey	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	a	
| column	shop	r	b	
| column	shop	r	id	
| column	shop	t2	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	r	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	r	fk1	
| index	shop	r	i	
| index	shop	t2	PRIMARY	
| procedure	shop		p	
| relation	q		nokey	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t2	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/RENAME lctn1#3
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
| 1050 42S01 Table 'r' already exists
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
| 1050 42S01 Table 'r' already exists
> exec
| ALTER TABLE t RENAME TO t2;
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN zz int;
< error
| 1146 42S02 Table 'shop.t' doesn't exist
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| RENAME TABLE nokey TO q.nokey;
< ok 0 0
> exec
| CREATE TABLE q.nokey (a int);
< error
| 1050 42S01 Table 'nokey' already exists

== TestWalkThroughAgainstMySQLFolded/SET_foreign_key_checks_in_each_assignment's_scope lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/SET_foreign_key_checks_in_each_assignment's_scope lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z2	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| foreign key	shop	z2	z2_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| index	shop	z2	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z2	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z2	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| foreign key	shop	z2	z2_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| index	shop	z2	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z2	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 41
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z2	a	
| column	shop	z3	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| foreign key	shop	z2	z2_ibfk_1	
| foreign key	shop	z3	z3_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| index	shop	z2	a	
| index	shop	z3	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z2	BASE TABLE
| relation	shop		z3	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 41
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z2	a	
| column	shop	z3	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| foreign key	shop	z2	z2_ibfk_1	
| foreign key	shop	z3	z3_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| index	shop	z2	a	
| index	shop	z3	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z2	BASE TABLE
| relation	shop		z3	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 41
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z2	a	
| column	shop	z3	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| foreign key	shop	z2	z2_ibfk_1	
| foreign key	shop	z3	z3_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| index	shop	z2	a	
| index	shop	z3	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z2	BASE TABLE
| relation	shop		z3	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 41
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z2	a	
| column	shop	z3	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| foreign key	shop	z2	z2_ibfk_1	
| foreign key	shop	z3	z3_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| index	shop	z2	a	
| index	shop	z3	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z2	BASE TABLE
| relation	shop		z3	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 41
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z2	a	
| column	shop	z3	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| foreign key	shop	z2	z2_ibfk_1	
| foreign key	shop	z3	z3_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| index	shop	z2	a	
| index	shop	z3	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z2	BASE TABLE
| relation	shop		z3	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 41
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| column	shop	z2	a	
| column	shop	z3	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| foreign key	shop	z2	z2_ibfk_1	
| foreign key	shop	z3	z3_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| index	shop	z2	a	
| index	shop	z3	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| relation	shop		z2	BASE TABLE
| relation	shop		z3	BASE TABLE
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/SET_foreign_key_checks_in_each_assignment's_scope lctn1#3
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
| SET GLOBAL foreign_key_checks = 1, @@foreign_key_checks = 1;
< ok 0 0
> exec
| CREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'
> exec
| SET @@global.foreign_key_checks = 1, foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE z2 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1;
< ok 0 0
> exec
| CREATE TABLE z3 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< ok 0 0
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE z4 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'
> exec
| SET GLOBAL foreign_key_checks = 1, foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE w (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'

== TestWalkThroughAgainstMySQLFolded/SET_foreign_key_checks_under_a_later_scope lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/SET_foreign_key_checks_under_a_later_scope lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	z	z_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	z	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		z	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	z	z_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	z	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		z	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	z	z_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	z	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		z	BASE TABLE
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/SET_foreign_key_checks_under_a_later_scope lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| SET @a = 1, GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'
> exec
| SET GLOBAL foreign_key_checks = 1, SESSION foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< ok 0 0
> exec
| SET SESSION foreign_key_checks = 1, GLOBAL foreign_key_checks = 1, foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE w (a int, FOREIGN KEY (a) REFERENCES nosuch (id));
< error
| 1824 HY000 Failed to open the referenced table 'nosuch'

== TestWalkThroughAgainstMySQLFolded/USE_moves_unqualified_names lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/USE_moves_unqualified_names lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	q	t	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	q	t	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	q	t	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	q	t	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/USE_moves_unqualified_names lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| USE q;
< ok 0 0
> exec
| CREATE TABLE t (a int);
< ok 0 0
> exec ^
< error
| 1050 42S01 Table 't' already exists
> exec
| USE shop;
< ok 0 0
> exec
| CREATE TABLE t (a int);
< error
| 1050 42S01 Table 't' already exists

== TestWalkThroughAgainstMySQLFolded/a_foreign_key's_name_and_its_index's lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_foreign_key's_name_and_its_index's lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	x	cx	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	cx	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	cx	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	cx	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	y_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 34
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	z	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	z	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	nokey	nk	
| foreign key	shop	t	fk1	
| foreign key	shop	z	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	nokey	nk	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	z	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 41
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	nokey	nk	
| foreign key	shop	nokey	nk2	
| foreign key	shop	nokey	nk3	
| foreign key	shop	t	fk1	
| foreign key	shop	z	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	nokey	nk3	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	z	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 45
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	w	a	
| column	shop	x	a	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	nokey	nk	
| foreign key	shop	nokey	nk2	
| foreign key	shop	nokey	nk3	
| foreign key	shop	t	fk1	
| foreign key	shop	w	cw	
| foreign key	shop	z	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	nokey	nk3	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	w	cw	
| index	shop	z	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		w	BASE TABLE
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 45
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	w	a	
| column	shop	x	a	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	nokey	nk	
| foreign key	shop	nokey	nk2	
| foreign key	shop	nokey	nk3	
| foreign key	shop	t	fk1	
| foreign key	shop	w	cw	
| foreign key	shop	z	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	nokey	nk3	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	w	cw	
| index	shop	z	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		w	BASE TABLE
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 45
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	w	a	
| column	shop	x	a	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	nokey	nk	
| foreign key	shop	nokey	nk2	
| foreign key	shop	nokey	nk3	
| foreign key	shop	t	fk1	
| foreign key	shop	w	cw	
| foreign key	shop	z	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	nokey	nk3	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	w	cw	
| index	shop	z	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		w	BASE TABLE
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/a_foreign_key's_name_and_its_index's lctn1#3
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
| ALTER TABLE y DROP FOREIGN KEY y_ibfk_1;
< ok 0 0
> exec
| ALTER TABLE y DROP INDEX fy;
< ok 0 0
> exec
| CREATE TABLE z (a int, CONSTRAINT fy FOREIGN KEY (a) REFERENCES r (id));
< ok 0 0
> exec
| ALTER TABLE nokey ADD CONSTRAINT nk FOREIGN KEY fk1 (a) REFERENCES r (id);
< ok 0 0
> exec
| ALTER TABLE nokey ADD CONSTRAINT nk2 FOREIGN KEY fk1 (a) REFERENCES r (id), ADD CONSTRAINT nk3 FOREIGN KEY fk1 (a) REFERENCES r (id);
< ok 0 0
> exec
| CREATE TABLE w (a int, CONSTRAINT cw FOREIGN KEY fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff (a) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE w2 (a int, FOREIGN KEY fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff (a) REFERENCES r (id));
< error
| 1059 42000 Identifier name 'fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff' is too long
> exec
| CREATE TABLE w3 (a int, CONSTRAINT fk1 FOREIGN KEY fw (a) REFERENCES r (id));
< error
| 1826 HY000 Duplicate foreign key constraint name 'fk1'

== TestWalkThroughAgainstMySQLFolded/a_function_a_view_or_another_database_holds_may_create_a_temporary_table lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| CREATE FUNCTION q.mk() RETURNS int DETERMINISTIC BEGIN CREATE TEMPORARY TABLE IF NOT EXISTS shop.made (a int); RETURN 1; END;
< ok 0 0
> exec
| CREATE FUNCTION q.mk2() RETURNS int DETERMINISTIC BEGIN CREATE TEMPORARY TABLE IF NOT EXISTS shop.made2 (a int); RETURN 1; END;
< ok 0 0
> exec
| CREATE FUNCTION q.mk3() RETURNS int DETERMINISTIC BEGIN CREATE TEMPORARY TABLE IF NOT EXISTS shop.made3 (a int); RETURN 1; END;
< ok 0 0
> exec
| CREATE FUNCTION q.mk4() RETURNS int DETERMINISTIC BEGIN CREATE TEMPORARY TABLE IF NOT EXISTS shop.made4 (a int); RETURN 1; END;
< ok 0 0
> exec
| CREATE VIEW vv AS SELECT q.mk() AS x;
< ok 0 0
> exec
| CREATE VIEW q.qv AS SELECT q.mk3() AS x;
< ok 0 0
> exec
| INSERT INTO nokey VALUES (1);
< ok 1 0

== TestWalkThroughAgainstMySQLFolded/a_function_a_view_or_another_database_holds_may_create_a_temporary_table lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
< row
: bytes
| vv
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| vv
: bytes
| x
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	q	qv	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	vv	x	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	q		mk	
| function	q		mk2	
| function	q		mk3	
| function	q		mk4	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		qv	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		vv	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	q	qv	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	vv	x	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	q		mk	
| function	q		mk2	
| function	q		mk3	
| function	q		mk4	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		qv	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		vv	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	q	qv	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	vv	x	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	q		mk	
| function	q		mk2	
| function	q		mk3	
| function	q		mk4	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		qv	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		vv	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	q	qv	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	vv	x	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	q		mk	
| function	q		mk2	
| function	q		mk3	
| function	q		mk4	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		qv	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		vv	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	q	qv	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	vv	x	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	q		mk	
| function	q		mk2	
| function	q		mk3	
| function	q		mk4	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		qv	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		vv	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	q	qv	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	vv	x	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	q		mk	
| function	q		mk2	
| function	q		mk3	
| function	q		mk4	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		qv	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		vv	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	q	qv	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	vv	x	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	q		mk	
| function	q		mk2	
| function	q		mk3	
| function	q		mk4	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		qv	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		vv	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	q	qv	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	vv	x	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	q		mk	
| function	q		mk2	
| function	q		mk3	
| function	q		mk4	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		qv	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		vv	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	q	qv	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	vv	x	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	q		mk	
| function	q		mk2	
| function	q		mk3	
| function	q		mk4	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		qv	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		vv	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	q	qv	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	vv	x	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	q		mk	
| function	q		mk2	
| function	q		mk3	
| function	q		mk4	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		qv	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		vv	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	q	qv	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	vv	x	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	q		mk	
| function	q		mk2	
| function	q		mk3	
| function	q		mk4	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		qv	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		vv	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	q	qv	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	vv	x	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	q		mk	
| function	q		mk2	
| function	q		mk3	
| function	q		mk4	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		qv	VIEW
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		vv	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/a_function_a_view_or_another_database_holds_may_create_a_temporary_table lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made2`
< error
| 1146 42S02 Table 'q.made2' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made3`
< error
| 1146 42S02 Table 'q.made3' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made4`
< error
| 1146 42S02 Table 'q.made4' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made2`
< error
| 1146 42S02 Table 'shop.made2' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made3`
< error
| 1146 42S02 Table 'shop.made3' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made4`
< error
| 1146 42S02 Table 'shop.made4' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made2`
< error
| 1146 42S02 Table 'test.made2' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made3`
< error
| 1146 42S02 Table 'test.made3' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made4`
< error
| 1146 42S02 Table 'test.made4' doesn't exist
> exec
| SELECT q.mk4() FROM nokey;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made2`
< error
| 1146 42S02 Table 'q.made2' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made3`
< error
| 1146 42S02 Table 'q.made3' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made4`
< error
| 1146 42S02 Table 'q.made4' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made2`
< error
| 1146 42S02 Table 'shop.made2' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made3`
< error
| 1146 42S02 Table 'shop.made3' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| made4
: bytes
| CREATE TEMPORARY TABLE `made4` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made2`
< error
| 1146 42S02 Table 'test.made2' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made3`
< error
| 1146 42S02 Table 'test.made3' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made4`
< error
| 1146 42S02 Table 'test.made4' doesn't exist
> exec
| DROP TEMPORARY TABLE made4;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made2`
< error
| 1146 42S02 Table 'q.made2' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made3`
< error
| 1146 42S02 Table 'q.made3' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made4`
< error
| 1146 42S02 Table 'q.made4' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made2`
< error
| 1146 42S02 Table 'shop.made2' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made3`
< error
| 1146 42S02 Table 'shop.made3' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made4`
< error
| 1146 42S02 Table 'shop.made4' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made2`
< error
| 1146 42S02 Table 'test.made2' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made3`
< error
| 1146 42S02 Table 'test.made3' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made4`
< error
| 1146 42S02 Table 'test.made4' doesn't exist
> exec
| SELECT * FROM vv;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made2`
< error
| 1146 42S02 Table 'q.made2' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made3`
< error
| 1146 42S02 Table 'q.made3' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made4`
< error
| 1146 42S02 Table 'q.made4' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| made
: bytes
| CREATE TEMPORARY TABLE `made` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `shop`.`made2`
< error
| 1146 42S02 Table 'shop.made2' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made3`
< error
| 1146 42S02 Table 'shop.made3' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made4`
< error
| 1146 42S02 Table 'shop.made4' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made2`
< error
| 1146 42S02 Table 'test.made2' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made3`
< error
| 1146 42S02 Table 'test.made3' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made4`
< error
| 1146 42S02 Table 'test.made4' doesn't exist
> exec
| DROP TEMPORARY TABLE made;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made2`
< error
| 1146 42S02 Table 'q.made2' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made3`
< error
| 1146 42S02 Table 'q.made3' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made4`
< error
| 1146 42S02 Table 'q.made4' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made2`
< error
| 1146 42S02 Table 'shop.made2' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made3`
< error
| 1146 42S02 Table 'shop.made3' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made4`
< error
| 1146 42S02 Table 'shop.made4' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made2`
< error
| 1146 42S02 Table 'test.made2' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made3`
< error
| 1146 42S02 Table 'test.made3' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made4`
< error
| 1146 42S02 Table 'test.made4' doesn't exist
> exec
| USE q;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made2`
< error
| 1146 42S02 Table 'q.made2' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made3`
< error
| 1146 42S02 Table 'q.made3' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made4`
< error
| 1146 42S02 Table 'q.made4' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made2`
< error
| 1146 42S02 Table 'shop.made2' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made3`
< error
| 1146 42S02 Table 'shop.made3' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made4`
< error
| 1146 42S02 Table 'shop.made4' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made2`
< error
| 1146 42S02 Table 'test.made2' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made3`
< error
| 1146 42S02 Table 'test.made3' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made4`
< error
| 1146 42S02 Table 'test.made4' doesn't exist
> exec
| SELECT mk2();
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made2`
< error
| 1146 42S02 Table 'q.made2' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made3`
< error
| 1146 42S02 Table 'q.made3' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made4`
< error
| 1146 42S02 Table 'q.made4' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made2`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| made2
: bytes
| CREATE TEMPORARY TABLE `made2` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `shop`.`made3`
< error
| 1146 42S02 Table 'shop.made3' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made4`
< error
| 1146 42S02 Table 'shop.made4' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made2`
< error
| 1146 42S02 Table 'test.made2' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made3`
< error
| 1146 42S02 Table 'test.made3' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made4`
< error
| 1146 42S02 Table 'test.made4' doesn't exist
> exec
| USE shop;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made2`
< error
| 1146 42S02 Table 'q.made2' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made3`
< error
| 1146 42S02 Table 'q.made3' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made4`
< error
| 1146 42S02 Table 'q.made4' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made2`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| made2
: bytes
| CREATE TEMPORARY TABLE `made2` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `shop`.`made3`
< error
| 1146 42S02 Table 'shop.made3' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made4`
< error
| 1146 42S02 Table 'shop.made4' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made2`
< error
| 1146 42S02 Table 'test.made2' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made3`
< error
| 1146 42S02 Table 'test.made3' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made4`
< error
| 1146 42S02 Table 'test.made4' doesn't exist
> exec
| DROP TEMPORARY TABLE made2;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made2`
< error
| 1146 42S02 Table 'q.made2' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made3`
< error
| 1146 42S02 Table 'q.made3' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made4`
< error
| 1146 42S02 Table 'q.made4' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made2`
< error
| 1146 42S02 Table 'shop.made2' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made3`
< error
| 1146 42S02 Table 'shop.made3' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made4`
< error
| 1146 42S02 Table 'shop.made4' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made2`
< error
| 1146 42S02 Table 'test.made2' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made3`
< error
| 1146 42S02 Table 'test.made3' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made4`
< error
| 1146 42S02 Table 'test.made4' doesn't exist
> exec
| INSERT INTO nokey SELECT x FROM q.qv;
< ok 1 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made2`
< error
| 1146 42S02 Table 'q.made2' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made3`
< error
| 1146 42S02 Table 'q.made3' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made4`
< error
| 1146 42S02 Table 'q.made4' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made2`
< error
| 1146 42S02 Table 'shop.made2' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made3`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| made3
: bytes
| CREATE TEMPORARY TABLE `made3` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `shop`.`made4`
< error
| 1146 42S02 Table 'shop.made4' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made2`
< error
| 1146 42S02 Table 'test.made2' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made3`
< error
| 1146 42S02 Table 'test.made3' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made4`
< error
| 1146 42S02 Table 'test.made4' doesn't exist
> exec
| DROP TEMPORARY TABLE made3;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made2`
< error
| 1146 42S02 Table 'q.made2' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made3`
< error
| 1146 42S02 Table 'q.made3' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made4`
< error
| 1146 42S02 Table 'q.made4' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made2`
< error
| 1146 42S02 Table 'shop.made2' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made3`
< error
| 1146 42S02 Table 'shop.made3' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made4`
< error
| 1146 42S02 Table 'shop.made4' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made2`
< error
| 1146 42S02 Table 'test.made2' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made3`
< error
| 1146 42S02 Table 'test.made3' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made4`
< error
| 1146 42S02 Table 'test.made4' doesn't exist
> exec
| DROP TEMPORARY TABLE made3;
< error
| 1051 42S02 Unknown table 'shop.made3'
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made2`
< error
| 1146 42S02 Table 'q.made2' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made3`
< error
| 1146 42S02 Table 'q.made3' doesn't exist
> query
| SHOW CREATE TABLE `q`.`made4`
< error
| 1146 42S02 Table 'q.made4' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made2`
< error
| 1146 42S02 Table 'shop.made2' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made3`
< error
| 1146 42S02 Table 'shop.made3' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made4`
< error
| 1146 42S02 Table 'shop.made4' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made2`
< error
| 1146 42S02 Table 'test.made2' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made3`
< error
| 1146 42S02 Table 'test.made3' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made4`
< error
| 1146 42S02 Table 'test.made4' doesn't exist

== TestWalkThroughAgainstMySQLFolded/a_new_table's_text lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_new_table's_text lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/a_new_table's_text lctn1#3
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

== TestWalkThroughAgainstMySQLFolded/a_new_table_holds_what_it_declares lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_new_table_holds_what_it_declares lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	x	fx	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	x	fx	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	x	fx	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	x	fx	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	x	fx	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	x	fx	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| column	shop	y	a	
| column	shop	y	b	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	x	fx	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	x	k	
| index	shop	y	b	
| index	shop	y	k	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/a_new_table_holds_what_it_declares lctn1#3
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
| 1091 42000 Can't DROP 'k'; check that column/key exists
> exec
| ALTER TABLE nokey ADD CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id);
< error
| 1826 HY000 Duplicate foreign key constraint name 'fx'
> exec
| CREATE TABLE y LIKE x;
< ok 0 0
> exec
| ALTER TABLE y DROP INDEX k;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_parent_beside_a_RENAME lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_parent_beside_a_RENAME lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 33
: bytes
| column	q	p2	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	q	p2	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		p2	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	q	p2	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	m	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	q	p2	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		p2	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		m	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
: bytes
| column	q	m	a	
| column	q	p2	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	q	m	m_ibfk_1	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	q	m	a	
| index	q	p2	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		m	BASE TABLE
| relation	q		p2	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 41
: bytes
| column	q	m	a	
| column	q	p2	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	m2	a	
| column	shop	m2	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	q	m	m_ibfk_1	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	q	m	a	
| index	q	p2	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	m2	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		m	BASE TABLE
| relation	q		p2	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		m2	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 43
: bytes
| column	q	m	a	
| column	q	p2	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	m3	a	
| column	shop	m3	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	q	m	m_ibfk_1	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	m3	m3_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	q	m	a	
| index	q	p2	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	m3	PRIMARY	
| index	shop	m3	a	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		m	BASE TABLE
| relation	q		p2	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		m3	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 47
: bytes
| column	q	m	a	
| column	q	p2	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	m3	a	
| column	shop	m3	id	
| column	shop	m4	a	
| column	shop	m4	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	q	m	m_ibfk_1	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	m3	m3_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	q	m	a	
| index	q	p2	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	m3	PRIMARY	
| index	shop	m3	a	
| index	shop	m4	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		m	BASE TABLE
| relation	q		p2	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		m3	BASE TABLE
| relation	shop		m4	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 47
: bytes
| column	q	m	a	
| column	q	p2	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	m3	a	
| column	shop	m3	id	
| column	shop	m4	a	
| column	shop	m4	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	q	m	m_ibfk_1	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	m3	m3_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	q	m	a	
| index	q	p2	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	m3	PRIMARY	
| index	shop	m3	a	
| index	shop	m4	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		m	BASE TABLE
| relation	q		p2	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		m3	BASE TABLE
| relation	shop		m4	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/a_parent_beside_a_RENAME lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| CREATE TABLE q.p2 (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE m (a int);
< ok 0 0
> exec
| ALTER TABLE m RENAME TO q.m, ADD FOREIGN KEY (a) REFERENCES p2 (id);
< ok 0 0
> exec
| CREATE TABLE m2 (id int PRIMARY KEY, a int);
< ok 0 0
> exec
| ALTER TABLE m2 RENAME TO m3, ADD FOREIGN KEY (a) REFERENCES m2 (id);
< ok 0 0
> exec
| CREATE TABLE m4 (id int PRIMARY KEY, a int);
< ok 0 0
> exec
| ALTER TABLE m4 RENAME TO m5, ADD FOREIGN KEY (a) REFERENCES m5 (id);
< error
| 1824 HY000 Failed to open the referenced table 'm5'

== TestWalkThroughAgainstMySQLFolded/a_query's_function_may_create_a_temporary_table lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION mk() RETURNS int DETERMINISTIC BEGIN CREATE TEMPORARY TABLE IF NOT EXISTS made (a int); RETURN 1; END;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_query's_function_may_create_a_temporary_table lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| FUNCTION
: bytes
| mk
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		mk	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	shop	c	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		mk	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		c	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	shop	c	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		mk	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		c	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	shop	c	x	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		mk	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		c	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/a_query's_function_may_create_a_temporary_table lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> exec
| CREATE TABLE c AS SELECT mk() AS x;
< ok 1 0
> query
| SHOW CREATE TABLE `shop`.`made`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| made
: bytes
| CREATE TEMPORARY TABLE `made` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> exec
| DROP TEMPORARY TABLE made;
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> exec
| DROP TEMPORARY TABLE made;
< error
| 1051 42S02 Unknown table 'shop.made'
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist

== TestWalkThroughAgainstMySQLFolded/a_table_with_no_primary_key_may_get_an_invisible_one lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_table_with_no_primary_key_may_get_an_invisible_one lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 31
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	g	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		g	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 31
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	g	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		g	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/a_table_with_no_primary_key_may_get_an_invisible_one lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE g (a int);
< ok 0 0
> exec
| ALTER TABLE g DROP PRIMARY KEY;
< error
| 1091 42000 Can't DROP 'PRIMARY'; check that column/key exists
> exec
| ALTER TABLE g ADD PRIMARY KEY (a);
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_temporary_table's_definition lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_temporary_table's_definition lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/a_temporary_table's_definition lctn1#3
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist
> exec
| CREATE TEMPORARY TABLE tmp (a int);
< error
| 1050 42S01 Table 'tmp' already exists
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `test`.`tmp`
< error
| 1146 42S02 Table 'test.tmp' doesn't exist

== TestWalkThroughAgainstMySQLFolded/a_temporary_table_is_no_parent lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_temporary_table_is_no_parent lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/a_temporary_table_is_no_parent lctn1#3
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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

== TestWalkThroughAgainstMySQLFolded/a_trigger_may_create_a_temporary_table lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_trigger_may_create_a_temporary_table lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
| trigger	shop		tt	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
| trigger	shop		tt	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
| trigger	shop		tt	

== TestWalkThroughAgainstMySQLFolded/a_trigger_may_create_a_temporary_table lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`shadow`
< error
| 1146 42S02 Table 'shop.shadow' doesn't exist
> query
| SHOW CREATE TABLE `test`.`shadow`
< error
| 1146 42S02 Table 'test.shadow' doesn't exist
> exec
| CREATE TRIGGER tt BEFORE INSERT ON nokey FOR EACH ROW CREATE TEMPORARY TABLE IF NOT EXISTS shadow (a int);
< ok 0 0
> query
| SHOW CREATE TABLE `shop`.`shadow`
< error
| 1146 42S02 Table 'shop.shadow' doesn't exist
> query
| SHOW CREATE TABLE `test`.`shadow`
< error
| 1146 42S02 Table 'test.shadow' doesn't exist
> exec
| INSERT INTO nokey VALUES (1);
< ok 1 0
> query
| SHOW CREATE TABLE `shop`.`shadow`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| shadow
: bytes
| CREATE TEMPORARY TABLE `shadow` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `test`.`shadow`
< error
| 1146 42S02 Table 'test.shadow' doesn't exist
> exec
| INSERT INTO shadow VALUES (1);
< ok 1 0
> query
| SHOW CREATE TABLE `shop`.`shadow`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| shadow
: bytes
| CREATE TEMPORARY TABLE `shadow` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `test`.`shadow`
< error
| 1146 42S02 Table 'test.shadow' doesn't exist
> exec
| DROP TABLE shadow;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_trigger_of_another_database_may_create_a_temporary_table lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0
> exec
| DROP TRIGGER tr;
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| CREATE TABLE q.o (a int);
< ok 0 0
> exec
| CREATE TRIGGER q.ot BEFORE INSERT ON q.o FOR EACH ROW CREATE TEMPORARY TABLE IF NOT EXISTS shop.made (a int);
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_trigger_of_another_database_may_create_a_temporary_table lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	q	o	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		o	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	q		ot	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	q	o	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		o	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	q		ot	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	q	o	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		o	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	q		ot	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	q	o	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		o	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	q		ot	

== TestWalkThroughAgainstMySQLFolded/a_trigger_of_another_database_may_create_a_temporary_table lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> exec
| INSERT INTO q.o VALUES (1);
< ok 1 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| made
: bytes
| CREATE TEMPORARY TABLE `made` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> exec
| DROP TEMPORARY TABLE made;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist
> exec
| DROP TEMPORARY TABLE made;
< error
| 1051 42S02 Unknown table 'shop.made'
> query
| SHOW CREATE TABLE `q`.`made`
< error
| 1146 42S02 Table 'q.made' doesn't exist
> query
| SHOW CREATE TABLE `shop`.`made`
< error
| 1146 42S02 Table 'shop.made' doesn't exist
> query
| SHOW CREATE TABLE `test`.`made`
< error
| 1146 42S02 Table 'test.made' doesn't exist

== TestWalkThroughAgainstMySQLFolded/a_view_stays_in_its_database lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/a_view_stays_in_its_database lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/a_view_stays_in_its_database lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| RENAME TABLE v TO q.v;
< error
| 1450 HY000 Changing schema from 'shop' to 'q' is not allowed.
> exec
| DROP VIEW v;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/an_ASCII_name_may_be_one_outside_ASCII lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0
> exec
| CREATE TABLE kel (`K` int, a int, KEY `K` (a));
< ok 0 0
> exec
| CREATE PROCEDURE `pé`() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION `fé`() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT `eé` ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0
> exec
| CREATE TRIGGER `tré` BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/an_ASCII_name_may_be_one_outside_ASCII lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| kel
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| kel
: bytes
| K
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| kel
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| kel
: bytes
| K
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
< row
: bytes
| nokey
: bytes
| tré
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| FUNCTION
: bytes
| fé
< row
: bytes
| PROCEDURE
: bytes
| p
< row
: bytes
| PROCEDURE
: bytes
| pé
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
< row
: bytes
| EVENT
: bytes
| eé
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| kel
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| event	shop		eé	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		fé	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	kel	K	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| procedure	shop		pé	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
| trigger	shop		tré	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 36
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| event	shop		eé	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		fé	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	kel	K	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| procedure	shop		pé	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
| trigger	shop		tré	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| event	shop		eé	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		fé	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| procedure	shop		pé	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
| trigger	shop		tré	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 34
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| event	shop		eé	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		fé	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
| trigger	shop		tré	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 33
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| event	shop		eé	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
| trigger	shop		tré	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
| trigger	shop		tré	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 31
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 31
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 31
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 31
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	kel	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		kel	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/an_ASCII_name_may_be_one_outside_ASCII lctn1#3
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
| ALTER TABLE kel DROP INDEX k;
< ok 0 0
> exec
| DROP PROCEDURE pe;
< ok 0 0
> exec
| DROP FUNCTION shop.fe;
< ok 0 0
> exec
| DROP EVENT ee;
< ok 0 0
> exec
| DROP TRIGGER tre;
< ok 0 0
> exec
| ALTER TABLE kel DROP COLUMN zz;
< error
| 1091 42000 Can't DROP 'zz'; check that column/key exists
> exec
| DROP PROCEDURE nosuch;
< error
| 1305 42000 PROCEDURE shop.nosuch does not exist
> exec
| ALTER TABLE t DROP COLUMN zz;
< error
| 1091 42000 Can't DROP 'zz'; check that column/key exists

== TestWalkThroughAgainstMySQLFolded/an_engine_that_ignores_foreign_keys lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0
> exec
| CREATE TABLE my (id int PRIMARY KEY, a int) ENGINE=MyISAM;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/an_engine_that_ignores_foreign_keys lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| my
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| my
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| my
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| my
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| my
: bytes
| MyISAM
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 33
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	my	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 33
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	my	PRIMARY	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 34
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	my	PRIMARY	
| index	shop	my	fk1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	m2	a	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	m2	mf	
| index	shop	my	PRIMARY	
| index	shop	my	fk1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		m2	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 41
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	m2	a	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	ok	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	ok	mf	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	m2	mf	
| index	shop	my	PRIMARY	
| index	shop	my	fk1	
| index	shop	ok	mf	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		m2	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		ok	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 46
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	m2	a	
| column	shop	m3	a	
| column	shop	m3	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	ok	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	ok	mf	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	m2	mf	
| index	shop	m3	PRIMARY	
| index	shop	m3	fk1	
| index	shop	my	PRIMARY	
| index	shop	my	fk1	
| index	shop	ok	mf	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		m2	BASE TABLE
| relation	shop		m3	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		ok	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 46
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	m2	a	
| column	shop	m3	a	
| column	shop	m3	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	ok	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	ok	mf	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	m2	mf	
| index	shop	m3	PRIMARY	
| index	shop	m3	fk1	
| index	shop	my	PRIMARY	
| index	shop	my	fk1	
| index	shop	ok	mf	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		m2	BASE TABLE
| relation	shop		m3	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		ok	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 50
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	inn	a	
| column	shop	inn	id	
| column	shop	m2	a	
| column	shop	m3	a	
| column	shop	m3	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	ok	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	ok	mf	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	inn	PRIMARY	
| index	shop	m2	mf	
| index	shop	m3	PRIMARY	
| index	shop	m3	fk1	
| index	shop	my	PRIMARY	
| index	shop	my	fk1	
| index	shop	ok	mf	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		inn	BASE TABLE
| relation	shop		m2	BASE TABLE
| relation	shop		m3	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		ok	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 50
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	inn	a	
| column	shop	inn	id	
| column	shop	m2	a	
| column	shop	m3	a	
| column	shop	m3	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	ok	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	ok	mf	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	inn	PRIMARY	
| index	shop	m2	mf	
| index	shop	m3	PRIMARY	
| index	shop	m3	fk1	
| index	shop	my	PRIMARY	
| index	shop	my	fk1	
| index	shop	ok	mf	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		inn	BASE TABLE
| relation	shop		m2	BASE TABLE
| relation	shop		m3	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		ok	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 51
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	inn	a	
| column	shop	inn	id	
| column	shop	m2	a	
| column	shop	m3	a	
| column	shop	m3	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	ok	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	ok	mf	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	inn	PRIMARY	
| index	shop	inn	fk1	
| index	shop	m2	mf	
| index	shop	m3	PRIMARY	
| index	shop	m3	fk1	
| index	shop	my	PRIMARY	
| index	shop	my	fk1	
| index	shop	ok	mf	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		inn	BASE TABLE
| relation	shop		m2	BASE TABLE
| relation	shop		m3	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		ok	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 51
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	inn	a	
| column	shop	inn	id	
| column	shop	m2	a	
| column	shop	m3	a	
| column	shop	m3	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	ok	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	ok	mf	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	inn	PRIMARY	
| index	shop	inn	fk1	
| index	shop	m2	mf	
| index	shop	m3	PRIMARY	
| index	shop	m3	fk1	
| index	shop	my	PRIMARY	
| index	shop	my	fk1	
| index	shop	ok	mf	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		inn	BASE TABLE
| relation	shop		m2	BASE TABLE
| relation	shop		m3	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		ok	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 51
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	inn	a	
| column	shop	inn	id	
| column	shop	m2	a	
| column	shop	m3	a	
| column	shop	m3	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	ok	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	ok	mf	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	inn	PRIMARY	
| index	shop	inn	fk1	
| index	shop	m2	mf	
| index	shop	m3	PRIMARY	
| index	shop	m3	fk1	
| index	shop	my	PRIMARY	
| index	shop	my	fk1	
| index	shop	ok	mf	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		inn	BASE TABLE
| relation	shop		m2	BASE TABLE
| relation	shop		m3	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		ok	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 51
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	inn	a	
| column	shop	inn	id	
| column	shop	m2	a	
| column	shop	m3	a	
| column	shop	m3	id	
| column	shop	my	a	
| column	shop	my	id	
| column	shop	nokey	a	
| column	shop	ok	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	ok	mf	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	inn	PRIMARY	
| index	shop	inn	fk1	
| index	shop	m2	mf	
| index	shop	m3	PRIMARY	
| index	shop	m3	fk1	
| index	shop	my	PRIMARY	
| index	shop	my	fk1	
| index	shop	ok	mf	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		inn	BASE TABLE
| relation	shop		m2	BASE TABLE
| relation	shop		m3	BASE TABLE
| relation	shop		my	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		ok	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/an_engine_that_ignores_foreign_keys lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE my DROP FOREIGN KEY nosuch;
< ok 0 0
> exec
| ALTER TABLE my ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES nosuch (id);
< ok 0 0
> exec
| CREATE TABLE m2 (a int, CONSTRAINT mf FOREIGN KEY (a) REFERENCES nosuch (id), CONSTRAINT mf FOREIGN KEY (a) REFERENCES r (id)) ENGINE=MyISAM;
< ok 0 0
> exec
| CREATE TABLE ok (a int, CONSTRAINT mf FOREIGN KEY (a) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE m3 LIKE my;
< ok 0 0
> exec
| ALTER TABLE m3 ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id);
< ok 0 0
> exec
| CREATE TABLE inn (id int PRIMARY KEY, a int);
< ok 0 0
> exec
| ALTER TABLE inn ENGINE=MyISAM;
< ok 0 0
> exec
| ALTER TABLE inn ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES nosuch (id);
< ok 0 0
> exec
| ALTER TABLE inn DROP FOREIGN KEY whatever;
< ok 0 0
> exec
| ALTER TABLE inn ENGINE=InnoDB;
< ok 0 0
> exec
| ALTER TABLE inn DROP FOREIGN KEY whatever;
< error
| 1091 42000 Can't DROP 'whatever'; check that column/key exists

== TestWalkThroughAgainstMySQLFolded/an_unqualified_parent_is_in_the_child's_database lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/an_unqualified_parent_is_in_the_child's_database lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 33
: bytes
| column	q	parent	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	q	parent	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		parent	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	q	child	child_ibfk_1	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	q	child	id	
| index	q	parent	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		child	BASE TABLE
| relation	q		parent	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	q	child	c2	
| foreign key	q	child	child_ibfk_1	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	q	child	c2	
| index	q	parent	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		child	BASE TABLE
| relation	q		parent	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	q	child	c2	
| foreign key	q	child	child_ibfk_1	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	q	child	c2	
| index	q	parent	PRIMARY	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		child	BASE TABLE
| relation	q		parent	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/an_unqualified_parent_is_in_the_child's_database lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 1 0
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

== TestWalkThroughAgainstMySQLFolded/column_names_a_statement_adds lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0
> exec
| CREATE TABLE cols (a int, b int, c int);
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/column_names_a_statement_adds lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| cols
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| cols
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| cols
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| cols
: bytes
| c
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| cols
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 33
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	c	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 33
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	c	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 33
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	c	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 33
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	c	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 33
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	c	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 33
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	c	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 34
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	c	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 34
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	c	
| column	shop	cols	y	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 36
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 39
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	cols	i	
| index	shop	cols	j	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 39
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	cols	i	
| index	shop	cols	j	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 39
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	cols	i	
| index	shop	cols	j	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	cols	j	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	a	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	A	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	d	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	A	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	d	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	A	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	d	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 38
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	cols	A	
| column	shop	cols	b	
| column	shop	cols	b2	
| column	shop	cols	b3	
| column	shop	cols	c	
| column	shop	cols	d	
| column	shop	cols	y	
| column	shop	cols	y2	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		cols	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/column_names_a_statement_adds lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE cols RENAME COLUMN a TO z, RENAME COLUMN b TO z;
< error
| 1060 42S21 Duplicate column name 'z'
> exec
| ALTER TABLE cols CHANGE a z int, CHANGE b z int;
< error
| 1060 42S21 Duplicate column name 'z'
> exec
| ALTER TABLE cols ADD COLUMN z int, RENAME COLUMN a TO z;
< error
| 1060 42S21 Duplicate column name 'z'
> exec
| ALTER TABLE cols CHANGE a z int, ADD COLUMN Z int;
< error
| 1060 42S21 Duplicate column name 'Z'
> exec
| ALTER TABLE cols RENAME COLUMN a TO b, RENAME COLUMN b TO a;
< ok 0 0
> exec
| ALTER TABLE cols ADD COLUMN b2 int AFTER b;
< ok 0 0
> exec
| ALTER TABLE cols RENAME COLUMN a TO b, RENAME COLUMN b TO y;
< ok 0 0
> exec
| ALTER TABLE cols ADD COLUMN y2 int AFTER y, ADD COLUMN b3 int AFTER b;
< ok 0 0
> exec
| ALTER TABLE cols ADD COLUMN a int;
< ok 0 0
> exec
| ALTER TABLE cols ADD KEY i (a), ADD KEY j (b);
< ok 0 0
> exec
| ALTER TABLE cols RENAME INDEX i TO j, RENAME INDEX j TO i;
< ok 0 0
> exec
| ALTER TABLE cols DROP INDEX i, DROP INDEX j;
< ok 0 0
> exec
| ALTER TABLE cols ADD KEY i (a), ADD KEY j (b);
< ok 0 0
> exec
| ALTER TABLE cols RENAME INDEX i TO j, DROP INDEX j;
< ok 0 0
> exec
| ALTER TABLE cols DROP INDEX j;
< ok 0 0
> exec ^
< error
| 1091 42000 Can't DROP 'j'; check that column/key exists
> exec
| ALTER TABLE cols CHANGE a A int, ADD COLUMN d int;
< ok 0 0
> exec
| ALTER TABLE cols MODIFY a int PRIMARY KEY, MODIFY b int PRIMARY KEY;
< error
| 1068 42000 Multiple primary key defined
> exec
| ALTER TABLE cols CHANGE a a int PRIMARY KEY, ADD PRIMARY KEY (b);
< error
| 1068 42000 Multiple primary key defined
> exec
| ALTER TABLE cols MODIFY a int NOT NULL PRIMARY KEY, ADD COLUMN e int PRIMARY KEY;
< error
| 1068 42000 Multiple primary key defined

== TestWalkThroughAgainstMySQLFolded/create_a_table_that_exists lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/create_a_table_that_exists lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/create_a_table_that_exists lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE t (a int);
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE v (a int);
< error
| 1050 42S01 Table 'v' already exists
> exec
| CREATE TABLE IF NOT EXISTS t (a int);
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS v (a int);
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS v LIKE nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| CREATE TABLE x LIKE nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist

== TestWalkThroughAgainstMySQLFolded/databases lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/databases lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/databases lctn1#3
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
< ok 1 0
> exec ^
< error
| 1007 HY000 Can't create database 'q'; database exists
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 1 0
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
| 1050 42S01 Table 'x' already exists
> exec
| DROP DATABASE q;
< ok 1 0
> exec
| CREATE TABLE q.y (a int);
< error
| 1049 42000 Unknown database 'q'

== TestWalkThroughAgainstMySQLFolded/foreign_keys lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/foreign_keys lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	y	y_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	y	y_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	y	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	y	y_ibfk_1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	y	a	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		y	BASE TABLE
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/foreign_keys lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE y0 (a int, FOREIGN KEY (a) REFERENCES v (id));
< error
| 1824 HY000 Failed to open the referenced table 'v'
> exec
| CREATE TABLE y00 (a int, CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id));
< error
| 1826 HY000 Duplicate foreign key constraint name 'fk1'
> exec
| ALTER TABLE nokey DROP FOREIGN KEY fk1;
< error
| 1091 42000 Can't DROP 'fk1'; check that column/key exists
> exec
| ALTER TABLE nokey ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);
< error
| 1826 HY000 Duplicate foreign key constraint name 'FK1'
> exec
| ALTER TABLE t DROP FOREIGN KEY nosuch;
< error
| 1091 42000 Can't DROP 'nosuch'; check that column/key exists
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

== TestWalkThroughAgainstMySQLFolded/foreign_keys_named_twice lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/foreign_keys_named_twice lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/foreign_keys_named_twice lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a int, b int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT f FOREIGN KEY (b) REFERENCES r (id));
< error
| 1061 42000 Duplicate key name 'f'
> exec
| CREATE TABLE x (a int, b int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT F FOREIGN KEY (a) REFERENCES r (id));
< error
| 1826 HY000 Duplicate foreign key constraint name 'F'
> exec
| CREATE TABLE x (a int, b int);
< ok 0 0
> exec
| ALTER TABLE x ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), ADD CONSTRAINT F FOREIGN KEY (b) REFERENCES r (id);
< error
| 1061 42000 Duplicate key name 'F'

== TestWalkThroughAgainstMySQLFolded/loadable_functions lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/loadable_functions lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		service_get_read_locks	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		service_get_read_locks	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/loadable_functions lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE FUNCTION service_get_read_locks RETURNS INT SONAME 'locking_service.so';
< ok 0 0
> exec
| DROP FUNCTION service_get_read_locks;
< ok 0 0
> exec
| DROP FUNCTION shop.service_get_read_locks;
< error
| 1305 42000 FUNCTION shop.service_get_read_locks does not exist
> exec
| CREATE FUNCTION service_get_read_locks RETURNS INT SONAME 'locking_service.so';
< ok 0 0
> exec
| CREATE FUNCTION service_get_read_locks() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| DROP FUNCTION service_get_read_locks;
< ok 0 0
> exec
| DROP FUNCTION shop.service_get_read_locks;
< ok 0 0
> exec ^
< error
| 1305 42000 FUNCTION shop.service_get_read_locks does not exist

== TestWalkThroughAgainstMySQLFolded/name_case lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/name_case lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 27
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 28
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 29
: bytes
| column	q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 28
: bytes
| column	q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 27
: bytes
| column	q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 26
: bytes
| column	q	x	a	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| relation	q		x	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW

== TestWalkThroughAgainstMySQLFolded/name_case lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE T (a int);
< error
| 1050 42S01 Table 't' already exists
> exec
| DROP TABLE NOKEY;
< ok 0 0
> exec
| CREATE DATABASE Q;
< ok 1 0
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
| 1050 42S01 Table 'x' already exists
> exec
| ALTER TABLE R ADD COLUMN ID int;
< error
| 1060 42S21 Duplicate column name 'ID'
> exec
| CREATE VIEW V AS SELECT 1 AS a;
< error
| 1050 42S01 Table 'V' already exists
> exec
| DROP PROCEDURE P;
< ok 0 0
> exec
| DROP FUNCTION F;
< ok 0 0
> exec
| DROP EVENT E;
< ok 0 0
> exec
| DROP TRIGGER TR;
< ok 0 0
> exec
| CREATE TRIGGER Tr BEFORE INSERT ON r FOR EACH ROW SET NEW.id = 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/names_outside_ASCII lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/names_outside_ASCII lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	café	
| column	shop	x	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	ké	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	ké	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/names_outside_ASCII lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
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
| 1091 42000 Can't DROP 'Ké'; check that column/key exists
> exec
| ALTER TABLE x DROP COLUMN zz;
< error
| 1091 42000 Can't DROP 'zz'; check that column/key exists

== TestWalkThroughAgainstMySQLFolded/names_too_long lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/names_too_long lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/names_too_long lctn1#3
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
| ALTER TABLE t ADD COLUMN ccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc int;
< error
| 1059 42000 Identifier name 'ccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc' is too long

== TestWalkThroughAgainstMySQLFolded/no_database_selected lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/no_database_selected lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 1
: bytes
| database	test			

== TestWalkThroughAgainstMySQLFolded/no_database_selected lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP DATABASE shop;
< ok 6 0
> exec
| ALTER DATABASE CHARACTER SET utf8mb4;
< error
| 1046 3D000 No database selected
> exec
| CREATE TABLE x (a int);
< error
| 1046 3D000 No database selected
> exec
| CREATE PROCEDURE p2() SELECT 1;
< error
| 1046 3D000 No database selected
> exec
| DROP PROCEDURE IF EXISTS p;
< error
| 1046 3D000 No database selected
> exec
| DROP FUNCTION shop.f;
< error
| 1305 42000 FUNCTION shop.f does not exist
> exec
| DROP TRIGGER IF EXISTS tr;
< error
| 1046 3D000 No database selected
> exec
| CREATE EVENT e2 ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< error
| 1046 3D000 No database selected
> exec
| ALTER EVENT e COMMENT 'x';
< error
| 1046 3D000 No database selected

== TestWalkThroughAgainstMySQLFolded/objects_in_a_dropped_database lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/objects_in_a_dropped_database lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| database	z			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/objects_in_a_dropped_database lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `z`.`tmp`
< error
| 1049 42000 Unknown database 'z'
> exec
| CREATE DATABASE z;
< ok 1 0
> query
| SHOW CREATE TABLE `z`.`tmp`
< error
| 1146 42S02 Table 'z.tmp' doesn't exist
> exec
| DROP DATABASE z;
< ok 0 0
> query
| SHOW CREATE TABLE `z`.`tmp`
< error
| 1049 42000 Unknown database 'z'
> exec
| CREATE PROCEDURE z.p() SELECT 1;
< error
| 1049 42000 Unknown database 'z'
> query
| SHOW CREATE TABLE `z`.`tmp`
< error
| 1049 42000 Unknown database 'z'
> exec
| CREATE FUNCTION z.f() RETURNS int DETERMINISTIC RETURN 1;
< error
| 1049 42000 Unknown database 'z'
> query
| SHOW CREATE TABLE `z`.`tmp`
< error
| 1049 42000 Unknown database 'z'
> exec
| CREATE VIEW z.w AS SELECT 1 AS a;
< error
| 1049 42000 Unknown database 'z'
> query
| SHOW CREATE TABLE `z`.`tmp`
< error
| 1049 42000 Unknown database 'z'
> exec
| CREATE TEMPORARY TABLE z.tmp (a int);
< error
| 1049 42000 Unknown database 'z'
> query
| SHOW CREATE TABLE `z`.`tmp`
< error
| 1049 42000 Unknown database 'z'

== TestWalkThroughAgainstMySQLFolded/qualified_events lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| CREATE EVENT q.e ON SCHEDULE EVERY 1 DAY DO SELECT 2;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/qualified_events lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	q		e	
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	q		f	
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| event	shop		g	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	q		h	
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	q		h	
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	q		h	
| event	shop		e	
| event	shop		h	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	q		h	
| event	shop		e	
| event	shop		h	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| event	shop		h	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/qualified_events lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP EVENT q.e;
< ok 0 0
> exec
| ALTER EVENT e COMMENT 'x';
< ok 0 0
> exec
| ALTER EVENT q.e COMMENT 'x';
< error
| 1539 HY000 Unknown event 'e'
> exec
| CREATE EVENT q.f ON SCHEDULE EVERY 1 DAY DO SELECT 3;
< ok 0 0
> exec
| ALTER EVENT q.f RENAME TO g;
< ok 0 0
> exec
| ALTER EVENT g RENAME TO q.h;
< ok 0 0
> exec
| ALTER EVENT q.h COMMENT 'y';
< ok 0 0
> exec
| CREATE EVENT h ON SCHEDULE EVERY 1 DAY DO SELECT 5;
< ok 0 0
> exec
| CREATE EVENT q.h ON SCHEDULE EVERY 1 DAY DO SELECT 4;
< error
| 1537 HY000 Event 'h' already exists
> exec
| DROP EVENT q.h;
< ok 0 0
> exec
| DROP EVENT h;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/qualified_triggers lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| CREATE TABLE q.only_q (id int);
< ok 0 0
> exec
| CREATE TABLE q.t (id int);
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/qualified_triggers lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 34
: bytes
| column	q	only_q	id	
| column	q	t	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		only_q	BASE TABLE
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	q	only_q	id	
| column	q	t	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		only_q	BASE TABLE
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	q		tr	
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	q	only_q	id	
| column	q	t	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		only_q	BASE TABLE
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	q		tr	
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 34
: bytes
| column	q	only_q	id	
| column	q	t	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		only_q	BASE TABLE
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	q	only_q	id	
| column	q	t	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		only_q	BASE TABLE
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	q		tr2	
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	q	only_q	id	
| column	q	t	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		only_q	BASE TABLE
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	q		tr2	
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
: bytes
| column	q	only_q	id	
| column	q	t	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		only_q	BASE TABLE
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	q		tr2	
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 34
: bytes
| column	q	only_q	id	
| column	q	t	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		only_q	BASE TABLE
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 34
: bytes
| column	q	only_q	id	
| column	q	t	id	
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	q		only_q	BASE TABLE
| relation	q		t	BASE TABLE
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/qualified_triggers lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TRIGGER q.tr BEFORE INSERT ON only_q FOR EACH ROW SET NEW.id = 1;
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.id = 1;
< error
| 1359 HY000 Trigger already exists
> exec
| DROP TRIGGER q.tr;
< ok 0 0
> exec
| CREATE TRIGGER q.tr2 BEFORE INSERT ON q.t FOR EACH ROW SET NEW.id = 1;
< ok 0 0
> exec
| CREATE TRIGGER q.tr3 BEFORE INSERT ON shop.t FOR EACH ROW SET NEW.id = 1;
< error
| 1435 HY000 Trigger in wrong schema
> exec
| CREATE TRIGGER tr4 BEFORE INSERT ON q.t FOR EACH ROW SET NEW.id = 1;
< error
| 1435 HY000 Trigger in wrong schema
> exec
| DROP TRIGGER q.tr2;
< ok 0 0
> exec
| DROP TRIGGER q.tr3;
< error
| 1360 HY000 Trigger does not exist
> exec
| DROP TRIGGER tr;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/renamed_generated_foreign_key_names lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/renamed_generated_foreign_key_names lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	g	g_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	g	PRIMARY	
| index	shop	g	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		g	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	h	h_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	h	PRIMARY	
| index	shop	h	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		h	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	h	h_ibfk_1	
| foreign key	shop	nokey	g_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	h	PRIMARY	
| index	shop	h	b	
| index	shop	nokey	g_ibfk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		h	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 36
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	nokey	g_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	h	PRIMARY	
| index	shop	h	b	
| index	shop	nokey	g_ibfk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		h	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 36
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	nokey	g_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	h	PRIMARY	
| index	shop	h	b	
| index	shop	nokey	g_ibfk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		h	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 36
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	kid	kid_ibfk_1	
| foreign key	shop	nokey	g_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	h	PRIMARY	
| index	shop	h	b	
| index	shop	kid	PRIMARY	
| index	shop	kid	b	
| index	shop	nokey	g_ibfk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		h	BASE TABLE
| relation	shop		kid	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	nokey	g_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	h	PRIMARY	
| index	shop	h	b	
| index	shop	kid	PRIMARY	
| index	shop	kid	b	
| index	shop	nokey	g_ibfk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		h	BASE TABLE
| relation	shop		kid	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	nokey	g_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	h	PRIMARY	
| index	shop	h	b	
| index	shop	kid	PRIMARY	
| index	shop	kid	b	
| index	shop	nokey	g_ibfk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		h	BASE TABLE
| relation	shop		kid	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 35
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	nokey	g_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	c2	PRIMARY	
| index	shop	c2	b	
| index	shop	h	PRIMARY	
| index	shop	h	b	
| index	shop	nokey	g_ibfk_1	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		c2	BASE TABLE
| relation	shop		h	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/renamed_generated_foreign_key_names lctn1#3
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
| ALTER TABLE nokey ADD CONSTRAINT g_ibfk_1 FOREIGN KEY (a) REFERENCES r (id);
< ok 0 0
> exec
| ALTER TABLE h DROP FOREIGN KEY h_ibfk_1;
< ok 0 0
> exec ^
< error
| 1091 42000 Can't DROP 'h_ibfk_1'; check that column/key exists
> exec
| RENAME TABLE child TO kid;
< ok 0 0
> exec
| ALTER TABLE kid DROP FOREIGN KEY kid_ibfk_1;
< ok 0 0
> exec ^
< error
| 1091 42000 Can't DROP 'kid_ibfk_1'; check that column/key exists
> exec
| ALTER TABLE kid RENAME TO c2;
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, CONSTRAINT child_ibfk_1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/repeated_subcommands lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/repeated_subcommands lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/repeated_subcommands lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN c int, ADD COLUMN c int;
< error
| 1060 42S21 Duplicate column name 'c'
> exec
| ALTER TABLE t DROP COLUMN a, DROP COLUMN a;
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| ALTER TABLE t ADD INDEX k (a), ADD INDEX k (b);
< error
| 1061 42000 Duplicate key name 'k'
> exec
| ALTER TABLE nokey ADD PRIMARY KEY (a), ADD PRIMARY KEY (a);
< error
| 1068 42000 Multiple primary key defined
> exec
| ALTER TABLE t DROP INDEX i, DROP INDEX i;
< error
| 1091 42000 Can't DROP 'i'; check that column/key exists
> exec
| ALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;
< error
| 1091 42000 Can't DROP 'fk1'; check that column/key exists
> exec
| ALTER TABLE t DROP PRIMARY KEY, DROP PRIMARY KEY;
< error
| 1091 42000 Can't DROP 'PRIMARY'; check that column/key exists
> exec
| ALTER TABLE t DROP PRIMARY KEY, DROP INDEX `PRIMARY`;
< error
| 1091 42000 Can't DROP 'PRIMARY'; check that column/key exists

== TestWalkThroughAgainstMySQLFolded/routines_and_events lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/routines_and_events lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e2	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e2	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e2	
| event	shop		shop.e2	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e2	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e2	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| function	shop		p	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/routines_and_events lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< error
| 1304 42000 PROCEDURE p already exists
> exec
| CREATE PROCEDURE P() SELECT 1;
< error
| 1304 42000 PROCEDURE P already exists
> exec
| CREATE FUNCTION p() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< error
| 1304 42000 FUNCTION f already exists
> exec
| DROP PROCEDURE nosuch;
< error
| 1305 42000 PROCEDURE shop.nosuch does not exist
> exec
| DROP FUNCTION shop.nosuch;
< error
| 1305 42000 FUNCTION shop.nosuch does not exist
> exec
| ALTER PROCEDURE nosuch COMMENT 'x';
< error
| 1305 42000 PROCEDURE shop.nosuch does not exist
> exec
| CREATE EVENT E ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< error
| 1537 HY000 Event 'E' already exists
> exec
| DROP EVENT nosuch;
< error
| 1539 HY000 Unknown event 'nosuch'
> exec
| ALTER EVENT nosuch COMMENT 'x';
< error
| 1539 HY000 Unknown event 'nosuch'
> exec
| ALTER EVENT e RENAME TO e2;
< ok 0 0
> exec
| DROP EVENT e;
< error
| 1539 HY000 Unknown event 'e'
> exec
| CREATE EVENT `shop.e2` ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0
> exec
| DROP EVENT `shop.e2`;
< ok 0 0
> exec ^
< error
| 1539 HY000 Unknown event 'shop.e2'

== TestWalkThroughAgainstMySQLFolded/temporary_tables lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/temporary_tables lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 30
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	nokey	z	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/temporary_tables lctn1#3
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
|   PRIMARY KEY (`id`),
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`),
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`),
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`),
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`),
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`),
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| 1050 42S01 Table 'tt' already exists
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
|   PRIMARY KEY (`id`),
|   KEY `i` (`a`),
|   KEY `fk1` (`b`),
|   CONSTRAINT `fk1` FOREIGN KEY (`b`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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

== TestWalkThroughAgainstMySQLFolded/temporary_tables_outlive_their_database lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/temporary_tables_outlive_their_database lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	q			
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/temporary_tables_outlive_their_database lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`tmp`
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE DATABASE q;
< ok 1 0
> query
| SHOW CREATE TABLE `q`.`tmp`
< error
| 1146 42S02 Table 'q.tmp' doesn't exist
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE q;
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE DATABASE q;
< ok 1 0
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TEMPORARY TABLE q.tmp;
< ok 0 0
> query
| SHOW CREATE TABLE `q`.`tmp`
< error
| 1146 42S02 Table 'q.tmp' doesn't exist
> exec
| DROP TEMPORARY TABLE q.tmp;
< error
| 1051 42S02 Unknown table 'q.tmp'
> query
| SHOW CREATE TABLE `q`.`tmp`
< error
| 1146 42S02 Table 'q.tmp' doesn't exist

== TestWalkThroughAgainstMySQLFolded/the_sync_may_not_list_an_invisible_primary_key lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0
> exec
| SET GLOBAL show_gipk_in_create_table_and_information_schema = OFF;
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = ON;
< ok 0 0
> exec
| CREATE TABLE gk (a int);
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/the_sync_may_not_list_an_invisible_primary_key lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| gk
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| gk
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| gk
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 31
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	gk	a	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		gk	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	gk	a	
| column	shop	gk	x	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		gk	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	gk	a	
| column	shop	gk	x	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		gk	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 32
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	gk	a	
| column	shop	gk	x	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		gk	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/the_sync_may_not_list_an_invisible_primary_key lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| ALTER TABLE gk ADD COLUMN x int AFTER my_row_id;
< ok 0 0
> exec
| ALTER TABLE gk ADD COLUMN my_row_id int;
< error
| 1060 42S21 Duplicate column name 'my_row_id'
> exec
| ALTER TABLE gk ADD PRIMARY KEY (a);
< error
| 1068 42000 Multiple primary key defined

== TestWalkThroughAgainstMySQLFolded/triggers lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/triggers lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW

== TestWalkThroughAgainstMySQLFolded/triggers lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;
< error
| 1359 HY000 Trigger already exists
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON nosuch FOR EACH ROW SET NEW.a = 1;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| DROP TRIGGER nosuch;
< error
| 1360 HY000 Trigger does not exist
> exec
| DROP TRIGGER tr;
< ok 0 0
> exec ^
< error
| 1360 HY000 Trigger does not exist
> exec
| CREATE TRIGGER TR BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/unknown_by_design lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/unknown_by_design lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
< row
: int64 37
: bytes
| column	shop	child	b	
| column	shop	child	id	
| column	shop	nokey	a	
| column	shop	r	id	
| column	shop	t	a	
| column	shop	t	b	
| column	shop	t	id	
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| column	shop	y	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	y	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	y	fy	
| index	shop	z	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	y	fy	
| index	shop	z	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	y	fy	
| index	shop	z	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	y	fy	
| index	shop	z	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	y	fy	
| index	shop	z	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| column	shop	x	a	
| column	shop	x	b	
| column	shop	y	a	
| column	shop	z	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| foreign key	shop	y	fy	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| index	shop	x	b	
| index	shop	y	fy	
| index	shop	z	fy	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| relation	shop		x	BASE TABLE
| relation	shop		y	BASE TABLE
| relation	shop		z	BASE TABLE
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/unknown_by_design lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP FUNCTION nosuch;
< error
| 1305 42000 FUNCTION shop.nosuch does not exist
> exec
| DROP DATABASE nosuch;
< error
| 1008 HY000 Can't drop database 'nosuch'; database doesn't exist
> exec
| RENAME TABLE t TO nodb.t;
< error
| 1049 42000 Unknown database 'nodb'
> exec
| CREATE VIEW nodb.w AS SELECT 1 AS a;
< error
| 1049 42000 Unknown database 'nodb'
> exec
| CREATE TABLE x AS SELECT * FROM nosuch;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| CREATE TABLE x (a int, b int UNIQUE);
< ok 0 0
> exec
| ALTER TABLE x DROP INDEX zz;
< error
| 1091 42000 Can't DROP 'zz'; check that column/key exists
> exec
| CREATE TABLE y (a int, CONSTRAINT fy FOREIGN KEY (a) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE z LIKE y;
< ok 0 0
> exec
| ALTER TABLE z DROP FOREIGN KEY fy;
< error
| 1091 42000 Can't DROP 'fy'; check that column/key exists
> exec
| ALTER TABLE nokey DROP PRIMARY KEY;
< error
| 1091 42000 Can't DROP 'PRIMARY'; check that column/key exists
> exec
| ALTER TABLE t ADD COLUMN z int, DROP COLUMN z;
< error
| 1091 42000 Can't DROP 'z'; check that column/key exists
> exec
| ALTER TABLE t ADD INDEX j (b), DROP INDEX j;
< error
| 1091 42000 Can't DROP 'j'; check that column/key exists
> exec
| ALTER TABLE t ADD COLUMN x int, RENAME COLUMN x TO y;
< error
| 1054 42S22 Unknown column 'x' in 't'

== TestWalkThroughAgainstMySQLFolded/views lctn1#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE r (id int PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
< ok 0 0
> exec
| CREATE TABLE nokey (a int);
< ok 0 0
> exec
| CREATE TABLE trg (a int);
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id FROM r;
< ok 0 0
> exec
| CREATE PROCEDURE p() SELECT 1;
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
< ok 0 0

== TestWalkThroughAgainstMySQLFolded/views lctn1#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| child
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| r
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| trg
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| child
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| child
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| r
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| trg
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| v
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| child
: bytes
| b
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| child
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| r
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk1
: int64 1
: bytes
| b
: int64 -1
< row
: bytes
| t
: bytes
| i
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| child
: bytes
| child_ibfk_1
: bytes
| r
: bytes
| b
: bytes
| id
< row
: bytes
| t
: bytes
| fk1
: bytes
| r
: bytes
| b
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| trg
: bytes
| tr
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "ENUM"
< column "ROUTINE_NAME" "VARCHAR"
< row
: bytes
| FUNCTION
: bytes
| f
< row
: bytes
| PROCEDURE
: bytes
| p
> query
| SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT" "VARCHAR"
< column "EVENT_NAME" "VARCHAR"
< row
: bytes
| EVENT
: bytes
| e
> query
| SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "ENGINE" "VARCHAR"
< row
: bytes
| child
: bytes
| InnoDB
< row
: bytes
| nokey
: bytes
| InnoDB
< row
: bytes
| r
: bytes
| InnoDB
< row
: bytes
| t
: bytes
| InnoDB
< row
: bytes
| trg
: bytes
| InnoDB
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, ''), EVENT_SCHEMA FROM information_schema.EVENTS
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	id	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "LONGTEXT"
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
| column	shop	trg	a	
| column	shop	v	a	
| database	shop			
| database	test			
| event	shop		e	
| foreign key	shop	child	child_ibfk_1	
| foreign key	shop	t	fk1	
| function	shop		f	
| index	shop	child	PRIMARY	
| index	shop	child	b	
| index	shop	r	PRIMARY	
| index	shop	t	PRIMARY	
| index	shop	t	fk1	
| index	shop	t	i	
| procedure	shop		p	
| relation	shop		child	BASE TABLE
| relation	shop		nokey	BASE TABLE
| relation	shop		r	BASE TABLE
| relation	shop		t	BASE TABLE
| relation	shop		trg	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	

== TestWalkThroughAgainstMySQLFolded/views lctn1#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE VIEW t AS SELECT 1 AS a;
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE VIEW v AS SELECT 1 AS a;
< error
| 1050 42S01 Table 'v' already exists
> exec
| CREATE OR REPLACE VIEW t AS SELECT 1 AS a;
< error
| 1347 HY000 'shop.t' is not VIEW
> exec
| CREATE OR REPLACE VIEW v AS SELECT 1 AS a;
< ok 0 0
> exec
| ALTER VIEW nosuch AS SELECT 1 AS a;
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> exec
| ALTER VIEW t AS SELECT 1 AS a;
< error
| 1347 HY000 'shop.t' is not VIEW
> exec
| DROP VIEW IF EXISTS t;
< ok 0 0
> exec
| DROP VIEW t;
< error
| 1347 HY000 'shop.t' is not VIEW
