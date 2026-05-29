# Oracle Production Completion Scenarios

Goal: Oracle completion is production-ready for Bytebase when omni exposes parser-native completion signals that let Bytebase replace its current ANTLR/C3 `backend/plugin/parser/plsql/completion.go` implementation without losing SQL editor behavior.

Status legend: `[ ]` pending, `[x]` passing, `[~]` partial/deferred.

## Phase 0: Public Parser Completion API

- [x] Empty input at cursor 0 returns statement starter candidates.
- [x] Whitespace-only input returns statement starter candidates.
- [x] Cursor after semicolon returns statement starter candidates for a new statement.
- [x] Cursor inside a partially typed keyword backs up to token start and prefix-filters candidates.
- [x] `Tokenize` returns non-EOF tokens with stable byte `Loc` and `End`.
- [x] `TokenName` returns uppercase Oracle keyword text for keyword token types.
- [x] `Collect` never panics on incomplete SQL.
- [x] `CollectCompletion` returns candidates, prefix, scope, CTEs, and object intent.

## Phase 1: Bytebase SELECT Object Completion

- [x] `SELECT * FROM |` emits table-reference intent.
- [x] `SELECT * FROM schema.|` emits table-reference intent qualified by schema.
- [x] `SELECT | FROM t` emits column-reference intent and table scope for `t`.
- [x] `SELECT t.| FROM t` emits column-reference intent qualified by table `t`.
- [x] `SELECT alias.| FROM t alias` resolves `alias` as a visible range reference.
- [x] `SELECT * FROM t JOIN |` emits table-reference intent.
- [x] `SELECT * FROM t JOIN u ON |` emits column-reference intent with both tables visible.
- [x] `SELECT * FROM t WHERE |` emits column-reference intent.
- [x] `SELECT c FROM t GROUP BY |` emits column-reference intent.
- [x] `SELECT c AS alias FROM t ORDER BY |` emits column-reference intent and Bytebase adapter returns select aliases.
- [x] `SELECT * FROM (SELECT c FROM t) x WHERE x.|` exposes virtual table `x`.
- [x] `WITH x AS (SELECT * FROM t) SELECT * FROM |` exposes CTE `x` as table-reference candidate.
- [x] `WITH x(c1, c2) AS (SELECT * FROM t) SELECT x.| FROM x` exposes explicit CTE columns.

## Phase 2: DML Completion

- [x] `INSERT INTO |` emits table-reference intent.
- [x] `INSERT INTO t (|)` emits column-reference intent scoped to table `t`.
- [x] `INSERT INTO t VALUES (|)` emits expression/column-reference context.
- [x] `INSERT INTO t SELECT | FROM u` emits column-reference intent scoped to `u`.
- [x] `UPDATE | SET c = 1` emits table-reference intent.
- [x] `UPDATE t SET |` emits column-reference intent scoped to table `t`.
- [x] `UPDATE t SET c = |` emits expression/column-reference context.
- [x] `UPDATE t SET c = 1 WHERE |` emits column-reference intent scoped to table `t`.
- [x] `DELETE FROM |` emits table-reference intent.
- [x] `DELETE FROM t WHERE |` emits column-reference intent scoped to table `t`.
- [x] `MERGE INTO |` emits table-reference intent.
- [x] `MERGE INTO t USING |` emits table-reference intent for source.
- [x] `MERGE INTO t USING u ON |` emits column-reference intent with target and source visible.

## Phase 3: DDL And Utility Completion

- [x] `CREATE |` suggests supported Oracle object-type keywords.
- [x] `CREATE TABLE t (c |)` emits datatype candidates.
- [x] `CREATE TABLE t (c NUMBER REFERENCES |)` emits table-reference intent.
- [x] `ALTER TABLE |` emits table-reference intent.
- [x] `ALTER TABLE t ADD |` suggests column/constraint action keywords.
- [x] `ALTER TABLE t DROP COLUMN |` emits column-reference intent scoped to table `t`.
- [x] `DROP TABLE |` emits table-reference intent.
- [x] `DROP VIEW |` emits view/table-reference intent.
- [x] `DROP SEQUENCE |` emits sequence-reference intent.
- [x] `TRUNCATE TABLE |` emits table-reference intent.
- [x] `COMMENT ON TABLE |` emits table-reference intent.
- [x] `COMMENT ON COLUMN t.|` emits column-reference intent scoped to table `t`.
- [x] `GRANT SELECT ON |` emits table-reference intent.
- [x] `REVOKE SELECT ON |` emits table-reference intent.

## Phase 4: Oracle-Specific Production Hardening

- [x] Quoted identifier prefix completion keeps the user's quoting mode.
- [x] Reserved keywords are not suggested as unquoted identifiers where parser rejects them.
- [x] Type keywords come from the Oracle keyword manifest.
- [x] Function-like keyword candidates appear only in call-capable expression positions.
- [x] Pseudo-column candidates appear in expression positions.
- [x] Case-sensitive metadata names can be quoted by the Bytebase adapter.
- [x] Multi-statement scripts isolate completion to the cursor statement.
- [x] Malformed earlier statements do not prevent completion in the current statement.
- [x] Large scripts avoid whole-file expensive parsing in completion hot path.
- [x] Oracle parser soft-fail, strictness, keyword, Loc, and corpus gates remain green.

## Phase 5: Bytebase Adapter Cutover

- [x] Bytebase Oracle completion calls omni Oracle completion APIs.
- [x] Bytebase Oracle completion no longer imports ANTLR or `github.com/bytebase/parser/plsql`.
- [x] Existing Bytebase `backend/plugin/parser/plsql/test-data/test_completion.yaml` passes unchanged or with only intended ordering updates.
- [x] Bytebase Oracle LSP completion continues returning `base.Candidate` schema/table/view/sequence/column/function/keyword values.
- [x] Bytebase Oracle completion preserves schema-as-database metadata behavior.
