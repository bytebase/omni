# T1.8 — Snowflake Statement-Classification Helper — Implementation Plan

## Step 1 — Create snowflake/analysis/classify.go

Implement:
1. `Category` type and constants (Unknown, Select, DML, DDL, Show, Describe, Other).
2. `(Category) String()` method.
3. `Classify(node ast.Node) Category` — type switch over all current statement nodes.
4. `ClassifySQL(sql string) Category` — parse first statement, delegate to Classify.

Imports: `github.com/bytebase/omni/snowflake/ast`, `github.com/bytebase/omni/snowflake/parser`.

## Step 2 — Create snowflake/analysis/classify_test.go

Test cases:
- `SELECT 1` → CategorySelect
- `SELECT 1 UNION ALL SELECT 2` → CategorySelect
- `WITH cte AS (SELECT 1) SELECT * FROM cte` → CategorySelect
- `INSERT INTO t VALUES (1)` → CategoryDML
- `UPDATE t SET c = 1 WHERE id = 1` → CategoryDML
- `DELETE FROM t WHERE id = 1` → CategoryDML
- `MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE` → CategoryDML
- `CREATE TABLE t (id INT)` → CategoryDDL
- `ALTER TABLE t ADD COLUMN c INT` → CategoryDDL
- `DROP TABLE t` → CategoryDDL
- `CREATE DATABASE d` → CategoryDDL
- `CREATE VIEW v AS SELECT 1` → CategoryDDL
- `nil` node → CategoryUnknown
- empty string → CategoryUnknown
- `Classify(nil)` → CategoryUnknown

## Step 3 — Run tests

```
cd /Users/h3n4l/OpenSource/omni/.worktrees/snowflake-stmt-classify
go test ./snowflake/... -count=1
```

## Step 4 — gofmt

```
gofmt -w /Users/h3n4l/OpenSource/omni/.worktrees/snowflake-stmt-classify/snowflake/analysis/
```

## Step 5 — Commit in 2 chunks

Chunk 1 (T1.8 step 1): docs + classify.go (impl)
Chunk 2 (T1.8 step 2): classify_test.go (tests)

## Step 6 — Push + open PR
