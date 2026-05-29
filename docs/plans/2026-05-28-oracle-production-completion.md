# Oracle Production Completion Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build Oracle parser-native completion in omni to the point that Bytebase can replace its production Oracle ANTLR/C3 completer without losing SQL editor behavior.

**Architecture:** omni owns parser-native completion signals: token candidates, rule candidates, prefix handling, visible range scope, CTEs, and qualified-name intent. Bytebase remains responsible for metadata-backed resolution to `base.Candidate` values because Oracle schemas are represented there as database names with empty internal schema metadata. Implementation is test-first: each parser signal is introduced by a failing `oracle/parser` test before production code.

**Tech Stack:** Go, recursive-descent Oracle parser under `oracle/parser`, Bytebase adapter target at `/Users/rebeliceyang/Github/bytebase/backend/plugin/parser/plsql`, existing oracle parser coverage gates.

## Task 1: Parser Completion API Foundation

**Files:**
- Create: `oracle/parser/complete.go`
- Create: `oracle/parser/complete_test.go`

**Steps:**
1. Write tests for empty input, semicolon boundary, `Tokenize`, `TokenName`, and prefix retry.
2. Run `go test ./oracle/parser -run 'TestOracleCompletionAPI' -count=1` and verify the tests fail because APIs are missing.
3. Implement `Tokenize`, `TokenName`, `CandidateSet`, `Collect`, prefix extraction, and statement-starter inference.
4. Run the same test and verify it passes.
5. Run `go test ./oracle/parser -run 'TestOracleParserProgress|TestOracleCoverage|TestVerifyCorpus|TestOracleKeywordManifestExhaustive|TestOracleLocNodeCoverage|TestOracleCompletionAPI' -count=1`.

## Task 2: SELECT Scope And Intent

**Files:**
- Create: `oracle/parser/complete_context.go`
- Create: `oracle/parser/complete_context_test.go`
- Modify: `oracle/parser/complete.go`

**Steps:**
1. Write tests for `SELECT * FROM |`, `SELECT | FROM t`, `SELECT a.| FROM t a`, joins, WHERE, GROUP BY, ORDER BY, and schema-qualified table refs.
2. Run `go test ./oracle/parser -run 'TestOracleCompletionSelect' -count=1` and verify failure.
3. Implement `CollectCompletion`, `CompletionContext`, `ScopeSnapshot`, `RangeReference`, `CompletionIntent`, and token-stream scope extraction.
4. Run the SELECT completion tests and verify pass.
5. Run parser coverage gates.

## Task 3: CTE And Subquery Scope

**Files:**
- Modify: `oracle/parser/complete_context.go`
- Modify: `oracle/parser/complete_context_test.go`

**Steps:**
1. Write tests for simple CTE, explicit CTE columns, subquery alias columns, and qualified CTE references.
2. Run targeted tests and verify failure.
3. Implement best-effort CTE extraction and subquery range references.
4. Run targeted tests and parser coverage gates.

## Task 4: DML Rule Signals

**Files:**
- Modify: `oracle/parser/complete.go`
- Modify: `oracle/parser/complete_test.go`
- Modify: `oracle/parser/complete_context.go`

**Steps:**
1. Write tests for INSERT, UPDATE, DELETE, and MERGE table/column/expression positions.
2. Verify failure.
3. Add DML token-pattern completion signals and DML target/source scope extraction.
4. Run targeted tests and parser coverage gates.

## Task 5: DDL And Utility Rule Signals

**Files:**
- Modify: `oracle/parser/complete.go`
- Modify: `oracle/parser/complete_test.go`

**Steps:**
1. Write tests for CREATE/ALTER/DROP/TRUNCATE/COMMENT/GRANT/REVOKE contexts.
2. Verify failure.
3. Add DDL/utility rule and keyword candidate inference.
4. Run targeted tests and parser coverage gates.

## Task 6: Oracle Keyword-Driven Candidates

**Files:**
- Modify: `oracle/parser/complete.go`
- Modify: `oracle/parser/complete_test.go`

**Steps:**
1. Write tests for datatype, function-like keyword, pseudo-column, and reserved identifier policies.
2. Verify failure.
3. Reuse Oracle keyword classification helpers and manifest policies for candidate generation.
4. Run targeted tests and parser coverage gates.

## Task 7: Bytebase Adapter Cutover

**Files:**
- Modify: `/Users/rebeliceyang/Github/bytebase/backend/plugin/parser/plsql/completion.go`
- Modify: `/Users/rebeliceyang/Github/bytebase/backend/plugin/parser/plsql/completion_test.go`

**Steps:**
1. Add a failing Bytebase test proving Oracle completion no longer depends on ANTLR/C3 imports.
2. Run `go test ./backend/plugin/parser/plsql -run 'TestCompletionDoesNotDependOnANTLR|TestCompletion' -count=1` from the Bytebase repo and verify failure.
3. Replace the ANTLR completion collector with omni Oracle `CollectCompletion`, preserving metadata resolution and quoting behavior.
4. Run Bytebase Oracle completion tests.
5. Run omni Oracle parser gates again.

## Task 8: Final Verification And Status

**Files:**
- Modify: `SCENARIOS-oracle-completion.md`
- Modify: `docs/PARSER-DEFENSE-MATRIX.md`
- Modify: `docs/engine-capability-guide.md`

**Steps:**
1. Mark scenario rows according to passing tests.
2. Update status docs only for proven behavior.
3. Run full verification:
   - `go test ./oracle/parser -count=1`
   - Bytebase: `go test ./backend/plugin/parser/plsql -count=1`
4. Summarize residual defects and deferred items.
