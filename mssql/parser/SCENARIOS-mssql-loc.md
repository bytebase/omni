# MSSQL Parser Position Tracking (Loc) Scenarios

> Goal: Bring MSSQL parser position tracking to full parity with PG parser — accurate Token.End, correct Loc.End on all AST nodes, public API with line:column, sentinel value alignment
> Verification: `go build ./mssql/... && go test ./mssql/... -count=1`; Loc accuracy tests verify `sql[node.Loc.Start:node.Loc.End]` matches expected text
> Reference sources: PG parser (`pg/parser/lexer.go` Token.End, `pg/parser/parser.go` p.prev.End, `pg/parse.go` Statement/Position/lineIndex)

Status: [ ] pending, [x] passing, [~] partial (needs upstream change)

---

## Phase 1: Foundation

Lexer and parser infrastructure changes that all subsequent phases depend on.

### 1.1 Token End Position

- [x] Token struct has `End int` field (exclusive byte offset past token end)
- [x] Lexer sets `Token.End` to lexer position after consuming each token
- [x] EOF token has `End == Loc` (zero-length)
- [x] Single-char tokens (`;`, `(`, `)`) have `End == Loc + 1`
- [x] Multi-char operators (`<=`, `>=`, `<>`, `!=`, `+=`, `-=`, `*=`, `/=`, `%=`, `&=`, `|=`, `^=`, `!<`, `!>`, `::`) have correct End
- [x] Integer literals have correct End (`123` → End = Loc + 3)
- [x] Hex literals have correct End (`0xFF` → End = Loc + 4)
- [x] Scientific notation literals have correct End (`1.5e10` → End = Loc + 6)
- [x] String literals include quotes in range (`'hello'` → End = Loc + 7)
- [x] Escaped quotes in strings have correct End (`'it''s'` → End = Loc + 7)
- [x] N-string literals include prefix (`N'hello'` → End = Loc + 8)
- [x] Identifiers have correct End (`myTable` → End = Loc + 7)
- [x] Quoted identifiers include brackets (`[my col]` → End = Loc + 8)
- [x] Double-quoted identifiers include quotes (`"my col"` → End = Loc + 8)
- [x] Escaped brackets in identifiers have correct End (`[col]]name]` → correct End)
- [x] Keywords have correct End (`SELECT` → End = Loc + 6)
- [x] Variables include @ prefix (`@var` → End = Loc + 4)
- [x] System variables include @@ prefix (`@@ROWCOUNT` → End = Loc + 10)
- [x] Float literals have correct End (`3.14` → End = Loc + 4)
- [x] Whitespace/comments are skipped — next token Loc starts after them

### 1.2 Sentinel Value Alignment

- [x] Loc struct uses -1 for "unknown" (not 0)
- [x] `NoLoc()` helper returns `Loc{Start: -1, End: -1}`
- [x] Existing AST node Loc defaults updated (0 → -1 where "not set")
- [~] Nodes that never get explicit Loc assignment use Go zero-value `Loc{}` which equals `NoLoc()` (-1, -1) — impossible: Go zero-value is `{0,0}`, not `{-1,-1}`; use explicit `NoLoc()` where needed
- [x] All existing tests pass with new sentinel value

### 1.3 Parser Prev.End Infrastructure

- [ ] Parser `advance()` preserves full Token (including End) in `p.prev`
- [ ] New `p.prevEnd()` helper returns `p.prev.End`
- [ ] `p.pos()` still returns `p.cur.Loc` (start of current token, unchanged)
- [ ] Build passes with infrastructure in place

---

## Phase 2: Loc.End Accuracy — DML & Expressions

Migrate `.End = p.pos()` → `.End = p.prevEnd()` across DML and expression files. Each scenario = one file fully migrated, build passes.

### 2.1 Core Expressions

- [ ] expr.go — 7 sites migrated to `p.prevEnd()`
- [ ] name.go — 3 sites migrated
- [ ] type.go — 1 site migrated
- [ ] rowset_functions.go — Loc.End sites migrated (if any `.End = p.pos()` exist)
- [ ] Build passes after 2.1

### 2.2 SELECT & Related

- [ ] select.go — 25 sites migrated to `p.prevEnd()`
- [ ] cursor.go — 9 sites migrated
- [ ] Build passes after 2.2

### 2.3 DML: INSERT, UPDATE, DELETE, MERGE

- [ ] insert.go — 3 sites migrated
- [ ] update_delete.go — 4 sites migrated
- [ ] merge.go — 3 sites migrated
- [ ] Build passes after 2.3

### 2.4 Control Flow & Declarations

- [ ] control_flow.go — 6 sites migrated
- [ ] declare_set.go — 9 sites migrated
- [ ] execute.go — 5 sites migrated
- [ ] transaction.go — 5 sites migrated
- [ ] Build passes after 2.4

---

## Phase 3: Loc.End Accuracy — DDL

### 3.1 CREATE: Tables, Indexes, Views

- [ ] create_table.go — 21 sites migrated to `p.prevEnd()`
- [ ] create_index.go — 8 sites migrated
- [ ] create_view.go — 3 sites migrated
- [ ] Build passes after 3.1

### 3.2 CREATE: Procs, Triggers, Types, Sequences

- [ ] create_proc.go — 5 sites migrated
- [ ] create_trigger.go — 2 sites migrated
- [ ] create_type.go — 3 sites migrated
- [ ] create_sequence.go — 2 sites migrated
- [ ] create_statistics.go — 3 sites migrated
- [ ] create_synonym.go — 1 site migrated
- [ ] Build passes after 3.2

### 3.3 CREATE: Database, Schema, Other

- [ ] create_database.go — 7 sites migrated
- [ ] create_schema.go — 2 sites migrated
- [ ] assembly.go — 2 sites migrated
- [ ] partition.go — 4 sites migrated
- [ ] xml_schema.go — 2 sites migrated
- [ ] Build passes after 3.3

### 3.4 ALTER

- [ ] alter_table.go — 20 sites migrated
- [ ] alter_objects.go — 2 sites migrated
- [ ] Build passes after 3.4

### 3.5 DROP & Utility

- [ ] drop.go — 1 site migrated
- [ ] utility.go — 36 sites migrated
- [ ] dbcc.go — 1 site migrated
- [ ] bulk_insert.go — 1 site migrated
- [ ] go_batch.go — 1 site migrated
- [ ] Build passes after 3.5

---

## Phase 4: Loc.End Accuracy — Security & Server

### 4.1 Security

- [ ] grant.go — 3 sites migrated
- [ ] security_principals.go — 11 sites migrated
- [ ] security_misc.go — 9 sites migrated
- [ ] security_keys.go — 13 sites migrated
- [ ] security_audit.go — 11 sites migrated
- [ ] Build passes after 4.1

### 4.2 Server & Admin

- [ ] server.go — 7 sites migrated
- [ ] backup_restore.go — 10 sites migrated
- [ ] service_broker.go — 23 sites migrated
- [ ] resource_governor.go — 13 sites migrated
- [ ] Build passes after 4.2

### 4.3 Network & External

- [ ] endpoint.go — 3 sites migrated
- [ ] availability.go — 3 sites migrated
- [ ] external.go — 16 sites migrated
- [ ] event.go — 6 sites migrated
- [ ] fulltext.go — 10 sites migrated
- [ ] Build passes after 4.3

### 4.4 Top-Level Dispatch

- [ ] parser.go — 3 sites migrated (top-level parse loop Loc.End)
- [ ] All 348 `.End = p.pos()` sites across 48 files now use `p.prevEnd()`
- [ ] Full test suite passes
- [ ] Build passes after 4.4

---

## Phase 5: Public API Wrapper

### 5.1 Statement & Position Types

- [ ] `mssql/parse.go` exists with `Statement` struct (Text, AST, ByteStart, ByteEnd, Start, End)
- [ ] `Position` struct with `Line int` (1-based) and `Column int` (1-based, byte offset)
- [ ] `Statement.Text` contains the SQL text for one statement (including trailing semicolon)
- [ ] `Statement.AST` contains the parsed AST node
- [ ] `Statement.ByteStart` is inclusive start byte offset in original SQL
- [ ] `Statement.ByteEnd` is exclusive end byte offset in original SQL

### 5.2 Line Index Infrastructure

- [ ] `buildLineIndex(sql)` returns slice of byte offsets where each line starts
- [ ] Line index handles empty string (single entry: 0)
- [ ] Line index handles no newlines (single entry: 0)
- [ ] Line index handles multiple newlines (`\n` at various positions)
- [ ] Line index handles trailing newline
- [ ] `offsetToPosition(idx, offset)` returns correct Position via binary search
- [ ] Offset 0 → `Position{Line: 1, Column: 1}`
- [ ] Offset at start of line N → `Position{Line: N, Column: 1}`
- [ ] Offset mid-line → correct column calculation
- [ ] Offset past end of input → reasonable behavior (last line, past-end column)

### 5.3 Parse() Function

- [ ] `Parse(sql string) ([]Statement, error)` is the public entry point
- [ ] Calls `parser.Parse()` internally
- [ ] Builds line index once for entire SQL input
- [ ] Each statement gets correct ByteStart/ByteEnd from AST Loc
- [ ] Each statement gets correct Start/End Position (line:column)
- [ ] Multi-statement SQL splits correctly (`SELECT 1; SELECT 2`)
- [ ] Statement.Text includes trailing semicolon when present
- [ ] Empty input returns empty slice, no error
- [ ] Single statement without semicolon works
- [ ] Error from parser is propagated correctly

---

## Phase 6: Loc Accuracy Verification Tests

### 6.1 DML Loc Verification

- [ ] `SELECT col FROM t` — SelectStmt Loc spans full statement
- [ ] `SELECT a, b FROM t WHERE x = 1 ORDER BY a` — SelectStmt Loc spans full statement
- [ ] `INSERT INTO t (a) VALUES (1)` — InsertStmt Loc spans full statement
- [ ] `UPDATE t SET a = 1 WHERE b = 2` — UpdateStmt Loc spans full statement
- [ ] `DELETE FROM t WHERE a = 1` — DeleteStmt Loc spans full statement
- [ ] `MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET a = s.a` — MergeStmt Loc spans full statement
- [ ] `WITH cte AS (SELECT 1) SELECT * FROM cte` — top-level Loc spans full CTE+SELECT
- [ ] `SELECT * FROM t1 JOIN t2 ON t1.id = t2.id` — JoinExpr Loc spans join clause

### 6.2 DDL Loc Verification

- [ ] `CREATE TABLE t (a INT, b VARCHAR(50))` — CreateTableStmt Loc spans full statement
- [ ] `CREATE INDEX ix ON t (a)` — CreateIndexStmt Loc spans full statement
- [ ] `CREATE VIEW v AS SELECT 1` — CreateViewStmt Loc spans full statement
- [ ] `CREATE PROCEDURE p AS SELECT 1` — CreateProcedureStmt Loc spans full statement
- [ ] `ALTER TABLE t ADD c INT` — AlterTableStmt Loc spans full statement
- [ ] `DROP TABLE t` — DropStmt Loc spans full statement
- [ ] `CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN 1 END` — CreateFunctionStmt Loc spans full statement
- [ ] `CREATE TRIGGER tr ON t FOR INSERT AS SELECT 1` — CreateTriggerStmt Loc spans full statement
- [ ] `TRUNCATE TABLE t` — TruncateStmt Loc spans full statement

### 6.3 Expression & Sub-node Loc Verification

- [ ] `SELECT CAST(1 AS INT)` — CastExpr Loc spans `CAST(1 AS INT)`
- [ ] `SELECT CASE WHEN 1=1 THEN 'a' ELSE 'b' END` — CaseExpr Loc spans CASE..END
- [ ] `SELECT a + b` — BinaryExpr Loc spans `a + b`
- [ ] `SELECT -x` — UnaryExpr Loc spans `-x`
- [ ] `SELECT COALESCE(a, b, c)` — FuncCall Loc spans full call
- [ ] `SELECT * FROM t WHERE EXISTS (SELECT 1)` — ExistsExpr Loc spans `EXISTS (...)`
- [ ] `SELECT CONVERT(INT, '1')` — ConvertExpr Loc spans `CONVERT(INT, '1')`
- [ ] `SELECT TRY_CAST(1 AS VARCHAR)` — TryCastExpr Loc spans full expression
- [ ] `SELECT TRY_CONVERT(INT, '1')` — TryConvertExpr Loc spans full expression

### 6.4 Multi-Statement & Edge Cases

- [ ] `SELECT 1; SELECT 2` — each statement has non-overlapping Loc ranges
- [ ] Multi-line SQL — Loc byte offsets account for newlines correctly
- [ ] Statement with leading whitespace — Loc.Start points to first keyword, not whitespace
- [ ] Statement with trailing semicolon — Loc.End points past the semicolon or past the last token (consistent choice)
- [ ] Empty statement (bare semicolons `;;`) — handled without panic
- [ ] Very long single-line SQL — positions remain accurate
- [ ] Statement ending with block comment — `SELECT 1 /* comment */; SELECT 2` — Loc.End behavior is consistent
- [ ] GO batch separator — Loc tracking works correctly across GO-separated batches
- [ ] Error-path Loc — nodes parsed before a syntax error have valid Loc; incomplete nodes get NoLoc

### 6.5 Public API Line:Column Tests

- [ ] Single-line `SELECT 1` → Start = {1,1}, End past statement
- [ ] Two-line `SELECT\n1` → correct line breaks
- [ ] Multi-statement multi-line → each statement has correct Start/End positions
- [ ] Tab characters count as 1 column (byte-based, not visual)
- [ ] Unicode content — column is byte offset, not character count
- [ ] `Parse("")` returns empty slice with no error
- [ ] `Parse("SELECT 1")` returns one Statement with correct positions
- [ ] `Parse("SELECT 1; SELECT 2; SELECT 3")` returns three Statements
