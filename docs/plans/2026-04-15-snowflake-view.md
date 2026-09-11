# T2.4 — Snowflake VIEW + MATERIALIZED VIEW — Implementation Plan

## Step 1 — AST nodes (snowflake/ast/parsenodes.go + nodetags.go)

1. Add ViewColumn helper struct (not a Node).
2. Add RowAccessPolicy helper struct (not a Node).
3. Add CreateViewStmt + Tag + compile-time assertion.
4. Add CreateMaterializedViewStmt + Tag + compile-time assertion.
5. Add AlterViewAction enum + AlterViewStmt + Tag.
6. Add AlterMaterializedViewAction enum + AlterMaterializedViewStmt + Tag.
7. Add NodeTag constants T_CreateViewStmt, T_CreateMaterializedViewStmt,
   T_AlterViewStmt, T_AlterMaterializedViewStmt in nodetags.go.
8. Add String() cases for the new tags.

## Step 2 — Walker (snowflake/ast/walk_generated.go)

Add cases for the four new statement types in walkChildren.

## Step 3 — Parser: CREATE VIEW + CREATE MATERIALIZED VIEW (snowflake/parser/view.go — new file)

Functions:
- parseCreateViewStmt(start, orReplace, secure, recursive)
- parseCreateMaterializedViewStmt(start, orReplace, secure)
- parseViewColumnList() — parses ( col [COMMENT 'x'], ... )
- parseViewCols() — parses view_col* entries
- parseRowAccessPolicy() — parses [WITH] ROW ACCESS POLICY name ON (cols)

## Step 4 — Parser: CREATE dispatch (snowflake/parser/create_table.go)

In parseCreateStmt, after consuming OR REPLACE and temporary modifiers,
add SECURE/RECURSIVE peek before the switch and add VIEW/MATERIALIZED cases.

## Step 5 — Parser: ALTER VIEW + ALTER MATERIALIZED VIEW (snowflake/parser/view.go)

Functions:
- parseAlterViewStmt()
- parseAlterMaterializedViewStmt()

## Step 6 — Parser: ALTER dispatch (snowflake/parser/database_schema.go)

In parseAlterStmt, add kwVIEW and kwMATERIALIZED cases.

## Step 7 — Tests (snowflake/parser/view_test.go — new file)

Cover all forms documented in the spec.

## Step 8 — Run tests + gofmt

```
cd /Users/h3n4l/OpenSource/omni/.worktrees/snowflake-view
go test ./snowflake/... -count=1
gofmt -w snowflake/
```

## Step 9 — Commit in chunks

Chunk 1: AST nodes + nodetags + walker
Chunk 2: Parser CREATE VIEW + CREATE MATERIALIZED VIEW + dispatch
Chunk 3: Parser ALTER VIEW + ALTER MATERIALIZED VIEW + dispatch  
Chunk 4: Tests + docs

## Step 10 — Open PR
