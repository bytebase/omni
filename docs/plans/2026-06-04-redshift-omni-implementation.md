# Redshift Omni Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement a hand-written Omni Redshift dialect that can replace Bytebase's ANTLR-backed Redshift parser without regressing SQL review, SQL Editor, masking, completion, syntax diagnostics, statement reporting, or changed-resource extraction.

**Architecture:** Create `redshift` as a controlled fork of Omni `pg`: retain the mature PostgreSQL recursive-descent parser, AST, catalog, query analysis, and completion foundations, then add Redshift-specific syntax and adapters as explicit dialect deltas. Bytebase registration remains in Bytebase; Omni exposes plain package APIs that Bytebase glue can call.

**Tech Stack:** Go 1.25+, Omni recursive-descent parsers, Redshift legacy parser examples, AWS Redshift SQL command reference, Bytebase parser base registry contracts.

## Source Material

- `pg/` in this worktree is the implementation base.
- `SCENARIOS-redshift-omni.md` is the coverage contract.
- Legacy examples live in the Go module cache under `github.com/bytebase/parser/redshift/examples/redshift`.
- Bytebase Redshift consumers live in `backend/plugin/parser/redshift`, `backend/common/engine.go`, SQL Editor, export resources, LSP completion, and statement report.

## Current Implementation Status

Implemented in this Omni worktree:

- `redshift.Parse`, parser/AST/walker fork, and 115-file legacy Redshift corpus harness.
- Redshift parser deltas for CREATE TABLE distribution/sort/encoding options, COPY/UNLOAD, SHOW/DESC, datashare/external-object/policy/model-style utility statements, `MERGE ... REMOVE DUPLICATES`, `QUALIFY`, and `SELECT * EXCLUDE`.
- Runtime APIs for diagnostics, UTF-16 statement ranges, statement classification including distinct `EXPLAIN_ANALYZE`, SQL Editor read-only validation, and changed-resource extraction.
- `redshift/analysis.GetQuerySpan` over `redshift/catalog`, including CTEs, subqueries, set operations, view expansion, `QUALIFY` predicate lineage, `EXCLUDE` shaping, and explicit rejection for unsupported `CONNECT BY` lineage.
- `redshift/completion.Complete`, including Redshift top-level keywords, metadata schemas/tables/columns, schema-qualified relation filtering, CTE/subquery syntax-derived output columns, CREATE TABLE/COPY/UNLOAD slots, SHOW subcommands, and negative coverage for unsupported expression outputs.
- `redshift/compat`, an Omni-side compatibility verification layer with an AWS command manifest, command parse/support matrix tests, legacy corpus statistics, statement-level legacy parser parity statistics, runtime semantic API checks, an optional reference Redshift DSN harness, a report generator, and a checked-in freshness-tested Markdown report.

Still downstream or environment-bound:

- `go test ./redshift/...` passes in this worktree; catalog container tests skip cleanly when Docker is unavailable.
- Bytebase parser registration and feature adapters must be switched in a Bytebase worktree after this Omni commit is available.

Current compatibility snapshot:

- AWS command manifest: 124 commands classified (`27` runtime-supported, `84` parse-supported, `13` explicitly unsupported, `0` not relevant).
- Legacy Redshift corpus: 115 files tracked (`76` passing, `39` expected failures, `0` new failures, `0` promoted expected failures).
- Legacy statement parity: 5075 statements compared against the legacy ANTLR parser (`4247` accepted by both, `828` legacy-accepts/Omni-rejects, `0` Omni-only accepts).
- Runtime semantic checks: 8 checks, 8 passing, 0 failing, covering parse, diagnostics, ranges, statement classification, SQL Editor validation, changed resources, and completion scope.
- Reference Redshift harness: available through `REDSHIFT_COMPAT_DSN`; skipped in local verification when the DSN is unset.
- Generate the report with `go run ./redshift/compat/cmd/redshift-compat-report`.

## Phase 0: Scaffold And Corpus

### Task 0.1: Create Redshift Package Scaffold

**Files:**
- Create: `redshift/**`
- Test: `redshift/parse_test.go`

**Steps:**
1. Write a failing test that imports package `redshift` and expects `Parse("SELECT 1")` to return one statement.
2. Run `go test ./redshift -run TestParsePostgresCompatibleSelect -v`; expected failure is missing package/API.
3. Copy `pg` to `redshift` as a mechanical fork.
4. Replace package/import paths from `pg` to `redshift`.
5. Run `go test ./redshift/...`; expected pass for PostgreSQL-compatible copied tests.

### Task 0.2: Add Legacy Corpus Harness

**Files:**
- Create: `redshift/parser/legacy_examples_test.go`
- Create: `redshift/parser/testdata/legacy/*.sql`

**Steps:**
1. Copy the 115 Redshift legacy example SQL files into `redshift/parser/testdata/legacy`.
2. Write a test that iterates files and checks only the initial P0 allowlist.
3. Record remaining files as expected failures with reason strings.
4. Run `go test ./redshift/parser -run TestLegacyExamples -v`.

## Phase 1: Parser Foundation

### Task 1.1: Redshift CREATE TABLE Options

**Files:**
- Modify: `redshift/ast/parsenodes.go`
- Modify: `redshift/parser/parser.go`
- Test: `redshift/parser/redshift_create_table_test.go`

**Steps:**
1. Add failing tests for `DISTSTYLE`, `DISTKEY`, `SORTKEY`, `ENCODE AUTO`, `BACKUP NO`, column `IDENTITY`, column `ENCODE`, column `DISTKEY`, and column `SORTKEY`.
2. Extend AST with Redshift table and column option fields.
3. Extend `parseCreateTable` tail parsing to consume Redshift options after the column list.
4. Extend column definition parsing for Redshift column attributes.
5. Run `go test ./redshift/parser -run TestRedshiftCreateTable -v`.

### Task 1.2: Redshift COPY And UNLOAD

**Files:**
- Modify: `redshift/ast/parsenodes.go`
- Modify: `redshift/parser/copy.go`
- Create: `redshift/parser/unload.go`
- Test: `redshift/parser/redshift_copy_unload_test.go`

**Steps:**
1. Add failing tests for `COPY ... IAM_ROLE`, `COPY ... CREDENTIALS`, `COPY ... FORMAT AS CSV`, and representative `UNLOAD` options.
2. Extend COPY parsing for Redshift authorization and format clauses.
3. Add `UnloadStmt` AST and parser.
4. Run `go test ./redshift/parser -run 'TestRedshift(Copy|Unload)' -v`.

### Task 1.3: Redshift SHOW/DESC And Object Statements

**Files:**
- Modify: `redshift/ast/parsenodes.go`
- Create: `redshift/parser/show.go`
- Create: `redshift/parser/redshift_objects.go`
- Test: `redshift/parser/redshift_utility_test.go`
- Test: `redshift/parser/redshift_objects_test.go`

**Steps:**
1. Add failing tests for SHOW DATABASES/SCHEMAS/TABLES/COLUMNS/GRANTS/DATASHARES/MODEL/PROCEDURE/VIEW.
2. Add failing tests for DESC DATASHARE and DESC IDENTITY PROVIDER.
3. Add failing tests for DATASHARE, EXTERNAL SCHEMA/TABLE/VIEW, MASKING/RLS policy, IDENTITY PROVIDER, and MODEL statements.
4. Implement parser support with typed AST for consumed names and raw option capture for complex option lists.
5. Run focused tests.

## Phase 2: Bytebase Runtime APIs

### Task 2.1: Diagnostics, Statement Ranges, And Types

**Files:**
- Create: `redshift/diagnose.go`
- Create: `redshift/ranges.go`
- Create: `redshift/classify.go`
- Test: `redshift/diagnose_test.go`
- Test: `redshift/ranges_test.go`
- Test: `redshift/classify_test.go`

**Steps:**
1. Write failing tests for valid SQL diagnostics, invalid SQL diagnostics, UTF-16 ranges, and statement classification.
2. Implement diagnostic conversion using parser errors.
3. Implement statement ranges from split results.
4. Implement statement classification.
5. Run `go test ./redshift -run 'Test(Diagnose|StatementRanges|Classify)' -v`.

### Task 2.2: SQL Editor Validation

**Files:**
- Create: `redshift/validate.go`
- Test: `redshift/validate_test.go`

**Steps:**
1. Add failing tests for SELECT, EXPLAIN, SHOW, SET, DML, DDL, and EXPLAIN ANALYZE behavior.
2. Implement conservative read-only/returns-data classification.
3. Run `go test ./redshift -run TestValidateSQLForEditor -v`.

### Task 2.3: Changed Resources

**Files:**
- Create: `redshift/changes.go`
- Test: `redshift/changes_test.go`

**Steps:**
1. Add failing tests for CREATE TABLE, DROP TABLE, ALTER TABLE, INSERT, UPDATE, DELETE, and `SET search_path`.
2. Implement changed-resource extraction over Redshift AST.
3. Run `go test ./redshift -run TestExtractChangedResources -v`.

## Phase 3: Query Span

### Task 3.1: SELECT Query Span Foundation

**Files:**
- Create: `redshift/analysis/query_span.go`
- Test: `redshift/analysis/query_span_test.go`

**Steps:**
1. Add failing tests for simple columns, `SELECT *`, table-qualified star, aliases, expressions, WHERE, JOIN ON, JOIN USING, CTEs, subqueries, set operations, and views.
2. Adapt `redshift/catalog` analyzer for Redshift metadata/search path defaults.
3. Convert analyzed query results to an Omni-native Redshift query span shape.
4. Run `go test ./redshift/analysis -run TestQuerySpan -v`.

### Task 3.2: Redshift SELECT Extensions

**Files:**
- Modify: `redshift/parser/select.go`
- Modify: `redshift/analysis/query_span.go`
- Test: `redshift/analysis/query_span_redshift_test.go`

**Steps:**
1. Add failing tests for `QUALIFY`, `EXCLUDE`, `SELECT INTO`, and unsupported `CONNECT BY`.
2. Parse extension clauses.
3. Implement lineage where correct; emit explicit unsupported errors where not.
4. Run focused query span tests.

## Phase 4: Completion

### Task 4.1: Core Completion

**Files:**
- Create: `redshift/completion/completion.go`
- Test: `redshift/completion/completion_test.go`

**Steps:**
1. Add failing tests for top-level candidates, schema/table/column candidates, table alias columns, CTE names, CTE columns, and subquery alias columns.
2. Adapt PG completion context and metadata resolution for Redshift.
3. Use parser completion scope and syntax-derived target lists for CTE/subquery column inference; do not invent names for unsupported expression outputs.
4. Run `go test ./redshift/completion -run TestCompletion -v`.

### Task 4.2: Redshift-Specific Completion Slots

**Files:**
- Modify: `redshift/completion/completion.go`
- Test: `redshift/completion/redshift_slots_test.go`

**Steps:**
1. Add failing tests for CREATE TABLE option slots, COPY slots, UNLOAD slots, and SHOW subcommands.
2. Implement slot-aware Redshift keyword candidates.
3. Run focused completion tests.

## Phase 5: Bytebase Switch

### Task 5.1: Downstream Adapter Switch

**Files:**
- Modify downstream Bytebase only in a Bytebase worktree, not this Omni worktree.

**Steps:**
1. Update Bytebase `go.mod` to the Omni commit containing Redshift support.
2. Replace Redshift parser package internals to call Omni Redshift APIs.
3. Keep Bytebase `base.Register*` blocks in place.
4. Run Bytebase focused Redshift tests.
5. Run SQL Editor, masking/query-span, statement report, and completion tests.

### Bytebase Adapter Mapping

Use these Omni APIs as the downstream contract:

- Parse/split/diagnose/ranges/types:
  - `redshift.Parse(sql)`
  - `redshift.Diagnose(sql)`
  - `redshift.StatementRanges(sql)`
  - `redshift.GetStatementTypes(sql)` / `redshift.ClassifyStatement(stmt)`
- SQL Editor read-only validation:
  - `redshift.ValidateSQLForEditor(sql)`
  - Preserve Bytebase's existing read-only policy surface; map Omni's `returnsData`, `readOnly`, and `error` into the current Redshift validator result type.
- Changed resources and statement report:
  - `redshift.ExtractChangedResources(sql, currentDatabase, currentSchema)`
  - Map `ChangeSummary.Resources` and DML counts into Bytebase's existing affected-resource/report structs.
- Masking/query span:
  - Build a `redshift/catalog.Catalog` from Bytebase metadata, then call `redshift/analysis.GetQuerySpan(cat, sql)`.
  - Treat explicit unsupported errors as "unsupported lineage" rather than returning partial or misleading masking data.
- Completion:
  - Build the same metadata catalog and call `redshift/completion.Complete(sql, cursorOffset, cat)`.
  - Convert Omni candidate kinds to Bytebase LSP completion item kinds.

Recommended downstream verification order:

1. `go test ./backend/plugin/parser/redshift`
2. Focused SQL Editor validation tests for Redshift SELECT/SHOW/EXPLAIN vs DML/DDL.
3. Focused changed-resource and statement-report tests for CREATE/ALTER/DROP/INSERT/UPDATE/DELETE/MERGE/SELECT INTO.
4. Focused masking/query-span tests for SELECT, `QUALIFY`, `EXCLUDE`, CTE, subquery, and unsupported `CONNECT BY`.
5. Focused LSP completion tests for metadata columns, CTE/subquery output columns, CREATE TABLE slots, COPY/UNLOAD options, and SHOW subcommands.

## Global Verification

Run these before claiming completion:

```bash
go test ./redshift/compat -v
go test ./redshift/...
go test ./pg/...
```

When a real Redshift cluster is available:

```bash
REDSHIFT_COMPAT_DSN='postgres://...' go test ./redshift/compat -run TestCompatibilityReportMarkdown -v
REDSHIFT_COMPAT_DSN='postgres://...' go run ./redshift/compat/cmd/redshift-compat-report
```

After downstream switch:

```bash
go test ./backend/plugin/parser/redshift
go test ./backend/runner/plancheck -run Redshift
go test ./backend/api/v1 -run Redshift
```
