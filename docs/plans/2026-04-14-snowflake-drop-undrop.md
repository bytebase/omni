# Plan: T2.5 Snowflake DROP / UNDROP DDL

## Steps

### Step 1 — AST types + node tags

Files: `snowflake/ast/parsenodes.go`, `snowflake/ast/nodetags.go`

- Add `DropObjectKind` enum with 16 values
- Add `DropStmt` struct (Kind, IfExists, Name, Cascade, Restrict, Loc)
- Add `UndropObjectKind` enum with 3 values
- Add `UndropStmt` struct (Kind, Name, Loc)
- Add `T_DropStmt` and `T_UndropStmt` to nodetags.go
- Add compile-time assertions `var _ Node = (*DropStmt)(nil)` etc.

### Step 2 — Regenerate walker

Run `go run ./snowflake/ast/cmd/genwalker` from the worktree root.
This auto-generates cases for `*DropStmt` and `*UndropStmt` in
`walk_generated.go`.

### Step 3 — Parser: drop.go + undrop dispatch

File: `snowflake/parser/drop.go` (new file)

Implement:
- `parseDropStmt()` — dispatch on object type keyword
- `parseDropObject(kind, name *ast.ObjectName, ifExists, cascadeOK bool)` — shared builder
- `parseUndropStmt()` — dispatch on object type keyword

Update `snowflake/parser/parser.go`:
- Replace `p.unsupported("DROP")` → `p.parseDropStmt()`
- Replace `p.unsupported("UNDROP")` → `p.parseUndropStmt()`

### Step 4 — Tests

File: `snowflake/parser/drop_test.go` (new file)

Cover all forms listed in the spec.

### Step 5 — Verify, format, commit

```
cd /Users/h3n4l/OpenSource/omni/.worktrees/snowflake-drop
go test ./snowflake/... -count=1
gofmt -w snowflake/
```

Commit in 3 logical chunks:
1. AST types + node tags
2. Parser implementation + parser.go dispatch wiring
3. Tests + docs
