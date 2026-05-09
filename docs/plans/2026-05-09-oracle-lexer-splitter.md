# Oracle Lexer Splitter Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a production-ready Oracle lexer-based SQL script splitter that handles SQL, PL/SQL blocks, stored units, and SQL*Plus-style slash delimiters without depending on full parsing.

**Architecture:** Implement `oracle/parser.Split(sql string) []Segment` as a soft-fail lexical scanner over Oracle tokens and raw byte offsets. It should reuse Oracle lexical rules by scanning token-like boundaries, preserve source ranges, and model script delimiters and PL/SQL block boundaries with a small state machine.

**Tech Stack:** Go, existing `oracle/parser` lexer/token constants, existing parser package test patterns, `go test ./oracle/...`.

### Task 1: Baseline Splitter API and Ordinary SQL

**Files:**
- Create: `oracle/parser/split.go`
- Create: `oracle/parser/split_test.go`

**Steps:**
1. Write failing tests for empty input, ordinary semicolon splitting, comments/whitespace-only filtering, and semicolons inside Oracle literals.
2. Run `go test ./oracle/parser -run 'TestSplit' -count=1` and verify failure is due to missing `Split`.
3. Implement `Segment`, `Empty`, and lexical splitting for ordinary SQL.
4. Re-run the focused tests until they pass.

### Task 2: SQL*Plus Slash Delimiter

**Files:**
- Modify: `oracle/parser/split.go`
- Modify: `oracle/parser/split_test.go`

**Steps:**
1. Add failing tests for slash on its own line, slash after whitespace, slash at EOF, division expressions, and slash inside comments/literals.
2. Implement line-alone slash detection using raw byte offsets.
3. Re-run `go test ./oracle/parser -run 'TestSplit' -count=1`.

### Task 3: PL/SQL Blocks and Stored Units

**Files:**
- Modify: `oracle/parser/split.go`
- Modify: `oracle/parser/split_test.go`

**Steps:**
1. Add failing tests for anonymous `BEGIN...END`, `DECLARE...BEGIN...END`, nested blocks, procedure, function, package spec/body, trigger, type body, and follow-up SQL statements.
2. Implement statement-mode tracking for PL/SQL starts and outer `END [name];`.
3. Re-run focused tests and fix edge cases.

### Task 4: Public Oracle Parse Facade

**Files:**
- Create: `oracle/parse.go`
- Create: `oracle/parse_test.go`

**Steps:**
1. Add failing tests for `oracle.Parse` returning statements with text, AST, byte ranges, and line/column positions.
2. Implement facade using `parser.Split` and per-segment `parser.Parse`.
3. Re-run `go test ./oracle/... -count=1`.

### Task 5: Verification

**Commands:**
- `go test ./oracle/parser -run 'TestSplit' -count=1`
- `go test ./oracle/... -count=1`
- If parser package failures are unrelated existing failures, capture exact failing tests and reason in the final report.

### Task 6: SQL*Plus Script Commands

**Files:**
- Modify: `oracle/parser/split.go`
- Modify: `oracle/parser/split_test.go`
- Modify: `oracle/parse_test.go` if skipped command ranges affect facade behavior

**Step 1: Write failing tests**

Cover line-oriented SQL*Plus commands that must not be sent to the database parser:

- Environment/output commands: `SET`, `SHOW`, `PROMPT`, `SPOOL`, `COLUMN`, `TTITLE`, `BTITLE`, `BREAK`, `COMPUTE`, `REPHEADER`, `REPFOOTER`, `CLEAR`, `PAUSE`, `HELP`.
- Variable/input commands: `DEFINE`, `UNDEFINE`, `ACCEPT`, `VARIABLE`, `PRINT`, `EXECUTE`.
- Script/session commands: `@`, `@@`, `START`, `CONNECT`, `DISCONNECT`, `PASSWORD`, `EXIT`, `QUIT`, `WHENEVER SQLERROR`, `WHENEVER OSERROR`, `HOST`, `!`.
- Comment commands: `REM`, `REMARK`.
- Buffer commands that execute or manipulate SQL buffer: `/`, `RUN`, `LIST`, `APPEND`, `CHANGE`, `DEL`, `INPUT`, `SAVE`, `GET`, `EDIT`.

Expected behavior: these commands are line-local script commands and are skipped from returned SQL segments. `/` and `RUN` also flush the current segment as executable SQL/PLSQL before skipping the command line.

**Step 2: Verify red**

Run `go test ./oracle/parser -run 'TestSplitSQLPlusCommands' -count=1` and confirm the failures show command text leaking into SQL segments.

**Step 3: Implement line-command scanner**

Add raw-line command detection before token processing. When a command line starts at the current statement boundary, skip it and advance `stmtStart`. When a flushing command is found after accumulated SQL, append the segment first, then skip the command. Keep quoted SQL contents untouched by only scanning physical lines outside tokenized strings/comments.

**Step 4: Verify green**

Run the focused splitter tests, then `go test ./oracle/... -count=1`.
