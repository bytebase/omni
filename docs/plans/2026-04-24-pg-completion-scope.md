# PostgreSQL Completion Scope Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a PostgreSQL completion-only API that returns parser candidates plus the visible FROM scope at the cursor.

**Architecture:** Keep ordinary `Parse` strict and avoid exposing partial ASTs. Add `parser.CollectCompletion(sql, cursorOffset)` as a completion context API that wraps existing `Collect` candidates and computes a best-effort `ScopeSnapshot` from grammar/token facts and parseable prefixes. ByteBase remains responsible for metadata and final candidate formatting.

**Tech Stack:** Go, `github.com/bytebase/omni/pg/parser`, `github.com/bytebase/omni/pg/ast`.

### Task 1: Scope API Tests

**Files:**
- Create: `pg/parser/complete_context_test.go`

**Steps:**
- Write failing tests for `CollectCompletion`.
- Cover incomplete `JOIN ... USING (`, aliases with alias columns, incomplete `JOIN`, subquery aliases, and CTE references.
- Run `go test ./pg/parser -run 'TestCollectCompletion' -count=1`; expected failure is missing API/types.

### Task 2: Minimal API and Scope Extraction

**Files:**
- Create: `pg/parser/complete_context.go`

**Steps:**
- Define `CompletionContext`, `ScopeSnapshot`, `RangeReference`, and `RangeReferenceKind`.
- Implement `CollectCompletion` by reusing `Collect` and extracting visible references for the cursor's SELECT scope.
- Preserve strict `Parse`; completion-only recovery may tolerate incomplete join qualifiers.
- Run focused tests until green.

### Task 3: Public Package Wrapper

**Files:**
- Create: `pg/completion.go`

**Steps:**
- Add a small `pg.CollectCompletion` wrapper if downstream users prefer the top-level pg package.
- Keep the parser package API as the primary implementation.
- Run `go test ./pg ./pg/parser -run 'TestCollectCompletion|TestParse' -count=1`.

### Task 4: Verification

**Steps:**
- Run `go test ./pg/parser ./pg/completion`.
- Confirm ordinary strict parse behavior remains covered by existing parser tests.
