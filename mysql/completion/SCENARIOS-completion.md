# MySQL SQL Autocompletion Scenarios

> Goal: Implement parser-native SQL autocompletion for MySQL, matching PG's architecture — parser instrumentation with collectMode, completion module with Complete(sql, cursor, catalog) API
> Verification: Table-driven Go tests — Complete(sql, cursorPos, catalog) returns expected candidates with correct types
> Reference: PG completion system (pg/completion/, pg/parser/complete.go)
> Out of scope (v1): XA transactions, SHUTDOWN, RESTART, CLONE, INSTALL/UNINSTALL PLUGIN, IMPORT TABLE, BINLOG, CACHE INDEX, PURGE BINARY LOGS, CHANGE REPLICATION, START/STOP REPLICA

Status legend: `[ ]` pending, `[x]` passing, `[~]` partial

---

## Phase 1: Parser Completion Infrastructure

### 1.1 Completion Mode Fields & Entry Point

```
[x] Parser struct has completion fields: completing bool, cursorOff int, candidates *CandidateSet, collecting bool, maxCollect int
[x] collectMode() returns true when completing && collecting
[x] checkCursor() sets collecting=true when p.cur.Loc >= cursorOff
[x] CandidateSet struct with Tokens []int, Rules []RuleCandidate, seen/seenR dedup maps
[x] RuleCandidate struct with Rule string (e.g., "columnref", "table_ref")
[x] addTokenCandidate(tokType int) adds to CandidateSet.Tokens with dedup
[x] addRuleCandidate(rule string) adds to CandidateSet.Rules with dedup
[x] Collect(sql string, cursorOffset int) *CandidateSet — entry point with panic recovery
[x] Collect returns non-nil CandidateSet even for empty input
[x] Collect returns keyword candidates at statement start position: SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, DROP, etc.
[x] All existing parser tests pass with completion fields present (completing=false) — regression gate
```

### 1.2 Basic Keyword Collection

```
[x] Empty input → keyword candidates for all top-level statements
[x] After semicolon → keyword candidates for new statement
[x] `SELECT |` → keyword candidates for select expressions (DISTINCT, ALL) + rule candidates (columnref, func_name)
[x] `CREATE |` → keyword candidates for object types (TABLE, INDEX, VIEW, DATABASE, FUNCTION, PROCEDURE, TRIGGER, EVENT)
[x] `ALTER |` → keyword candidates (TABLE, DATABASE, VIEW, FUNCTION, PROCEDURE, EVENT)
[x] `DROP |` → keyword candidates (TABLE, INDEX, VIEW, DATABASE, FUNCTION, PROCEDURE, TRIGGER, EVENT, IF)
```

## Phase 2: Completion Module (API & Resolution)

Build the public API and candidate resolution before instrumenting grammar rules,
so that Phase 3+ instrumentation can be tested end-to-end.

### 2.1 Public API & Core Logic

```
[x] Complete(sql, cursorOffset, catalog) returns []Candidate
[x] Candidate struct has Text, Type, Definition, Comment fields
[x] CandidateType enum: Keyword, Database, Table, View, Column, Function, Procedure, Index, Trigger, Event, Variable, Charset, Engine, Type
[x] Complete with nil catalog returns keyword-only candidates
[x] Complete with empty sql returns top-level statement keywords
[x] Prefix filtering: `SEL|` matches SELECT keyword
[x] Prefix filtering is case-insensitive
[x] Deduplication: same candidate not returned twice
```

### 2.2 Candidate Resolution

```
[x] Token candidates → keyword strings (from token type mapping)
[x] "table_ref" rule → catalog tables + views
[x] "columnref" rule → columns from tables in scope
[x] "database_ref" rule → catalog databases
[x] "function_ref" / "func_name" rule → catalog functions + built-in function names
[x] "procedure_ref" rule → catalog procedures
[x] "index_ref" rule → indexes from relevant table
[x] "trigger_ref" rule → catalog triggers
[x] "event_ref" rule → catalog events
[x] "view_ref" rule → catalog views
[x] "charset" rule → known charset names (utf8mb4, latin1, utf8, ascii, binary)
[x] "engine" rule → known engine names (InnoDB, MyISAM, MEMORY, CSV, ARCHIVE)
[x] "type_name" rule → MySQL type keywords (INT, VARCHAR, TEXT, BLOB, DATE, etc.)
```

### 2.3 Table Reference Extraction

```
[x] Extract table refs from simple SELECT: `SELECT * FROM t` → [{Table: "t"}]
[x] Extract table refs with alias: `SELECT * FROM t AS x` → [{Table: "t", Alias: "x"}]
[x] Extract table refs from JOIN: `FROM t1 JOIN t2 ON ...` → [{Table: "t1"}, {Table: "t2"}]
[x] Extract table refs with database: `FROM db.t` → [{Database: "db", Table: "t"}]
[x] Extract table refs from subquery: inner tables don't leak to outer scope
[x] Extract table refs from UPDATE: `UPDATE t SET ...` → [{Table: "t"}]
[x] Extract table refs from INSERT: `INSERT INTO t ...` → [{Table: "t"}]
[x] Extract table refs from DELETE: `DELETE FROM t ...` → [{Table: "t"}]
[x] Fallback to lexer-based extraction when AST parsing fails (incomplete SQL)
```

### 2.4 Tricky Completion (Fallback)

```
[ ] Incomplete SQL: `SELECT * FROM ` (trailing space) → insert placeholder, re-collect
[ ] Truncated mid-keyword: `SELE` → prefix-filter against keywords
[ ] Truncated after comma: `SELECT a,` → insert placeholder column
[ ] Truncated after operator: `WHERE a >` → insert placeholder expression
[ ] Multiple placeholder strategies tried in order
[ ] Fallback returns best-effort results when no strategy succeeds
[ ] Placeholder insertion does not corrupt the candidate set from the initial pass
```

## Phase 3: SELECT Statement Instrumentation

### 3.1 SELECT Target List

```
[ ] `SELECT |` → columnref, func_name, literal keywords (DISTINCT, ALL, *)
[ ] `SELECT a, |` → columnref, func_name after comma
[ ] `SELECT a, b, |` → columnref after multiple commas
[ ] `SELECT * FROM t WHERE a > (SELECT |)` → columnref in subquery
[ ] `SELECT DISTINCT |` → columnref after DISTINCT
```

### 3.2 FROM Clause

```
[ ] `SELECT * FROM |` → table_ref (tables, views, databases)
[ ] `SELECT * FROM db.|` → table_ref qualified with database
[ ] `SELECT * FROM t1, |` → table_ref after comma (multi-table)
[ ] `SELECT * FROM (SELECT * FROM |)` → table_ref in derived table
[ ] `SELECT * FROM t |` → keyword candidates (WHERE, JOIN, LEFT, RIGHT, CROSS, NATURAL, STRAIGHT_JOIN, ORDER, GROUP, HAVING, LIMIT, UNION, FOR)
[ ] `SELECT * FROM t AS |` → no specific candidates (alias context)
```

### 3.3 JOIN Clauses

```
[ ] `SELECT * FROM t1 JOIN |` → table_ref after JOIN
[ ] `SELECT * FROM t1 LEFT JOIN |` → table_ref after LEFT JOIN
[ ] `SELECT * FROM t1 RIGHT JOIN |` → table_ref after RIGHT JOIN
[ ] `SELECT * FROM t1 CROSS JOIN |` → table_ref after CROSS JOIN
[ ] `SELECT * FROM t1 NATURAL JOIN |` → table_ref after NATURAL JOIN
[ ] `SELECT * FROM t1 STRAIGHT_JOIN |` → table_ref after STRAIGHT_JOIN
[ ] `SELECT * FROM t1 JOIN t2 ON |` → columnref after ON
[ ] `SELECT * FROM t1 JOIN t2 USING (|)` → columnref after USING (
[ ] `SELECT * FROM t1 |` → JOIN keywords (JOIN, LEFT, RIGHT, INNER, CROSS, NATURAL, STRAIGHT_JOIN)
```

### 3.4 WHERE, GROUP BY, HAVING

```
[ ] `SELECT * FROM t WHERE |` → columnref after WHERE
[ ] `SELECT * FROM t WHERE a = 1 AND |` → columnref after AND
[ ] `SELECT * FROM t WHERE a = 1 OR |` → columnref after OR
[ ] `SELECT * FROM t GROUP BY |` → columnref after GROUP BY
[ ] `SELECT * FROM t GROUP BY a, |` → columnref after comma
[ ] `SELECT * FROM t GROUP BY a |` → keyword candidates (HAVING, ORDER, LIMIT, WITH ROLLUP)
[ ] `SELECT * FROM t HAVING |` → columnref after HAVING
```

### 3.5 ORDER BY, LIMIT, DISTINCT

```
[ ] `SELECT * FROM t ORDER BY |` → columnref after ORDER BY
[ ] `SELECT * FROM t ORDER BY a, |` → columnref after comma
[ ] `SELECT * FROM t ORDER BY a |` → keyword candidates (ASC, DESC, LIMIT, comma)
[ ] `SELECT * FROM t LIMIT |` → no specific candidates (numeric context)
[ ] `SELECT * FROM t LIMIT 10 OFFSET |` → no specific candidates
```

### 3.6 Set Operations & FOR UPDATE

```
[ ] `SELECT a FROM t UNION |` → keyword candidates (ALL, SELECT)
[ ] `SELECT a FROM t UNION ALL |` → keyword candidate (SELECT)
[ ] `SELECT a FROM t INTERSECT |` → keyword candidates (ALL, SELECT)
[ ] `SELECT a FROM t EXCEPT |` → keyword candidates (ALL, SELECT)
[ ] `SELECT * FROM t FOR |` → keyword candidates (UPDATE, SHARE)
[ ] `SELECT * FROM t FOR UPDATE |` → keyword candidates (OF, NOWAIT, SKIP)
```

### 3.7 CTE (WITH Clause)

```
[ ] `WITH |` → keyword candidate (RECURSIVE) + identifier context for CTE name
[ ] `WITH cte AS (|)` → keyword candidate (SELECT)
[ ] `WITH cte AS (SELECT * FROM t) SELECT |` → columnref (CTE columns available)
[ ] `WITH cte AS (SELECT * FROM t) SELECT * FROM |` → table_ref (CTE name available)
[ ] `WITH RECURSIVE cte(|)` → identifier context for column names
```

### 3.8 Window Functions & Index Hints

```
[ ] `SELECT a, ROW_NUMBER() OVER (|)` → keyword candidates (PARTITION, ORDER)
[ ] `SELECT a, SUM(b) OVER (PARTITION BY |)` → columnref
[ ] `SELECT a, SUM(b) OVER (ORDER BY |)` → columnref
[ ] `SELECT a, SUM(b) OVER (ORDER BY a ROWS |)` → keyword candidates (BETWEEN, UNBOUNDED, CURRENT)
[ ] `SELECT * FROM t USE INDEX (|)` → index_ref
[ ] `SELECT * FROM t FORCE INDEX (|)` → index_ref
[ ] `SELECT * FROM t IGNORE INDEX (|)` → index_ref
```

## Phase 4: DML Instrumentation

### 4.1 INSERT

```
[ ] `INSERT INTO |` → table_ref after INTO
[ ] `INSERT INTO t (|)` → columnref for table t
[ ] `INSERT INTO t (a, |)` → columnref after comma
[ ] `INSERT INTO t VALUES (|)` → no specific candidates (value context)
[ ] `INSERT INTO t |` → keyword candidates (VALUES, SET, SELECT, PARTITION)
[ ] `INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE |` → columnref
[ ] `INSERT INTO t SET |` → columnref (assignment context)
[ ] `INSERT INTO t SELECT |` → columnref (INSERT SELECT)
[ ] `REPLACE INTO |` → table_ref
```

### 4.2 UPDATE

```
[ ] `UPDATE |` → table_ref
[ ] `UPDATE t SET |` → columnref for table t
[ ] `UPDATE t SET a = 1, |` → columnref after comma
[ ] `UPDATE t SET a = 1 WHERE |` → columnref
[ ] `UPDATE t SET a = 1 ORDER BY |` → columnref
[ ] `UPDATE t1 JOIN t2 ON t1.a = t2.a SET |` → columnref from both tables
```

### 4.3 DELETE & LOAD DATA & CALL

```
[ ] `DELETE FROM |` → table_ref
[ ] `DELETE FROM t WHERE |` → columnref for table t
[ ] `DELETE FROM t ORDER BY |` → columnref
[ ] `DELETE t1 FROM t1 JOIN t2 ON t1.a = t2.a WHERE |` → columnref from both tables
[ ] `LOAD DATA INFILE 'f' INTO TABLE |` → table_ref
[ ] `CALL |` → procedure_ref
```

## Phase 5: DDL Instrumentation

### 5.1 CREATE TABLE

```
[ ] `CREATE TABLE |` → identifier context (no specific candidates)
[ ] `CREATE TABLE t (a INT, |)` → keyword candidates for column/constraint start (PRIMARY, UNIQUE, INDEX, KEY, FOREIGN, CHECK, CONSTRAINT)
[ ] `CREATE TABLE t (a INT |)` → keyword candidates for column options (NOT, NULL, DEFAULT, AUTO_INCREMENT, PRIMARY, UNIQUE, COMMENT, COLLATE, REFERENCES, CHECK, GENERATED)
[ ] `CREATE TABLE t (a INT) |` → keyword candidates for table options (ENGINE, DEFAULT, CHARSET, COLLATE, COMMENT, AUTO_INCREMENT, ROW_FORMAT, PARTITION)
[ ] `CREATE TABLE t (a INT) ENGINE=|` → keyword candidates for engines (InnoDB, MyISAM, MEMORY, etc.)
[ ] `CREATE TABLE t (a |)` → type candidates (INT, VARCHAR, TEXT, BLOB, DATE, DATETIME, DECIMAL, etc.)
[ ] `CREATE TABLE t (FOREIGN KEY (a) REFERENCES |)` → table_ref
[ ] `CREATE TABLE t LIKE |` → table_ref
[ ] `CREATE TABLE t (a INT GENERATED ALWAYS AS (|))` → expression context (columnref, func_name)
```

### 5.2 ALTER TABLE

```
[ ] `ALTER TABLE |` → table_ref
[ ] `ALTER TABLE t |` → keyword candidates (ADD, DROP, MODIFY, CHANGE, RENAME, ALTER, CONVERT, ENGINE, DEFAULT, ORDER, ALGORITHM, LOCK, FORCE, ADD PARTITION, DROP PARTITION)
[ ] `ALTER TABLE t ADD |` → keyword candidates (COLUMN, INDEX, KEY, UNIQUE, PRIMARY, FOREIGN, CONSTRAINT, CHECK, PARTITION, SPATIAL, FULLTEXT)
[ ] `ALTER TABLE t ADD COLUMN |` → identifier context
[ ] `ALTER TABLE t DROP |` → keyword candidates (COLUMN, INDEX, KEY, FOREIGN, PRIMARY, CHECK, CONSTRAINT, PARTITION)
[ ] `ALTER TABLE t DROP COLUMN |` → columnref for table t
[ ] `ALTER TABLE t DROP INDEX |` → index_ref for table t
[ ] `ALTER TABLE t DROP FOREIGN KEY |` → constraint_ref
[ ] `ALTER TABLE t DROP CONSTRAINT |` → constraint_ref (generic, MySQL 8.0.16+)
[ ] `ALTER TABLE t MODIFY |` → columnref
[ ] `ALTER TABLE t MODIFY COLUMN |` → columnref
[ ] `ALTER TABLE t CHANGE |` → columnref (old name)
[ ] `ALTER TABLE t RENAME TO |` → identifier context
[ ] `ALTER TABLE t RENAME COLUMN |` → columnref
[ ] `ALTER TABLE t RENAME INDEX |` → index_ref
[ ] `ALTER TABLE t ADD INDEX idx (|)` → columnref
[ ] `ALTER TABLE t CONVERT TO CHARACTER SET |` → charset candidates
[ ] `ALTER TABLE t ALGORITHM=|` → keyword candidates (DEFAULT, INPLACE, COPY, INSTANT)
```

### 5.3 CREATE/DROP Index, View, Database

```
[ ] `CREATE INDEX idx ON |` → table_ref
[ ] `CREATE INDEX idx ON t (|)` → columnref for table t
[ ] `CREATE UNIQUE INDEX idx ON |` → table_ref
[ ] `DROP INDEX |` → index_ref
[ ] `DROP INDEX idx ON |` → table_ref
[ ] `CREATE VIEW |` → identifier context
[ ] `CREATE VIEW v AS |` → keyword candidate (SELECT)
[ ] `CREATE DEFINER=|` → keyword candidate (CURRENT_USER) + user context
[ ] `ALTER VIEW v AS |` → keyword candidate (SELECT)
[ ] `DROP VIEW |` → view_ref
[ ] `CREATE DATABASE |` → identifier context
[ ] `DROP DATABASE |` → database_ref
[ ] `DROP TABLE |` → table_ref
[ ] `DROP TABLE IF EXISTS |` → table_ref
[ ] `TRUNCATE TABLE |` → table_ref
[ ] `RENAME TABLE |` → table_ref
[ ] `RENAME TABLE t TO |` → identifier context
[ ] `DESCRIBE |` → table_ref
[ ] `DESC |` → table_ref
```

## Phase 6: Routine/Trigger/Event Instrumentation

### 6.1 Functions & Procedures

```
[ ] `CREATE FUNCTION |` → identifier context
[ ] `CREATE FUNCTION f(|)` → keyword candidates for param direction (IN, OUT, INOUT) + type context
[ ] `CREATE FUNCTION f() RETURNS |` → type candidates
[ ] `CREATE FUNCTION f() |` → keyword candidates for characteristics (DETERMINISTIC, NO SQL, READS SQL DATA, MODIFIES SQL DATA, COMMENT, LANGUAGE, SQL SECURITY)
[ ] `DROP FUNCTION |` → function_ref
[ ] `DROP FUNCTION IF EXISTS |` → function_ref
[ ] `CREATE PROCEDURE |` → identifier context
[ ] `DROP PROCEDURE |` → procedure_ref
[ ] `ALTER FUNCTION |` → function_ref
[ ] `ALTER PROCEDURE |` → procedure_ref
```

### 6.2 Triggers & Events

```
[ ] `CREATE TRIGGER |` → identifier context
[ ] `CREATE TRIGGER trg |` → keyword candidates (BEFORE, AFTER)
[ ] `CREATE TRIGGER trg BEFORE |` → keyword candidates (INSERT, UPDATE, DELETE)
[ ] `CREATE TRIGGER trg BEFORE INSERT ON |` → table_ref
[ ] `DROP TRIGGER |` → trigger_ref
[ ] `CREATE EVENT |` → identifier context
[ ] `CREATE EVENT ev ON SCHEDULE |` → keyword candidates (AT, EVERY)
[ ] `DROP EVENT |` → event_ref
[ ] `ALTER EVENT |` → event_ref
```

### 6.3 Transaction, LOCK & Table Maintenance

```
[ ] `BEGIN |` → keyword candidates (WORK)
[ ] `START TRANSACTION |` → keyword candidates (WITH CONSISTENT SNAPSHOT, READ ONLY, READ WRITE)
[ ] `COMMIT |` → keyword candidates (AND, WORK)
[ ] `ROLLBACK |` → keyword candidates (TO, WORK)
[ ] `ROLLBACK TO |` → keyword candidate (SAVEPOINT)
[ ] `SAVEPOINT |` → identifier context
[ ] `RELEASE SAVEPOINT |` → identifier context
[ ] `LOCK TABLES |` → table_ref
[ ] `LOCK TABLES t |` → keyword candidates (READ, WRITE)
[ ] `ANALYZE TABLE |` → table_ref
[ ] `OPTIMIZE TABLE |` → table_ref
[ ] `CHECK TABLE |` → table_ref
[ ] `REPAIR TABLE |` → table_ref
[ ] `FLUSH |` → keyword candidates (PRIVILEGES, TABLES, LOGS, STATUS, HOSTS)
```

## Phase 7: Session/Utility Instrumentation

### 7.1 SET & SHOW

```
[ ] `SET |` → variable candidates (@@, @, GLOBAL, SESSION, NAMES, CHARACTER, PASSWORD) + keyword candidates
[ ] `SET GLOBAL |` → variable candidates (system variables)
[ ] `SET SESSION |` → variable candidates
[ ] `SET NAMES |` → charset candidates
[ ] `SET CHARACTER SET |` → charset candidates
[ ] `SHOW |` → keyword candidates (TABLES, COLUMNS, INDEX, DATABASES, CREATE, STATUS, VARIABLES, PROCESSLIST, GRANTS, WARNINGS, ERRORS, ENGINE, etc.)
[ ] `SHOW CREATE TABLE |` → table_ref
[ ] `SHOW CREATE VIEW |` → view_ref
[ ] `SHOW CREATE FUNCTION |` → function_ref
[ ] `SHOW CREATE PROCEDURE |` → procedure_ref
[ ] `SHOW COLUMNS FROM |` → table_ref
[ ] `SHOW INDEX FROM |` → table_ref
[ ] `SHOW TABLES FROM |` → database_ref
```

### 7.2 USE, GRANT, EXPLAIN

```
[ ] `USE |` → database_ref
[ ] `GRANT |` → keyword candidates (ALL, SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, DROP, INDEX, EXECUTE, etc.)
[ ] `GRANT SELECT ON |` → table_ref (or database.*)
[ ] `GRANT SELECT ON t TO |` → user context
[ ] `REVOKE SELECT ON |` → table_ref
[ ] `EXPLAIN |` → keyword candidates (SELECT, INSERT, UPDATE, DELETE, FORMAT)
[ ] `EXPLAIN SELECT |` → same as SELECT instrumentation
[ ] `PREPARE stmt FROM |` → string context
[ ] `EXECUTE |` → prepared statement name
[ ] `DEALLOCATE PREPARE |` → prepared statement name
[ ] `DO |` → expression context (columnref, func_name)
```

## Phase 8: Expression Instrumentation

### 8.1 Function & Type Names

```
[ ] `SELECT |()` context → func_name candidates (built-in functions: COUNT, SUM, AVG, MAX, MIN, CONCAT, SUBSTRING, TRIM, NOW, IF, IFNULL, COALESCE, CAST, CONVERT, etc.)
[ ] `SELECT CAST(a AS |)` → type candidates (CHAR, SIGNED, UNSIGNED, DECIMAL, DATE, DATETIME, TIME, JSON, BINARY)
[ ] `SELECT CONVERT(a, |)` → type candidates
[ ] `SELECT CONVERT(a USING |)` → charset candidates
```

### 8.2 Expression Contexts

```
[ ] `SELECT a + |` → columnref, func_name (expression continuation)
[ ] `SELECT CASE WHEN |` → columnref (CASE WHEN condition)
[ ] `SELECT CASE WHEN a THEN |` → columnref, literal (CASE THEN result)
[ ] `SELECT CASE a WHEN |` → literal context (CASE WHEN value)
[ ] `SELECT * FROM t WHERE a IN (|)` → columnref, literal (IN list or subquery)
[ ] `SELECT * FROM t WHERE a BETWEEN |` → columnref, literal (BETWEEN lower bound)
[ ] `SELECT * FROM t WHERE a BETWEEN 1 AND |` → columnref, literal (BETWEEN upper bound)
[ ] `SELECT * FROM t WHERE EXISTS (|)` → keyword candidate (SELECT)
[ ] `SELECT * FROM t WHERE a IS |` → keyword candidates (NULL, NOT, TRUE, FALSE, UNKNOWN)
```

## Phase 9: Integration Tests

### 9.1 Multi-Table Schema Tests

```
[ ] Column completion scoped to correct table in JOIN: `SELECT t1.| FROM t1 JOIN t2 ON ...` → only t1 columns
[ ] Column completion from all tables when unqualified: `SELECT | FROM t1 JOIN t2 ON ...` → columns from both tables
[ ] Table alias completion: `SELECT x.| FROM t AS x` → columns of t via alias x
[ ] View column completion: `SELECT | FROM v` → columns of view v
[ ] CTE column completion: `WITH cte AS (SELECT a FROM t) SELECT | FROM cte` → column a
[ ] Database-qualified table: `SELECT * FROM db.| ` → tables in database db
```

### 9.2 Edge Cases

```
[ ] Cursor at beginning of SQL: `|SELECT * FROM t` → top-level keywords
[ ] Cursor in middle of identifier: `SELECT us|ers FROM t` → prefix "us" filters candidates
[ ] Cursor after semicolon (multi-statement): `SELECT 1; SELECT |` → new statement context
[ ] Empty SQL: `|` → top-level keywords
[ ] Whitespace only: `   |` → top-level keywords
[ ] Very long SQL with cursor in middle
[ ] SQL with syntax errors before cursor: completion still works for valid prefix
[ ] Backtick-quoted identifiers: `SELECT \`| FROM t` → column candidates
```

### 9.3 Complex SQL Patterns

```
[ ] Nested subquery column completion: `SELECT * FROM t WHERE a IN (SELECT | FROM t2)` → t2 columns
[ ] Correlated subquery: `SELECT *, (SELECT | FROM t2 WHERE t2.a = t1.a) FROM t1` → t2 columns
[ ] UNION: `SELECT a FROM t1 UNION SELECT | FROM t2` → t2 columns
[ ] Multiple JOINs: `SELECT | FROM t1 JOIN t2 ON ... JOIN t3 ON ...` → columns from all 3 tables
[ ] INSERT ... SELECT: `INSERT INTO t1 SELECT | FROM t2` → t2 columns
[ ] Complex ALTER: `ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (|) REFERENCES ...` → t columns
```
