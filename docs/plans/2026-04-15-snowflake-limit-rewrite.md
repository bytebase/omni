# T3.3 — Snowflake LIMIT Injection Rewrite — Implementation Plan

## Steps

### Step 1 — Write rewrite.go

File: `snowflake/deparse/rewrite.go`

Functions:
1. `InjectLimit(sql string, maxRows int) (string, error)` — public entry point
2. `rewriteStmt(node ast.Node, maxRows int) ast.Node` — dispatch on node type
3. `rewriteSelect(s *ast.SelectStmt, maxRows int) ast.Node` — handle SelectStmt
4. `wrapInLimit(inner ast.Node, maxRows int) *ast.SelectStmt` — build wrapper SELECT

### Step 2 — Write rewrite_test.go

File: `snowflake/deparse/rewrite_test.go`

Test cases (12):
1. No LIMIT → add LIMIT
2. Existing LIMIT n <= maxRows → unchanged
3. Existing LIMIT n > maxRows → lowered
4. LIMIT with bind variable → wrap
5. FETCH FIRST n > maxRows → lowered in-place
6. UNION → wrap
7. WITH CTE + SELECT → add LIMIT at outer SELECT
8. INSERT → unchanged
9. CREATE TABLE → unchanged
10. Multi-statement (SELECT;INSERT;SELECT) → apply only to SELECTs
11. Invalid SQL → error returned
12. maxRows <= 0 → error returned

### Step 3 — Run tests

```
go test ./snowflake/... -count=1
```

### Step 4 — gofmt

```
gofmt -w snowflake/deparse/rewrite.go snowflake/deparse/rewrite_test.go
```

### Step 5 — Commit

Commit 1: spec + plan docs
Commit 2: rewrite.go implementation
Commit 3: rewrite_test.go tests

### Step 6 — Push + PR

```
git push -u origin feat/snowflake/limit-rewrite
gh pr create ...
```
