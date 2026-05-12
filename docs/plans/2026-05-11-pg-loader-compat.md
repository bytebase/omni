# PG Loader Compatibility Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Close PostgreSQL-supported syntax and analysis gaps that currently fail in omni's PG parser/catalog loader.

**Architecture:** Add focused regression tests in `pg/catalog` or `pg/parser` for each PG-supported construct, verify they fail for the expected reason, then make the narrowest parser/catalog/analyzer change for each root cause. Keep each change scoped to the component that rejects behavior PostgreSQL accepts.

**Tech Stack:** Go, `go test`, omni `pg/parser`, omni `pg/catalog`.

### Task 1: Zero-column table and `MATCH SIMPLE`

**Files:**
- Modify: `pg/catalog/tablecmds.go`
- Modify: `pg/parser/create_table.go`
- Test: `pg/catalog/loader_compat_test.go`

**Step 1: Write the failing tests**

Add tests that call `LoadSQL` for `CREATE TABLE t ();` and a foreign key using `MATCH SIMPLE`.

**Step 2: Run tests to verify they fail**

Run: `go test ./pg/catalog -run 'TestLoaderCompat' -count=1`
Expected: zero-column table fails with `tables must have at least one column`; `MATCH SIMPLE` fails near `SIMPLE`.

**Step 3: Implement minimal fixes**

Remove the catalog rejection for zero-column regular tables. Update `parseKeyMatch` to consume the `SIMPLE` token as well as identifier text.

**Step 4: Run tests to verify they pass**

Run: `go test ./pg/catalog -run 'TestLoaderCompat' -count=1`
Expected: PASS.

### Task 2: Function comments and `RETURNS TABLE`

**Files:**
- Modify: `pg/parser/define.go` or relevant function argument parser
- Modify: `pg/catalog/functioncmds.go`
- Test: `pg/catalog/loader_compat_test.go`

**Step 1: Write failing tests**

Add tests for `COMMENT ON FUNCTION f(arg_name integer)` and `RETURNS TABLE(...) LANGUAGE plpgsql`.

**Step 2: Verify red**

Run targeted `go test` and confirm failures are from named argument parsing or return validation.

**Step 3: Implement minimal fixes**

Accept optional argument names in `ObjectWithArgs` type lists and keep locating comments by argument type OIDs. Align `RETURNS TABLE` return validation with PostgreSQL's OUT parameter behavior.

**Step 4: Verify green**

Run the targeted loader compatibility test.

### Task 3: View analyzer expression gaps

**Files:**
- Modify: `pg/catalog/analyze.go`
- Test: `pg/catalog/loader_compat_test.go`

**Step 1: Write failing tests**

Add view tests for `concat_ws(text, variadic any)`, `jsonb ->> c.col_name` with `unnest(text[])` alias, CTE range resolution, and later LATERAL alias scope.

**Step 2: Verify red**

Run targeted tests and record exact failing analyzer path for each case.

**Step 3: Implement one analyzer fix per failure**

Use PostgreSQL-compatible type inference and range-table scope handling. Keep changes narrow and do not add broad fallback typing unless the test proves the exact missing behavior.

**Step 4: Verify green**

Run targeted tests, then relevant analyzer/view regression tests.

### Task 4: Index partition attach and FK type compatibility

**Files:**
- Modify: `pg/catalog/alter.go`
- Modify: `pg/catalog/constraint.go`
- Test: `pg/catalog/loader_compat_test.go`

**Step 1: Write failing tests**

Add tests for `ALTER INDEX parent ATTACH PARTITION child` and bigint FK referencing integer PK.

**Step 2: Verify red**

Run targeted tests and confirm whether failure is relation lookup, index parent metadata, or type compatibility direction.

**Step 3: Implement minimal fixes**

Fix partitioned parent index lookup/metadata and align FK compatibility with PostgreSQL's accepted implicit coercion direction.

**Step 4: Verify green**

Run targeted tests and partition/FK regression tests.
