# MySQL Functional Indexes Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement MySQL 8.0 functional index catalog behavior for C1.11/C1.12 and C19: hidden generated columns, type inference, validation, visibility, JSON expression normalization, and index lifecycle cleanup.

**Architecture:** Treat each functional key part as MySQL does: synthesize a hidden virtual generated `Column` and keep the public `IndexColumn.Expr` for SHOW/deparse fidelity. Add a hidden-column visibility model to `Table`/`Column`, then make CREATE TABLE and CREATE INDEX share one synthesis pipeline. Layer expression type inference and validation on top of that pipeline rather than embedding ad hoc checks in deparse.

**Tech Stack:** Go, `mysql/catalog`, existing `mysql/ast` expression nodes, existing catalog analyzer `ResolvedType`, Docker-backed MySQL oracle scenario tests.

## Current State

Parser support is mostly present: table-level and standalone indexes carry `ast.IndexColumn.Expr`, and catalog `IndexColumn` already has `Expr string`. The gap is below parsing:

- `mysql/catalog/index.go`: `IndexColumn` stores either `Name` or `Expr`, but there is no link to a synthesized hidden generated column.
- `mysql/catalog/table.go`: `Column` has `Invisible` and GIPK-specific fields, but no `HiddenBySystem` model for MySQL `HT_HIDDEN_SQL`.
- `mysql/catalog/tablecmds.go` and `mysql/catalog/indexcmds.go`: functional key parts are accepted as raw expression strings; no hidden columns are synthesized and no functional-index validation runs.
- `mysql/catalog/show.go`: `showIndex` can render `Expr`, but `ShowCreateTable` does not understand system-hidden columns because none exist.
- `mysql/catalog/analyze_expr.go` / `function_types.go`: partial expression typing exists and can be reused, but VarExprQ currently does not get populated from catalog columns in the standalone table-expression context needed for index expressions.

## Phase 1: Hidden Column Architecture and Functional Index Synthesis

### Task 1.1: Add Hidden Column Model

**Files:**
- Modify: `mysql/catalog/table.go`
- Modify: `mysql/catalog/scope.go`
- Modify: `mysql/catalog/show.go`
- Test: `mysql/catalog/functional_index_test.go`

**Step 1: Write failing test**

Add `TestFunctionalIndexHiddenColumns_VisibleAndHiddenViews`:

```go
func TestFunctionalIndexHiddenColumns_VisibleAndHiddenViews(t *testing.T) {
	c := mustLoadSQL(t, `
		CREATE DATABASE testdb;
		USE testdb;
		CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(64));
	`)
	db := c.GetDatabase("testdb")
	tbl := db.GetTable("t")
	tbl.Columns = append(tbl.Columns, &Column{Name: "!hidden!idx_lower!0!0", DataType: "varchar", ColumnType: "varchar(64)", Hidden: ColumnHiddenSystem})
	rebuildColIndex(tbl)

	if got := len(tbl.VisibleColumns()); got != 2 { t.Fatalf("visible=%d", got) }
	if got := len(tbl.HiddenColumns()); got != 1 { t.Fatalf("hidden=%d", got) }
	if strings.Contains(c.ShowCreateTable("testdb", "t"), "!hidden!") { t.Fatal("hidden column leaked") }
}
```

**Step 2: Run RED**

Run: `go test ./mysql/catalog -run 'TestFunctionalIndexHiddenColumns_VisibleAndHiddenViews' -count=1`

Expected: compile failure because `Column.Hidden`, `ColumnHiddenSystem`, `VisibleColumns`, and `HiddenColumns` do not exist.

**Step 3: Implement minimal model**

Add:

```go
type ColumnHiddenKind int
const (
	ColumnHiddenNone ColumnHiddenKind = iota
	ColumnHiddenSystem
)
```

Add `Hidden ColumnHiddenKind` to `Column`. Add `VisibleColumns()` and `HiddenColumns()` on `Table`. Update analyzer scope registration and `ShowCreateTable` column loop to use visible columns where user SQL visibility matters. Keep `GetColumn` internal and unchanged for now.

**Step 4: Run GREEN**

Run: `go test ./mysql/catalog -run 'TestFunctionalIndexHiddenColumns_VisibleAndHiddenViews|TestShowCreateTableBasic|TestWalkThrough_3_1_ColumnPositionsSequential' -count=1`

Expected: PASS.

**Step 5: Commit**

```bash
git add mysql/catalog/table.go mysql/catalog/scope.go mysql/catalog/show.go mysql/catalog/functional_index_test.go
git commit -m "add mysql system hidden column model"
```

### Task 1.2: Synthesize Hidden Columns for CREATE TABLE Functional Indexes

**Files:**
- Modify: `mysql/catalog/tablecmds.go`
- Modify: `mysql/catalog/index.go`
- Test: `mysql/catalog/functional_index_test.go`
- Update: `mysql/catalog/scenarios_c1_test.go`

**Step 1: Write failing tests**

Add tests for:

- Unnamed functional indexes become `functional_index`, `functional_index_2`.
- `INDEX fx ((a + 1), (a * 2))` synthesizes `!hidden!fx!0!0` and `!hidden!fx!1!0`.
- The index keeps `IndexColumn.Expr` for SHOW rendering and stores `IndexColumn.Name` as the hidden column name for internal lifecycle linkage.

**Step 2: Run RED**

Run: `go test ./mysql/catalog -run 'TestFunctionalIndexCreateTable|TestScenario_C1/1_11|TestScenario_C1/1_12' -count=1 -timeout=20m`

Expected: fail because functional indexes do not synthesize hidden columns and C1 subtests are skipped.

**Step 3: Implement synthesis helper**

Add shared helpers in `tablecmds.go` or a new `functional_index.go`:

- `isFunctionalIndexColumn(*IndexColumn) bool`
- `allocFunctionalIndexName(tbl *Table) string`
- `makeFunctionalIndexColumnName(tbl *Table, indexName string, keyPart int) string`
- `synthesizeFunctionalIndexColumns(tbl *Table, idx *Index) error`

Rules:

- Unnamed functional indexes use `functional_index`, then suffix `_2`, `_3`, ...
- Hidden column name format is `!hidden!<index>!<part>!<count>`.
- Count starts at `0` and increments on column-name collision.
- Truncate only the `<index>` portion if the final name would exceed 64 characters.
- Hidden columns are virtual generated columns: `Generated={Expr: idxCol.Expr, Stored:false}`, `Hidden=ColumnHiddenSystem`.
- Set `idxCol.Name` to the hidden column name and preserve `idxCol.Expr`.

**Step 4: Unskip targeted C1 tests**

Remove `t.Skip` from C1.11 and C1.12 only after the RED failure is observed.

**Step 5: Run GREEN**

Run: `go test ./mysql/catalog -run 'TestFunctionalIndexCreateTable|TestScenario_C1/1_11|TestScenario_C1/1_12' -count=1 -timeout=20m`

Expected: PASS.

**Step 6: Commit**

```bash
git add mysql/catalog/tablecmds.go mysql/catalog/index.go mysql/catalog/functional_index_test.go mysql/catalog/scenarios_c1_test.go
git commit -m "synthesize mysql functional index hidden columns"
```

### Task 1.3: Share Synthesis with Standalone CREATE INDEX

**Files:**
- Modify: `mysql/catalog/indexcmds.go`
- Test: `mysql/catalog/functional_index_test.go`
- Update: `mysql/catalog/scenarios_c19_test.go`

**Step 1: Write failing test**

Add `TestFunctionalIndexCreateIndexSynthesizesHiddenColumn` for:

```sql
CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(64));
CREATE INDEX idx_lower ON t ((LOWER(name)));
```

Assert one hidden column named `!hidden!idx_lower!0!0`, `IndexColumn.Expr` is non-empty, `IndexColumn.Name` points to the hidden column, and `ShowCreateTable` renders `KEY `idx_lower` ((lower(`name`)))` without rendering the hidden column definition.

**Step 2: Run RED**

Run: `go test ./mysql/catalog -run 'TestFunctionalIndexCreateIndexSynthesizesHiddenColumn|TestScenario_C19/19_1' -count=1 -timeout=20m`

Expected: fail because CREATE INDEX path does not call synthesis and C19 is skipped.

**Step 3: Implement**

Call the shared synthesis helper after building `idx` and before appending it to `tbl.Indexes`. Ensure rollback on synthesis error so failed CREATE INDEX does not leave hidden columns behind.

**Step 4: Unskip C19.1**

Replace the whole-test C19 skip with per-subtest skips for still-unimplemented scenarios, and unskip 19.1.

**Step 5: Run GREEN**

Run: `go test ./mysql/catalog -run 'TestFunctionalIndexCreateIndexSynthesizesHiddenColumn|TestScenario_C19/19_1' -count=1 -timeout=20m`

Expected: PASS.

**Step 6: Commit**

```bash
git add mysql/catalog/indexcmds.go mysql/catalog/scenarios_c19_test.go mysql/catalog/functional_index_test.go
git commit -m "support mysql create index functional key parts"
```

## Phase 2: Expression Type Inference

### Task 2.1: Populate VarExprQ Types from Catalog Columns

**Files:**
- Modify: `mysql/catalog/scope.go`
- Modify: `mysql/catalog/analyze_expr.go`
- Create/Modify: `mysql/catalog/type_resolve.go`
- Test: `mysql/catalog/functional_index_test.go`

**Steps:**

1. Write tests for converting `Column` metadata to `ResolvedType`: `INT`, `BIGINT UNSIGNED`, `VARCHAR(64)` with charset/collation, `JSON`.
2. Run RED.
3. Add `resolvedTypeFromColumn(col *Column) *ResolvedType`.
4. Update `analyzeColumnRef` to set `VarExprQ.Type` and `Collation` using the `analyzerScope` column pointer.
5. Run targeted analyzer and functional-index tests.
6. Commit: `git commit -m "infer mysql column reference expression types"`.

### Task 2.2: Infer Functional Hidden Column Types

**Files:**
- Modify: `mysql/catalog/functional_index.go`
- Modify: `mysql/catalog/function_types.go`
- Modify: `mysql/catalog/analyze_expr.go`
- Test: `mysql/catalog/functional_index_test.go`
- Update: `mysql/catalog/scenarios_c19_test.go`

**Steps:**

1. Write failing tests for C19.2 examples: `(a+b) -> bigint`, `LOWER(name) -> varchar(64)` with inherited charset/collation, `CAST(payload->'$.age' AS UNSIGNED) -> bigint unsigned`.
2. Run RED.
3. Build a one-table analyzer scope from the table and analyze the expression AST before converting to hidden `Column`.
4. Add `columnFromResolvedType(rt *ResolvedType)`.
5. Improve function/operator typing minimally: arithmetic integer promotion to BIGINT, LOWER/UPPER preserve string length/collation, CAST returns explicit target.
6. Unskip C19.2.
7. Run GREEN: `go test ./mysql/catalog -run 'TestFunctionalIndexTypeInference|TestScenario_C19/19_2' -count=1 -timeout=20m`.
8. Commit.

## Phase 3: Functional Expression Validation

### Task 3.1: Reject Bare Column Functional Indexes

**Files:** `mysql/catalog/functional_index.go`, `mysql/catalog/scenarios_c19_test.go`

**Steps:**

1. Write failing test for `CREATE TABLE t (a INT, INDEX ((a)))`.
2. Return MySQL-like error 3756.
3. Run targeted test and commit.

### Task 3.2: Reject Disallowed Functions and Subqueries

**Files:** `mysql/catalog/functional_index.go`, `mysql/catalog/validation.go`

**Steps:**

1. Write failing tests for `RAND()`, `DEFAULT()`, variables, and subqueries inside functional key parts.
2. Add expression walker `validateFunctionalIndexExpr`.
3. Return MySQL-like error 3757.
4. Run targeted test and commit.

### Task 3.3: Reject LOB/JSON Result Types Unless Cast

**Files:** `mysql/catalog/functional_index.go`, `mysql/catalog/function_types.go`

**Steps:**

1. Write failing tests for `payload->'$.name'` and `doc->>'$.name'` without CAST.
2. Classify `JSON`, `TEXT`, `BLOB`, and variants as LOB-for-functional-index.
3. Return MySQL-like error 3754 or 3753 where currently asserted.
4. Unskip C19.4 bad-case subtests.
5. Commit.

## Phase 4: JSON Expression Normalization and SHOW Fidelity

### Task 4.1: Preserve/Normalize JSON Operator Forms

**Files:**
- Modify: `mysql/catalog/deparse.go` or `mysql/catalog/deparse_expr.go`
- Modify: `mysql/catalog/tablecmds.go`
- Test: `mysql/catalog/functional_index_test.go`
- Update: `mysql/catalog/scenarios_c19_test.go`

**Steps:**

1. Write failing test using oracle expectation from C19.5:
   `CAST(doc->>'$.name' AS CHAR(64))` renders as
   `cast(json_unquote(json_extract(`doc`,_utf8mb4'$.name')) as char(64) charset utf8mb4)`.
2. Run RED.
3. Normalize `->>` to `json_unquote(json_extract(...))` in functional-index expression rendering, with `_utf8mb4` introducer for JSON path string literals.
4. Ensure non-functional normal deparse is not changed unless tests require it.
5. Unskip C19.5.
6. Run GREEN and commit.

## Phase 5: Lifecycle Cleanup

### Task 5.1: DROP INDEX Removes Attached Hidden Columns

**Files:** `mysql/catalog/indexcmds.go`, `mysql/catalog/table.go`, `mysql/catalog/scenarios_c19_test.go`

**Steps:**

1. Write failing test: create functional index, drop it, assert both index and hidden columns are gone.
2. Implement `removeFunctionalIndexHiddenColumns(tbl, idx)`.
3. Run targeted tests and commit.

### Task 5.2: Reject Direct DROP COLUMN of System Hidden Columns

**Files:** `mysql/catalog/altercmds.go`, `mysql/catalog/errors.go`

**Steps:**

1. Write failing test for `ALTER TABLE t DROP COLUMN `!hidden!idx_lower!0!0`` returning code 3108 after a functional index exists.
2. Implement hidden-column guard before normal drop.
3. Run targeted tests and commit.

### Task 5.3: RENAME INDEX Cascades Hidden Column Names

**Files:** `mysql/catalog/altercmds.go`, `mysql/catalog/functional_index.go`

**Steps:**

1. Write failing test for `ALTER TABLE t RENAME INDEX idx_lower TO idx_lc`.
2. Rename attached hidden columns using the same naming helper and update `IndexColumn.Name`.
3. Run targeted tests and commit.

## Verification Ladder

Use this ladder after each phase:

```bash
go test ./mysql/catalog -run 'TestFunctionalIndex|TestScenario_C1/1_11|TestScenario_C1/1_12|TestScenario_C19' -count=1 -timeout=20m
go test ./mysql/catalog -run 'TestShowCreateTable|TestWalkThrough_3_1_ColumnPositionsSequential|TestWalkThrough_11_2' -count=1
go test ./mysql/catalog -run 'TestScenario_C' -count=1 -timeout=30m
```

Full `go test ./mysql/catalog -count=1 -timeout=30m` is desirable before final handoff, but it may enter long Docker-backed suites. If it exceeds practical turn time, record it explicitly rather than treating it as passed.

## Execution Order

1. Phase 1 is mandatory first. It creates the hidden-column architecture and makes the basic C1/C19.1 shape pass.
2. Phase 2 must precede LOB validation because LOB rejection depends on resolved expression types.
3. Phase 3 can be split into independent validation commits once the type resolver exists.
4. Phase 4 should wait until JSON operator parsing/deparse details are isolated by tests.
5. Phase 5 should wait until hidden columns exist, otherwise lifecycle tests pass for the wrong reason.

