# Snowflake ALTER TABLE Implementation Plan (T2.3)

## Steps

### Step 1 — AST types (parsenodes.go + nodetags.go)
- Add `AlterTableStmt`, `AlterTableAction`, `AlterTableActionKind`, `ColumnAlter`, `ColumnAlterKind`, `TableProp` to `parsenodes.go`
- Add `T_AlterTableStmt` to `nodetags.go`
- Run `go build ./snowflake/ast/...` to verify

### Step 2 — Regenerate walker
- Run `cd snowflake/ast && go run ./cmd/genwalker`
- Verify `walk_generated.go` has case for `*AlterTableStmt`

### Step 3 — Parser: alter_table.go
- Create `snowflake/parser/alter_table.go`
- Implement `parseAlterTableStmt()`
- Implement all action sub-parsers
- Update `parseAlterStmt()` in `database_schema.go` to dispatch `kwTABLE`

### Step 4 — Tests: alter_table_test.go
- Lint-critical paths: ADD COLUMN, DROP COLUMN, RENAME TO, ADD CONSTRAINT, DROP PRIMARY KEY
- Other paths: SWAP WITH, RENAME COLUMN, ALTER COLUMN variants, CLUSTER BY, SET, UNSET, TAG, ROW ACCESS POLICY, SEARCH OPTIMIZATION

### Step 5 — Verify and format
- `go test ./snowflake/... -count=1`
- `gofmt -w snowflake/`

### Step 6 — Commits (3-5 logical)
1. `feat(snowflake): AlterTableStmt AST types + NodeTag (T2.3 step 1)`
2. `feat(snowflake): regen walker for AlterTableStmt (T2.3 step 2)`
3. `feat(snowflake): ALTER TABLE parser — all major actions (T2.3 step 3)`
4. `feat(snowflake): ALTER TABLE tests — lint-critical + coverage (T2.3 step 4)`

### Step 7 — Push + PR
- `git push -u origin feat/snowflake/alter-table`
- `gh pr create ...`
