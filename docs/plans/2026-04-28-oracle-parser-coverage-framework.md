# Oracle Parser Coverage Framework Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Close the remaining Oracle parser-layer coverage gap against MySQL, PG, and MSSQL after the foundation/completion slice by making grammar coverage, strictness breadth, keyword classification, Loc coverage, and reference-oracle behavior measurable and enforceable.

**Architecture:** Add coverage accounting before adding more parser behavior. The first phase creates manifests and guard tests that explain what is covered, partial, missing, unsupported, or intentionally deferred. Later phases expand table-driven tests and optional reference-oracle checks until Oracle has MySQL/MSSQL-style exhaustive keyword and strictness coverage, PG-style scenario attribution, and a clear readiness gate for completion work.

**Tech Stack:** Go tests under `oracle/parser`, Oracle AST node types under `oracle/ast`, local BNF files under `oracle/parser/bnf`, corpus verifier in `oracle/parser/verify_parse_test.go`, optional Oracle reference harness behind a build tag and environment variable, status docs in `docs/`.

## Initial Baseline

Parser-layer state before this coverage-framework plan:

- `TestOracleParserProgress` reports `445/445` parser `parse*` methods with error returns.
- Production parser silent error discards are `0`.
- Soft-fail coverage is `62` scenarios.
- Strictness coverage is `29` scenarios.
- Keyword golden coverage is `62` scenarios.
- Corpus verification covers `128` statements with `0` parse violations, `0` Loc violations, and `0` crashes.
- Oracle has `171` BNF files under `oracle/parser/bnf`, but no BNF-to-test coverage map.
- Oracle completion is still not started; this plan makes completion a later phase, not the next edit.

The current gap is coverage depth and reference quality, not basic parser error propagation.

## Execution Status

This plan has been implemented through Phase 7. Current parser-layer gates:

- Soft-fail coverage is `62` scenarios.
- Strictness coverage is `121` scenarios.
- Keyword coverage is an exhaustive `344`-row manifest for the declared Oracle lexer keyword set, with a pinned Oracle 26ai SQL reserved-word audit manifest.
- BNF coverage has `171` classified rows with `0` unknown rows.
- High-value BNF families have no `missing` or `unknown` rows.
- Loc-node coverage has `249` classified rows with `0` unknown rows and `152` direct SQL fixtures.
- Optional reference-oracle coverage has a `20`-row manifest and skips cleanly without `ORACLE_PARSER_REF_DSN` or explicit `ORACLE_PARSER_REF_CONTAINER=1`.
- Completion remains a separate scoped plan: `docs/plans/2026-04-28-oracle-completion-scope.md`.

## Non-Goals

- Do not include catalog, resolver, schema semantics, or execution behavior.
- Do not require a local Oracle database for the default test suite.
- Do not rewrite the parser with a parser generator.
- Do not make every BNF production a hard failure on day one; first remove unknown coverage, then raise thresholds.
- Do not start completion implementation until parser coverage gates are stable.

## Canonical Proof Commands

Focused Oracle parser proof:

```bash
go test -run 'TestOracleParserProgress|TestVerifyCorpus|TestOracleKeyword|TestSoftFail|TestStrict|TestOracleLoc|TestOracleCoverage' -count=1 -v ./oracle/parser
go test -count=1 ./oracle/parser
git diff --check
```

Optional reference-oracle proof, skipped by default:

```bash
ORACLE_PARSER_REF_DSN="$ORACLE_PARSER_REF_DSN" go test -tags oracle_ref -run 'TestOracleReference|TestOracleVReservedWordsKeywordAudit' -count=1 -v ./oracle/parser
ORACLE_PARSER_REF_CONTAINER=1 go test -tags oracle_ref -run 'TestOracleReference|TestOracleVReservedWordsKeywordAudit' -count=1 -v ./oracle/parser
```

Repository-wide proof is informational for this plan because non-Oracle packages may have unrelated container/catalog dependencies:

```bash
go test ./...
```

## Coverage Model

Use a single vocabulary across manifests, tests, and docs:

| Status | Meaning | Allowed At Final Gate |
|---|---|---|
| `covered` | Positive and negative tests exist for the parser surface. | Yes |
| `partial` | Some scenarios exist, but important variants are missing. | Yes, with issue text |
| `missing` | No meaningful parser test coverage exists yet. | No |
| `unsupported` | Parser intentionally does not support the Oracle feature yet. | Yes, with reason |
| `catalog` | Requires catalog or semantic resolution, outside parser scope. | Yes |
| `deferred` | Parser-layer work exists but belongs to a later named phase. | Yes, with owner |
| `unknown` | Not classified. | No |

Final parser coverage does not mean every Oracle feature parses. It means every discovered parser surface is classified, measured, and either tested or explicitly out of scope.

## Phase 0: Coverage Accounting

### Task 0.1: Add BNF Coverage Manifest

**Files:**
- Create: `oracle/parser/testdata/coverage/bnf_coverage.tsv`
- Create: `oracle/parser/coverage_matrix_test.go`

**Design:**

Create a tab-separated manifest with one row per `oracle/parser/bnf/*.bnf` file:

```text
bnf_file	family	status	tests	notes
select.bnf	select	covered	oracle/parser/select_test.go;oracle/parser/soft_fail_strict_test.go	core SELECT has positive and strictness coverage
```

`TestOracleBNFCoverageManifestCompleteness` must:

- List all BNF files from `oracle/parser/bnf`.
- Fail if any BNF file is absent from the manifest.
- Fail if any row has `unknown`.
- Fail if status is not one of the approved coverage statuses.
- Log counts by family and status.

It should not initially fail on `missing` or `partial`; those become tracked debt for later phases.

**Proof:**

```bash
go test -run 'TestOracleBNFCoverageManifestCompleteness' -count=1 -v ./oracle/parser
```

Expected after Task 0.1: PASS, with a status summary showing current `covered`, `partial`, `missing`, `unsupported`, `catalog`, and `deferred` counts.

### Task 0.2: Add Scenario Coverage Counters

**Files:**
- Modify: `oracle/parser/parser_progress_test.go`
- Create or modify: `oracle/parser/testdata/coverage/scenario_targets.tsv`

**Design:**

Track scenario counts separately from test function counts:

- soft-fail scenarios
- strictness scenarios
- keyword golden scenarios
- Loc fixture scenarios
- reference-oracle scenarios
- BNF rows classified

Use thresholds as monotonic guards. Final parser-coverage thresholds are:

- soft-fail `>= 62`
- strictness `>= 100`
- keyword `>= 344`
- BNF rows classified `>= len(oracle/parser/bnf)`
- unknown BNF rows `== 0`
- Loc-node unknown rows `== 0`
- Loc-node direct fixture rows `>= 80`
- reference-oracle rows `>= 20`

Later phases can raise these thresholds after new scenarios land.

**Proof:**

```bash
go test -run 'TestOracleParserProgress|TestOracleCoverage' -count=1 -v ./oracle/parser
```

Expected: PASS and print the expanded coverage counters.

## Phase 1: Strictness V2

### Task 1.1: Expand Strictness Taxonomy

**Files:**
- Modify: `oracle/parser/soft_fail_strict_test.go`
- Create: `oracle/parser/testdata/coverage/strictness_v2.tsv`

**Design:**

Grow strictness from representative coverage to MySQL-like breadth. Target at least `100` strictness scenarios, grouped by parser surface:

| Family | Target Count | Examples |
|---|---:|---|
| SELECT clause order and duplication | 18 | duplicate `WHERE`, `GROUP BY`, `HAVING`, `ORDER BY`, bad clause order |
| SELECT operands and predicates | 16 | missing RHS for `BETWEEN`, `IN`, `LIKE`, `IS`, arithmetic, concatenation |
| Set operations | 8 | trailing `UNION`, duplicate modifiers, missing RHS query |
| INSERT / UPDATE / DELETE / MERGE | 18 | missing target, missing value source, duplicate `SET`, incomplete `USING` |
| CREATE TABLE / INDEX / VIEW | 18 | unknown options, duplicate exclusive options, illegal keyword position |
| ALTER / DROP / GRANT / SESSION | 16 | unknown action, malformed option payload, incomplete name |
| PL/SQL block structure | 10 | missing `END`, missing `THEN`, malformed exception section |

Every negative strictness scenario should have a nearby positive control when the grammar supports it. Unsupported valid Oracle syntax should be classified separately, not hidden as strictness.

**Proof:**

```bash
go test -run 'TestStrict' -count=1 -v ./oracle/parser
go test -run 'TestOracleParserProgress' -count=1 -v ./oracle/parser
```

Expected after Phase 1: strictness counter `>= 100`.

## Phase 2: Keyword Exhaustiveness

### Task 2.1: Replace Representative Keyword Golden With Exhaustive Manifest

**Files:**
- Modify: `oracle/parser/keyword_class_test.go`
- Modify: `oracle/parser/keywords.go`
- Create: `oracle/parser/testdata/coverage/oracle_keywords.tsv`

**Design:**

Move keyword tests from representative samples to a full manifest. Each row should include:

```text
word	token	category	table_name	column_name	alias	dotted_name	function_call	pseudocolumn	notes
SELECT	SELECT	reserved	reject	reject	reject	reject	reject	no	Oracle reserved keyword
LEVEL	LEVEL	pseudocolumn	allow	allow	allow	allow	reject	yes	Expression keyword in pseudo-column position
```

Categories:

- `reserved`
- `nonreserved`
- `context`
- `function`
- `pseudocolumn`
- `datatype`
- `plsql`
- `operator`

Guards:

- Every lexer keyword token must appear in the manifest.
- Every manifest row must map to a known token or be marked `documented_only`.
- Reserved words must reject unquoted object and column identifiers.
- Quoted identifiers must override keyword classification.
- Context keywords must be accepted only in contexts where Oracle permits them.
- Function-like keywords must distinguish call position from identifier position.
- Pseudo-columns must distinguish expression position from quoted identifier position.

**Reference Source Policy:**

Pin the keyword manifest to a named Oracle SQL Language Reference version in a header comment. Do not let the test fetch documentation at runtime.

**Proof:**

```bash
go test -run 'TestOracleKeyword' -count=1 -v ./oracle/parser
go test -run 'TestOracleParserProgress' -count=1 -v ./oracle/parser
```

Expected after Phase 2: keyword golden counter reflects manifest rows, not only the current `62` representative scenarios.

## Phase 3: Loc Node-Type Coverage

### Task 3.1: Add AST Loc Coverage Manifest

**Files:**
- Create: `oracle/parser/testdata/coverage/loc_node_coverage.tsv`
- Create: `oracle/parser/loc_node_coverage_test.go`
- Modify: `oracle/parser/loc_contract_test.go`

**Design:**

The corpus Loc verifier proves parsed statements do not emit invalid spans. It does not prove every important AST node type has fixture coverage.

Add a node-type manifest:

```text
node_type	family	status	fixture	notes
*nodes.SelectStmt	select	covered	SELECT a FROM t WHERE a > 1	basic select span
*nodes.CreateTableStmt	ddl	covered	CREATE TABLE t (a NUMBER)	create table span
```

The test should:

- Enumerate AST node types that implement the local node/Loc contract.
- Fail if a Loc-bearing node type is absent from the manifest.
- Parse each fixture and verify the expected node type appears.
- Verify `Loc.Start >= 0`, `Loc.End >= Loc.Start`, and child spans stay within parent spans when both are known.
- Allow `unsupported`, `catalog`, or `deferred` rows with a reason.

**Proof:**

```bash
go test -run 'TestOracleLoc|TestOracleLocNodeCoverage' -count=1 -v ./oracle/parser
```

Expected after Phase 3: no unknown Loc-bearing node types, corpus Loc verifier still reports `0` violations.

## Phase 4: Optional Oracle Reference Harness

### Task 4.1: Add Gated Reference Runner

**Files:**
- Create: `oracle/parser/reference_oracle_test.go`
- Create: `oracle/parser/testdata/coverage/reference_oracle.tsv`

**Design:**

Add a build-tagged test file:

```go
//go:build oracle_ref
```

The test runs only when `ORACLE_PARSER_REF_DSN` is set or `ORACLE_PARSER_REF_CONTAINER=1` is explicitly requested. It should skip by default.

Reference strategy:

- Use the real Oracle parser as the accept/reject oracle.
- Keep DDL in a disposable schema or random object namespace.
- Prefer `DBMS_SQL.PARSE` for parse classification where possible.
- Keep destructive or catalog-dependent SQL out of the default reference set.
- Classify each reference row as `accept`, `reject`, `catalog`, or `unsafe`.

Comparison rule:

- If Oracle accepts and omni rejects, record `missing_valid_syntax`.
- If Oracle rejects and omni accepts, record `over_accept`.
- If both agree, pass.
- If Oracle requires catalog state, mark the row `catalog` and keep it out of parser-only gates.

**Proof:**

```bash
ORACLE_PARSER_REF_DSN="$ORACLE_PARSER_REF_DSN" go test -tags oracle_ref -run 'TestOracleReference|TestOracleVReservedWordsKeywordAudit' -count=1 -v ./oracle/parser
ORACLE_PARSER_REF_CONTAINER=1 go test -tags oracle_ref -run 'TestOracleReference|TestOracleVReservedWordsKeywordAudit' -count=1 -v ./oracle/parser
```

Expected without DSN/container: skipped. Expected with DSN or explicit container: mismatch report grouped by family plus a `V$RESERVED_WORDS` keyword audit.

## Phase 5: BNF Gap Closure

### Task 5.1: Close High-Value BNF Families

**Files:**
- Modify: `oracle/parser/testdata/coverage/bnf_coverage.tsv`
- Modify: relevant `oracle/parser/*_test.go`
- Modify: `oracle/parser/soft_fail_strict_test.go`

**Design:**

Use the BNF manifest to close gaps in this order:

1. Query core: SELECT, expressions, joins, subqueries, set operations.
2. DML: INSERT, UPDATE, DELETE, MERGE.
3. DDL high-frequency: CREATE TABLE, ALTER TABLE, CREATE INDEX, CREATE VIEW, DROP, COMMENT.
4. PL/SQL parser surface: block structure, procedure/function/package skeletons, triggers.
5. Administrative statements already represented in parser code.
6. Low-value or catalog-dependent grammar marked `catalog`, `unsupported`, or `deferred`.

For each BNF family, add:

- at least one positive parse case
- at least one negative strictness or truncation case when applicable
- Loc coverage if the family creates distinct AST nodes
- keyword coverage when the family has context keywords

**Proof:**

```bash
go test -run 'TestOracleBNFCoverageManifestCompleteness|TestSoftFail|TestStrict|TestVerifyCorpus' -count=1 -v ./oracle/parser
```

Expected after Phase 5: no `missing` rows for high-value parser families.

## Phase 6: Progress And Defense Matrix Integration

### Task 6.1: Make Coverage Status Visible

**Files:**
- Modify: `oracle/parser/parser_progress_test.go`
- Modify: `docs/PARSER-DEFENSE-MATRIX.md`
- Modify: `docs/engine-capability-guide.md`
- Modify: `oracle/parser/SCENARIOS-oracle-parser-completion.md` or create a follow-up scenario file if the completion file should remain frozen.

**Design:**

Update progress reporting so the Oracle row can say more than `Done` or `Partial`:

- `parse_methods_with_error`
- `silent_error_discards`
- `soft_fail_scenarios`
- `strict_scenarios`
- `keyword_manifest_rows`
- `bnf_rows_total`
- `bnf_rows_missing`
- `loc_node_rows_unknown`
- `reference_oracle_rows`
- `reference_oracle_mismatches`

Docs should distinguish:

- foundation complete
- coverage framework complete
- strictness breadth target met
- keyword exhaustive target met
- reference oracle optional and gated
- completion not started until explicitly scoped

**Proof:**

```bash
go test -run 'TestOracleParserProgress|TestOracleCoverage' -count=1 -v ./oracle/parser
git diff --check
```

Expected: progress test emits all new counters and docs match the measured state.

## Phase 7: Completion Readiness Gate

### Task 7.1: Define Completion Scope After Parser Coverage Stabilizes

**Files:**
- Create: `docs/plans/2026-04-28-oracle-completion-scope.md`
- Create: `oracle/completion/SCENARIOS-oracle-completion.md` if Oracle completion package exists, otherwise document the package boundary first.

**Design:**

Do not start completion until these parser gates are true:

- BNF manifest has no `unknown` rows.
- High-value parser families have no `missing` rows.
- Strictness counter is at least `100`.
- Keyword manifest is exhaustive for lexer keywords.
- Loc node manifest has no unknown Loc-bearing node types.
- Optional reference-oracle harness exists and skips cleanly without DSN/container.

Completion scope should then follow the MySQL/PG/MSSQL pattern:

- lexical completion at statement start
- identifier and keyword completion in SELECT/DML/DDL contexts
- function and datatype suggestions
- clause-order-aware suggestions
- no catalog dependency in parser-layer completion tests unless explicitly tagged

## Execution Shape

Recommended order is sequential for Phase 0 because every later task depends on the manifests and counters.

After Phase 0:

- Phase 1 strictness and Phase 2 keyword can run independently.
- Phase 3 Loc node coverage can run independently once the Loc manifest format is fixed.
- Phase 4 reference oracle can run independently because it is build-tagged and skipped by default.
- Phase 5 BNF gap closure should run after the first four phases, using their manifests as the work queue.
- Phase 6 docs/status integration runs after each substantial batch.
- Phase 7 completion scope is last.

## Done Definition

This plan is complete when:

- `oracle/parser` tests pass.
- BNF coverage manifest has every BNF file classified and no `unknown` rows.
- Strictness coverage is at least `100` scenarios.
- Keyword coverage is exhaustive for lexer keywords and context behavior is tested from the manifest.
- Loc coverage has no unknown Loc-bearing AST node types.
- Optional reference-oracle tests skip cleanly without DSN/container and produce grouped mismatches with DSN or explicit container.
- Defense matrix and engine capability docs report measured Oracle status.
- Oracle completion has a separate scoped plan and is no longer mixed into parser foundation work.
