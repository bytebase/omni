# Oracle Completion Scope Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Define the Oracle parser-embedded completion scope that can begin after the parser coverage gates are green, without pulling in catalog, resolver, or execution semantics.

**Architecture:** Completion must be parser-embedded. It should reuse the Oracle lexer, token positions, parser clause state, and soft-fail paths to collect candidates at the cursor. The first implementation should return syntactic candidates only: statement starters, clause keywords, operators, datatype/function names, and placeholder object-shape suggestions. Catalog-backed table, column, schema, synonym, procedure, and privilege suggestions remain out of scope until Oracle catalog/resolver work exists.

**Tech Stack:** Go implementation under `oracle/parser` unless a sibling `oracle/completion` package already exists by the time this plan is executed; AST and Loc helpers under `oracle/ast`; parser status gates in `oracle/parser/*_test.go`; coverage manifests under `oracle/parser/testdata/coverage`; status docs under `docs/`.

## Readiness Gate

Do not start implementation unless these parser-layer gates pass:

- `TestOracleParserProgress`: all parser parse methods return `error`; silent parser error discards are zero.
- `TestOracleCoverage`: soft-fail, strictness, keyword, BNF, Loc-node, and reference-oracle minimums are met.
- `TestVerifyCorpus`: parser mismatches, Loc violations, and crashes are zero for the current corpus.
- `TestOracleKeywordManifestExhaustive`: every declared Oracle lexer keyword is classified.
- `TestOracleLocNodeCoverage`: every Loc-bearing node type is classified, and covered fixtures pass Loc validation.

Current measured floor:

- Soft-fail scenarios: 62
- Strictness scenarios: 121
- Keyword manifest rows: 344
- Oracle 26ai SQL reserved-word audit rows: 110
- BNF rows classified: 171
- BNF unknown rows: 0
- Loc-node rows classified: 249
- Direct Loc fixture rows: 152
- Loc-node unknown rows: 0
- Optional reference-oracle rows: 20

## Non-Goals

- No catalog lookup, schema introspection, privilege resolution, synonym expansion, or type compatibility analysis.
- No external completion algorithm or parser generator.
- No guarantee that every unsupported Oracle grammar feature has completion candidates.
- No user-facing SQL formatting or rewriting.
- No dependency on an Oracle server in the default test suite.

## Phase 0: Completion Boundary

**Files:**

- Inspect existing completion packages for PG, MySQL, and MSSQL.
- Create Oracle completion entry points only after matching the local package boundary.

**Tasks:**

- Identify the common completion API shape: request fields, cursor offset, candidate shape, replacement span, and tests.
- Decide whether Oracle completion belongs in `oracle/parser` or a sibling package by matching existing engine conventions.
- Add a minimal Oracle completion request/response type if no shared type exists.
- Add a smoke test that an empty input at cursor 0 returns statement-start candidates.

**Acceptance:**

- Oracle has a callable completion entry point with no catalog dependency.
- The entry point preserves cursor byte offsets and returns stable replacement spans.

## Phase 1: Statement Starters

**Tasks:**

- Return top-level candidates for `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `MERGE`, `CREATE`, `ALTER`, `DROP`, `TRUNCATE`, `COMMENT`, `GRANT`, `REVOKE`, `COMMIT`, `ROLLBACK`, `SAVEPOINT`, `SET`, `BEGIN`, and `DECLARE`.
- Keep candidates case-normalized while preserving caller-provided prefix matching.
- Treat semicolon-separated multi-statement input as a fresh statement context after the latest statement boundary.

**Acceptance:**

- Empty input, whitespace-only input, and input after semicolon produce statement starters.
- Truncated statement starts such as `SE` and `CRE` filter candidates without parse panics.

## Phase 2: SELECT Clause Completion

**Tasks:**

- Track SELECT clause order with parser state: target list, FROM, WHERE, GROUP BY, HAVING, ORDER BY, OFFSET, FETCH, set operators.
- Suggest only legal next clauses for the current clause position.
- Exclude duplicate clauses already rejected by strictness tests.
- Suggest expression operators in expression contexts where the soft-fail matrix proves truncated input is recoverable.

**Acceptance:**

- `SELECT 1 ` suggests `FROM`, set operators, and legal row-limiting/order clauses.
- `SELECT * FROM t WHERE ` suggests expression starters, not top-level statement starters.
- Duplicate `WHERE`, `GROUP BY`, `HAVING`, `ORDER BY`, and `FETCH` are not suggested after already appearing.

## Phase 3: DML Completion

**Tasks:**

- Add syntactic candidates for `INSERT INTO`, column lists, `VALUES`, `SELECT` source, and `RETURNING`.
- Add syntactic candidates for `UPDATE ... SET`, `WHERE`, and `RETURNING`.
- Add syntactic candidates for `DELETE FROM`, `WHERE`, and `RETURNING`.
- Add MERGE candidates for `INTO`, `USING`, `ON`, `WHEN MATCHED`, `WHEN NOT MATCHED`, `THEN`, `UPDATE SET`, `INSERT`, and `VALUES`.

**Acceptance:**

- DML completion follows the same mandatory operand rules as strictness tests.
- MERGE completion never suggests action clauses before mandatory `ON`.

## Phase 4: DDL And Utility Completion

**Tasks:**

- Add CREATE object-type candidates for supported parser families: table, view, index, sequence, synonym, procedure, function, package, trigger, type, user, role.
- Add ALTER object-type and action candidates for supported parser families.
- Add DROP object-type candidates and known parser options.
- Add GRANT/REVOKE privilege-position and grantee-position syntactic candidates without catalog validation.
- Add COMMENT, TRUNCATE, SET, ALTER SESSION, transaction, and SAVEPOINT utility candidates.

**Acceptance:**

- Unknown option positions use explicit allowlists rather than accepting arbitrary keywords.
- Unsupported or catalog-only BNF families do not leak misleading object-specific candidates.

## Phase 5: Expressions, Functions, And Datatypes

**Tasks:**

- Reuse the 344-row keyword manifest to drive keyword/function/pseudo-column candidates.
- Suggest function-like keywords only in call-capable expression positions.
- Suggest pseudo-columns only in expression positions where the parser accepts them.
- Suggest datatype keywords in datatype contexts.
- Respect quoted identifiers and dotted-name keyword rules from keyword golden tests.

**Acceptance:**

- Keyword candidates match the manifest classification policy.
- Completion does not suggest reserved words as unquoted identifiers in contexts where the parser rejects them.

## Phase 6: Completion Coverage And Status

**Tasks:**

- Add a completion scenario manifest that mirrors the parser coverage style.
- Add tests for empty input, truncated SQL, legal prefix filtering, illegal duplicate clauses, and multi-statement cursor positions.
- Add parser-progress logging for Oracle completion scenario counts once implementation exists.
- Update `docs/PARSER-DEFENSE-MATRIX.md` and `docs/engine-capability-guide.md` only after tests prove completion behavior.

**Acceptance:**

- Oracle L6 can move from `--` to measured status only with passing completion tests.
- Completion failures do not weaken parser strictness or corpus gates.
