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
[x] Incomplete SQL: `SELECT * FROM ` (trailing space) → insert placeholder, re-collect
[x] Truncated mid-keyword: `SELE` → prefix-filter against keywords
[x] Truncated after comma: `SELECT a,` → insert placeholder column
[x] Truncated after operator: `WHERE a >` → insert placeholder expression
[x] Multiple placeholder strategies tried in order
[x] Fallback returns best-effort results when no strategy succeeds
[x] Placeholder insertion does not corrupt the candidate set from the initial pass
```

## Phase 3: SELECT Statement Instrumentation

### 3.1 SELECT Target List

```
[x] `SELECT |` → columnref, func_name, literal keywords (DISTINCT, ALL, *)
[x] `SELECT a, |` → columnref, func_name after comma
[x] `SELECT a, b, |` → columnref after multiple commas
[x] `SELECT * FROM t WHERE a > (SELECT |)` → columnref in subquery
[x] `SELECT DISTINCT |` → columnref after DISTINCT
```

### 3.2 FROM Clause

```
[x] `SELECT * FROM |` → table_ref (tables, views, databases)
[x] `SELECT * FROM db.|` → table_ref qualified with database
[x] `SELECT * FROM t1, |` → table_ref after comma (multi-table)
[x] `SELECT * FROM (SELECT * FROM |)` → table_ref in derived table
[x] `SELECT * FROM t |` → keyword candidates (WHERE, JOIN, LEFT, RIGHT, CROSS, NATURAL, STRAIGHT_JOIN, ORDER, GROUP, HAVING, LIMIT, UNION, FOR)
[x] `SELECT * FROM t AS |` → no specific candidates (alias context)
```

### 3.3 JOIN Clauses

```
[x] `SELECT * FROM t1 JOIN |` → table_ref after JOIN
[x] `SELECT * FROM t1 LEFT JOIN |` → table_ref after LEFT JOIN
[x] `SELECT * FROM t1 RIGHT JOIN |` → table_ref after RIGHT JOIN
[x] `SELECT * FROM t1 CROSS JOIN |` → table_ref after CROSS JOIN
[x] `SELECT * FROM t1 NATURAL JOIN |` → table_ref after NATURAL JOIN
[x] `SELECT * FROM t1 STRAIGHT_JOIN |` → table_ref after STRAIGHT_JOIN
[x] `SELECT * FROM t1 JOIN t2 ON |` → columnref after ON
[x] `SELECT * FROM t1 JOIN t2 USING (|)` → columnref after USING (
[x] `SELECT * FROM t1 |` → JOIN keywords (JOIN, LEFT, RIGHT, INNER, CROSS, NATURAL, STRAIGHT_JOIN)
```

### 3.4 WHERE, GROUP BY, HAVING

```
[x] `SELECT * FROM t WHERE |` → columnref after WHERE
[x] `SELECT * FROM t WHERE a = 1 AND |` → columnref after AND
[x] `SELECT * FROM t WHERE a = 1 OR |` → columnref after OR
[x] `SELECT * FROM t GROUP BY |` → columnref after GROUP BY
[x] `SELECT * FROM t GROUP BY a, |` → columnref after comma
[x] `SELECT * FROM t GROUP BY a |` → keyword candidates (HAVING, ORDER, LIMIT, WITH ROLLUP)
[x] `SELECT * FROM t HAVING |` → columnref after HAVING
```

### 3.5 ORDER BY, LIMIT, DISTINCT

```
[x] `SELECT * FROM t ORDER BY |` → columnref after ORDER BY
[x] `SELECT * FROM t ORDER BY a, |` → columnref after comma
[x] `SELECT * FROM t ORDER BY a |` → keyword candidates (ASC, DESC, LIMIT, comma)
[x] `SELECT * FROM t LIMIT |` → no specific candidates (numeric context)
[x] `SELECT * FROM t LIMIT 10 OFFSET |` → no specific candidates
```

### 3.6 Set Operations & FOR UPDATE

```
[x] `SELECT a FROM t UNION |` → keyword candidates (ALL, SELECT)
[x] `SELECT a FROM t UNION ALL |` → keyword candidate (SELECT)
[x] `SELECT a FROM t INTERSECT |` → keyword candidates (ALL, SELECT)
[x] `SELECT a FROM t EXCEPT |` → keyword candidates (ALL, SELECT)
[x] `SELECT * FROM t FOR |` → keyword candidates (UPDATE, SHARE)
[x] `SELECT * FROM t FOR UPDATE |` → keyword candidates (OF, NOWAIT, SKIP)
```

### 3.7 CTE (WITH Clause)

```
[x] `WITH |` → keyword candidate (RECURSIVE) + identifier context for CTE name
[x] `WITH cte AS (|)` → keyword candidate (SELECT)
[x] `WITH cte AS (SELECT * FROM t) SELECT |` → columnref (CTE columns available)
[x] `WITH cte AS (SELECT * FROM t) SELECT * FROM |` → table_ref (CTE name available)
[x] `WITH RECURSIVE cte(|)` → identifier context for column names
```

### 3.8 Window Functions & Index Hints

```
[x] `SELECT a, ROW_NUMBER() OVER (|)` → keyword candidates (PARTITION, ORDER)
[x] `SELECT a, SUM(b) OVER (PARTITION BY |)` → columnref
[x] `SELECT a, SUM(b) OVER (ORDER BY |)` → columnref
[x] `SELECT a, SUM(b) OVER (ORDER BY a ROWS |)` → keyword candidates (BETWEEN, UNBOUNDED, CURRENT)
[x] `SELECT * FROM t USE INDEX (|)` → index_ref
[x] `SELECT * FROM t FORCE INDEX (|)` → index_ref
[x] `SELECT * FROM t IGNORE INDEX (|)` → index_ref
```

## Phase 4: DML Instrumentation

### 4.1 INSERT

```
[x] `INSERT INTO |` → table_ref after INTO
[x] `INSERT INTO t (|)` → columnref for table t
[x] `INSERT INTO t (a, |)` → columnref after comma
[x] `INSERT INTO t VALUES (|)` → no specific candidates (value context)
[x] `INSERT INTO t |` → keyword candidates (VALUES, SET, SELECT, PARTITION)
[x] `INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE |` → columnref
[x] `INSERT INTO t SET |` → columnref (assignment context)
[x] `INSERT INTO t SELECT |` → columnref (INSERT SELECT)
[x] `REPLACE INTO |` → table_ref
```

### 4.2 UPDATE

```
[x] `UPDATE |` → table_ref
[x] `UPDATE t SET |` → columnref for table t
[x] `UPDATE t SET a = 1, |` → columnref after comma
[x] `UPDATE t SET a = 1 WHERE |` → columnref
[x] `UPDATE t SET a = 1 ORDER BY |` → columnref
[x] `UPDATE t1 JOIN t2 ON t1.a = t2.a SET |` → columnref from both tables
```

### 4.3 DELETE & LOAD DATA & CALL

```
[x] `DELETE FROM |` → table_ref
[x] `DELETE FROM t WHERE |` → columnref for table t
[x] `DELETE FROM t ORDER BY |` → columnref
[x] `DELETE t1 FROM t1 JOIN t2 ON t1.a = t2.a WHERE |` → columnref from both tables
[x] `LOAD DATA INFILE 'f' INTO TABLE |` → table_ref
[x] `CALL |` → procedure_ref
```

## Phase 5: DDL Instrumentation

### 5.1 CREATE TABLE

```
[x] `CREATE TABLE |` → identifier context (no specific candidates)
[x] `CREATE TABLE t (a INT, |)` → keyword candidates for column/constraint start (PRIMARY, UNIQUE, INDEX, KEY, FOREIGN, CHECK, CONSTRAINT)
[x] `CREATE TABLE t (a INT |)` → keyword candidates for column options (NOT, NULL, DEFAULT, AUTO_INCREMENT, PRIMARY, UNIQUE, COMMENT, COLLATE, REFERENCES, CHECK, GENERATED)
[x] `CREATE TABLE t (a INT) |` → keyword candidates for table options (ENGINE, DEFAULT, CHARSET, COLLATE, COMMENT, AUTO_INCREMENT, ROW_FORMAT, PARTITION)
[x] `CREATE TABLE t (a INT) ENGINE=|` → keyword candidates for engines (InnoDB, MyISAM, MEMORY, etc.)
[x] `CREATE TABLE t (a |)` → type candidates (INT, VARCHAR, TEXT, BLOB, DATE, DATETIME, DECIMAL, etc.)
[x] `CREATE TABLE t (FOREIGN KEY (a) REFERENCES |)` → table_ref
[x] `CREATE TABLE t LIKE |` → table_ref
[x] `CREATE TABLE t (a INT GENERATED ALWAYS AS (|))` → expression context (columnref, func_name)
```

### 5.2 ALTER TABLE

```
[x] `ALTER TABLE |` → table_ref
[x] `ALTER TABLE t |` → keyword candidates (ADD, DROP, MODIFY, CHANGE, RENAME, ALTER, CONVERT, ENGINE, DEFAULT, ORDER, ALGORITHM, LOCK, FORCE, ADD PARTITION, DROP PARTITION)
[x] `ALTER TABLE t ADD |` → keyword candidates (COLUMN, INDEX, KEY, UNIQUE, PRIMARY, FOREIGN, CONSTRAINT, CHECK, PARTITION, SPATIAL, FULLTEXT)
[x] `ALTER TABLE t ADD COLUMN |` → identifier context
[x] `ALTER TABLE t DROP |` → keyword candidates (COLUMN, INDEX, KEY, FOREIGN, PRIMARY, CHECK, CONSTRAINT, PARTITION)
[x] `ALTER TABLE t DROP COLUMN |` → columnref for table t
[x] `ALTER TABLE t DROP INDEX |` → index_ref for table t
[x] `ALTER TABLE t DROP FOREIGN KEY |` → constraint_ref
[x] `ALTER TABLE t DROP CONSTRAINT |` → constraint_ref (generic, MySQL 8.0.16+)
[x] `ALTER TABLE t MODIFY |` → columnref
[x] `ALTER TABLE t MODIFY COLUMN |` → columnref
[x] `ALTER TABLE t CHANGE |` → columnref (old name)
[x] `ALTER TABLE t RENAME TO |` → identifier context
[x] `ALTER TABLE t RENAME COLUMN |` → columnref
[x] `ALTER TABLE t RENAME INDEX |` → index_ref
[x] `ALTER TABLE t ADD INDEX idx (|)` → columnref
[x] `ALTER TABLE t CONVERT TO CHARACTER SET |` → charset candidates
[x] `ALTER TABLE t ALGORITHM=|` → keyword candidates (DEFAULT, INPLACE, COPY, INSTANT)
```

### 5.3 CREATE/DROP Index, View, Database

```
[x] `CREATE INDEX idx ON |` → table_ref
[x] `CREATE INDEX idx ON t (|)` → columnref for table t
[x] `CREATE UNIQUE INDEX idx ON |` → table_ref
[x] `DROP INDEX |` → index_ref
[x] `DROP INDEX idx ON |` → table_ref
[x] `CREATE VIEW |` → identifier context
[x] `CREATE VIEW v AS |` → keyword candidate (SELECT)
[x] `CREATE DEFINER=|` → keyword candidate (CURRENT_USER) + user context
[x] `ALTER VIEW v AS |` → keyword candidate (SELECT)
[x] `DROP VIEW |` → view_ref
[x] `CREATE DATABASE |` → identifier context
[x] `DROP DATABASE |` → database_ref
[x] `DROP TABLE |` → table_ref
[x] `DROP TABLE IF EXISTS |` → table_ref
[x] `TRUNCATE TABLE |` → table_ref
[x] `RENAME TABLE |` → table_ref
[x] `RENAME TABLE t TO |` → identifier context
[x] `DESCRIBE |` → table_ref
[x] `DESC |` → table_ref
```

## Phase 6: Routine/Trigger/Event Instrumentation

### 6.1 Functions & Procedures

```
[x] `CREATE FUNCTION |` → identifier context
[x] `CREATE FUNCTION f(|)` → keyword candidates for param direction (IN, OUT, INOUT) + type context
[x] `CREATE FUNCTION f() RETURNS |` → type candidates
[x] `CREATE FUNCTION f() |` → keyword candidates for characteristics (DETERMINISTIC, NO SQL, READS SQL DATA, MODIFIES SQL DATA, COMMENT, LANGUAGE, SQL SECURITY)
[x] `DROP FUNCTION |` → function_ref
[x] `DROP FUNCTION IF EXISTS |` → function_ref
[x] `CREATE PROCEDURE |` → identifier context
[x] `DROP PROCEDURE |` → procedure_ref
[x] `ALTER FUNCTION |` → function_ref
[x] `ALTER PROCEDURE |` → procedure_ref
```

### 6.2 Triggers & Events

```
[x] `CREATE TRIGGER |` → identifier context
[x] `CREATE TRIGGER trg |` → keyword candidates (BEFORE, AFTER)
[x] `CREATE TRIGGER trg BEFORE |` → keyword candidates (INSERT, UPDATE, DELETE)
[x] `CREATE TRIGGER trg BEFORE INSERT ON |` → table_ref
[x] `DROP TRIGGER |` → trigger_ref
[x] `CREATE EVENT |` → identifier context
[x] `CREATE EVENT ev ON SCHEDULE |` → keyword candidates (AT, EVERY)
[x] `DROP EVENT |` → event_ref
[x] `ALTER EVENT |` → event_ref
```

### 6.3 Transaction, LOCK & Table Maintenance

```
[x] `BEGIN |` → keyword candidates (WORK)
[x] `START TRANSACTION |` → keyword candidates (WITH CONSISTENT SNAPSHOT, READ ONLY, READ WRITE)
[x] `COMMIT |` → keyword candidates (AND, WORK)
[x] `ROLLBACK |` → keyword candidates (TO, WORK)
[x] `ROLLBACK TO |` → keyword candidate (SAVEPOINT)
[x] `SAVEPOINT |` → identifier context
[x] `RELEASE SAVEPOINT |` → identifier context
[x] `LOCK TABLES |` → table_ref
[x] `LOCK TABLES t |` → keyword candidates (READ, WRITE)
[x] `ANALYZE TABLE |` → table_ref
[x] `OPTIMIZE TABLE |` → table_ref
[x] `CHECK TABLE |` → table_ref
[x] `REPAIR TABLE |` → table_ref
[x] `FLUSH |` → keyword candidates (PRIVILEGES, TABLES, LOGS, STATUS, HOSTS)
```

## Phase 7: Session/Utility Instrumentation

### 7.1 SET & SHOW

```
[x] `SET |` → variable candidates (@@, @, GLOBAL, SESSION, NAMES, CHARACTER, PASSWORD) + keyword candidates
[x] `SET GLOBAL |` → variable candidates (system variables)
[x] `SET SESSION |` → variable candidates
[x] `SET NAMES |` → charset candidates
[x] `SET CHARACTER SET |` → charset candidates
[x] `SHOW |` → keyword candidates (TABLES, COLUMNS, INDEX, DATABASES, CREATE, STATUS, VARIABLES, PROCESSLIST, GRANTS, WARNINGS, ERRORS, ENGINE, etc.)
[x] `SHOW CREATE TABLE |` → table_ref
[x] `SHOW CREATE VIEW |` → view_ref
[x] `SHOW CREATE FUNCTION |` → function_ref
[x] `SHOW CREATE PROCEDURE |` → procedure_ref
[x] `SHOW COLUMNS FROM |` → table_ref
[x] `SHOW INDEX FROM |` → table_ref
[x] `SHOW TABLES FROM |` → database_ref
```

### 7.2 USE, GRANT, EXPLAIN

```
[x] `USE |` → database_ref
[x] `GRANT |` → keyword candidates (ALL, SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, DROP, INDEX, EXECUTE, etc.)
[x] `GRANT SELECT ON |` → table_ref (or database.*)
[x] `GRANT SELECT ON t TO |` → user context
[x] `REVOKE SELECT ON |` → table_ref
[x] `EXPLAIN |` → keyword candidates (SELECT, INSERT, UPDATE, DELETE, FORMAT)
[x] `EXPLAIN SELECT |` → same as SELECT instrumentation
[x] `PREPARE stmt FROM |` → string context
[x] `EXECUTE |` → prepared statement name
[x] `DEALLOCATE PREPARE |` → prepared statement name
[x] `DO |` → expression context (columnref, func_name)
```

## Phase 8: Expression Instrumentation

### 8.1 Function & Type Names

```
[x] `SELECT |()` context → func_name candidates (built-in functions: COUNT, SUM, AVG, MAX, MIN, CONCAT, SUBSTRING, TRIM, NOW, IF, IFNULL, COALESCE, CAST, CONVERT, etc.)
[x] `SELECT CAST(a AS |)` → type candidates (CHAR, SIGNED, UNSIGNED, DECIMAL, DATE, DATETIME, TIME, JSON, BINARY)
[x] `SELECT CONVERT(a, |)` → type candidates
[x] `SELECT CONVERT(a USING |)` → charset candidates
```

### 8.2 Expression Contexts

```
[x] `SELECT a + |` → columnref, func_name (expression continuation)
[x] `SELECT CASE WHEN |` → columnref (CASE WHEN condition)
[x] `SELECT CASE WHEN a THEN |` → columnref, literal (CASE THEN result)
[x] `SELECT CASE a WHEN |` → literal context (CASE WHEN value)
[x] `SELECT * FROM t WHERE a IN (|)` → columnref, literal (IN list or subquery)
[x] `SELECT * FROM t WHERE a BETWEEN |` → columnref, literal (BETWEEN lower bound)
[x] `SELECT * FROM t WHERE a BETWEEN 1 AND |` → columnref, literal (BETWEEN upper bound)
[x] `SELECT * FROM t WHERE EXISTS (|)` → keyword candidate (SELECT)
[x] `SELECT * FROM t WHERE a IS |` → keyword candidates (NULL, NOT, TRUE, FALSE, UNKNOWN)
```

## Phase 9: Integration Tests

### 9.1 Multi-Table Schema Tests

```
[x] Column completion scoped to correct table in JOIN: `SELECT t1.| FROM t1 JOIN t2 ON ...` → only t1 columns
[x] Column completion from all tables when unqualified: `SELECT | FROM t1 JOIN t2 ON ...` → columns from both tables
[x] Table alias completion: `SELECT x.| FROM t AS x` → columns of t via alias x
[x] View column completion: `SELECT | FROM v` → columns of view v
[x] CTE column completion: `WITH cte AS (SELECT a FROM t) SELECT | FROM cte` → column a
[x] Database-qualified table: `SELECT * FROM db.| ` → tables in database db
```

### 9.2 Edge Cases

```
[x] Cursor at beginning of SQL: `|SELECT * FROM t` → top-level keywords
[x] Cursor in middle of identifier: `SELECT us|ers FROM t` → prefix "us" filters candidates
[x] Cursor after semicolon (multi-statement): `SELECT 1; SELECT |` → new statement context
[x] Empty SQL: `|` → top-level keywords
[x] Whitespace only: `   |` → top-level keywords
[x] Very long SQL with cursor in middle
[x] SQL with syntax errors before cursor: completion still works for valid prefix
[x] Backtick-quoted identifiers: `SELECT \`| FROM t` → column candidates
```

### 9.3 Complex SQL Patterns

```
[x] Nested subquery column completion: `SELECT * FROM t WHERE a IN (SELECT | FROM t2)` → t2 columns
[x] Correlated subquery: `SELECT *, (SELECT | FROM t2 WHERE t2.a = t1.a) FROM t1` → t2 columns
[x] UNION: `SELECT a FROM t1 UNION SELECT | FROM t2` → t2 columns
[x] Multiple JOINs: `SELECT | FROM t1 JOIN t2 ON ... JOIN t3 ON ...` → columns from all 3 tables
[x] INSERT ... SELECT: `INSERT INTO t1 SELECT | FROM t2` → t2 columns
[x] Complex ALTER: `ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (|) REFERENCES ...` → t columns
```
