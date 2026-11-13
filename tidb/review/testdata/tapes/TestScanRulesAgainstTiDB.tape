-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScanRulesAgainstTiDB/CREATE_TABLE_IF_NOT_EXISTS,_not_judged default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| CREATE TABLE IF NOT EXISTS a (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| a
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| a
: int64 0

== TestScanRulesAgainstTiDB/CREATE_TABLE_IF_NOT_EXISTS_of_a_table_with_a_key default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| CREATE TABLE IF NOT EXISTS plain (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0

== TestScanRulesAgainstTiDB/DROP_CHECK default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| ALTER TABLE t DROP CHECK ck_n
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0

== TestScanRulesAgainstTiDB/DROP_INDEX_of_the_primary_key default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| DROP INDEX `PRIMARY` ON plain
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| plain
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| plain
: int64 0

== TestScanRulesAgainstTiDB/a_key_added_later default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| CREATE TABLE a (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| a
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| a
: int64 0
> exec
| ALTER TABLE a ADD PRIMARY KEY (id)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| a
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| a
: int64 1
< row
: bytes
| nokey
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| a
: int64 1
< row
: bytes
| nokey
: int64 0

== TestScanRulesAgainstTiDB/a_primary_key_dropped,_then_added_back default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| ALTER TABLE plain DROP PRIMARY KEY
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| plain
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| plain
: int64 0
> exec
| ALTER TABLE plain ADD PRIMARY KEY (v)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| v
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0

== TestScanRulesAgainstTiDB/a_table_created_without_a_key default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| CREATE TABLE a (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| a
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| a
: int64 0
> exec
| CREATE TABLE b (id INT PRIMARY KEY)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| b
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| b
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| a
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| b
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| a
: int64 0

== TestScanRulesAgainstTiDB/a_table_renamed_after_losing_its_key default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| ALTER TABLE plain DROP PRIMARY KEY
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| plain
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| plain
: int64 0
> exec
| RENAME TABLE plain TO plain2
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| plain2
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| plain2
: int64 0

== TestScanRulesAgainstTiDB/a_table_whose_database_the_change_drops default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| CREATE DATABASE d
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| CREATE TABLE d.e (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| DROP DATABASE d
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0

== TestScanRulesAgainstTiDB/a_unique_index_dropped_by_name default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| ALTER TABLE t DROP INDEX uk_code
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| DROP INDEX idx_n ON t
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0

== TestScanRulesAgainstTiDB/a_unique_key_the_change_makes,_then_drops default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| ALTER TABLE nokey ADD UNIQUE KEY uk_a (a)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| nokey
: bytes
| UNIQUE
: bytes
| uk_a
: bytes
| a
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| nokey
: bytes
| UNIQUE
: bytes
| uk_a
: bytes
| a
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| ALTER TABLE nokey DROP INDEX uk_a
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0

== TestScanRulesAgainstTiDB/an_index_of_a_table_CREATE_TABLE_IF_NOT_EXISTS_leaves default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| CREATE TABLE IF NOT EXISTS t (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| DROP INDEX uk_code ON t
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0

== TestScanRulesAgainstTiDB/an_index_of_a_table_the_change_recreates default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| DROP TABLE t
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| CREATE TABLE t (code INT, KEY uk_code (code))
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| t
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| t
: int64 0
> exec
| ALTER TABLE t DROP INDEX uk_code
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| t
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| t
: int64 0

== TestScanRulesAgainstTiDB/constraint_drops_the_SQL_names default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| t
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| nokey
: int64 0
> exec
| ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY fk_p
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
< row
: bytes
| nokey
: int64 0
> exec
| ALTER TABLE t DROP CONSTRAINT ck_n
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
< row
: bytes
| nokey
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
< row
: bytes
| nokey
: int64 0
> exec
| ALTER TABLE t DROP INDEX uk_code
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			'')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		UNION ALL
| 		SELECT TABLE_NAME, 'CHECK', CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t'')" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
< row
: bytes
| nokey
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "BIGINT"
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
< row
: bytes
| nokey
: int64 0
