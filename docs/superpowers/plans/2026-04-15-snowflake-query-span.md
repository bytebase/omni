# Snowflake Query Span Extractor — Implementation Plan (T3.1)

**Date**: 2026-04-15
**Spec**: `docs/superpowers/specs/2026-04-15-snowflake-query-span-design.md`
**Branch**: `feat/snowflake/query-span`
**Worktree**: `.worktrees/snowflake-query-span`

---

## Step 1 — Define types + public API skeleton

File: `snowflake/analysis/query_span.go`

- `QuerySpan`, `ResultColumn`, `SourceColumn` structs
- `Extract(ast.Node) (*QuerySpan, error)` — dispatch on node type
- `ExtractSQL(sql string) (*QuerySpan, error)` — parse + extract

## Step 2 — Implement queryScope + helpers

- `queryScope` struct with `ctes map[string]*QuerySpan`
- `newScope() *queryScope`
- `childScope(parent *queryScope) *queryScope`
- `deduplicateSources(srcs []*SourceColumn) []*SourceColumn`
- `resultColumnName(target *ast.SelectTarget) string`

## Step 3 — Implement FROM scope builder

- `tableEntry` struct: `{alias string, columns []*SourceColumn}`
- `buildFromScope(froms []ast.Node, scope *queryScope) []tableEntry`
- `resolveTableRef(ref *ast.TableRef, scope *queryScope) tableEntry`
- `resolveJoin(join *ast.JoinExpr, scope *queryScope) []tableEntry`

## Step 4 — Implement expression column ref collector

- `collectRefs(expr ast.Node, fromEntries []tableEntry) []*SourceColumn`
- `isDerivedExpr(expr ast.Node) bool`

## Step 5 — Implement extractSelectStmt

- Handles WITH clause → builds child scope with CTE spans
- Iterates over `Targets`
- For Star targets: emits `*` pseudo-column with sources from FROM tables
- For expression targets: calls collectRefs, builds ResultColumn

## Step 6 — Implement extractSetOperationStmt

- Positional merge (default)
- Named merge (ByName)
- Sources = union

## Step 7 — Write tests (query_span_test.go)

15+ tests covering all required cases:
1. Single column from table
2. Two columns from table
3. Aliased column
4. COUNT(*) — derived
5. SELECT * from table
6. SELECT * EXCLUDE(b)
7. JOIN two tables
8. Subquery in FROM
9. CTE — source resolution
10. UNION ALL — positional merge
11. UNION ALL BY NAME
12. Qualified column ref db.schema.t.col
13. CTE with column aliases
14. SELECT constant — derived, no sources
15. Multi-table JOIN with qualified refs

## Step 8 — Run tests, fix failures

```
cd .worktrees/snowflake-query-span && go test ./snowflake/... -count=1
```

## Step 9 — Format + commit

```
gofmt -w snowflake/analysis/
git add snowflake/analysis/ docs/superpowers/specs/... docs/superpowers/plans/...
git commit ...
```

## Step 10 — Push + open PR

```
git push -u origin feat/snowflake/query-span
gh pr create --title "feat(snowflake): T3.1 query span extractor" ...
```
