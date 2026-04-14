# Plan: T2.1 — DATABASE + SCHEMA DDL (Snowflake)

**Date:** 2026-04-14
**Spec:** `docs/superpowers/specs/2026-04-14-snowflake-db-schema-design.md`

## Steps

### Step 1 — AST types + NodeTags (commit 1)

1. Edit `snowflake/ast/parsenodes.go`:
   - Add `AlterDatabaseAction` enum + constants
   - Add `AlterSchemaAction` enum + constants
   - Add `DBSchemaProps` helper struct
   - Add `CreateDatabaseStmt`, `CreateSchemaStmt`
   - Add `AlterDatabaseStmt`, `AlterSchemaStmt`
   - Add `DropDatabaseStmt`, `DropSchemaStmt`
   - Add `UndropDatabaseStmt`, `UndropSchemaStmt`
   - Add compile-time var _ Node assertions

2. Edit `snowflake/ast/nodetags.go`:
   - Add 8 T_* constants to the iota block
   - Add 8 String() cases

3. Run `go run ./snowflake/ast/cmd/genwalker` to regenerate walker

4. Run `go test ./snowflake/...` — must pass (no parser tests yet)

### Step 2 — Parser implementation (commit 2)

1. Edit `snowflake/parser/create_table.go`:
   - Add `kwDATABASE` and `kwSCHEMA` cases to `parseCreateStmt`'s switch

2. Edit `snowflake/parser/parser.go`:
   - Replace `unsupported("ALTER")` with `p.parseAlterStmt()`
   - Replace `unsupported("DROP")` with `p.parseDropStmt()`
   - Replace `unsupported("UNDROP")` with `p.parseUndropStmt()`

3. Create `snowflake/parser/database_schema.go` with:
   - `parseCreateDatabaseStmt`
   - `parseCreateSchemaStmt`
   - `parseAlterStmt` dispatcher
   - `parseAlterDatabaseStmt`
   - `parseAlterSchemaStmt`
   - `parseDropStmt` dispatcher
   - `parseDropDatabaseStmt`
   - `parseDropSchemaStmt`
   - `parseUndropStmt` dispatcher
   - `parseUndropDatabaseStmt`
   - `parseUndropSchemaStmt`
   - `parseDBSchemaProps` helper
   - `parsePropertyNameList` helper
   - `parseUnsetTagList` helper

4. Run `go build ./snowflake/...` — must pass

### Step 3 — Tests (commit 3)

1. Create `snowflake/parser/database_schema_test.go` with ~45 test cases
2. Run `go test ./snowflake/... -count=1` — must pass
3. Run `gofmt -w snowflake/`

### Step 4 — Docs + commit (commit 4 or merge with commit 1)

- Spec and plan already written before code changes
- Combine docs commit with AST commit if clean

## Commit messages

```
feat(snowflake): T2.1 step 1 — AST nodes for DATABASE/SCHEMA DDL
feat(snowflake): T2.1 step 2 — parser for DATABASE/SCHEMA DDL
feat(snowflake): T2.1 step 3 — tests for DATABASE/SCHEMA DDL
feat(snowflake): T2.1 step 4 — spec+plan docs for DATABASE/SCHEMA DDL
```
