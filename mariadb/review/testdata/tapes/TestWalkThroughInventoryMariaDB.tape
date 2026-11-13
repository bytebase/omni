-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestWalkThroughInventoryMariaDB default#1
> exec
| SET SESSION group_concat_max_len = 16777216
< ok 0 0
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE
> query
| SELECT COUNT(*), GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\n') FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(SCHEMA_NAME ORDER BY BINARY SCHEMA_NAME SEPARATOR '\\n')" "TEXT"
< row
: int64 1
: bytes
| shop
> query
| SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
| SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) AS line FROM information_schema.TABLES WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, '') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN (?) AND SEQ_IN_INDEX = 1
| UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, '') FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, '') FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, '') FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN (?)
| UNION ALL SELECT CONCAT_WS('\t', 'event', EVENT_SCHEMA, '', EVENT_NAME, '') FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN (?)
| ) o
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
: string
| shop
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\\n')" "TEXT"
< row
: int64 12
: bytes
| column	shop	A	x	
| column	shop	A	y	
| column	shop	B	x	
| column	shop	a	x	
| column	shop	a	y	
| column	shop	b	x	
| index	shop	B	x	
| index	shop	b	x	
| relation	shop		A	SYSTEM VERSIONED
| relation	shop		B	BASE TABLE
| relation	shop		a	BASE TABLE
| relation	shop		b	BASE TABLE

== TestWalkThroughInventoryMariaDB default#2
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE a (x int, y int)
< ok 0 0
> exec
| CREATE TABLE A (x int, y int) WITH SYSTEM VERSIONING
< ok 0 0
> exec
| CREATE TABLE b (x int, KEY x (x))
< ok 0 0
> exec
| CREATE TABLE B (x int, KEY x (x))
< ok 0 0
> exec
| CREATE TEMPORARY TABLE c (x int)
< ok 0 0
> exec
| CREATE TEMPORARY TABLE C (x int)
< ok 0 0
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| DROP SEQUENCE b, b
< error
| 1066 42000 Not unique table/alias: 'b'
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| TRUNCATE TABLE nosuch
< error
| 1146 42S02 Table 'shop.nosuch' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
> exec
| RENAME TABLE s TO a, s TO b
< error
| 1146 42S02 Table 'shop.s' doesn't exist
> query
| SELECT COUNT(*), GROUP_CONCAT(CONCAT_WS('\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\n') FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%'
< result
< column "COUNT(*)" "BIGINT"
< column "GROUP_CONCAT(CONCAT_WS('\\t', 'temporary', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE) ORDER BY BINARY TABLE_SCHEMA, BINARY TABLE_NAME, TABLE_TYPE SEPARATOR '\\n')" "TEXT"
< row
: int64 2
: bytes
| temporary	shop		C	TEMPORARY
| temporary	shop		c	TEMPORARY
