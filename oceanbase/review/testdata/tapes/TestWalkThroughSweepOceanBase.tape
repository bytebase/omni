-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestWalkThroughSweepOceanBase/change0 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE INDEX i ON v (y);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE A (x int, y int, KEY I (x));
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE INDEX j (x);
< ok 0 0
> exec
| CREATE TEMPORARY TABLE A (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| DROP INDEX x ON a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE v DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP TRIGGER IF EXISTS TR;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS a_obfk_1 LIKE a;
< error
| 1347 HY000 'shop.a' is not BASE TABLE

== TestWalkThroughSweepOceanBase/change0 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| a
: bytes
| I
: int64 1
: bytes
| x
: int64 -1
< row
: bytes
| a
: bytes
| j
: int64 0
: bytes
| x
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| index	shop	a	I	
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
| index	shop	a	I	
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
| index	shop	a	I	
| index	shop	a	j	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	a	x	
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	q		a	BASE TABLE
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	a	x	
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	q		a	BASE TABLE
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	a	x	
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	q		a	BASE TABLE
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	a	x	
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	q		a	BASE TABLE
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	a	x	
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	q		a	BASE TABLE
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	a	x	
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	q		a	BASE TABLE
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	a	x	
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	q		a	BASE TABLE
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	a	x	
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	q		a	BASE TABLE
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	q	a	x	
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	j	
| relation	q		a	BASE TABLE
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	a	I	
| index	shop	a	j	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	a	I	
| index	shop	a	j	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	a	I	
| index	shop	a	j	
| index	shop	c	PRIMARY	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	a	I	
| index	shop	a	j	
| index	shop	c	PRIMARY	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	a	I	
| index	shop	a	j	
| index	shop	c	PRIMARY	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	a	I	
| index	shop	a	j	
| index	shop	c	PRIMARY	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 22
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	a	I	
| index	shop	a	j	
| index	shop	c	PRIMARY	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	c	PRIMARY	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	c	PRIMARY	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	c	PRIMARY	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	c	PRIMARY	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	c	PRIMARY	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	q	a	x	
| column	q	a_obfk_1	id	
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	q			
| database	shop			
| index	q	a_obfk_1	PRIMARY	
| index	q	a_obfk_1	x	
| index	shop	c	PRIMARY	
| relation	q		a	BASE TABLE
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepOceanBase/change0 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| DROP INDEX x ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE a ADD INDEX i (y);
< error
| 1061 42000 Duplicate key name 'i'
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 1 0
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE b ADD UNIQUE INDEX x (z);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE q.A AS SELECT 1 AS x;
< ok 1 0
> exec
| INSERT INTO c (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE a_obfk_1 LIKE a_obfk_1;
< error
| 1066 42000 Not unique table/alias: 'a_obfk_1'
> exec
| ALTER TABLE b RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE c DROP INDEX j, ADD INDEX j (x);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE q.b CHANGE y y int;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE a_obfk_1 ADD FOREIGN KEY g (x) REFERENCES a_obfk_1 (id);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| INSERT INTO a_obfk_1 (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS q.a_obfk_1 (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr2;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE VIEW c AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'c' already exists
> exec
| DROP INDEX x ON a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TABLE q.b LIKE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE v AS SELECT 1 AS x;
< ok 1 0
> exec
| DROP TABLE IF EXISTS b, a;
< ok 0 0
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec
| ALTER TABLE q.v ADD COLUMN y int;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| INSERT INTO v (x) VALUES (1);
< ok 1 0
> exec
| ALTER TABLE c ENGINE = InnoDB;
< ok 0 0
> exec
| CREATE INDEX j ON c (z);
< error
| 1072 42000 Key column 'z' doesn't exist in table
> exec
| SET GLOBAL foreign_key_checks = 0, foreign_key_checks = 1;
< ok 0 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change1 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| SET GLOBAL foreign_key_checks = 0, foreign_key_checks = 1;
< ok 0 0
> exec
| DROP TRIGGER TR;
< error
| 1360 HY000 Trigger does not exist
> exec
| ALTER TABLE q.c ADD CONSTRAINT f FOREIGN KEY (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP INDEX j ON q.c;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DROP VIEW a_obfk_1;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1'
> exec
| CREATE TABLE v (x int, y int, KEY j (x));
< ok 0 0
> exec
| DROP INDEX j ON a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE DATABASE q;
< error
| 1007 HY000 Can't create database 'q'; database exists
> exec
| DROP TABLE IF EXISTS q.a_obfk_1;
< ok 0 0
> exec
| ALTER TABLE q.c DROP FOREIGN KEY g;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< ok 1 0

== TestWalkThroughSweepOceanBase/change1 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
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
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| b
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| v
: bytes
| j
: int64 1
: bytes
| x
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	q	b	id	
| column	q	b	x	
| column	q	b	y	
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	q	b	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	q	b	id	
| column	q	b	x	
| column	q	b	y	
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	q	b	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	q	b	id	
| column	q	b	x	
| column	q	b	y	
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	q	b	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	q	b	id	
| column	q	b	x	
| column	q	b	y	
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	q	b	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	q	b	id	
| column	q	b	x	
| column	q	b	y	
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	q	b	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	q	b	id	
| column	q	b	x	
| column	q	b	y	
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	q	b	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	q		b	BASE TABLE
| relation	shop		a	BASE TABLE
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
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
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
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
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
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
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
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
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
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
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
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
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
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
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
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	j	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepOceanBase/change1 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| DROP TRIGGER TR;
< error
| 1360 HY000 Trigger does not exist
> exec
| CREATE TABLE A LIKE c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE c DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE v AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'v' already exists
> exec
| ALTER TABLE a DROP INDEX I;
< error
| 1091 42000 Can't DROP 'I'; check that column/key exists
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE q.b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 DROP FOREIGN KEY a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'a' already exists
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec
| RENAME TABLE q.a TO a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| CREATE UNIQUE INDEX I ON c (z);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 1 0
> exec
| TRUNCATE TABLE v;
< ok 0 0
> exec
| RENAME TABLE v TO q.A;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE INDEX x ON q.c (z);
< error
| 1049 42000 Unknown database
> exec
| ALTER TABLE v ADD FOREIGN KEY a_obfk_1 (x) REFERENCES a (id);
< error
| 5031 HY000 Column not found
> exec
| DROP INDEX i ON a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE q.v DROP COLUMN z, ADD COLUMN X int;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| TRUNCATE TABLE q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP TABLE IF EXISTS a, c;
< ok 0 0
> exec
| RENAME TABLE A TO c;
< error
| 1017 HY000 File not exist
> exec
| ALTER TABLE q.a_obfk_1 DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| DROP TEMPORARY TABLE a_obfk_1;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1'
> exec
| CREATE INDEX j ON A (X);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| ALTER TABLE q.c CHANGE y x int;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE a DROP INDEX I;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE a_obfk_1 LIKE b;
< ok 0 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change10 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE q.c RENAME TO a_obfk_1;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DROP TABLE IF EXISTS b, A;
< ok 0 0
> exec
| CREATE TEMPORARY TABLE c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| DROP TABLE IF EXISTS b;
< ok 0 0
> exec
| ALTER TABLE q.A DROP INDEX I;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE a_obfk_1 ADD COLUMN z int AFTER y;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| TRUNCATE TABLE c;
< ok 0 0
> exec
| ALTER TABLE q.b DROP INDEX I, ADD INDEX i (x);
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| RENAME TABLE q.v TO A, q.v TO b;
< error
| 1017 HY000 File not exist
> exec
| ALTER TABLE b ADD INDEX I (y);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE a_obfk_1 RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP TRIGGER tr2;
< error
| 1360 HY000 Trigger does not exist

== TestWalkThroughSweepOceanBase/change10 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	q			
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	q			
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	q			
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	q	a	x	
| column	q	a	y	
| database	q			
| database	shop			
| index	q	a	I	
| relation	q		a	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	q	a	x	
| column	q	a	y	
| database	q			
| database	shop			
| index	q	a	I	
| relation	q		a	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	q	a	x	
| column	q	a	y	
| database	q			
| database	shop			
| index	q	a	I	
| relation	q		a	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	q	a	x	
| column	q	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a	I	
| index	shop	c	I	
| relation	q		a	BASE TABLE
| relation	shop		c	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	q	a	x	
| column	q	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a	I	
| index	shop	c	I	
| relation	q		a	BASE TABLE
| relation	shop		c	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	q	a	x	
| column	q	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a	I	
| index	shop	c	I	
| relation	q		a	BASE TABLE
| relation	shop		c	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	q	a	x	
| column	q	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a	I	
| index	shop	c	I	
| relation	q		a	BASE TABLE
| relation	shop		c	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	q	a	x	
| column	q	a	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a	I	
| index	shop	c	I	
| relation	q		a	BASE TABLE
| relation	shop		c	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	q	a	x	
| column	q	a	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a	I	
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	q		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	q	a	x	
| column	q	a	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a	I	
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	q		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	q	a	x	
| column	q	a	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a	I	
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	q		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| temporary	shop		a	

== TestWalkThroughSweepOceanBase/change10 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| CREATE TRIGGER tr BEFORE INSERT ON v FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP VIEW a;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| CREATE TEMPORARY TABLE a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE a_obfk_1 LIKE A;
< error
| 1347 HY000 'shop.a' is not BASE TABLE
> exec
| ALTER TABLE q.b DROP COLUMN X, ADD COLUMN X int;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| CREATE TRIGGER TR BEFORE INSERT ON v FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE q.b ENGINE = InnoDB;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| CREATE TABLE q.A (x int, y int, KEY I (x));
< ok 0 0
> exec
| DROP INDEX I ON A;
< error
| 1091 42000 Can't DROP 'I'; check that column/key exists
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| CREATE TABLE c (x int, y int, KEY I (x));
< ok 0 0
> exec
| INSERT INTO a (x) VALUES (1);
< error
| 1364 HY000 Field 'id' doesn't have a default value
> exec
| TRUNCATE TABLE q.c;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TABLE c (x int, y int, KEY I (x));
< error
| 1050 42S01 Table 'c' already exists
> exec
| CREATE TRIGGER TR BEFORE INSERT ON v FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP INDEX I ON q.a_obfk_1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change11 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE q.b (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE q.A DROP COLUMN z, ADD COLUMN X int;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP INDEX i ON q.c;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON a_obfk_1 FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE a ADD INDEX I (x);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE A ENGINE = InnoDB;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| CREATE INDEX j ON b (z);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TABLE q.v, c;
< error
| 1051 42S02 Unknown table 'q.v,shop.c'
> exec
| RENAME TABLE q.v TO q.b;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE c DROP INDEX I, ADD INDEX i (x);
< error
| 1146 42S02 Table 'shop.c' doesn't exist

== TestWalkThroughSweepOceanBase/change11 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	b	x	
| database	shop			
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	b	x	
| database	shop			
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	b	x	
| database	shop			
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	b	x	
| database	shop			
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE

== TestWalkThroughSweepOceanBase/change11 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| ALTER TABLE c ADD CONSTRAINT g FOREIGN KEY g (x) REFERENCES q.a (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS v LIKE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE INDEX i ON A (x);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS a_obfk_1 LIKE A;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| RENAME TABLE a TO b, q.A TO c;
< error
| 1017 HY000 File not exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| CREATE TABLE q.c LIKE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS q.v, q.b;
< ok 0 0
> exec
| ALTER TABLE b DROP INDEX j, ADD INDEX i (x);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE v DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS q.a (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP TEMPORARY TABLE c, c;
< error
| 1066 42000 Not unique table/alias: 'c'
> exec
| ALTER TABLE a_obfk_1 DROP COLUMN y;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE q.A ADD PRIMARY KEY (x);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE q.a CHANGE x X int;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec ^
< ok 0 0
> exec
| ALTER TABLE b ADD FOREIGN KEY a_obfk_1 (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE a DROP INDEX j;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP INDEX j ON q.c;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TABLE b AS SELECT 1 AS x;
< ok 1 0
> exec
| ALTER TABLE v DROP COLUMN z;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE a ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS v;
< ok 0 0
> exec
| CREATE TABLE a_obfk_1 AS SELECT 1 AS x;
< ok 1 0
> exec
| ALTER TABLE b RENAME TO q.b;
< error
| 1210 HY000 Invalid argument
> exec
| ALTER TABLE a DROP INDEX I;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| RENAME TABLE v TO a_obfk_1, a_obfk_1 TO a;
< error
| 1050 42S01 Table 'a_obfk_1' already exists
> exec
| ALTER TABLE A DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change12 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| USE shop;
< ok 0 0
> exec
| CREATE VIEW q.c AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE q.a ADD CONSTRAINT g FOREIGN KEY g (x) REFERENCES v (id);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE A AS SELECT 1 AS x;
< ok 1 0
> exec
| DROP TRIGGER IF EXISTS TR;
< ok 0 0
> exec
| DROP INDEX x ON b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE q.b LIKE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TEMPORARY TABLE b;
< error
| 1051 42S02 Unknown table 'shop.b'
> exec
| DROP TABLE v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| CREATE TABLE IF NOT EXISTS q.a LIKE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP TRIGGER TR;
< error
| 1360 HY000 Trigger does not exist
> exec
| ALTER TABLE q.a_obfk_1 DROP FOREIGN KEY f;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist

== TestWalkThroughSweepOceanBase/change12 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	v	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	q	b	id	
| column	q	b	x	
| column	q	b	y	
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	q	b	PRIMARY	
| index	shop	v	PRIMARY	
| relation	q		b	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	q	b	id	
| column	q	b	x	
| column	q	b	y	
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	q	b	PRIMARY	
| index	shop	v	PRIMARY	
| relation	q		b	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	q	b	id	
| column	q	b	x	
| column	q	b	y	
| column	shop	a	x	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	q	b	PRIMARY	
| index	shop	v	PRIMARY	
| relation	q		b	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	

== TestWalkThroughSweepOceanBase/change12 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| CREATE TRIGGER tr BEFORE INSERT ON a FOR EACH ROW SET NEW.x = 1;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP INDEX i ON v;
< error
| 1091 42000 Can't DROP 'i'; check that column/key exists
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec ^
< ok 0 0
> exec
| CREATE TABLE q.v (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| DROP INDEX x ON a;
< error
| 1091 42000 Can't DROP 'x'; check that column/key exists
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 1 0
> exec
| DROP VIEW c;
< error
| 1051 42S02 Unknown table 'shop.c'
> exec
| TRUNCATE TABLE q.v;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| CREATE TABLE A LIKE a;
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| CREATE TABLE q.b LIKE v;
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1364 HY000 Field 'id' doesn't have a default value
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change13 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE q.b (x int, y int, KEY x (x), CONSTRAINT F FOREIGN KEY (y) REFERENCES a (id));
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE a_obfk_1 ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TABLE v LIKE v;
< error
| 1066 42000 Not unique table/alias: 'v'
> exec
| ALTER TABLE q.c DROP COLUMN X, ADD COLUMN z int;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< ok 1 0
> exec
| CREATE TABLE A AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'a' already exists
> exec
| CREATE VIEW a_obfk_1 AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 CHANGE x X int;
< error
| 1347 HY000 'shop.a_obfk_1' is not BASE TABLE
> exec
| DROP TABLE a_obfk_1, A;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1'
> exec
| CREATE TABLE b LIKE q.a_obfk_1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist

== TestWalkThroughSweepOceanBase/change13 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| a_obfk_1
: bytes
| VIEW
: bytes
| MEMORY
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| x
< row
: bytes
| a_obfk_1
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	q			
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	q			
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	q			
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	q			
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	q			
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	q			
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	q			
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	q			
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW

== TestWalkThroughSweepOceanBase/change13 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| SET GLOBAL foreign_key_checks = 1, @@session.foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE c DROP INDEX I;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE c ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE q.v AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TEMPORARY TABLE A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE v DROP COLUMN z;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP VIEW c;
< error
| 1051 42S02 Unknown table 'shop.c'
> exec
| CREATE UNIQUE INDEX I ON q.a_obfk_1 (z);
< error
| 1049 42000 Unknown database
> exec
| ALTER TABLE A ADD INDEX i (y);
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS q.a_obfk_1 (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| TRUNCATE TABLE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE q.b ADD FOREIGN KEY a_obfk_1 (x) REFERENCES a_obfk_1 (id);
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 1 0
> exec
| CREATE TABLE A (x int, y int, KEY x (x));
< error
| 1050 42S01 Table 'a' already exists
> exec
| TRUNCATE TABLE q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP TEMPORARY TABLE b;
< error
| 1051 42S02 Unknown table 'shop.b'
> exec
| DROP TABLE A;
< ok 0 0
> exec
| ALTER TABLE v ADD PRIMARY KEY (x);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP VIEW v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| CREATE TEMPORARY TABLE v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE c DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE a_obfk_1 DROP FOREIGN KEY a_obfk_1;
< error
| 1347 HY000 'shop.a_obfk_1' is not BASE TABLE
> exec
| DROP TABLE v;
< ok 0 0
> exec
| CREATE TABLE v LIKE q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change14 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE b DROP COLUMN y;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TABLE a;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| ALTER TABLE a_obfk_1 DROP INDEX I;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON A FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| CREATE TEMPORARY TABLE b (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS a_obfk_1 (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE VIEW a_obfk_1 AS SELECT 1 AS x;
< ok 0 0
> exec
| DROP VIEW IF EXISTS q.a;
< ok 0 0
> exec
| CREATE TEMPORARY TABLE c (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0

== TestWalkThroughSweepOceanBase/change14 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a_obfk_1
: bytes
| VIEW
: bytes
| MEMORY
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a_obfk_1
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a_obfk_1	VIEW
| relation	shop		v	VIEW

== TestWalkThroughSweepOceanBase/change14 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| TRUNCATE TABLE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP TABLE IF EXISTS a_obfk_1;
< ok 0 0
> exec
| TRUNCATE TABLE a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP VIEW q.b;
< error
| 1051 42S02 Unknown table 'q.b'
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| ALTER TABLE a DROP FOREIGN KEY f;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| INSERT INTO A (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP TABLE IF EXISTS q.A, c;
< ok 0 0
> exec
| ALTER TABLE b DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TABLE b;
< error
| 1051 42S02 Unknown table 'shop.b'
> exec
| ALTER TABLE q.a_obfk_1 ADD FOREIGN KEY a_obfk_1 (x) REFERENCES A (id);
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| CREATE OR REPLACE VIEW v AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE VIEW q.c AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE b ADD COLUMN x int AFTER X;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE a CHANGE z X int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| RENAME TABLE q.a_obfk_1 TO c;
< error
| 1017 HY000 File not exist
> exec
| DROP VIEW IF EXISTS q.b;
< ok 0 0
> exec
| ALTER TABLE b ADD CONSTRAINT f FOREIGN KEY (x) REFERENCES q.v (id);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE OR REPLACE VIEW q.A AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE q.v ADD COLUMN z int AFTER X;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS v;
< ok 0 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change15 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE q.b ADD FOREIGN KEY g (x) REFERENCES v (id);
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| INSERT INTO q.c (x) VALUES (1);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON b FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TEMPORARY TABLE a_obfk_1 (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 RENAME INDEX i TO j;
< error
| 1176 42000 Key 'i' doesn't exist in table 'a_obfk_1'
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON q.a_obfk_1 FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS a;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE A RENAME TO a_obfk_1;
< error
| 1050 42S01 Table 'a_obfk_1' already exists
> exec
| DROP INDEX j ON a_obfk_1;
< error
| 1091 42000 Can't DROP 'j'; check that column/key exists
> exec
| ALTER TABLE a DROP INDEX x, ADD INDEX i (x);
< error
| 1091 42000 Can't DROP 'x'; check that column/key exists
> exec
| DROP TABLE q.a_obfk_1;
< error
| 1051 42S02 Unknown table 'q.a_obfk_1'

== TestWalkThroughSweepOceanBase/change15 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
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
| MEMORY
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| v
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| a
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
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
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
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
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
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
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
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
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
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
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
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
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
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
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
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
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 20
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 21
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	id	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr2	

== TestWalkThroughSweepOceanBase/change15 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| CREATE TABLE c (x int, y int, KEY j (x));
< ok 0 0
> exec
| INSERT INTO q.c (x) VALUES (1);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TRIGGER tr BEFORE INSERT ON v FOR EACH ROW SET NEW.x = 1;
< error
| 1235 0A000 simple dml trigger isn't used on user table not supported
> exec
| ALTER TABLE q.c DROP COLUMN X;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TABLE a_obfk_1 (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE TABLE q.a (x int, y int, KEY I (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP INDEX j ON q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE q.A ADD CONSTRAINT a_obfk_1 FOREIGN KEY (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON a FOR EACH ROW SET NEW.x = 1;
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 ENGINE = InnoDB;
< ok 0 0
> exec
| DROP TABLE IF EXISTS q.v, b;
< ok 0 0
> exec
| ALTER TABLE q.a_obfk_1 ADD PRIMARY KEY (x);
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| ALTER TABLE b DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE a ADD UNIQUE INDEX I (x);
< ok 0 0
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE UNIQUE INDEX j ON A (X);
< ok 0 0
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| INSERT INTO A (x) VALUES (1);
< error
| 1364 HY000 Field 'id' doesn't have a default value
> exec
| SET GLOBAL foreign_key_checks = 1, @@session.foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE q.a RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| TRUNCATE TABLE q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE a DROP COLUMN X;
< ok 0 0
> exec
| CREATE TABLE A AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'a' already exists
> exec
| CREATE TABLE IF NOT EXISTS q.A LIKE a;
< error
| 1049 42000 Unknown database 'q'
> exec
| INSERT INTO A (x) VALUES (1);
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec
| ALTER TABLE q.a_obfk_1 CHANGE z y int;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| INSERT INTO a (x) VALUES (1);
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec
| ALTER TABLE c DROP FOREIGN KEY F;
< error
| 1091 42000 Can't DROP 'F'; check that column/key exists
> exec
| ALTER TABLE b ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change16 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE c DROP INDEX I;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE q.c (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| TRUNCATE TABLE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP VIEW A;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| DROP VIEW a_obfk_1;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1'
> exec
| INSERT INTO q.b (x) VALUES (1);
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| CREATE TABLE q.b AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE b CHANGE y z int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TRIGGER IF EXISTS tr2;
< ok 0 0

== TestWalkThroughSweepOceanBase/change16 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| temporary	shop		b	
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
| index	shop	v	PRIMARY	
| relation	shop		v	BASE TABLE
| temporary	shop		b	
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
| index	shop	v	PRIMARY	
| relation	shop		v	BASE TABLE
| temporary	shop		b	
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
| index	shop	v	PRIMARY	
| relation	shop		v	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	c	x	
| index	shop	v	PRIMARY	
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	x	
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			

== TestWalkThroughSweepOceanBase/change16 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| ALTER TABLE v ADD CONSTRAINT j UNIQUE INDEX i (z);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE v ADD COLUMN X int AFTER x;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| RENAME TABLE a TO q.c;
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP VIEW c;
< error
| 1051 42S02 Unknown table 'shop.c'
> exec
| CREATE TEMPORARY TABLE b (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE UNIQUE INDEX x ON v (z);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE A DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0, foreign_key_checks = 1;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE c DROP COLUMN x, ADD COLUMN y int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP TABLE c;
< error
| 1051 42S02 Unknown table 'shop.c'
> exec
| CREATE TABLE c (x int, y int, KEY x (x));
< ok 0 0
> exec
| DROP TABLE IF EXISTS v, b;
< ok 0 0
> exec
| DROP TABLE c;
< ok 0 0
> exec
| ALTER TABLE q.c DROP INDEX i, ADD INDEX j (x);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TABLE q.a (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change17 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| DROP TABLE a_obfk_1, c;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1,shop.c'
> exec
| DROP TABLE v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| CREATE TABLE q.v (x int, y int, KEY I (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| TRUNCATE TABLE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| RENAME TABLE a_obfk_1 TO v;
< error
| 1017 HY000 File not exist
> exec
| CREATE INDEX i ON q.v (x);
< error
| 1049 42000 Unknown database
> exec
| ALTER TABLE a_obfk_1 ADD COLUMN z int AFTER y;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE INDEX x ON b (x);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE a_obfk_1 AS SELECT 1 AS x;
< ok 1 0
> exec
| DROP VIEW c;
< error
| 1051 42S02 Unknown table 'shop.c'
> exec
| CREATE TEMPORARY TABLE A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON b FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS b LIKE b;
< error
| 1066 42000 Not unique table/alias: 'b'

== TestWalkThroughSweepOceanBase/change17 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a_obfk_1
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a_obfk_1
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
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
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
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
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
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
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
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
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		v	
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
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		v	
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
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		v	
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
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		v	
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
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		v	
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
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		v	
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
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		v	
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
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		v	

== TestWalkThroughSweepOceanBase/change17 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| RENAME TABLE q.A TO a, a TO b;
< error
| 1017 HY000 File not exist
> exec
| CREATE TABLE v (x int, y int, KEY I (x), CONSTRAINT g FOREIGN KEY (y) REFERENCES a (id));
< error
| 1146 42S02 Table doesn't exist
> exec
| CREATE OR REPLACE VIEW b AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE c ADD CONSTRAINT a_obfk_1 FOREIGN KEY (x) REFERENCES v (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| INSERT INTO A (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE v DROP INDEX x, ADD INDEX i (x);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TEMPORARY TABLE q.c (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE IF NOT EXISTS v LIKE c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE VIEW q.c AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| INSERT INTO a (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS q.v (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| TRUNCATE TABLE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TABLE A;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| CREATE UNIQUE INDEX i ON A (y);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE v CHANGE y y int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE v CHANGE z z int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE a_obfk_1 LIKE a;
< error
| 1050 42S01 Table 'a_obfk_1' already exists
> exec
| CREATE TEMPORARY TABLE v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE a DROP INDEX I, ADD INDEX i (x);
< error
| 1091 42000 Can't DROP 'I'; check that column/key exists
> exec
| DROP TABLE a_obfk_1;
< ok 0 0
> exec
| ALTER TABLE q.c DROP FOREIGN KEY a_obfk_1;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE INDEX x ON a (x);
< ok 0 0
> exec
| CREATE TABLE q.a_obfk_1 (x int, y int, KEY j (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TRIGGER tr BEFORE INSERT ON q.a FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE INDEX I ON b (y);
< error
| 1072 42000 Key column 'y' doesn't exist in table
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change18 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| DROP TEMPORARY TABLE q.a_obfk_1, v;
< error
| 1051 42S02 Unknown table 'q.a_obfk_1,shop.v'
> exec
| DROP VIEW a;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| ALTER TABLE a ADD UNIQUE INDEX i (x);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE c DROP FOREIGN KEY f;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE c (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE TABLE b (x int, y int, KEY I (x), CONSTRAINT g FOREIGN KEY a_obfk_1 (y) REFERENCES a (id));
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| ALTER TABLE b ENGINE = InnoDB;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE a_obfk_1 ADD COLUMN y int;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP INDEX I ON q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE A ADD CONSTRAINT j UNIQUE INDEX I (X);
< ok 0 0
> exec
| DROP VIEW b;
< error
| 1051 42S02 Unknown table 'shop.b'
> exec
| ALTER TABLE v DROP INDEX i;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TRIGGER TR BEFORE INSERT ON c FOR EACH ROW SET NEW.x = 1;
< ok 0 0

== TestWalkThroughSweepOceanBase/change18 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| c
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| c
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| c
: bytes
| x
: int64 0
: bytes
| x
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| c
: bytes
| TR
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
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
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
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
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
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
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
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
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	c	PRIMARY	
| index	shop	c	x	
| index	shop	v	x	
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	c	PRIMARY	
| index	shop	c	x	
| index	shop	v	x	
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	c	PRIMARY	
| index	shop	c	x	
| index	shop	v	x	
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	c	PRIMARY	
| index	shop	c	x	
| index	shop	v	x	
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	c	PRIMARY	
| index	shop	c	x	
| index	shop	v	x	
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		TR	
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
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
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
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
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
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
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
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	b	i	
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	b	i	
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	b	i	
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	b	i	
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	b	i	
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	

== TestWalkThroughSweepOceanBase/change18 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| CREATE TABLE q.a AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE IF NOT EXISTS q.A LIKE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP INDEX I ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE v (x int, y int, KEY x (x));
< ok 0 0
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS q.A (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP VIEW IF EXISTS v;
< error
| 1347 HY000 'shop.v' is not VIEW
> exec
| ALTER TABLE b ADD CONSTRAINT a_obfk_1 FOREIGN KEY a_obfk_1 (x) REFERENCES c (id);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TABLE v;
< ok 0 0
> exec
| CREATE TABLE q.a (x int, y int, KEY I (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| TRUNCATE TABLE q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| RENAME TABLE q.a TO a;
< error
| 1017 HY000 File not exist
> exec
| CREATE TABLE b (x int, y int, KEY i (x));
< ok 0 0
> exec
| RENAME TABLE b TO b;
< error
| 1050 42S01 Table 'b' already exists
> exec
| DROP VIEW c;
< error
| 1347 HY000 'shop.c' is not VIEW
> exec
| ALTER TABLE a ADD CONSTRAINT a_obfk_1 FOREIGN KEY (x) REFERENCES q.a_obfk_1 (id);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| INSERT INTO c (x) VALUES (1);
< error
| 1364 HY000 Field 'id' doesn't have a default value
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change19 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE a (x int, y int, KEY I (x));
< ok 0 0
> exec
| RENAME TABLE c TO a, q.A TO A;
< error
| 1050 42S01 Table 'a' already exists
> exec
| CREATE UNIQUE INDEX I ON b (y);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE UNIQUE INDEX I ON q.v (X);
< error
| 1049 42000 Unknown database
> exec
| CREATE TABLE c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE v CHANGE x x int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE q.v RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| CREATE TABLE a (id int PRIMARY KEY, x int, y int);
< error
| 1050 42S01 Table 'a' already exists
> exec
| ALTER TABLE q.b DROP COLUMN z;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE q.a DROP COLUMN X;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE UNIQUE INDEX x ON a (z);
< error
| 1072 42000 Key column 'z' doesn't exist in table

== TestWalkThroughSweepOceanBase/change19 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| a
: bytes
| I
: int64 1
: bytes
| x
: int64 -1
< row
: bytes
| c
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	I	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	I	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	a	I	
| index	shop	a	x	
| index	shop	c	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE

== TestWalkThroughSweepOceanBase/change19 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| CREATE UNIQUE INDEX i ON q.a (X);
< error
| 1049 42000 Unknown database
> exec
| CREATE INDEX x ON A (y);
< ok 0 0
> exec
| ALTER TABLE q.a DROP INDEX x;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS b LIKE q.v;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| INSERT INTO a (x) VALUES (1);
< ok 1 0
> exec
| ALTER TABLE A ADD CONSTRAINT x UNIQUE INDEX i (z);
< error
| 1061 42000 Duplicate key name 'i'
> exec
| ALTER TABLE b ADD UNIQUE INDEX i (x);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE A (x int, y int, KEY i (x));
< error
| 1050 42S01 Table 'a' already exists
> exec
| RENAME TABLE a TO q.a;
< error
| 1049 42000 Unknown database 'q'
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec ^
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| ALTER TABLE c ADD FOREIGN KEY g (x) REFERENCES a (id);
< error
| 5031 HY000 Column not found
> exec
| DROP INDEX x ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TEMPORARY TABLE a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change2 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| DROP INDEX I ON a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE a_obfk_1 DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE A ENGINE = InnoDB;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TEMPORARY TABLE A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP TABLE v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| RENAME TABLE c TO a;
< error
| 1017 HY000 File not exist
> exec
| DROP TABLE IF EXISTS b;
< ok 0 0
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| INSERT INTO A (x) VALUES (1);
< error
| 1364 HY000 Field 'id' doesn't have a default value
> exec
| DROP INDEX j ON a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE a DROP INDEX I, ADD INDEX x (x);
< error
| 1091 42000 Can't DROP 'I'; check that column/key exists
> exec
| ALTER TABLE a_obfk_1 ADD CONSTRAINT I UNIQUE INDEX j (z);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| USE shop;
< ok 0 0

== TestWalkThroughSweepOceanBase/change2 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	q			
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	q			
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	q			
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| temporary	q		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| temporary	q		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| temporary	q		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| temporary	q		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| temporary	q		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| database	q			
| database	shop			
| temporary	q		a	
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
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
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
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
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
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
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
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
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
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
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
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
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
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
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
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
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
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
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
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
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
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
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
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
| temporary	q		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
| temporary	q		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
| temporary	q		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
| temporary	q		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| temporary	q		a	
| temporary	q		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	q		a	
| temporary	q		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	q			
| database	shop			
| index	shop	c	PRIMARY	
| index	shop	v	PRIMARY	
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	q		a	
| temporary	q		v	

== TestWalkThroughSweepOceanBase/change2 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| CREATE TRIGGER tr BEFORE INSERT ON A FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN z int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TEMPORARY TABLE q.a (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| RENAME TABLE q.a_obfk_1 TO a;
< error
| 1017 HY000 File not exist
> exec
| ALTER TABLE A CHANGE x z int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TRIGGER tr BEFORE INSERT ON c FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP TEMPORARY TABLE q.v, A;
< error
| 1051 42S02 Unknown table 'q.v,shop.a'
> exec
| TRUNCATE TABLE A;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 CHANGE x y int;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE c DROP INDEX j, ADD INDEX i (x);
< error
| 1091 42000 Can't DROP 'j'; check that column/key exists
> exec
| DROP VIEW IF EXISTS A;
< ok 0 0
> exec
| CREATE TABLE b (x int, y int, KEY j (x), FOREIGN KEY g (y) REFERENCES a (id));
< error
| 1146 42S02 Table doesn't exist
> exec
| DROP TABLE a, b;
< error
| 1051 42S02 Unknown table 'shop.a,shop.b'
> exec
| RENAME TABLE c TO c, v TO v;
< error
| 1050 42S01 Table 'c' already exists
> exec
| ALTER TABLE q.a ADD COLUMN X int;
< error
| 1060 42S21 Duplicate column name 'X'
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| CREATE TABLE q.a AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'a' already exists
> exec
| ALTER TABLE b CHANGE x x int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP INDEX j ON a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TEMPORARY TABLE q.v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE b DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| RENAME TABLE q.c TO a_obfk_1, A TO v;
< error
| 1017 HY000 File not exist
> exec
| DROP TABLE IF EXISTS q.a_obfk_1, A;
< ok 0 0
> exec
| CREATE TABLE a_obfk_1 (x int, y int, KEY i (x), CONSTRAINT a_obfk_1 FOREIGN KEY a_obfk_1 (y) REFERENCES a (id));
< error
| 1146 42S02 Table doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr2;
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS b, A;
< ok 0 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change20 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr2;
< ok 0 0
> exec
| DROP TABLE IF EXISTS b;
< ok 0 0
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| ALTER TABLE c ADD PRIMARY KEY (x);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE q.b ADD COLUMN X int AFTER x;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| DROP TRIGGER tr2;
< error
| 1360 HY000 Trigger does not exist
> exec
| DROP INDEX j ON c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE c ADD CONSTRAINT f FOREIGN KEY F (x) REFERENCES v (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE q.A ADD CONSTRAINT i UNIQUE INDEX i (z);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE b DROP INDEX j;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE INDEX j ON b (z);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE q.b ADD COLUMN x int AFTER X;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE q.c CHANGE z x int;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DROP INDEX j ON c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist

== TestWalkThroughSweepOceanBase/change20 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	VIEW

== TestWalkThroughSweepOceanBase/change20 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| CREATE VIEW c AS SELECT 1 AS x;
< ok 0 0
> exec
| DROP VIEW q.c;
< error
| 1051 42S02 Unknown table 'q.c'
> exec
| DROP TABLE a;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| CREATE INDEX j ON A (x);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP TABLE IF EXISTS a_obfk_1;
< ok 0 0
> exec
| DROP TABLE A, a_obfk_1;
< error
| 1051 42S02 Unknown table 'shop.a,shop.a_obfk_1'
> exec
| DROP TRIGGER TR;
< error
| 1360 HY000 Trigger does not exist
> exec
| INSERT INTO a_obfk_1 (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE A ADD CONSTRAINT f FOREIGN KEY a_obfk_1 (x) REFERENCES b (id);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE a_obfk_1 AS SELECT 1 AS x;
< ok 1 0
> exec
| CREATE TABLE A AS SELECT 1 AS x;
< ok 1 0
> exec
| DROP TABLE IF EXISTS a_obfk_1;
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS a_obfk_1;
< ok 0 0
> exec
| ALTER TABLE c DROP COLUMN x;
< error
| 1347 HY000 'shop.c' is not BASE TABLE
> exec
| CREATE INDEX I ON b (x);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE a (x int, y int, KEY i (x), CONSTRAINT g FOREIGN KEY g (y) REFERENCES a (id));
< error
| 5031 HY000 Column not found
> exec
| ALTER TABLE q.a DROP INDEX x, ADD INDEX I (x);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP TRIGGER IF EXISTS tr2;
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN x int;
< error
| 1060 42S21 Duplicate column name 'x'
> exec
| ALTER TABLE A ADD CONSTRAINT g FOREIGN KEY F (x) REFERENCES q.a (id);
< error
| 1146 42S02 Table doesn't exist
> exec
| DROP TRIGGER IF EXISTS tr2;
< ok 0 0
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| CREATE TRIGGER tr BEFORE INSERT ON c FOR EACH ROW SET NEW.x = 1;
< error
| 1235 0A000 simple dml trigger isn't used on user table not supported
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| ALTER TABLE q.v DROP INDEX j;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| DROP TABLE c;
< error
| 1051 42S02 Unknown table 'shop.c'
> exec
| ALTER TABLE v CHANGE z z int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE a_obfk_1 CHANGE y z int;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change21 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON c FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE a_obfk_1 CHANGE x x int;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP TRIGGER TR;
< error
| 1360 HY000 Trigger does not exist
> exec
| CREATE VIEW a AS SELECT 1 AS x;
< ok 0 0
> exec
| RENAME TABLE c TO a_obfk_1;
< error
| 1017 HY000 File not exist
> exec
| ALTER TABLE q.v ADD CONSTRAINT a_obfk_1 FOREIGN KEY (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| DROP TABLE IF EXISTS q.v, a_obfk_1;
< ok 0 0
> exec
| ALTER TABLE q.A DROP INDEX i, ADD INDEX I (x);
< error
| 1146 42S02 Table 'q.a' doesn't exist

== TestWalkThroughSweepOceanBase/change21 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| VIEW
: bytes
| MEMORY
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	a	x	
| database	q			
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	a	x	
| database	q			
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	a	x	
| database	q			
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	a	x	
| database	q			
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	a	x	
| database	q			
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	a	x	
| database	q			
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	a	x	
| database	q			
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	a	x	
| database	q			
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	a	x	
| database	q			
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		a	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		a	VIEW
| relation	shop		c	BASE TABLE
> query ^
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
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		c	BASE TABLE
> query ^
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
| index	shop	c	PRIMARY	
| index	shop	c	x	
| relation	shop		c	BASE TABLE

== TestWalkThroughSweepOceanBase/change21 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| ALTER TABLE a CHANGE y X int;
< error
| 1054 42S22 Unknown column 'y' in 'a'
> exec
| DROP VIEW q.a;
< error
| 1051 42S02 Unknown table 'q.a'
> exec
| DROP INDEX i ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE q.c ADD CONSTRAINT f FOREIGN KEY a_obfk_1 (x) REFERENCES A (id);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TABLE q.A (x int, y int, KEY x (x), CONSTRAINT g FOREIGN KEY (y) REFERENCES a (id));
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE UNIQUE INDEX I ON b (x);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS v, q.b;
< ok 0 0
> exec
| DROP INDEX i ON q.c;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE c DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP TABLE a_obfk_1;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1'
> exec
| INSERT INTO c (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE A ADD UNIQUE INDEX j (X);
< error
| 1347 HY000 'shop.a' is not BASE TABLE
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| TRUNCATE TABLE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE INDEX i ON q.a_obfk_1 (X);
< error
| 1049 42000 Unknown database
> exec
| CREATE TABLE c (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| DROP INDEX j ON c;
< error
| 1091 42000 Can't DROP 'j'; check that column/key exists
> exec
| DROP VIEW A;
< ok 0 0
> exec
| DROP TRIGGER tr2;
< error
| 1360 HY000 Trigger does not exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change22 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1, @@session.foreign_key_checks = 0;
< ok 0 0
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| ALTER TABLE b DROP COLUMN y;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TRIGGER IF EXISTS TR;
< ok 0 0
> exec
| DROP TEMPORARY TABLE q.a;
< error
| 1051 42S02 Unknown table 'q.a'
> exec
| TRUNCATE TABLE q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| INSERT INTO c (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE INDEX I ON a (y);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE VIEW q.A AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP TABLE IF EXISTS v, a;
< ok 0 0
> exec
| CREATE VIEW q.c AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON b FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist

== TestWalkThroughSweepOceanBase/change22 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	a_obfk_1	x	
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	

== TestWalkThroughSweepOceanBase/change22 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| SET GLOBAL foreign_key_checks = 0, foreign_key_checks = 1;
< ok 0 0
> exec
| DROP INDEX x ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TRIGGER tr BEFORE INSERT ON A FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| DROP VIEW A;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| ALTER TABLE q.b DROP INDEX x;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE q.b DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE b DROP COLUMN y, ADD COLUMN y int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS a (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE TABLE a_obfk_1 (x int, y int, KEY i (x), CONSTRAINT f FOREIGN KEY (y) REFERENCES a (id));
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| SET GLOBAL foreign_key_checks = 0, foreign_key_checks = 1;
< ok 0 0
> exec
| CREATE TABLE a_obfk_1 (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE TEMPORARY TABLE a (id int PRIMARY KEY, x int, y int);
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP VIEW IF EXISTS q.A;
< ok 0 0
> exec
| CREATE OR REPLACE VIEW q.a AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE q.c ADD CONSTRAINT f FOREIGN KEY (x) REFERENCES q.v (id);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS a_obfk_1 (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP TABLE IF EXISTS A;
< ok 0 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change23 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE b (x int, y int, KEY j (x));
< ok 0 0
> exec
| TRUNCATE TABLE a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| INSERT INTO a (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| CREATE TABLE v AS SELECT 1 AS x;
< ok 1 0
> exec
| ALTER TABLE b DROP FOREIGN KEY g;
< error
| 1091 42000 Can't DROP 'g'; check that column/key exists
> exec
| INSERT INTO b (x) VALUES (1);
< ok 1 0
> exec
| DROP INDEX x ON q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE INDEX I ON a_obfk_1 (z);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE a_obfk_1 CHANGE y x int;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE b DROP FOREIGN KEY g;
< error
| 1091 42000 Can't DROP 'g'; check that column/key exists

== TestWalkThroughSweepOceanBase/change23 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| b
: bytes
| j
: int64 1
: bytes
| x
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	b	j	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepOceanBase/change23 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| ALTER TABLE a ADD CONSTRAINT j UNIQUE INDEX i (y);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE q.a DROP INDEX I;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE v (x int, y int, KEY x (x));
< error
| 1050 42S01 Table 'v' already exists
> exec
| CREATE TABLE c AS SELECT 1 AS x;
< ok 1 0
> exec
| ALTER TABLE v DROP FOREIGN KEY g;
< error
| 1091 42000 Can't DROP 'g'; check that column/key exists
> exec
| CREATE TEMPORARY TABLE v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP INDEX I ON a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE q.b ADD COLUMN x int;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE A DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| TRUNCATE TABLE q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP INDEX x ON a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE v ADD UNIQUE INDEX x (x);
< ok 0 0
> exec
| DROP TABLE v;
< ok 0 0
> exec
| DROP VIEW q.a_obfk_1;
< error
| 1051 42S02 Unknown table 'q.a_obfk_1'
> exec
| ALTER TABLE A CHANGE X X int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec
| ALTER TABLE q.c CHANGE z x int;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change24 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE INDEX i ON c (z);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP TRIGGER tr2;
< error
| 1360 HY000 Trigger does not exist
> exec
| RENAME TABLE v TO a, A TO a;
< error
| 1017 HY000 File not exist
> exec
| ALTER TABLE q.a DROP INDEX j, ADD INDEX I (x);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE OR REPLACE VIEW a_obfk_1 AS SELECT 1 AS x;
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr;
< ok 0 0
> exec
| ALTER TABLE c ADD FOREIGN KEY F (x) REFERENCES q.v (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE VIEW c AS SELECT 1 AS x;
< ok 0 0
> exec
| INSERT INTO q.v (x) VALUES (1);
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| CREATE TABLE q.c LIKE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE A AS SELECT 1 AS x;
< ok 1 0

== TestWalkThroughSweepOceanBase/change24 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| a_obfk_1
: bytes
| VIEW
: bytes
| MEMORY
< row
: bytes
| c
: bytes
| VIEW
: bytes
| MEMORY
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| x
< row
: bytes
| a_obfk_1
: bytes
| x
< row
: bytes
| c
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		c	VIEW
| relation	shop		v	VIEW
| temporary	shop		a_obfk_1	

== TestWalkThroughSweepOceanBase/change24 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| CREATE OR REPLACE VIEW q.c AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS a_obfk_1 (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE INDEX I ON b (x);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE VIEW q.a_obfk_1 AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TRIGGER tr BEFORE INSERT ON q.b FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| RENAME TABLE a TO A, b TO c;
< error
| 1050 42S01 Table 'a' already exists
> exec
| CREATE VIEW v AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE q.c ADD INDEX x (x);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| ALTER TABLE b DROP INDEX I;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE OR REPLACE VIEW c AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE A CHANGE y x int;
< error
| 1054 42S22 Unknown column 'y' in 'a'
> exec
| ALTER TABLE c RENAME TO q.A;
< error
| 1347 HY000 'shop.c' is not BASE TABLE
> exec
| ALTER TABLE a ADD INDEX x (y);
< error
| 1072 42000 Key column 'y' doesn't exist in table
> exec
| CREATE OR REPLACE VIEW c AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE a ADD INDEX x (z);
< error
| 1072 42000 Key column 'z' doesn't exist in table
> exec
| CREATE TABLE q.A (x int, y int, KEY i (x), CONSTRAINT g FOREIGN KEY g (y) REFERENCES a (id));
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE IF NOT EXISTS v LIKE q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE A CHANGE X y int;
< ok 0 0
> exec
| ALTER TABLE q.c DROP COLUMN x, ADD COLUMN y int;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE a_obfk_1 DROP FOREIGN KEY f;
< error
| 1091 42000 Can't DROP 'f'; check that column/key exists
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change25 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE q.A ADD CONSTRAINT i UNIQUE INDEX i (y);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE a DROP FOREIGN KEY g;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE a_obfk_1 DROP FOREIGN KEY g;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE a DROP FOREIGN KEY a_obfk_1;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE v ADD PRIMARY KEY (x);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TRIGGER TR BEFORE INSERT ON q.A FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| CREATE INDEX I ON a (X);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| CREATE INDEX i ON A (X);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE a_obfk_1 DROP INDEX I;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TABLE c (x int, y int, KEY I (x), CONSTRAINT a_obfk_1 FOREIGN KEY (y) REFERENCES a (id));
< error
| 1146 42S02 Table doesn't exist

== TestWalkThroughSweepOceanBase/change25 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	shop			
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	q			
| database	shop			
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| database	q			
| database	shop			
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW

== TestWalkThroughSweepOceanBase/change25 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| ALTER TABLE c ADD CONSTRAINT i UNIQUE INDEX j (z);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE b DROP COLUMN X;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE v ENGINE = InnoDB;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP TABLE IF EXISTS b;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS A LIKE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE A ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| TRUNCATE TABLE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE a (x int, y int, KEY i (x), FOREIGN KEY F (y) REFERENCES a (id));
< error
| 5031 HY000 Column not found
> exec
| CREATE VIEW A AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE OR REPLACE VIEW a_obfk_1 AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE b DROP INDEX i;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE a_obfk_1 (x int, y int, KEY I (x));
< error
| 1050 42S01 Table 'a_obfk_1' already exists
> exec
| DROP TABLE q.a;
< error
| 1051 42S02 Unknown table 'q.a'
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON q.b FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change26 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| INSERT INTO q.c (x) VALUES (1);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE OR REPLACE VIEW b AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE TABLE b (x int, y int, KEY i (x));
< error
| 1050 42S01 Table 'b' already exists
> exec
| CREATE TABLE c AS SELECT 1 AS x;
< ok 1 0
> exec
| CREATE TABLE IF NOT EXISTS c LIKE q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE UNIQUE INDEX j ON q.c (z);
< error
| 1049 42000 Unknown database
> exec
| DROP TABLE a;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON b FOR EACH ROW SET NEW.x = 1;
< error
| 1235 0A000 simple dml trigger isn't used on user table not supported

== TestWalkThroughSweepOceanBase/change26 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| VIEW
: bytes
| MEMORY
< row
: bytes
| c
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
| MEMORY
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| b
: bytes
| x
< row
: bytes
| c
: bytes
| x
< row
: bytes
| v
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	a	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	a	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	a	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	a	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	a	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	a	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	a	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	a	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	a	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	x	
| column	shop	v	x	
| database	shop			
| index	shop	a	j	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| temporary	shop		a	

== TestWalkThroughSweepOceanBase/change26 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| DROP TABLE A;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| ALTER TABLE v ADD CONSTRAINT F FOREIGN KEY (x) REFERENCES c (id);
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| ALTER TABLE v ADD CONSTRAINT g FOREIGN KEY (x) REFERENCES c (id);
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TEMPORARY TABLE A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE b AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'b' already exists
> exec
| INSERT INTO a (x) VALUES (1);
< error
| 1364 HY000 Field 'id' doesn't have a default value
> exec
| CREATE TABLE IF NOT EXISTS b (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| ALTER TABLE A ADD UNIQUE INDEX x (X);
< ok 0 0
> exec
| ALTER TABLE b CHANGE x x int;
< error
| 1347 HY000 'shop.b' is not BASE TABLE
> exec
| DROP INDEX x ON a;
< ok 0 0
> exec
| DROP VIEW IF EXISTS q.c;
< ok 0 0
> exec
| CREATE UNIQUE INDEX I ON v (z);
< error
| 1072 42000 Key column 'z' doesn't exist in table
> exec
| CREATE TABLE a (x int, y int, KEY j (x));
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN X int;
< error
| 1060 42S21 Duplicate column name 'X'
> exec
| ALTER TABLE A DROP INDEX I, ADD INDEX I (x);
< error
| 1091 42000 Can't DROP 'I'; check that column/key exists
> exec
| RENAME TABLE q.b TO A;
< error
| 1050 42S01 Table 'a' already exists
> exec
| TRUNCATE TABLE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS q.a_obfk_1 LIKE b;
< error
| 1347 HY000 'shop.b' is not BASE TABLE
> exec
| TRUNCATE TABLE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE VIEW v AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'v' already exists
> exec
| ALTER TABLE b CHANGE x x int;
< error
| 1347 HY000 'shop.b' is not BASE TABLE
> exec
| DROP TABLE a_obfk_1, A;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1'
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change27 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE b ADD COLUMN y int AFTER y;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE A AS SELECT 1 AS x;
< ok 1 0
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| USE shop;
< ok 0 0
> exec
| ALTER TABLE v DROP COLUMN y, ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE UNIQUE INDEX j ON a_obfk_1 (z);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP VIEW q.a;
< error
| 1051 42S02 Unknown table 'q.a'
> exec
| ALTER TABLE c DROP FOREIGN KEY f;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE v DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE a LIKE c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE A ADD CONSTRAINT g FOREIGN KEY (x) REFERENCES a_obfk_1 (id);
< error
| 1146 42S02 Table doesn't exist

== TestWalkThroughSweepOceanBase/change27 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| b
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| b
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	b	X	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	b	X	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	b	X	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	b	X	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	b	X	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	b	X	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 13
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 18
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepOceanBase/change27 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| ALTER TABLE b RENAME TO q.b;
< error
| 1210 HY000 Invalid argument
> exec
| ALTER TABLE b CHANGE X X int;
< ok 0 0
> exec
| DROP TABLE IF EXISTS q.a, c;
< ok 0 0
> exec
| DROP TABLE q.a_obfk_1, v;
< error
| 1051 42S02 Unknown table 'q.a_obfk_1,shop.v'
> exec
| ALTER TABLE a RENAME INDEX i TO j;
< error
| 1176 42000 Key 'i' doesn't exist in table 'a'
> exec
| CREATE TABLE a (x int, y int, KEY j (x), CONSTRAINT g FOREIGN KEY f (y) REFERENCES a (id));
< error
| 5031 HY000 Column not found
> exec
| ALTER TABLE b ADD COLUMN x int;
< error
| 1060 42S21 Duplicate column name 'x'
> exec
| ALTER TABLE b DROP COLUMN X;
< ok 0 0
> exec
| CREATE INDEX x ON A (z);
< error
| 1072 42000 Key column 'z' doesn't exist in table
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON v FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP TRIGGER IF EXISTS TR;
< ok 0 0
> exec
| ALTER TABLE c DROP FOREIGN KEY f;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| INSERT INTO q.b (x) VALUES (1);
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE b DROP INDEX j, ADD INDEX i (x);
< error
| 1091 42000 Can't DROP 'j'; check that column/key exists
> exec
| CREATE TABLE v (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE TABLE A (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1050 42S01 Table 'a' already exists
> exec
| CREATE TABLE a_obfk_1 (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE INDEX i ON c (y);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec
| USE shop;
< ok 0 0
> exec
| CREATE OR REPLACE VIEW A AS SELECT 1 AS x;
< error
| 1347 HY000 'shop.a' is not VIEW
> exec
| RENAME TABLE q.a TO A;
< error
| 1050 42S01 Table 'a' already exists
> exec
| CREATE INDEX j ON A (y);
< error
| 1072 42000 Key column 'y' doesn't exist in table
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change28 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| DROP INDEX j ON A;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TRIGGER tr BEFORE INSERT ON c FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON a_obfk_1 FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP VIEW c;
< error
| 1051 42S02 Unknown table 'shop.c'
> exec
| ALTER TABLE q.b DROP COLUMN x;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE q.v DROP INDEX j, ADD INDEX x (x);
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| ALTER TABLE A CHANGE X z int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS a_obfk_1 (id int PRIMARY KEY, x int, y int);
< ok 0 0

== TestWalkThroughSweepOceanBase/change28 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 2
: bytes
| database	shop			
| temporary	shop		a	
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
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| temporary	shop		a	
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
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| temporary	shop		a	
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
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| temporary	shop		a	
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
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| temporary	shop		a	
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
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| temporary	shop		a	
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
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| temporary	shop		a	
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
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| temporary	shop		a	
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
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| temporary	shop		a	
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
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| temporary	shop		a	
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
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| temporary	shop		a	
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
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| temporary	shop		a	
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
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	b	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| temporary	shop		a	
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
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	b	x	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| temporary	shop		a	

== TestWalkThroughSweepOceanBase/change28 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| ALTER TABLE a DROP INDEX x;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TEMPORARY TABLE a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE b (x int, y int, KEY x (x));
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE A CHANGE z z int;
< error
| 1054 42S22 Unknown column 'z' in 'a'
> exec
| RENAME TABLE A TO v;
< error
| 1017 HY000 File not exist
> exec
| DROP TABLE IF EXISTS v;
< ok 0 0
> exec
| ALTER TABLE q.c RENAME TO b;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DROP TEMPORARY TABLE q.A;
< error
| 1051 42S02 Unknown table 'q.a'
> exec
| ALTER TABLE a ADD FOREIGN KEY f (x) REFERENCES c (id);
< error
| 1146 42S02 Table doesn't exist
> exec
| INSERT INTO q.a_obfk_1 (x) VALUES (1);
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| ALTER TABLE b DROP INDEX I, ADD INDEX j (x);
< error
| 1091 42000 Can't DROP 'I'; check that column/key exists
> exec
| DROP VIEW v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| CREATE TABLE a (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 ADD INDEX i (x);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change29 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| CREATE UNIQUE INDEX j ON A (y);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE q.c DROP FOREIGN KEY g;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE INDEX j ON v (z);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TRIGGER tr BEFORE INSERT ON a FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE b ADD COLUMN z int AFTER z;
< error
| 1054 42S22 Unknown column 'z' in 'b'
> exec
| CREATE TABLE b LIKE q.a_obfk_1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| ALTER TABLE A DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| TRUNCATE TABLE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE b ADD CONSTRAINT F FOREIGN KEY a_obfk_1 (x) REFERENCES v (id);
< error
| 1146 42S02 Table doesn't exist
> exec
| ALTER TABLE q.a_obfk_1 ADD INDEX x (y);
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist

== TestWalkThroughSweepOceanBase/change29 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| b
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| b
: bytes
| x
: int64 0
: bytes
| x
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
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
| index	shop	b	PRIMARY	
| index	shop	b	x	
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
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| index	shop	b	PRIMARY	
| index	shop	b	x	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE

== TestWalkThroughSweepOceanBase/change29 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| TRUNCATE TABLE A;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE c AS SELECT 1 AS x;
< ok 1 0
> exec
| ALTER TABLE q.b DROP COLUMN X;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| TRUNCATE TABLE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE c ENGINE = InnoDB;
< ok 0 0
> exec
| TRUNCATE TABLE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE A CHANGE y y int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE VIEW q.a_obfk_1 AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE IF NOT EXISTS a_obfk_1 (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE INDEX i ON a (z);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE c (id int PRIMARY KEY, x int, y int);
< error
| 1050 42S01 Table 'c' already exists
> exec
| TRUNCATE TABLE A;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS q.b (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP TABLE q.a_obfk_1, a;
< error
| 1051 42S02 Unknown table 'q.a_obfk_1,shop.a'
> exec
| CREATE TRIGGER TR BEFORE INSERT ON q.b FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| CREATE INDEX j ON A (y);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON a_obfk_1 FOR EACH ROW SET NEW.x = 1;
< ok 0 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change3 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< ok 1 0
> exec
| CREATE UNIQUE INDEX x ON q.a (y);
< error
| 1049 42000 Unknown database
> exec
| ALTER TABLE v ADD COLUMN z int AFTER X;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP INDEX j ON q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0, foreign_key_checks = 1;
< ok 0 0
> exec
| CREATE TABLE a (x int, y int, KEY i (x));
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP VIEW b;
< error
| 1051 42S02 Unknown table 'shop.b'
> exec
| CREATE VIEW b AS SELECT 1 AS x;
< ok 0 0

== TestWalkThroughSweepOceanBase/change3 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| b
: bytes
| VIEW
: bytes
| MEMORY
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| x
< row
: bytes
| b
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	b	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	I	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	I	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	I	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	I	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	I	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	I	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	I	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	I	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	b	x	
| database	q			
| database	shop			
| index	shop	a	I	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	b	x	
| database	q			
| database	shop			
| index	shop	a	I	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	b	x	
| database	q			
| database	shop			
| relation	shop		b	VIEW
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
| column	shop	b	x	
| database	q			
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
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
| column	shop	b	x	
| database	q			
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
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
| column	shop	b	x	
| database	q			
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
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
| column	shop	b	x	
| database	q			
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
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
| column	shop	b	x	
| database	q			
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
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
| column	shop	b	x	
| database	q			
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
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
| column	shop	b	x	
| database	q			
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW
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
| column	shop	b	x	
| database	q			
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		b	VIEW

== TestWalkThroughSweepOceanBase/change3 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| DROP TRIGGER IF EXISTS tr;
< ok 0 0
> exec
| CREATE INDEX I ON a (x);
< ok 0 0
> exec
| ALTER TABLE b CHANGE z x int;
< error
| 1054 42S22 Unknown column 'z' in 'b'
> exec
| ALTER TABLE v ADD CONSTRAINT j UNIQUE INDEX i (y);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| INSERT INTO c (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE INDEX i ON a_obfk_1 (X);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP TEMPORARY TABLE v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| RENAME TABLE b TO b;
< error
| 1050 42S01 Table 'b' already exists
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| ALTER TABLE c RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP TABLE a;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP INDEX x ON q.b;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| CREATE INDEX j ON a_obfk_1 (x);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| TRUNCATE TABLE c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int, y int);
< error
| 1050 42S01 Table 'b' already exists
> exec
| RENAME TABLE a TO q.v, q.a TO q.b;
< error
| 1017 HY000 File not exist
> exec
| ALTER TABLE a_obfk_1 DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change30 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE q.v ADD COLUMN x int;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| ALTER TABLE c CHANGE x z int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE c ADD CONSTRAINT F FOREIGN KEY (x) REFERENCES v (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0, foreign_key_checks = 1;
< ok 0 0
> exec
| DROP TABLE a, q.a;
< error
| 1051 42S02 Unknown table 'shop.a,q.a'
> exec
| ALTER TABLE A ADD CONSTRAINT a_obfk_1 FOREIGN KEY g (x) REFERENCES q.v (id);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP VIEW IF EXISTS a_obfk_1;
< ok 0 0
> exec
| ALTER TABLE v ADD UNIQUE INDEX i (X);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP TRIGGER TR;
< error
| 1360 HY000 Trigger does not exist
> exec
| ALTER TABLE q.c ADD COLUMN z int;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| TRUNCATE TABLE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP INDEX j ON a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist

== TestWalkThroughSweepOceanBase/change30 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
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
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
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
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
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
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
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
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
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
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
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
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
| index	shop	c	i	
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
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
| index	shop	c	i	
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
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
| index	shop	c	i	
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
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
| index	shop	c	i	
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
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
| index	shop	c	i	
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
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	c	PRIMARY	
| index	shop	c	i	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE

== TestWalkThroughSweepOceanBase/change30 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| ALTER TABLE q.a DROP INDEX x, ADD INDEX i (x);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| TRUNCATE TABLE A;
< ok 0 0
> exec
| ALTER TABLE A ADD CONSTRAINT a_obfk_1 FOREIGN KEY a_obfk_1 (x) REFERENCES v (id);
< error
| 1146 42S02 Table doesn't exist
> exec
| DROP TRIGGER tr2;
< error
| 1360 HY000 Trigger does not exist
> exec
| CREATE TEMPORARY TABLE q.c (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE INDEX x ON c (x);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS c LIKE a;
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| DROP VIEW q.v;
< error
| 1051 42S02 Unknown table 'q.v'
> exec
| DROP TABLE q.c;
< error
| 1051 42S02 Unknown table 'q.c'
> exec
| ALTER TABLE c ADD CONSTRAINT x UNIQUE INDEX i (X);
< ok 0 0
> exec
| ALTER TABLE q.c ADD CONSTRAINT F FOREIGN KEY (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE TRIGGER tr BEFORE INSERT ON v FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP TABLE q.a_obfk_1;
< error
| 1051 42S02 Unknown table 'q.a_obfk_1'
> exec
| CREATE TABLE q.a_obfk_1 AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE v ADD CONSTRAINT F FOREIGN KEY f (x) REFERENCES c (id);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 1, @@session.foreign_key_checks = 0;
< ok 0 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change31 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE c AS SELECT 1 AS x;
< ok 1 0
> exec
| INSERT INTO q.A (x) VALUES (1);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE a_obfk_1 ADD UNIQUE INDEX x (z);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE q.A ADD INDEX x (y);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE v CHANGE z z int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE a DROP INDEX I, ADD INDEX x (x);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| TRUNCATE TABLE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE a CHANGE X z int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP VIEW a_obfk_1;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1'
> exec
| CREATE TABLE a (x int, y int, KEY I (x));
< ok 0 0
> exec
| ALTER TABLE c DROP FOREIGN KEY a_obfk_1;
< error
| 1091 42000 Can't DROP 'a_obfk_1'; check that column/key exists
> exec
| INSERT INTO c (x) VALUES (1);
< ok 1 0
> exec
| CREATE INDEX I ON c (x);
< ok 0 0
> exec
| TRUNCATE TABLE a;
< ok 0 0
> exec
| SET foreign_key_checks = 0;
< ok 0 0

== TestWalkThroughSweepOceanBase/change31 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| a
: bytes
| I
: int64 1
: bytes
| x
: int64 -1
< row
: bytes
| c
: bytes
| I
: int64 1
: bytes
| x
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	a	I	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	c	x	
| database	shop			
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		tr	

== TestWalkThroughSweepOceanBase/change31 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| DROP VIEW a_obfk_1;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1'
> exec
| ALTER TABLE q.A CHANGE x y int;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| ALTER TABLE A ADD CONSTRAINT F FOREIGN KEY f (x) REFERENCES q.b (id);
< error
| 1146 42S02 Table doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| DROP TABLE IF EXISTS q.a_obfk_1, b;
< ok 0 0
> exec
| DROP VIEW IF EXISTS q.A;
< ok 0 0
> exec
| TRUNCATE TABLE A;
< ok 0 0
> exec
| DROP VIEW a_obfk_1;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1'
> exec
| DROP VIEW A;
< error
| 1347 HY000 'shop.a' is not VIEW
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE a_obfk_1 CHANGE X X int;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE q.a DROP INDEX j;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS q.a_obfk_1 (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP INDEX i ON A;
< ok 0 0
> exec
| RENAME TABLE q.a_obfk_1 TO c, b TO a;
< error
| 1050 42S01 Table 'c' already exists
> exec
| CREATE TRIGGER tr BEFORE INSERT ON A FOR EACH ROW SET NEW.x = 1;
< ok 0 0
> exec
| CREATE TABLE b AS SELECT 1 AS x;
< ok 1 0
> exec
| CREATE INDEX i ON A (y);
< ok 0 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change32 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE q.A (x int, y int, KEY x (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE c DROP INDEX j;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TEMPORARY TABLE A (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| CREATE UNIQUE INDEX I ON a_obfk_1 (y);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP INDEX I ON q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TRIGGER TR BEFORE INSERT ON v FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON q.b FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| ALTER TABLE A ADD COLUMN z int;
< ok 0 0
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 1, @@session.foreign_key_checks = 0;
< ok 0 0

== TestWalkThroughSweepOceanBase/change32 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| b
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| b
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
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
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| temporary	shop		b	
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
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| temporary	shop		b	
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
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| temporary	shop		b	
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
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| temporary	shop		b	
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
| index	shop	b	PRIMARY	
| relation	shop		b	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	v	x	
| database	shop			
| index	shop	b	PRIMARY	
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		v	BASE TABLE
| temporary	shop		b	

== TestWalkThroughSweepOceanBase/change32 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| ALTER TABLE c ADD CONSTRAINT i UNIQUE INDEX i (X);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| RENAME TABLE a_obfk_1 TO a_obfk_1, q.A TO q.A;
< error
| 1017 HY000 File not exist
> exec
| CREATE TABLE q.v (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE v LIKE q.a_obfk_1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| ALTER TABLE a_obfk_1 ADD CONSTRAINT f FOREIGN KEY (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP TEMPORARY TABLE v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| ALTER TABLE c DROP INDEX I;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP INDEX I ON a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE c DROP COLUMN z, ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP INDEX i ON q.a_obfk_1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| INSERT INTO a (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE a_obfk_1 AS SELECT 1 AS x;
< ok 1 0
> exec
| TRUNCATE TABLE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE v AS SELECT 1 AS x;
< ok 1 0
> exec
| ALTER TABLE c ADD CONSTRAINT j UNIQUE INDEX i (y);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| RENAME TABLE A TO q.b, q.b TO a_obfk_1;
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP TRIGGER IF EXISTS tr;
< ok 0 0
> exec
| ALTER TABLE A DROP FOREIGN KEY f;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE c DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| USE shop;
< ok 0 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change33 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE a DROP FOREIGN KEY f;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE v CHANGE X x int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE c CHANGE z z int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE OR REPLACE VIEW q.v AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP TABLE a_obfk_1, c;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1,shop.c'
> exec
| ALTER TABLE q.v ADD COLUMN X int;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| CREATE TRIGGER tr BEFORE INSERT ON A FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| TRUNCATE TABLE a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| DROP VIEW a;
< error
| 1051 42S02 Unknown table 'shop.a'

== TestWalkThroughSweepOceanBase/change33 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a_obfk_1	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	VIEW

== TestWalkThroughSweepOceanBase/change33 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| SET GLOBAL foreign_key_checks = 1, @@session.foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE q.b (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| TRUNCATE TABLE q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| RENAME TABLE b TO b;
< error
| 1017 HY000 File not exist
> exec
| CREATE VIEW c AS SELECT 1 AS x;
< ok 0 0
> exec
| DROP TABLE v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| ALTER TABLE v DROP INDEX j, ADD INDEX j (x);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE a_obfk_1 AS SELECT 1 AS x;
< ok 1 0
> exec
| CREATE TEMPORARY TABLE v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE q.a ADD INDEX x (z);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP TABLE IF EXISTS v;
< ok 0 0
> exec
| CREATE TABLE q.v LIKE A;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS q.a_obfk_1, a;
< ok 0 0
> exec
| ALTER TABLE A DROP INDEX x, ADD INDEX j (x);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP TEMPORARY TABLE q.a;
< error
| 1051 42S02 Unknown table 'q.a'
> exec
| CREATE TEMPORARY TABLE q.c (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE c (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1050 42S01 Table 'c' already exists
> exec
| DROP TEMPORARY TABLE a, b;
< error
| 1051 42S02 Unknown table 'shop.a,shop.b'
> exec
| CREATE INDEX x ON q.v (z);
< error
| 1049 42000 Unknown database
> exec
| DROP VIEW b;
< error
| 1051 42S02 Unknown table 'shop.b'
> exec
| CREATE TABLE a (x int, y int, KEY x (x), CONSTRAINT F FOREIGN KEY f (y) REFERENCES a (id));
< error
| 5031 HY000 Column not found
> exec
| DROP VIEW IF EXISTS q.c;
< ok 0 0
> exec
| CREATE VIEW a AS SELECT 1 AS x;
< ok 0 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change34 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| DROP INDEX j ON a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE q.v AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0, foreign_key_checks = 1;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS q.b (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE A CHANGE y x int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TRIGGER tr BEFORE INSERT ON q.v FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| ALTER TABLE a DROP FOREIGN KEY g;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE q.c (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON q.a_obfk_1 FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| TRUNCATE TABLE a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| RENAME TABLE A TO v;
< error
| 1017 HY000 File not exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| RENAME TABLE c TO A;
< error
| 1017 HY000 File not exist

== TestWalkThroughSweepOceanBase/change34 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	x	
| database	shop			
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
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
| index	shop	a	I	
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
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
| index	shop	a	I	
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
| index	shop	a	I	
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
| index	shop	a	I	
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
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
| column	shop	b	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	b	I	
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
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
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	b	I	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
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
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	b	I	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
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
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	b	I	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
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
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	b	I	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
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
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	b	y	
| database	shop			
| index	shop	a	I	
| index	shop	a	PRIMARY	
| index	shop	b	I	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE

== TestWalkThroughSweepOceanBase/change34 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| RENAME TABLE c TO q.a_obfk_1;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE a_obfk_1 LIKE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP TEMPORARY TABLE v;
< error
| 1051 42S02 Unknown table 'shop.v'
> exec
| TRUNCATE TABLE A;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| DROP INDEX j ON a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr;
< ok 0 0
> exec
| CREATE TABLE v AS SELECT 1 AS x;
< ok 1 0
> exec
| CREATE VIEW q.A AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP INDEX j ON v;
< error
| 1091 42000 Can't DROP 'j'; check that column/key exists
> exec
| CREATE INDEX i ON q.A (x);
< error
| 1049 42000 Unknown database
> exec
| ALTER TABLE q.c ADD COLUMN z int AFTER y;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE v DROP COLUMN z;
< error
| 1091 42000 Can't DROP 'z'; check that column/key exists
> exec
| CREATE OR REPLACE VIEW q.b AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE IF NOT EXISTS a (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE q.A DROP COLUMN y, ADD COLUMN z int;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| INSERT INTO q.a (x) VALUES (1);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| RENAME TABLE a TO q.a;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE a_obfk_1 ADD INDEX j (z);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE A ADD INDEX I (x);
< ok 0 0
> exec
| DROP TABLE v;
< ok 0 0
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TABLE b (x int, y int, KEY I (x));
< ok 0 0
> exec
| CREATE VIEW a_obfk_1 AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE v RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE UNIQUE INDEX x ON q.v (y);
< error
| 1049 42000 Unknown database
> exec
| DROP TRIGGER IF EXISTS tr2;
< ok 0 0
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change35 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| DROP INDEX j ON a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TABLE c LIKE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| USE shop;
< ok 0 0
> exec
| ALTER TABLE b ADD UNIQUE INDEX I (y);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE A ADD COLUMN x int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TRIGGER TR BEFORE INSERT ON c FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< ok 1 0
> exec
| DROP TABLE IF EXISTS b;
< ok 0 0
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS a (id int PRIMARY KEY, x int, y int);
< ok 0 0

== TestWalkThroughSweepOceanBase/change35 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	a	x	
| database	shop			
| relation	shop		a	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	c	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 11
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	c	x	
| database	shop			
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		c	BASE TABLE
| trigger	shop		TR	

== TestWalkThroughSweepOceanBase/change35 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| DROP INDEX i ON A;
< error
| 1091 42000 Can't DROP 'i'; check that column/key exists
> exec
| DROP VIEW IF EXISTS q.A;
< ok 0 0
> exec
| CREATE TABLE c AS SELECT 1 AS x;
< ok 1 0
> exec
| ALTER TABLE q.c ADD FOREIGN KEY F (x) REFERENCES c (id);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec ^
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN x int;
< error
| 1060 42S21 Duplicate column name 'x'
> exec
| ALTER TABLE a_obfk_1 CHANGE z y int;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS a_obfk_1 (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| DROP VIEW c;
< error
| 1347 HY000 'shop.c' is not VIEW
> exec
| CREATE TRIGGER TR BEFORE INSERT ON A FOR EACH ROW SET NEW.x = 1;
< ok 0 0
> exec
| DROP INDEX I ON a;
< error
| 1091 42000 Can't DROP 'I'; check that column/key exists
> exec
| RENAME TABLE b TO c;
< error
| 1050 42S01 Table 'c' already exists
> exec
| RENAME TABLE a TO q.b;
< error
| 1049 42000 Unknown database 'q'
> exec
| TRUNCATE TABLE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP INDEX I ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE q.a ADD INDEX I (y);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE q.a ADD CONSTRAINT j UNIQUE INDEX I (x);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| INSERT INTO a_obfk_1 (x) VALUES (1);
< error
| 1364 HY000 Field 'id' doesn't have a default value
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change36 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| INSERT INTO q.c (x) VALUES (1);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE b DROP INDEX I;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE a_obfk_1 ADD FOREIGN KEY F (x) REFERENCES c (id);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TABLE q.v (id int PRIMARY KEY, x int UNIQUE, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| RENAME TABLE A TO c;
< error
| 1017 HY000 File not exist
> exec
| CREATE TABLE a_obfk_1 AS SELECT 1 AS x;
< ok 1 0
> exec
| ALTER TABLE v CHANGE X z int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DROP TEMPORARY TABLE a;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| ALTER TABLE q.a_obfk_1 ADD COLUMN y int;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| DROP TABLE IF EXISTS c;
< ok 0 0
> exec
| CREATE TABLE b AS SELECT 1 AS x;
< ok 1 0
> exec
| CREATE TABLE A (x int, y int, KEY x (x));
< ok 0 0
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| ALTER TABLE q.c ADD PRIMARY KEY (x);
< error
| 1146 42S02 Table 'q.c' doesn't exist

== TestWalkThroughSweepOceanBase/change36 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| a_obfk_1
: bytes
| BASE TABLE
: bytes
| InnoDB
< row
: bytes
| b
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| a_obfk_1
: bytes
| x
< row
: bytes
| b
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| a
: bytes
| x
: int64 1
: bytes
| x
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 12
: bytes
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	x	
| column	shop	b	x	
| column	shop	c	x	
| database	shop			
| index	shop	a	x	
| index	shop	b	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE

== TestWalkThroughSweepOceanBase/change36 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| DROP INDEX I ON q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE c ADD CONSTRAINT g FOREIGN KEY a_obfk_1 (x) REFERENCES q.a (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE q.A (x int, y int, KEY i (x), FOREIGN KEY F (y) REFERENCES a (id));
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE IF NOT EXISTS b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE b ADD PRIMARY KEY (x);
< ok 0 0
> exec
| ALTER TABLE A DROP INDEX i, ADD INDEX I (x);
< error
| 1091 42000 Can't DROP 'i'; check that column/key exists
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS q.v (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE a_obfk_1 ADD UNIQUE INDEX I (y);
< error
| 1072 42000 Key column 'y' doesn't exist in table
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| CREATE TABLE IF NOT EXISTS c LIKE a_obfk_1;
< ok 0 0
> exec
| ALTER TABLE q.a RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE b RENAME INDEX i TO j;
< error
| 1176 42000 Key 'i' doesn't exist in table 'b'
> exec
| ALTER TABLE c CHANGE z X int;
< error
| 1054 42S22 Unknown column 'z' in 'c'
> exec
| ALTER TABLE c DROP FOREIGN KEY a_obfk_1;
< error
| 1091 42000 Can't DROP 'a_obfk_1'; check that column/key exists
> exec
| ALTER TABLE q.A ADD FOREIGN KEY g (x) REFERENCES a (id);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE q.v (x int, y int, KEY x (x));
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP INDEX i ON q.a_obfk_1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| DROP TRIGGER IF EXISTS tr;
< ok 0 0
> exec
| DROP VIEW q.A;
< error
| 1051 42S02 Unknown table 'q.a'
> exec
| TRUNCATE TABLE a;
< ok 0 0
> exec
| INSERT INTO b (x) VALUES (1);
< error
| 1062 23000 Duplicate entry '1' for key 'PRIMARY'
> exec
| ALTER TABLE c DROP INDEX x;
< error
| 1091 42000 Can't DROP 'x'; check that column/key exists
> exec
| TRUNCATE TABLE c;
< ok 0 0
> exec
| CREATE VIEW q.a AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change37 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE v DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| CREATE TABLE v LIKE q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE a (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| RENAME TABLE A TO a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP TABLE IF EXISTS a_obfk_1;
< ok 0 0
> exec
| CREATE TABLE A (id int PRIMARY KEY, x int, y int);
< error
| 1050 42S01 Table 'a' already exists
> exec
| ALTER TABLE b CHANGE y y int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TRIGGER tr BEFORE INSERT ON A FOR EACH ROW SET NEW.x = 1;
< ok 0 0
> exec
| CREATE UNIQUE INDEX i ON a (z);
< error
| 1072 42000 Key column 'z' doesn't exist in table
> exec
| DROP VIEW a;
< error
| 1347 HY000 'shop.a' is not VIEW
> exec
| DROP INDEX x ON c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist

== TestWalkThroughSweepOceanBase/change37 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| a
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| a
: bytes
| x
: int64 0
: bytes
| x
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| tr
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 8
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		v	BASE TABLE
| trigger	shop		tr	

== TestWalkThroughSweepOceanBase/change37 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| ALTER TABLE q.A ADD PRIMARY KEY (x);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| ALTER TABLE b DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE v LIKE A;
< ok 0 0
> exec
| TRUNCATE TABLE a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP INDEX j ON A;
< error
| 1091 42000 Can't DROP 'j'; check that column/key exists
> exec
| CREATE OR REPLACE VIEW v AS SELECT 1 AS x;
< error
| 1347 HY000 'shop.v' is not VIEW
> exec
| CREATE TABLE A AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'a' already exists
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON q.A FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE b DROP FOREIGN KEY F;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1364 HY000 Field 'id' doesn't have a default value
> exec
| DROP VIEW q.a;
< error
| 1051 42S02 Unknown table 'q.a'
> exec
| ALTER TABLE q.c DROP FOREIGN KEY g;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DROP TEMPORARY TABLE b;
< error
| 1051 42S02 Unknown table 'shop.b'
> exec
| ALTER TABLE a_obfk_1 ADD INDEX I (z);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP TABLE a_obfk_1, a;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1'
> exec
| ALTER TABLE q.A DROP INDEX I;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| USE shop;
< ok 0 0
> exec
| CREATE TABLE v (x int, y int, KEY x (x));
< error
| 1050 42S01 Table 'v' already exists
> exec
| SET GLOBAL foreign_key_checks = 1, @@session.foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE b ADD COLUMN y int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE INDEX j ON b (X);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TABLE a_obfk_1;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1'
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change38 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| RENAME TABLE q.A TO A;
< error
| 1017 HY000 File not exist
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< ok 1 0
> exec
| CREATE OR REPLACE VIEW v AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE b CHANGE X y int;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE v DROP FOREIGN KEY f;
< error
| 1347 HY000 'shop.v' is not BASE TABLE
> exec
| ALTER TABLE v DROP INDEX I, ADD INDEX i (x);
< error
| 1347 HY000 'shop.v' is not BASE TABLE
> exec
| CREATE TRIGGER TR BEFORE INSERT ON a_obfk_1 FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE UNIQUE INDEX j ON c (z);
< error
| 1146 42S02 Table 'shop.c' doesn't exist

== TestWalkThroughSweepOceanBase/change38 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
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
| MEMORY
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| x
< row
: bytes
| v
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	a	x	
| column	shop	v	x	
| database	shop			
| relation	shop		a	BASE TABLE
| relation	shop		v	VIEW

== TestWalkThroughSweepOceanBase/change38 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| INSERT INTO v (x) VALUES (1);
< error
| 1471 HY000 The target table v of the INSERT is not insertable-into
> exec
| RENAME TABLE c TO c;
< error
| 1017 HY000 File not exist
> exec
| INSERT INTO A (x) VALUES (1);
< ok 1 0
> exec
| CREATE UNIQUE INDEX j ON a (X);
< error
| 1062 23000 Duplicate entry '1' for key 'j'
> exec
| ALTER TABLE a_obfk_1 ADD CONSTRAINT I UNIQUE INDEX I (z);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TABLE q.v AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TRIGGER tr BEFORE INSERT ON v FOR EACH ROW SET NEW.x = 1;
< error
| 1235 0A000 simple dml trigger isn't used on user table not supported
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE INDEX x ON b (x);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| DROP TABLE b, q.a_obfk_1;
< error
| 1051 42S02 Unknown table 'shop.b,q.a_obfk_1'
> exec
| DROP INDEX I ON q.a_obfk_1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| ALTER TABLE v ADD COLUMN y int AFTER x;
< error
| 1347 HY000 'shop.v' is not BASE TABLE
> exec
| INSERT INTO q.a_obfk_1 (x) VALUES (1);
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| CREATE UNIQUE INDEX I ON a (y);
< error
| 1072 42000 Key column 'y' doesn't exist in table
> exec
| CREATE TABLE v (x int, y int, KEY j (x), FOREIGN KEY g (y) REFERENCES a (id));
< error
| 5031 HY000 Column not found
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change39 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE INDEX i ON q.A (y);
< error
| 1049 42000 Unknown database
> exec
| CREATE OR REPLACE VIEW q.a_obfk_1 AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE c ADD COLUMN X int AFTER y;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE c (x int, y int, KEY j (x));
< ok 0 0
> exec
| CREATE INDEX j ON v (y);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE VIEW A AS SELECT 1 AS x;
< ok 0 0
> exec
| DROP TRIGGER tr;
< error
| 1360 HY000 Trigger does not exist
> exec
| DROP TABLE a, A;
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| CREATE TABLE IF NOT EXISTS a LIKE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 1 0
> exec
| ALTER TABLE a_obfk_1 ADD COLUMN y int AFTER y;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE A DROP INDEX x, ADD INDEX j (x);
< error
| 1347 HY000 'shop.a' is not BASE TABLE
> exec
| ALTER TABLE q.c ADD FOREIGN KEY a_obfk_1 (x) REFERENCES v (id);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| ALTER TABLE c ADD CONSTRAINT i UNIQUE INDEX I (X);
< ok 0 0

== TestWalkThroughSweepOceanBase/change39 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| VIEW
: bytes
| MEMORY
< row
: bytes
| c
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| x
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| c
: bytes
| I
: int64 0
: bytes
| x
: int64 -1
< row
: bytes
| c
: bytes
| j
: int64 1
: bytes
| x
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	I	
| index	shop	c	j	
| relation	shop		a	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	I	
| index	shop	c	j	
| relation	shop		a	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	I	
| index	shop	c	j	
| relation	shop		a	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 9
: bytes
| column	shop	a	x	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	I	
| index	shop	c	j	
| relation	shop		a	VIEW
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| index	shop	c	I	
| index	shop	c	j	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| index	shop	c	I	
| index	shop	c	j	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| index	shop	c	I	
| index	shop	c	j	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| index	shop	c	I	
| index	shop	c	j	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| index	shop	c	I	
| index	shop	c	j	
| relation	shop		a	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	b	x	
| index	shop	c	I	
| index	shop	c	j	
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	c	I	
| index	shop	c	j	
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 16
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	c	I	
| index	shop	c	j	
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 15
: bytes
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	I	
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	I	
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 17
: bytes
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	I	
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	I	
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	I	
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	I	
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	I	
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 19
: bytes
| column	q	a_obfk_1	x	
| column	q	a_obfk_1	y	
| column	shop	a	x	
| column	shop	a_obfk_1	x	
| column	shop	b	id	
| column	shop	b	x	
| column	shop	b	y	
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	q	a_obfk_1	I	
| index	shop	b	PRIMARY	
| index	shop	c	I	
| relation	q		a_obfk_1	BASE TABLE
| relation	shop		a	VIEW
| relation	shop		a_obfk_1	VIEW
| relation	shop		b	BASE TABLE
| relation	shop		c	BASE TABLE

== TestWalkThroughSweepOceanBase/change39 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| DROP TRIGGER TR;
< error
| 1360 HY000 Trigger does not exist
> exec
| DROP TABLE b;
< error
| 1051 42S02 Unknown table 'shop.b'
> exec
| CREATE TRIGGER tr BEFORE INSERT ON b FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE b (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr;
< ok 0 0
> exec
| CREATE UNIQUE INDEX i ON a_obfk_1 (x);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE A ADD CONSTRAINT F FOREIGN KEY (x) REFERENCES a_obfk_1 (id);
< error
| 1146 42S02 Table doesn't exist
> exec
| DROP VIEW b;
< error
| 1347 HY000 'shop.b' is not VIEW
> exec
| CREATE VIEW a_obfk_1 AS SELECT 1 AS x;
< ok 0 0
> exec
| ALTER TABLE b DROP INDEX x;
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 ENGINE = InnoDB;
< error
| 1347 HY000 'shop.a_obfk_1' is not BASE TABLE
> exec
| DROP INDEX j ON c;
< ok 0 0
> exec
| RENAME TABLE v TO a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| ALTER TABLE q.c CHANGE y X int;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DROP TEMPORARY TABLE a_obfk_1, q.a_obfk_1;
< error
| 1051 42S02 Unknown table 'shop.a_obfk_1,q.a_obfk_1'
> exec
| CREATE VIEW A AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP INDEX I ON a_obfk_1;
< error
| 1091 42000 Can't DROP 'I'; check that column/key exists
> exec
| DROP TRIGGER IF EXISTS tr2;
< ok 0 0
> exec
| ALTER TABLE b ADD COLUMN z int AFTER z;
< error
| 1054 42S22 Unknown column 'z' in 'b'
> exec
| CREATE TABLE q.a_obfk_1 LIKE c;
< ok 0 0
> exec
| ALTER TABLE b DROP FOREIGN KEY g;
< error
| 1091 42000 Can't DROP 'g'; check that column/key exists
> exec
| DROP VIEW A;
< ok 0 0
> exec
| CREATE OR REPLACE VIEW A AS SELECT 1 AS x;
< ok 0 0
> exec
| DROP VIEW IF EXISTS q.a;
< ok 0 0
> exec
| DROP TABLE q.v;
< error
| 1051 42S02 Unknown table 'q.v'
> exec
| ALTER TABLE a_obfk_1 ADD CONSTRAINT i UNIQUE INDEX x (z);
< error
| 1347 HY000 'shop.a_obfk_1' is not BASE TABLE
> exec
| CREATE TABLE c AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'c' already exists
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change4 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE c (x int, y int, KEY I (x));
< ok 0 0
> exec
| DROP INDEX x ON q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS a, A;
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| ALTER TABLE v ADD COLUMN X int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TRIGGER tr BEFORE INSERT ON q.A FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| CREATE OR REPLACE VIEW c AS SELECT 1 AS x;
< error
| 1347 HY000 'shop.c' is not VIEW
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS a_obfk_1 (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE b DROP INDEX x;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE v DROP INDEX j;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| SET foreign_key_checks = 1;
< ok 0 0
> exec
| ALTER TABLE c ADD INDEX i (z);
< error
| 1061 42000 Duplicate key name 'i'
> exec
| RENAME TABLE q.v TO q.c;
< error
| 1049 42000 Unknown database 'q'
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE A (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| DROP INDEX I ON a_obfk_1;
< error
| 1091 42000 Can't DROP 'I'; check that column/key exists

== TestWalkThroughSweepOceanBase/change4 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| a
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| a
: bytes
| x
: int64 0
: bytes
| x
: int64 -1
< row
: bytes
| c
: bytes
| I
: int64 1
: bytes
| x
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	I	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	id	
| column	shop	v	x	
| column	shop	v	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	j	
| index	shop	a	x	
| index	shop	c	I	
| index	shop	v	PRIMARY	
| index	shop	v	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
| relation	shop		v	BASE TABLE

== TestWalkThroughSweepOceanBase/change4 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| CREATE TABLE IF NOT EXISTS v (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| ALTER TABLE q.c CHANGE y X int;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| INSERT INTO A (x) VALUES (1);
< error
| 1364 HY000 Field 'id' doesn't have a default value
> exec
| ALTER TABLE A DROP FOREIGN KEY f;
< error
| 1091 42000 Can't DROP 'f'; check that column/key exists
> exec
| CREATE TABLE A LIKE c;
< error
| 1050 42S01 Table 'a' already exists
> exec
| ALTER TABLE q.a CHANGE z z int;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE q.c LIKE q.a;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| CREATE TABLE c (x int, y int, KEY I (x), CONSTRAINT F FOREIGN KEY (y) REFERENCES a (id));
< error
| 1050 42S01 Table 'c' already exists
> exec
| CREATE INDEX j ON A (x);
< ok 0 0
> exec
| TRUNCATE TABLE v;
< ok 0 0
> exec
| CREATE TABLE q.c LIKE c;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE IF NOT EXISTS A LIKE c;
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS TR;
< ok 0 0
> exec
| CREATE INDEX j ON q.c (z);
< error
| 1049 42000 Unknown database
> exec
| TRUNCATE TABLE q.a_obfk_1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| ALTER TABLE q.v CHANGE y y int;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| ALTER TABLE q.a_obfk_1 DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| DROP TEMPORARY TABLE a, c;
< error
| 1051 42S02 Unknown table 'shop.a,shop.c'
> exec
| CREATE TABLE c (x int, y int, KEY i (x), CONSTRAINT g FOREIGN KEY (y) REFERENCES a (id));
< error
| 1050 42S01 Table 'c' already exists
> exec
| CREATE OR REPLACE VIEW v AS SELECT 1 AS x;
< error
| 1347 HY000 'shop.v' is not VIEW
> exec
| DROP TRIGGER TR;
< error
| 1360 HY000 Trigger does not exist
> exec
| SET GLOBAL foreign_key_checks = 0;
< ok 0 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change5 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE c ADD COLUMN y int;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP VIEW IF EXISTS c;
< ok 0 0
> exec
| CREATE INDEX i ON c (X);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP TABLE A;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< ok 1 0
> exec
| CREATE TABLE IF NOT EXISTS b LIKE q.c;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DROP TABLE IF EXISTS q.A, a;
< ok 0 0
> exec
| CREATE INDEX x ON q.A (y);
< error
| 1049 42000 Unknown database
> exec
| CREATE TABLE A LIKE c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| CREATE TABLE c AS SELECT 1 AS x;
< ok 1 0
> exec
| CREATE TABLE q.c (x int, y int, KEY j (x));
< error
| 1049 42000 Unknown database 'q'

== TestWalkThroughSweepOceanBase/change5 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| c
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| c
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE
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
| column	shop	c	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
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
| index	shop	a	PRIMARY	
| index	shop	a	x	
| index	shop	c	j	
| relation	shop		a	BASE TABLE
| relation	shop		c	BASE TABLE

== TestWalkThroughSweepOceanBase/change5 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| RENAME TABLE a TO b, c TO a;
< error
| 1017 HY000 File not exist
> exec
| CREATE TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| RENAME TABLE c TO q.a;
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP DATABASE IF EXISTS q;
< ok 0 0
> exec
| CREATE INDEX j ON c (X);
< ok 0 0
> exec
| DROP VIEW A;
< error
| 1347 HY000 'shop.a' is not VIEW
> exec
| ALTER TABLE c ADD COLUMN y int;
< ok 0 0
> exec
| DROP INDEX x ON c;
< error
| 1091 42000 Can't DROP 'x'; check that column/key exists
> exec
| CREATE TABLE c AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'c' already exists
> exec
| DROP TABLE v, q.b;
< error
| 1051 42S02 Unknown table 'shop.v,q.b'
> exec
| ALTER TABLE v ADD FOREIGN KEY g (x) REFERENCES v (id);
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| TRUNCATE TABLE q.a_obfk_1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| ALTER TABLE q.c DROP FOREIGN KEY a_obfk_1;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| CREATE UNIQUE INDEX j ON a_obfk_1 (X);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE A ADD PRIMARY KEY (x);
< error
| 1068 42000 Multiple primary key defined
> exec
| CREATE TABLE a AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'a' already exists
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS q.a_obfk_1 (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| RENAME TABLE c TO a_obfk_1, b TO q.c;
< error
| 1049 42000 Unknown database 'q'
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change6 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE A ADD CONSTRAINT F FOREIGN KEY (x) REFERENCES a_obfk_1 (id);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE v DROP COLUMN x;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| TRUNCATE TABLE q.A;
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DROP VIEW c;
< error
| 1051 42S02 Unknown table 'shop.c'
> exec
| CREATE TABLE c (x int, y int, KEY i (x));
< ok 0 0
> exec
| ALTER TABLE a_obfk_1 DROP INDEX j, ADD INDEX i (x);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE b DROP INDEX j, ADD INDEX j (x);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE A DROP COLUMN x, ADD COLUMN y int;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TRIGGER TR BEFORE INSERT ON A FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.A' doesn't exist
> exec
| ALTER TABLE A ADD FOREIGN KEY f (x) REFERENCES v (id);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE b DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'shop.b' doesn't exist

== TestWalkThroughSweepOceanBase/change6 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| c
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| c
: bytes
| i
: int64 1
: bytes
| x
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
| temporary	shop		b	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 7
: bytes
| column	shop	c	x	
| column	shop	c	y	
| database	q			
| database	shop			
| index	shop	c	i	
| relation	shop		c	BASE TABLE
| temporary	shop		b	

== TestWalkThroughSweepOceanBase/change6 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| USE shop;
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS q.a_obfk_1 (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP TRIGGER tr2;
< error
| 1360 HY000 Trigger does not exist
> exec
| CREATE TABLE q.c AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON b FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE VIEW q.a_obfk_1 AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE q.v CHANGE z X int;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| ALTER TABLE q.b DROP FOREIGN KEY g;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| CREATE TEMPORARY TABLE b (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE DATABASE IF NOT EXISTS q;
< ok 1 0
> exec
| CREATE TRIGGER TR BEFORE INSERT ON q.b FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| CREATE TABLE v LIKE v;
< error
| 1066 42000 Not unique table/alias: 'v'
> exec
| ALTER TABLE a_obfk_1 ADD CONSTRAINT x UNIQUE INDEX x (z);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE OR REPLACE VIEW c AS SELECT 1 AS x;
< error
| 1347 HY000 'shop.c' is not VIEW
> exec
| ALTER TABLE c DROP INDEX j, ADD INDEX x (x);
< error
| 1091 42000 Can't DROP 'j'; check that column/key exists
> exec
| DROP INDEX j ON A;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE INDEX i ON a_obfk_1 (x);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| ALTER TABLE a_obfk_1 ADD CONSTRAINT g FOREIGN KEY (x) REFERENCES A (id);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE INDEX i ON q.a (x);
< error
| 1146 42S02 Table 'q.a' doesn't exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change7 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| INSERT INTO a_obfk_1 (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TABLE v LIKE v;
< error
| 1066 42000 Not unique table/alias: 'v'
> exec
| RENAME TABLE v TO a_obfk_1;
< error
| 1017 HY000 File not exist
> exec
| ALTER TABLE b ENGINE = InnoDB;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS q.a (id int PRIMARY KEY, x int, y int);
< error
| 1049 42000 Unknown database 'q'
> exec
| RENAME TABLE v TO a_obfk_1, c TO c;
< error
| 1017 HY000 File not exist
> exec
| CREATE INDEX i ON A (y);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE v CHANGE x y int;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| ALTER TABLE a_obfk_1 RENAME INDEX i TO j;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist

== TestWalkThroughSweepOceanBase/change7 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 1
: bytes
| database	shop			
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 3
: bytes
| column	shop	c	x	
| database	shop			
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
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 4
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
| temporary	shop		v	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 5
: bytes
| column	shop	c	x	
| database	shop			
| relation	shop		c	VIEW
| temporary	shop		a	
| temporary	shop		v	

== TestWalkThroughSweepOceanBase/change7 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| CREATE VIEW c AS SELECT 1 AS x;
< ok 0 0
> exec
| CREATE TEMPORARY TABLE IF NOT EXISTS A (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE A RENAME INDEX i TO j;
< error
| 1176 42000 Key 'i' doesn't exist in table 'a'
> exec
| ALTER TABLE b ADD CONSTRAINT I UNIQUE INDEX j (X);
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| ALTER TABLE a CHANGE x X int;
< ok 0 0
> exec
| USE q;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TABLE c (x int, y int, KEY I (x));
< error
| 1050 42S01 Table 'c' already exists
> exec
| CREATE TABLE a LIKE a;
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| ALTER TABLE q.v DROP FOREIGN KEY f;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| ALTER TABLE c ADD CONSTRAINT f FOREIGN KEY (x) REFERENCES a (id);
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON b FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| TRUNCATE TABLE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE q.A LIKE q.a_obfk_1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| CREATE UNIQUE INDEX j ON A (y);
< ok 0 0
> exec
| CREATE VIEW c AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'c' already exists
> exec
| ALTER TABLE A DROP COLUMN y;
< ok 0 0
> exec
| ALTER TABLE v DROP FOREIGN KEY a_obfk_1;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| INSERT INTO a (x) VALUES (1);
< error
| 1364 HY000 Field 'id' doesn't have a default value
> exec
| ALTER TABLE a_obfk_1 CHANGE X y int;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| CREATE TEMPORARY TABLE v (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE c CHANGE x z int;
< error
| 1347 HY000 'shop.c' is not BASE TABLE
> exec
| RENAME TABLE a_obfk_1 TO q.a_obfk_1;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON q.A FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.A' doesn't exist
> exec
| ALTER TABLE a DROP COLUMN x, ADD COLUMN y int;
< ok 0 0
> exec
| DROP INDEX I ON q.c;
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change8 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE A ADD INDEX i (X);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE q.c ADD UNIQUE INDEX I (z);
< error
| 1146 42S02 Table 'q.c' doesn't exist
> exec
| DROP INDEX I ON v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE OR REPLACE VIEW q.a_obfk_1 AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE UNIQUE INDEX j ON a_obfk_1 (X);
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP TEMPORARY TABLE IF EXISTS q.a_obfk_1;
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE INDEX i (z);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TABLE a (id int PRIMARY KEY, x int, y int);
< ok 0 0

== TestWalkThroughSweepOceanBase/change8 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
: bytes
| InnoDB
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| a
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 6
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| database	shop			
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
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| temporary	shop		a	
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
| index	shop	a	PRIMARY	
| relation	shop		a	BASE TABLE
| temporary	shop		a	
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
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
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
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
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
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
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
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
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
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
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
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
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
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
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
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		a	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 14
: bytes
| column	shop	a	id	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	a_obfk_1	id	
| column	shop	a_obfk_1	x	
| column	shop	a_obfk_1	y	
| column	shop	b	x	
| database	shop			
| index	shop	a	PRIMARY	
| index	shop	a_obfk_1	PRIMARY	
| relation	shop		a	BASE TABLE
| relation	shop		a_obfk_1	BASE TABLE
| relation	shop		b	VIEW
| temporary	shop		a	

== TestWalkThroughSweepOceanBase/change8 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| CREATE TEMPORARY TABLE A (id int PRIMARY KEY, x int UNIQUE, y int);
< ok 0 0
> exec
| TRUNCATE TABLE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| CREATE TABLE IF NOT EXISTS a_obfk_1 (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| ALTER TABLE c ADD CONSTRAINT F FOREIGN KEY f (x) REFERENCES c (id);
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| ALTER TABLE b ENGINE = InnoDB;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TABLE a_obfk_1 AS SELECT 1 AS x;
< error
| 1050 42S01 Table 'a_obfk_1' already exists
> exec
| ALTER TABLE a DROP PRIMARY KEY;
< ok 0 0
> exec
| TRUNCATE TABLE c;
< error
| 1146 42S02 Table 'shop.c' doesn't exist
> exec
| DROP VIEW a_obfk_1;
< error
| 1347 HY000 'shop.a_obfk_1' is not VIEW
> exec
| CREATE TRIGGER TR BEFORE INSERT ON q.b FOR EACH ROW SET NEW.x = 1;
< error
| 1146 42S02 Table 'q.b' doesn't exist
> exec
| CREATE VIEW b AS SELECT 1 AS x;
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr;
< ok 0 0
> exec
| CREATE TABLE A (x int, y int, KEY x (x), CONSTRAINT g FOREIGN KEY (y) REFERENCES a (id));
< error
| 5031 HY000 Column not found
> exec
| ALTER TABLE a_obfk_1 RENAME INDEX i TO j;
< error
| 1176 42000 Key 'i' doesn't exist in table 'a_obfk_1'
> exec
| DROP TABLE q.c;
< error
| 1051 42S02 Unknown table 'q.c'
> exec
| CREATE TABLE a_obfk_1 LIKE v;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0

== TestWalkThroughSweepOceanBase/change9 default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE c (id int PRIMARY KEY, x int, y int);
< ok 0 0
> exec
| CREATE TABLE a LIKE a;
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP TABLE A;
< error
| 1051 42S02 Unknown table 'shop.a'
> exec
| ALTER TABLE A DROP INDEX i, ADD INDEX x (x);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON c FOR EACH ROW SET NEW.x = 1;
< ok 0 0
> exec
| RENAME TABLE a TO a, q.c TO v;
< error
| 1017 HY000 File not exist
> exec
| CREATE TABLE a LIKE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| CREATE TRIGGER tr BEFORE INSERT ON c FOR EACH ROW SET NEW.x = 1;
< ok 0 0
> exec
| CREATE TABLE q.a_obfk_1 AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| TRUNCATE TABLE q.a_obfk_1;
< error
| 1146 42S02 Table 'q.a_obfk_1' doesn't exist
> exec
| ALTER TABLE c ADD COLUMN y int;
< error
| 1060 42S21 Duplicate column name 'y'
> exec
| USE shop;
< ok 0 0
> exec
| CREATE VIEW v AS SELECT 1 AS x;
< ok 0 0

== TestWalkThroughSweepOceanBase/change9 default#2
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "VARCHAR"
< column "COALESCE(ENGINE, '')" "VARCHAR"
< row
: bytes
| c
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
| MEMORY
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
< row
: bytes
| v
: bytes
| x
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COALESCE(COLUMN_NAME, '')" "VARCHAR"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| c
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| c
: bytes
| tr
< row
: bytes
| c
: bytes
| tr2
> query
| SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME
< result
< column "ROUTINE_TYPE" "VARCHAR"
< column "ROUTINE_NAME" "VARCHAR"
> query
| SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME
< result
< column "EVENT_NAME" "VARCHAR"
> query
| SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME
< result
< column "MVIEW_NAME" "VARCHAR"
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
| UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
| UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
| UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('omni_walk')
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
| UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
| UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
| ) o WHERE o.db NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
| trigger	shop		tr2	
> query ^
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "VARCHAR"
< row
: int64 10
: bytes
| column	shop	c	id	
| column	shop	c	x	
| column	shop	c	y	
| column	shop	v	x	
| database	shop			
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| relation	shop		v	VIEW
| trigger	shop		tr	
| trigger	shop		tr2	
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
| trigger	shop		tr2	
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
| trigger	shop		tr2	
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
| trigger	shop		tr2	
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
| trigger	shop		tr2	
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
| trigger	shop		tr2	
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
| trigger	shop		tr2	
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
| trigger	shop		tr2	
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
| trigger	shop		tr2	
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
| trigger	shop		tr2	
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
| trigger	shop		tr2	
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
| trigger	shop		tr2	
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr	
| trigger	shop		tr2	
> query ^
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr2	
> query ^
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr2	
> query ^
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr2	
> query ^
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr2	
> query ^
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr2	
> query ^
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
| index	shop	c	PRIMARY	
| relation	shop		c	BASE TABLE
| trigger	shop		tr2	

== TestWalkThroughSweepOceanBase/change9 default#3
> exec
| USE shop
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SELECT GET_LOCK('omni_walk', 60)
< result
< column "GET_LOCK('omni_walk', 60)" "BIGINT"
< row
: int64 1
> exec
| DROP INDEX x ON A;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE a_obfk_1 RENAME TO q.a_obfk_1;
< error
| 1146 42S02 Table 'shop.a_obfk_1' doesn't exist
> exec
| DROP VIEW IF EXISTS v;
< ok 0 0
> exec
| DROP VIEW IF EXISTS q.a;
< ok 0 0
> exec
| CREATE TABLE q.A AS SELECT 1 AS x;
< error
| 1049 42000 Unknown database 'q'
> exec
| TRUNCATE TABLE a;
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| ALTER TABLE q.v DROP COLUMN X;
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| DROP TABLE A, q.v;
< error
| 1051 42S02 Unknown table 'shop.a,q.v'
> exec
| CREATE INDEX x ON A (z);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| INSERT INTO q.v (x) VALUES (1);
< error
| 1146 42S02 Table 'q.v' doesn't exist
> exec
| RENAME TABLE c TO q.c;
< error
| 1049 42000 Unknown database 'q'
> exec
| ALTER TABLE v DROP INDEX i;
< error
| 1146 42S02 Table 'shop.v' doesn't exist
> exec
| RENAME TABLE b TO v;
< error
| 1017 HY000 File not exist
> exec
| CREATE TABLE q.b (x int, y int, KEY x (x), CONSTRAINT f FOREIGN KEY a_obfk_1 (y) REFERENCES a (id));
< error
| 1049 42000 Unknown database 'q'
> exec
| DROP TRIGGER tr;
< ok 0 0
> exec
| TRUNCATE TABLE b;
< error
| 1146 42S02 Table 'shop.b' doesn't exist
> exec
| SET GLOBAL foreign_key_checks = 0, foreign_key_checks = 1;
< ok 0 0
> exec
| DROP VIEW IF EXISTS b;
< ok 0 0
> exec
| INSERT INTO a (x) VALUES (1);
< error
| 1146 42S02 Table 'shop.a' doesn't exist
> exec
| RENAME TABLE a_obfk_1 TO q.c;
< error
| 1049 42000 Unknown database 'q'
> exec
| CREATE DATABASE q;
< ok 1 0
> exec
| DO RELEASE_LOCK('omni_walk')
< ok 1 0
