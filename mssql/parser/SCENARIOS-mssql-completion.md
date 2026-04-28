# MSSQL SQL Autocompletion Scenarios

> Goal: Implement parser-native SQL autocompletion for MSSQL, matching PG/MySQL architecture — parser instrumentation with collectMode, completion module with Complete(sql, cursor, catalog) API
> Verification: Table-driven Go tests — Complete(sql, cursorPos, catalog) returns expected candidates with correct types
> Reference: PG completion (pg/completion/, pg/parser/complete.go), MySQL completion (mysql/completion/)
> Out of scope (v1): Service Broker, Resource Governor, Availability Groups, Endpoints, Full-text search, DBCC, BULK INSERT, CLR Assembly, XML Schema, External Tables

Status legend: `[ ]` pending, `[x]` passing, `[~]` partial

---

## Phase 1: Parser Completion Infrastructure

### 1.1 Completion Mode Fields & Entry Point

- [x] Parser struct has completion fields: completing bool, cursorOff int, candidates *CandidateSet, collecting bool, maxCollect int
- [x] collectMode() returns true when completing && collecting
- [x] checkCursor() sets collecting=true when p.cur.Loc >= cursorOff
- [x] CandidateSet struct with Tokens []int, Rules []RuleCandidate, seen/seenR dedup maps
- [x] RuleCandidate struct with Rule string (e.g., "columnref", "table_ref")
- [x] addTokenCandidate(tokType int) adds to CandidateSet.Tokens with dedup
- [x] addRuleCandidate(rule string) adds to CandidateSet.Rules with dedup
- [x] errCollecting sentinel error for unwinding parser stack
- [x] Collect(sql string, cursorOffset int) *CandidateSet — entry point with panic recovery
- [x] Collect returns non-nil CandidateSet even for empty input
- [x] All existing parser tests pass with completion fields present (completing=false)

### 1.2 TokenName & Basic Keyword Collection

- [x] TokenName(tokType int) returns SQL keyword string for token type
- [x] TokenName handles single-char tokens (';', '(', ')', ',', etc.)
- [x] TokenName handles all 221 MSSQL keywords
- [x] Empty input → keyword candidates for all top-level statements
- [x] After semicolon → keyword candidates for new statement
- [x] Cursor at EOF of empty string → top-level keywords (SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, DROP, etc.)

### 1.3 Statement-Level Dispatch Instrumentation

- [x] `|` (empty) → top-level statement keywords (SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, DROP, MERGE, WITH, EXEC, DECLARE, SET, IF, WHILE, BEGIN, PRINT, GRANT, REVOKE, DENY, TRUNCATE, USE, BACKUP, RESTORE)
- [x] `CREATE |` → object type keywords (TABLE, INDEX, VIEW, PROCEDURE, FUNCTION, TRIGGER, DATABASE, SCHEMA, TYPE, SEQUENCE, SYNONYM, STATISTICS, LOGIN, USER, ROLE)
- [x] `ALTER |` → object type keywords (TABLE, INDEX, VIEW, PROCEDURE, FUNCTION, TRIGGER, DATABASE, SCHEMA, SEQUENCE, LOGIN, USER, ROLE)
- [x] `DROP |` → object type keywords (TABLE, INDEX, VIEW, PROCEDURE, FUNCTION, TRIGGER, DATABASE, SCHEMA, TYPE, SEQUENCE, SYNONYM, STATISTICS, LOGIN, USER, ROLE, IF)
- [x] Build passes after 1.3

### 1.4 Bytebase-Facing Parser API Compatibility

- [x] `Tokenize(sql)` returns non-EOF tokens with stable `Loc`/`End` byte offsets for caret-token and quoted-identifier handling
- [x] `IsIdentTokenType(tokenType)` matches T-SQL parser identifier semantics: regular/bracketed/quoted identifiers and context keywords are accepted; core keywords are rejected unless quoted
- [x] Exported completion token aliases cover PG/MySQL adapter-style keyword checks (`SELECT`, `FROM`, `JOIN`, `ORDER`, `GO`, etc.)
- [x] Exported aliases include Bytebase object/qualifier and clause tokens (`DATABASE`, `SCHEMA`, `TABLE`, `VIEW`, `SEQUENCE`, `COLUMN`, `BY`, `ASC`, `DESC`, `OFFSET`, `FETCH`, etc.)
- [x] Prefix retry flow works: tokenize `SELECT * FROM dbo.Us|`, retry `Collect` at the `Us` token start, and receive `table_ref`
- [x] Prefix retry flow works for column references: `SELECT Na| FROM Employees` and `SELECT tableAlias.Na| ...` produce `columnref`
- [x] `CTEPositions` records `WITH` offsets for Bytebase virtual table extraction
- [x] `SelectAliasPositions` records alias token offsets for `AS alias`, implicit aliases, and T-SQL `alias = expr`
- [x] `GO` batch separator followed by cursor returns new top-level statement candidates
- [x] Public API tests use `parser_test` so Bytebase-facing behavior only depends on exported parser APIs

---

## Phase 2: Completion Module (API & Resolution)

### 2.1 Public API & Core Logic

- [x] mssql/completion/ package exists
- [x] Complete(sql, cursorOffset, catalog) returns []Candidate
- [x] Candidate struct has Text, Type (CandidateType), Definition, Comment fields
- [x] CandidateType enum: Keyword, Schema, Table, View, Column, Function, Procedure, Index, Trigger, Sequence, Type, Login, User, Role
- [x] Complete with nil catalog returns keyword-only candidates
- [x] Complete with empty sql returns top-level statement keywords
- [x] Prefix filtering: `SEL|` matches SELECT keyword
- [x] Prefix filtering is case-insensitive
- [x] Deduplication: same candidate not returned twice

### 2.2 Candidate Resolution

- [x] Token candidates → keyword strings (from TokenName mapping)
- [~] "table_ref" rule → parser emits rule; standalone `mssql/completion` catalog resolver is still a TODO
- [~] "columnref" rule → parser emits rule; table refs are extracted, but catalog-backed column resolution is still a TODO
- [~] "schema_ref" rule → parser emits rule; standalone catalog resolver is still a TODO
- [~] "func_name" rule → parser emits rule; standalone built-in/catalog function resolution is still a TODO
- [~] "proc_ref" rule → parser emits rule; standalone catalog resolver is still a TODO
- [~] "index_ref" rule → parser emits rule; standalone catalog resolver is still a TODO
- [~] "trigger_ref" rule → parser emits rule; standalone catalog resolver is still a TODO
- [x] "type_name" rule → T-SQL type keywords (INT, VARCHAR, NVARCHAR, TEXT, DATETIME, DECIMAL, BIT, etc.)
- [~] "database_ref" rule → parser emits rule; standalone catalog resolver is still a TODO
- [~] "sequence_ref" rule → parser emits rule; standalone catalog resolver is still a TODO
- [~] "login_ref" rule → parser emits rule; standalone catalog resolver is still a TODO
- [~] "user_ref" rule → parser emits rule; standalone catalog resolver is still a TODO
- [~] "role_ref" rule → parser emits rule; standalone catalog resolver is still a TODO

### 2.3 Table Reference Extraction

- [x] Extract table refs from simple SELECT: `SELECT * FROM t` → [{Table: "t"}]
- [x] Extract table refs with alias: `SELECT * FROM t AS x` → [{Table: "t", Alias: "x"}]
- [x] Extract table refs from JOIN: `FROM t1 JOIN t2 ON ...` → [{Table: "t1"}, {Table: "t2"}]
- [x] Extract table refs with schema: `FROM dbo.t` → [{Schema: "dbo", Table: "t"}]
- [x] Extract table refs from UPDATE: `UPDATE t SET ...` → [{Table: "t"}]
- [x] Extract table refs from INSERT: `INSERT INTO t ...` → [{Table: "t"}]
- [x] Extract table refs from DELETE: `DELETE FROM t ...` → [{Table: "t"}]
- [x] Extract table refs from MERGE: `MERGE INTO t ...` → [{Table: "t"}]
- [x] Fallback to lexer-based extraction when AST parsing fails

### 2.4 Tricky Completion (Fallback)

- [x] Incomplete SQL: `SELECT * FROM ` (trailing space) → insert placeholder, re-collect
- [x] Truncated mid-keyword: `SELE` → prefix-filter against keywords
- [x] Truncated after comma: `SELECT a,` → insert placeholder column
- [x] Truncated after operator: `WHERE a >` → insert placeholder expression
- [x] Multiple placeholder strategies tried in order
- [x] Fallback returns best-effort results when no strategy succeeds

### 2.5 Parser-Only Incomplete SQL Signals

- [x] `SELECT * FROM|` and `SELECT * FROM |` → `table_ref`
- [x] `SELECT * FROM dbo.|` → `table_ref`
- [x] `SELECT * FROM Employees JOIN|` → `table_ref`
- [x] `SELECT Id,|` and `SELECT Id, |` → `columnref`
- [x] `SELECT * FROM Employees WHERE Id =|` and trailing-space variant → `columnref`
- [x] `SELECT Id FROM Employees ORDER BY|` and trailing-space variant → `columnref`
- [x] `WITH cte AS (|` → `SELECT`
- [x] `WITH cte AS (...) SELECT * FROM|` → `table_ref`

---

## Phase 3: SELECT Statement Instrumentation

### 3.1 SELECT Target List

- [x] `SELECT |` → columnref, func_name, literal keywords (DISTINCT, TOP, ALL, *)
- [x] `SELECT a, |` → columnref, func_name after comma
- [x] `SELECT a, b, |` → columnref after multiple commas
- [x] `SELECT * FROM t WHERE a > (SELECT |)` → columnref in subquery
- [x] `SELECT DISTINCT |` → columnref after DISTINCT
- [x] `SELECT TOP |` → no specific candidates (numeric context)
- [x] `SELECT TOP 10 |` → columnref after TOP N

### 3.2 FROM Clause

- [x] `SELECT * FROM |` → table_ref (tables, views, schemas)
- [x] `SELECT * FROM dbo.|` → table_ref qualified with schema
- [x] `SELECT * FROM t1, |` → table_ref after comma
- [x] `SELECT * FROM (SELECT * FROM |)` → table_ref in derived table
- [x] `SELECT * FROM t |` → keyword candidates (WHERE, JOIN, LEFT, RIGHT, CROSS, INNER, FULL, OUTER, APPLY, ORDER, GROUP, HAVING, UNION, FOR, OPTION)
- [x] `SELECT * FROM t AS |` → no specific candidates (alias context)

### 3.3 JOIN Clauses

- [x] `SELECT * FROM t1 JOIN |` → table_ref after JOIN
- [x] `SELECT * FROM t1 LEFT JOIN |` → table_ref after LEFT JOIN
- [x] `SELECT * FROM t1 RIGHT JOIN |` → table_ref after RIGHT JOIN
- [x] `SELECT * FROM t1 CROSS JOIN |` → table_ref after CROSS JOIN
- [x] `SELECT * FROM t1 FULL OUTER JOIN |` → table_ref after FULL OUTER JOIN
- [x] `SELECT * FROM t1 CROSS APPLY |` → table_ref after CROSS APPLY
- [x] `SELECT * FROM t1 OUTER APPLY |` → table_ref after OUTER APPLY
- [x] `SELECT * FROM t1 JOIN t2 ON |` → columnref after ON
- [x] `SELECT * FROM t1 |` → JOIN keywords (JOIN, LEFT, RIGHT, INNER, CROSS, FULL, OUTER, APPLY)

### 3.4 WHERE, GROUP BY, HAVING

- [x] `SELECT * FROM t WHERE |` → columnref after WHERE
- [x] `SELECT * FROM t WHERE a = 1 AND |` → columnref after AND
- [x] `SELECT * FROM t WHERE a = 1 OR |` → columnref after OR
- [x] `SELECT * FROM t GROUP BY |` → columnref after GROUP BY
- [x] `SELECT * FROM t GROUP BY a, |` → columnref after comma
- [x] `SELECT * FROM t GROUP BY a |` → keyword candidates (HAVING, ORDER, FOR, OPTION)
- [x] `SELECT * FROM t HAVING |` → columnref after HAVING

### 3.5 ORDER BY, OFFSET-FETCH, TOP

- [x] `SELECT * FROM t ORDER BY |` → columnref after ORDER BY
- [x] `SELECT * FROM t ORDER BY a, |` → columnref after comma
- [x] `SELECT * FROM t ORDER BY a |` → keyword candidates (ASC, DESC, OFFSET, comma, FOR, OPTION)
- [x] `SELECT * FROM t ORDER BY a OFFSET |` → no specific (numeric context)
- [x] `SELECT * FROM t ORDER BY a OFFSET 10 ROWS FETCH |` → keyword candidates (NEXT, FIRST)

### 3.6 Set Operations & FOR Clause

- [x] `SELECT a FROM t UNION |` → keyword candidates (ALL, SELECT)
- [x] `SELECT a FROM t UNION ALL |` → keyword candidate (SELECT)
- [x] `SELECT a FROM t INTERSECT |` → keyword candidates (ALL, SELECT)
- [x] `SELECT a FROM t EXCEPT |` → keyword candidates (ALL, SELECT)
- [x] `SELECT * FROM t FOR |` → keyword candidates (XML, JSON, BROWSE)

### 3.7 CTE (WITH Clause)

- [x] `WITH |` → identifier context for CTE name
- [x] `WITH cte AS (|)` → keyword candidate (SELECT)
- [x] `WITH cte AS (SELECT * FROM t) SELECT |` → columnref (CTE columns available)
- [x] `WITH cte AS (SELECT * FROM t) SELECT * FROM |` → table_ref (CTE name available)
- [x] `WITH cte (|)` → identifier context for column names

### 3.8 Window Functions & Table Hints

- [x] `SELECT ROW_NUMBER() OVER (|)` → keyword candidates (PARTITION, ORDER)
- [x] `SELECT SUM(b) OVER (PARTITION BY |)` → columnref
- [x] `SELECT SUM(b) OVER (ORDER BY |)` → columnref
- [x] `SELECT SUM(b) OVER (ORDER BY a ROWS |)` → keyword candidates (BETWEEN, UNBOUNDED, CURRENT)
- [x] `SELECT * FROM t WITH (|)` → table hint keywords (NOLOCK, READUNCOMMITTED, READCOMMITTED, REPEATABLEREAD, SERIALIZABLE, HOLDLOCK, UPDLOCK, TABLOCK, TABLOCKX, ROWLOCK, PAGLOCK, INDEX, FORCESEEK)
- [x] `SELECT * FROM t WITH (NOLOCK, |)` → table hint keywords after comma

### 3.9 PIVOT/UNPIVOT, FOR XML/JSON, OPTION

- [x] `SELECT * FROM t PIVOT (|)` → aggregate function context
- [x] `SELECT * FROM t UNPIVOT (|)` → column context
- [x] `SELECT * FROM t FOR XML |` → keyword candidates (PATH, RAW, AUTO, EXPLICIT)
- [x] `SELECT * FROM t FOR JSON |` → keyword candidates (PATH, AUTO)
- [x] `SELECT * FROM t OPTION (|)` → query hint keywords (RECOMPILE, OPTIMIZE FOR, MAXDOP, HASH JOIN, MERGE JOIN, LOOP JOIN, FORCE ORDER)

---

## Phase 4: DML Instrumentation

### 4.1 INSERT

- [x] `INSERT INTO |` → table_ref after INTO
- [x] `INSERT INTO t (|)` → columnref for table t
- [x] `INSERT INTO t (a, |)` → columnref after comma
- [x] `INSERT INTO t VALUES (|)` → no specific candidates (value context)
- [x] `INSERT INTO t |` → keyword candidates (VALUES, SELECT, DEFAULT, EXEC, OUTPUT)
- [x] `INSERT INTO t SELECT |` → columnref (INSERT SELECT)
- [x] `INSERT INTO t OUTPUT |` → keyword candidates (inserted.*, deleted.*)

### 4.2 UPDATE

- [x] `UPDATE |` → table_ref
- [x] `UPDATE t SET |` → columnref for table t
- [x] `UPDATE t SET a = 1, |` → columnref after comma
- [x] `UPDATE t SET a = 1 WHERE |` → columnref
- [x] `UPDATE t SET a = 1 OUTPUT |` → keyword candidates (inserted.*, deleted.*)
- [x] `UPDATE t SET a = 1 FROM t JOIN |` → table_ref (UPDATE with JOIN)

### 4.3 DELETE

- [x] `DELETE FROM |` → table_ref
- [x] `DELETE FROM t WHERE |` → columnref for table t
- [x] `DELETE FROM t OUTPUT |` → keyword candidates (deleted.*)
- [x] `DELETE t FROM t JOIN |` → table_ref (DELETE with JOIN)

### 4.4 MERGE

- [x] `MERGE INTO |` → table_ref
- [x] `MERGE INTO t USING |` → table_ref for source
- [x] `MERGE INTO t USING s ON |` → columnref for join condition
- [x] `MERGE INTO t USING s ON t.id = s.id WHEN |` → keyword candidates (MATCHED, NOT)
- [x] `MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN |` → keyword candidates (UPDATE, DELETE)
- [x] `MERGE INTO t USING s ON t.id = s.id WHEN NOT MATCHED THEN |` → keyword candidate (INSERT)

---

## Phase 5: DDL Instrumentation — CREATE/ALTER/DROP

### 5.1 CREATE TABLE

- [x] `CREATE TABLE |` → identifier context (schema-qualified name)
- [x] `CREATE TABLE t (a INT, |)` → keyword candidates for column/constraint (PRIMARY, UNIQUE, INDEX, FOREIGN, CHECK, CONSTRAINT, DEFAULT)
- [x] `CREATE TABLE t (a INT |)` → keyword candidates for column options (NOT, NULL, DEFAULT, IDENTITY, PRIMARY, UNIQUE, REFERENCES, CHECK, CONSTRAINT, COLLATE)
- [x] `CREATE TABLE t (a |)` → type candidates (INT, BIGINT, VARCHAR, NVARCHAR, CHAR, NCHAR, TEXT, NTEXT, DATETIME, DATETIME2, DATE, TIME, DECIMAL, NUMERIC, FLOAT, REAL, BIT, MONEY, UNIQUEIDENTIFIER, XML, VARBINARY, IMAGE)
- [x] `CREATE TABLE t (a INT IDENTITY(|)` → no specific (numeric seed/increment)
- [x] `CREATE TABLE t (a INT REFERENCES |)` → table_ref
- [x] `CREATE TABLE t (a INT DEFAULT |)` → expression context
- [x] `CREATE TABLE t (a INT) |` → keyword candidates (ON, WITH, ;)
- [x] `CREATE TABLE t (CONSTRAINT pk PRIMARY KEY (|)` → columnref

### 5.2 ALTER TABLE

- [x] `ALTER TABLE |` → table_ref
- [x] `ALTER TABLE t |` → keyword candidates (ADD, DROP, ALTER, ENABLE, DISABLE, SWITCH, REBUILD, SET, WITH)
- [x] `ALTER TABLE t ADD |` → keyword candidates (COLUMN, CONSTRAINT, PRIMARY, UNIQUE, FOREIGN, CHECK, DEFAULT, INDEX)
- [x] `ALTER TABLE t ADD COLUMN |` → identifier context
- [x] `ALTER TABLE t DROP |` → keyword candidates (COLUMN, CONSTRAINT, INDEX)
- [x] `ALTER TABLE t DROP COLUMN |` → columnref for table t
- [x] `ALTER TABLE t DROP CONSTRAINT |` → constraint name context
- [x] `ALTER TABLE t ALTER COLUMN |` → columnref for table t
- [x] `ALTER TABLE t ALTER COLUMN a |` → type candidates
- [x] `ALTER TABLE t ENABLE TRIGGER |` → trigger_ref (ALL keyword also valid)
- [x] `ALTER TABLE t DISABLE TRIGGER |` → trigger_ref (ALL keyword also valid)

### 5.3 CREATE/DROP Index, View, Database

- [x] `CREATE INDEX idx ON |` → table_ref
- [x] `CREATE INDEX idx ON t (|)` → columnref for table t
- [x] `CREATE UNIQUE INDEX idx ON |` → table_ref
- [x] `CREATE CLUSTERED INDEX idx ON |` → table_ref
- [x] `CREATE NONCLUSTERED INDEX idx ON |` → table_ref
- [x] `DROP INDEX |` → index name context
- [x] `CREATE VIEW |` → identifier context
- [x] `CREATE VIEW v AS |` → keyword candidate (SELECT)
- [x] `ALTER VIEW v AS |` → keyword candidate (SELECT)
- [x] `DROP VIEW |` → view name context
- [x] `CREATE DATABASE |` → identifier context
- [x] `DROP DATABASE |` → database_ref
- [x] `DROP TABLE |` → table_ref
- [x] `DROP TABLE IF EXISTS |` → table_ref
- [x] `TRUNCATE TABLE |` → table_ref
- [x] `CREATE SCHEMA |` → identifier context

### 5.4 CREATE/ALTER PROCEDURE & FUNCTION

- [x] `CREATE PROCEDURE |` → identifier context
- [x] `CREATE PROC |` → identifier context
- [x] `CREATE PROCEDURE p |` → keyword candidates (AS, WITH)
- [x] `CREATE PROCEDURE p AS |` → statement keywords (SELECT, INSERT, etc.)
- [x] `ALTER PROCEDURE |` → procedure name context
- [x] `DROP PROCEDURE |` → procedure name context
- [x] `CREATE FUNCTION |` → identifier context
- [x] `CREATE FUNCTION f (|)` → keyword for param direction (@ variable)
- [x] `CREATE FUNCTION f () RETURNS |` → type candidates
- [x] `ALTER FUNCTION |` → function name context
- [x] `DROP FUNCTION |` → function name context

### 5.5 CREATE/DROP TRIGGER

- [x] `CREATE TRIGGER |` → identifier context
- [x] `CREATE TRIGGER tr ON |` → table_ref
- [x] `CREATE TRIGGER tr ON t |` → keyword candidates (FOR, AFTER, INSTEAD OF)
- [x] `CREATE TRIGGER tr ON t FOR |` → keyword candidates (INSERT, UPDATE, DELETE)
- [x] `DROP TRIGGER |` → trigger name context
- [x] `ALTER TRIGGER |` → trigger name context

---

## Phase 6: Control Flow & Declarations

### 6.1 DECLARE & SET

- [x] `DECLARE |` → keyword candidate (@) for variable declaration
- [x] `DECLARE @v |` → type candidates
- [x] `SET |` → keyword candidates (@variable, ANSI_NULLS, QUOTED_IDENTIFIER, NOCOUNT, XACT_ABORT, etc.)
- [x] `SET @v = |` → expression context (columnref, func_name)

### 6.2 Control Flow

- [x] `IF |` → expression context (columnref, func_name)
- [x] `WHILE |` → expression context
- [x] `BEGIN |` → keyword candidates (TRANSACTION, TRY, statement keywords)
- [x] `BEGIN TRY |` → statement keywords
- [x] `BEGIN CATCH |` → statement keywords
- [x] `RETURN |` → expression context (optional value)
- [x] `EXEC |` → proc_ref (procedure names)
- [x] `EXECUTE |` → proc_ref
- [x] `EXEC p @param = |` → expression context
- [x] `WAITFOR |` → keyword candidates (DELAY, TIME)

### 6.4 Cursor Operations

- [x] `DECLARE c CURSOR FOR |` → keyword candidate (SELECT)
- [x] `OPEN |` → cursor name context
- [x] `FETCH NEXT FROM |` → cursor name context
- [x] `CLOSE |` → cursor name context
- [x] `DEALLOCATE |` → cursor name context

### 6.3 Transaction

- [x] `BEGIN TRANSACTION |` → identifier context (optional name)
- [x] `COMMIT |` → keyword candidates (TRANSACTION, WORK)
- [x] `ROLLBACK |` → keyword candidates (TRANSACTION, WORK)
- [x] `SAVE TRANSACTION |` → identifier context (savepoint name)

---

## Phase 7: Security & Utility

### 7.1 GRANT, REVOKE, DENY

- [x] `GRANT |` → keyword candidates (SELECT, INSERT, UPDATE, DELETE, EXECUTE, ALTER, CONTROL, VIEW, REFERENCES, ALL)
- [x] `GRANT SELECT ON |` → table_ref (object name)
- [x] `GRANT SELECT ON t TO |` → user_ref / role_ref
- [x] `REVOKE |` → keyword candidates (same as GRANT)
- [x] `DENY |` → keyword candidates (same as GRANT)

### 7.2 USE & Utility

- [x] `USE |` → database_ref
- [x] `PRINT |` → expression context
- [x] `RAISERROR(|)` → expression context
- [x] `THROW |` → expression context
- [x] `BACKUP DATABASE |` → database_ref
- [x] `RESTORE DATABASE |` → database_ref

---

## Phase 8: Expression Context Instrumentation

### 8.1 Expression Positions

- [x] After `=` operator → columnref, func_name, literal keywords
- [x] After `AND` / `OR` → columnref, func_name
- [x] After `(` in expression → columnref, func_name, subquery keywords
- [x] After `CASE WHEN |` → columnref, func_name
- [x] After `CASE WHEN x THEN |` → columnref, func_name
- [x] After `CASE WHEN x THEN y ELSE |` → columnref, func_name
- [x] After `IN (|)` → columnref, func_name, SELECT
- [x] After `BETWEEN |` → columnref, func_name
- [x] After `LIKE |` → expression context
- [x] After `CAST(|` → columnref, func_name
- [x] After `CAST(x AS |)` → type candidates
- [x] After `CONVERT(|` → type candidates
- [x] After `COALESCE(|` → columnref, func_name
- [x] After `NULLIF(|` → columnref, func_name
- [x] After `IIF(|` → expression context (condition)
- [x] After `WHERE a IS |` → keyword candidates (NULL, NOT NULL)
- [x] After `WHERE EXISTS (|)` → keyword candidate (SELECT)
- [x] After `BETWEEN 1 AND |` → columnref, func_name
- [x] `CASE a WHEN |` → expression context (simple CASE value matching)
- [x] After `TRY_CAST(x AS |)` → type candidates
- [x] After `TRY_CONVERT(|` → type candidates

---

## Phase 9: Integration & Edge Cases

### 9.1 Qualified Name Resolution

- [x] `SELECT t1.| FROM t1 JOIN t2` → columnref qualified to t1
- [x] `SELECT x.| FROM t AS x` → columnref with alias resolution
- [x] `SELECT dbo.t.| FROM dbo.t` → columnref with schema-qualified table
- [x] `SELECT * FROM dbo.|` → table_ref within schema
- [x] `SELECT |.* FROM Employees AS e` → columnref signal for Bytebase table/alias star qualifier completion
- [x] `SELECT e.|* FROM Employees AS e` → columnref signal for qualified star completion
- [x] CTE columns available: `WITH cte(a,b) AS (...) SELECT cte.| FROM cte` → columnref from CTE

### 9.2 Edge Cases

- [x] Cursor at position 0 in non-empty SQL → top-level keywords
- [x] Cursor mid-identifier: `SEL|ECT` → prefix "SEL" matches SELECT
- [x] After semicolon in multi-statement: `SELECT 1; |` → top-level keywords
- [x] Empty input `""` → top-level keywords
- [x] Whitespace only `"   "` → top-level keywords
- [x] Syntax error before cursor: `SELECT FROM WHERE |` → best-effort candidates
- [x] Bracket-quoted identifiers: `SELECT [col| FROM t` → column completion with bracket context
- [x] Very long SQL (10K+ chars) → no panic, returns candidates

### 9.3 Complex SQL Patterns

- [x] Nested subqueries: `SELECT * FROM (SELECT * FROM (SELECT | FROM t))` → columnref at inner level
- [x] Correlated subquery: `SELECT * FROM t1 WHERE EXISTS (SELECT * FROM t2 WHERE t2.a = t1.|)` → columnref from t1
- [x] Multi-JOIN: `SELECT | FROM t1 JOIN t2 ON ... JOIN t3 ON ...` → columnref from all three tables
- [x] UNION column scoping: `SELECT a FROM t1 UNION SELECT | FROM t2` → columnref from t2
- [x] INSERT...SELECT: `INSERT INTO t1 SELECT | FROM t2` → columnref from t2
- [x] Complex ALTER: `ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (|)` → columnref for t
