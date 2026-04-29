# Oracle Parser Foundation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Bring Oracle parser foundation, soft-fail/strictness, and keyword classification guardrails closer to MySQL/PG/MSSQL parser quality.

**Architecture:** Work in three ordered passes: foundation error/location plumbing, behavior-level soft-fail and strictness checks, then keyword classification golden tests. Each pass starts with failing parser tests, then uses the smallest parser changes that make those tests pass without changing catalog behavior.

**Tech Stack:** Go recursive-descent parser in `oracle/parser`, AST nodes in `oracle/ast`, tests under `oracle/parser`.

### Task 1: Foundation Errors And Locations

**Files:**
- Modify: `oracle/parser/parser.go`
- Modify: `oracle/ast/node.go`
- Modify: `oracle/parser/*_test.go`

**Steps:**
1. Add failing tests for PG-style syntax errors, lexer errors, strict statement separators, RawStmt `Loc.End`, and `NoLoc` sentinel behavior.
2. Run the focused tests and confirm they fail for current behavior.
3. Implement `syntaxErrorAtCur`, `syntaxErrorAtTok`, `lexerError`, statement separator enforcement, lexer error propagation, and `p.prev.End` RawStmt ranges.
4. Normalize `Loc` docs/checks around `-1` as the only unknown sentinel.
5. Run focused tests, then `go test ./oracle/parser`.

### Task 2: Soft-Fail And Strictness

**Files:**
- Modify: `oracle/parser/expr.go`
- Modify: `oracle/parser/select.go`
- Modify: `oracle/parser/create_table.go`
- Modify: `oracle/parser/alter_table.go`
- Modify: `oracle/parser/*_test.go`

**Steps:**
1. Add failing tests for truncated expressions, missing operands, duplicate clauses, unknown ALTER/option leads, and reserved keyword misuse in identifier positions.
2. Run focused tests and confirm they fail.
3. Add parser-level error recording/propagation where current helper signatures cannot yet return errors directly.
4. Replace silent nil/placeholder parses on required grammar positions with syntax errors.
5. Run focused tests, corpus verifier, then `go test ./oracle/parser`.

### Task 3: Keyword Classification Golden Tests

**Files:**
- Create: `oracle/parser/keyword_classification_test.go`
- Modify: `oracle/parser/lexer.go`
- Modify: `oracle/parser/name.go`

**Steps:**
1. Add failing golden tests for Oracle reserved/context/nonreserved keyword behavior in table, column, alias, dotted-name, and quoted identifier positions.
2. Run tests and confirm current gaps.
3. Introduce keyword category metadata and identifier helpers aligned with Oracle behavior.
4. Remove keyword string-matching gaps where practical.
5. Run keyword tests and full parser tests.
